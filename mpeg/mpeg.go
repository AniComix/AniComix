package mpeg

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// ffprobeShowStreamEntriesResult represents the overall structure of the JSON output
type ffprobeShowStreamEntriesResult struct {
	Streams []AVStreamInfo `json:"streams"`
}

func GetFullMediaStreamInfo(path string) []AVStreamInfo {
	info, err := ffprobe().
		setLogLevel("error").
		setInput(path).
		setShowEntry("stream").
		setOf("json").combinedOutput()
	if err != nil {
		log.Printf("Could not get media info: %s,%v", string(info), err)
		return nil
	}
	var result ffprobeShowStreamEntriesResult
	err = json.Unmarshal(info, &result)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
		return nil
	}
	return result.Streams
}

// width of the output video is calculated automatically
func transformVideoResolutionGeneric(inputPath, outputPath string, height int32) bool {
	err := ffmpeg().
		useHardwareAcceleration("cuda").
		setInput(inputPath).
		setLogLevel("error").
		setVideoFilter(fmt.Sprintf("scale=-2:%d", height)).
		setAudioCodec("copy").
		setOutput(outputPath).
		run()
	if err != nil {
		log.Printf("Could not transform video: %v", err)
		return false
	}
	return true
}

func TransformVideoResolution2160p(inputPath, outputPath string) {
	transformVideoResolutionGeneric(inputPath, outputPath, 2160)
}
func TransformVideoResolution1440p(inputPath, outputPath string) {
	transformVideoResolutionGeneric(inputPath, outputPath, 1440)
}
func TransformVideoResolution1080p(inputPath, outputPath string) {
	transformVideoResolutionGeneric(inputPath, outputPath, 1080)
}
func TransformVideoResolution720p(inputPath, outputPath string) {
	transformVideoResolutionGeneric(inputPath, outputPath, 720)
}
func TransformVideoResolution360p(inputPath, outputPath string) {
	transformVideoResolutionGeneric(inputPath, outputPath, 360)
}
func TransformVideoResolution144p(inputPath, outputPath string) {
	transformVideoResolutionGeneric(inputPath, outputPath, 144)
}

func adjustBitrate(source, target string) string {
	targetRate, err := strconv.ParseInt(strings.TrimRight(target, "k"), 10, 32)
	if err != nil {
		log.Printf("Error parsing bitrate: %v", err)
		return source
	}
	targetRate *= 1000
	sourceRate, _ := strconv.ParseInt(source, 10, 32)
	log.Printf("source: %d, target: %d", sourceRate, targetRate)
	if sourceRate >= targetRate {
		return target
	}
	return source
}

func TransformVideoToDASHMultipleResolution(inputPath, mpdPath string) bool {
	streams := GetFullMediaStreamInfo(inputPath)
	if len(streams) == 0 {
		log.Printf("Could not find a video stream")
		return false
	}
	height := streams[0].Height
	log.Println(streams[0])
	cmd := ffmpeg().
		useHardwareAcceleration("cuda").
		setBenchmark().
		setInput(inputPath)
	source := streams[0].Bitrate
	if height >= 360 {
		cmd.setMap("0:v", "-vf", "scale=-2:360", "-b:v:0", adjustBitrate(source, "800k"), "-c:v:0", LibX264)
	}
	if height >= 720 {
		cmd.setMap("0:v", "-vf", "scale=-2:720", "-b:v:1", adjustBitrate(source, "1500k"), "-c:v:1", LibX264)
	}
	if height >= 1080 {
		cmd.setMap("0:v", "-vf", "scale=-2:1080", "-b:v:2", adjustBitrate(source, "3000k"), "-c:v:2", LibX264)
	}
	if height >= 1440 {
		cmd.setMap("0:v", "-vf", "scale=-2:1440", "-b:v:3", adjustBitrate(source, "4000k"), "-c:v:3", LibX264)
	}
	cmd.setMap("0:a").setAudioCodec(AAC).setAudioBitrate("128k")
	err := cmd.arg("-f", "dash", mpdPath).run()
	if err != nil {
		log.Printf("Could not transform video: %v", err)
		return false
	}
	return true
}

func CheckVideoFileIntegrity(path string) bool {
	streams := GetFullMediaStreamInfo(path)
	return len(streams) > 0 && streams[0].CodecType == "video"
}

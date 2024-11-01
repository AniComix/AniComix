package mpeg

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// ffprobeShowStreamEntriesResult represents the overall structure of the JSON output
type ffprobeShowStreamEntriesResult struct {
	Streams []AVStreamInfo `json:"streams"`
}

func GetFullMediaStreamInfo(path string) []AVStreamInfo {

	cmd := ffprobe().
		setLogLevel("error").
		setInput(path).
		setShowEntry("stream").
		setOf("json")
	info, err := exec.Command(cmd.args[0], cmd.args[1:]...).CombinedOutput()
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
	result, err := ffmpeg().
		useHardwareAcceleration("cuda").
		setInput(inputPath).
		setLogLevel("error").
		setVideoFilter(fmt.Sprintf("scale=-2:%d", height)).
		setAudioCodec("copy").
		setOutput(outputPath).
		run()
	if err != nil {
		log.Printf("Could not transform video: %s,%v", result, err)
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

func TransformVideoToDASHMultipleResolution(inputPath, mpdPath string) {
	info, err := ffmpeg().
		useHardwareAcceleration("cuda").
		setInput(inputPath).
		setMap("0:v", "-s:v:0", "1920x1080", "-b:v:0", "3000k", "-c:v:0", LibX264).
		setMap("0:v", "-s:v:1", "1280x720", "-b:v:1", "1500k", "-c:v:1", LibX264).
		setMap("0:v", "-s:v:2", "854x480", "-b:v:2", "800k", "-c:v:2", LibX264).
		setMap("0:a").setAudioCodec(AAC).setAudioBitrate("128k").
		arg("-f", "dash", mpdPath).run()
	if err != nil {
		log.Printf("Could not transform video: %s,%v", info, err)
		return
	}
	log.Printf("Transformed video to: %s", info)
}

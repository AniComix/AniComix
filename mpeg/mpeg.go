package mpeg

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

const (
	ffprobe = "ffprobe"
	ffmpeg  = "ffmpeg"
	ffplay  = "ffplay"
	// help / information / capabilities:
	optHardwareAcceleration = "-hwaccel"
	optShowEntries          = "-show_entries"
	optOf                   = "-of"
	optSegDuration          = "-seg_duration"
	optMap                  = "-map"
	optInput                = "-i"
	optVersion              = "-version"
	optMuxers               = "-muxers"
	optDemuxers             = "-demuxers"
	optDevices              = "-devices"
	optEncoders             = "-encoders"
	optDecoders             = "-decoders"
	optPixelFormats         = "-pix_fmts"
	optFilters              = "-filters"
	optLayouts              = "-layouts"
	optSampleFormats        = "-sample_fmts"

	optSetLogLevel              = "-v"
	optOverwriteOutputFile      = "-y"
	optNeverOverwriteOutputFile = "-n"
	optPrintEncodingProgress    = "-stats"

	optSetAudioCodec   = "-c:a"
	optSetVideoCodec   = "-c:v"
	optSetAudioBitrate = "-b:a"
	optSetVideoBitrate = "-b:v"
	optSetVideoFilter  = "-vf"
)

func Version() (string, error) {
	cmd := exec.Command(ffmpeg, optVersion)
	version, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(version), nil
}

// ffprobeShowStreamEntriesResult represents the overall structure of the JSON output
type ffprobeShowStreamEntriesResult struct {
	Streams []AVStreamInfo `json:"streams"`
}

func GetFullMediaStreamInfo(path string) []AVStreamInfo {
	cmd := exec.Command(
		ffprobe, optSetLogLevel, "error",
		optShowEntries, "stream",
		optOf, "json",
		path)
	info, err := cmd.CombinedOutput()
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

func transformVideoResolutionGeneric(path string, height int32) bool {
	outputName := fmt.Sprintf("./output_%dp.mp4", height)
	cmd := exec.Command(
		ffmpeg,
		optHardwareAcceleration, "cuda",
		optInput, path,
		optSetLogLevel, "error",
		optSetVideoFilter, fmt.Sprintf("scale=-2:%d", height),
		optSetAudioCodec, "copy",
		optOverwriteOutputFile,
		outputName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Could not transform video: %s,%v", string(result), err)
		return false
	}
	return true
}

func TransformVideoResolution2160p(path string) {
	transformVideoResolutionGeneric(path, 2160)
}
func TransformVideoResolution1440p(path string) {
	transformVideoResolutionGeneric(path, 1440)
}
func TransformVideoResolution1080p(path string) {
	transformVideoResolutionGeneric(path, 1080)
}
func TransformVideoResolution720p(path string) {
	transformVideoResolutionGeneric(path, 720)
}
func TransformVideoResolution360p(path string) {
	transformVideoResolutionGeneric(path, 360)
}
func TransformVideoResolution144p(path string) {
	transformVideoResolutionGeneric(path, 144)
}

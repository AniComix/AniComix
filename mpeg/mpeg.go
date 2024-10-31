package mpeg

import (
	"log"
	"os/exec"
)

const (
	ffprobe = "ffprobe"
	ffmpeg  = "ffmpeg"
	ffplay  = "ffplay"
	// help / information / capabilities:
	optShowEntries   = "-show_entries"
	optOf            = "-of"
	optSegDuration   = "-seg_duration"
	optMap           = "-map"
	optInput         = "-i"
	optVersion       = "-version"
	optMuxers        = "-muxers"
	optDemuxers      = "-demuxers"
	optDevices       = "-devices"
	optEncoders      = "-encoders"
	optDecoders      = "-decoders"
	optPixelFormats  = "-pix_fmts"
	optFilters       = "-filters"
	optLayouts       = "-layouts"
	optSampleFormats = "-sample_fmts"

	optSetLogLevel              = "-v"
	optOverwriteOutputFile      = "-y"
	optNeverOverwriteOutputFile = "-n"
	optPrintEncodingProgress    = "-stats"
)

func Version() (string, error) {
	cmd := exec.Command(ffmpeg, optVersion)
	version, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(version), nil
}

func GetFullMediaInfo(path string) {
	cmd := exec.Command(ffprobe, optSetLogLevel, "error", optShowEntries, "stream", optOf, "json", path)
	info, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Could not get media info: %v", err)
	}
	log.Printf("Media info: %s", string(info))
}

func TransformMpeg4IntoDash(path string) {
	cmd := exec.Command(ffmpeg, optInput, path, optMap, "0", optSegDuration, "4", "-f", "dash", "./data/output.mpd")
	info, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Could not transform media: %v", err)
	}
	log.Printf("%s", string(info))
}

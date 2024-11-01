package mpeg

import (
	"log"
	"os/exec"
)

type FFCommand struct {
	args []string
}

// ffmpeg command line options
const (
	optHardwareAcceleration = "-hwaccel"
	optSegDuration          = "-seg_duration"
	optMap                  = "-map"
	optInput                = "-i"
	optVersion              = "-version"

	optSetLogLevel              = "-v"
	optOverwriteOutputFile      = "-y"
	optNeverOverwriteOutputFile = "-n"

	optSetAudioCodec   = "-c:a"
	optSetVideoCodec   = "-c:v"
	optSetAudioBitrate = "-b:a"
	optSetVideoBitrate = "-b:v"
	optSetVideoFilter  = "-vf"
)

const (
	LibX264 = "libx264"
	AAC     = "aac"
)

func ffmpeg() *FFCommand {
	return &FFCommand{args: []string{"ffmpeg"}}
}
func (c *FFCommand) useHardwareAcceleration(accelerateDevice string) *FFCommand {
	c.args = append(c.args, optHardwareAcceleration, accelerateDevice)
	return c
}
func (c *FFCommand) setVideoCodec(videoCodec string) *FFCommand {
	c.args = append(c.args, optSetVideoCodec, videoCodec)
	return c
}
func (c *FFCommand) setAudioCodec(audioCodec string) *FFCommand {
	c.args = append(c.args, optSetAudioCodec, audioCodec)
	return c
}
func (c *FFCommand) setAudioBitrate(audioBitrate string) *FFCommand {
	c.args = append(c.args, optSetAudioBitrate, audioBitrate)
	return c
}
func (c *FFCommand) setVideoBitrate(videoBitrate string) *FFCommand {
	c.args = append(c.args, optSetVideoBitrate, videoBitrate)
	return c
}
func (c *FFCommand) setVideoFilter(filters ...string) *FFCommand {
	c.args = append(c.args, optSetVideoFilter)
	c.args = append(c.args, filters...)
	return c
}
func (c *FFCommand) setLogLevel(level string) *FFCommand {
	c.args = append(c.args, optSetLogLevel, level)
	return c
}
func (c *FFCommand) setInput(inputPath string) *FFCommand {
	c.args = append(c.args, optInput, inputPath)
	return c
}
func (c *FFCommand) setOutput(outputPath string) *FFCommand {
	c.args = append(c.args, outputPath)
	return c
}
func (c *FFCommand) setOutputOverwrite(overwrite bool) *FFCommand {
	if overwrite {
		c.args = append(c.args, optOverwriteOutputFile)
	} else {
		c.args = append(c.args, optNeverOverwriteOutputFile)
	}
	return c
}

func (c *FFCommand) setSegDuration(segDuration string) *FFCommand {
	c.args = append(c.args, optSegDuration, segDuration)
	return c
}

func (c *FFCommand) setMap(args ...string) *FFCommand {
	c.args = append(c.args, optMap)
	c.args = append(c.args, args...)
	return c
}
func (c *FFCommand) arg(args ...string) *FFCommand {
	c.args = append(c.args, args...)
	return c
}

func (c *FFCommand) run() (string, error) {
	cmd := exec.Command(c.args[0], c.args[1:]...)
	log.Printf("Running command: %s\n", cmd.String())
	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Could not run command: %s\n", cmd.String())
	}
	return string(result), err
}

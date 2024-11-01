package mpeg

import (
	"fmt"
	"strconv"
	"strings"
)

type AVRational struct {
	Numerator   int32
	Denominator int32
}

// UnmarshalJSON is a custom unmarshal method for AVRational
func (r *AVRational) UnmarshalJSON(data []byte) error {
	// Remove the surrounding quotes from the JSON string
	str := strings.Trim(string(data), "\"")
	// Split the string on "/"
	parts := strings.Split(str, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid AVRational format: %s", str)
	}

	// Parse the numerator and denominator as integers
	numerator, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid numerator: %v", err)
	}
	denominator, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid denominator: %v", err)
	}

	// Set the values in the AVRational struct
	r.Numerator = int32(numerator)
	r.Denominator = int32(denominator)
	return nil
}

type AVStreamInfo struct {
	Index         int32
	CodecName     string     `json:"codec_name"`
	CodecLongName string     `json:"codec_long_name"`
	CodecType     string     `json:"codec_type"`
	Duration      string     `json:"duration"`
	SampleRate    string     `json:"sample_rate"`
	SampleFormat  string     `json:"sample_fmt"`
	Width         int32      `json:"width"`
	Height        int32      `json:"height"`
	Bitrate       string     `json:"bit_rate"`
	AverageFPS    AVRational `json:"avg_frame_rate"`
	TimeBase      AVRational `json:"time_base"`
}

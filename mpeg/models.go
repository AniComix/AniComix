package mpeg

type AVStream struct {
	CodecName      string
	CodecLongName  string
	Profile        string
	CodecType      string
	CodecTagString string
	CodecTag       string
	Width          int32
	Height         int32
	CodecWidth     int32
	CodecHeight    int32
	Bitrate        int32
	AverageFPS     string
}

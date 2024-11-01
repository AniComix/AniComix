package main

import (
	"AniComix/mpeg"
	"fmt"
)

func main() {
	filename := "D:\\Anime\\间谍过家家\\[DMG&LoliHouse] Spy x Family - 29 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv"
	streams := mpeg.GetFullMediaStreamInfo(filename)
	for _, stream := range streams {
		fmt.Printf("Codec Name: %s\n", stream.CodecLongName)
	}
	mpeg.TransformVideoInto720p(filename)
}

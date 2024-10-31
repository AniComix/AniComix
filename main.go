package main

import (
	"fmt"
	"github.com/AniComix/server"
)

func main() {
	fmt.Println("hello")
	mpeg.Transform_MP4_to_DASH("a.mp4")
	server.Run()
}

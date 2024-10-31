package main

import (
	"AniComix/mpeg"
	"fmt"
)

func main() {
	fmt.Println("hello go")
	mpeg.Hello()

	mpeg.Transform_MP4_to_DASH("../abc.mp4")
}

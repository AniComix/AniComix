package mpeg

import (
	"fmt"
	"os/exec"
)

func Hello() {
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}

package mpeg

/*
#include "mpeg.h"
#cgo CFLAGS: -I ./ffmpeg/include
#cgo LDFLAGS: -L ./ffmpeg/lib -lavcodec -lavformat -lavutil
*/
import "C"

func Transform_MP4_to_DASH(name *C.char) {
	C.transform_MP4_to_DASH(name)
}

func Hello() {
	C.hello()
}

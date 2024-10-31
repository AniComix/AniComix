package mpeg

/*
#cgo CFLAGS: -I /opt/homebrew/Cellar/ffmpeg/7.1_2/include
#cgo LDFLAGS: -L /opt/homebrew/Cellar/ffmpeg/7.1_2/lib -lavcodec -lavformat -lavutil
#include "mpeg.h"
*/
import "C"

func Transform_MP4_to_DASH(name string) {
	cstr := C.CString(name)
	C.transform_MP4_to_DASH(cstr)
}

func Hello() {
	C.hello()
}

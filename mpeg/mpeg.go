package c

/*
#include "mpeg.h"
#cgo LDFLAGS: -lavcodec -lavformat -lswscale -lavutil -lavfilter -lm
*/
import "C"

func Transform_MP4_to_DASH(name *C.char) {
	C.transform_MP4_to_DASH(name)
}

//
// Created by delta on 30 Oct 2024.
//

#ifndef MPEG__MPEG_H_
#define MPEG__MPEG_H_

#include "FFmpeg/libavcodec/avcodec.h"
#include "FFmpeg/libavformat/avformat.h"
#include "FFmpeg/libavutil/avutil.h"
#include "FFmpeg/libavutil/dict.h"

int transform_MP4_to_DASH(const char *);

#endif // MPEG__MPEG_H_

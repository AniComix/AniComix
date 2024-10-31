
#ifndef _ANICOMIX_MPEG_H
#define _ANICOMIX_MPEG_H
#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libavutil/avutil.h>
#include <libavutil/dict.h>
#include <stdio.h>
int transform_MP4_to_DASH(const char *);
void hello();
#endif
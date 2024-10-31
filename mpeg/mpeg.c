//
// Created by delta on 30 Oct 2024.
//

#include "mpeg.h"

int transform_MP4_to_DASH(const char *input_path) {
  AVFormatContext *input_context = NULL;
  if (avformat_open_input(&input_context, input_path, NULL, NULL) < 0) {
    fprintf(stderr, "Could not open input file: %s\n", input_path);
    return -1;
  }
  if (avformat_find_stream_info(input_context, NULL) < 0) {
    fprintf(stderr, "Could not find stream information.\n");
    avformat_close_input(&input_context);
    return -1;
  }
  AVFormatContext *output_context = NULL;
  avformat_alloc_output_context2(&output_context, NULL, "dash", NULL);
  if (!output_context) {
    fprintf(stderr, "Could not create output context.\n");
    avformat_close_input(&input_context);
    return -1;
  }

  for (int i = 0; i < input_context->nb_streams; i++) {
    AVStream *in_stream = input_context->streams[i];
    AVStream *out_stream = avformat_new_stream(output_context, NULL);
    if (!out_stream) {
      fprintf(stderr, "Failed allocating output stream\n");
      avformat_close_input(&input_context);
      avformat_free_context(output_context);
      return -1;
    }
    avcodec_parameters_copy(out_stream->codecpar, in_stream->codecpar);
    out_stream->codecpar->codec_tag = 0;
  }

  // Set DASH muxer options
  av_dict_set(&output_context->metadata, "media_segmentation_time", "2", 0);
  av_dict_set(&output_context->metadata, "use_template", "1", 0);
  av_dict_set(&output_context->metadata, "use_timeline", "1", 0);
  av_dict_set(&output_context->metadata, "init_seg_name",
              "init-stream$RepresentationID$.m4s", 0);
  av_dict_set(&output_context->metadata, "media_seg_name",
              "chunk-stream$RepresentationID$-$Number$.m4s", 0);
  av_dict_set(&output_context->metadata, "dash_playlist_name", "manifest.mpd",
              0);

  // Create the output DASH directory
  char dash_output[512];
  snprintf(dash_output, sizeof(dash_output), "%s/manifest.mpd", "output");

  if (!(output_context->oformat->flags & AVFMT_NOFILE)) {
    if (avio_open(&output_context->pb, dash_output, AVIO_FLAG_WRITE) < 0) {
      fprintf(stderr, "Could not open output file '%s'\n", dash_output);
      avformat_close_input(&input_context);
      avformat_free_context(output_context);
      return -1;
    }
  }

  // Write the stream header
  if (avformat_write_header(output_context, NULL) < 0) {
    fprintf(stderr, "Error occurred when opening output file\n");
    avformat_close_input(&input_context);
    if (!(output_context->oformat->flags & AVFMT_NOFILE))
      avio_closep(&output_context->pb);
    avformat_free_context(output_context);
    return -1;
  }

  AVPacket packet;
  while (av_read_frame(input_context, &packet) >= 0) {
    AVStream *in_stream = input_context->streams[packet.stream_index];
    AVStream *out_stream = output_context->streams[packet.stream_index];

    // Rescale timestamps from input to output stream timebase
    packet.pts = av_rescale_q_rnd(packet.pts, in_stream->time_base,
                                  out_stream->time_base,
                                  AV_ROUND_NEAR_INF | AV_ROUND_PASS_MINMAX);
    packet.dts = av_rescale_q_rnd(packet.dts, in_stream->time_base,
                                  out_stream->time_base,
                                  AV_ROUND_NEAR_INF | AV_ROUND_PASS_MINMAX);
    packet.duration = av_rescale_q(packet.duration, in_stream->time_base,
                                   out_stream->time_base);
    packet.pos = -1;

    // Write packet
    if (av_interleaved_write_frame(output_context, &packet) < 0) {
      fprintf(stderr, "Error muxing packet\n");
      break;
    }
    av_packet_unref(&packet);
  }

  // Write the trailer (end of file)
  av_write_trailer(output_context);

  // Close input and output files
  avformat_close_input(&input_context);
  if (!(output_context->oformat->flags & AVFMT_NOFILE))
    avio_closep(&output_context->pb);
  avformat_free_context(output_context);

  printf("DASH file successfully created at %s\n", dash_output);
  return 0;
}

void hello(){
  printf("hello from c\n");
}
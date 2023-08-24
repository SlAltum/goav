// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// Giorgis (habtom@giorgis.io)

//Package avformat provides some generic global options, which can be set on all the muxers and demuxers.
//In addition each muxer or demuxer may support so-called private options, which are specific for that component.
//Supported formats (muxers and demuxers) provided by the libavformat library
package avformat

//#cgo pkg-config: libavformat libavcodec libavutil libavdevice libavfilter libswresample libswscale
//#include <stdio.h>
//#include <stdlib.h>
//#include <inttypes.h>
//#include <stdint.h>
//#include <string.h>
//#include <libavformat/avformat.h>
//#include <libavformat/avio.h>
import "C"

const (
	AVIO_FLAG_READ       = int(C.AVIO_FLAG_READ)
	AVIO_FLAG_WRITE      = int(C.AVIO_FLAG_WRITE)
	AVIO_FLAG_READ_WRITE = int(C.AVIO_FLAG_READ_WRITE)
	// codec context flags
	AV_CODEC_FLAG_UNALIGNED      = int(C.AV_CODEC_FLAG_UNALIGNED)
	AV_CODEC_FLAG_4MV            = int(C.AV_CODEC_FLAG_4MV)
	AV_CODEC_FLAG_OUTPUT_CORRUPT = int(C.AV_CODEC_FLAG_OUTPUT_CORRUPT)
	AV_CODEC_FLAG_QPEL           = int(C.AV_CODEC_FLAG_QPEL)
	AV_CODEC_FLAG_DROPCHANGED    = int(C.AV_CODEC_FLAG_DROPCHANGED)
	// AV_CODEC_FLAG_RECON_FRAME    = int(C.AV_CODEC_FLAG_RECON_FRAME)
	// AV_CODEC_FLAG_COPY_OPAQUE    = int(C.AV_CODEC_FLAG_COPY_OPAQUE)
	// AV_CODEC_FLAG_FRAME_DURATION = int(C.AV_CODEC_FLAG_FRAME_DURATION)
	AV_CODEC_FLAG_PASS1          = int(C.AV_CODEC_FLAG_PASS1)
	AV_CODEC_FLAG_PASS2          = int(C.AV_CODEC_FLAG_PASS2)
	AV_CODEC_FLAG_LOOP_FILTER    = int(C.AV_CODEC_FLAG_LOOP_FILTER)
	AV_CODEC_FLAG_GRAY           = int(C.AV_CODEC_FLAG_GRAY)
	AV_CODEC_FLAG_PSNR           = int(C.AV_CODEC_FLAG_PSNR)
	AV_CODEC_FLAG_INTERLACED_DCT = int(C.AV_CODEC_FLAG_INTERLACED_DCT)
	AV_CODEC_FLAG_LOW_DELAY      = int(C.AV_CODEC_FLAG_LOW_DELAY)
	AV_CODEC_FLAG_GLOBAL_HEADER  = int(C.AV_CODEC_FLAG_GLOBAL_HEADER)
	AV_CODEC_FLAG_BITEXACT       = int(C.AV_CODEC_FLAG_BITEXACT)
	AV_CODEC_FLAG_AC_PRED        = int(C.AV_CODEC_FLAG_AC_PRED)
	AV_CODEC_FLAG_INTERLACED_ME  = int(C.AV_CODEC_FLAG_INTERLACED_ME)
	AV_CODEC_FLAG_CLOSED_GOP     = int(C.AV_CODEC_FLAG_CLOSED_GOP)
	// codec context flags2
	AV_CODEC_FLAG2_FAST          = int(C.AV_CODEC_FLAG2_FAST)
	AV_CODEC_FLAG2_NO_OUTPUT     = int(C.AV_CODEC_FLAG2_NO_OUTPUT)
	AV_CODEC_FLAG2_LOCAL_HEADER  = int(C.AV_CODEC_FLAG2_LOCAL_HEADER)
	AV_CODEC_FLAG2_CHUNKS        = int(C.AV_CODEC_FLAG2_CHUNKS)
	AV_CODEC_FLAG2_IGNORE_CROP   = int(C.AV_CODEC_FLAG2_IGNORE_CROP)
	AV_CODEC_FLAG2_SHOW_ALL      = int(C.AV_CODEC_FLAG2_SHOW_ALL)
	AV_CODEC_FLAG2_EXPORT_MVS    = int(C.AV_CODEC_FLAG2_EXPORT_MVS)
	AV_CODEC_FLAG2_SKIP_MANUAL   = int(C.AV_CODEC_FLAG2_SKIP_MANUAL)
	AV_CODEC_FLAG2_RO_FLUSH_NOOP = int(C.AV_CODEC_FLAG2_RO_FLUSH_NOOP)
	// AV_CODEC_FLAG2_ICC_PROFILES  = int(C.AV_CODEC_FLAG2_ICC_PROFILES)
)

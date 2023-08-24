// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// Giorgis (habtom@giorgis.io)

// Package avutil is a utility library to aid portable multimedia programming.
// It contains safe portable string functions, random number generators, data structures,
// additional mathematics functions, cryptography and multimedia related functionality.
// Some generic features and utilities provided by the libavutil library
package avutil

//#cgo pkg-config: libavutil
//#include <libavutil/avutil.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

type (
	Options       C.struct_AVOptions
	AvTree        C.struct_AVTree
	Rational      C.struct_AVRational
	MediaType     C.enum_AVMediaType
	AvPictureType C.enum_AVPictureType
	PixelFormat   C.enum_AVPixelFormat
	File          C.FILE
	AVRounding    C.enum_AVRounding
)

const (
	AV_ROUND_ZERO        = int(C.AV_ROUND_ZERO)
	AV_ROUND_INF         = int(C.AV_ROUND_INF)
	AV_ROUND_DOWN        = int(C.AV_ROUND_DOWN)
	AV_ROUND_UP          = int(C.AV_ROUND_UP)
	AV_ROUND_NEAR_INF    = int(C.AV_ROUND_NEAR_INF)
	AV_ROUND_PASS_MINMAX = int(C.AV_ROUND_PASS_MINMAX)
)

//Return the LIBAvUTIL_VERSION_INT constant.
func AvutilVersion() uint {
	return uint(C.avutil_version())
}

//Return the libavutil build-time configuration.
func AvutilConfiguration() string {
	return C.GoString(C.avutil_configuration())
}

//Return the libavutil license.
func AvutilLicense() string {
	return C.GoString(C.avutil_license())
}

//Return a string describing the media_type enum, NULL if media_type is unknown.
func AvGetMediaTypeString(mt MediaType) string {
	return C.GoString(C.av_get_media_type_string((C.enum_AVMediaType)(mt)))
}

//Return a single letter to describe the given picture type pict_type.
func AvGetPictureTypeChar(pt AvPictureType) string {
	return string(C.av_get_picture_type_char((C.enum_AVPictureType)(pt)))
}

//Return x default pointer in case p is NULL.
func AvXIfNull(p, x int) {
	C.av_x_if_null(unsafe.Pointer(&p), unsafe.Pointer(&x))
}

//Compute the length of an integer list.
func AvIntListLengthForSize(e uint, l int, t uint64) uint {
	return uint(C.av_int_list_length_for_size(C.uint(e), unsafe.Pointer(&l), (C.uint64_t)(t)))
}

//Open a file using a UTF-8 filename.
func AvFopenUtf8(p, m string) *File {
	f := C.av_fopen_utf8(C.CString(p), C.CString(m))
	return (*File)(f)
}

//Return the fractional representation of the internal time base.
func AvGetTimeBaseQ() Rational {
	return (Rational)(C.av_get_time_base_q())
}

func AvRescaleQRnd(a int64, bq Rational, cq Rational, rnd AVRounding) int64 {
	return int64(C.av_rescale_q_rnd(C.long(a), C.struct_AVRational(bq), C.struct_AVRational(cq), uint32(rnd)))
}

func AvRescaleQ(a int64, bq Rational, cq Rational) int64 {
	return int64(C.av_rescale_q(C.long(a), C.struct_AVRational(bq), C.struct_AVRational(cq)))
}

package main

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
	"github.com/giorgisio/goav/swscale"

	"log"
)

const (
	OUTPUT_PIX_FMT avcodec.PixelFormat = avcodec.AV_PIX_FMT_YUV
	// OUTPUT_PIX_FMT avcodec.PixelFormat = avcodec.PixelFormat(avcodec.AV_CODEC_ID_MPEG4)
)

// SaveFrame writes a single frame to disk as a PPM file
func SaveFrame(frame *avutil.Frame, width, height, frameNumber int) {
	// Open file
	fileName := fmt.Sprintf("frame%d.ppm", frameNumber)
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error Reading")
	}
	defer file.Close()

	// Write header
	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
	file.Write([]byte(header))

	// Write pixel data
	for y := 0; y < height; y++ {
		data0 := avutil.Data(frame)[0]
		buf := make([]byte, width*3)
		startPos := uintptr(unsafe.Pointer(data0)) + uintptr(y)*uintptr(avutil.Linesize(frame)[0])
		for i := 0; i < width*3; i++ {
			element := *(*uint8)(unsafe.Pointer(startPos + uintptr(i)))
			buf[i] = element
		}
		file.Write(buf)
	}
}

func Encode(ctx *avcodec.Context, frame *avcodec.Frame, pkt *avcodec.Packet, f *os.File) error {
	if ret := ctx.AvcodecSendFrame(frame); ret < 0 {
		return errors.New("Error sending a frame for encoding")
	}
	for ret := -1; ret < 0; {
		if ret = ctx.AvcodecReceivePacket(pkt); ret == avutil.AvErrorEAGAIN || ret == avutil.AvErrorEOF {
			return avutil.ErrorFromCode(ret)
		} else if ret < 0 {
			return errors.New("Error during encoding")
		}
		if _, err := f.Write((*[1 << 30]byte)(unsafe.Pointer(pkt.Data()))[:pkt.Size()]); err != nil {
			return err
		}
	}
	return nil
}

// static void encode(AVCodecContext *enc_ctx, AVFrame *frame, AVPacket *pkt,
// 	FILE *outfile)
// {
// int ret;

// /* send the frame to the encoder */
// if (frame)
// printf("Send frame %3"PRId64"\n", frame->pts);

// ret = avcodec_send_frame(enc_ctx, frame);
// if (ret < 0) {
// fprintf(stderr, "Error sending a frame for encoding\n");
// exit(1);
// }

// while (ret >= 0) {
// ret = avcodec_receive_packet(enc_ctx, pkt);
// if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
// return;
// else if (ret < 0) {
// fprintf(stderr, "Error during encoding\n");
// exit(1);
// }

// printf("Write packet %3"PRId64" (size=%5d)\n", pkt->pts, pkt->size);
// fwrite(pkt->data, 1, pkt->size, outfile);
// av_packet_unref(pkt);
// }
// }

func main() {

	// 打开输入文件

	inputCtx := avformat.AvformatAllocContext()

	if ret := avformat.AvformatOpenInput(&inputCtx, "sample.mp4", nil, nil); ret < 0 {
		log.Fatalf("open input fail %d", ret)
	}
	defer inputCtx.AvformatCloseInput()

	if ret := inputCtx.AvformatFindStreamInfo(nil); ret < 0 {
		log.Fatalf("find stream fail %d", ret)
	}

	videoIndex := -1
	for i := 0; i < int(inputCtx.NbStreams()); i++ {
		codecCtx := inputCtx.Streams()[i].Codec()
		if codecCtx.GetCodecType() == avformat.AVMEDIA_TYPE_VIDEO {
			videoIndex = i
			break
		}
	}

	if videoIndex < 0 {
		log.Fatal("No audio stream found")
	}

	// 打开解码器
	codecCtx := inputCtx.Streams()[videoIndex].Codec()
	codec := avcodec.AvcodecFindDecoder(avcodec.CodecId(codecCtx.GetCodecId()))
	if codec == nil {
		log.Fatal("Unsupported codec")
	}

	codecCtx2 := (*avcodec.Context)(unsafe.Pointer(codecCtx))

	if ret := codecCtx2.AvcodecOpen2(codec, nil); ret < 0 {
		log.Fatalf("codecopen fail %d", ret)
	}
	defer codecCtx2.AvcodecClose()

	// 打开输出文件
	outputFileName := "output.mp4"
	outputFileFormat := "mp4"
	// outputFileName := "output.h264"
	// outputFileFormat := "h264"
	outputFmt := avformat.AvGuessFormat(outputFileFormat, outputFileName, "")
	if outputFmt == nil {
		log.Fatal("Failed to guess output format")
	}

	outputCtx := avformat.AvformatAllocContext()
	if ret := avformat.AvformatAllocOutputContext2(&outputCtx, outputFmt, outputFileFormat, outputFileName); ret < 0 {
		log.Fatalf("fail to allocate output context %d", ret)
	}

	if pb, err := avformat.AvIOOpen(outputFileName, avformat.AVIO_FLAG_WRITE); err != nil {
		log.Fatalf("fail to open avio %s", err)
	} else {
		outputCtx.SetPb(pb)
	}

	// 打开编码器
	outputCodec := avcodec.AvcodecFindEncoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
	var stream *avformat.Stream
	// if stream = outputCtx.AvformatNewStream((*avformat.AvCodec)(unsafe.Pointer(codec))); stream == nil {
	if stream = outputCtx.AvformatNewStream((*avformat.AvCodec)(unsafe.Pointer(outputCodec))); stream == nil {
		log.Fatalln("add stream fail")
	}
	outputCodecPar := stream.CodecParameters()
	outputCodecPar.AvCodecSetId(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
	outputCodecPar.AvCodecSetType(avformat.AVMEDIA_TYPE_VIDEO)
	outputCodecPar.AvCodecSetBitRate(200000)
	outputCodecPar.AvCodecSetWidth(640)
	outputCodecPar.AvCodecSetHeight(480)
	outputCodecPar.AvCodecSetFormat(OUTPUT_PIX_FMT)
	// if ret := outputCodecPar.AvCodecParametersCopy(inputCtx.Streams()[videoIndex].CodecParameters()); ret < 0 {
	// 	log.Fatal("copy codec parameters error")
	// }

	outputCodecCtx := stream.Codec()
	// outputCodecCtx := outputCtx.Streams()[videoIndex].Codec()
	// outputCodec := avcodec.AvcodecFindEncoder(avcodec.CodecId(outputCodecCtx.GetCodecId()))
	if outputCodec == nil {
		log.Fatal("Unsupported codec")
	}

	outputCodecCtx2 := (*avcodec.Context)(unsafe.Pointer(outputCodecCtx))
	// outputCodecCtx2 := (*avcodec.Context)(unsafe.Pointer(codecCtx))
	outputCodecCtx.SetCodecId(avformat.CodecId(avcodec.AV_CODEC_ID_H264))
	outputCodecCtx.SetCodecType(avformat.AVMEDIA_TYPE_VIDEO)
	outputCodecCtx.SetBitRate(200000)
	outputCodecCtx.SetWidth(640)
	outputCodecCtx.SetHeight(480)
	timeBase := avcodec.Rational{}
	timeBase.Set(1, 25)
	outputCodecCtx.SetTimeBase(timeBase)
	outputCodecCtx.SetPktTimeBase(timeBase)
	timeBase.Set(25, 1)
	outputCodecCtx.SetFramerate(timeBase)
	// outputCodecCtx.SetPixelFormat(avcodec.AV_PIX_FMT_YUVA420P9)
	outputCodecCtx.SetPixelFormat(OUTPUT_PIX_FMT)
	outputCodecCtx.SetQMin(10)
	outputCodecCtx.SetQMax(51)
	outputCodecCtx.SetGopSize(codecCtx.GetGopSize())
	var flags int = 0
	// flags |= avformat.AV_CODEC_FLAG_FRAME_DURATION
	outputCodecCtx.SetFlags(flags)
	var flags2 int = 0
	outputCodecCtx.SetFlags2(flags2)

	if ret := outputCodecCtx2.AvcodecOpen2(outputCodec, nil); ret < 0 {
		log.Fatalf("codecopen fail %s", avutil.ErrorFromCode(ret))
	}
	defer outputCodecCtx2.AvcodecClose()

	// 写入输出头
	var param *avutil.Dictionary
	param.AvDictSet("preset", "low", 0)
	param.AvDictSet("tune", "zerolatency", 0)
	if ret := outputCtx.AvformatWriteHeader(&param); ret < 0 {
		log.Fatalf("fail to write header %d", ret)
	}

	pkt := avcodec.AvPacketAlloc()
	defer pkt.AvFreePacket()
	frame := avutil.AvFrameAlloc()
	defer avutil.AvFrameFree(frame)

	// 编码后的包
	encPkt := avcodec.AvPacketAlloc()
	defer encPkt.AvFreePacket()

	outputFrame := avutil.AvFrameAlloc()
	if outputFrame == nil {
		log.Println("Unable to allocate RGB Frame")
		return
	}

	// 分配重新编码后的帧
	// Determine required buffer size and allocate buffer
	numBytes := uintptr(avcodec.AvpictureGetSize(OUTPUT_PIX_FMT, outputCodecCtx2.Width(),
		outputCodecCtx2.Height()))
	buffer := avutil.AvMalloc(numBytes)

	// Assign appropriate parts of buffer to image planes in pFrameRGB
	// Note that pFrameRGB is an AVFrame, but AVFrame is a superset
	// of AVPicture
	avp := (*avcodec.Picture)(unsafe.Pointer(outputFrame))
	avp.AvpictureFill((*uint8)(buffer), OUTPUT_PIX_FMT, outputCodecCtx2.Width(), outputCodecCtx2.Height())

	swsCtx := swscale.SwsGetcontext(
		codecCtx2.Width(),
		codecCtx2.Height(),
		(swscale.PixelFormat)(codecCtx2.PixFmt()),
		outputCodecCtx2.Width(),
		outputCodecCtx2.Height(),
		// avcodec.AV_PIX_FMT_RGB24,
		// avcodec.AV_PIX_FMT_YUV420P10,
		swscale.PixelFormat(OUTPUT_PIX_FMT),
		avcodec.SWS_BILINEAR,
		nil,
		nil,
		nil,
	)

	// outputFile, err := os.Create("output.mp4")
	// if err != nil {
	// 	log.Println("Error Reading")
	// }
	// defer outputFile.Close()

	for frameNumber, errCount, bq, cq := 1,
		0,
		inputCtx.Streams()[videoIndex].TimeBase(),
		stream.TimeBase(); inputCtx.AvReadFrame(pkt) >= 0; {
		if pkt.StreamIndex() != videoIndex {
			continue
		}

		// var g int
		// if ret := codecCtx2.AvcodecDecodeVideo2((*avcodec.Frame)(unsafe.Pointer(frame)), &g, pkt); ret < 0 || g == 0 {
		// 	continue
		// }

		for response := codecCtx2.AvcodecSendPacket(pkt); response >= 0; {
			response = codecCtx2.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(frame)))
			if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
				break
			} else if response < 0 {
				log.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
				if errCount++; errCount > 50 {
					return
				}
				break
				// return
			}
			errCount = 0

			swscale.SwsScale2(swsCtx, avutil.Data(frame),
				avutil.Linesize(frame), 0, codecCtx2.Height(),
				avutil.Data(outputFrame), avutil.Linesize(outputFrame))

			if frameNumber <= 5 {
				SaveFrame(outputFrame, outputCodecCtx2.Width(), outputCodecCtx2.Height(), frameNumber)
				frameNumber++
			}

			// Save the frame to disk
			// SaveFrame(outputFrame, codecCtx2.Width(), codecCtx2.Height(), frameNumber)

			// 编码
			avOutputFrame := (*avcodec.Frame)(unsafe.Pointer(outputFrame))
			avOutputFrame.CopyFrameInfo((*avcodec.Frame)(unsafe.Pointer(frame)))
			avOutputFrame.SetWidth(int32(outputCodecCtx2.Width()))
			avOutputFrame.SetHeight(int32(outputCodecCtx2.Height()))
			avOutputFrame.SetFormat(OUTPUT_PIX_FMT)

			// 写入输出帧
			// if err := Encode(outputCodecCtx2, avOutputFrame, encPkt, outputFile); err != nil {
			// 	log.Printf("write frame error: %s\n", err)
			// }

			var gp int
			if ret := outputCodecCtx2.AvcodecEncodeVideo2(encPkt, avOutputFrame, &gp); ret < 0 || gp == 0 {
				continue
			}

			// pts := (*avcodec.Frame)(unsafe.Pointer(frame)).Pts()
			// // dts := (*avcodec.Frame)(unsafe.Pointer(frame)).PktDts()

			pts := avutil.AvRescaleQRnd(
				pkt.Pts(),
				*(*avutil.Rational)(unsafe.Pointer(&bq)),
				*(*avutil.Rational)(unsafe.Pointer(&cq)),
				avutil.AVRounding(int(avutil.AV_ROUND_NEAR_INF)|int(avutil.AV_ROUND_PASS_MINMAX)),
			)
			dts := avutil.AvRescaleQRnd(
				pkt.Dts(),
				*(*avutil.Rational)(unsafe.Pointer(&bq)),
				*(*avutil.Rational)(unsafe.Pointer(&cq)),
				avutil.AVRounding(int(avutil.AV_ROUND_NEAR_INF)|int(avutil.AV_ROUND_PASS_MINMAX)),
			)
			duration := avutil.AvRescaleQ(
				pkt.Duration(),
				*(*avutil.Rational)(unsafe.Pointer(&bq)),
				*(*avutil.Rational)(unsafe.Pointer(&cq)),
			)
			encPkt.SetPts(pts)
			encPkt.SetDts(dts)
			encPkt.SetDuration(duration)
			// encPkt.AvPacketRescaleTs(codecCtx.GetTimeBase(), outputCodecCtx.GetTimeBase())

			// outputCtx.AvInterleavedWriteFrame(encPkt)
			outputCtx.AvWriteFrame(encPkt)

		}

	}

	// 写入输出尾

	// if err := Encode(outputCodecCtx2, nil, encPkt, outputFile); err != nil {
	// 	log.Printf("write frame error: %s\n", err)
	// }

	if ret := outputCtx.AvWriteTrailer(); ret < 0 {
		log.Fatalf("write trailer error %d", ret)
	}

}

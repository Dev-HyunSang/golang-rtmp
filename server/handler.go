package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

var _ rtmp.Handler = (*Handler)(nil)

type Handler struct {
	rtmp.DefaultHandler
	flvFile *os.File
	flvEnc  *flv.Encoder
}

func (h *Handler) OnServe(conn *rtmp.Conn) {
}

func (h *Handler) OnConnect(timestamp uint32, cmd *rtmpmsg.NetConnectionConnect) error {
	log.Printf("OnConnect: %#v", cmd)
	return nil
}

func (h *Handler) OnCreateStream(timestamp uint32, cmd *rtmpmsg.NetConnectionCreateStream) error {
	log.Printf("OnCreateStream: %#v", cmd)
	return nil
}

// func (h *Handler) OnPublish(_ *rtmp.StreamContext, timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error
// rmtp.StreamContext 패키지 내에 누락되어 있음.
func (h *Handler) OnPublish(timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	log.Printf("onPublish: %#v", cmd)

	if cmd.PublishingName == "" {
		return errors.New("PublishingName is empty")
	}

	p := filepath.Join(
		os.TempDir(),
		filepath.Clean(filepath.Join("/", fmt.Sprintf("%s.flv", cmd.PublishingName))))
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.New("Failed to create flv file")
	}
	h.flvFile = f

	enc, err := flv.NewEncoder(f, flv.FlagsAudio|flv.FlagsVideo)
	if err != nil {
		_ = f.Close()
		return errors.New("Failed to create flv encoder")
	}
	h.flvEnc = enc

	return nil
}

func (h *Handler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	r := bytes.NewBuffer(data.Payload)

	var script flvtag.ScriptData
	if err := flvtag.DecodeScriptData(r, &script); err != nil {
		log.Printf("Failed to decode script data: Err = %+v", err)
		return nil
	}

	log.Printf("SetDataFrame: Script = %#v", script)

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeScriptData,
		Timestamp: timestamp,
		Data:      &script,
	}); err != nil {
		log.Printf("Failed to write script data: Err %+v", err)
	}

	return nil
}

func (h *Handler) OnVideo(timestamp uint32, payload io.Reader) error {
	var video flvtag.VideoData
	if err := flvtag.DecodeVideoData(payload, &video); err != nil {
		return err
	}

	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, video.Data); err != nil {
		return err
	}

	video.Data = flvBody

	log.Printf("FLV Video Data: Timestamp = %d, FrameType = %+v, CodecID = %+v, AVCPacketType = %+v, CT = %+v, Data length = %+v",
		timestamp,
		video.FrameType,
		video.CodecID,
		video.AVCPacketType,
		video.CompositionTime,
		len(flvBody.Bytes()),
	)

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      &video,
	}); err != nil {
		log.Printf("Failed to write video: Err = %+v", err)
	}

	return nil
}

func (h *Handler) OnClose() {
	log.Printf("OnClose")

	if h.flvFile != nil {
		_ = h.flvFile.Close()
	}
}

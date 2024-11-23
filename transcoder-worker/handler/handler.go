package handler

import (
	"log"
	"os"

	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/communicator"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/transcoder"
)

type Handler struct {
	destsMap     map[string]*os.File
	communicator *communicator.Communicator
}

func NewHandler(communicator *communicator.Communicator) *Handler {
	return &Handler{
		destsMap:     make(map[string]*os.File),
		communicator: communicator,
	}
}

func (h *Handler) Init(videoId string) {
	_, err := os.Stat("uploads")
	if os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}
	dest, err := os.Create("uploads/" + videoId)
	if err != nil {
		log.Printf("Error creating file: %v", err)
	}
	h.destsMap[videoId] = dest
}

func (h *Handler) Process(videoId string, data []byte) {
	dest := h.destsMap[videoId]
	if dest != nil {
		dest.Write(data)
	}
}

func (h *Handler) End(videoId string) {
	dest := h.destsMap[videoId]
	if dest == nil {
		return
	}
	dest.Close()
	delete(h.destsMap, videoId)

	scales := transcoder.Transcode(videoId)

	h.communicator.PublishVideo(videoId, scales)
}

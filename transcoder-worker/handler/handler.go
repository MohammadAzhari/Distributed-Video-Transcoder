package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/transcoder"
)

type Handler struct {
	destsMap            map[string]*os.File
	videoServiceAddress string
}

func NewHandler(videoServiceAddress string) *Handler {
	return &Handler{
		destsMap:            make(map[string]*os.File),
		videoServiceAddress: videoServiceAddress,
	}
}

func (h *Handler) Init(key string) {
	dest, err := os.Create("uploads/" + key)
	if err != nil {
		log.Printf("Error creating file: %v", err)
	}
	h.destsMap[key] = dest
}

func (h *Handler) Process(key string, data []byte) {
	dest := h.destsMap[key]
	if dest != nil {
		dest.Write(data)
	}
}

func (h *Handler) End(key string) {
	dest := h.destsMap[key]
	if dest == nil {
		return
	}
	dest.Close()
	delete(h.destsMap, key)

	scales := transcoder.Transcode(key)

	data := map[string]any{
		"scales": scales,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}
	// send http request to the video service that the transcoding is done
	res, err := http.Post(h.videoServiceAddress+"/prossess-completed/"+key, "Application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending request to video service: %v", err)
	}
	log.Printf("Response from video service: %v", res.Status)
	res.Body.Close()
}

package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/transcoder"
)

type Handler struct {
	destsMap map[string]*os.File
}

func NewHandler() *Handler {
	return &Handler{
		destsMap: make(map[string]*os.File),
	}
}

func (h *Handler) Init(key string) {
	dest, err := os.Create("uploads/" + key)
	if err != nil {
		log.Fatal(err)
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

	transcoder.Transcode(key)
	// send http request to the video service that the transcoding is done
	res, err := http.Post("http://localhost:8080/prossess-completed/"+key, "", nil)
	if err != nil {
		log.Printf("Error sending request to video service: %v", err)
	}
	log.Printf("Response from video service: %v", res.Status)
	res.Body.Close()
}

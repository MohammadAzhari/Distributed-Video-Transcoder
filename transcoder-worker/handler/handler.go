package handler

import (
	"log"
	"os"
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
	dest, err := os.Create(string(key))
	if err != nil {
		log.Fatal(err)
	}
	h.destsMap[string(key)] = dest
}

func (h *Handler) Process(key string, data []byte) {
	dest := h.destsMap[string(key)]
	if dest != nil {
		dest.Write(data)
	}
}

func (h *Handler) End(key string) {
	dest := h.destsMap[string(key)]
	if dest != nil {
		dest.Close()
		delete(h.destsMap, string(key))
	}
}

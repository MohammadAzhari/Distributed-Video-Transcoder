package communicator

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Communicator struct {
	videoServiceAddress string
	currentPort         string
}

func NewCommunicator(videoServiceAddress string, currentPort string) *Communicator {
	return &Communicator{
		videoServiceAddress: videoServiceAddress,
		currentPort:         currentPort,
	}
}

func (c *Communicator) PublishVideo(videoId string, scales []string) {
	data := map[string]any{
		"scales": scales,
		"port":   c.currentPort,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}
	// send http request to the video service that the transcoding is done
	res, err := http.Post(c.videoServiceAddress+"/prossess-completed/"+videoId, "Application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending request to video service: %v", err)
	}
	log.Printf("Response from video service: %v", res.Status)
	res.Body.Close()
}

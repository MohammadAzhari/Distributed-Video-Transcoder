package main

import (
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/api"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/producer"
)

func main() {

	producer := producer.NewProducer()

	server := api.NewServer(producer)
	server.Start(":8080")
}

package main

import (
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/api"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/producer"
)

const (
	kafkaHost = "localhost:9092"
	topic     = "video"
)

func main() {
	producer := producer.NewProducer(kafkaHost, topic)
	defer producer.Close()

	server := api.NewServer(producer)
	server.Start(":8080")
}

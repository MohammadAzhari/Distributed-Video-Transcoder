package main

import "github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/consumer"

const (
	kafkaHost = "localhost:9092"
	topic     = "video"
)

func main() {
	consumer.NewConsumer(kafkaHost, topic)
}

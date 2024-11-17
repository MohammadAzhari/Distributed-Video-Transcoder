package main

import (
	"net/http"

	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/consumer"
	"github.com/gin-gonic/gin"
)

const (
	kafkaHost = "localhost:9092"
	topic     = "video"
)

func main() {
	consumer := consumer.NewConsumer(kafkaHost, topic)
	defer consumer.Close()

	router := gin.Default()
	router.StaticFS("/", http.Dir("uploads"))
	router.Run(":8081")
}

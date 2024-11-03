package api

import (
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/producer"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router   *gin.Engine
	producer *producer.Producer
}

func NewServer(producer *producer.Producer) *Server {
	router := gin.Default()

	server := &Server{
		router: router,
	}

	server.setupRoutes()
	server.producer = producer

	return server
}

func (s *Server) setupRoutes() {
	s.router.POST("/upload-video", s.uploadVideo)
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

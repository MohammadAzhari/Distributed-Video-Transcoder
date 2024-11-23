package api

import (
	db "github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/db/sqlc"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/producer"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router   *gin.Engine
	producer *producer.Producer
	store    *db.Store
}

func NewServer(producer *producer.Producer, store *db.Store) *Server {
	router := gin.Default()

	server := &Server{
		router: router,
		store:  store,
	}

	server.setupRoutes()
	server.producer = producer

	return server
}

func (s *Server) setupRoutes() {
	s.router.POST("/upload-video", s.uploadVideo)
	s.router.POST("/prossess-completed/:videoId", s.processCompleted)
	s.router.GET("/video/:videoId", s.getVideo)

	s.router.GET("/", func(ctx *gin.Context) {
		ctx.File("frontend/index.html")
	})
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

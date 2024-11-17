package api

import (
	"errors"
	"io"
	"strings"

	db "github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/db/sqlc"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/producer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) uploadVideo(ctx *gin.Context) {
	// receive bytes from request let's say n bytes at a time
	file, header, err := ctx.Request.FormFile("video")

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !strings.HasSuffix(header.Filename, ".mp4") {
		ctx.JSON(400, gin.H{"error": "Only mp4 files are allowed"})
		return
	}

	buffer := make([]byte, 1024*8)

	s.producer.SendMessage(&producer.Message{
		Key:   header.Filename,
		Value: "new file",
	})

	for {
		n, err := file.Read(buffer)

		if err != nil && !errors.Is(err, io.EOF) {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if n == 0 {
			break
		}

		_, _, err = s.producer.SendMessage(&producer.Message{
			Key:   header.Filename,
			Value: string(buffer[:n]),
		})
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	s.producer.SendMessage(&producer.Message{
		Key:   header.Filename,
		Value: "close file",
	})

	video, err := s.store.CreateVideo(ctx, header.Filename)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, video)
}

func (s *Server) processCompleted(ctx *gin.Context) {
	videoId, ok := ctx.Params.Get("videoId")
	ip := ctx.ClientIP()

	if !ok {
		ctx.JSON(400, gin.H{"error": "Invalid video id"})
		return
	}

	arg := db.PublishVideoParams{
		WorkerIp: pgtype.Text{
			String: ip,
			Valid:  true,
		},
		ID: uuid.UUID([]byte(videoId)),
	}

	video, err := s.store.PublishVideo(ctx, arg)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, video)
}

func (s *Server) getVideo(ctx *gin.Context) {
	videoId, ok := ctx.Params.Get("videoId")

	if !ok {
		ctx.JSON(400, gin.H{"error": "Invalid video id"})
		return
	}

	video, err := s.store.GetVideo(ctx, uuid.UUID([]byte(videoId)))

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, video)
}

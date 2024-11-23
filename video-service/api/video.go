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

const (
	newFile   = "new file"
	closeFile = "close file"
)

func (s *Server) uploadVideo(ctx *gin.Context) {
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

	arg := db.CreateVideoParams{
		Filename: header.Filename,
		ID:       uuid.New(),
	}

	video, err := s.store.CreateVideo(ctx, arg)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	s.producer.SendMessage(&producer.Message{
		Key:   video.ID.String(),
		Value: newFile,
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
			Key:   video.ID.String(),
			Value: string(buffer[:n]),
		})
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	s.producer.SendMessage(&producer.Message{
		Key:   video.ID.String(),
		Value: closeFile,
	})

	ctx.JSON(200, video)
}

type ProccessCompletedRequest struct {
	Scales []string `json:"scales"`
	Port   string   `json:"port"`
}

func (s *Server) processCompleted(ctx *gin.Context) {
	videoId, ok := ctx.Params.Get("videoId")

	if !ok {
		ctx.JSON(400, gin.H{"error": "Invalid video id"})
		return
	}

	var request ProccessCompletedRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	workerIp := ctx.ClientIP() + request.Port

	uuid, err := uuid.Parse(videoId)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	arg := db.PublishVideoParams{
		WorkerIp: pgtype.Text{
			String: workerIp,
			Valid:  true,
		},
		ID:     uuid,
		Scales: request.Scales,
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

	uuid, err := uuid.Parse(videoId)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	video, err := s.store.GetVideo(ctx, uuid)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, video)
}

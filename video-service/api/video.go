package api

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/producer"
	"github.com/gin-gonic/gin"
)

func (s *Server) uploadVideo(ctx *gin.Context) {
	// receive bytes from request let's say n bytes at a time
	file, header, err := ctx.Request.FormFile("video")

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("File name: %s\n", header.Filename)
	fmt.Printf("File size: %d\n", header.Size)

	buffer := make([]byte, 1024*8)

	s.producer.SendMessage(&producer.Message{
		Key:   header.Filename,
		Value: "new file",
	})

	for {
		n, err := file.Read(buffer)
		fmt.Println("n:", n)

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
			log.Println("Here", err)
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	s.producer.SendMessage(&producer.Message{
		Key:   header.Filename,
		Value: "close file",
	})

	ctx.JSON(200, gin.H{"message": "Video uploaded successfully"})
}

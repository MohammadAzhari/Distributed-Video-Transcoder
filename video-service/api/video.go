package api

import (
	"errors"
	"fmt"
	"io"

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

	// dest, err := os.Create("api/uploads/" + header.Filename)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// defer dest.Close()

	buffer := make([]byte, 1024*8)

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

		s.producer.Produce("video", 1, "video", string(buffer))
	}

	ctx.JSON(200, gin.H{"message": "Video uploaded successfully"})
}

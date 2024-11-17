package transcoder

import (
	"log"
	"os/exec"
	"sync"
)

var scalesMap = map[string]string{
	"480p":  "scale=854:480",
	"720p":  "scale=1280:720",
	"1080p": "scale=1920:1080",
}
var scales = make([]string, 0, len(scalesMap))

func init() {
	for scale := range scalesMap {
		scales = append(scales, scale)
	}
}

func transcodeVideo(inputFileName, scale string) {
	outputFileName := inputFileName[:len(inputFileName)-4] + "_" + scale + ".mp4"
	cmd := exec.Command("ffmpeg", "-i", "uploads/"+inputFileName,
		"-vf", scalesMap[scale], "-c:v", "libx264", "-crf", "23",
		"-c:a", "copy", "uploads/"+outputFileName)

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running ffmpeg for %s: %v", inputFileName, err)
	}
}

func Transcode(key string) {
	var wg sync.WaitGroup
	for _, scale := range scales {
		wg.Add(1)
		go func(scale string) {
			defer wg.Done()
			transcodeVideo(key, scale)
		}(scale)
	}
	wg.Wait()
}

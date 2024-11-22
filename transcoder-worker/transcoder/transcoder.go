package transcoder

import (
	"log"
	"os"
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

func transcodeVideo(inputFileName, scale string, ch chan string) {
	outputFileName := inputFileName + "_" + scale + ".mp4"
	cmd := exec.Command("ffmpeg", "-i", "uploads/"+inputFileName,
		"-vf", scalesMap[scale], "-c:v", "libx264", "-crf", "23",
		"-c:a", "copy", "uploads/"+outputFileName)

	if scale == "1080p" {
		cmd.Run()
	}
	// check if the video exist in the disk
	if _, err := os.Stat("uploads/" + outputFileName); err != nil {
		log.Printf("Error transcoding video: %v, scale: %v", err, scale)
		ch <- scale
	}
}

func Transcode(key string) []string {
	var wg sync.WaitGroup
	errChan := make(chan string, len(scales))

	for _, scale := range scales {
		wg.Add(1)
		go func(scale string) {
			defer wg.Done()
			transcodeVideo(key, scale, errChan)
		}(scale)
	}

	wg.Wait()
	close(errChan)
	erroredScales := make(map[string]bool)
	for scale := range errChan {
		erroredScales[scale] = true
	}

	transcodedScales := make([]string, 0)
	for _, scale := range scales {
		if _, ok := erroredScales[scale]; !ok {
			transcodedScales = append(transcodedScales, scale)
		}
	}

	return transcodedScales
}

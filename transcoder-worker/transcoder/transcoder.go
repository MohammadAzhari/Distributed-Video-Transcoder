package transcoder

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
)

var scalesMap = make(map[string]string)
var scales = []string{"480p", "720p", "1080p"}

func init() {
	scalesMap["480p"] = "scale=854:480"
	scalesMap["720p"] = "scale=1280:720"
	scalesMap["1080p"] = "scale=1920:1080"
}

func transcode(key string, scale string, wg *sync.WaitGroup) {
	wg.Add(1)
	cmd := exec.Command("ffmpeg", "-i", key, "-vf", scalesMap[scale], "-c:v", "libx264", "-crf", "23", key+"_"+scale+".mp4", "-c:a", "copy", key+"_"+scale+".mp4")

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error creating stderr pipe for %s: %v", key, err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating stdout pipe for %s: %v", key, err)
	}

	go logOutput("stderr", stderr)
	go logOutput("stdout", stdout)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running ffmpeg for %s: %v", key, err)
	}
	wg.Done()
}

func Transcode(key string) {
	wg := new(sync.WaitGroup)
	for _, scale := range scales {
		go transcode(key, scale, wg)
	}
	wg.Wait()
}

func logOutput(prefix string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.Printf("[%s] %s", prefix, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading %s: %v", prefix, err)
	}
}

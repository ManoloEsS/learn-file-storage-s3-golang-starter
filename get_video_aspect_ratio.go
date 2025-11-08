package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

func getVideoAspectRatio(filepath string) (string, error) {
	var videoData bytes.Buffer
	var JsonVideoData JSONvideoData

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filepath)
	cmd.Stdout = &videoData

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %v", err)
	}

	if err := json.Unmarshal(videoData.Bytes(), &JsonVideoData); err != nil {
		return "", fmt.Errorf("Could not parse ffprobe output: %v", err)
	}

	if len(JsonVideoData.Streams) == 0 {
		return "", errors.New("no video streams found")
	}

	var videoStream Stream
	for _, stream := range JsonVideoData.Streams {
		if stream.CodecType == "video" {
			videoStream = stream
		}
	}

	return calculateAspectRatio(videoStream.Width, videoStream.Height), nil

}

type JSONvideoData struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	CodecType string `json:"codec_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func calculateAspectRatio(w, h int) (aspectRatio string) {
	ratio := float32(w) / float32(h)

	if ratio < 1 {
		if ratio*16 > 8 && ratio*16 < 10 {
			return "9:16"
		}
	} else if ratio > 1 {
		if ratio*9 > 15 && ratio*9 < 17 {
			return "16:9"
		}
	}

	return "other"

}

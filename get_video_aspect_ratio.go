package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

func getVideoAspectRatio(filepath string) (string, error) {
	var videoData bytes.Buffer
	var JsonVideoData JSONvideoData

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filepath)
	cmd.Stdout = &videoData

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(videoData.Bytes(), JsonVideoData)
	if err != nil {
		return "", err
	}

}

type JSONvideoData struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	CodecType string `json:"codec_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

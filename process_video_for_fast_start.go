package main

import (
	"fmt"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {
	outputFile := fmt.Sprintf("%s.processing", filePath)

	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFile)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Couldn't process video with ffmpeg: %w", err)
	}

	return outputFile, nil
}

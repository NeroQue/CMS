package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func IsVideoFile(filename string) bool {
	filename = strings.ToLower(filename)
	return strings.HasSuffix(filename, ".mp4") ||
		strings.HasSuffix(filename, ".avi") ||
		strings.HasSuffix(filename, ".mov") ||
		strings.HasSuffix(filename, ".mkv") ||
		strings.HasSuffix(filename, ".wmv") ||
		strings.HasSuffix(filename, ".m4v") ||
		strings.HasSuffix(filename, ".m4a") ||
		strings.HasSuffix(filename, ".webm") ||
		strings.HasSuffix(filename, ".flv") ||
		strings.HasSuffix(filename, ".mpeg") ||
		strings.HasSuffix(filename, ".mpg")
}

func ExtractVideoDuration(filePath string) (int, error) {

	type FFprobeOutput struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}

	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", filePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}
	var result FFprobeOutput
	err = json.Unmarshal(output, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}
	durationFloat, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %w", err)
	}
	duration := int(durationFloat)
	log.Printf("Extracted video duration: %d seconds from %s", duration, filePath)
	return duration, nil
}

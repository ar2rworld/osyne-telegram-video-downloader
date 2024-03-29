package downloader

import (
	"os"
	"os/exec"
)

func Convert(path string) (string, error) {
	outputFile := path + ".mp4"

	cmd := exec.Command("ffmpeg", "-y", "-i", path, "-b:v", "800k", "-c:v", "libx264", "-c:a", "aac", "-b:a", "128k", outputFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outputFile, nil
}

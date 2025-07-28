package downloader

import (
	"context"
	"os"
	"os/exec"
)

func Convert(ctx context.Context, path string) (string, error) {
	outputFile := path + ".mp4"

	cmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", path, "-b:v", "800k", "-c:v", "libx264", "-c:a", "aac", "-b:a", "128k", outputFile)
	cmd.Stderr = os.Stderr

	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

package downloader

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/wader/goutubedl"
)

func DownloadVideo(url string, opts goutubedl.Options) (string, error) {
	goutubedl.Path = "yt-dlp"

	result, err := goutubedl.New(context.Background(), url, opts)
	if err != nil {
		return "", err
	}
	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		return "", err
	}
	defer downloadResult.Close()
	f, err := os.Create("output")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, downloadResult)
	if err != nil {
		return "", err
	}

	outputFile := "output.mp4"

	cmd := exec.Command("ffmpeg", "-y", "-i", "output", "-b:v", "800k", "-c:v", "libx264", "-c:a", "aac", "-b:a", "128k", outputFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return outputFile, nil
}

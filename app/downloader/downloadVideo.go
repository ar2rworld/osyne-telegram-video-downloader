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

	output := "output"

	f, err := os.CreateTemp("", output)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, downloadResult)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func DownloadWithCookies(url, cookiesPath string) (string, error) {
	fileName := "videoDownloadedWithCookies"
	cmd := exec.Command("yt-dlp", "-f", "mp4", "-o", fileName, "--cookies", cookiesPath, url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return fileName, nil
}

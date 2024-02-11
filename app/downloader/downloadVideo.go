package downloader

import (
	"context"
	"io"
	"os"

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

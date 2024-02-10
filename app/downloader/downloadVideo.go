package downloader

import (
	"context"
	"io"
	"os"

	"github.com/wader/goutubedl"
)

func DownloadVideo(url string, opts goutubedl.Options) error {
	goutubedl.Path = "yt-dlp"

	result, err := goutubedl.New(context.Background(), url, opts)
	if err != nil {
		return err
	}
	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		return err
	}
	defer downloadResult.Close()
	f, err := os.Create("output")
	if err != nil {
		return err
	}
	defer f.Close()
	_, copyErr := io.Copy(f, downloadResult)
	if copyErr != nil {
		return copyErr
	}
	return nil
}

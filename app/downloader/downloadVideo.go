package downloader

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/wader/goutubedl"
)

func DownloadVideo (url string, client *http.Client) error {
	goutubedl.Path = "yt-dlp"
	result, err := goutubedl.New(context.Background(), url, goutubedl.Options{HTTPClient: client})
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
	io.Copy(f, downloadResult)
	return nil
}

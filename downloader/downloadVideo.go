package downloader

import (
	"context"
	"io"
	"os"

	"github.com/wader/goutubedl"
)

func DownloadVideo (url string) error {
	goutubedl.Path = os.Getenv("YT_DLP_PATH")

	result, err := goutubedl.New(context.Background(), url, goutubedl.Options{})
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
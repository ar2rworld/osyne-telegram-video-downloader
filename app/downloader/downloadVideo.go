package downloader

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/wader/goutubedl"
)

const magicn = 5

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

	id, err := GenerateUniqueString(magicn)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("output-%s", id)

	f, err := os.Create(output)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, downloadResult)
	if err != nil {
		return "", err
	}

	return output, nil
}

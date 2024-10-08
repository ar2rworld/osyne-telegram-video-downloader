package downloader

import (
	"context"
	"io"
	"math"
	"os"
	"os/exec"

	"github.com/wader/goutubedl"
)

const MaxFileNameLength = 90

func DownloadVideo(url string, opts goutubedl.Options, do *goutubedl.DownloadOptions) (string, error) {
	goutubedl.Path = "yt-dlp"

	result, err := goutubedl.New(context.Background(), url, opts)
	if err != nil {
		return "", err
	}

	if do == nil {
		do = &goutubedl.DownloadOptions{}
	}

	// Weird stuff to make yt-dlp download only audio
	if !do.DownloadAudioOnly {
		do.Filter = "best"
	}

	downloadResult, err := result.DownloadWithOptions(context.Background(), *do)
	if err != nil {
		return "", err
	}
	defer downloadResult.Close()

	output := cutString(result.Info.Title, MaxFileNameLength) + result.Info.Ext

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

func cutString(s string, maxLength int) string {
	absMaxLength := int(math.Abs(float64(maxLength)))
	if len(s) <= absMaxLength {
		return s
	}
	return s[:absMaxLength]
}

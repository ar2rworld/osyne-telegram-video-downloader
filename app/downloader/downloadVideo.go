package downloader

import (
	"context"
	"io"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/wader/goutubedl"
)

const MaxFileNameLength = 90

func DownloadVideo(ctx context.Context, url string, opts goutubedl.Options, do *goutubedl.DownloadOptions) (string, error) {
	goutubedl.Path = "yt-dlp"

	result, err := goutubedl.New(ctx, url, opts)
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

	downloadResult, err := result.DownloadWithOptions(ctx, *do)
	if err != nil {
		return "", err
	}
	defer downloadResult.Close()

	title := RemoveNonAlphanumericRegex(ConvertToUTF8(result.Info.Title))
	output := cutString(title, MaxFileNameLength) + result.Info.Ext

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

func DownloadWithCookies(ctx context.Context, url, cookiesPath string) (string, error) {
	fileName := "videoDownloadedWithCookies"
	cmd := exec.CommandContext(ctx, "yt-dlp", "-f", "mp4", "-o", fileName, "--cookies", cookiesPath, url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
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

func ConvertToUTF8(s string) string {
	result := make([]byte, 0, len(s))
	for s != "" {
		r, size := utf8.DecodeRuneInString(s)
		if r != utf8.RuneError || size > 1 {
			result = append(result, s[:size]...)
		}

		s = s[size:]
	}

	return string(result)
}

// RemoveNonAlphanumericRegex uses a regular expression to remove non-alphanumeric characters.
// This version also preserves spaces but trims leading and trailing whitespace.
func RemoveNonAlphanumericRegex(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	return strings.TrimSpace(reg.ReplaceAllString(s, ""))
}

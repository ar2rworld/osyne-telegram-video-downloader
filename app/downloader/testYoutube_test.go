package downloader

import (
	"context"
	"testing"

	"github.com/wader/goutubedl"
)

func TestYoutube(t *testing.T) {
	goutubedl.Path = "yt-dlp"
	// Attempt to download section of a video, but currently goutubedl does not support this
	result, err := goutubedl.New(context.Background(), "https://www.youtube.com/watch?v=OyuL5biOQ94",
		goutubedl.Options{
			// DownloadSections: "*0:0-0:10",
		})
	if err != nil {
		t.Errorf("new goutubedl err:\n%s", err)
	}
	downloadResult, err := result.Download(context.Background(), "best[ext=mp4]/best")
	if err != nil {
		t.Errorf("download err:\n%s", err)
	}
	defer downloadResult.Close()

	mybytes := make([]byte, 128)
	_, downloadErr := downloadResult.Read(mybytes)
	if downloadErr != nil {
		t.Errorf("err while reading download result:\n%s", downloadErr)
	}
}

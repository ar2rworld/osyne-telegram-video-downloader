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
			DownloadSections: "*0:0-0:10",
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

func TestCutString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{"Empty string", "", 5, ""},
		{"String shorter than max", "Hello", 10, "Hello"},
		{"String equal to max", "Hello", 5, "Hello"},
		{"String longer than max", "Hello, World!", 5, "Hello"},
		{"Max length 0", "Hello", 0, ""},
		{"Unicode characters", "你好世界", 2, "\xe4\xbd"},
		{"Max length 1", "Hello", 1, "H"},
		{"Spaces in string", "Hello World", 6, "Hello "},
		{"Very long string", "This is a very long string", 10, "This is a "},
		{"Negative max length", "Hello", -5, "Hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cutString(tt.input, tt.maxLength)
			if result != tt.expected {
				t.Errorf("cutString(%q, %d) = %q; want %q", tt.input, tt.maxLength, result, tt.expected)
			}
		})
	}
}

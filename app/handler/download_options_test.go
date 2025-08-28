package handler

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
)

func TestConvertSecondsToMinSec(t *testing.T) {
	testcases := []struct {
		name string
		have int
		want string
	}{
		{"0", 0, "0:0"},
		{"59", 59, "0:59"},
		{"60", 60, "1:0"},
		{"61", 61, "1:1"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := downloader.ConvertSecondsToMinSec(tc.have)
			if got != tc.want {
				t.Errorf("want %s, but got %s", tc.want, got)
			}
		})
	}
}

func TestChangeDownloadSectionsStart(t *testing.T) {
	testcases := []struct {
		name  string
		start int
		want  string
	}{
		{
			"",
			0,
			"*0:0-0:30",
		},
		{
			"",
			119,
			"*1:59-2:29",
		},
		{
			"",
			120,
			"*2:0-2:30",
		},
		{
			"",
			121,
			"*2:1-2:31",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &goutubedl.Options{}
			got := setDownloadSections(opts, tc.start, tc.start+30)
			assert.Equal(t, tc.want, got.DownloadSections)
		})
	}
}

func TestAlterDownloadOptions(t *testing.T) { //nolint: funlen
	// Claude 3.5 Sonnet.
	emptyDownloadOptions := &goutubedl.DownloadOptions{}
	tests := []struct {
		name                    string
		inputURL                string
		inputMessage            string
		expectedOutput          string
		expectedDownloadOptions *goutubedl.DownloadOptions
	}{
		{
			name:                    "Youtube URL",
			inputURL:                "https://www.youtube.com/watch?v=T_JKIkSf93Y",
			inputMessage:            "https://www.youtube.com/watch?v=T_JKIkSf93Y",
			expectedOutput:          "*0:0-0:30",
			expectedDownloadOptions: emptyDownloadOptions,
		},
		{
			name:                    "Youtube URL with sections argument",
			inputURL:                "https://www.youtube.com/watch?v=T_JKIkSf93Y",
			inputMessage:            "-s *1:10-2:10 https://www.youtube.com/watch?v=T_JKIkSf93Y",
			expectedOutput:          "*1:10-2:10",
			expectedDownloadOptions: emptyDownloadOptions,
		},
		{
			name:                    "Youtube URL with current time",
			inputURL:                "https://youtu.be/T_JKIkSf93Y?t=2289",
			inputMessage:            "https://youtu.be/T_JKIkSf93Y?t=2289",
			expectedOutput:          "*38:9-38:39",
			expectedDownloadOptions: emptyDownloadOptions,
		},
		{
			name:                    "Youtube URL with sections arg and current time",
			inputURL:                "https://youtu.be/T_JKIkSf93Y?t=2289",
			inputMessage:            "-s *1:10-2:10 https://youtu.be/T_JKIkSf93Y?t=2289",
			expectedOutput:          "*1:10-2:10",
			expectedDownloadOptions: emptyDownloadOptions,
		},
		{
			name:                    "Youtube URL with current time (alternate format)",
			inputURL:                "https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s",
			inputMessage:            "https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s",
			expectedOutput:          "*3:19-3:49",
			expectedDownloadOptions: emptyDownloadOptions,
		},
		{
			name:                    "Youtube URL with current time (alternate format)",
			inputURL:                "https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s",
			inputMessage:            "https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s",
			expectedOutput:          "*3:19-3:49",
			expectedDownloadOptions: emptyDownloadOptions,
		},
		{
			name:                    "Youtube URL with timestamp and -x",
			inputURL:                "https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s",
			inputMessage:            "-x https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s",
			expectedOutput:          "*3:19-3:49",
			expectedDownloadOptions: &goutubedl.DownloadOptions{DownloadAudioOnly: true},
		},
		{
			name:                    "Youtube URL and -x",
			inputURL:                "https://www.youtube.com/watch?v=AWVUp12XPpU",
			inputMessage:            "-x https://www.youtube.com/watch?v=AWVUp12XPpU&",
			expectedOutput:          "",
			expectedDownloadOptions: &goutubedl.DownloadOptions{DownloadAudioOnly: true},
		},
		{
			name:                    "test if no x and no s",
			inputURL:                "https://www.youtube.com/watch?v=AWVUp12XPpU",
			inputMessage:            "https://www.youtube.com/watch?v=AWVUp12XPpU",
			expectedOutput:          "*0:0-0:30",
			expectedDownloadOptions: &goutubedl.DownloadOptions{DownloadAudioOnly: false},
		},
		{
			name:                    "test if no x and s",
			inputURL:                "https://www.youtube.com/watch?v=AWVUp12XPpU",
			inputMessage:            "-s *0:5-0:35 https://www.youtube.com/watch?v=AWVUp12XPpU",
			expectedOutput:          "*0:5-0:35",
			expectedDownloadOptions: &goutubedl.DownloadOptions{DownloadAudioOnly: false},
		},
		{
			name:                    "test if x and no s",
			inputURL:                "https://www.youtube.com/watch?v=AWVUp12XPpU",
			inputMessage:            "-x https://www.youtube.com/watch?v=AWVUp12XPpU",
			expectedOutput:          "",
			expectedDownloadOptions: &goutubedl.DownloadOptions{DownloadAudioOnly: true},
		},
		{
			name:                    "test if x and s",
			inputURL:                "https://www.youtube.com/watch?v=AWVUp12XPpU",
			inputMessage:            "-s  *0:0-0:40 -x https://www.youtube.com/watch?v=AWVUp12XPpU",
			expectedOutput:          "*0:0-0:40",
			expectedDownloadOptions: &goutubedl.DownloadOptions{DownloadAudioOnly: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &tgbotapi.Update{Message: &tgbotapi.Message{Text: tt.inputMessage}}
			opts := &goutubedl.Options{}
			do := AlterDownloadOptions(u, tt.inputURL, opts)
			assert.Equal(t, tt.expectedOutput, opts.DownloadSections)
			assert.Equal(t, tt.expectedDownloadOptions.DownloadAudioOnly, do.DownloadAudioOnly)
		})
	}
}

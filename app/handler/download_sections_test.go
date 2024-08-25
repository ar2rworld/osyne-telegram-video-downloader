package handler

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/wader/goutubedl"
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
			got := convertSecondsToMinSec(tc.have)
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

// func TestAlterDownloadSections(t *testing.T) {
// 	t.Run("Testing alterDownloadSections Youtube", func(t *testing.T) {
// 		url := "https://www.youtube.com/watch?v=T_JKIkSf93Y"
// 		u := &tgbotapi.Update{Message: &tgbotapi.Message{Text: url}}
// 		opts := &goutubedl.Options{}
// 		alterDownloadSections(u, url, opts)
// 		assert.Equal(t, "*0:0-0:30", opts.DownloadSections)
// 	})

// 	t.Run("Testing alterDownloadSections with sections argument Youtube", func(t *testing.T) {
// 		url := "-s *1:10-2:10 https://www.youtube.com/watch?v=T_JKIkSf93Y"
// 		u := &tgbotapi.Update{Message: &tgbotapi.Message{Text: url}}
// 		opts := &goutubedl.Options{}
// 		alterDownloadSections(u, url, opts)
// 		assert.Equal(t, "*1:10-2:10", opts.DownloadSections)
// 	})

// 	t.Run("Testing alterDownloadSections with current time in url Youtube", func(t *testing.T) {
// 		url := "https://youtu.be/T_JKIkSf93Y?t=2289"
// 		u := &tgbotapi.Update{Message: &tgbotapi.Message{Text: url}}
// 		opts := &goutubedl.Options{}
// 		alterDownloadSections(u, url, opts)
// 		assert.Equal(t, "*38:9-38:39", opts.DownloadSections)
// 	})

// 	t.Run("Testing alterDownloadSections with sections arg and with current time in url Youtube", func(t *testing.T) {
// 		url := "https://youtu.be/T_JKIkSf93Y?t=2289"
// 		message := "-s *1:10-2:10 " + url
// 		u := &tgbotapi.Update{Message: &tgbotapi.Message{Text: message}}
// 		opts := &goutubedl.Options{}
// 		alterDownloadSections(u, url, opts)
// 		assert.Equal(t, "*1:10-2:10", opts.DownloadSections)
// 	})

// 	t.Run("Testing alterDownloadSections with current time in url Youtube", func(t *testing.T) {
// 		url := "https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s"
// 		u := &tgbotapi.Update{Message: &tgbotapi.Message{Text: url}}
// 		opts := &goutubedl.Options{}
// 		alterDownloadSections(u, url, opts)
// 		assert.Equal(t, "*3:19-3:49", opts.DownloadSections)
// 		// if opts.DownloadSections != "*3:20-3:49" {

// 		// }
// 	})
// }

func TestAlterDownloadSections(t *testing.T) {
	// Claude 3.5 Sonnet.
	tests := []struct {
		name           string
		inputURL       string
		inputMessage   string
		expectedOutput string
	}{
		{
			name:           "Youtube URL",
			inputURL:       "https://www.youtube.com/watch?v=T_JKIkSf93Y",
			inputMessage:   "https://www.youtube.com/watch?v=T_JKIkSf93Y",
			expectedOutput: "*0:0-0:30",
		},
		{
			name:           "Youtube URL with sections argument",
			inputURL:       "https://www.youtube.com/watch?v=T_JKIkSf93Y",
			inputMessage:   "-s *1:10-2:10 https://www.youtube.com/watch?v=T_JKIkSf93Y",
			expectedOutput: "*1:10-2:10",
		},
		{
			name:           "Youtube URL with current time",
			inputURL:       "https://youtu.be/T_JKIkSf93Y?t=2289",
			inputMessage:   "https://youtu.be/T_JKIkSf93Y?t=2289",
			expectedOutput: "*38:9-38:39",
		},
		{
			name:           "Youtube URL with sections arg and current time",
			inputURL:       "https://youtu.be/T_JKIkSf93Y?t=2289",
			inputMessage:   "-s *1:10-2:10 https://youtu.be/T_JKIkSf93Y?t=2289",
			expectedOutput: "*1:10-2:10",
		},
		{
			name:           "Youtube URL with current time (alternate format)",
			inputURL:       "https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s",
			inputMessage:   "https://www.youtube.com/watch?v=AWVUp12XPpU&t=199s",
			expectedOutput: "*3:19-3:49",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &tgbotapi.Update{Message: &tgbotapi.Message{Text: tt.inputMessage}}
			opts := &goutubedl.Options{}
			alterDownloadSections(u, tt.inputURL, opts)
			assert.Equal(t, tt.expectedOutput, opts.DownloadSections)
		})
	}
}

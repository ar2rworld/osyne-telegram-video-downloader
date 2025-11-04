package utilities_test

import (
	"testing"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/utilities"
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
			got := utilities.ConvertSecondsToMinSec(tc.have)
			if got != tc.want {
				t.Errorf("want %s, but got %s", tc.want, got)
			}
		})
	}
}

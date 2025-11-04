package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
)

func TestSetupCookies(t *testing.T) {
	h := NewHandler(nil, nil, nil, "c", "i", "g", int64(0))
	url := "https://youtube.com/shorts/id"
	opts := &goutubedl.Options{}
	prms := &downloader.Parameters{}
	isYoutubeVideo := false
	h.setupCookies(url, opts, prms, isYoutubeVideo)
	assert.Equal(t, "g", opts.Cookies)
}

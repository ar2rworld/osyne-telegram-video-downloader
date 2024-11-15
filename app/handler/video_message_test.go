package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wader/goutubedl"
)

func TestSetupCookies(t *testing.T) {
	h := NewHandler(nil, nil, "c", "i", "g")
	url := "https://youtube.com/shorts/id"
	opts := &goutubedl.Options{}
	isYoutubeVideo := false
	h.setupCookies(url, opts, isYoutubeVideo)
	assert.Equal(t, "g", opts.Cookies)
}

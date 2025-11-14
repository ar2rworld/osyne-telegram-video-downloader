package platform

import (
	"fmt"

	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
)

type Shorts struct {
	cookiesPath string
}

func NewYoutubeShorts(cookiesPath string) *Shorts {
	return &Shorts{cookiesPath: cookiesPath}
}

// RemuxRequired implements Platform.
func (y *Shorts) RemuxRequired() bool {
	return true
}

func (i *Shorts) RemuxVideoCodec() string {
	return "copy"
}

// ConfigureDownload implements Platform.
func (y *Shorts) ConfigureDownload(_ string, opts *goutubedl.Options) {
	opts.Cookies = y.cookiesPath
}

// Match implements Platform.
func (y *Shorts) Match(url string) bool {
	return match.YoutubeShorts(url) != ""
}

// Name implements Platform.
func (y *Shorts) Name() string {
	return "Shorts"
}

func (y *Shorts) SelectFormat(_ []goutubedl.Format) (string, error) {
	return c.BestFormat, nil
}

func (y *Shorts) MaxDuration(_ *goutubedl.Result) (string, error) {
	return "", fmt.Errorf("%w shouldn't be called", myerrors.ErrPlatform)
}

func (y *Shorts) NeedCut(_ *goutubedl.Result) (bool, error) {
	return false, nil
}

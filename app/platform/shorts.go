package platform

import (
	"fmt"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/wader/goutubedl"
)

type Shorts struct {
	cookiesPath string
}

// RemuxRequired implements Platform.
func (y *Shorts) RemuxRequired() bool {
	return true
}

// ConfigureDownload implements Platform.
func (y *Shorts) ConfigureDownload(url string, opts *goutubedl.Options) error {
	opts.Cookies = y.cookiesPath
	return nil
}

// Match implements Platform.
func (y *Shorts) Match(url string) bool {
	return match.YoutubeShorts(url) != ""
}

// Name implements Platform.
func (y *Shorts) Name() string {
	return "Shorts"
}

func NewYoutubeShorts(cookiesPath string) *Shorts {
	return &Shorts{cookiesPath: cookiesPath}
}

func (y *Shorts) SelectFormat(formats []goutubedl.Format) (format string, err error) {
	return c.BestFormat, nil
}

func (y *Shorts) MaxDuration(r *goutubedl.Result) (string, error) {
	return "", fmt.Errorf("error max duration shouldn't be called")
}

func (y *Shorts) NeedCut(r *goutubedl.Result) (bool, error) {
	return false, nil
}

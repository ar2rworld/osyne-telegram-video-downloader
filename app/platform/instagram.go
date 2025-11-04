package platform

import (
	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
)

type Instagram struct {
	cookiesPath string
}

func NewInstagram(cookiesPath string) *Instagram {
	return &Instagram{cookiesPath: cookiesPath}
}

// RemuxRequired implements Platform.
func (i *Instagram) RemuxRequired() bool {
	return false
}

// NeedCut implements Platform.
func (i *Instagram) NeedCut(*goutubedl.Result) (bool, error) {
	return false, nil
}

// MaxDuration implements Platform.
func (i *Instagram) MaxDuration(*goutubedl.Result) (string, error) {
	return c.DefaultSections, nil
}

// SelectFormat implements Platform.
func (i *Instagram) SelectFormat(formats []goutubedl.Format) (string, error) {
	return c.BestFormat, nil
}

func (i *Instagram) Name() string {
	return "instagram"
}

func (i *Instagram) Match(url string) bool {
	return match.Instagram(url) != ""
}

func (i *Instagram) ConfigureDownload(url string, opts *goutubedl.Options) {
	opts.Cookies = i.cookiesPath
}

package platform

import (
	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
)

type FacebookReels struct{}

func NewFacebookReels() *FacebookReels {
	return &FacebookReels{}
}

// RemuxRequired implements Platform.
func (i *FacebookReels) RemuxRequired() bool {
	return true
}

// NeedCut implements Platform.
func (i *FacebookReels) NeedCut(*goutubedl.Result) (bool, error) {
	return false, nil
}

// MaxDuration implements Platform.
func (i *FacebookReels) MaxDuration(*goutubedl.Result) (string, error) {
	return c.DefaultSections, nil
}

// SelectFormat implements Platform.
func (i *FacebookReels) SelectFormat(formats []goutubedl.Format) (string, error) {
	return c.BestFormat, nil
}

func (i *FacebookReels) Name() string {
	return "FacebookReels"
}

func (i *FacebookReels) Match(url string) bool {
	return match.FacebookReels(url) != ""
}

func (i *FacebookReels) ConfigureDownload(url string, opts *goutubedl.Options) {}

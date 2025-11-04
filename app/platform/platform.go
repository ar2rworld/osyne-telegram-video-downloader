package platform

import (
	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
)

type Platform interface {
	Name() string
	Match(url string) bool
	ConfigureDownload(url string, opts *goutubedl.Options) error
	SelectFormat(formats []goutubedl.Format) (string, error)
	NeedCut(*goutubedl.Result) (bool, error)
	MaxDuration(*goutubedl.Result) (string, error)
	RemuxRequired() bool
}

type DefaultPlatform struct{}

// RemuxRequired implements Platform.
func (i *DefaultPlatform) RemuxRequired() bool {
	return false
}

// NeedCut implements Platform.
func (i *DefaultPlatform) NeedCut(*goutubedl.Result) (bool, error) {
	return false, nil
}

// MaxDuration implements Platform.
func (i *DefaultPlatform) MaxDuration(*goutubedl.Result) (string, error) {
	return c.DefaultSections, nil
}

// SelectFormat implements Platform.
func (i *DefaultPlatform) SelectFormat(formats []goutubedl.Format) (string, error) {
	return c.BestFormat, nil
}

func (i *DefaultPlatform) Name() string {
	return "DefaultPlatform"
}

func (i *DefaultPlatform) Match(url string) bool {
	return true
}

func (i *DefaultPlatform) ConfigureDownload(url string, opts *goutubedl.Options) error {
	return nil
}

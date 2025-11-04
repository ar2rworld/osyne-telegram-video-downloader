package platform

import (
	"github.com/wader/goutubedl"
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

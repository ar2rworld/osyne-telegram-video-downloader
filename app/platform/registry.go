package platform

import c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"

type Registry struct {
	platforms []Platform
}

func NewRegistry(options *c.Options) *Registry {
	registry := &Registry{}
	instagram := NewInstagram(options.InstagramCookiesPath)
	youtube := NewYoutube(options.YouTubeCookiesPath)
	shorts := NewYoutubeShorts(options.YouTubeCookiesPath)
	facebookreels := NewFacebookReels()

	registry.Register(instagram)
	registry.Register(youtube)
	registry.Register(shorts)
	registry.Register(facebookreels)

	return registry
}

func (r *Registry) Register(p Platform) {
	r.platforms = append(r.platforms, p)
}

func (r *Registry) FindPlatform(url string) Platform {
	for _, p := range r.platforms {
		if p.Match(url) {
			return p
		}
	}

	return &DefaultPlatform{}
}

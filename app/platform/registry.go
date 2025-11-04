package platform

type Registry struct {
	platforms []Platform
}

func NewRegistry() *Registry {
	return &Registry{}
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

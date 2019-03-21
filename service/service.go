package service

import "github.com/loivis/prunusavium-go/pavium"

type Service struct {
	lefts  map[pavium.SiteName]pavium.Left
	rights map[pavium.SiteName]pavium.Right
	sites  map[pavium.SiteName]pavium.Site
}

type opts func(*Service)

func New(opts ...opts) *Service {
	svc := &Service{
		lefts:  make(map[pavium.SiteName]pavium.Left),
		rights: make(map[pavium.SiteName]pavium.Right),
		sites:  make(map[pavium.SiteName]pavium.Site),
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func WithLefts(lefts map[pavium.SiteName]pavium.Left) opts {
	return func(svc *Service) {
		svc.lefts = lefts
		for name, left := range lefts {
			if site, ok := left.(pavium.Site); ok {
				svc.sites[name] = site
			}
		}
	}
}

func WithRights(rights map[pavium.SiteName]pavium.Right) opts {
	return func(svc *Service) {
		svc.rights = rights
		for name, right := range rights {
			if site, ok := right.(pavium.Site); ok {
				svc.sites[name] = site
			}
		}
	}
}

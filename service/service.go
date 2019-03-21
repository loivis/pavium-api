package service

import "github.com/loivis/prunusavium-go/pavium"

type Service struct {
	sites map[pavium.SiteName]pavium.Site
}

func New(sites map[pavium.SiteName]pavium.Site) *Service {
	return &Service{
		sites: sites,
	}
}

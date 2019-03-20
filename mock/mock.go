package mock

import "github.com/loivis/prunusavium-go/pavium"

type Site struct {
	ChaptersFunc func(string) []pavium.Chapter
}

func (s *Site) Chapters(link string) []pavium.Chapter {
	return s.ChaptersFunc(link)
}

package mock

import "github.com/loivis/prunusavium-go/pavium"

type Site struct {
	ChaptersFunc func(string) []pavium.Chapter
	TextFunc     func(string) string
}

func (s *Site) Chapters(link string) []pavium.Chapter {
	return s.ChaptersFunc(link)
}

func (s *Site) Text(link string) string {
	return s.TextFunc(link)
}

package mock

import "github.com/loivis/pavium-api/pavium"

type Site struct {
	NameFunc           func() string
	ChaptersFunc       func(string) []pavium.Chapter
	TextFunc           func(string) string
	SearchBookFunc     func(author, title string) (pavium.Book, error)
	SearchKeywordsFunc func(string) []pavium.Book
}

func (s *Site) Name() string {
	return s.NameFunc()
}

func (s *Site) Chapters(link string) []pavium.Chapter {
	return s.ChaptersFunc(link)
}

func (s *Site) Text(link string) string {
	return s.TextFunc(link)
}

func (s *Site) SearchBook(author, title string) (pavium.Book, error) {
	return s.SearchBookFunc(author, title)
}

func (s *Site) SearchKeywords(keywords string) []pavium.Book {
	return s.SearchKeywordsFunc(keywords)
}

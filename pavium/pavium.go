package pavium

import "time"

type Book struct {
	Author string    `json:"author,omitempty"`
	ID     string    `json:"id,omitempty"`
	Site   string    `json:"site,omitempty"`
	Title  string    `json:"title,omitempty"`
	Update time.Time `json:"update,omitempty"`
}

type Chapter struct {
	Name string `json:"name,omitempty"`
	Link string `json:"link,omitempty"`
}

type SiteName string

const (
	Piaotian SiteName = "飘天文学网"
)

type Site interface {
	Name() string
	Chapters(link string) []Chapter
	Text(link string) (text string)
	SearchBook(author, title string) (Book, error)
}

type Left interface {
}

type Right interface {
	Name() string
	SearchKeywords(keywords string) []Book
}

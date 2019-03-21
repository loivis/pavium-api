package pavium

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

type Left interface {
}

type Right interface {
	Name() string
	SearchKeywords(keywords string) []Book
}
type Site interface {
	Name() string
	SearchBook(author, title string) (Book, error)
	Chapters(link string) []Chapter
	Text(link string) (text string)
}

type Store interface {
	Delete(ctx context.Context, fav *Favorite) error
	List(ctx context.Context) []Favorite
	Put(ctx context.Context, fav *Favorite) error
}

type Book struct {
	Author      string     `json:"author,omitempty"`
	ChapterLink string     `json:"chapterLink,omitempty"`
	ID          string     `json:"id,omitempty"`
	Site        string     `json:"site,omitempty"`
	Title       string     `json:"title,omitempty"`
	Update      *time.Time `json:"update,omitempty"`
}

type Chapter struct {
	Name string `json:"name,omitempty"`
	Link string `json:"link,omitempty"`
}

type Favorite struct {
	Author string `json:"author,omitempty" firestore:"author"`
	BookID string `json:"bookID,omitempty" firestore:"bookID"`
	Site   string `json:"site,omitempty" firestore:"site"`
	Title  string `json:"title,omitempty" firestore:"title"`
}

func (fav *Favorite) ID() string {
	s := fmt.Sprintf("%s-%s-%s", fav.Site, fav.Author, fav.Title)
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:8])
}

type SiteName string

const (
	Piaotian SiteName = "飘天文学网"
)

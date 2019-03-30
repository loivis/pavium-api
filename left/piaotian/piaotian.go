package piaotian

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/loivis/prunusavium-utils/http"
	"github.com/loivis/prunusavium-api/pavium"
)

type Site struct {
	name       string
	home       string
	chapterURL string
}

func New() *Site {
	return &Site{
		name:       string(pavium.Piaotian),
		home:       "https://www.ptwxz.com/",
		chapterURL: "https://www.ptwxz.com/html/",
	}
}

func (s *Site) Name() string {
	return s.name
}

func (s *Site) SearchBook(author, title string) (pavium.Book, error) {
	book := pavium.Book{Author: author, Site: s.name, Title: title}

	q := fmt.Sprintf("%s %s site:%s", author, title, s.home)
	link, err := http.Search(q)
	if err != nil {
		return book, err
	}

	link = s.parseChapterLink(link)

	if link == "" {
		return book, fmt.Errorf("%v not found", book)
	}

	book.ChapterLink = link

	return book, nil
}

func (s *Site) Chapters(link string) []pavium.Chapter {
	chapters := []pavium.Chapter{}

	doc, err := http.GetDoc(link)
	if err != nil {
		log.Println(err)
		return chapters
	}

	doc.Find("div.centent").Find("li").Each(func(i int, sel *goquery.Selection) {
		chapter := pavium.Chapter{}

		name := sel.Text()
		switch name {
		case "分享到twitter", "分享到facebook", "分享到Google+", "\u00a0":
			return
		default:
			chapter.Name = name
		}

		a := sel.Find("a")
		if href, ok := a.Attr("href"); ok {
			chapter.Link = link + href
		}

		if chapter.Name != "" && chapter.Link != "" {
			chapters = append(chapters, chapter)
		}
	})

	return chapters
}

func (s *Site) Text(link string) string {
	doc, err := http.GetDoc(link)
	if err != nil {
		log.Println(err)
		return ""
	}

	replacer := strings.NewReplacer(
		"<br/>", "\n",
		"\u00a0", "",
		"\b", "",
		"\t", "",
	)

	html, _ := doc.Html()
	html = replacer.Replace(html)
	ndoc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	text := ndoc.Find("body").Children().Remove().End().Text()

	return strings.TrimSpace(text)
}

func (s *Site) parseChapterLink(link string) string {
	if link == "" {
		return ""
	}

	link = strings.Trim(link, ".html")
	ss := strings.Split(link, "/")

	if len(ss) < 6 || (len(ss) == 6 && ss[5] == "") {
		return ""
	}

	return s.chapterURL + ss[4] + "/" + ss[5] + "/"
}

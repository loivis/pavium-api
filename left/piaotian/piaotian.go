package piaotian

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/loivis/convolvulus-utils/http"
	"github.com/loivis/prunusavium-go/pavium"
)

type Site struct{}

func New() *Site {
	return &Site{}
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

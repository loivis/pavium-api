package qidian

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/loivis/prunusavium-api/pavium"
	"github.com/loivis/prunusavium-utils/http"
)

type Site struct {
	name       string
	home       string
	bookURL    string
	chapterURL string
	searchURL  string
}

func New() *Site {
	return &Site{
		name:       string(pavium.Qidian),
		home:       "https://book.qidian.com",
		bookURL:    "https://book.qidian.com/info/%s",
		chapterURL: "https://book.qidian.com/info/%s#Catalog",
		searchURL:  "https://www.qidian.com/search?kw=%s",
	}
}

func (s *Site) Name() string {
	return s.name
}

func (s *Site) SearchKeywords(keywords string) []pavium.Book {
	url := fmt.Sprintf(s.searchURL, url.QueryEscape(keywords))

	doc, err := http.GetDoc(url)
	if err != nil {
		log.Printf("failed to search %q at %q: %v", keywords, s.name, err)
		return nil
	}

	return s.findBooks(doc)
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

	doc.Find("div.volume-wrap").Find("li").Each(func(i int, sel *goquery.Selection) {
		a := sel.Find("a")

		chapter := pavium.Chapter{
			Name: a.Text(),
		}

		if href, ok := a.Attr("href"); ok {
			chapter.Link = "https:" + href
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
		"<p>　　", "",
		"</p>", "\n",
	)

	html, _ := doc.Find("div.read-content").Html()
	html = replacer.Replace(html)
	ndoc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	text := ndoc.Text()

	if strings.Contains(link, "vipreader") {
		text = "###VIP章节### 需要订阅后才能阅读全文\n\n" + text
	}

	return text
}

func (s *Site) findBooks(doc *goquery.Document) (books []pavium.Book) {
	doc.Find("div#result-list").Find("li.res-book-item").Each(func(i int, sel *goquery.Selection) {
		a := sel.Find("div.book-mid-info").Find("h4").Find("a")
		title := a.Text()
		href, _ := a.Attr("href")
		id := strings.Split(href, "/")[4]
		author := sel.Find("div.book-mid-info").Find("a.name").Text()

		books = append(books, pavium.Book{
			Author: author,
			ID:     id,
			Link:   "https:" + href,
			Site:   s.name,
			Title:  title,
		})
	})

	if len(books) == 0 {
		return books
	}

	var wg sync.WaitGroup
	wg.Add(len(books))

	for _, book := range books {
		go func(book *pavium.Book) {
			s.lastUpdate(book)
			wg.Done()
		}(&book)
	}

	wg.Wait()

	return books
}

func (s *Site) lastUpdate(b *pavium.Book) {
	url := fmt.Sprintf(s.chapterURL, b.ID)

	doc, err := http.GetDoc(url)
	if err != nil {
		log.Printf("failed to fetch content for %q at %q: %v", b.Title, s.name, err)
		return
	}

	found := doc.Find("div.book-info").Find("h1").Find("em").Text()
	if b.Title != found {
		log.Printf("%q not found at %v, found %q", b.Title, s.name, found)
		return
	}

	// TODO: it doesn't manage to load all chapters if too many (like 1000+)
	// 首发时间：2017-04-25 15:25:00 章节字数：4227
	sel := doc.Find("div.volume-wrap").Find("li").Last().Find("a")
	re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}")
	update, _ := sel.Attr("title")
	updates := re.FindStringSubmatch(update)
	if len(updates) > 0 {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", updates[0], loc)
		utc := t.UTC()
		b.Update = &utc
	}
}

func (s *Site) parseChapterLink(link string) string {
	if !strings.HasPrefix(link, s.home+"/info/") {
		log.Printf("%s: invalid book info page %q", s.name, link)
		return ""
	}

	return link + "#Catalog"
}

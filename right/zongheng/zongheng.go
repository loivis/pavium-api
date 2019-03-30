package zongheng

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
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
	bookLink   string
	chapterURL string
	searchURL  string
}

func New() *Site {
	return &Site{
		name:       string(pavium.Zongheng),
		home:       "http://book.zongheng.com",
		bookLink:   "http://book.zongheng.com/book/%s.html",
		chapterURL: "http://book.zongheng.com/showchapter/%v.html",
		searchURL:  "http://search.zongheng.com/s?keyword=%s",
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

	book.ChapterLink = s.parseChapterLink(link)

	return book, nil
}

func (s *Site) Chapters(link string) []pavium.Chapter {
	chapters := []pavium.Chapter{}

	doc, err := http.GetDoc(link)
	if err != nil {
		log.Println(err)
		return chapters
	}

	doc.Find("div.volume-list").Find("li").Each(func(i int, sel *goquery.Selection) {
		a := sel.Find("a")

		chapter := pavium.Chapter{
			Name: a.Text(),
		}

		if href, ok := a.Attr("href"); ok {
			chapter.Link = href
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

	html, _ := doc.Find("div.content").Html()
	html = replacer.Replace(html)
	ndoc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	text := ndoc.Text()

	if doc.Find("div.reader_order").HasClass("reader_order") {
		text = "###VIP章节### 需要订阅后才能阅读全文\n\n" + text
	}

	return text
}

func (s *Site) findBooks(doc *goquery.Document) (books []pavium.Book) {
	doc.Find("div.search-result-list").Each(func(i int, sel *goquery.Selection) {
		title := sel.Find("h2.tit").Find("a").Text()
		author := sel.Find("div.bookinfo").Find("a:first-of-type").Text()
		// image, _ := sel.Find("div.imgbox").Find("img").Attr("src")
		// intro := strings.TrimSpace(sel.Find("p").Text())
		link, _ := sel.Find("h2.tit").Find("a").Attr("href")
		id := strings.Split((strings.Split(link, "/")[4]), ".")[0]

		books = append(books, pavium.Book{
			Author: author,
			ID:     id,
			Link:   link,
			Site:   s.name,
			Title:  title,
		})
	})

	if len(books) == 0 {
		return
	}

	// return only the first 10 results
	if len(books) > 10 {
		books = books[:10]
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

	return
}

func (s *Site) lastUpdate(b *pavium.Book) (t time.Time) {
	url := fmt.Sprintf(s.chapterURL, b.ID)

	doc, err := http.GetDoc(url)
	if err != nil {
		log.Printf("failed to fetch content for %q at %q: %v", b.Title, s.name, err)
		return
	}

	pageTitle := doc.Find("head").Find("title").Text()

	switch {
	case strings.Contains(pageTitle, s.name):
		t = lastUpdateAtZongheng(b, s.name, doc)
	case strings.Contains(pageTitle, "花语女生网"):
		site := "花语女生网"
		url := fmt.Sprintf("http://huayu.baidu.com/showchapter/%s.html", b.ID)
		doc, err := http.GetDoc(url)
		if err != nil {
			log.Printf("failed to fetch content for %q at %q: %v", b.Title, site, err)
			return
		}
		t = lastUpdateAtHuayu(b, site, doc)
	default:
		log.Printf("unknown site for %q: %q", b.Title, pageTitle)
	}

	return
}

func lastUpdateAtZongheng(b *pavium.Book, site string, doc *goquery.Document) (t time.Time) {
	found := doc.Find("div.book-meta").Find("h1").Text()
	if b.Title != found {
		log.Printf("%q not found at %v, found %q", b.Title, site, found)
		return
	}

	// 第一章：恐怖初始 字数：3462 更新时间：2018-10-10 15:43
	sel := doc.Find("ul.chapter-list").Find("li.col-4").Last().Find("a")
	re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}")
	update, _ := sel.Attr("title")
	updates := re.FindStringSubmatch(update)
	if len(updates) > 0 {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		t, _ = time.ParseInLocation("2006-01-02 15:04", updates[0], loc)
		utc := t.UTC()
		b.Update = &utc
	}

	return
}

func lastUpdateAtHuayu(b *pavium.Book, site string, doc *goquery.Document) (t time.Time) {
	found := doc.Find("div.book_title").Find("h1").Children().Remove().End().Text()

	if b.Title != found {
		log.Printf("%q not found at %v, found %q\n", b.Title, site, found)
		return
	}

	// 2018-08-04 17:38:01
	update := doc.Find("div.book_chapter").Find("span.chaptime").Last().Text()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ = time.ParseInLocation("2006-01-02 15:04:05", update, loc)
	utc := t.UTC()
	b.Update = &utc

	return
}

func (s *Site) parseChapterLink(link string) string {
	// http://book.zongheng.com/book/685640.html?fr=pc_alading
	prefix := "http://book.zongheng.com/book/"

	if !strings.HasPrefix(link, prefix) {
		return ""
	}

	ss := strings.Split(strings.TrimPrefix(link, "http://book.zongheng.com/book/"), ".")

	if len(ss) <= 1 {
		return ""
	}

	id, err := strconv.Atoi(ss[0])
	if err != nil {
		log.Printf("%q doesn't contain a valid book id: %v", link, err)
		return ""
	}

	return fmt.Sprintf(s.chapterURL, id)
}

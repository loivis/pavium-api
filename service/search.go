package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/loivis/pavium-api/pavium"
)

func (svc *Service) Search(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		svc.getSearch(w, r)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (svc *Service) getSearch(w http.ResponseWriter, r *http.Request) {
	var books []pavium.Book

	q := r.URL.Query()
	author := q.Get("author")
	title := q.Get("title")
	keywords := q.Get("keywords")

	switch {
	case author != "" && title != "" && keywords == "":
		books = svc.searchBook(author, title)
	case keywords != "" && author == "" && title == "":
		books = svc.searchKeywords(keywords)
	default:
		http.Error(w, "missing required parameters: author/title or keywords", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(books)
	return
}

func (svc *Service) searchBook(author, title string) []pavium.Book {
	ch := make(chan *pavium.Book, len(svc.sites))
	books := []pavium.Book{}

	var wg sync.WaitGroup
	wg.Add(len(svc.sites))

	go func() {
		wg.Wait()
		close(ch)
	}()

	for _, site := range svc.sites {
		go func(site pavium.Site) {
			defer wg.Done()

			tick := time.Now()

			book, err := site.SearchBook(author, title)
			if err != nil {
				log.Println(err)
				return
			}

			ch <- &book
			log.Printf("%v: fetched %q by %q in %v\n", site.Name(), title, author, time.Since(tick))
		}(site)
	}

	for book := range ch {
		books = append(books, *book)
	}

	return books
}

func (svc *Service) searchKeywords(keywords string) []pavium.Book {
	ch := make(chan []pavium.Book, len(svc.rights))
	books := []pavium.Book{}

	var wg sync.WaitGroup
	wg.Add(len(svc.rights))

	go func() {
		wg.Wait()
		close(ch)
	}()

	for _, right := range svc.rights {
		go func(right pavium.Right) {
			defer wg.Done()

			tick := time.Now()

			ret := right.SearchKeywords(keywords)
			ch <- ret
			log.Printf("%v: %d results on %q in %v\n", right.Name(), len(ret), keywords, time.Since(tick))
		}(right)
	}

	for slice := range ch {
		books = append(books, slice...)
	}

	return books
}

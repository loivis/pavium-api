package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loivis/prunusavium-go/pavium"

	"github.com/loivis/prunusavium-go/mock"
)

func TestService_Search(t *testing.T) {
	t.Run("WrongMethod", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/foo", nil)

		New().Search(w, r)

		if got, want := w.Code, 405; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})

}

func TestService_GetSearch(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		lefts := map[pavium.SiteName]pavium.Left{
			"left": &mock.Site{
				NameFunc: func() string { return "left" },
				SearchBookFunc: func(author, title string) (pavium.Book, error) {
					return pavium.Book{Title: "left book"}, nil
				},
			}}
		rights := map[pavium.SiteName]pavium.Right{
			"right": &mock.Site{
				NameFunc: func() string { return "right" },
				SearchBookFunc: func(author, title string) (pavium.Book, error) {
					return pavium.Book{Title: "right book"}, nil
				},
				SearchKeywordsFunc: func(keywords string) []pavium.Book {
					return []pavium.Book{
						{Title: "right keywords foo"},
						{Title: "right keywords bar"},
					}
				},
			}}

		for _, tc := range []struct {
			desc string
			path string
		}{
			{
				desc: "WithAuthorAndTitle",
				path: "/search?author=foo&title=bar",
			},
			{
				desc: "WithKeywords",
				path: "/search?keywords=foo",
			},
		} {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tc.path, nil)

			svc := New(
				WithLefts(lefts),
				WithRights(rights),
			)
			svc.getSearch(w, r)

			if got, want := w.Code, 200; got != want {
				t.Fatalf("[%s] w.Code = %d, want %d", tc.desc, got, want)
			}

			if got, want := w.Header().Get("Content-Type"), "application/json; charset=utf-8"; got != want {
				t.Errorf("[%s] %q, %q", tc.desc, got, want)
			}

			var books []pavium.Book
			if err := json.NewDecoder(w.Body).Decode(&books); err != nil {
				t.Fatalf("[%s] failed to encode response: %q", tc.desc, err)
			}

			if got, want := len(books), 2; got != want {
				t.Fatalf("[%s] got response = %v, want %v", tc.desc, got, want)
			}
		}
	})

	t.Run("MissingParameters", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)

		New().getSearch(w, r)

		if got, want := w.Code, 400; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})
}

func TestService_SearchBook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		lefts := map[pavium.SiteName]pavium.Left{
			"l1": &mock.Site{
				NameFunc: func() string { return "l1" },
				SearchBookFunc: func(author, title string) (pavium.Book, error) {
					return pavium.Book{Title: "left foo"}, nil
				},
			},
			"l2": &mock.Site{
				NameFunc: func() string { return "l2" },
				SearchBookFunc: func(author, title string) (pavium.Book, error) {
					return pavium.Book{}, errors.New("error foo")
				},
			},
		}

		rights := map[pavium.SiteName]pavium.Right{
			"r1": &mock.Site{
				NameFunc: func() string { return "r1" },
				SearchBookFunc: func(author, title string) (pavium.Book, error) {
					return pavium.Book{Title: "right foo"}, nil
				},
			},
			"r2": &mock.Site{
				NameFunc: func() string { return "r2" },
				SearchBookFunc: func(author, title string) (pavium.Book, error) {
					return pavium.Book{Title: "right bar"}, nil
				},
			},
			"r3": &mock.Site{
				NameFunc: func() string { return "r3" },
				SearchBookFunc: func(author, title string) (pavium.Book, error) {
					return pavium.Book{}, errors.New("error foo")
				},
			},
		}

		svc := New(
			WithLefts(lefts),
			WithRights(rights),
		)

		books := svc.searchBook("foo", "bar")

		if got, want := len(books), 3; got != want {
			t.Fatalf("got response = %v, want %v", got, want)
		}
	})
}

func TestService_SearchKeywords(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		rights := map[pavium.SiteName]pavium.Right{
			"r1": &mock.Site{
				NameFunc: func() string { return "r1" },
				SearchKeywordsFunc: func(keywords string) []pavium.Book {
					return []pavium.Book{
						{Title: "keywords foo"},
						{Title: "keywords bar"},
					}
				},
			},

			"r2": &mock.Site{
				NameFunc: func() string { return "r2" },
				SearchKeywordsFunc: func(keywords string) []pavium.Book {
					return []pavium.Book{
						{Title: "keywords baz"},
						{Title: "keywords qux"},
					}
				},
			},
		}

		svc := New(
			WithRights(rights),
		)

		books := svc.searchKeywords("foo")

		if got, want := len(books), 4; got != want {
			t.Fatalf("got response = %v, want %v", got, want)
		}
	})
}

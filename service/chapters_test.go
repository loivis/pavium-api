package service

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loivis/pavium-api/mock"
	"github.com/loivis/pavium-api/pavium"
)

func TestService_Chapters(t *testing.T) {
	t.Run("WrongMethod", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodTrace, "/foo", nil)

		New().Chapters(w, r)

		if got, want := w.Code, http.StatusMethodNotAllowed; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})
}

func TestService_GetChapters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo?site=foo&link=https%3A%2F%2Ffoo.org", nil)

		wantChapters := []pavium.Chapter{
			{
				Name: "bar",
				Link: "https://foo.org/bar",
			},
			{
				Name: "baz",
				Link: "https://foo.org/baz",
			},
		}

		site := &mock.Site{
			ChaptersFunc: func(string) []pavium.Chapter {
				return wantChapters
			},
		}
		lefts := map[pavium.SiteName]pavium.Left{
			"foo": site,
		}

		svc := New(
			WithLefts(lefts),
		)

		log.Printf("svc: %+v", svc)

		svc.getChapters(w, r)

		if got, want := w.Code, http.StatusOK; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}

		if got, want := w.Header().Get("Content-Type"), "application/json; charset=utf-8"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}

		var gotChapters []pavium.Chapter

		if err := json.NewDecoder(w.Body).Decode(&gotChapters); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := len(gotChapters), len(wantChapters); got != want {
			t.Fatalf("got chapters = %v, want %v", gotChapters, wantChapters)
		}

		for n := range wantChapters {
			if got, want := gotChapters[n], wantChapters[n]; got != want {
				t.Errorf("gotChapter[%d] = %q, want %q", n, got, want)
			}
		}
	})

	t.Run("MissingParameterSite", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo?link=https%3A%2F%2Ffoo.org", nil)
		svc := New()

		svc.getChapters(w, r)

		if got, want := w.Code, http.StatusBadRequest; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})

	t.Run("MissingParameterLink", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo?site=bar", nil)
		svc := New()

		svc.getChapters(w, r)

		if got, want := w.Code, http.StatusBadRequest; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})

	t.Run("UnknownSite", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo?site=bar&link=https%3A%2F%2Ffoo.org", nil)
		svc := New()

		svc.getChapters(w, r)

		if got, want := w.Code, http.StatusNoContent; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})
}

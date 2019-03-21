package service

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loivis/prunusavium-go/mock"
	"github.com/loivis/prunusavium-go/pavium"
)

func TestService_Text(t *testing.T) {
	t.Run("WrongMethod", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodTrace, "/foo", nil)

		New(nil).Text(w, r)

		if got, want := w.Code, http.StatusMethodNotAllowed; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})
}

func TestService_GetText(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo?site=foo&link=https%3A%2F%2Ffoo.org", nil)

		wantText := "keep quiet"

		site := &mock.Site{
			TextFunc: func(string) string {
				return wantText
			},
		}

		svc := New(map[pavium.SiteName]pavium.Site{
			"foo": site,
		})

		svc.getText(w, r)

		if got, want := w.Code, http.StatusOK; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}

		if got, want := w.Header().Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}

		var gotText string
		b, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := string(b), wantText; got != want {
			t.Fatalf("got text = %v, want %v", gotText, wantText)
		}
	})

	t.Run("MissingParameterSite", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo?link=https%3A%2F%2Ffoo.org", nil)
		svc := New(nil)

		svc.getText(w, r)

		if got, want := w.Code, http.StatusBadRequest; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})

	t.Run("MissingParameterLink", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo?site=bar", nil)
		svc := New(nil)

		svc.getText(w, r)

		if got, want := w.Code, http.StatusBadRequest; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})

	t.Run("UnknownSite", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo?site=bar&link=https%3A%2F%2Ffoo.org", nil)
		svc := New(nil)

		svc.getText(w, r)

		if got, want := w.Code, http.StatusNoContent; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})
}

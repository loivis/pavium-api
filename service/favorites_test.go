package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/loivis/prunusavium-go/mock"
	"github.com/loivis/prunusavium-go/pavium"
)

func TestService_Favorites(t *testing.T) {
	t.Run("WrongMethod", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodTrace, "/foo", nil)

		New().Favorites(w, r)

		if got, want := w.Code, http.StatusMethodNotAllowed; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})
}

func TestService_GetFavorites(t *testing.T) {
	store := &mock.Store{
		ListFunc: func(ctx context.Context) []pavium.Favorite {
			return []pavium.Favorite{
				{Title: "foo"},
				{Title: "bar"},
			}
		},
	}

	svc := New(WithFavStore(store))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo", nil)

	svc.getFavorites(w, r)

	if got, want := w.Code, http.StatusOK; got != want {
		t.Fatalf("w.Code = %d, want %d", got, want)
	}

	if got, want := w.Header().Get("Content-Type"), "application/json; charset=utf-8"; got != want {
		t.Errorf("w.Header.Content-Type = %q, want %q", got, want)
	}

	var res []pavium.Favorite

	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("unexpected error: %q", err)
	}

	if got, want := len(res), 2; got != want {
		t.Fatalf("got %d favorites, want %d", got, want)
	}
}

func TestService_PostFavorites(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		store := &mock.Store{
			PutFunc: func(context.Context, *pavium.Favorite) error { return nil },
		}
		svc := New(WithFavStore(store))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/foo", bytes.NewReader([]byte(`
		{
			"author": "foo",
			"bookID": "123",
			"site": "bar",
			"title": "baz"
		}`)))

		svc.postFavorites(w, r)

		if got, want := w.Code, http.StatusOK; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})

	t.Run("InvalidJSONPayload", func(t *testing.T) {
		store := &mock.Store{}
		svc := New(WithFavStore(store))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader("foo"))

		svc.postFavorites(w, r)

		if got, want := w.Code, http.StatusBadRequest; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}

		if got, want := w.Body.String(), "invalid json payload\n"; got != want {
			t.Fatalf("w.Body = %q, want %q", got, want)
		}
	})

	t.Run("InvalidFavoritePayload", func(t *testing.T) {
		store := &mock.Store{}
		svc := New(WithFavStore(store))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/foo", bytes.NewReader([]byte(`{}`)))

		svc.postFavorites(w, r)

		if got, want := w.Code, http.StatusBadRequest; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}

		if got, want := w.Body.String(), "invalid favorite payload\n"; got != want {
			t.Fatalf("w.Body = %q, want %q", got, want)
		}
	})

	t.Run("ErrorFromStore", func(t *testing.T) {
		store := &mock.Store{
			PutFunc: func(context.Context, *pavium.Favorite) error { return errors.New("foo") },
		}
		svc := New(WithFavStore(store))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/foo", bytes.NewReader([]byte(`
		{
			"author": "foo",
			"bookID": "123",
			"site": "bar",
			"title": "baz"
		}`)))

		svc.postFavorites(w, r)

		if got, want := w.Code, http.StatusInternalServerError; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})
}

func TestService_DeleteFavorites(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		store := &mock.Store{
			DeleteFunc: func(context.Context, *pavium.Favorite) error { return nil },
		}
		svc := New(WithFavStore(store))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/foo", bytes.NewReader([]byte(`
		{
			"author": "foo",
			"site": "bar",
			"title": "baz"
		}`)))

		svc.deleteFavorites(w, r)

		if got, want := w.Code, http.StatusOK; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})

	t.Run("InvalidJSONPayload", func(t *testing.T) {
		store := &mock.Store{}
		svc := New(WithFavStore(store))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/foo", strings.NewReader("foo"))

		svc.deleteFavorites(w, r)

		if got, want := w.Code, http.StatusBadRequest; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}

		if got, want := w.Body.String(), "invalid json payload\n"; got != want {
			t.Fatalf("w.Body = %q, want %q", got, want)
		}
	})

	t.Run("InvalidFavoritePayload", func(t *testing.T) {
		store := &mock.Store{}
		svc := New(WithFavStore(store))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/foo", bytes.NewReader([]byte(`{}`)))

		svc.deleteFavorites(w, r)

		if got, want := w.Code, http.StatusBadRequest; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}

		if got, want := w.Body.String(), "invalid favorite payload\n"; got != want {
			t.Fatalf("w.Body = %q, want %q", got, want)
		}
	})

	t.Run("ErrorFromStore", func(t *testing.T) {
		store := &mock.Store{
			DeleteFunc: func(context.Context, *pavium.Favorite) error { return errors.New("foo") },
		}
		svc := New(WithFavStore(store))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/foo", bytes.NewReader([]byte(`
		{
			"author": "foo",
			"site": "bar",
			"title": "baz"
		}`)))

		svc.deleteFavorites(w, r)

		if got, want := w.Code, http.StatusInternalServerError; got != want {
			t.Fatalf("w.Code = %d, want %d", got, want)
		}
	})
}

package service

import (
	"encoding/json"
	"net/http"

	"github.com/loivis/prunusavium-api/pavium"
)

func (svc *Service) Favorites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		svc.getFavorites(w, r)
	case http.MethodPost:
		svc.postFavorites(w, r)
	case http.MethodDelete:
		svc.deleteFavorites(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (svc *Service) getFavorites(w http.ResponseWriter, r *http.Request) {
	favs := svc.favStore.List(r.Context())

	_ = json.NewEncoder(w).Encode(favs)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (svc *Service) postFavorites(w http.ResponseWriter, r *http.Request) {
	var fav *pavium.Favorite

	if err := json.NewDecoder(r.Body).Decode(&fav); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}

	// TODO: idiomatic way?
	if fav.Author == "" || fav.Site == "" || fav.Title == "" || fav.BookID == "" {
		http.Error(w, "invalid favorite payload", http.StatusBadRequest)
		return
	}

	err := svc.favStore.Put(r.Context(), fav)

	if err != nil {
		http.Error(w, "failed to add/update favorite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (svc *Service) deleteFavorites(w http.ResponseWriter, r *http.Request) {
	var fav *pavium.Favorite

	if err := json.NewDecoder(r.Body).Decode(&fav); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}

	// TODO: idiomatic way?
	if fav.Author == "" || fav.Site == "" || fav.Title == "" {
		http.Error(w, "invalid favorite payload", http.StatusBadRequest)
		return
	}

	err := svc.favStore.Delete(r.Context(), fav)

	if err != nil {
		http.Error(w, "failed to delete favorite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

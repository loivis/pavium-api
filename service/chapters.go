package service

import (
	"encoding/json"
	"net/http"

	"github.com/loivis/pavium-api/pavium"
)

func (svc *Service) Chapters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		svc.getChapters(w, r)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (svc *Service) getChapters(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	st := q.Get("site")
	lk := q.Get("link")
	if st == "" || lk == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	site, ok := svc.sites[pavium.SiteName(st)]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}

	links := site.Chapters(lk)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(links)
}

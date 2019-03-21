package service

import (
	"net/http"

	"github.com/loivis/prunusavium-go/pavium"
)

func (svc *Service) Text(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		svc.getText(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (svc *Service) getText(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte(site.Text(lk)))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

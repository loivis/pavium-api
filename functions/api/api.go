package api

import (
	"net/http"

	"github.com/loivis/prunusavium-go/left/piaotian"
	"github.com/loivis/prunusavium-go/pavium"
	"github.com/loivis/prunusavium-go/service"
)

var mux = http.NewServeMux()
var svc = setupService()

func init() {
	mux.HandleFunc("/search", http.HandlerFunc(svc.Search))
	mux.HandleFunc("/chapters", http.HandlerFunc(svc.Chapters))
	mux.HandleFunc("/text", http.HandlerFunc(svc.Text))
}

func V1(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}

func setupService() *service.Service {
	lefts := map[pavium.SiteName]pavium.Left{
		pavium.Piaotian: piaotian.New(),
	}

	return service.New(
		service.WithLefts(lefts),
	)
}

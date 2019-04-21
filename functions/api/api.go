package api

import (
	"net/http"
	"os"

	"github.com/loivis/pavium-api/left/piaotian"
	"github.com/loivis/pavium-api/pavium"
	"github.com/loivis/pavium-api/right/qidian"
	"github.com/loivis/pavium-api/right/zongheng"
	"github.com/loivis/pavium-api/service"
	"github.com/loivis/pavium-api/store"
)

var mux = http.NewServeMux()
var svc = setupService()

func init() {
	mux.HandleFunc("/search", http.HandlerFunc(svc.Search))
	mux.HandleFunc("/favorites", http.HandlerFunc(svc.Favorites))
	mux.HandleFunc("/chapters", http.HandlerFunc(svc.Chapters))
	mux.HandleFunc("/text", http.HandlerFunc(svc.Text))
}

func V1(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}

func setupService() *service.Service {
	project := os.Getenv("GCP_PROJECT")

	lefts := map[pavium.SiteName]pavium.Left{
		pavium.Piaotian: piaotian.New(),
	}
	rights := map[pavium.SiteName]pavium.Right{
		pavium.Qidian:   qidian.New(),
		pavium.Zongheng: zongheng.New(),
	}
	favStore := store.New(project, "favorites")

	return service.New(
		service.WithLefts(lefts),
		service.WithRights(rights),
		service.WithFavStore(favStore),
	)
}

package api

import (
	"net/http"

	"github.com/loivis/prunusavium-go/service"
)

var mux = http.NewServeMux()
var svc = setupService()

func init() {
	mux.HandleFunc("/chapters", http.HandlerFunc(svc.Chapters))
}

func V1(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}

func setupService() *service.Service {
	return service.New(nil)
}

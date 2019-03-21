package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/loivis/prunusavium-go/functions/api"
	"github.com/loivis/prunusavium-go/pavium"
)

func main() {
	ts := httptest.NewServer(http.HandlerFunc(api.V1))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/chapters?site=飘天文学网&link=https://www.ptwxz.com/html/9/9189")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var data []pavium.Chapter

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
	}

	log.Println(data)
}

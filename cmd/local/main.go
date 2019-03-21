package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/loivis/prunusavium-go/functions/api"
)

func main() {
	ts := httptest.NewServer(http.HandlerFunc(api.V1))
	defer ts.Close()

	piaotianSearchBook := ts.URL + "/search?author=贪睡的龙&title=万世为王"
	// piaotianChapters := ts.URL + "/chapters?site=飘天文学网&link=https://www.ptwxz.com/html/9/9189"
	// piaotianText := ts.URL + "/text?site=飘天文学网&link=https://www.ptwxz.com/html/9/9189/6162982.html"

	get(piaotianSearchBook)
}

func get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response body:\n%s\n", raw)
}

package main

import (
	"log"
	"net/http"

	"github.com/loivis/pavium-api/functions/api"
)

func main() {
	s := http.Server{
		Addr:    ":1234",
		Handler: http.HandlerFunc(api.V1),
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

// /search?author=贪睡的龙&title=万世为王
// /chapters?site=飘天文学网&link=https://www.ptwxz.com/html/9/9189
// /text?site=飘天文学网&link=https://www.ptwxz.com/html/9/9189/6162982.html

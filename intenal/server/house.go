package server

import (
	"log"
	"net/http"
	"net/url"
	"simple2gis/intenal/sqlite"
	"strings"
)

//Получение всех организаций в здании
func (s *HTTPServer) house(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.URL)

	geo := strings.Split(url.Values.Get(r.Form, "geo"), ",")

	sqlite.GetHouse(w, url.Values.Get(r.Form, "addr"), geo[0], geo[1])
}

package server

import (
	"log"
	"net/http"
	"net/url"
	"simple2gis/intenal/sqlite"
	"strings"
)

//Получает организации по категории
func (s *HTTPServer) category(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.URL)
	sqlite.GetCategory(w, strings.Split(url.Values.Get(r.Form, "category"), ","))
}

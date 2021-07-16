package server

import (
	"log"
	"net/http"
	"net/url"
	"simple2gis/intenal/sqlite"
	"strings"
)

// Добавляет организацию в здание
func (s *HTTPServer) newOrganization(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.URL)
	//Считывание параметров в URL
	houseAddr := url.Values.Get(r.Form, "houseAddr")
	houseGeo := strings.Split(url.Values.Get(r.Form, "houseGeo"), ",")
	numbers := strings.Split(url.Values.Get(r.Form, "numbers"), ",")
	category := strings.Split(url.Values.Get(r.Form, "category"), ",")
	name := url.Values.Get(r.Form, "name")

	//Валидация телефонных номеров
	for _, v := range numbers {
		_, err := checkTelnum(v)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	sqlite.SetHouse(w, name, houseAddr, houseGeo[0], houseGeo[1], numbers, category)
}

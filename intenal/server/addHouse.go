package server

import (
	"fmt"
	"net/http"
	"simple2gis/intenal/sqlite"
	"strings"
)

// Добавляет организацию в здание
func (s *HTTPServer) addHouse(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.URL)
	httpArgs := s.agregator(r.Form)
	//Считывание параметров в URL
	houseAddr := httpArgs["houseAddr"]
	houseGeo := strings.Split(httpArgs["houseGeo"], "/")
	numbers := strings.Split(httpArgs["numbers"], "/")
	category := strings.Split(httpArgs["category"], "/")
	name := httpArgs["name"]

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

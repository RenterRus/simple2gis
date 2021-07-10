package server

import (
	"fmt"
	"net/http"
	"simple2gis/intenal/sqlite"
	"strings"
)

//Получение всех организаций в здании
func (s *HTTPServer) house(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.URL)
	httpArgs := s.agregator(r.Form)

	geo := strings.Split(httpArgs["geo"], "/")

	sqlite.GetHouse(w, httpArgs["addr"], geo[0], geo[1])
}

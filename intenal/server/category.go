package server

import (
	"fmt"
	"net/http"
	"simple2gis/intenal/sqlite"
	"strings"
)

//Получает организации по категории
func (s *HTTPServer) category(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.URL)
	httpArgs := s.agregator(r.Form)

	sqlite.GetCategory(w, strings.Split(httpArgs["category"], "/"))
}

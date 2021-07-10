package server

import (
	"fmt"
	"net/http"
	"simple2gis/intenal/sqlite"
)

//Получение организации по ID
func (s *HTTPServer) organization(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.URL)
	httpArgs := s.agregator(r.Form)

	id, err := CheckID(httpArgs["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sqlite.GetOrganization(w, id)
}

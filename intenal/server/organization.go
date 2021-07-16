package server

import (
	"log"
	"net/http"
	"simple2gis/intenal/sqlite"
)

//Получение организации по ID
func (s *HTTPServer) organization(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.URL)

	id, err := CheckID(r.Form.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sqlite.GetOrganization(w, id)
}

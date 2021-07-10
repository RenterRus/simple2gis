package server

import (
	"context"
	"net/http"
	"time"
)

//HTTP Сервер

type HTTPServer struct {
	s http.Server
}

func NewServer(addr string) *HTTPServer {
	s := new(HTTPServer)
	s.s = http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}

func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/house", s.house)
	mux.HandleFunc("/category", s.category)
	mux.HandleFunc("/organization", s.organization)
	mux.HandleFunc("/addHouse", s.addHouse)
	mux.HandleFunc("/genFakeBase", s.genFakeBase)
	s.s.Handler = mux
	return s.s.ListenAndServe()
}

//Красивое завершение
func (s *HTTPServer) GraceShutdown() {
	s.s.Shutdown(context.Background())
}

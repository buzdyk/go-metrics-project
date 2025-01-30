package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	config Config
}

func (s *Server) Run() {
	router := chi.NewRouter()

	router.Handle("GET /", http.HandlerFunc(GetIndex))
	router.Handle("POST /update/{type}/{metric}/{value}", http.HandlerFunc(StoreMetric))
	router.Handle("GET /value/{type}/{metric}", http.HandlerFunc(GetMetric))

	fmt.Println("Started metrics server on", s.config.Address)

	if err := http.ListenAndServe(s.config.Address, router); err != nil {
		panic(err)
	}
}

func NewServer(config Config) Server {
	return Server{
		config,
	}
}

package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct{}

func (s *Server) Run() {
	router := chi.NewRouter()
	router.Handle("POST /update/{type}/{metric}/{value}", http.HandlerFunc(StoreMetric))
	router.Handle("GET /value/{type}/{metric}", http.HandlerFunc(GetMetric))
	fmt.Println("Started metrics server on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

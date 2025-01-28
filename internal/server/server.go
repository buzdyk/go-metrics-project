package server

import (
	"fmt"
	"net/http"
)

type Server struct{}

func (s *Server) Run() {
	router := http.NewServeMux()
	router.Handle("POST /update/{type}/{metric}/{value}", metricExists(http.HandlerFunc(StoreMetric)))

	fmt.Println("Started metrics server on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

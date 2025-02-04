package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	config Config
}

func (s *Server) Run(ctx context.Context) {
	router := chi.NewRouter()

	router.Handle("GET /", http.HandlerFunc(GetIndex))
	router.Handle("POST /update/{type}/{metric}/{value}", http.HandlerFunc(StoreMetric))
	router.Handle("GET /value/{type}/{metric}", http.HandlerFunc(GetMetric))

	server := &http.Server{
		Addr:    s.config.Address,
		Handler: router,
	}

	go func() {
		fmt.Println("Started metrics server on", s.config.Address)
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	<-ctx.Done()

	fmt.Println("context is Done()")
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("shutdown server")
	}
}

func NewServer(config Config) Server {
	return Server{
		config,
	}
}

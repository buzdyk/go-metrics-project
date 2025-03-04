package server

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	config Config
}

func (s *Server) Run(ctx context.Context) {
	handler := NewMetricHandler(storage.NewCounterMemStorage(), storage.NewGaugeMemStorage())

	router := chi.NewRouter()
	router.Handle("GET /", http.HandlerFunc(handler.GetIndex))
	router.Handle("POST /update/{type}/{metric}/{value}", http.HandlerFunc(handler.StoreMetric))
	router.Handle("GET /value/{type}/{metric}", http.HandlerFunc(handler.GetMetric))

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

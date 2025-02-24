package server

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/server/handlers"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	config Config
}

func (s *Server) Run(ctx context.Context) {
	logger, _ := zap.NewProduction()
	handler := handlers.NewMetricHandler(storage.NewCounterMemStorage(), storage.NewGaugeMemStorage())

	router := chi.NewRouter()
	router.Handle("GET /", withMiddleware(logger, handler.GetIndex))
	router.Handle("POST /update/", withMiddleware(logger, handler.StoreMetricJSON))
	router.Handle("POST /value/", withMiddleware(logger, handler.GetMetricJSON))

	router.Handle("POST /update/{type}/{metric}/{value}", withMiddleware(logger, handler.StoreMetric))
	router.Handle("GET /value/{type}/{metric}", withMiddleware(logger, handler.GetMetric))

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

func withMiddleware(logger *zap.Logger, handler func(rw http.ResponseWriter, r *http.Request)) http.Handler {
	h := handlers.DecompressRequestMiddleware()(http.HandlerFunc(handler))
	h = handlers.CompressResponseMiddleware()(h)
	h = handlers.LoggingMiddleware(logger)(h)

	return h
}

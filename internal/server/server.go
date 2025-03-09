package server

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/buzdyk/go-metrics-project/internal/server/config"
	"github.com/buzdyk/go-metrics-project/internal/server/handlers"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type Server struct{}

func (s *Server) Run(ctx context.Context) {
	logger, _ := zap.NewProduction()
	cfg := config.GetConfig()

	var handler *handlers.MetricHandler

	if cfg.FileStoragePath != "" {
		cs := storage.NewFileStorage[metrics.Counter](cfg.FileStoragePath)
		gs := storage.NewFileStorage[metrics.Gauge](cfg.FileStoragePath)
		handler = handlers.NewMetricHandler(cs, gs)
	} else {
		cs := storage.NewCounterMemStorage()
		gs := storage.NewGaugeMemStorage()
		handler = handlers.NewMetricHandler(cs, gs)
	}

	router := chi.NewRouter()
	router.Handle("GET /", withMiddleware(logger, handler.GetIndex))
	router.Handle("GET /ping", withMiddleware(logger, handler.Ping))
	router.Handle("POST /update/", withMiddleware(logger, handler.StoreMetricJSON))
	router.Handle("POST /value/", withMiddleware(logger, handler.GetMetricJSON))

	router.Handle("POST /update/{type}/{metric}/{value}", withMiddleware(logger, handler.StoreMetric))
	router.Handle("GET /value/{type}/{metric}", withMiddleware(logger, handler.GetMetric))

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		fmt.Println("Started metrics server on", cfg.Address)
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	<-ctx.Done()
	fmt.Println("context is Done()")
}

func NewServer() Server {
	return Server{}
}

func withMiddleware(logger *zap.Logger, handler func(rw http.ResponseWriter, r *http.Request)) http.Handler {
	h := handlers.DecompressRequestMiddleware()(http.HandlerFunc(handler))
	h = handlers.CompressResponseMiddleware()(h)
	h = handlers.LoggingMiddleware(logger)(h)

	return h
}

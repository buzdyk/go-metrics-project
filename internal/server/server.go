package server

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/database"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/buzdyk/go-metrics-project/internal/server/config"
	"github.com/buzdyk/go-metrics-project/internal/server/handlers"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Server struct{}

func (s *Server) Run(ctx context.Context) {
	cfg := config.GetConfig()

	if cfg.PgDsn != "" {
		if err := database.GetClient().RunMigrations(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(err)
		}
	}

	mux := setupMux(cfg)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
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

func setupMux(cfg *config.Config) *chi.Mux {
	logger, _ := zap.NewProduction()

	cs, gs := getStorage(cfg)

	mux := chi.NewRouter()
	handler := handlers.NewMetricHandler(cs, gs)

	mux.Handle("GET /", withMiddleware(logger, handler.GetIndex))
	mux.Handle("GET /ping", withMiddleware(logger, handler.Ping))
	mux.Handle("POST /update/", withMiddleware(logger, handler.StoreMetricJSON))
	mux.Handle("POST /value/", withMiddleware(logger, handler.GetMetricJSON))

	mux.Handle("POST /update/{type}/{metric}/{value}", withMiddleware(logger, handler.StoreMetric))
	mux.Handle("GET /value/{type}/{metric}", withMiddleware(logger, handler.GetMetric))

	return mux
}

func getStorage(cfg *config.Config) (storage.Storage[metrics.Counter], storage.Storage[metrics.Gauge]) {
	if cfg.FileStoragePath != "" {
		cs := storage.NewFileStorage[metrics.Counter](cfg.FileStoragePath)
		gs := storage.NewFileStorage[metrics.Gauge](cfg.FileStoragePath)
		return cs, gs
	} else {
		cs := storage.NewCounterMemStorage()
		gs := storage.NewGaugeMemStorage()
		return cs, gs
	}
}

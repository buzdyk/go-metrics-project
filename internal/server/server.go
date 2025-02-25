package server

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/server/handlers"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	config Config
}

func (s *Server) Run(ctx context.Context) {
	logger, _ := zap.NewProduction()

	counterStore := storage.NewCounterMemStorage()
	gaugeStore := storage.NewGaugeMemStorage()

	b := Backup{
		s.config.FileStoragePath,
		gaugeStore,
		counterStore,
	}

	if s.config.Restore {
		if err := b.Restore(); err != nil {
			fmt.Println("backup restore error: ", err)
		}
	}

	handler := handlers.NewMetricHandler(counterStore, gaugeStore)

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

	backupTicker := time.NewTicker(time.Duration(s.config.StoreInterval) * time.Second)
	defer backupTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("context is Done()")

			if err := b.Backup(); err != nil {
				fmt.Println("server backup error: ", err)
			} else {
				fmt.Println("Backed up before shutdown")
			}

			if err := server.Shutdown(ctx); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("shutdown server")
			}

			return
		case <-backupTicker.C:
			go func() {
				if err := b.Backup(); err != nil {
					fmt.Print("server backup error: ", err)
				}
			}()
		}
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

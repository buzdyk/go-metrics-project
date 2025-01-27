package main

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/handlers"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"net/http"
)

var store *storage.MemStorage

func init() {
	store = storage.NewMemStorage()
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /update/{type}/{metric}/{value}", handlers.StoreMetric)

	fmt.Println("Started metrics server on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

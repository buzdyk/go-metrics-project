package main

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/server"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /update/{type}/{metric}/{value}", server.StoreMetric)

	fmt.Println("Started metrics server on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

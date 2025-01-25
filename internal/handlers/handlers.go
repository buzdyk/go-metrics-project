package handlers

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"log"
	"net/http"
)

var store = storage.NewMemStorage()

var StoreMetric = func(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	//metricName := r.PathValue("metric")
	//metricValue, _ := strconv.Atoi(r.PathValue("value"))

	g := metrics.Gauge{}

	switch metricType {
	case "gauge":
		store.StoreGauge(&g)
	case "counter":
		//store.StoreCounter()
	}

	log.Default().Println("type:", metricType, "metric", metricName)

	rw.WriteHeader(200)
	rw.Write([]byte("ok"))

	fmt.Println(storage)
}

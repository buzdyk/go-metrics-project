package server

import (
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"log"
	"net/http"
	"strconv"
)

var store = storage.NewMemStorage()

var StoreMetric = func(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")
	metricValue := r.PathValue("value")

	// todo move to middleware
	if metrics.Exists(r.PathValue("metric")) == false {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	v, err := strconv.ParseFloat(metricValue, 64)

	if err != nil {
		http.Error(rw, "metric value is not convertible to Float64", http.StatusBadRequest)
	}

	store.StoreGauge(metricName, metrics.Gauge(v))

	log.Default().Println("type:", metricType, "metric", metricName, metricValue)

	rw.WriteHeader(200)
	rw.Write([]byte("ok"))
}

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

	//if metrics.Exists(r.PathValue("metric")) == false {
	//	http.Error(rw, "metric does not exist", http.StatusBadRequest)
	//	return
	//}

	switch metricType {
	case "gauge":
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "metric value is not convertible to float64", http.StatusBadRequest)
		}
		store.StoreGauge(metricName, metrics.Gauge(v))
	case "counter":
		v, err := strconv.Atoi(metricValue)
		if err != nil {
			http.Error(rw, "metric value is not convertible to int", http.StatusBadRequest)
		}
		store.StoreCounter(metricName, metrics.Counter(v))
	}

	log.Default().Println("type:", metricType, "metric", metricName, metricValue)

	rw.WriteHeader(200)
	rw.Write([]byte("ok"))
}

var GetMetric = func(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")

	//if metrics.Exists(r.PathValue("metric")) == false {
	//	http.Error(rw, "metric does not exist", http.StatusBadRequest)
	//	return
	//}

	switch metricType {
	case "gauge":
		if v, err := store.GetGauge(metricName); err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		}
	case "counter":
		v, err := store.GetCounter(metricName)
		if err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.Itoa(int(v))))
		}
	}
}

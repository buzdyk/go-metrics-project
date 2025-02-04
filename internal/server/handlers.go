package server

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
)

var store = storage.NewMemStorage()

var StoreMetric = func(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")
	metricValue := r.PathValue("value")

	// todo move this to middleware
	if metrics.IsValidType(metricType) == false {
		rw.WriteHeader(400)
		return
	}

	if metrics.Exists(r.PathValue("metric")) == false {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch metricType {
	case metrics.GaugeName:
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "metric value is not convertible to float64", http.StatusBadRequest)
		}
		store.StoreGauge(metricName, metrics.Gauge(v))
	case metrics.CounterName:
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

	if metrics.IsValidType(metricType) == false {
		rw.WriteHeader(400)
		return
	}

	if metrics.Exists(metricName) == false {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch metricType {
	case "gauge":
		if v, err := store.Gauge(metricName); err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		}
	case metrics.CounterName:
		v, err := store.Counter(metricName)
		if err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.Itoa(int(v))))
		}
	}
}

var GetIndex = func(rw http.ResponseWriter, r *http.Request) {
	data := struct {
		Gauges   map[string]metrics.Gauge
		Counters map[string]metrics.Counter
	}{
		store.Gauges(),
		store.Counters(),
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	tmpl, err := template.ParseFiles(dir + "/templates/index.html")
	if err != nil {
		http.Error(rw, "Failed to parse HTML template", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(rw, data); err != nil {
		fmt.Println(err)
		http.Error(rw, "Failed to render metrics page", http.StatusInternalServerError)
	}

}

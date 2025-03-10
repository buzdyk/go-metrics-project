package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"net/http"
	"strconv"
)

func (mh *MetricHandler) StoreMetric(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")
	metricValue := r.PathValue("value")

	if !metrics.IsValidType(metricType) {
		http.Error(rw, "metric type is invalid", http.StatusBadRequest)
		return
	}

	if !metrics.Exists(metricName) {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch metricType {
	case metrics.GaugeName:
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "metric value is not convertible to int", http.StatusBadRequest)
		}
		err = mh.gaugeStore.Store(metricName, metrics.Gauge(v))
		fmt.Println("error storing gauge metric: ", err)
	case metrics.CounterName:
		v, err := strconv.Atoi(metricValue)
		if err != nil {
			http.Error(rw, "metric value is not convertible to int", http.StatusBadRequest)
		}
		currentValue, _ := mh.counterStore.Value(metricName)
		err = mh.counterStore.Store(metricName, metrics.Counter(v)+currentValue)
		fmt.Println("error storing counter metric: ", err)
	}

	rw.WriteHeader(200)
	rw.Write([]byte("ok"))
}

func (mh *MetricHandler) StoreMetricJSON(rw http.ResponseWriter, r *http.Request) {
	var m metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if !metrics.IsValidType(m.MType) {
		http.Error(rw, "metric type is invalid", http.StatusBadRequest)
		return
	}

	if !metrics.Exists(m.ID) {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch m.MType {
	case metrics.GaugeName:
		err := mh.gaugeStore.Store(m.ID, *m.Value)
		fmt.Println("error storing gauge metric: ", err)
	case metrics.CounterName:
		currentValue, _ := mh.counterStore.Value(m.ID)
		newValue := *m.Delta + currentValue
		err := mh.counterStore.Store(m.ID, newValue)
		m.Delta = &newValue
		fmt.Println("error storing gauge metric: ", err)
	}

	resp, _ := json.Marshal(m)

	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(resp)
}

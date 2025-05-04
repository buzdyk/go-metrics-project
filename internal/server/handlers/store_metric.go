package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/collector"
	"github.com/buzdyk/go-metrics-project/internal/models"
	"net/http"
	"strconv"
)

func (mh *MetricHandler) StoreMetric(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")
	metricValue := r.PathValue("value")

	if !collector.IsValidType(metricType) {
		http.Error(rw, "metric type is invalid", http.StatusBadRequest)
		return
	}

	if !collector.Exists(metricName) {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch metricType {
	case models.GaugeName:
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "metric value is not convertible to int", http.StatusBadRequest)
		}
		err = mh.gaugeStore.Store(r.Context(), metricName, models.Gauge(v))
		fmt.Println("error storing gauge metric: ", err)
	case models.CounterName:
		v, err := strconv.Atoi(metricValue)
		if err != nil {
			http.Error(rw, "metric value is not convertible to int", http.StatusBadRequest)
		}
		currentValue, _ := mh.counterStore.Value(r.Context(), metricName)
		err = mh.counterStore.Store(r.Context(), metricName, models.Counter(v)+currentValue)
		fmt.Println("error storing counter metric: ", err)
	}

	rw.WriteHeader(200)
	rw.Write([]byte("ok"))
}

func (mh *MetricHandler) StoreMetricJSON(rw http.ResponseWriter, r *http.Request) {
	var m models.Metric

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if !collector.IsValidType(m.MType) {
		http.Error(rw, "metric type is invalid", http.StatusBadRequest)
		return
	}

	if !collector.Exists(m.ID) {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch m.MType {
	case models.GaugeName:
		err := mh.gaugeStore.Store(r.Context(), m.ID, *m.Value)
		fmt.Println("error storing gauge metric: ", err)
	case models.CounterName:
		currentValue, _ := mh.counterStore.Value(r.Context(), m.ID)
		newValue := *m.Delta + currentValue
		err := mh.counterStore.Store(r.Context(), m.ID, newValue)
		m.Delta = &newValue
		fmt.Println("error storing gauge metric: ", err)
	}

	resp, _ := json.Marshal(m)

	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(resp)
}

func (mh *MetricHandler) Updates(rw http.ResponseWriter, r *http.Request) {
	var data []models.Metric

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	for _, metric := range data {
		if !collector.IsValidType(metric.MType) {
			http.Error(rw, "metric type is invalid", http.StatusBadRequest)
			return
		}

		if !collector.Exists(metric.ID) {
			http.Error(rw, "metric does not exist", http.StatusBadRequest)
			return
		}
	}

	gauges := make(map[string]models.Gauge)

	for _, metric := range data {
		switch metric.MType {
		case models.GaugeName:
			gauges[metric.ID] = *metric.Value
		case models.CounterName:
			currentValue, _ := mh.counterStore.Value(r.Context(), metric.ID)
			newValue := *metric.Delta + currentValue
			mh.counterStore.Store(r.Context(), metric.ID, newValue)
		}
	}

	if err := mh.gaugeStore.StoreMany(r.Context(), gauges); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
}

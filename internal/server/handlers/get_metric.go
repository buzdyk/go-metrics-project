package handlers

import (
	"encoding/json"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"net/http"
	"strconv"
)

func (mh *MetricHandler) GetMetric(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")

	if !metrics.IsValidType(metricType) {
		rw.WriteHeader(400)
		return
	}

	if !metrics.Exists(metricName) {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch metricType {
	case metrics.GaugeName:
		if v, err := mh.gaugeStore.Value(metricName); err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		}
	case metrics.CounterName:
		v, err := mh.counterStore.Value(metricName)
		if err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.Itoa(int(v))))
		}
	}
}

func (mh *MetricHandler) GetMetricJSON(rw http.ResponseWriter, r *http.Request) {
	var m metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if !metrics.IsValidType(m.MType) {
		rw.WriteHeader(400)
		return
	}

	if !metrics.Exists(m.ID) {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch m.MType {
	case metrics.GaugeName:
		if v, err := mh.gaugeStore.Value(m.ID); err != nil {
			rw.WriteHeader(404)
		} else {
			m.Value = &v
			resp, _ := json.Marshal(m)
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(resp)
		}
	case metrics.CounterName:
		v, err := mh.counterStore.Value(m.ID)
		if err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Header().Set("Content-Type", "application/json")
			m.Delta = &v
			resp, _ := json.Marshal(m)
			rw.Write(resp)
		}
	}
}

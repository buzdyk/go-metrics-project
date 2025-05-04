package handlers

import (
	"github.com/buzdyk/go-metrics-project/internal/collector"
	"github.com/buzdyk/go-metrics-project/internal/models"
	"net/http"
	"strconv"
)

func (mh *MetricHandler) GetMetric(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")

	if !collector.IsValidType(metricType) {
		rw.WriteHeader(400)
		return
	}

	if !collector.Exists(metricName) {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch metricType {
	case models.GaugeName:
		if v, err := mh.gaugeStore.Value(r.Context(), metricName); err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		}
	case models.CounterName:
		v, err := mh.counterStore.Value(r.Context(), metricName)
		if err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.Itoa(int(v))))
		}
	}
}

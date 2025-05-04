package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/collector"
	"github.com/buzdyk/go-metrics-project/internal/models"
	"net/http"
)

func (mh *MetricHandler) GetMetricJSON(rw http.ResponseWriter, r *http.Request) {
	var m models.Metric

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if !collector.IsValidType(m.MType) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !collector.Exists(m.ID) {
		http.Error(rw, "metric does not exist", http.StatusBadRequest)
		return
	}

	switch m.MType {
	case models.GaugeName:
		if v, err := mh.gaugeStore.Value(r.Context(), m.ID); err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusNotFound)
		} else {
			m.Value = &v
			resp, _ := json.Marshal(m)
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(resp)
		}
	case models.CounterName:
		v, err := mh.counterStore.Value(r.Context(), m.ID)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.Header().Set("Content-Type", "application/json")
			m.Delta = &v
			resp, _ := json.Marshal(m)
			rw.Write(resp)
		}
	}
}

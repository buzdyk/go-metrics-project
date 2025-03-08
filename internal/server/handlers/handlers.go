package handlers

import (
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/buzdyk/go-metrics-project/internal/storage"
)

type MetricHandler struct {
	counterStore storage.Storage[metrics.Counter]
	gaugeStore   storage.Storage[metrics.Gauge]
}

func NewMetricHandler(
	counterStore storage.Storage[metrics.Counter],
	gaugeStore storage.Storage[metrics.Gauge],
) *MetricHandler {
	return &MetricHandler{
		counterStore: counterStore,
		gaugeStore:   gaugeStore,
	}
}

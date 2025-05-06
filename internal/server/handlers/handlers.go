package handlers

import (
	"github.com/buzdyk/go-metrics-project/internal/models"
	"github.com/buzdyk/go-metrics-project/internal/storage"
)

type MetricHandler struct {
	counterStore storage.Storage[models.Counter]
	gaugeStore   storage.Storage[models.Gauge]
}

func NewMetricHandler(
	counterStore storage.Storage[models.Counter],
	gaugeStore storage.Storage[models.Gauge],
) *MetricHandler {
	return &MetricHandler{
		counterStore: counterStore,
		gaugeStore:   gaugeStore,
	}
}

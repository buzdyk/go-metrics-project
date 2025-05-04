package collector

import (
	"github.com/buzdyk/go-metrics-project/internal/models"
)

type Gauge = models.Gauge
type Counter = models.Counter
type Metric = models.Metric

const (
	GaugeName   = models.GaugeName
	CounterName = models.CounterName
)

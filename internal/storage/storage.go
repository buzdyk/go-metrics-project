package storage

import (
	"github.com/buzdyk/go-metrics-project/internal/metrics"
)

type AllowedTypes interface {
	metrics.Gauge | metrics.Counter
}

type Storage[T AllowedTypes] interface {
	Store(name string, value T) error
	StoreMany(map[string]T) error
	Value(name string) (T, error)
	Values() (map[string]T, error)
}

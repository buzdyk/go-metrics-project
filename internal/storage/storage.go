package storage

import (
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"sync"
)

type Storage interface {
	StoreGauge(m *metrics.Gauge) (bool, error)
}

type MemStorage struct {
	g map[string][]metrics.Gauge
	c map[string]metrics.Counter
	m sync.Mutex
}

func (s *MemStorage) StoreGauge(name string, v metrics.Gauge) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	s.g[name] = append(s.g[name], v)

	return true, nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		g: make(map[string][]metrics.Gauge),
	}
}

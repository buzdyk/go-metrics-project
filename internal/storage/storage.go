package storage

import (
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"sync"
)

type Storage interface {
	StoreGauge(m *metrics.Gauge) (bool, error)
}

type MemStorage struct {
	g  map[string][]metrics.Gauge
	c  map[string]metrics.Counter
	gm sync.Mutex
	cm sync.Mutex
}

func (s *MemStorage) StoreGauge(name string, v metrics.Gauge) (bool, error) {
	s.gm.Lock()
	defer s.gm.Unlock()

	s.g[name] = append(s.g[name], v)

	return true, nil
}

func (s *MemStorage) StoreCounter(name string, v metrics.Counter) (bool, error) {
	s.cm.Lock()
	defer s.cm.Unlock()

	curr, ok := s.c[name]
	if !ok {
		curr = 0
	}

	s.c[name] = v + curr

	return true, nil
}

func (s *MemStorage) GetGauge(name string) (metrics.Gauge, error) {
	if v, ok := s.g[name]; !ok {
		return 0, errors.New("unknown gauge metric:" + name)
	} else {
		return v[len(v)-1], nil
	}
}

func (s *MemStorage) GetCounter(name string) (metrics.Counter, error) {
	if v, ok := s.c[name]; !ok {
		return 0, errors.New("unknown counter metric:" + name)
	} else {
		return v, nil
	}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		g: make(map[string][]metrics.Gauge),
		c: make(map[string]metrics.Counter),
	}
}

package storage

import (
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"sync"
)

type Storage interface {
	StoreGauge(m *metrics.Gauge) (bool, error)
	StoreCounter(m *metrics.Counter) (bool, error)
	Gauges() map[string]metrics.Gauge
	Counters() map[string]metrics.Counter
	Gauge(name string) (metrics.Gauge, error)
	Counter(name string) (metrics.Counter, error)
}

type MemStorage struct {
	g  map[string]metrics.Gauge
	c  map[string]metrics.Counter
	mu sync.RWMutex
}

func (s *MemStorage) Gauges() map[string]metrics.Gauge {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m := make(map[string]metrics.Gauge, len(s.g))
	for k, v := range s.g {
		m[k] = v
	}
	return m
}

func (s *MemStorage) Counters() map[string]metrics.Counter {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m := make(map[string]metrics.Counter, len(s.c))
	for k, v := range s.c {
		m[k] = v
	}
	return m
}

func (s *MemStorage) StoreGauge(name string, v metrics.Gauge) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.g[name] = v

	return true, nil
}

func (s *MemStorage) StoreCounter(name string, v metrics.Counter) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	curr, ok := s.c[name]
	if !ok {
		curr = 0
	}

	s.c[name] = v + curr

	return true, nil
}

func (s *MemStorage) Gauge(name string) (metrics.Gauge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.g[name]; !ok {
		return 0, errors.New("unknown gauge metric:" + name)
	} else {
		return v, nil
	}
}

func (s *MemStorage) Counter(name string) (metrics.Counter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.c[name]; !ok {
		return 0, errors.New("unknown counter metric:" + name)
	} else {
		return v, nil
	}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		g: make(map[string]metrics.Gauge),
		c: make(map[string]metrics.Counter),
	}
}

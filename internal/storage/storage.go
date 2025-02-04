package storage

import (
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"sync"
)

type AllowedTypes interface {
	metrics.Gauge | metrics.Counter
}

type Storage[T AllowedTypes] interface {
	Store(m *T) (bool, error)
	Values() map[string]T
	Value(name string) (T, error)
}

type MemStorage[T AllowedTypes] struct {
	c  map[string]T
	mu sync.RWMutex
}

func (s *MemStorage[T]) Store(name string, v T) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.c[name] = v

	return true, nil
}

func (s *MemStorage[T]) Values() map[string]T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m := make(map[string]T, len(s.c))
	for k, v := range s.c {
		m[k] = v
	}

	return m
}

func (s *MemStorage[T]) Value(name string) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.c[name]; !ok {
		return 0, errors.New("unknown gauge metric:" + name)
	} else {
		return v, nil
	}
}

func NewMemStorage[T AllowedTypes]() *MemStorage[T] {
	return &MemStorage[T]{
		c: make(map[string]T),
	}
}

func NewGaugeMemStorage() *MemStorage[metrics.Gauge] {
	return NewMemStorage[metrics.Gauge]()
}

func NewCounterMemStorage() *MemStorage[metrics.Counter] {
	return NewMemStorage[metrics.Counter]()
}

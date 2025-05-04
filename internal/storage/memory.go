package storage

import (
	"context"
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/models"
	"sync"
)

type MemStorage[T AllowedTypes] struct {
	c  map[string]T
	mu sync.RWMutex
}

func (s *MemStorage[T]) Store(ctx context.Context, name string, v T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.c[name] = v

	return nil
}

func (s *MemStorage[T]) StoreMany(ctx context.Context, m map[string]T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range m {
		s.c[i] = v
	}

	return nil
}

func (s *MemStorage[T]) Values(ctx context.Context) (map[string]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m := make(map[string]T, len(s.c))
	for k, v := range s.c {
		m[k] = v
	}

	return m, nil
}

func (s *MemStorage[T]) Value(ctx context.Context, name string) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.c[name]; !ok {
		return 0, errors.New("unknown metric:" + name)
	} else {
		return v, nil
	}
}

func NewMemStorage[T AllowedTypes]() *MemStorage[T] {
	return &MemStorage[T]{
		c: make(map[string]T),
	}
}

func NewGaugeMemStorage() *MemStorage[models.Gauge] {
	return NewMemStorage[models.Gauge]()
}

func NewCounterMemStorage() *MemStorage[models.Counter] {
	return NewMemStorage[models.Counter]()
}

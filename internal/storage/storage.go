package storage

import "github.com/buzdyk/go-metrics-project/internal/metrics"

type Storage interface {
	StoreCounter(m *metrics.Counter) (bool, error)
	StoreGauge(m *metrics.Gauge) (bool, error)
}

type MemStorage struct {
	g map[string][]metrics.Gauge
	c map[string]metrics.Counter
}

func (s *MemStorage) StoreGauge(g *metrics.Gauge) (bool, error) {
	v, _ := g.Value()
	s.g[g.ID()] = append(s.g[g.ID()], v)

	return true, nil
}

func (s *MemStorage) StoreCounter(_ *metrics.Counter) (bool, error) {
	//if _, ok := s.c[id]; !ok {
	//	s.c[id] = 0
	//}
	//s.c[id] += value

	return true, nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		g: make(map[string][]metrics.Gauge),
		c: make(map[string]metrics.Counter),
	}
}

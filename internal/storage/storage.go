package storage

import "github.com/buzdyk/go-metrics-project/internal/metrics"

type Storage interface {
	StoreGauge(m *metrics.Gauge) (bool, error)
}

type MemStorage struct {
	g map[string][]metrics.Gauge
}

func (s *MemStorage) StoreGauge(g *metrics.Gauge) (bool, error) {
	//v, _ := g.Value()
	//s.g[g.ID()] = append(s.g[g.ID()], v)

	return true, nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		g: make(map[string][]metrics.Gauge),
	}
}

package agent

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"net/http"
	"sync"
	"time"
)

type Syncer interface {
	SyncMetric(name string, value any) (*http.Response, error)
	SyncMetrics([]metrics.Metric) (*http.Response, error)
}

type MetricsCollector interface {
	Collect(out map[string]any)
}

type Agent struct {
	config    Config
	collector MetricsCollector
	syncer    Syncer
	data      map[string]interface{}
	mu        sync.RWMutex
}

func (a *Agent) collect() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.collector.Collect(a.data)
}

func (a *Agent) sync() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var wg sync.WaitGroup
	wg.Add(len(a.data))

	var data []metrics.Metric

	for id, value := range a.data {
		switch v := value.(type) {
		case metrics.Gauge:
			data = append(data, metrics.Metric{
				ID:    id,
				MType: metrics.GaugeName,
				Value: v,
			})
		case metrics.Counter:
			data = append(data, metrics.Metric{
				ID:    id,
				MType: metrics.CounterName,
				Delta: v,
			})
		}
	}

	if len(data) > 0 {
		a.syncer.SyncMetrics(data)
	}
}

func (a *Agent) Run(ctx context.Context) {
	pollTicker := time.NewTicker(time.Duration(a.config.Collect) * time.Second)
	syncTicker := time.NewTicker(time.Duration(a.config.Report) * time.Second)

	defer pollTicker.Stop()
	defer syncTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("context is Done()")
			return
		case <-pollTicker.C:
			go a.collect()
		case <-syncTicker.C:
			go a.sync()
		}
	}
}

func NewAgent(config *Config, collector MetricsCollector, client Syncer) *Agent {
	return &Agent{
		config:    *config,
		collector: collector,
		syncer:    client,
		data:      make(map[string]any),
	}
}

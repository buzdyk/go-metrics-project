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

	var data []metrics.Metric

	for id, value := range a.data {
		switch v := value.(type) {
		case metrics.Gauge:
			data = append(data, metrics.Metric{
				ID:    id,
				MType: metrics.GaugeName,
				Value: &v,
			})
		case metrics.Counter:
			data = append(data, metrics.Metric{
				ID:    id,
				MType: metrics.CounterName,
				Delta: &v,
			})
		}
	}

	if len(data) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	retryBackoffs := []time.Duration{1, 3, 5} // seconds
	
	for attempt := 0; attempt <= len(retryBackoffs); attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(retryBackoffs[attempt-1] * time.Second):
			}
		}
		
		res, err := a.syncer.SyncMetrics(data)
		if err == nil && res.StatusCode < 500 {
			defer res.Body.Close()
			return
		}
		
		if err != nil {
			fmt.Printf("Attempt %d failed: %v\n", attempt+1, err)
		} else {
			fmt.Printf("Attempt %d failed with status: %d\n", attempt+1, res.StatusCode)
			res.Body.Close()
		}
		
		if attempt == len(retryBackoffs) {
			fmt.Println("All retry attempts failed")
		}
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

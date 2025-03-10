package agent

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Syncer interface {
	SyncMetric(name string, value any) (*http.Response, error)
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

	for id, value := range a.data {
		go func(id string, value any) {
			defer wg.Done()

			resp, err := a.syncer.SyncMetric(id, value)
			if err == nil {
				defer resp.Body.Close()
			}

			if err != nil {
				fmt.Println(err)
				return
			}
		}(id, value)
	}

	wg.Wait()
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

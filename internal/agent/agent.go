package agent

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/agent/config"
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
	CollectSystem(out map[string]any)
}

type Agent struct {
	config    config.Config
	collector MetricsCollector
	syncer    Syncer
	data      map[string]interface{}
	mu        sync.RWMutex
	metricsCh chan []metrics.Metric
	workersWg sync.WaitGroup
}

func (a *Agent) collect() {
	a.collector.Collect(a.data)
}

func (a *Agent) collectSystem() {
	a.collector.CollectSystem(a.data)
}

func (a *Agent) prepareMetrics() {
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

	select {
	case a.metricsCh <- data:
	default:
		// канал заполнен, пропускаем эту партию метрик
		fmt.Println("Warning: metrics channel is full, skipping batch")
	}
}

func (a *Agent) syncWorker(ctx context.Context, id int) {
	defer a.workersWg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-a.metricsCh:
			if !ok {
				return
			}

			reqCtx, cancel := context.WithTimeout(ctx, time.Duration(a.config.RequestTimeout)*time.Second)
			retryBackoffs := a.config.RetryBackoffs

			for attempt := 0; attempt <= len(retryBackoffs); attempt++ {
				if attempt > 0 {
					select {
					case <-reqCtx.Done():
						cancel()
						return
					case <-time.After(retryBackoffs[attempt-1]):
					}
				}

				res, err := a.syncer.SyncMetrics(data)
				if err == nil && res.StatusCode < 500 {
					res.Body.Close()
					break
				}

				if err != nil {
					fmt.Printf("Worker %d: Attempt %d failed: %v\n", id, attempt+1, err)
				} else {
					fmt.Printf("Worker %d: Attempt %d failed with status: %d\n", id, attempt+1, res.StatusCode)
					res.Body.Close()
				}

				if attempt == len(retryBackoffs) {
					fmt.Printf("Worker %d: All retry attempts failed\n", id)
				}
			}

			cancel()
		}
	}
}

func (a *Agent) Run(ctx context.Context) {
	for i := 0; i < a.config.RateLimit; i++ {
		a.workersWg.Add(1)
		go a.syncWorker(ctx, i+1)
	}

	pollTicker := time.NewTicker(time.Duration(a.config.Collect) * time.Second)
	reportTicker := time.NewTicker(time.Duration(a.config.Report) * time.Second)

	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Agent context cancelled, shutting down...")
			close(a.metricsCh)
			a.workersWg.Wait()
			return
		case <-pollTicker.C:
			go a.collect()
			go a.collectSystem()
		case <-reportTicker.C:
			go a.prepareMetrics()
		}
	}
}

func NewAgent(config *config.Config, collector MetricsCollector, client Syncer) *Agent {
	return &Agent{
		config:    *config,
		collector: collector,
		syncer:    client,
		data:      make(map[string]any),
		metricsCh: make(chan []metrics.Metric, config.RateLimit),
	}
}

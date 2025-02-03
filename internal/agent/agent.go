package agent

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type HTTPClient interface {
	Post(name string, value interface{}) (*http.Response, error)
}

type MetricsCollector interface {
	Collect(out map[string]interface{})
}

type Agent struct {
	config    Config
	collector MetricsCollector
	client    HTTPClient
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

	for id, value := range a.data {
		wg.Add(1)
		go func(id string, value interface{}) {
			defer wg.Done()
			var err error
			var r *http.Response
			r, err = a.client.Post(id, value)

			if err != nil {
				fmt.Println(err)
				return
			}

			err = r.Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(id, value)
	}

	wg.Wait()
}

func (a *Agent) Run() {
	pollTicker := time.NewTicker(time.Duration(a.config.Collect) * time.Second)
	syncTicker := time.NewTicker(time.Duration(a.config.Report) * time.Second)

	defer pollTicker.Stop()
	defer syncTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			go a.collect()
		case <-syncTicker.C:
			go a.sync()
		}
	}
}

func NewAgent(config *Config, collector MetricsCollector, client HTTPClient) (*Agent, error) {
	return &Agent{
		*config,
		collector,
		client,
		make(map[string]interface{}),
		sync.RWMutex{},
	}, nil
}

package agent

import (
	"fmt"
	"net/http"
	"time"
)

type HttpClient interface {
	Post(name string, value interface{}) (*http.Response, error)
}

type MetricsCollector interface {
	Collect(map[string]interface{})
}

type Agent struct {
	data      map[string]interface{}
	collector MetricsCollector
	client    HttpClient
}

func (a *Agent) collect() {
	a.collector.Collect(a.data)
}

func (a *Agent) sync() {
	for id, value := range a.data {
		go func() {
			if _, err := a.client.Post(id, value); err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func (a *Agent) Run() {
	pollTicker := time.NewTicker(2 * time.Second)
	syncTicker := time.NewTicker(10 * time.Second)

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

func NewAgent(collector MetricsCollector, client HttpClient) (*Agent, error) {
	return &Agent{
		make(map[string]interface{}),
		collector,
		client,
	}, nil
}

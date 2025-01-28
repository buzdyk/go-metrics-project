package agent

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"io"
	"net/http"
	"time"
)

type HttpClient interface {
	Post(endpoint, contentType string, body io.Reader) (*http.Response, error)
}

type Agent struct {
	data       map[string]interface{}
	collectors map[string]func() metrics.Gauge
	client     HttpClient
}

func (a *Agent) collect() {
	for id, getter := range a.collectors {
		a.data[id] = getter()
	}
}

func (a *Agent) sync() {
	for id, value := range a.data {
		go func(id string, value interface{}) {
			if err := a.syncOne(id, value); err != nil {
				fmt.Println(err)
			}
		}(id, value)
	}
}

func (a *Agent) syncOne(id string, value interface{}) error {
	endpoint := fmt.Sprintf("http://127.0.0.1:8080/update/gauge/%v/%v", id, value)

	_, err := a.client.Post(endpoint, "text/plain", nil)

	if err != nil {
		return err
	}

	return nil
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

func NewAgent(collectors map[string]func() metrics.Gauge, client HttpClient) (*Agent, error) {
	return &Agent{
		make(map[string]interface{}),
		collectors,
		client,
	}, nil
}

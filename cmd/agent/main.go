package main

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"net/http"
	"time"
)

type agent struct {
	data       map[string]interface{}
	collectors map[string]func() metrics.Gauge
}

func (a *agent) collect() {
	for id, getter := range a.collectors {
		a.data[id] = getter()
	}
}

func (a *agent) sync() {
	for id, value := range a.data {
		go func(id string, value interface{}) {
			endpoint := fmt.Sprintf("http://127.0.0.1:8080/update/gauge/%v/%v", id, value)

			_, err := http.Post(endpoint, "text/plain", nil)

			if err != nil {
				fmt.Println(err)
				return
			}
		}(id, value)
	}
}

func (a *agent) run() {
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

func newAgent(collectors map[string]func() metrics.Gauge) (*agent, error) {
	return &agent{
		make(map[string]interface{}),
		collectors,
	}, nil
}

func main() {
	a, err := newAgent(metrics.Collectors)

	if err != nil {
		panic(err)
	}

	a.run()
}

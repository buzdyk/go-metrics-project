package main

import (
	"github.com/buzdyk/go-metrics-project/internal/agent"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
)

func main() {
	collector := &metrics.Collector{}
	client := &agent.RealHttpClient{
		Host: "http://127.0.0.1:8080",
	}

	a, err := agent.NewAgent(collector, client)

	if err != nil {
		panic(err)
	}

	a.Run()
}

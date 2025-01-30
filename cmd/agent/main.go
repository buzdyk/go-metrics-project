package main

import (
	"github.com/buzdyk/go-metrics-project/internal/agent"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
)

func main() {
	collector := &metrics.Collector{}
	client := &agent.RealHTTPClient{
		Host: config.Address,
	}

	agentConfig := agent.Config{
		Report:  config.Report,
		Collect: config.Collect,
	}

	a, err := agent.NewAgent(agentConfig, collector, client)

	if err != nil {
		panic(err)
	}

	a.Run()
}

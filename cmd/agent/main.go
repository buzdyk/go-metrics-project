package main

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/agent"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
)

func main() {
	config := agent.NewConfigFromCLI()

	collector := &metrics.Collector{}
	client := &agent.RealHTTPClient{
		Host: config.Address,
	}

	a := agent.NewAgent(config, collector, client)

	fmt.Println("started agent")
	fmt.Println("  with config: ", config)

	a.Run()
}

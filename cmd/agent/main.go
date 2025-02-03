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

	a, err := agent.NewAgent(config, collector, client)

	if err != nil {
		panic(err)
	}

	fmt.Println("started agent")
	fmt.Println("  with config: ", config)

	a.Run()
}

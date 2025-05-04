package main

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/agent"
	"github.com/buzdyk/go-metrics-project/internal/agent/config"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := config.NewConfigFromCLI()
	collector := metrics.NewCollector()
	client := agent.NewHTTPSyncer(config.Address, config.Key)

	a := agent.NewAgent(config, collector, client)

	ctx, cancel := context.WithCancel(context.Background())

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh
		fmt.Println("received termination signal")
		cancel()
	}()

	fmt.Println("starting agent")
	fmt.Println("  with config: ", config)

	a.Run(ctx)
}

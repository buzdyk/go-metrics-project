package main

import (
	"context"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := server.NewConfigFromCLI()
	s := server.NewServer(c)

	ctx, cancel := context.WithCancel(context.Background())

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh
		fmt.Println("received termination signal")
		cancel()
	}()

	s.Run(ctx)
}

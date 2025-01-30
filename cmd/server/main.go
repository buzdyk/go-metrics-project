package main

import "github.com/buzdyk/go-metrics-project/internal/server"

func main() {
	s := server.NewServer(server.Config{
		Address: config.Address,
	})

	s.Run()
}

package main

import "github.com/buzdyk/go-metrics-project/internal/server"

func main() {
	c := server.NewConfigFromCLI()
	s := server.NewServer(c)
	s.Run()
}

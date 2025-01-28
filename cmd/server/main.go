package main

import "github.com/buzdyk/go-metrics-project/internal/server"

func main() {
	s := server.Server{}
	s.Run()
}

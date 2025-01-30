package main

import (
	"flag"
	"os"
)

type Config struct {
	Address string
}

var config Config

func init() {
	addrFlag := flag.String("a", "", "Address to listen on")
	flag.Parse()

	envAddr := os.Getenv("ADDRESS")

	switch {
	case envAddr != "":
		config.Address = envAddr
	case *addrFlag != "":
		config.Address = *addrFlag
	default:
		config.Address = "0.0.0.0:8080"
	}
}

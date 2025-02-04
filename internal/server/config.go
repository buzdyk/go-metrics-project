package server

import (
	"flag"
	"os"
)

type Config struct {
	Address string
}

func NewConfig() Config {
	return Config{
		"0.0.0.0:8080",
	}
}

func NewConfigFromCLI() Config {
	config := NewConfig()

	address := flag.String("a", config.Address, "Address to listen on")

	flag.Parse()

	config.Address = *address

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		config.Address = envAddr
	}

	return config
}

package server

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
)

type Config struct {
	Address string `env:ADDRESS`
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

	if err := env.Parse(&config); err != nil {
		panic(fmt.Sprintf("error parsing env: %v", err))
	}

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		config.Address = envAddr
	}

	return config
}

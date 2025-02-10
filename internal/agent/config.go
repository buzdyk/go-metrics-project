package agent

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"strings"
)

type Config struct {
	Address string `env:"ADDRESS"`
	Report  int    `env:"REPORT"`
	Collect int    `env:"COLLECT"`
}

func NewConfig() Config {
	return Config{"0.0.0.0:8080", 10, 2}
}

func NewConfigFromCLI() *Config {
	config := NewConfig()

	address := flag.String("a", config.Address, "Address of the server that collects metrics")
	report := flag.Int("r", config.Report, "Report interval in seconds")
	collect := flag.Int("p", config.Collect, "Poll interval in seconds")

	flag.Parse()

	config.Address = *address
	config.Report = *report
	config.Collect = *collect

	if err := env.Parse(&config); err != nil {
		panic(fmt.Sprintf("error parsing env %v", err))
	}

	if !strings.HasPrefix(config.Address, "http://") {
		config.Address = "http://" + config.Address
	}

	return &config
}

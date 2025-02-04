package agent

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Address string
	Report  int
	Collect int
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

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		config.Address = envAddr
	}

	if !strings.HasPrefix(config.Address, "http://") {
		config.Address = "http://" + config.Address
	}

	if envReport := os.Getenv("REPORT"); envReport != "" {
		v, err := strconv.Atoi(envReport)
		if err != nil {
			panic(err)
		}
		config.Report = v
	}

	if envCollect := os.Getenv("COLLECT"); envCollect != "" {
		v, err := strconv.Atoi(envCollect)
		if err != nil {
			panic(err)
		}
		config.Collect = v
	}

	return &config
}

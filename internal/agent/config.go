package agent

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"strings"
	"time"
)

type Config struct {
	Address        string          `env:"ADDRESS"`
	Report         int             `env:"REPORT"`
	Collect        int             `env:"COLLECT"`
	RequestTimeout int             `env:"REQUEST_TIMEOUT"`
	RetryBackoffs  []time.Duration `env:"RETRY_BACKOFFS"`
	Key            string          `env:"KEY"`
}

func NewConfig() Config {
	return Config{
		Address:        "0.0.0.0:8080",
		Report:         10,
		Collect:        2,
		RequestTimeout: 15,
		RetryBackoffs:  []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
	}
}

func NewConfigFromCLI() *Config {
	config := NewConfig()

	address := flag.String("a", config.Address, "Address of the server that collects metrics")
	report := flag.Int("r", config.Report, "Report interval in seconds")
	collect := flag.Int("p", config.Collect, "Poll interval in seconds")
	timeout := flag.Int("t", config.RequestTimeout, "Request timeout in seconds")
	key := flag.String("k", "", "Key for SHA256 signature")

	defaultBackoffs := "1,3,5"
	backoffsStr := flag.String("b", defaultBackoffs, "Retry backoffs in seconds (comma-separated, e.g., '1,3,5')")

	flag.Parse()

	config.Address = *address
	config.Report = *report
	config.Collect = *collect
	config.RequestTimeout = *timeout
	config.Key = *key

	if *backoffsStr != "" {
		backoffValues := strings.Split(*backoffsStr, ",")
		config.RetryBackoffs = make([]time.Duration, 0, len(backoffValues))

		for _, v := range backoffValues {
			v = strings.TrimSpace(v)
			if seconds, err := time.ParseDuration(v + "s"); err == nil {
				config.RetryBackoffs = append(config.RetryBackoffs, seconds)
			} else {
				fmt.Printf("Warning: invalid backoff value '%s', using defaults\n", v)
				config.RetryBackoffs = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
				break
			}
		}
	}

	if err := env.Parse(&config); err != nil {
		panic(fmt.Sprintf("error parsing env %v", err))
	}

	if !strings.HasPrefix(config.Address, "http://") {
		config.Address = "http://" + config.Address
	}

	return &config
}

package server

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func NewConfig() Config {
	return Config{
		Address:         "0.0.0.0:8080",
		StoreInterval:   300,
		FileStoragePath: "backup.json",
		Restore:         false,
	}
}

func NewConfigFromCLI() Config {
	config := NewConfig()

	address := flag.String("a", config.Address, "Address to listen on")
	storeInterval := flag.Int("i", config.StoreInterval, "Backup to file every n seconds")
	fileStoragePath := flag.String("f", config.FileStoragePath, "Backup file name")
	restore := flag.Bool("r", config.Restore, "Hydrate stats from file on startup? (true/false)")

	flag.Parse()

	config.Address = *address
	config.StoreInterval = *storeInterval
	config.FileStoragePath = *fileStoragePath
	config.Restore = *restore

	if err := env.Parse(&config); err != nil {
		panic(fmt.Sprintf("error parsing env: %v", err))
	}

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		config.Address = envAddr
	}

	return config
}

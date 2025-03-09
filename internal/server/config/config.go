package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"sync"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	PgDsn           string `env:"DATABASE_DSN"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			Address:         "0.0.0.0:8080",
			StoreInterval:   300,
			FileStoragePath: "backup.json",
			Restore:         false,
			PgDsn:           "",
		}

		address := flag.String("a", instance.Address, "Address to listen on")
		storeInterval := flag.Int("i", instance.StoreInterval, "Backup to file every n seconds")
		fileStoragePath := flag.String("f", instance.FileStoragePath, "Backup file name")
		restore := flag.Bool("r", instance.Restore, "Hydrate stats from file on startup? (true/false)")
		pgDsn := flag.String("d", instance.PgDsn, "Postgres DSN string")

		flag.Parse()

		instance.Address = *address
		instance.StoreInterval = *storeInterval
		instance.FileStoragePath = *fileStoragePath
		instance.Restore = *restore
		instance.PgDsn = *pgDsn

		if err := env.Parse(instance); err != nil {
			panic(fmt.Sprintf("error parsing env: %v", err))
		}
	})

	return instance
}

// helper for tests
func resetConfig() {
	instance = nil
	once = sync.Once{}
}

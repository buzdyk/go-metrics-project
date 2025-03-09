package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/server/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"sync"
)

type Client struct {
	dsn string
}

var (
	once     sync.Once
	instance *Client
)

func GetClient() *Client {
	once.Do(func() {
		dsn := config.GetConfig().PgDsn
		instance = &Client{dsn}
	})

	return instance
}

func (pg *Client) RunMigrations() error {
	fmt.Println(pg.dsn)
	pgsql, err := sql.Open("postgres", pg.dsn)
	if err != nil {
		fmt.Println("1")
		return err
	}
	defer pgsql.Close()

	cfg := &postgres.Config{}

	db, err := postgres.WithInstance(pgsql, cfg)

	if err != nil {
		fmt.Println("2")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres",
		db,
	)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create migration instance: %v", err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.New(fmt.Sprintf("Migration failed: %v", err))
	}

	return nil
}

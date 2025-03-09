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

func (pg *Client) DB() (*sql.DB, error) {
	db, err := sql.Open("postgres", pg.dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (pg *Client) RunMigrations() error {
	db, err := pg.DB()

	if err != nil {
		return err
	}
	defer db.Close()

	cfg := &postgres.Config{}

	driver, err := postgres.WithInstance(db, cfg)

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create migration instance: %v", err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.New(fmt.Sprintf("Migration failed: %v", err))
	}

	return nil
}

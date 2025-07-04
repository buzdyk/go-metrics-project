package database

import (
	"context"
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

func (pg *Client) Ping() error {
	ctx := context.Background()
	return pg.PingContext(ctx)
}

func (pg *Client) PingContext(ctx context.Context) error {
	db, err := pg.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.PingContext(ctx)
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
		"file://migrations",
		"postgres",
		driver,
	)

	if err != nil {
		return fmt.Errorf("failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %v", err)
	}

	return nil
}

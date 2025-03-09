package storage

import (
	"database/sql"
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/database"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"sync"
)

type PgStorage[T AllowedTypes] struct {
	client *database.Client
}

var mu2 *sync.Mutex

func (s *PgStorage[T]) Store(name string, v T) error {
	mu2.Lock()
	defer mu2.Unlock()

	db, err := s.client.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO "+s.table()+" (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value", name, v)
	return err
}

func (s *PgStorage[T]) StoreMany(m map[string]T) error {
	mu2.Lock()
	defer mu2.Unlock()

	db, err := s.client.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO " + s.table() + " (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for name, value := range m {
		if _, err := stmt.Exec(name, value); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PgStorage[T]) Values() (map[string]T, error) {
	mu.Lock()
	defer mu.Unlock()

	db, err := s.client.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, value FROM " + s.table())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]T)
	for rows.Next() {
		var name string
		var value T
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		data[name] = value
	}

	return data, rows.Err()
}

func (s *PgStorage[T]) Value(name string) (T, error) {
	mu.Lock()
	defer mu.Unlock()

	db, err := s.client.DB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var value T
	err = db.QueryRow("SELECT value FROM "+s.table()+" WHERE name = $1", name).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("unknown metric: " + name)
	}

	return value, err
}

func (s *PgStorage[T]) table() string {
	var zero T
	switch any(zero).(type) {
	case metrics.Gauge:
		return "gauges"
	case metrics.Counter:
		return "counters"
	default:
		panic("unsupported type")
	}
}

func NewPgStorage[T AllowedTypes](client *database.Client) *PgStorage[T] {
	return &PgStorage[T]{
		client: client,
	}
}

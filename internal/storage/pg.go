package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/database"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
)

type PgStorage[T AllowedTypes] struct {
	client *database.Client
}

func (s *PgStorage[T]) Store(ctx context.Context, name string, v T) error {

	db, err := s.client.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(SQLInsertOrUpdate, s.table())
	_, err = db.ExecContext(ctx, query, name, v)
	return err
}

func (s *PgStorage[T]) StoreMany(ctx context.Context, m map[string]T) error {

	db, err := s.client.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf(SQLInsertOrUpdate, s.table())
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for name, value := range m {
		if _, err := stmt.ExecContext(ctx, name, value); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PgStorage[T]) Values(ctx context.Context) (map[string]T, error) {

	db, err := s.client.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(SQLSelectAll, s.table())
	rows, err := db.QueryContext(ctx, query)
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

func (s *PgStorage[T]) Value(ctx context.Context, name string) (T, error) {

	db, err := s.client.DB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var value T
	query := fmt.Sprintf(SQLSelectByName, s.table())
	err = db.QueryRowContext(ctx, query, name).Scan(&value)
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

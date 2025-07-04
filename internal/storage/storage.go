package storage

import (
	"context"
	"github.com/buzdyk/go-metrics-project/internal/models"
)

type AllowedTypes interface {
	models.Gauge | models.Counter
}

type Storage[T AllowedTypes] interface {
	Store(ctx context.Context, name string, value T) error
	StoreMany(ctx context.Context, m map[string]T) error
	Value(ctx context.Context, name string) (T, error)
	Values(ctx context.Context) (map[string]T, error)
}

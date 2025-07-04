package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/models"
	"os"
	"sync"
)

var mu sync.Mutex

type FileEntry struct {
	Type    string          `json:"type"`
	Counter models.Counter  `json:"counter"`
	Gauge   models.Gauge    `json:"gauge"`
}

type FileStorage[T AllowedTypes] struct {
	filepath string
}

func NewFileStorage[T AllowedTypes](filepath string) *FileStorage[T] {
	return &FileStorage[T]{filepath: filepath}
}

func (b *FileStorage[T]) StoreMany(ctx context.Context, m map[string]T) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := b.readFile(ctx)
	if err != nil {
		return err
	}

	var zero T
	switch any(zero).(type) {
	case models.Gauge:
		for name, value := range m {
			data[name] = FileEntry{
				Type:  models.GaugeName,
				Gauge: models.Gauge(value),
			}
		}
	case models.Counter:
		for name, value := range m {
			data[name] = FileEntry{
				Type:    models.CounterName,
				Counter: models.Counter(value),
			}
		}
	}

	file, err := os.OpenFile(b.filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}

func (b *FileStorage[T]) Values(ctx context.Context) (map[string]T, error) {
	data, err := b.readFile(ctx)

	if err != nil {
		return nil, err
	}

	m := make(map[string]T)

	var zero T
	switch any(zero).(type) {
	case models.Gauge:
		for name, entry := range data {
			if entry.Type == models.GaugeName {
				m[name] = T(entry.Gauge)
			}
		}
	case models.Counter:
		for name, entry := range data {
			if entry.Type == models.CounterName {
				m[name] = T(entry.Counter)
			}
		}
	}

	return m, nil
}

func (b *FileStorage[T]) Value(ctx context.Context, name string) (T, error) {
	var zero T

	data, err := b.Values(ctx)

	if err != nil {
		return zero, err
	}

	if v, ok := data[name]; !ok {
		return zero, errors.New("value not found")
	} else {
		return v, nil
	}
}

func (b *FileStorage[T]) Store(ctx context.Context, name string, value T) error {
	if err := b.StoreMany(ctx, map[string]T{name: value}); err != nil {
		return err
	}

	return nil
}

func (b *FileStorage[T]) readFile(ctx context.Context) (map[string]FileEntry, error) {
	// Check if context is canceled before performing I/O
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Continue with file operations
	}

	file, err := os.Open(b.filepath)
	if err != nil {
		return make(map[string]FileEntry), nil
	}
	defer file.Close()

	var data map[string]FileEntry
	if err = json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

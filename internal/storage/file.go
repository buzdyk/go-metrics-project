package storage

import (
	"encoding/json"
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"os"
	"sync"
)

var mu sync.Mutex

type FileEntry struct {
	Type    string          `json:"type"`
	Counter metrics.Counter `json:"counter"`
	Gauge   metrics.Gauge   `json:"gauge"`
}

type FileStorage[T AllowedTypes] struct {
	filepath string
}

func NewFileStorage[T AllowedTypes](filepath string) *FileStorage[T] {
	return &FileStorage[T]{filepath: filepath}
}

func (b *FileStorage[T]) StoreMany(m map[string]T) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := b.readFile()
	if err != nil {
		return err
	}

	var zero T
	switch any(zero).(type) {
	case metrics.Gauge:
		for name, value := range m {
			data[name] = FileEntry{
				Type:  metrics.GaugeName,
				Gauge: metrics.Gauge(value),
			}
		}
	case metrics.Counter:
		for name, value := range m {
			data[name] = FileEntry{
				Type:    metrics.CounterName,
				Counter: metrics.Counter(value),
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

func (b *FileStorage[T]) Values() (map[string]T, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := b.readFile()

	if err != nil {
		return nil, err
	}

	m := make(map[string]T)

	var zero T
	switch any(zero).(type) {
	case metrics.Gauge:
		for name, entry := range data {
			if entry.Type == metrics.GaugeName {
				m[name] = T(entry.Gauge)
			}
		}
	case metrics.Counter:
		for name, entry := range data {
			if entry.Type == metrics.CounterName {
				m[name] = T(entry.Counter)
			}
		}
	}

	return m, nil
}

func (b *FileStorage[T]) Value(name string) (T, error) {
	var zero T

	data, err := b.Values()

	if err != nil {
		return zero, err
	}

	if v, ok := data[name]; !ok {
		return zero, errors.New("value not found")
	} else {
		return v, nil
	}
}

func (b *FileStorage[T]) Store(name string, value T) error {
	if err := b.StoreMany(map[string]T{name: value}); err != nil {
		return err
	}

	return nil
}

func (b *FileStorage[T]) readFile() (map[string]FileEntry, error) {
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

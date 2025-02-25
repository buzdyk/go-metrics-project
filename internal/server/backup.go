package server

import (
	"encoding/json"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"os"
)

type Backup struct {
	filepath     string
	gaugeStore   storage.Storage[metrics.Gauge]
	counterStore storage.Storage[metrics.Counter]
}

type BackupEntry struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

func (b *Backup) Backup() error {
	var backupData []BackupEntry

	for n, v := range b.gaugeStore.Values() {
		backupData = append(backupData, BackupEntry{
			Name:  n,
			Value: v,
			Type:  "gauge",
		})
	}

	for n, v := range b.counterStore.Values() {
		backupData = append(backupData, BackupEntry{
			Name:  n,
			Value: v,
			Type:  "counter",
		})
	}

	file, err := os.OpenFile(b.filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err2 := json.NewEncoder(file).Encode(backupData); err2 != nil {
		return err2
	}

	return nil
}

func (b *Backup) Restore() error {
	file, err := os.Open(b.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var backupData []BackupEntry
	if err := json.NewDecoder(file).Decode(&backupData); err != nil {
		return err
	}

	for _, entry := range backupData {
		switch entry.Type {
		case "gauge":
			if value, ok := entry.Value.(float64); ok {
				b.gaugeStore.Store(entry.Name, metrics.Gauge(value))
			}
		case "counter":
			if value, ok := entry.Value.(float64); ok {
				b.counterStore.Store(entry.Name, metrics.Counter(value))
			}
		}
	}

	return nil
}

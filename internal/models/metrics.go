package models

// Gauge represents a floating point metric value
type Gauge float64

// Counter represents an integer metric counter
type Counter uint64

// MetricType constants
const (
	GaugeName   = "gauge"
	CounterName = "counter"
)

// Metric represents a single metric with its metadata and value
type Metric struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // counter or gauge
	Delta *Counter `json:"delta,omitempty"` // value for counter
	Value *Gauge   `json:"value,omitempty"` // value for gauge
}

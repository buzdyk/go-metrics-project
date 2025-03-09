package metrics

import (
	"math/rand"
	"reflect"
	"runtime"
)

type Gauge float64
type Counter uint64

const GaugeName = "gauge"
const CounterName = "counter"

type Metric struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // counter or gauge
	Delta *Counter `json:"delta,omitempty"` // value for counter
	Value *Gauge   `json:"value,omitempty"` // value for gauge
}

var memStats = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

type Collector struct {
	pollCount   Counter
	randomValue Gauge
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(out map[string]any) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	r := reflect.ValueOf(m)

	for _, stat := range memStats {
		field := r.FieldByName(stat)
		switch field.Kind() {
		case reflect.Uint64, reflect.Uint32:
			out[stat] = Gauge(field.Uint())
		case reflect.Float64:
			out[stat] = Gauge(field.Float())
		default:
		}
	}

	c.pollCount += 1
	c.randomValue = Gauge(rand.NormFloat64())

	out["PollCount"] = c.pollCount
	out["RandomValue"] = c.randomValue
}

func Exists(metric string) bool {
	return metric != "unknown"
}

func IsValidType(t string) bool {
	return t == GaugeName || t == CounterName
}

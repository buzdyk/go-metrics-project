package collector

import (
	"math/rand"
	"reflect"
	"runtime"
	"sync"
)

// List of memory statistics to collect
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

// Collector collects runtime metrics
type Collector struct {
	pollCount   Counter
	randomValue Gauge
	mu          sync.Mutex
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{}
}

// Collect gathers runtime metrics and stores them in the provided map
func (c *Collector) Collect(out map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()

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

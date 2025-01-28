package metrics

import (
	"math/rand"
	"reflect"
	"runtime"
)

type Gauge float64
type Counter uint64

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
	pollCount Counter
}

func (c *Collector) Collect(out map[string]interface{}) {
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

	out["PollCount"] = c.pollCount
	out["RandomValue"] = Gauge(rand.NormFloat64())
}

func Exists(metric string) bool {
	if metric == "PollCount" || metric == "RandomValue" {
		return true
	}

	for _, v := range memStats {
		if v == metric {
			return true
		}
	}

	return false
}

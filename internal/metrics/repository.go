package metrics

import (
	"fmt"
	"runtime"
)

var gauges map[string]*Gauge

func init() {
	gauges = make(map[string]*Gauge)

	c := Collector{func() (Gauge, error) {
		var m *runtime.MemStats
		runtime.ReadMemStats(m)

		return Gauge{
			ID:    "alloc",
			Name:  "alloc",
			Value: float64(m.Alloc),
		}, nil
	}}

	fmt.Println(c)
}

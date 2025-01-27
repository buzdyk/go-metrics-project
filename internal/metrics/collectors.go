package metrics

import "runtime"

var Collectors map[string]func() Gauge

func init() {
	Collectors = map[string]func() Gauge{
		"Alloc":       Alloc,
		"BuckHashSys": BuckHashSys,
		"Frees":       Frees,
	}
}

func Alloc() Gauge {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return Gauge(m.Alloc)
}

func BuckHashSys() Gauge {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return Gauge(m.BuckHashSys)
}

func Frees() Gauge {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return Gauge(m.Frees)
}

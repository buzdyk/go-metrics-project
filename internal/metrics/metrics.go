package metrics

import (
	rt "runtime/metrics"
)

func init() {
	rt.All()
}

type Gauge struct {
	ID    string
	Name  string
	Value float64
}

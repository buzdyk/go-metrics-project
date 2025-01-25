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

type Counter struct {
	id    string
	value int64
}

//func (g *Gauge) ID() string {
//	return g.id
//}
//
//func (g *Gauge) Value() float64 {
//	return rand.ExpFloat64()
//}

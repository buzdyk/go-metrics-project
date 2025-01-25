package metrics

type metrics interface {
	ID() string
	Value() float64
}

type Gauge struct{}

type Counter struct{}

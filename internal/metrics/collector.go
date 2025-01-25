package metrics

type Collector struct {
	Collect func() (Gauge, error)
}

type Collectors []Collector

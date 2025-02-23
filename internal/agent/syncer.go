package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"net/http"
)

type UnknownTypeError struct{}

func (t UnknownTypeError) Error() string {
	return "unknown variable type"
}

type HTTPSyncer struct {
	Host string
}

func NewHTTPSyncer(host string) *HTTPSyncer {
	return &HTTPSyncer{
		Host: host,
	}
}

func (hc *HTTPSyncer) SyncMetric(id string, value any) (*http.Response, error) {
	switch v := value.(type) {
	case metrics.Gauge:
		res, err := hc.syncGauge(id, v)
		if err != nil {
			return nil, err
		}

		return res, nil
	case metrics.Counter:
		res, err := hc.syncCounter(id, v)

		if err != nil {
			return nil, err
		}

		return res, nil
	default:
		return nil, UnknownTypeError{}
	}
}

func (hc *HTTPSyncer) syncGauge(name string, g metrics.Gauge) (*http.Response, error) {
	metric := metrics.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &g,
	}

	data, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%v/update", hc.Host)
	res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (hc *HTTPSyncer) syncCounter(name string, c metrics.Counter) (*http.Response, error) {
	metric := metrics.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &c,
	}

	data, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%v/update", hc.Host)
	res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return res, nil
}

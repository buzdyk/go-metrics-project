package agent

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"net/http"
)

type UnknownTypeError struct{}

func (t UnknownTypeError) Error() string {
	return "unknown variable type"
}

type RealHTTPClient struct {
	Host string
}

func NewHTTPClient(host string) *RealHTTPClient {
	return &RealHTTPClient{
		Host: host,
	}
}

func (hc *RealHTTPClient) Post(id string, value any) (*http.Response, error) {
	switch v := value.(type) {
	case metrics.Gauge:
		res, err := hc.postGauge(id, v)
		if err != nil {
			return nil, err
		}

		return res, nil
	case metrics.Counter:
		res, err := hc.postCounter(id, v)

		if err != nil {
			return nil, err
		}

		return res, nil
	default:
		return nil, UnknownTypeError{}
	}
}

func (hc *RealHTTPClient) postGauge(name string, g metrics.Gauge) (*http.Response, error) {
	endpoint := fmt.Sprintf("%v/update/gauge/%v/%v", hc.Host, name, g)
	if res, err := http.Post(endpoint, "text/plain", nil); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (hc *RealHTTPClient) postCounter(name string, c metrics.Counter) (*http.Response, error) {
	endpoint := fmt.Sprintf("%v/update/counter/%v/%v", hc.Host, name, c)
	if res, err := http.Post(endpoint, "text/plain", nil); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

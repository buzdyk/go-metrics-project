package agent

import (
	"errors"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"log"
	"net/http"
)

type HttpClient interface {
	Post(name string, value interface{}) (*http.Response, error)
}

type RealHttpClient struct {
	Endpoint string
	Port     int
}

func (hc *RealHttpClient) Post(id string, value interface{}) (*http.Response, error) {
	switch v := value.(type) {
	case metrics.Gauge:
		res, err := hc.PostGauge(id, v)
		if err != nil {
			return nil, err
		}

		return res, nil
	case metrics.Counter:
		res, err := hc.PostCounter(id, v)

		if err != nil {
			return nil, err
		}

		return res, nil
	default:
		log.Printf("WARNING: Unknown type %v %t %T\n", v, v, v)
		return nil, errors.New("unknown type in real http client")
	}
}

func (hc *RealHttpClient) PostGauge(name string, g metrics.Gauge) (*http.Response, error) {
	endpoint := fmt.Sprintf("%v:%v/update/gauge/%v/%v", hc.Endpoint, hc.Port, name, g)
	if res, err := http.Post(endpoint, "text-plain", nil); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (hc *RealHttpClient) PostCounter(name string, c metrics.Counter) (*http.Response, error) {
	endpoint := fmt.Sprintf("%v:%v/update/counter/%v/%v", hc.Endpoint, hc.Port, name, c)
	if res, err := http.Post(endpoint, "text-plain", nil); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

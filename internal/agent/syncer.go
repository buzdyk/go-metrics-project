package agent

import (
	"bytes"
	"compress/gzip"
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
	endpoint := fmt.Sprintf("%v/update/gauge/%v/%v", hc.Host, name, g)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	return client.Do(req)
}

func (hc *HTTPSyncer) syncCounter(name string, c metrics.Counter) (*http.Response, error) {
	endpoint := fmt.Sprintf("%v/update/counter/%v/%v", hc.Host, name, c)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	return client.Do(req)
}

func (hc *HTTPSyncer) SyncMetrics(ms []metrics.Metric) (*http.Response, error) {
	endpoint := fmt.Sprintf("%v/updates/", hc.Host)

	jsonData, err := json.Marshal(ms)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	var buf bytes.Buffer

	gzipWriter := gzip.NewWriter(&buf)
	if _, err := gzipWriter.Write(jsonData); err != nil {
		return nil, err
	}
	defer gzipWriter.Close()

	req, err := http.NewRequest("POST", endpoint, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

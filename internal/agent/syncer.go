package agent

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/models"
	"net/http"
)

type UnknownTypeError struct{}

func (t UnknownTypeError) Error() string {
	return "unknown variable type"
}

type HTTPSyncer struct {
	Host string
	Key  string
}

func NewHTTPSyncer(host string, key string) *HTTPSyncer {
	return &HTTPSyncer{
		Host: host,
		Key:  key,
	}
}

func (hc *HTTPSyncer) calculateHash(value string) string {
	if hc.Key == "" {
		return ""
	}
	
	h := sha256.New()
	h.Write([]byte(value + hc.Key))
	return hex.EncodeToString(h.Sum(nil))
}

func (hc *HTTPSyncer) SyncMetric(id string, value any) (*http.Response, error) {
	switch v := value.(type) {
	case models.Gauge:
		res, err := hc.syncGauge(id, v)
		if err != nil {
			return nil, err
		}

		return res, nil
	case models.Counter:
		res, err := hc.syncCounter(id, v)

		if err != nil {
			return nil, err
		}

		return res, nil
	default:
		return nil, UnknownTypeError{}
	}
}

func (hc *HTTPSyncer) syncGauge(name string, g models.Gauge) (*http.Response, error) {
	gaugeValue := fmt.Sprintf("%v", g)
	endpoint := fmt.Sprintf("%v/update/gauge/%v/%v", hc.Host, name, gaugeValue)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain")
	
	// Calculate hash of the string value with the key and set the header
	hashValue := hc.calculateHash(gaugeValue)
	if hashValue != "" {
		req.Header.Set("HashSHA256", hashValue)
	}

	client := &http.Client{}
	return client.Do(req)
}

func (hc *HTTPSyncer) syncCounter(name string, c models.Counter) (*http.Response, error) {
	counterValue := fmt.Sprintf("%v", c)
	endpoint := fmt.Sprintf("%v/update/counter/%v/%v", hc.Host, name, counterValue)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain")
	
	// Calculate hash of the string value with the key and set the header
	hashValue := hc.calculateHash(counterValue)
	if hashValue != "" {
		req.Header.Set("HashSHA256", hashValue)
	}

	client := &http.Client{}
	return client.Do(req)
}

func (hc *HTTPSyncer) SyncMetrics(ms []models.Metric) (*http.Response, error) {
	endpoint := fmt.Sprintf("%v/updates/", hc.Host)

	jsonData, err := json.Marshal(ms)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	var buf bytes.Buffer

	if _, err := buf.Write(jsonData); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", endpoint, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	// Calculate hash of the JSON payload with the key and set the header
	hashValue := hc.calculateHash(string(jsonData))
	if hashValue != "" {
		req.Header.Set("HashSHA256", hashValue)
	}

	return client.Do(req)
}

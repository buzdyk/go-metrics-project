package main

import (
	"github.com/buzdyk/go-metrics-project/internal/agent"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"io"
	"net/http"
)

type RealHttpClient struct{}

func (hc *RealHttpClient) Post(endpoint, contentType string, body io.Reader) (res *http.Response, err error) {
	res, err = http.Post(endpoint, contentType, body)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	a, err := agent.NewAgent(metrics.Collectors, &RealHttpClient{})

	if err != nil {
		panic(err)
	}

	a.Run()
}

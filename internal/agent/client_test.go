package agent

import (
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestRealHttpClient_Post(t *testing.T) {
	var counter int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&counter, 1)

		switch counter {
		case 1:
			assert.Equal(t, "/update/gauge/metric/42", r.URL.Path)
			assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		case 2:
			assert.Equal(t, "/update/counter/metric/42", r.URL.Path)
			assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		case 3:
			t.Error("expected 2 http requests, got 3rd")
		}
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := &RealHttpClient{
		Host: server.URL,
	}

	client.Post("metric", metrics.Gauge(42))
	client.Post("metric", metrics.Counter(42))

	_, err := client.Post("metric", int64(20))

	var want UnknownTypeError
	assert.True(t, errors.As(err, &want))

	require.Equal(t, counter, int64(2), "expected 2 http requests to be made, got %v", counter)
}

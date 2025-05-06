package agent

import (
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestHTTPSyncer_SyncMetric(t *testing.T) {
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

	syncer := &HTTPSyncer{
		Host: server.URL,
	}

	var (
		r   *http.Response
		err error
	)

	r, _ = syncer.SyncMetric("metric", models.Gauge(42))
	r.Body.Close()
	r, _ = syncer.SyncMetric("metric", models.Counter(42))
	r.Body.Close()
	r, err = syncer.SyncMetric("metric", int64(20))
	if r != nil {
		r.Body.Close()
	}

	var want UnknownTypeError
	assert.True(t, errors.As(err, &want))

	require.Equal(t, counter, int64(2), "expected 2 http requests to be made, got %v", counter)
}

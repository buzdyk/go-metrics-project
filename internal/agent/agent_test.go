package agent

import (
	"bytes"
	"context"
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPSyncer mocks the Syncer interface
type MockHTTPSyncer struct {
	mock.Mock
}

func (m *MockHTTPSyncer) SyncMetric(name string, value any) (*http.Response, error) {
	args := m.Called("SyncMetric", name, value)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

func (m *MockHTTPSyncer) SyncMetrics(ms []metrics.Metric) (*http.Response, error) {
	args := m.Called("SyncMetrics", ms)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

// MockMetricsCollector mocks the MetricsCollector interface
type MockMetricsCollector struct {
	mock.Mock
}

func (m *MockMetricsCollector) Collect(out map[string]any) {
	m.Called(out)
}

// TestNewAgent checks initialization
func TestNewAgent(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockSyncer := new(MockHTTPSyncer)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}

	agent := NewAgent(config, mockCollector, mockSyncer)

	assert.NotNil(t, agent)
	assert.Equal(t, *config, agent.config)
	assert.Equal(t, mockCollector, agent.collector)
	assert.Equal(t, mockSyncer, agent.syncer)
	assert.NotNil(t, agent.data)
}

// TestCollect ensures collect() correctly calls Collect()
func TestAgentCollect(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockSyncer := new(MockHTTPSyncer)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockSyncer)

	mockCollector.On("Collect", mock.Anything).Once()

	agent.collect()

	mockCollector.AssertCalled(t, "Collect", agent.data)
}

// TestSync ensures sync() correctly sends HTTP requests
func TestAgentSync(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockSyncer := new(MockHTTPSyncer)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockSyncer)

	agent.data = map[string]any{
		"metric1": 100,
		"metric2": 200,
	}

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}

	mockSyncer.On("SyncMetric", "metric1", 100).Return(mockResponse, nil).Once()
	mockSyncer.On("SyncMetric", "metric2", 200).Return(mockResponse, nil).Once()

	agent.sync()

	mockSyncer.AssertCalled(t, "SyncMetric", "metric1", 100)
	mockSyncer.AssertCalled(t, "SyncMetric", "metric2", 200)
}

// TestSyncErrorHandling ensures sync() handles HTTP syncer errors properly
func TestAgentSyncWithErrors(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockSyncer := new(MockHTTPSyncer)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockSyncer)

	agent.data = map[string]any{
		"metric1": 100,
		"metric2": 200,
	}

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}

	mockSyncer.On("SyncMetric", "metric1", 100).Return(mockResponse, nil).Once()
	mockSyncer.On("SyncMetric", "metric2", 200).Return(nil, errors.New("network error")).Once()

	agent.sync()

	mockSyncer.AssertCalled(t, "SyncMetric", "metric1", 100)
	mockSyncer.AssertCalled(t, "SyncMetric", "metric2", 200)
}

// TestSyncEmptyData ensures sync() does nothing if there are no metrics
func TestAgentSyncWithEmptyData(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockSyncer := new(MockHTTPSyncer)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockSyncer)

	agent.sync()

	mockSyncer.AssertNotCalled(t, "SyncMetric", mock.Anything, mock.Anything)
}

// TestConcurrency checks collect() and sync() can run concurrently
func TestAgentConcurrency(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockSyncer := new(MockHTTPSyncer)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockSyncer)

	mockCollector.On("Collect", mock.Anything).Maybe()
	mockSyncer.On("SyncMetric", mock.Anything, mock.Anything).Maybe()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		agent.collect()
	}()

	go func() {
		defer wg.Done()
		agent.sync()
	}()

	wg.Wait()
}

// TestRun ensures Run() triggers collect() and sync() at intervals
func TestAgentRun(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockSyncer := new(MockHTTPSyncer)
	config := &Config{Address: "http://localhost:8080", Report: 1, Collect: 1}
	agent := NewAgent(config, mockCollector, mockSyncer)
	agent.data["value"] = metrics.Gauge(1)

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}

	mockCollector.On("Collect", mock.Anything).Maybe()
	mockSyncer.On("SyncMetric", mock.Anything, mock.Anything).Maybe().Return(mockResponse, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	go agent.Run(context.Background())

	go func() {
		defer wg.Done()
		time.Sleep(2*time.Second + 100*time.Millisecond)
	}()

	wg.Wait()

	mockCollector.AssertNumberOfCalls(t, "Collect", 2)
	mockSyncer.AssertNumberOfCalls(t, "SyncMetric", 2)
}

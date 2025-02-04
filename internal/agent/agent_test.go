package agent

import (
	"bytes"
	"errors"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient mocks the HTTPClient interface
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(name string, value interface{}) (*http.Response, error) {
	args := m.Called(name, value)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

// MockMetricsCollector mocks the MetricsCollector interface
type MockMetricsCollector struct {
	mock.Mock
}

func (m *MockMetricsCollector) Collect(out map[string]interface{}) {
	m.Called(out)
}

// TestNewAgent checks initialization
func TestNewAgent(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockClient := new(MockHTTPClient)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}

	agent := NewAgent(config, mockCollector, mockClient)

	assert.NotNil(t, agent)
	assert.Equal(t, *config, agent.config)
	assert.Equal(t, mockCollector, agent.collector)
	assert.Equal(t, mockClient, agent.client)
	assert.NotNil(t, agent.data)
}

// TestCollect ensures collect() correctly calls Collect()
func TestAgentCollect(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockClient := new(MockHTTPClient)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockClient)

	mockCollector.On("Collect", mock.Anything).Once()

	agent.collect()

	mockCollector.AssertCalled(t, "Collect", agent.data)
}

// TestSync ensures sync() correctly sends HTTP requests
func TestAgentSync(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockClient := new(MockHTTPClient)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockClient)

	agent.data = map[string]interface{}{
		"metric1": 100,
		"metric2": 200,
	}

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}

	mockClient.On("Post", "metric1", 100).Return(mockResponse, nil).Once()
	mockClient.On("Post", "metric2", 200).Return(mockResponse, nil).Once()

	agent.sync()

	mockClient.AssertCalled(t, "Post", "metric1", 100)
	mockClient.AssertCalled(t, "Post", "metric2", 200)
}

// TestSyncErrorHandling ensures sync() handles HTTP client errors properly
func TestAgentSyncWithErrors(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockClient := new(MockHTTPClient)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockClient)

	agent.data = map[string]interface{}{
		"metric1": 100,
		"metric2": 200,
	}

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
	}

	mockClient.On("Post", "metric1", 100).Return(mockResponse, nil).Once()
	mockClient.On("Post", "metric2", 200).Return(nil, errors.New("network error")).Once()

	agent.sync()

	mockClient.AssertCalled(t, "Post", "metric1", 100)
	mockClient.AssertCalled(t, "Post", "metric2", 200)
}

// TestSyncEmptyData ensures sync() does nothing if there are no metrics
func TestAgentSyncWithEmptyData(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockClient := new(MockHTTPClient)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockClient)

	agent.sync()

	mockClient.AssertNotCalled(t, "Post", mock.Anything, mock.Anything)
}

// TestConcurrency checks collect() and sync() can run concurrently
func TestAgentConcurrency(t *testing.T) {
	mockCollector := new(MockMetricsCollector)
	mockClient := new(MockHTTPClient)
	config := &Config{Address: "http://localhost:8080", Report: 10, Collect: 2}
	agent := NewAgent(config, mockCollector, mockClient)

	mockCollector.On("Collect", mock.Anything).Maybe()
	mockClient.On("Post", mock.Anything, mock.Anything).Maybe()

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
	mockClient := new(MockHTTPClient)
	config := &Config{Address: "http://localhost:8080", Report: 1, Collect: 1}
	agent := NewAgent(config, mockCollector, mockClient)
	agent.data["value"] = metrics.Gauge(1)

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}

	mockCollector.On("Collect", mock.Anything).Maybe()
	mockClient.On("Post", mock.Anything, mock.Anything).Maybe().Return(mockResponse, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	go agent.Run()

	go func() {
		defer wg.Done()
		time.Sleep(2*time.Second + 100*time.Millisecond)
	}()

	wg.Wait()

	mockCollector.AssertNumberOfCalls(t, "Collect", 2)
	mockClient.AssertNumberOfCalls(t, "Post", 2)
}

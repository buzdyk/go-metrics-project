package metrics

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCollector_Collect ensures that Collect() populates memory stats correctly
func TestCollector_Collect(t *testing.T) {
	collector := &Collector{}
	out := make(map[string]interface{})

	collector.Collect(out)

	for _, stat := range memStats {
		assert.Contains(t, out, stat, "Missing metric: %s", stat)
	}

	assert.Contains(t, out, "PollCount", "PollCount should be present")
	assert.Contains(t, out, "RandomValue", "RandomValue should be present")

	initialPollCount := collector.pollCount
	collector.Collect(out)
	assert.Equal(t, initialPollCount+1, out["PollCount"], "PollCount did not increment")
}

// TestCollector_RandomValue ensures RandomValue is a float64
func TestCollector_RandomValue(t *testing.T) {
	collector := &Collector{}
	out := make(map[string]interface{})

	collector.Collect(out)

	_, ok := out["RandomValue"].(Gauge)
	assert.True(t, ok, "RandomValue should be a float64")
}

// TestCollector_PollCountIncrements ensures PollCount increments correctly
func TestCollector_PollCountIncrements(t *testing.T) {
	collector := &Collector{}
	out := make(map[string]interface{})

	collector.Collect(out)
	firstPollCount := collector.pollCount

	collector.Collect(out)
	secondPollCount := collector.pollCount

	assert.Equal(t, firstPollCount+1, secondPollCount, "PollCount did not increment properly")
}

// TestExists_ValidMetrics ensures Exists() returns true for valid metrics
func TestExists_ValidMetrics(t *testing.T) {
	for _, metric := range memStats {
		assert.True(t, Exists(metric), "Exists should return true for valid metric: %s", metric)
	}
	assert.True(t, Exists("PollCount"), "PollCount should exist")
	assert.True(t, Exists("RandomValue"), "RandomValue should exist")
}

// TestExists_InvalidMetrics ensures Exists() returns false for unknown metrics
func TestExists_InvalidMetrics(t *testing.T) {
	assert.False(t, Exists("UnknownMetric"), "Exists should return false for unknown metric")
}

// TestCollector_CorrectDataTypes ensures all collected values are of correct types
func TestCollector_CorrectDataTypes(t *testing.T) {
	collector := &Collector{}
	out := make(map[string]interface{})
	collector.Collect(out)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	r := reflect.ValueOf(mem)

	for _, stat := range memStats {
		field := r.FieldByName(stat)
		switch field.Kind() {
		case reflect.Uint64, reflect.Uint32, reflect.Float64:
			assert.IsType(t, Gauge(0), out[stat], "Metric %s has an unexpected type %T", stat, stat)
		}
	}
}

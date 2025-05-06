package collector

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/buzdyk/go-metrics-project/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestCollector_Collect ensures that Collect() populates memory stats correctly
func TestCollector_Collect(t *testing.T) {
	collector := NewCollector()
	out := make(map[string]any)

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

// TestCollector_RandomValue ensures RandomValue is a Gauge type
func TestCollector_RandomValue(t *testing.T) {
	collector := NewCollector()
	out := make(map[string]any)

	collector.Collect(out)

	_, ok := out["RandomValue"].(models.Gauge)
	assert.True(t, ok, "RandomValue should be a Gauge type")
}

// TestCollector_PollCountIncrements ensures PollCount increments correctly
func TestCollector_PollCountIncrements(t *testing.T) {
	collector := NewCollector()
	out := make(map[string]any)

	collector.Collect(out)
	firstPollCount := collector.pollCount

	collector.Collect(out)
	secondPollCount := collector.pollCount

	assert.Equal(t, firstPollCount+1, secondPollCount, "PollCount did not increment properly")
}

// TestCollector_CorrectDataTypes ensures all collected values are of correct types
func TestCollector_CorrectDataTypes(t *testing.T) {
	collector := NewCollector()
	out := make(map[string]any)
	collector.Collect(out)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	r := reflect.ValueOf(mem)

	for _, stat := range memStats {
		field := r.FieldByName(stat)
		switch field.Kind() {
		case reflect.Uint64, reflect.Uint32, reflect.Float64:
			assert.IsType(t, models.Gauge(0), out[stat], "Metric %s has an unexpected type %T", stat, out[stat])
		}
	}
}

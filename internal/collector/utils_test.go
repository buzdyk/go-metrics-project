package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.False(t, Exists("unknown"), "Exists should return false for unknown metric")
}

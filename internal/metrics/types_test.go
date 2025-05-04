package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsValidType ensures that metric type validation works correctly
func TestIsValidType(t *testing.T) {
	assert.True(t, IsValidType(GaugeName), "Gauge should be a valid metric type")
	assert.True(t, IsValidType(CounterName), "Counter should be a valid metric type")
	assert.False(t, IsValidType("invalid"), "Invalid type should return false")
}
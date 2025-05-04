package collector

import (
	"testing"

	"github.com/buzdyk/go-metrics-project/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestCollectSystem ensures that system metrics collection works
func TestCollectSystem(t *testing.T) {
	collector := NewCollector()
	out := make(map[string]any)

	collector.CollectSystem(out)

	// Check that some system metrics were collected
	assert.Contains(t, out, "TotalMemory", "TotalMemory metric should be present")
	assert.Contains(t, out, "FreeMemory", "FreeMemory metric should be present")

	// CPU utilization metrics might vary by system, so we just check the types
	for key, value := range out {
		if key != "TotalMemory" && key != "FreeMemory" {
			_, ok := value.(models.Gauge)
			assert.True(t, ok, "System metric %s should be of type Gauge", key)
		}
	}
}

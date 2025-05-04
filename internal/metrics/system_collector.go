package metrics

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

// CollectSystem gathers system metrics and stores them in the provided map
func (c *Collector) CollectSystem(out map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	collectMemoryMetrics(out)
	collectCPUMetrics(out)
}

// collectMemoryMetrics gathers memory-related system metrics
func collectMemoryMetrics(out map[string]any) {
	vmStat, err := mem.VirtualMemory()
	if err == nil {
		out["TotalMemory"] = Gauge(vmStat.Total)
		out["FreeMemory"] = Gauge(vmStat.Free)
	}
}

// collectCPUMetrics gathers CPU-related system metrics
func collectCPUMetrics(out map[string]any) {
	percentages, err := cpu.Percent(time.Second, true)
	if err == nil {
		for i, cpuPercent := range percentages {
			out[fmt.Sprintf("CPUutilization%d", i+1)] = Gauge(cpuPercent)
		}
	}
}
package metrics

// Exists checks if a metric name is valid
func Exists(metric string) bool {
	return metric != "unknown"
}

// IsValidType checks if a metric type is valid
func IsValidType(t string) bool {
	return t == GaugeName || t == CounterName
}
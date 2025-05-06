package storage

// SQL queries for metric storage operations
const (
	// Insert or update a metric
	SQLInsertOrUpdate = "INSERT INTO %s (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value"

	// Select all metrics
	SQLSelectAll = "SELECT name, value FROM %s"

	// Select a metric by name
	SQLSelectByName = "SELECT value FROM %s WHERE name = $1"
)
package agent

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewConfig checks the default values of NewConfig
func TestNewConfig(t *testing.T) {
	config := NewConfig()

	assert.Equal(t, "0.0.0.0:8080", config.Address)
	assert.Equal(t, 10, config.Report)
	assert.Equal(t, 2, config.Collect)
}

// TestNewConfigFromCLI_Defaults ensures CLI args are parsed correctly
func TestNewConfigFromCLI_Defaults(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{"cmd", "-a", "127.0.0.1:9000", "-r", "15", "-p", "5"}
	config := NewConfigFromCLI()

	assert.Equal(t, "http://127.0.0.1:9000", config.Address)
	assert.Equal(t, 15, config.Report)
	assert.Equal(t, 5, config.Collect)
}

// TestNewConfigFromCLI_EnvVariables ensures environment variables override CLI args
func TestNewConfigFromCLI_EnvVariables(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Setenv("ADDRESS", "192.168.1.100:9090")
	os.Setenv("REPORT", "20")
	os.Setenv("COLLECT", "7")
	defer os.Clearenv()

	os.Args = []string{"cmd", "-a", "127.0.0.1:9000", "-r", "15", "-p", "5"}
	config := NewConfigFromCLI()

	assert.Equal(t, "http://192.168.1.100:9090", config.Address)
	assert.Equal(t, 20, config.Report)
	assert.Equal(t, 7, config.Collect)
}

// TestNewConfigFromCLI_EnvVariables_MissingValues ensures missing env vars don't override defaults
func TestNewConfigFromCLI_EnvVariables_MissingValues(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Setenv("ADDRESS", "192.168.1.100:9090")
	defer os.Clearenv()

	os.Args = []string{"cmd"}
	config := NewConfigFromCLI()

	assert.Equal(t, "http://192.168.1.100:9090", config.Address) // Overridden by env var
	assert.Equal(t, 10, config.Report)                           // Default value
	assert.Equal(t, 2, config.Collect)                           // Default value
}

// TestNewConfigFromCLI_InvalidEnvValues ensures invalid env vars cause panic
func TestNewConfigFromCLI_InvalidEnvValues(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set invalid environment variables
	os.Setenv("REPORT", "invalid_number")
	os.Setenv("COLLECT", "invalid_number")
	defer os.Clearenv()

	assert.Panics(t, func() {
		NewConfigFromCLI()
	})
}

// TestNewConfigFromCLI_WithoutHTTPPrefix ensures http:// is automatically added
func TestNewConfigFromCLI_WithoutHTTPPrefix(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{"cmd", "-a", "localhost:8000"}
	config := NewConfigFromCLI()

	assert.Equal(t, "http://localhost:8000", config.Address)
}

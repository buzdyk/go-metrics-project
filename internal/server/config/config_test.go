package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewConfig ensures the default config is correctly initialized
func TestNewConfig(t *testing.T) {
	resetConfig()

	config := GetConfig()

	assert.Equal(t, "0.0.0.0:8080", config.Address, "Default address should be 0.0.0.0:8080")
}

// TestNewConfigFromCLI_Defaults ensures CLI args are parsed correctly
func TestNewConfigFromCLI_Defaults(t *testing.T) {
	resetConfig()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{"cmd", "-a", "127.0.0.1:9000"}
	config := GetConfig()

	assert.Equal(t, "127.0.0.1:9000", config.Address, "CLI flag should override default address")
}

// TestNewConfigFromCLI_EnvVariable ensures environment variables override CLI args
func TestNewConfigFromCLI_EnvVariable(t *testing.T) {
	resetConfig()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Setenv("ADDRESS", "192.168.1.100:9090")
	defer os.Clearenv() // Clean up after test

	os.Args = []string{"cmd", "-a", "127.0.0.1:9000"}
	config := GetConfig()

	assert.Equal(t, "192.168.1.100:9090", config.Address, "Environment variable should override CLI flag")
}

// TestNewConfigFromCLI_EnvVariable_Only ensures that if no CLI flag is set, env var is used
func TestNewConfigFromCLI_EnvVariable_Only(t *testing.T) {
	resetConfig()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Setenv("ADDRESS", "192.168.1.200:9091")
	defer os.Clearenv()

	os.Args = []string{"cmd"}
	config := GetConfig()

	assert.Equal(t, "192.168.1.200:9091", config.Address, "Environment variable should be used if no CLI flag is set")
}

// TestNewConfigFromCLI_NoArgs ensures the default config is used when no CLI flags or env vars are provided
func TestNewConfigFromCLI_NoArgs(t *testing.T) {
	resetConfig()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Clearenv()

	os.Args = []string{"cmd"}
	config := GetConfig()

	assert.Equal(t, "0.0.0.0:8080", config.Address, "Default config should be used if no flags or env vars are set")
}

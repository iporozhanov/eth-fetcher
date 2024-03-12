package config_test

import (
	"os"
	"testing"

	"eth-fetcher/config"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables
	os.Setenv("API_PORT", "8080")
	os.Setenv("ETH_NODE_URL", "https://example.com")
	os.Setenv("DB_CONNECTION_URL", "postgres://user:password@localhost:5432/db")

	defer os.Unsetenv("API_PORT")
	defer os.Unsetenv("ETH_NODE_URL")
	defer os.Unsetenv("DB_CONNECTION_URL")
	// Load config
	cfg := config.LoadConfig()

	// Assert values
	assert.Equal(t, "8080", cfg.APIPort)
	assert.Equal(t, "https://example.com", cfg.ETHNodeURL)
	assert.Equal(t, "postgres://user:password@localhost:5432/db", cfg.DBConnectionURL)

	// Clean up environment variables
	os.Unsetenv("API_PORT")
	os.Unsetenv("ETH_NODE_URL")
	os.Unsetenv("DB_CONNECTION_URL")
}

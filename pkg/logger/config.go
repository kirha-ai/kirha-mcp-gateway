package logger

import (
	"github.com/caarlos0/env/v11"
)

// Config holds the logger configuration options.
type Config struct {
	// EnableLogs indicates whether logging is enabled.
	// By default, logging is disabled to avoid issues with stdio usage.
	EnableLogs bool `env:"ENABLE_LOGS" envDefault:"false"`
}

var defaultConfig = &Config{
	EnableLogs: false,
}

// LoadConfig loads the logger configuration from environment variables.
func LoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return defaultConfig
	}

	return cfg
}

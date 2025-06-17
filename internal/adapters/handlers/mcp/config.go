package mcp

import (
	"github.com/caarlos0/env/v11"
)

type HTTPConfig struct {
	Port        int `env:"MCP_PORT" envDefault:"8022"`
	ToolTimeout int `env:"MCP_TOOL_CALL_TIMEOUT_SECONDS" envDefault:"120"`
}

func LoadHTTPConfig() (*HTTPConfig, error) {
	cfg := &HTTPConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

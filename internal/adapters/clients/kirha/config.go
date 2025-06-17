package kirha

import (
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/errors"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds the configuration for the Kirha client.
type Config struct {
	// APIKey is the API key for authenticating with the Kirha API.
	APIKey string `env:"KIRHA_API_KEY,required"`

	// VerticalID is the ID of the vertical to use with the Kirha API.
	VerticalID string `env:"KIRHA_VERTICAL,required"`

	// BaseURL is the base URL of the Kirha API.
	BaseURL string `env:"KIRHA_BASE_URL" envDefault:"https://api.kirha.ai"`

	// Timeout is the timeout for requests to the Kirha API.
	Timeout time.Duration `env:"KIRHA_TIMEOUT" envDefault:"120s"`
}

// LoadConfig loads the Kirha client configuration from environment variables.
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.APIKey == "" {
		return errors.ErrAPIKeyMissing
	}

	if c.VerticalID == "" {
		return errors.ErrVerticalMissing
	}

	if c.Timeout <= 0 {
		return errors.ErrInvalidTimeout
	}

	return nil
}

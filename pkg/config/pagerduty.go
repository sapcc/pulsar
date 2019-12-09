package config

import (
	"fmt"
	"os"
)

const (
	authToken    = "PAGERDUTY_AUTH_TOKEN"
	defaultEmail = "PAGERDUTY_DEFAULT_EMAIL"
)

// PagerdutyConfig ...
type PagerdutyConfig struct {
	AuthToken,
	DefaultEmail string
}

// NewPagerdutyConfigFromEnv returns a new PagerdutyConfig or an error.
func NewPagerdutyConfigFromEnv() (*PagerdutyConfig, error) {
	c := &PagerdutyConfig{
		AuthToken:    os.Getenv(authToken),
		DefaultEmail: os.Getenv(defaultEmail),
	}

	return c, c.validate()
}

func (c *PagerdutyConfig) validate() error {
	if c.AuthToken == "" {
		return fmt.Errorf("missing %s", authToken)
	}

	if c.DefaultEmail == "" {
		return fmt.Errorf("missing %s", defaultEmail)
	}

	return nil
}

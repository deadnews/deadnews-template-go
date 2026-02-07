package main

import (
	"errors"
	"os"
)

// Config holds application configuration from environment variables.
type Config struct {
	DSN string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DSN: os.Getenv("SERVICE_DSN"),
	}

	if cfg.DSN == "" {
		return nil, errors.New("SERVICE_DSN environment variable is required")
	}

	return cfg, nil
}

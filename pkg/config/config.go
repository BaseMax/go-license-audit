package config

import (
	"encoding/json"
	"os"
)

// Config represents the configuration for license auditing
type Config struct {
	// IncompatibleLicenses are licenses that are not allowed
	IncompatibleLicenses []string `json:"incompatible_licenses"`
	// RiskyLicenses are licenses that require review
	RiskyLicenses []string `json:"risky_licenses"`
	// AllowedLicenses are explicitly allowed licenses (optional)
	AllowedLicenses []string `json:"allowed_licenses"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		IncompatibleLicenses: []string{
			"GPL-2.0",
			"GPL-3.0",
			"AGPL-3.0",
			"AGPL-1.0",
			"LGPL-2.0",
			"LGPL-2.1",
			"LGPL-3.0",
			"SSPL",
			"Commons Clause",
		},
		RiskyLicenses: []string{
			"MPL-2.0",
			"EPL-1.0",
			"EPL-2.0",
			"CDDL-1.0",
			"CPL-1.0",
			"OSL-3.0",
		},
		AllowedLicenses: []string{
			"MIT",
			"Apache-2.0",
			"BSD-2-Clause",
			"BSD-3-Clause",
			"ISC",
			"0BSD",
			"Unlicense",
			"CC0-1.0",
		},
	}
}

// Load loads configuration from a file, or returns default config if path is empty
func Load(path string) (*Config, error) {
	if path == "" {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// IsIncompatible checks if a license is in the incompatible list
func (c *Config) IsIncompatible(license string) bool {
	for _, l := range c.IncompatibleLicenses {
		if l == license {
			return true
		}
	}
	return false
}

// IsRisky checks if a license is in the risky list
func (c *Config) IsRisky(license string) bool {
	for _, l := range c.RiskyLicenses {
		if l == license {
			return true
		}
	}
	return false
}

// IsAllowed checks if a license is in the allowed list
func (c *Config) IsAllowed(license string) bool {
	// If no allowed licenses specified, allow by default
	if len(c.AllowedLicenses) == 0 {
		return true
	}

	for _, l := range c.AllowedLicenses {
		if l == license {
			return true
		}
	}
	return false
}

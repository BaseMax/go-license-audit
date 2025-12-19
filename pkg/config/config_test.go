package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if len(cfg.IncompatibleLicenses) == 0 {
		t.Error("Default config should have incompatible licenses")
	}

	if len(cfg.RiskyLicenses) == 0 {
		t.Error("Default config should have risky licenses")
	}

	if len(cfg.AllowedLicenses) == 0 {
		t.Error("Default config should have allowed licenses")
	}
}

func TestLoadConfig(t *testing.T) {
	// Test loading with empty path (should return default)
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	if len(cfg.IncompatibleLicenses) == 0 {
		t.Error("Default config should have incompatible licenses")
	}

	// Test loading from file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	testConfig := &Config{
		IncompatibleLicenses: []string{"GPL-3.0", "AGPL-3.0"},
		RiskyLicenses:        []string{"MPL-2.0"},
		AllowedLicenses:      []string{"MIT", "Apache-2.0"},
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err = Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config from file: %v", err)
	}

	if len(cfg.IncompatibleLicenses) != 2 {
		t.Errorf("Expected 2 incompatible licenses, got %d", len(cfg.IncompatibleLicenses))
	}

	if len(cfg.RiskyLicenses) != 1 {
		t.Errorf("Expected 1 risky license, got %d", len(cfg.RiskyLicenses))
	}
}

func TestIsIncompatible(t *testing.T) {
	cfg := &Config{
		IncompatibleLicenses: []string{"GPL-3.0", "AGPL-3.0"},
	}

	if !cfg.IsIncompatible("GPL-3.0") {
		t.Error("GPL-3.0 should be incompatible")
	}

	if cfg.IsIncompatible("MIT") {
		t.Error("MIT should not be incompatible")
	}
}

func TestIsRisky(t *testing.T) {
	cfg := &Config{
		RiskyLicenses: []string{"MPL-2.0", "EPL-1.0"},
	}

	if !cfg.IsRisky("MPL-2.0") {
		t.Error("MPL-2.0 should be risky")
	}

	if cfg.IsRisky("MIT") {
		t.Error("MIT should not be risky")
	}
}

func TestIsAllowed(t *testing.T) {
	// Test with explicit allowed licenses
	cfg := &Config{
		AllowedLicenses: []string{"MIT", "Apache-2.0"},
	}

	if !cfg.IsAllowed("MIT") {
		t.Error("MIT should be allowed")
	}

	if cfg.IsAllowed("GPL-3.0") {
		t.Error("GPL-3.0 should not be allowed")
	}

	// Test with empty allowed licenses (should allow all)
	cfg2 := &Config{
		AllowedLicenses: []string{},
	}

	if !cfg2.IsAllowed("GPL-3.0") {
		t.Error("Any license should be allowed when allowed list is empty")
	}
}

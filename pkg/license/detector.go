package license

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Detector handles license detection for dependencies
type Detector struct {
	client *http.Client
}

// NewDetector creates a new license detector
func NewDetector() *Detector {
	return &Detector{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// DetectLicense attempts to detect the license for a dependency
func (d *Detector) DetectLicense(name, version, ecosystem string) (string, string) {
	var license string
	var licenseType string

	switch ecosystem {
	case "npm":
		license, licenseType = d.detectNPMLicense(name, version)
	case "pypi":
		license, licenseType = d.detectPyPILicense(name)
	case "go":
		license, licenseType = d.detectGoLicense(name)
	case "packagist":
		license, licenseType = d.detectPackagistLicense(name)
	default:
		license = "Unknown"
		licenseType = "unknown"
	}

	return license, licenseType
}

// detectNPMLicense detects license from npm registry
func (d *Detector) detectNPMLicense(name, version string) (string, string) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", name)
	resp, err := d.client.Get(url)
	if err != nil {
		return "Unknown", "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Unknown", "unknown"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown", "unknown"
	}

	var pkg struct {
		License  interface{} `json:"license"`
		Versions map[string]struct {
			License interface{} `json:"license"`
		} `json:"versions"`
	}

	if err := json.Unmarshal(body, &pkg); err != nil {
		return "Unknown", "unknown"
	}

	// Try to get license from specific version first
	if version != "" && pkg.Versions != nil {
		cleanVersion := strings.TrimPrefix(version, "^")
		cleanVersion = strings.TrimPrefix(cleanVersion, "~")
		if versionData, ok := pkg.Versions[cleanVersion]; ok {
			if license := extractLicense(versionData.License); license != "" {
				return license, "detected"
			}
		}
	}

	// Fall back to latest license
	if license := extractLicense(pkg.License); license != "" {
		return license, "detected"
	}

	return "Unknown", "unknown"
}

// detectPyPILicense detects license from PyPI
func (d *Detector) detectPyPILicense(name string) (string, string) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", name)
	resp, err := d.client.Get(url)
	if err != nil {
		return "Unknown", "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Unknown", "unknown"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown", "unknown"
	}

	var pkg struct {
		Info struct {
			License string `json:"license"`
		} `json:"info"`
	}

	if err := json.Unmarshal(body, &pkg); err != nil {
		return "Unknown", "unknown"
	}

	if pkg.Info.License != "" && pkg.Info.License != "UNKNOWN" {
		return pkg.Info.License, "detected"
	}

	return "Unknown", "unknown"
}

// detectGoLicense detects license from Go packages
func (d *Detector) detectGoLicense(name string) (string, string) {
	// For Go packages, we use pkg.go.dev API
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/go/packages/%s", name)
	resp, err := d.client.Get(url)
	if err != nil {
		return "Unknown", "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Unknown", "unknown"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown", "unknown"
	}

	var pkg struct {
		Versions []struct {
			Licenses []string `json:"licenses"`
		} `json:"versions"`
	}

	if err := json.Unmarshal(body, &pkg); err != nil {
		return "Unknown", "unknown"
	}

	if len(pkg.Versions) > 0 && len(pkg.Versions[0].Licenses) > 0 {
		return pkg.Versions[0].Licenses[0], "detected"
	}

	return "Unknown", "unknown"
}

// detectPackagistLicense detects license from Packagist (PHP)
func (d *Detector) detectPackagistLicense(name string) (string, string) {
	url := fmt.Sprintf("https://repo.packagist.org/p2/%s.json", name)
	resp, err := d.client.Get(url)
	if err != nil {
		return "Unknown", "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Unknown", "unknown"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown", "unknown"
	}

	var pkg struct {
		Packages map[string][]struct {
			License []string `json:"license"`
		} `json:"packages"`
	}

	if err := json.Unmarshal(body, &pkg); err != nil {
		return "Unknown", "unknown"
	}

	if packages, ok := pkg.Packages[name]; ok && len(packages) > 0 {
		if len(packages[0].License) > 0 {
			return strings.Join(packages[0].License, " OR "), "detected"
		}
	}

	return "Unknown", "unknown"
}

// extractLicense extracts license string from various formats
func extractLicense(license interface{}) string {
	switch v := license.(type) {
	case string:
		return v
	case map[string]interface{}:
		if typeVal, ok := v["type"].(string); ok {
			return typeVal
		}
	}
	return ""
}

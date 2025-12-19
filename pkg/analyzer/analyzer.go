package analyzer

import (
	"fmt"
	"time"

	"github.com/BaseMax/go-license-audit/pkg/config"
	"github.com/BaseMax/go-license-audit/pkg/license"
	"github.com/BaseMax/go-license-audit/pkg/parser"
	"github.com/BaseMax/go-license-audit/pkg/types"
)

// Analyze analyzes a project for license compliance
func Analyze(projectPath string, cfg *config.Config) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		ProjectPath:      projectPath,
		Dependencies:     []types.Dependency{},
		IncompatibleDeps: []types.Dependency{},
		RiskyDeps:        []types.Dependency{},
		Timestamp:        time.Now().Format(time.RFC3339),
	}

	// Get all parsers
	parsers := parser.GetParsers()

	// Detect and parse dependencies from all available formats
	allDeps := []types.Dependency{}
	foundParsers := []string{}

	for _, p := range parsers {
		if p.Detect(projectPath) {
			deps, err := p.Parse(projectPath)
			if err != nil {
				return nil, fmt.Errorf("error parsing dependencies: %w", err)
			}
			allDeps = append(allDeps, deps...)

			// Get parser type for logging
			switch p.(type) {
			case *parser.GoModParser:
				foundParsers = append(foundParsers, "go.mod")
			case *parser.PackageJSONParser:
				foundParsers = append(foundParsers, "package.json")
			case *parser.RequirementsTxtParser:
				foundParsers = append(foundParsers, "requirements.txt")
			case *parser.ComposerJSONParser:
				foundParsers = append(foundParsers, "composer.json")
			}
		}
	}

	if len(allDeps) == 0 {
		return nil, fmt.Errorf("no dependency files found in %s", projectPath)
	}

	fmt.Printf("Found dependency files: %v\n", foundParsers)
	fmt.Printf("Analyzing %d dependencies...\n", len(allDeps))

	// Detect licenses for all dependencies
	detector := license.NewDetector()

	for i, dep := range allDeps {
		fmt.Printf("  [%d/%d] %s (%s)... ", i+1, len(allDeps), dep.Name, dep.Ecosystem)

		licenseName, licenseType := detector.DetectLicense(dep.Name, dep.Version, dep.Ecosystem)
		dep.License = licenseName
		dep.LicenseType = licenseType

		// Determine risk level
		if cfg.IsIncompatible(licenseName) {
			dep.Risk = "incompatible"
			result.IncompatibleDeps = append(result.IncompatibleDeps, dep)
			fmt.Printf("INCOMPATIBLE (%s)\n", licenseName)
		} else if cfg.IsRisky(licenseName) {
			dep.Risk = "risky"
			result.RiskyDeps = append(result.RiskyDeps, dep)
			fmt.Printf("RISKY (%s)\n", licenseName)
		} else if licenseName == "Unknown" {
			dep.Risk = "unknown"
			fmt.Printf("UNKNOWN\n")
		} else {
			dep.Risk = "safe"
			fmt.Printf("OK (%s)\n", licenseName)
		}

		result.Dependencies = append(result.Dependencies, dep)
	}

	return result, nil
}

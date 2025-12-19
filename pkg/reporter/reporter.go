package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BaseMax/go-license-audit/pkg/types"
)

// GenerateJSON generates a JSON report
func GenerateJSON(result *types.AnalysisResult, outputDir string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	outputPath := filepath.Join(outputDir, "license-report.json")
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	return nil
}

// GenerateMarkdown generates a Markdown report
func GenerateMarkdown(result *types.AnalysisResult, outputDir string) error {
	var sb strings.Builder

	sb.WriteString("# License Compliance Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", result.Timestamp))
	sb.WriteString(fmt.Sprintf("**Project Path:** %s\n\n", result.ProjectPath))

	// Summary section
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Dependencies:** %d\n", len(result.Dependencies)))
	sb.WriteString(fmt.Sprintf("- **Incompatible Licenses:** %d\n", len(result.IncompatibleDeps)))
	sb.WriteString(fmt.Sprintf("- **Risky Licenses:** %d\n", len(result.RiskyDeps)))
	sb.WriteString(fmt.Sprintf("- **Safe/Unknown:** %d\n\n", len(result.Dependencies)-len(result.IncompatibleDeps)-len(result.RiskyDeps)))

	// Incompatible licenses section
	if len(result.IncompatibleDeps) > 0 {
		sb.WriteString("## ⛔ Incompatible Licenses\n\n")
		sb.WriteString("These dependencies have licenses that are incompatible with your project:\n\n")
		sb.WriteString("| Package | Version | License | Ecosystem |\n")
		sb.WriteString("|---------|---------|---------|----------|\n")
		for _, dep := range result.IncompatibleDeps {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				dep.Name, dep.Version, dep.License, dep.Ecosystem))
		}
		sb.WriteString("\n")
	}

	// Risky licenses section
	if len(result.RiskyDeps) > 0 {
		sb.WriteString("## ⚠️ Risky Licenses\n\n")
		sb.WriteString("These dependencies have licenses that require review:\n\n")
		sb.WriteString("| Package | Version | License | Ecosystem |\n")
		sb.WriteString("|---------|---------|---------|----------|\n")
		for _, dep := range result.RiskyDeps {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				dep.Name, dep.Version, dep.License, dep.Ecosystem))
		}
		sb.WriteString("\n")
	}

	// All dependencies section
	sb.WriteString("## 📦 All Dependencies\n\n")
	sb.WriteString("| Package | Version | License | Risk | Ecosystem |\n")
	sb.WriteString("|---------|---------|---------|------|----------|\n")
	for _, dep := range result.Dependencies {
		riskIcon := getRiskIcon(dep.Risk)
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s %s | %s |\n",
			dep.Name, dep.Version, dep.License, riskIcon, dep.Risk, dep.Ecosystem))
	}
	sb.WriteString("\n")

	// License distribution
	sb.WriteString("## 📊 License Distribution\n\n")
	licenseCount := make(map[string]int)
	for _, dep := range result.Dependencies {
		licenseCount[dep.License]++
	}

	sb.WriteString("| License | Count |\n")
	sb.WriteString("|---------|-------|\n")
	for license, count := range licenseCount {
		sb.WriteString(fmt.Sprintf("| %s | %d |\n", license, count))
	}
	sb.WriteString("\n")

	outputPath := filepath.Join(outputDir, "license-report.md")
	if err := os.WriteFile(outputPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("error writing Markdown file: %w", err)
	}

	return nil
}

func getRiskIcon(risk string) string {
	switch risk {
	case "incompatible":
		return "⛔"
	case "risky":
		return "⚠️"
	case "safe":
		return "✅"
	case "unknown":
		return "❓"
	default:
		return "➖"
	}
}

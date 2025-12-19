package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BaseMax/go-license-audit/pkg/analyzer"
	"github.com/BaseMax/go-license-audit/pkg/config"
	"github.com/BaseMax/go-license-audit/pkg/reporter"
)

func main() {
	var (
		projectPath  = flag.String("path", ".", "Path to the project directory")
		configFile   = flag.String("config", "", "Path to configuration file")
		outputFormat = flag.String("format", "both", "Output format: json, markdown, or both")
		outputDir    = flag.String("output", ".", "Output directory for reports")
		strictMode   = flag.Bool("strict", false, "Exit with error code if incompatible licenses found")
	)
	
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Analyze project
	fmt.Printf("Analyzing project at: %s\n", *projectPath)
	result, err := analyzer.Analyze(*projectPath, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing project: %v\n", err)
		os.Exit(1)
	}

	// Generate reports
	if *outputFormat == "json" || *outputFormat == "both" {
		if err := reporter.GenerateJSON(result, *outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JSON report: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("JSON report generated: %s/license-report.json\n", *outputDir)
	}

	if *outputFormat == "markdown" || *outputFormat == "both" {
		if err := reporter.GenerateMarkdown(result, *outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating Markdown report: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Markdown report generated: %s/license-report.md\n", *outputDir)
	}

	// Print summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total dependencies: %d\n", len(result.Dependencies))
	fmt.Printf("  Incompatible licenses: %d\n", len(result.IncompatibleDeps))
	fmt.Printf("  Risky licenses: %d\n", len(result.RiskyDeps))

	// Exit with error code in strict mode if issues found
	if *strictMode && (len(result.IncompatibleDeps) > 0 || len(result.RiskyDeps) > 0) {
		os.Exit(1)
	}
}

package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoModParser(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a test go.mod file
	goModContent := `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.0
)
`
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create test go.mod: %v", err)
	}

	parser := &GoModParser{}

	// Test detection
	if !parser.Detect(tmpDir) {
		t.Error("GoModParser should detect go.mod")
	}

	// Test parsing
	deps, err := parser.Parse(tmpDir)
	if err != nil {
		t.Fatalf("Failed to parse go.mod: %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(deps))
	}

	// Verify first dependency
	if len(deps) > 0 {
		if deps[0].Name != "github.com/pkg/errors" {
			t.Errorf("Expected name 'github.com/pkg/errors', got '%s'", deps[0].Name)
		}
		if deps[0].Version != "v0.9.1" {
			t.Errorf("Expected version 'v0.9.1', got '%s'", deps[0].Version)
		}
		if deps[0].Ecosystem != "go" {
			t.Errorf("Expected ecosystem 'go', got '%s'", deps[0].Ecosystem)
		}
	}
}

func TestPackageJSONParser(t *testing.T) {
	tmpDir := t.TempDir()

	packageJSONContent := `{
  "name": "test-app",
  "dependencies": {
    "express": "^4.18.2",
    "lodash": "^4.17.21"
  },
  "devDependencies": {
    "jest": "^29.0.0"
  }
}`
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644); err != nil {
		t.Fatalf("Failed to create test package.json: %v", err)
	}

	parser := &PackageJSONParser{}

	if !parser.Detect(tmpDir) {
		t.Error("PackageJSONParser should detect package.json")
	}

	deps, err := parser.Parse(tmpDir)
	if err != nil {
		t.Fatalf("Failed to parse package.json: %v", err)
	}

	if len(deps) != 3 {
		t.Errorf("Expected 3 dependencies, got %d", len(deps))
	}

	// Check ecosystem
	for _, dep := range deps {
		if dep.Ecosystem != "npm" {
			t.Errorf("Expected ecosystem 'npm', got '%s'", dep.Ecosystem)
		}
	}
}

func TestRequirementsTxtParser(t *testing.T) {
	tmpDir := t.TempDir()

	requirementsContent := `# Test requirements
requests==2.28.0
flask>=2.0.0
numpy~=1.24.0

# Comment line
django>3.0
`
	requirementsPath := filepath.Join(tmpDir, "requirements.txt")
	if err := os.WriteFile(requirementsPath, []byte(requirementsContent), 0644); err != nil {
		t.Fatalf("Failed to create test requirements.txt: %v", err)
	}

	parser := &RequirementsTxtParser{}

	if !parser.Detect(tmpDir) {
		t.Error("RequirementsTxtParser should detect requirements.txt")
	}

	deps, err := parser.Parse(tmpDir)
	if err != nil {
		t.Fatalf("Failed to parse requirements.txt: %v", err)
	}

	if len(deps) != 4 {
		t.Errorf("Expected 4 dependencies, got %d", len(deps))
	}

	// Verify parsing of different version specifiers
	expectedDeps := map[string]string{
		"requests": "2.28.0",
		"flask":    "2.0.0",
		"numpy":    "1.24.0",
		"django":   "3.0",
	}

	for _, dep := range deps {
		if expectedVersion, ok := expectedDeps[dep.Name]; ok {
			if dep.Version != expectedVersion {
				t.Errorf("For %s, expected version '%s', got '%s'", dep.Name, expectedVersion, dep.Version)
			}
		}
		if dep.Ecosystem != "pypi" {
			t.Errorf("Expected ecosystem 'pypi', got '%s'", dep.Ecosystem)
		}
	}
}

func TestComposerJSONParser(t *testing.T) {
	tmpDir := t.TempDir()

	composerContent := `{
  "name": "test/app",
  "require": {
    "php": "^8.0",
    "symfony/http-foundation": "^6.0",
    "monolog/monolog": "^3.0",
    "ext-json": "*"
  },
  "require-dev": {
    "phpunit/phpunit": "^10.0"
  }
}`
	composerPath := filepath.Join(tmpDir, "composer.json")
	if err := os.WriteFile(composerPath, []byte(composerContent), 0644); err != nil {
		t.Fatalf("Failed to create test composer.json: %v", err)
	}

	parser := &ComposerJSONParser{}

	if !parser.Detect(tmpDir) {
		t.Error("ComposerJSONParser should detect composer.json")
	}

	deps, err := parser.Parse(tmpDir)
	if err != nil {
		t.Fatalf("Failed to parse composer.json: %v", err)
	}

	// Should skip php and ext-json
	if len(deps) != 3 {
		t.Errorf("Expected 3 dependencies (excluding php and extensions), got %d", len(deps))
	}

	// Verify none of the dependencies are PHP or extensions
	for _, dep := range deps {
		if dep.Name == "php" || dep.Name[:4] == "ext-" {
			t.Errorf("Should not include PHP or extensions, got: %s", dep.Name)
		}
		if dep.Ecosystem != "packagist" {
			t.Errorf("Expected ecosystem 'packagist', got '%s'", dep.Ecosystem)
		}
	}
}

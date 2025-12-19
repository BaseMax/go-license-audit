package parser

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/BaseMax/go-license-audit/pkg/types"
)

// Parser interface for different package managers
type Parser interface {
	Parse(projectPath string) ([]types.Dependency, error)
	Detect(projectPath string) bool
}

// GoModParser parses go.mod files
type GoModParser struct{}

func (p *GoModParser) Detect(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, "go.mod"))
	return err == nil
}

func (p *GoModParser) Parse(projectPath string) ([]types.Dependency, error) {
	filePath := filepath.Join(projectPath, "go.mod")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []types.Dependency
	scanner := bufio.NewScanner(file)
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}

		if inRequire && line == ")" {
			inRequire = false
			continue
		}

		if strings.HasPrefix(line, "require ") || inRequire {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				if name == "require" {
					if len(parts) >= 3 {
						name = parts[1]
						version := parts[2]
						deps = append(deps, types.Dependency{
							Name:      name,
							Version:   version,
							Ecosystem: "go",
						})
					}
				} else {
					version := parts[1]
					deps = append(deps, types.Dependency{
						Name:      name,
						Version:   version,
						Ecosystem: "go",
					})
				}
			}
		}
	}

	return deps, scanner.Err()
}

// PackageJSONParser parses package.json files
type PackageJSONParser struct{}

func (p *PackageJSONParser) Detect(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, "package.json"))
	return err == nil
}

func (p *PackageJSONParser) Parse(projectPath string) ([]types.Dependency, error) {
	filePath := filepath.Join(projectPath, "package.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var deps []types.Dependency
	for name, version := range pkg.Dependencies {
		deps = append(deps, types.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "npm",
		})
	}
	for name, version := range pkg.DevDependencies {
		deps = append(deps, types.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "npm",
		})
	}

	return deps, nil
}

// RequirementsTxtParser parses requirements.txt files
type RequirementsTxtParser struct{}

func (p *RequirementsTxtParser) Detect(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, "requirements.txt"))
	return err == nil
}

func (p *RequirementsTxtParser) Parse(projectPath string) ([]types.Dependency, error) {
	filePath := filepath.Join(projectPath, "requirements.txt")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []types.Dependency
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse package specification
		// Format: package==version or package>=version or package
		name := line
		version := ""

		for _, sep := range []string{"==", ">=", "<=", "~=", ">", "<"} {
			if idx := strings.Index(line, sep); idx != -1 {
				name = line[:idx]
				version = strings.TrimSpace(line[idx+len(sep):])
				break
			}
		}

		name = strings.TrimSpace(name)
		if name != "" {
			deps = append(deps, types.Dependency{
				Name:      name,
				Version:   version,
				Ecosystem: "pypi",
			})
		}
	}

	return deps, scanner.Err()
}

// ComposerJSONParser parses composer.json files for PHP projects
type ComposerJSONParser struct{}

func (p *ComposerJSONParser) Detect(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, "composer.json"))
	return err == nil
}

func (p *ComposerJSONParser) Parse(projectPath string) ([]types.Dependency, error) {
	filePath := filepath.Join(projectPath, "composer.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var composer struct {
		Require    map[string]string `json:"require"`
		RequireDev map[string]string `json:"require-dev"`
	}

	if err := json.Unmarshal(data, &composer); err != nil {
		return nil, err
	}

	var deps []types.Dependency
	for name, version := range composer.Require {
		// Skip PHP itself and extensions
		if name == "php" || strings.HasPrefix(name, "ext-") {
			continue
		}
		deps = append(deps, types.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "packagist",
		})
	}
	for name, version := range composer.RequireDev {
		// Skip PHP itself and extensions
		if name == "php" || strings.HasPrefix(name, "ext-") {
			continue
		}
		deps = append(deps, types.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "packagist",
		})
	}

	return deps, nil
}

// GetParsers returns all available parsers
func GetParsers() []Parser {
	return []Parser{
		&GoModParser{},
		&PackageJSONParser{},
		&RequirementsTxtParser{},
		&ComposerJSONParser{},
	}
}

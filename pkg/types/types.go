package types

// Dependency represents a single dependency in a project
type Dependency struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	License     string `json:"license"`
	LicenseType string `json:"license_type"` // detected, unknown, etc.
	Ecosystem   string `json:"ecosystem"`    // go, npm, pypi, packagist
	Risk        string `json:"risk"`         // safe, risky, incompatible
}

// AnalysisResult represents the complete analysis result
type AnalysisResult struct {
	ProjectPath      string       `json:"project_path"`
	Dependencies     []Dependency `json:"dependencies"`
	IncompatibleDeps []Dependency `json:"incompatible_dependencies"`
	RiskyDeps        []Dependency `json:"risky_dependencies"`
	Timestamp        string       `json:"timestamp"`
}

// ProjectType represents the type of project being analyzed
type ProjectType string

const (
	ProjectTypeGo     ProjectType = "go"
	ProjectTypeNodeJS ProjectType = "nodejs"
	ProjectTypePython ProjectType = "python"
	ProjectTypePHP    ProjectType = "php"
)

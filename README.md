# go-license-audit

A fast and extensible license compliance scanner for project dependencies across multiple ecosystems.

## Features

- 🔍 **Multi-Ecosystem Support**: Analyze dependencies from:
  - Go (go.mod)
  - Node.js (package.json)
  - Python (requirements.txt)
  - PHP (composer.json)
- 📊 **Comprehensive Reporting**: Generate reports in JSON and Markdown formats
- ⚙️ **Configurable Rules**: Define incompatible and risky licenses via config file
- 🚀 **CI/CD Friendly**: Exit codes and strict mode for pipeline integration
- 🔌 **Extensible Architecture**: Easy to add support for new package managers

## Installation

### From Source

```bash
git clone https://github.com/BaseMax/go-license-audit.git
cd go-license-audit
go build -o go-license-audit .
```

### Using Go Install

```bash
go install github.com/BaseMax/go-license-audit@latest
```

## Usage

### Basic Usage

Analyze the current directory:

```bash
./go-license-audit
```

Analyze a specific project:

```bash
./go-license-audit -path /path/to/project
```

### Output Formats

Generate both JSON and Markdown reports (default):

```bash
./go-license-audit -format both -output ./reports
```

Generate only JSON:

```bash
./go-license-audit -format json -output ./reports
```

Generate only Markdown:

```bash
./go-license-audit -format markdown -output ./reports
```

### CI/CD Integration

Use strict mode to fail the build if incompatible or risky licenses are found:

```bash
./go-license-audit -strict -path . -format json -output ./reports
```

Exit codes:
- `0`: Success (no incompatible licenses or not in strict mode)
- `1`: Failure (incompatible/risky licenses found in strict mode, or error occurred)

### Configuration

Create a `config.json` file to customize license rules:

```json
{
  "incompatible_licenses": [
    "GPL-2.0",
    "GPL-3.0",
    "AGPL-3.0"
  ],
  "risky_licenses": [
    "MPL-2.0",
    "EPL-1.0",
    "EPL-2.0"
  ],
  "allowed_licenses": [
    "MIT",
    "Apache-2.0",
    "BSD-2-Clause",
    "BSD-3-Clause",
    "ISC"
  ]
}
```

Use the config file:

```bash
./go-license-audit -config config.json
```

## Command Line Options

| Option | Default | Description |
|--------|---------|-------------|
| `-path` | `.` | Path to the project directory |
| `-config` | (none) | Path to configuration file (uses defaults if not specified) |
| `-format` | `both` | Output format: `json`, `markdown`, or `both` |
| `-output` | `.` | Output directory for reports |
| `-strict` | `false` | Exit with error code if incompatible licenses found |

## Report Examples

### JSON Report

```json
{
  "project_path": "/path/to/project",
  "dependencies": [
    {
      "name": "github.com/example/package",
      "version": "v1.2.3",
      "license": "MIT",
      "license_type": "detected",
      "ecosystem": "go",
      "risk": "safe"
    }
  ],
  "incompatible_dependencies": [],
  "risky_dependencies": [],
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Markdown Report

The Markdown report includes:
- Summary statistics
- Incompatible licenses table (⛔)
- Risky licenses table (⚠️)
- All dependencies with risk indicators
- License distribution chart

## Supported Ecosystems

### Go (go.mod)

Parses `go.mod` files and detects licenses via deps.dev API.

### Node.js (package.json)

Parses `package.json` including both dependencies and devDependencies. Detects licenses from npm registry.

### Python (requirements.txt)

Parses `requirements.txt` with support for various version specifiers (`==`, `>=`, `~=`, etc.). Detects licenses from PyPI.

### PHP (composer.json)

Parses `composer.json` including both require and require-dev sections. Detects licenses from Packagist.

## License Detection

The tool automatically detects licenses by querying public registries:
- **npm**: registry.npmjs.org
- **PyPI**: pypi.org
- **Go**: deps.dev API
- **Packagist**: repo.packagist.org

## Default License Rules

### Incompatible Licenses

- GPL-2.0, GPL-3.0
- AGPL-3.0, AGPL-1.0
- LGPL-2.0, LGPL-2.1, LGPL-3.0
- SSPL
- Commons Clause

### Risky Licenses (Require Review)

- MPL-2.0
- EPL-1.0, EPL-2.0
- CDDL-1.0
- CPL-1.0
- OSL-3.0

### Allowed Licenses

- MIT
- Apache-2.0
- BSD-2-Clause, BSD-3-Clause
- ISC
- 0BSD
- Unlicense
- CC0-1.0

## Contributing

Contributions are welcome! To add support for a new package manager:

1. Implement the `Parser` interface in `pkg/parser/parser.go`
2. Add license detection logic in `pkg/license/detector.go`
3. Register the parser in `GetParsers()` function
4. Update documentation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Max Base**

- GitHub: [@BaseMax](https://github.com/BaseMax)

## Acknowledgments

- Uses public package registries for license detection
- Built with extensibility and CI/CD integration in mind

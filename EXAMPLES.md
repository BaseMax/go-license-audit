# Usage Examples

This document demonstrates the usage of go-license-audit across different ecosystems.

## Example 1: Analyzing a Node.js Project

```bash
$ ./go-license-audit -path examples/node-project
Analyzing project at: examples/node-project
Found dependency files: [package.json]
Analyzing 3 dependencies...
  [1/3] express (npm)... OK (MIT)
  [2/3] lodash (npm)... OK (MIT)
  [3/3] jest (npm)... OK (MIT)
JSON report generated: ./license-report.json
Markdown report generated: ./license-report.md

Summary:
  Total dependencies: 3
  Incompatible licenses: 0
  Risky licenses: 0
```

## Example 2: Analyzing a PHP Project

```bash
$ ./go-license-audit -path examples/php-project -format markdown
Analyzing project at: examples/php-project
Found dependency files: [composer.json]
Analyzing 3 dependencies...
  [1/3] symfony/http-foundation (packagist)... OK (MIT)
  [2/3] monolog/monolog (packagist)... OK (MIT)
  [3/3] phpunit/phpunit (packagist)... OK (BSD-3-Clause)
Markdown report generated: ./license-report.md

Summary:
  Total dependencies: 3
  Incompatible licenses: 0
  Risky licenses: 0
```

## Example 3: Multi-Ecosystem Project

```bash
$ ./go-license-audit -path examples/mixed-project -format both
Analyzing project at: examples/mixed-project
Found dependency files: [go.mod package.json requirements.txt]
Analyzing 3 dependencies...
  [1/3] github.com/sirupsen/logrus (go)... UNKNOWN
  [2/3] express (npm)... OK (MIT)
  [3/3] requests (pypi)... OK (Apache-2.0)
JSON report generated: ./license-report.json
Markdown report generated: ./license-report.md

Summary:
  Total dependencies: 3
  Incompatible licenses: 0
  Risky licenses: 0
```

## Example 4: Using Strict Mode for CI/CD

```bash
$ ./go-license-audit -path . -strict -format json -output ./reports
# Exit code 0 if no issues
# Exit code 1 if incompatible or risky licenses found
```

## Example 5: Custom Configuration

Create a `config.json`:

```json
{
  "incompatible_licenses": ["GPL-3.0", "AGPL-3.0"],
  "risky_licenses": ["MPL-2.0"],
  "allowed_licenses": ["MIT", "Apache-2.0", "BSD-3-Clause"]
}
```

Run with config:

```bash
$ ./go-license-audit -config config.json -path .
```

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: License Audit
  run: |
    go install github.com/BaseMax/go-license-audit@latest
    go-license-audit -strict -format json -output ./reports
```

### GitLab CI Example

```yaml
license-audit:
  script:
    - go install github.com/BaseMax/go-license-audit@latest
    - go-license-audit -strict -format json -output ./reports
  artifacts:
    paths:
      - reports/
```

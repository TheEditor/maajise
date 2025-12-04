package templates

import "fmt"

func init() {
	Register(&GoTemplate{})
}

// GoTemplate is a Go project template
type GoTemplate struct{}

func (t *GoTemplate) Name() string {
	return "go"
}

func (t *GoTemplate) Description() string {
	return "Go project with Standard Go Project Layout (cmd/, internal/, pkg/)"
}

func (t *GoTemplate) Dependencies() []string {
	return []string{"git", "bd", "go"}
}

func (t *GoTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore":                      t.gitignore(),
		".ubsignore":                      t.ubsignore(),
		"README.md":                       t.readme(projectName),
		"go.mod":                          t.goMod(projectName),
		"api/.gitkeep":                    "",
		"build/.gitkeep":                  "",
		"cmd/" + projectName + "/main.go": t.mainGo(),
		"configs/.gitkeep":                "",
		"docs/.gitkeep":                   "",
		"internal/.gitkeep":               "",
		"pkg/.gitkeep":                    "",
		"scripts/.gitkeep":                "",
		"test/.gitkeep":                   "",
	}
}

func (t *GoTemplate) gitignore() string {
	return `# Go
/bin/
/dist/
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# Dependency directories
/vendor/

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env

# Coverage
coverage.out
coverage.html
`
}

func (t *GoTemplate) ubsignore() string {
	return `# UBS Scanner Ignore File
bin/
dist/
vendor/
.git/
.vscode/
.idea/
.beads/
.claude/
*.md
*.sum
*.mod
`
}

func (t *GoTemplate) readme(projectName string) string {
	return fmt.Sprintf(`# %s

## Project Structure

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

| Directory | Purpose |
|-----------|---------|
| cmd/%s/ | Main application entry point |
| internal/ | Private application code (enforced by Go) |
| pkg/ | Public library code (can be imported externally) |
| api/ | API definitions (OpenAPI, Protocol Buffers, JSON Schema) |
| configs/ | Configuration file templates |
| scripts/ | Build, install, and management scripts |
| test/ | Additional test apps and test data |
| docs/ | Documentation beyond README |
| build/ | Packaging and CI configuration |

## Development

### Build
` + "```bash" + `
go build -o bin/%s ./cmd/%s
` + "```" + `

### Run directly
` + "```bash" + `
go run ./cmd/%s
` + "```" + `

### Test all packages
` + "```bash" + `
go test ./...
` + "```" + `

### Format code
` + "```bash" + `
go fmt ./...
` + "```" + `

## Getting Started

1. Add your business logic to internal/
2. Create public APIs in pkg/ if needed
3. Import internal packages in cmd/%s/main.go

## Issue Tracking

` + "```bash" + `
bd list           # View issues
bd create --title "Task"  # Create issue
` + "```" + `

## Code Quality

` + "```bash" + `
ubs .             # Scan for bugs
go vet ./...      # Go static analysis
` + "```" + `
`, projectName, projectName, projectName, projectName, projectName, projectName)
}

func (t *GoTemplate) goMod(projectName string) string {
	return fmt.Sprintf(`module %s

go 1.23
`, projectName)
}

func (t *GoTemplate) mainGo() string {
	return `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, Go!")
	// TODO: Import and use packages from internal/ as your app grows
	// Example:
	// import "your-module/internal/config"
	// import "your-module/internal/services"
}
`
}

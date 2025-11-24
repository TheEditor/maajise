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
	return "Go project with module configuration"
}

func (t *GoTemplate) Dependencies() []string {
	return []string{"git", "bd", "go"}
}

func (t *GoTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore": t.gitignore(),
		".ubsignore": t.ubsignore(),
		"README.md":  t.readme(projectName),
		"go.mod":     t.goMod(projectName),
		"main.go":    t.mainGo(),
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

A Go project.

## Setup

` + "```bash" + `
# Download dependencies
go mod download
` + "```" + `

## Development

` + "```bash" + `
# Run the application
go run .

# Build
go build -o bin/%s .

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
` + "```" + `

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
`, projectName, projectName)
}

func (t *GoTemplate) goMod(projectName string) string {
	return fmt.Sprintf(`module %s

go 1.23
`, projectName)
}

func (t *GoTemplate) mainGo() string {
	return `package main

import "fmt"

func main() {
	fmt.Println("Hello, Go!")
}
`
}

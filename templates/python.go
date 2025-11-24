package templates

import "fmt"

func init() {
	Register(&PythonTemplate{})
}

// PythonTemplate is a Python project template
type PythonTemplate struct{}

func (t *PythonTemplate) Name() string {
	return "python"
}

func (t *PythonTemplate) Description() string {
	return "Python project with pip/venv configuration"
}

func (t *PythonTemplate) Dependencies() []string {
	return []string{"git", "bd", "python3", "pip"}
}

func (t *PythonTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore":       t.gitignore(),
		".ubsignore":       t.ubsignore(),
		"README.md":        t.readme(projectName),
		"pyproject.toml":   t.pyproject(projectName),
		"requirements.txt": t.requirements(),
		"src/__init__.py":  "",
		"src/main.py":      t.mainPy(),
	}
}

func (t *PythonTemplate) gitignore() string {
	return `# Python
__pycache__/
*.py[cod]
*$py.class
*.so

# Virtual environments
venv/
.venv/
ENV/
env/

# Distribution
dist/
build/
*.egg-info/
*.egg

# IDE
.vscode/
.idea/
*.swp
*.swo

# Testing
.pytest_cache/
.coverage
htmlcov/
.tox/

# Environment
.env
.env.local

# OS
.DS_Store
Thumbs.db

# Type checking
.mypy_cache/
`
}

func (t *PythonTemplate) ubsignore() string {
	return `# UBS Scanner Ignore File
__pycache__/
venv/
.venv/
dist/
build/
*.egg-info/
.git/
.vscode/
.idea/
.beads/
.claude/
.pytest_cache/
htmlcov/
*.md
*.txt
*.toml
*.cfg
`
}

func (t *PythonTemplate) readme(projectName string) string {
	return fmt.Sprintf(`# %s

A Python project.

## Setup

` + "```bash" + `
# Create virtual environment
python -m venv venv

# Activate (Linux/Mac)
source venv/bin/activate

# Activate (Windows)
venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
` + "```" + `

## Development

` + "```bash" + `
# Run the application
python src/main.py

# Run tests
pytest

# Type checking
mypy src/
` + "```" + `

## Issue Tracking

` + "```bash" + `
bd list           # View issues
bd create --title "Task"  # Create issue
` + "```" + `

## Code Quality

` + "```bash" + `
ubs .             # Scan for bugs
` + "```" + `
`, projectName)
}

func (t *PythonTemplate) pyproject(projectName string) string {
	return fmt.Sprintf(`[build-system]
requires = ["setuptools>=61.0"]
build-backend = "setuptools.build_meta"

[project]
name = "%s"
version = "0.1.0"
description = ""
readme = "README.md"
requires-python = ">=3.9"
dependencies = []

[project.optional-dependencies]
dev = [
    "pytest",
    "mypy",
]
`, projectName)
}

func (t *PythonTemplate) requirements() string {
	return `# Core dependencies
# Add your dependencies here

# Development dependencies
pytest>=7.0.0
mypy>=1.0.0
`
}

func (t *PythonTemplate) mainPy() string {
	return `"""Main entry point."""


def main() -> None:
    """Run the application."""
    print("Hello, Python!")


if __name__ == "__main__":
    main()
`
}

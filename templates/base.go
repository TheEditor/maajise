package templates

import "fmt"

func init() {
	Register(&BaseTemplate{})
}

// BaseTemplate is the default language-agnostic template
type BaseTemplate struct{}

func (t *BaseTemplate) Name() string {
	return "base"
}

func (t *BaseTemplate) Description() string {
	return "Base template (language-agnostic)"
}

func (t *BaseTemplate) Dependencies() []string {
	return []string{"git", "bd"}
}

func (t *BaseTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore": t.gitignore(),
		".ubsignore": t.ubsignore(),
		"README.md":  t.readme(projectName),
	}
}

func (t *BaseTemplate) gitignore() string {
	return `# Dependencies
node_modules/
vendor/
packages/

# Build outputs
.next/
build/
dist/
target/
out/
*.o
*.exe

# Logs
*.log
logs/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Environment variables
.env
.env.local
.env*.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
desktop.ini

# Testing
coverage/
.nyc_output/

# Temporary files
*.tmp
*.temp
.cache/
`
}

func (t *BaseTemplate) ubsignore() string {
	return `# UBS Scanner Ignore File
# Excludes non-source files from bug scanning

# Dependencies
node_modules/
vendor/
packages/

# Build outputs
.next/
build/
dist/
target/
out/
bin/
obj/

# Version control & IDE
.git/
.vscode/
.idea/
.beads/
.claude/

# Documentation & metadata
docs/
history/
openspec/

# Scripts (usually not app code)
scripts/

# Static assets
public/
static/
assets/

# File types to skip
*.md
*.json
*.config.*
*.log
*.txt
*.lock
*.sum

# Environment & secrets
.env*
*.key
*.pem
*.cert
`
}

func (t *BaseTemplate) readme(projectName string) string {
	return fmt.Sprintf(`# %s

## Description

[Add project description here]

## Setup

`+"```bash"+`
# Clone the repository
git clone <repository-url>
cd %s

# [Add setup instructions here]
`+"```"+`

## Usage

[Add usage instructions here]

## Development

### Prerequisites

- [List prerequisites here]

### Running Locally

`+"```bash"+`
# [Add development commands here]
`+"```"+`

### Testing

`+"```bash"+`
# [Add testing commands here]
`+"```"+`

## Issue Tracking

This project uses [Beads](https://github.com/jfischoff/beads) for issue tracking.

`+"```bash"+`
# View all issues
bd list

# Create new issue
bd create --title "Issue title" --description "Issue description"

# View issue details
bd show <issue-id>
`+"```"+`

## Code Quality

This project uses [UBS (Ultimate Bug Scanner)](https://github.com/Dicklesworthstone/ultimate_bug_scanner) for static analysis.

`+"```bash"+`
# Run scanner on source code
ubs .

# Run with strict mode (fail on warnings)
ubs . --fail-on-warning
`+"```"+`

## Contributing

[Add contribution guidelines here]

## License

[Add license information here]
`, projectName, projectName)
}

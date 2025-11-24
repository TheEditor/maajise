package templates

import "fmt"

func init() {
	Register(&PHPTemplate{})
}

// PHPTemplate is a PHP project template
type PHPTemplate struct{}

func (t *PHPTemplate) Name() string {
	return "php"
}

func (t *PHPTemplate) Description() string {
	return "PHP project with Composer configuration"
}

func (t *PHPTemplate) Dependencies() []string {
	return []string{"git", "bd", "php", "composer"}
}

func (t *PHPTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore":    t.gitignore(),
		".ubsignore":    t.ubsignore(),
		"README.md":     t.readme(projectName),
		"composer.json": t.composerJSON(projectName),
		"src/index.php": t.indexPHP(),
	}
}

func (t *PHPTemplate) gitignore() string {
	return `# PHP
/vendor/
composer.lock

# IDE
.vscode/
.idea/
*.swp
*.swo
.phpunit.result.cache

# Environment
.env
.env.local

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Cache
/cache/
`
}

func (t *PHPTemplate) ubsignore() string {
	return `# UBS Scanner Ignore File
vendor/
.git/
.vscode/
.idea/
.beads/
.claude/
cache/
*.md
*.json
*.lock
*.log
`
}

func (t *PHPTemplate) readme(projectName string) string {
	return fmt.Sprintf(`# %s

A PHP project.

## Setup

` + "```bash" + `
# Install dependencies
composer install
` + "```" + `

## Development

` + "```bash" + `
# Run the application
php src/index.php

# Run tests
composer test

# Run PHP linter
composer lint
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

func (t *PHPTemplate) composerJSON(projectName string) string {
	return fmt.Sprintf(`{
    "name": "project/%s",
    "description": "",
    "type": "project",
    "require": {
        "php": ">=8.1"
    },
    "require-dev": {
        "phpunit/phpunit": "^10.0"
    },
    "autoload": {
        "psr-4": {
            "App\\": "src/"
        }
    },
    "scripts": {
        "test": "phpunit",
        "lint": "php -l src/"
    }
}
`, projectName)
}

func (t *PHPTemplate) indexPHP() string {
	return `<?php

declare(strict_types=1);

echo "Hello, PHP!\n";
`
}

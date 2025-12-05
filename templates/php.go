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
	return "PHP project with PSR-4 structure (public/, src/, MVC separation)"
}

func (t *PHPTemplate) Dependencies() []string {
	return []string{"git", "bd", "php", "composer"}
}

func (t *PHPTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore":                 t.gitignore(),
		".ubsignore":                 t.ubsignore(),
		"README.md":                  t.readme(projectName),
		"composer.json":              t.composerJSON(projectName),
		"bin/.gitkeep":               "",
		"config/.gitkeep":            "",
		"public/index.php":           t.indexPHP(projectName),
		"public/assets/.gitkeep":     "",
		"src/Controllers/.gitkeep":   "",
		"src/Models/.gitkeep":        "",
		"src/Services/.gitkeep":      "",
		"src/Middleware/.gitkeep":    "",
		"tests/.gitkeep":             "",
		"var/cache/.gitkeep":         "",
		"var/logs/.gitkeep":          "",
		"vendor/.gitkeep":            "",
	}
}

func (t *PHPTemplate) gitignore() string {
	return `# PHP
/vendor/

# Runtime data
/var/cache/*
!/var/cache/.gitkeep
/var/logs/*
!/var/logs/.gitkeep

# IDE
.vscode/
.idea/
*.swp
*.swo
.phpunit.result.cache

# Environment
.env
.env.local
.env.*.local

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Composer
composer.lock
`
}

func (t *PHPTemplate) ubsignore() string {
	return `# Dependencies
vendor/

# Runtime data
var/

# Static assets
public/assets/

# Tests (scan separately if needed)
tests/
`
}

func (t *PHPTemplate) readme(projectName string) string {
	return fmt.Sprintf(`# %s

## Project Structure

This project follows PSR-4 autoloading and modern PHP best practices.

| Directory | Purpose |
|-----------|---------|
| public/ | Web server document root (only publicly accessible files) |
| public/assets/ | Static files (CSS, JS, images) |
| src/ | Application source code (PSR-4 namespace: App\) |
| src/Controllers/ | HTTP request handlers |
| src/Models/ | Data models and entities |
| src/Services/ | Business logic layer |
| src/Middleware/ | Request/response middleware |
| config/ | Configuration files |
| tests/ | PHPUnit tests (namespace: Tests\) |
| var/cache/ | Application cache (gitignored) |
| var/logs/ | Application logs (gitignored) |
| bin/ | CLI scripts and tools |
| vendor/ | Composer dependencies (gitignored) |

## Setup

### Install dependencies
` + "```bash" + `
composer install
` + "```" + `

### Start development server
` + "```bash" + `
php -S localhost:8000 -t public
` + "```" + `

Then visit http://localhost:8000

## PSR-4 Autoloading

Classes in src/ use the App\ namespace:

- src/Controllers/HomeController.php → App\Controllers\HomeController
- src/Services/UserService.php → App\Services\UserService
- src/Models/User.php → App\Models\User

Example:
` + "```php" + `
<?php
namespace App\Controllers;

class HomeController {
    public function index(): void {
        echo "Hello from HomeController!";
    }
}
` + "```" + `

## Web Server Configuration

### Apache
Point DocumentRoot to the public/ directory.
Ensure mod_rewrite is enabled.

### Nginx
Set root to the public/ directory.
Configure try_files for routing.

## Development

### Run tests
` + "```bash" + `
composer test
` + "```" + `

### Lint PHP files
` + "```bash" + `
composer lint
` + "```" + `

### Check syntax
` + "```bash" + `
php -l src/Controllers/HomeController.php
` + "```" + `

## Issue Tracking

Track bugs and features with [bd](https://github.com/davefojtik/beads):

` + "```bash" + `
bd create --title "Bug: Fix user validation"
bd list --status open
bd ready  # Show issues ready to work on
` + "```" + `

## Code Quality

Scan for security issues with [ubs](https://github.com/davefojtik/ubs):

` + "```bash" + `
ubs scan .  # Analyze entire project
ubs fix <file>  # Interactive fix mode
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
    "autoload-dev": {
        "psr-4": {
            "Tests\\": "tests/"
        }
    },
    "scripts": {
        "test": "phpunit",
        "lint": "php -l src/"
    }
}
`, projectName)
}

func (t *PHPTemplate) indexPHP(projectName string) string {
	return `<?php

declare(strict_types=1);

// Load Composer autoloader
require_once __DIR__ . '/../vendor/autoload.php';

// Enable error reporting for development
error_reporting(E_ALL);
ini_set('display_errors', '1');

// Example: Load configuration
// $config = require __DIR__ . '/../config/app.php';

// Example: Initialize your application
echo "Hello, PHP!\n";
echo "Your application is running.\n";

// TODO: Set up routing, load controllers, etc.
// Example:
// use App\Controllers\HomeController;
// $controller = new HomeController();
// $controller->index();
`
}

# Maajise - Repository Initialization Tool

Go CLI for creating repos with Git, Beads, UBS.

## Installation

Open PowerShell as Administrator:

```powershell
cd "C:\Users\davef\Documents\Projects\Claude\Vibe Coding\tooling\maajise\maajise"
.\install.ps1
```

Restart terminal, use anywhere: `maajise my-project`

### Updating Existing Installations

If you have multiple Maajise installations (e.g., in both `Program Files` and `AppData`), use the update script:

```powershell
.\update-installations.ps1
```

To remove duplicate installations and keep only one:

```powershell
Remove-Item -Path "C:\Users\$env:USERNAME\AppData\Local\Programs\Maajise" -Recurse -Force
```

To check which Maajise is in your PATH:

```powershell
where maajise
```

## Usage

### Create New Project (Standard)

```bash
# Creates my-project/my-project/ structure
maajise my-project
```

### Initialize Current Directory

```bash
# Run inside existing folder
cd existing-project
maajise --in-place
```

### Safe Mode (No Overwrites)

```bash
# Won't overwrite existing files
maajise --in-place --no-overwrite
```

### Skip Options

```bash
# Skip interactive prompts
maajise my-project --skip-remote

# Skip Git initialization
maajise --in-place --skip-git

# Skip Beads
maajise my-project --skip-beads

# Skip initial commit
maajise my-project --skip-commit
```

### Verbose Output

```bash
maajise my-project -v
```

## Flags

```
--in-place          Initialize current directory (no nested folders)
--no-overwrite      Skip files that already exist
--template=<name>   Project template (base, typescript, python, rust, php, go)
--skip-git          Don't initialize Git
--skip-beads        Don't initialize Beads
--skip-commit       Don't create initial commit
--skip-remote       Don't prompt for remote setup
-v, --verbose       Verbose output
-h, --help          Show help
```

## Templates

Maajise supports multiple project templates:

| Template   | Description                          |
|------------|--------------------------------------|
| base       | Language-agnostic project structure  |
| typescript | TypeScript with npm/package.json     |
| python     | Python with pip/pyproject.toml       |
| rust       | Rust with Cargo                      |
| php        | PHP with Composer                    |
| go         | Go with go.mod                       |

### Using Templates

```bash
# List available templates
maajise templates

# Create project with specific template
maajise my-app --template=typescript
maajise my-service --template=python
maajise my-tool --template=rust
maajise my-project --template=php
maajise my-cli --template=go
```

## Examples

```bash
# Standard new project (default template)
maajise awesome-app

# TypeScript project
maajise my-app --template=typescript

# Python service
maajise my-service --template=python

# Rust CLI tool
maajise my-tool --template=rust

# Initialize maajise itself
cd maajise
maajise --in-place --no-overwrite

# Quick setup without prompts
maajise test-project --skip-remote --skip-beads

# Maximum control
maajise my-app --no-overwrite --skip-commit -v
```

## Requirements

- `git`
- `bd` (Beads) - optional with `--skip-beads`
- `ubs` - optional for scanning

## Configuration

Maajise can be configured via `~/.maajiserc` (YAML format):

```yaml
# Default values for init command
defaults:
  template: typescript
  git_name: Your Name
  git_email: you@example.com
  skip_remote: true
  main_branch: main

# Template variables (used in README, LICENSE, etc.)
variables:
  author: Your Name
  email: you@example.com
  year: "2025"
  license: MIT
  github: yourusername

# Custom templates directory (default: ~/.maajise/templates/)
templates_dir: ~/.maajise/templates
```

### Template Variables

Templates can include variables that are replaced during project creation:

- `{{.ProjectName}}` - Project name
- `{{.Author}}` - From config or empty
- `{{.Email}}` - From config or empty
- `{{.Year}}` - Current year (default) or from config
- `{{.License}}` - From config or "MIT"
- `{{.GitHub}}` - From config or empty

## Custom Templates

Create custom templates in `~/.maajise/templates/` as YAML files:

```yaml
# ~/.maajise/templates/my-template.yaml
name: my-template
description: My custom project template
dependencies:
  - git
  - node
files:
  .gitignore: |
    node_modules/
    .env
  README.md: |
    # {{.ProjectName}}

    Created by {{.Author}}
  package.json: |
    {
      "name": "{{.ProjectName}}",
      "version": "1.0.0"
    }
```

Use with: `maajise init my-project --template=my-template`

## Commands

| Command    | Description                              |
|------------|------------------------------------------|
| init       | Initialize a new project                 |
| add        | Add files or tooling to existing project |
| update     | Update configuration files               |
| validate   | Validate project setup                   |
| status     | Show quick project status                |
| templates  | List available templates                 |
| doctor     | Check system dependencies                |
| help       | Show help                                |
| version    | Show version                             |

### init

Initialize a new project with Git, Beads, and configuration files.

```bash
maajise init <project-name> [flags]
maajise <project-name>  # shorthand

Flags:
  --template=<name>   Project template (default: base)
  --in-place          Initialize in current directory
  --no-overwrite      Don't overwrite existing files
  --skip-git          Skip Git initialization
  --skip-beads        Skip Beads initialization
  --skip-commit       Skip initial commit
  --skip-remote       Skip remote setup prompt
  --skip-git-user     Skip Git user configuration
  --git-name=<name>   Git user.name (non-interactive)
  --git-email=<email> Git user.email (non-interactive)
  --dry-run           Preview what would be created without making changes
  --interactive, -i   Interactive mode with prompts
  -v, --verbose       Verbose output

Examples:
  maajise init my-project
  maajise init my-app --template=typescript
  maajise init my-project --dry-run
  maajise init --interactive
  maajise init --in-place --template=python
```

### add

Add files or tooling to an existing project.

```bash
maajise add <item> [flags]

Tooling:
  git                Initialize Git repository
  beads              Initialize Beads issue tracking
  ubs                Add .ubsignore file

Files:
  .gitignore         Add .gitignore from template
  .ubsignore         Add .ubsignore from template
  readme             Add README.md from template

Flags:
  --force            Overwrite existing files
  --template=<name>  Template for file content (auto-detects if not specified)
  --dry-run          Preview without making changes
  -v, --verbose      Verbose output

Examples:
  maajise add git
  maajise add beads
  maajise add .gitignore
  maajise add .gitignore --template=typescript
  maajise add readme --force
  maajise add --dry-run git
```

### doctor

Check system dependencies and configuration.

```bash
maajise doctor [flags]

Flags:
  -v, --verbose   Verbose output (show error details)

Examples:
  maajise doctor
  maajise doctor --verbose
```

Reports on:
- Required dependencies: `git`, `bd` (Beads)
- Optional dependencies: `ubs`, `go`
- Configuration file status
- Custom templates directory

### update

Update configuration files in an existing project.

```bash
maajise update [flags] [files...]

Flags:
  --force             Overwrite existing files
  --template=<name>   Template to use (auto-detects if not specified)
  --dry-run           Show what would be updated
  -v, --verbose       Verbose output

Examples:
  maajise update                    # Update all files (skip existing)
  maajise update --force            # Update all files (overwrite)
  maajise update .gitignore         # Update specific file
  maajise update --dry-run          # Preview changes
```

### validate

Validate project setup and configuration.

```bash
maajise validate [flags]

Flags:
  --strict            Treat warnings as failures
  -v, --verbose       Verbose output

Examples:
  maajise validate
  maajise validate --strict
  maajise validate --verbose
```

### status

Show quick project status.

```bash
maajise status

Examples:
  maajise status
```

### templates

List available project templates.

```bash
maajise templates
```

## Quick Reference

```bash
# Create new project
maajise init my-project --template=typescript

# Check project health
maajise validate

# Quick status
maajise status

# Update config files
maajise update --force

# List templates
maajise templates
```

## What It Creates

- `.git/` - Git repo
- `.beads/` - Issue tracking
- `.ubsignore` - UBS scanner config
- `.gitignore` - Git ignores
- `README.md` - Template

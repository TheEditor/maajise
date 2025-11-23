# Maajise - Repository Initialization Tool

Go CLI for creating repos with Git, Beads, UBS.

## Installation

Open PowerShell as Administrator:

```powershell
cd "C:\Users\davef\Documents\Projects\Claude\Vibe Coding\tooling\maajise\maajise"
.\install.ps1
```

Restart terminal, use anywhere: `maajise my-project`

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
--skip-git          Don't initialize Git
--skip-beads        Don't initialize Beads
--skip-commit       Don't create initial commit
--skip-remote       Don't prompt for remote setup
-v, --verbose       Verbose output
-h, --help          Show help
```

## Examples

```bash
# Standard new project
maajise awesome-app

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

## What It Creates

- `.git/` - Git repo
- `.beads/` - Issue tracking
- `.ubsignore` - UBS scanner config
- `.gitignore` - Git ignores
- `README.md` - Template

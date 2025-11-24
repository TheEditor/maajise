#!/usr/bin/env bash
set -euo pipefail

# Repository Initialization Script
# Sets up new repo with UBS, Beads, and Git
# Usage: ./init-repo.sh <project-name>

VERSION="1.0.0"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RESET='\033[0m'

# Output helpers
error() { echo -e "${RED}✗${RESET} $*" >&2; }
success() { echo -e "${GREEN}✓${RESET} $*"; }
info() { echo -e "${BLUE}→${RESET} $*"; }
warn() { echo -e "${YELLOW}⚠${RESET} $*"; }

# Show usage
show_usage() {
    cat << EOF
Repository Initialization Script v${VERSION}

Usage: $0 <project-name>

Creates a new repository with:
  - Nested folder structure (project/project/)
  - Git initialization
  - Beads issue tracking
  - UBS scanner configuration
  - Standard files (.gitignore, README.md, etc.)

Example:
  $0 my-awesome-project

  Creates:
    my-awesome-project/
    └── my-awesome-project/
        ├── .git/
        ├── .beads/
        ├── .ubsignore
        ├── .gitignore
        └── README.md

Requirements:
  - git
  - beads (bd command)
  - ubs (optional, for scanning)

EOF
}

# Check if required tools are installed
check_dependencies() {
    local missing=()
    
    if ! command -v git >/dev/null 2>&1; then
        missing+=("git")
    fi
    
    if ! command -v bd >/dev/null 2>&1; then
        missing+=("beads (bd command)")
    fi
    
    if [ ${#missing[@]} -gt 0 ]; then
        error "Missing required tools: ${missing[*]}"
        info "Install them before running this script"
        return 1
    fi
    
    return 0
}

# Validate project name
validate_project_name() {
    local name="$1"
    
    # Check for empty name
    if [ -z "$name" ]; then
        error "Project name cannot be empty"
        return 1
    fi
    
    # Check for invalid characters (allow letters, numbers, hyphens, underscores)
    if [[ ! "$name" =~ ^[a-zA-Z0-9_-]+$ ]]; then
        error "Project name can only contain letters, numbers, hyphens, and underscores"
        return 1
    fi
    
    # Check if outer directory already exists
    if [ -d "$name" ]; then
        error "Directory '$name' already exists"
        return 1
    fi
    
    return 0
}

# Create directory structure
create_structure() {
    local project_name="$1"
    
    info "Creating directory structure..."
    
    # Create outer folder
    mkdir -p "$project_name"
    
    # Create inner folder (repo root)
    mkdir -p "$project_name/$project_name"
    
    success "Created $project_name/$project_name/"
}

# Initialize Git repository
init_git() {
    local repo_dir="$1"
    
    info "Initializing Git repository..."
    
    cd "$repo_dir"
    
    if git init >/dev/null 2>&1; then
        success "Git repository initialized"
        return 0
    else
        error "Failed to initialize Git repository"
        return 1
    fi
}

# Initialize Beads issue tracker
init_beads() {
    local repo_dir="$1"
    
    info "Initializing Beads issue tracker..."
    
    cd "$repo_dir"
    
    # Check if bd init works
    if bd init >/dev/null 2>&1; then
        success "Beads initialized"
        return 0
    else
        warn "Beads initialization failed (may need to run 'bd init' manually)"
        return 0  # Don't fail the whole script
    fi
}

# Create .ubsignore file
create_ubsignore() {
    local repo_dir="$1"
    
    info "Creating .ubsignore..."
    
    cat > "$repo_dir/.ubsignore" << 'EOF'
# UBS Scanner Ignore File
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
EOF
    
    success "Created .ubsignore"
}

# Create .gitignore file
create_gitignore() {
    local repo_dir="$1"
    
    info "Creating .gitignore..."
    
    cat > "$repo_dir/.gitignore" << 'EOF'
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
EOF
    
    success "Created .gitignore"
}

# Create README.md
create_readme() {
    local repo_dir="$1"
    local project_name="$2"
    
    info "Creating README.md..."
    
    cat > "$repo_dir/README.md" << EOF
# ${project_name}

## Description

[Add project description here]

## Setup

\`\`\`bash
# Clone the repository
git clone <repository-url>
cd ${project_name}

# [Add setup instructions here]
\`\`\`

## Usage

[Add usage instructions here]

## Development

### Prerequisites

- [List prerequisites here]

### Running Locally

\`\`\`bash
# [Add development commands here]
\`\`\`

### Testing

\`\`\`bash
# [Add testing commands here]
\`\`\`

## Issue Tracking

This project uses [Beads](https://github.com/jfischoff/beads) for issue tracking.

\`\`\`bash
# View all issues
bd list

# Create new issue
bd create --title "Issue title" --description "Issue description"

# View issue details
bd show <issue-id>
\`\`\`

## Code Quality

This project uses [UBS (Ultimate Bug Scanner)](https://github.com/Dicklesworthstone/ultimate_bug_scanner) for static analysis.

\`\`\`bash
# Run scanner on source code
ubs .

# Run with strict mode (fail on warnings)
ubs . --fail-on-warning
\`\`\`

## Contributing

[Add contribution guidelines here]

## License

[Add license information here]
EOF
    
    success "Created README.md"
}

# Create initial git commit
create_initial_commit() {
    local repo_dir="$1"
    
    info "Creating initial commit..."
    
    cd "$repo_dir"
    
    git add .ubsignore .gitignore README.md
    
    if git commit -m "Initial commit

- Add .ubsignore for UBS scanner
- Add .gitignore for version control
- Add README.md with project structure
- Initialize Beads issue tracking" >/dev/null 2>&1; then
        success "Initial commit created"
        return 0
    else
        error "Failed to create initial commit"
        return 1
    fi
}

# Setup git remote
setup_git_remote() {
    local repo_dir="$1"
    local project_name="$2"
    
    cd "$repo_dir"
    
    echo ""
    info "Git remote setup (optional)"
    echo ""
    
    read -p "$(echo -e "${BLUE}?${RESET} Add git remote? (y/N): ")" -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo ""
        info "Enter remote URL (e.g., https://github.com/username/${project_name}.git)"
        read -p "$(echo -e "${BLUE}→${RESET} Remote URL: ")" -r remote_url
        
        if [ -n "$remote_url" ]; then
            if git remote add origin "$remote_url" 2>/dev/null; then
                success "Added remote: origin → $remote_url"
                echo ""
                info "Push to remote with:"
                echo "  cd $project_name/$project_name"
                echo "  git push -u origin main"
                echo ""
            else
                warn "Failed to add remote (may already exist)"
            fi
        else
            info "Skipped remote setup"
        fi
    else
        info "Skipped remote setup"
    fi
}

# Show completion summary
show_summary() {
    local project_name="$1"
    local repo_path="$2"
    
    echo ""
    echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${RESET}"
    echo -e "${GREEN}║  ✓ Repository initialized successfully!                   ║${RESET}"
    echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${RESET}"
    echo ""
    success "Project: $project_name"
    success "Location: $repo_path"
    echo ""
    info "Next steps:"
    echo "  1. cd $project_name/$project_name"
    echo "  2. Create your project files"
    echo "  3. Run 'ubs .' to scan for issues"
    echo "  4. Run 'bd list' to manage tasks"
    echo ""
    info "Quick commands:"
    echo "  bd create --title \"Task name\"    # Create new task"
    echo "  bd list                           # View all tasks"
    echo "  ubs .                             # Scan for bugs"
    echo ""
}

# Main function
main() {
    # Parse arguments
    if [ $# -eq 0 ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        show_usage
        exit 0
    fi
    
    local project_name="$1"
    
    echo ""
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════╗${RESET}"
    echo -e "${BLUE}║  Repository Initialization Script v${VERSION}                ║${RESET}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${RESET}"
    echo ""
    
    # Validate
    info "Validating project name: $project_name"
    if ! validate_project_name "$project_name"; then
        exit 1
    fi
    success "Project name valid"
    echo ""
    
    # Check dependencies
    info "Checking dependencies..."
    if ! check_dependencies; then
        exit 1
    fi
    success "All dependencies available"
    echo ""
    
    # Create structure
    if ! create_structure "$project_name"; then
        exit 1
    fi
    
    local repo_path="$PWD/$project_name/$project_name"
    
    # Initialize Git
    if ! init_git "$repo_path"; then
        exit 1
    fi
    
    # Initialize Beads
    if ! init_beads "$repo_path"; then
        warn "Beads initialization may need manual intervention"
    fi
    
    # Create standard files
    create_ubsignore "$repo_path"
    create_gitignore "$repo_path"
    create_readme "$repo_path" "$project_name"
    
    # Initial commit
    if ! create_initial_commit "$repo_path"; then
        exit 1
    fi
    
    # Setup remote (optional)
    setup_git_remote "$repo_path" "$project_name"
    
    # Show summary
    show_summary "$project_name" "$repo_path"
}

# Run main function
main "$@"

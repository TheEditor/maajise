package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"maajise/internal/git"
	"maajise/internal/validate"
)

type InitCommand struct {
	fs     *flag.FlagSet
	config Config
}

func NewInitCommand() *InitCommand {
	ic := &InitCommand{
		fs:     flag.NewFlagSet("init", flag.ContinueOnError),
		config: DefaultConfig(),
	}

	// Define flags
	ic.fs.BoolVar(&ic.config.InPlace, "in-place", false, "Initialize in current directory")
	ic.fs.BoolVar(&ic.config.NoOverwrite, "no-overwrite", false, "Don't overwrite existing files")
	ic.fs.BoolVar(&ic.config.SkipGit, "skip-git", false, "Skip Git initialization")
	ic.fs.BoolVar(&ic.config.SkipBeads, "skip-beads", false, "Skip Beads initialization")
	ic.fs.BoolVar(&ic.config.SkipCommit, "skip-commit", false, "Skip initial commit")
	ic.fs.BoolVar(&ic.config.SkipRemote, "skip-remote", false, "Skip remote setup")
	ic.fs.BoolVar(&ic.config.SkipGitUser, "skip-git-user", false, "Skip Git user configuration")
	ic.fs.StringVar(&ic.config.GitName, "git-name", "", "Git user.name (non-interactive)")
	ic.fs.StringVar(&ic.config.GitEmail, "git-email", "", "Git user.email (non-interactive)")
	ic.fs.BoolVar(&ic.config.Verbose, "v", false, "Verbose output")
	ic.fs.BoolVar(&ic.config.Verbose, "verbose", false, "Verbose output")

	return ic
}

func (ic *InitCommand) Name() string {
	return "init"
}

func (ic *InitCommand) Description() string {
	return "Initialize a new project with Git, Beads, and configuration files"
}

func (ic *InitCommand) Usage() string {
	return "maajise init [flags] [project-name]"
}

func (ic *InitCommand) Examples() []string {
	return []string{
		"maajise init my-project",
		"maajise init --in-place",
		"maajise init my-project --no-git",
	}
}

func (ic *InitCommand) Execute(args []string) error {
	// Parse flags
	if err := ic.fs.Parse(args); err != nil {
		return err
	}

	// Get project name from remaining args
	remainingArgs := ic.fs.Args()
	if !ic.config.InPlace {
		if len(remainingArgs) == 0 {
			return fmt.Errorf("project name required")
		}
		ic.config.ProjectName = remainingArgs[0]
	} else {
		// In-place mode: use current directory name
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		ic.config.ProjectName = filepath.Base(cwd)
	}

	// Validate project name
	if err := ic.validateProjectName(ic.config.ProjectName); err != nil {
		return err
	}

	// Execute initialization
	return ic.runInit()
}

func (ic *InitCommand) validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("empty name")
	}

	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	if !match {
		return fmt.Errorf("invalid characters in project name")
	}

	return nil
}

func (ic *InitCommand) runInit() error {
	// Create/determine project directory
	repoPath, err := ic.createStructure()
	if err != nil {
		return err
	}

	// Initialize Git
	if err := ic.initGit(repoPath); err != nil {
		return err
	}

	// Configure Git user
	if err := ic.configureGitUser(repoPath); err != nil {
		return err
	}

	// Initialize Beads
	ic.initBeads(repoPath)

	// Create standard files
	ic.createUbsignore(repoPath)
	ic.createGitignore(repoPath)
	ic.createReadme(repoPath)

	// Initial commit
	if err := ic.createInitialCommit(repoPath); err != nil {
		return err
	}

	// Setup remote (optional)
	ic.setupGitRemote(repoPath)

	// Show summary
	ic.showSummary(repoPath)

	return nil
}

func (ic *InitCommand) createStructure() (string, error) {
	if ic.config.InPlace {
		// Use current directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		if ic.config.Verbose {
			fmt.Printf("→ Using current directory: %s\n", cwd)
		}
		return cwd, nil
	}

	// Standard mode: create nested structure
	fmt.Println("→ Creating directory structure...")
	innerPath := filepath.Join(ic.config.ProjectName, ic.config.ProjectName)

	if ic.fileExists(ic.config.ProjectName) {
		return "", fmt.Errorf("directory '%s' already exists", ic.config.ProjectName)
	}

	if err := os.MkdirAll(innerPath, 0755); err != nil {
		return "", err
	}

	fmt.Printf("✓ Created %s/\n", innerPath)
	return innerPath, nil
}

func (ic *InitCommand) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (ic *InitCommand) initGit(repoDir string) error {
	if ic.config.SkipGit {
		if ic.config.Verbose {
			fmt.Println("→ Skipping Git (--skip-git)")
		}
		return nil
	}

	if err := git.Init(repoDir, ic.config.Verbose); err != nil {
		return err
	}

	return nil
}

func (ic *InitCommand) configureGitUser(repoDir string) error {
	if ic.config.SkipGit || ic.config.SkipGitUser {
		if ic.config.Verbose {
			fmt.Println("→ Skipping Git user configuration")
		}
		return nil
	}

	var userName, userEmail string

	// Check if values provided via flags
	if ic.config.GitName != "" && ic.config.GitEmail != "" {
		userName = ic.config.GitName
		userEmail = ic.config.GitEmail

		// Validate email format
		if !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") {
			return fmt.Errorf("invalid email format")
		}

		if ic.config.Verbose {
			fmt.Printf("→ Using provided Git config: %s <%s>\n", userName, userEmail)
		}
	} else {
		// Interactive prompts
		fmt.Println("→ Configuring Git user...")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)

		// Get user.name
		fmt.Print("→ Git user.name (full name): ")
		var err error
		userName, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user name: %w", err)
		}
		userName = strings.TrimSpace(userName)

		if userName == "" {
			return fmt.Errorf("git user.name cannot be empty")
		}

		// Get user.email
		fmt.Print("→ Git user.email (email address): ")
		userEmail, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user email: %w", err)
		}
		userEmail = strings.TrimSpace(userEmail)

		if userEmail == "" {
			return fmt.Errorf("git user.email cannot be empty")
		}

		// Validate email format (basic check)
		if !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") {
			return fmt.Errorf("invalid email format")
		}
	}

	// Set user.name
	if err := git.SetConfig(repoDir, "user.name", userName, ic.config.Verbose); err != nil {
		return err
	}

	// Set user.email
	if err := git.SetConfig(repoDir, "user.email", userEmail, ic.config.Verbose); err != nil {
		return err
	}

	fmt.Printf("✓ Git user configured: %s <%s>\n", userName, userEmail)
	fmt.Println()
	return nil
}

func (ic *InitCommand) initBeads(repoDir string) {
	if ic.config.SkipBeads {
		if ic.config.Verbose {
			fmt.Println("→ Skipping Beads (--skip-beads)")
		}
		return
	}

	// Check if already initialized
	beadsDir := filepath.Join(repoDir, ".beads")
	if ic.fileExists(beadsDir) {
		if ic.config.NoOverwrite {
			fmt.Println("⚠ Beads already initialized (--no-overwrite)")
			return
		}
		fmt.Println("⚠ Beads already initialized")
		return
	}

	fmt.Println("→ Initializing Beads...")

	cmd := exec.Command("bd", "init")
	cmd.Dir = repoDir
	if ic.config.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("⚠ Beads init failed (run 'bd init' manually)")
		return
	}

	fmt.Println("✓ Beads initialized")
}

func (ic *InitCommand) writeFileIfNotExists(path, content string) error {
	if ic.fileExists(path) {
		if ic.config.NoOverwrite {
			fmt.Printf("⚠ Skipped %s (exists, --no-overwrite)\n", filepath.Base(path))
			return nil
		}
		fmt.Printf("⚠ Overwriting %s\n", filepath.Base(path))
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	fmt.Printf("✓ Created %s\n", filepath.Base(path))
	return nil
}

func (ic *InitCommand) createUbsignore(repoDir string) {
	fmt.Println("→ Creating .ubsignore...")

	content := `# UBS Scanner Ignore File
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

	path := filepath.Join(repoDir, ".ubsignore")
	_ = ic.writeFileIfNotExists(path, content)
}

func (ic *InitCommand) createGitignore(repoDir string) {
	fmt.Println("→ Creating .gitignore...")

	content := `# Dependencies
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

	path := filepath.Join(repoDir, ".gitignore")
	_ = ic.writeFileIfNotExists(path, content)
}

func (ic *InitCommand) createReadme(repoDir string) {
	fmt.Println("→ Creating README.md...")

	content := fmt.Sprintf(`# %s

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
`, ic.config.ProjectName, ic.config.ProjectName)

	path := filepath.Join(repoDir, "README.md")
	_ = ic.writeFileIfNotExists(path, content)
}

func (ic *InitCommand) createInitialCommit(repoDir string) error {
	if ic.config.SkipGit || ic.config.SkipCommit {
		if ic.config.Verbose {
			fmt.Println("→ Skipping initial commit")
		}
		return nil
	}

	fmt.Println("→ Creating initial commit...")

	// Check if there's anything to commit
	hasChanges, err := git.HasChanges(repoDir)
	if err != nil {
		return err
	}

	if !hasChanges {
		fmt.Println("⚠ Nothing to commit")
		return nil
	}

	// Add files
	if err := git.AddFiles(repoDir, []string{".ubsignore", ".gitignore", "README.md"}, ic.config.Verbose); err != nil {
		return err
	}

	commitMsg := `Initial commit

- Add .ubsignore for UBS scanner
- Add .gitignore for version control
- Add README.md with project structure
- Initialize Beads issue tracking`

	// Create commit
	if err := git.CreateCommit(repoDir, commitMsg, ic.config.Verbose); err != nil {
		return err
	}

	fmt.Println("✓ Initial commit created")
	return nil
}

func (ic *InitCommand) setupGitRemote(repoDir string) {
	if ic.config.SkipGit || ic.config.SkipRemote {
		if ic.config.Verbose {
			fmt.Println("→ Skipping remote setup")
		}
		return
	}

	// Check if remote already exists
	existingURL, err := git.GetRemote(repoDir, "origin")
	if err == nil {
		// Remote already exists
		fmt.Printf("⚠ Remote 'origin' already exists: %s\n", existingURL)
		return
	}

	fmt.Println()
	fmt.Println("→ Git remote setup (optional)")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("? Add git remote? (y/N): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	response = strings.TrimSpace(response)

	if strings.ToLower(response) != "y" {
		fmt.Println("→ Skipped remote setup")
		return
	}

	fmt.Println()
	fmt.Printf("→ Enter remote URL (e.g., https://github.com/username/%s.git)\n", ic.config.ProjectName)
	fmt.Print("→ Remote URL: ")
	remoteURL, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	remoteURL = strings.TrimSpace(remoteURL)

	if remoteURL == "" {
		fmt.Println("→ Skipped remote setup")
		return
	}

	// Validate the git remote URL
	if err := validate.ValidateGitURL(remoteURL); err != nil {
		fmt.Printf("✗ Invalid git remote URL: %v\n", err)
		return
	}

	if err := git.AddRemote(repoDir, "origin", remoteURL, ic.config.Verbose); err != nil {
		fmt.Println("⚠ Failed to add remote (may already exist)")
		return
	}

	fmt.Printf("✓ Added remote: origin → %s\n", remoteURL)
	fmt.Println()
	fmt.Println("→ Push to remote with:")
	if ic.config.InPlace {
		fmt.Println("  git push -u origin main")
	} else {
		fmt.Printf("  cd %s/%s\n", ic.config.ProjectName, ic.config.ProjectName)
		fmt.Println("  git push -u origin main")
	}
	fmt.Println()
}

func (ic *InitCommand) showSummary(repoPath string) {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  ✓ Repository initialized successfully!                   ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("✓ Project: %s\n", ic.config.ProjectName)
	fmt.Printf("✓ Location: %s\n", repoPath)
	fmt.Println()

	if !ic.config.InPlace {
		fmt.Println("→ Next steps:")
		fmt.Printf("  1. cd %s/%s\n", ic.config.ProjectName, ic.config.ProjectName)
		fmt.Println("  2. Create your project files")
	} else {
		fmt.Println("→ Next steps:")
		fmt.Println("  1. Create your project files")
	}

	fmt.Println("  2. Run 'ubs .' to scan for issues")
	fmt.Println("  3. Run 'bd list' to manage tasks")
	fmt.Println()
	fmt.Println("→ Quick commands:")
	fmt.Println("  bd create --title \"Task name\"    # Create new task")
	fmt.Println("  bd list                           # View all tasks")
	fmt.Println("  ubs .                             # Scan for bugs")
	fmt.Println()
}

func init() {
	Register(NewInitCommand())
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"maajise/internal/validate"
)

const VERSION = "2.0.0"

// ANSI colors
const (
	red    = "\033[0;31m"
	green  = "\033[0;32m"
	yellow = "\033[1;33m"
	blue   = "\033[0;34m"
	reset  = "\033[0m"
)

// Config holds runtime configuration
type Config struct {
	ProjectName  string
	InPlace      bool
	NoOverwrite  bool
	SkipGit      bool
	SkipBeads    bool
	SkipCommit   bool
	SkipRemote   bool
	SkipGitUser  bool
	GitName      string
	GitEmail     string
	Verbose      bool
}

func errorMsg(msg string) {
	fmt.Fprintf(os.Stderr, "%s✗%s %s\n", red, reset, msg)
}

func success(msg string) {
	fmt.Printf("%s✓%s %s\n", green, reset, msg)
}

func info(msg string) {
	fmt.Printf("%s→%s %s\n", blue, reset, msg)
}

func warn(msg string) {
	fmt.Printf("%s⚠%s %s\n", yellow, reset, msg)
}

func showUsage() {
	fmt.Printf(`Maajise - Repository Initialization Tool v%s

Usage: 
  maajise [flags] <project-name>
  maajise [flags] --in-place

Flags:
  --in-place          Initialize in current directory (no nested folders)
  --no-overwrite      Skip files that already exist
  --skip-git          Don't initialize Git
  --skip-beads        Don't initialize Beads
  --skip-commit       Don't create initial commit
  --skip-remote       Don't prompt for remote setup
  --skip-git-user     Don't prompt for Git user config
  --git-name NAME     Set Git user.name (non-interactive)
  --git-email EMAIL   Set Git user.email (non-interactive)
  -v, --verbose       Verbose output
  -h, --help          Show this help

Examples:
  # Standard: Create project/project/ structure
  maajise my-project

  # Initialize current directory
  maajise --in-place

  # Initialize with safety (don't overwrite)
  maajise --in-place --no-overwrite

  # Quick init without interactive prompts
  maajise my-project --skip-remote

Creates:
  - Git repository (.git/)
  - Beads issue tracking (.beads/)
  - UBS ignore config (.ubsignore)
  - Git ignore (.gitignore)
  - README template

Requirements:
  - git
  - bd (Beads)
  - ubs (optional)

`, VERSION)
}

func parseFlags() *Config {
	cfg := &Config{}

	flag.BoolVar(&cfg.InPlace, "in-place", false, "Initialize in current directory")
	flag.BoolVar(&cfg.NoOverwrite, "no-overwrite", false, "Skip existing files")
	flag.BoolVar(&cfg.SkipGit, "skip-git", false, "Skip Git initialization")
	flag.BoolVar(&cfg.SkipBeads, "skip-beads", false, "Skip Beads initialization")
	flag.BoolVar(&cfg.SkipCommit, "skip-commit", false, "Skip initial commit")
	flag.BoolVar(&cfg.SkipRemote, "skip-remote", false, "Skip remote setup")
	flag.BoolVar(&cfg.SkipGitUser, "skip-git-user", false, "Skip Git user configuration")
	flag.StringVar(&cfg.GitName, "git-name", "", "Git user.name (non-interactive)")
	flag.StringVar(&cfg.GitEmail, "git-email", "", "Git user.email (non-interactive)")
	flag.BoolVar(&cfg.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Verbose output")

	help := flag.Bool("h", false, "Show help")
	helpLong := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help || *helpLong {
		showUsage()
		os.Exit(0)
	}

	// Get project name from args
	args := flag.Args()
	if !cfg.InPlace {
		if len(args) == 0 {
			errorMsg("Project name required (or use --in-place)")
			fmt.Println()
			showUsage()
			os.Exit(1)
		}
		cfg.ProjectName = args[0]
	} else {
		// In-place mode: use current directory name
		cwd, err := os.Getwd()
		if err != nil {
			errorMsg("Failed to get current directory")
			os.Exit(1)
		}
		cfg.ProjectName = filepath.Base(cwd)
	}

	return cfg
}

func checkDependencies(cfg *Config) error {
	missing := []string{}

	if !cfg.SkipGit {
		if _, err := exec.LookPath("git"); err != nil {
			missing = append(missing, "git")
		}
	}

	if !cfg.SkipBeads {
		if _, err := exec.LookPath("bd"); err != nil {
			missing = append(missing, "beads (bd)")
		}
	}

	if len(missing) > 0 {
		errorMsg(fmt.Sprintf("Missing: %s", strings.Join(missing, ", ")))
		info("Install them or use --skip-* flags")
		return fmt.Errorf("missing dependencies")
	}

	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("empty name")
	}

	// Pattern is a compile-time constant, error is safe to ignore
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	if !match {
		errorMsg("Project name: letters, numbers, hyphens, underscores only")
		return fmt.Errorf("invalid characters")
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createStructure(cfg *Config) (string, error) {
	if cfg.InPlace {
		// Use current directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		if cfg.Verbose {
			info(fmt.Sprintf("Using current directory: %s", cwd))
		}
		return cwd, nil
	}

	// Standard mode: create nested structure
	info("Creating directory structure...")
	innerPath := filepath.Join(cfg.ProjectName, cfg.ProjectName)

	if fileExists(cfg.ProjectName) {
		errorMsg(fmt.Sprintf("Directory '%s' already exists", cfg.ProjectName))
		return "", fmt.Errorf("directory exists")
	}

	if err := os.MkdirAll(innerPath, 0755); err != nil {
		return "", err
	}

	success(fmt.Sprintf("Created %s/", innerPath))
	return innerPath, nil
}

func initGit(repoDir string, cfg *Config) error {
	if cfg.SkipGit {
		if cfg.Verbose {
			info("Skipping Git (--skip-git)")
		}
		return nil
	}

	// Check if already a git repo
	gitDir := filepath.Join(repoDir, ".git")
	if fileExists(gitDir) {
		if cfg.NoOverwrite {
			warn("Git already initialized (--no-overwrite)")
			return nil
		}
		warn("Git already initialized")
		return nil
	}

	info("Initializing Git...")

	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	if cfg.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		errorMsg("Git init failed")
		return err
	}

	success("Git initialized")
	return nil
}

func configureGitUser(repoDir string, cfg *Config) error {
	if cfg.SkipGit || cfg.SkipGitUser {
		if cfg.Verbose {
			info("Skipping Git user configuration")
		}
		return nil
	}

	var userName, userEmail string

	// Check if values provided via flags
	if cfg.GitName != "" && cfg.GitEmail != "" {
		userName = cfg.GitName
		userEmail = cfg.GitEmail

		// Validate email format
		if !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") {
			errorMsg("Invalid email format")
			return fmt.Errorf("invalid email format")
		}

		if cfg.Verbose {
			info(fmt.Sprintf("Using provided Git config: %s <%s>", userName, userEmail))
		}
	} else {
		// Interactive prompts
		info("Configuring Git user...")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)

		// Get user.name
		fmt.Printf("%s→%s Git user.name (full name): ", blue, reset)
		var err error
		userName, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user name: %w", err)
		}
		userName = strings.TrimSpace(userName)

		if userName == "" {
			errorMsg("Git user.name cannot be empty")
			return fmt.Errorf("empty user name")
		}

		// Get user.email
		fmt.Printf("%s→%s Git user.email (email address): ", blue, reset)
		userEmail, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user email: %w", err)
		}
		userEmail = strings.TrimSpace(userEmail)

		if userEmail == "" {
			errorMsg("Git user.email cannot be empty")
			return fmt.Errorf("empty user email")
		}

		// Validate email format (basic check)
		if !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") {
			errorMsg("Invalid email format")
			return fmt.Errorf("invalid email format")
		}
	}

	// Set user.name
	nameCmd := exec.Command("git", "config", "user.name", userName)
	nameCmd.Dir = repoDir
	if cfg.Verbose {
		nameCmd.Stdout = os.Stdout
		nameCmd.Stderr = os.Stderr
	}
	if err := nameCmd.Run(); err != nil {
		errorMsg("Failed to set user.name")
		return err
	}

	// Set user.email
	emailCmd := exec.Command("git", "config", "user.email", userEmail)
	emailCmd.Dir = repoDir
	if cfg.Verbose {
		emailCmd.Stdout = os.Stdout
		emailCmd.Stderr = os.Stderr
	}
	if err := emailCmd.Run(); err != nil {
		errorMsg("Failed to set user.email")
		return err
	}

	success(fmt.Sprintf("Git user configured: %s <%s>", userName, userEmail))
	fmt.Println()
	return nil
}

func initBeads(repoDir string, cfg *Config) error {
	if cfg.SkipBeads {
		if cfg.Verbose {
			info("Skipping Beads (--skip-beads)")
		}
		return nil
	}

	// Check if already initialized
	beadsDir := filepath.Join(repoDir, ".beads")
	if fileExists(beadsDir) {
		if cfg.NoOverwrite {
			warn("Beads already initialized (--no-overwrite)")
			return nil
		}
		warn("Beads already initialized")
		return nil
	}

	info("Initializing Beads...")

	cmd := exec.Command("bd", "init")
	cmd.Dir = repoDir
	if cfg.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		warn("Beads init failed (run 'bd init' manually)")
		return nil
	}

	success("Beads initialized")
	return nil
}

func writeFileIfNotExists(path, content string, cfg *Config) error {
	if fileExists(path) {
		if cfg.NoOverwrite {
			warn(fmt.Sprintf("Skipped %s (exists, --no-overwrite)", filepath.Base(path)))
			return nil
		}
		warn(fmt.Sprintf("Overwriting %s", filepath.Base(path)))
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	success(fmt.Sprintf("Created %s", filepath.Base(path)))
	return nil
}

func createUbsignore(repoDir string, cfg *Config) error {
	info("Creating .ubsignore...")

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
	return writeFileIfNotExists(path, content, cfg)
}

func createGitignore(repoDir string, cfg *Config) error {
	info("Creating .gitignore...")

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
	return writeFileIfNotExists(path, content, cfg)
}

func createReadme(repoDir string, cfg *Config) error {
	info("Creating README.md...")

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
`, cfg.ProjectName, cfg.ProjectName)

	path := filepath.Join(repoDir, "README.md")
	return writeFileIfNotExists(path, content, cfg)
}

func createInitialCommit(repoDir string, cfg *Config) error {
	if cfg.SkipGit || cfg.SkipCommit {
		if cfg.Verbose {
			info("Skipping initial commit")
		}
		return nil
	}

	info("Creating initial commit...")

	// Check if there's anything to commit
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = repoDir
	output, err := statusCmd.Output()
	if err != nil {
		errorMsg("Failed to check git status")
		return err
	}

	if len(output) == 0 {
		warn("Nothing to commit")
		return nil
	}

	addCmd := exec.Command("git", "add", ".ubsignore", ".gitignore", "README.md")
	addCmd.Dir = repoDir
	if cfg.Verbose {
		addCmd.Stdout = os.Stdout
		addCmd.Stderr = os.Stderr
	}
	if err := addCmd.Run(); err != nil {
		errorMsg("Failed to add files")
		return err
	}

	commitMsg := `Initial commit

- Add .ubsignore for UBS scanner
- Add .gitignore for version control
- Add README.md with project structure
- Initialize Beads issue tracking`

	commitCmd := exec.Command("git", "commit", "-m", commitMsg)
	commitCmd.Dir = repoDir
	if cfg.Verbose {
		commitCmd.Stdout = os.Stdout
		commitCmd.Stderr = os.Stderr
	}

	if err := commitCmd.Run(); err != nil {
		errorMsg("Failed to create commit")
		return err
	}

	success("Initial commit created")
	return nil
}

func setupGitRemote(repoDir string, cfg *Config) error {
	if cfg.SkipGit || cfg.SkipRemote {
		if cfg.Verbose {
			info("Skipping remote setup")
		}
		return nil
	}

	// Check if remote already exists
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	
	if err == nil {
		// Remote already exists
		existingURL := strings.TrimSpace(string(output))
		warn(fmt.Sprintf("Remote 'origin' already exists: %s", existingURL))
		return nil
	}

	fmt.Println()
	info("Git remote setup (optional)")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s?%s Add git remote? (y/N): ", blue, reset)
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil
	}
	response = strings.TrimSpace(response)

	if strings.ToLower(response) != "y" {
		info("Skipped remote setup")
		return nil
	}

	fmt.Println()
	info(fmt.Sprintf("Enter remote URL (e.g., https://github.com/username/%s.git)", cfg.ProjectName))
	fmt.Printf("%s→%s Remote URL: ", blue, reset)
	remoteURL, err := reader.ReadString('\n')
	if err != nil {
		return nil
	}
	remoteURL = strings.TrimSpace(remoteURL)

	if remoteURL == "" {
		info("Skipped remote setup")
		return nil
	}

	// Validate the git remote URL before using it
	if err := validate.ValidateGitURL(remoteURL); err != nil {
		errorMsg(fmt.Sprintf("Invalid git remote URL: %v", err))
		return fmt.Errorf("invalid git remote URL: %w", err)
	}

	addRemoteCmd := exec.Command("git", "remote", "add", "origin", remoteURL)
	addRemoteCmd.Dir = repoDir
	if err := addRemoteCmd.Run(); err != nil {
		warn("Failed to add remote (may already exist)")
		return nil
	}

	success(fmt.Sprintf("Added remote: origin → %s", remoteURL))
	fmt.Println()
	info("Push to remote with:")
	if cfg.InPlace {
		fmt.Println("  git push -u origin main")
	} else {
		fmt.Printf("  cd %s/%s\n", cfg.ProjectName, cfg.ProjectName)
		fmt.Println("  git push -u origin main")
	}
	fmt.Println()

	return nil
}

func showSummary(repoPath string, cfg *Config) {
	fmt.Println()
	fmt.Printf("%s╔═══════════════════════════════════════════════════════════╗%s\n", green, reset)
	fmt.Printf("%s║  ✓ Repository initialized successfully!                   ║%s\n", green, reset)
	fmt.Printf("%s╚═══════════════════════════════════════════════════════════╝%s\n", green, reset)
	fmt.Println()
	success(fmt.Sprintf("Project: %s", cfg.ProjectName))
	success(fmt.Sprintf("Location: %s", repoPath))
	fmt.Println()
	
	if !cfg.InPlace {
		info("Next steps:")
		fmt.Printf("  1. cd %s/%s\n", cfg.ProjectName, cfg.ProjectName)
		fmt.Println("  2. Create your project files")
	} else {
		info("Next steps:")
		fmt.Println("  1. Create your project files")
	}
	
	fmt.Println("  2. Run 'ubs .' to scan for issues")
	fmt.Println("  3. Run 'bd list' to manage tasks")
	fmt.Println()
	info("Quick commands:")
	fmt.Println("  bd create --title \"Task name\"    # Create new task")
	fmt.Println("  bd list                           # View all tasks")
	fmt.Println("  ubs .                             # Scan for bugs")
	fmt.Println()
}

func main() {
	cfg := parseFlags()

	fmt.Println()
	fmt.Printf("%s╔═══════════════════════════════════════════════════════════╗%s\n", blue, reset)
	fmt.Printf("%s║  Maajise v%s - Repository Initialization                 ║%s\n", blue, VERSION, reset)
	fmt.Printf("%s╚═══════════════════════════════════════════════════════════╝%s\n", blue, reset)
	fmt.Println()

	// Validate
	if !cfg.InPlace {
		info(fmt.Sprintf("Validating: %s", cfg.ProjectName))
		if err := validateProjectName(cfg.ProjectName); err != nil {
			os.Exit(1)
		}
		success("Valid project name")
		fmt.Println()
	}

	// Check dependencies
	info("Checking dependencies...")
	if err := checkDependencies(cfg); err != nil {
		os.Exit(1)
	}
	success("Dependencies available")
	fmt.Println()

	// Create/determine repo directory
	repoPath, err := createStructure(cfg)
	if err != nil {
		errorMsg(err.Error())
		os.Exit(1)
	}
	fmt.Println()

	// Initialize Git
	if err := initGit(repoPath, cfg); err != nil {
		os.Exit(1)
	}

	// Configure Git user
	if err := configureGitUser(repoPath, cfg); err != nil {
		os.Exit(1)
	}

	// Initialize Beads
	initBeads(repoPath, cfg)

	// Create standard files
	createUbsignore(repoPath, cfg)
	createGitignore(repoPath, cfg)
	createReadme(repoPath, cfg)

	// Initial commit
	fmt.Println()
	if err := createInitialCommit(repoPath, cfg); err != nil {
		os.Exit(1)
	}

	// Setup remote (optional)
	setupGitRemote(repoPath, cfg)

	// Show summary
	showSummary(repoPath, cfg)
}

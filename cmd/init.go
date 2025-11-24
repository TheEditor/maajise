package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"maajise/internal/beads"
	"maajise/internal/config"
	"maajise/internal/git"
	"maajise/internal/ui"
	"maajise/internal/validate"
	"maajise/templates"
)

type InitCommand struct {
	fs       *flag.FlagSet
	config   config.Config
	template string
}

func NewInitCommand() *InitCommand {
	ic := &InitCommand{
		fs:     flag.NewFlagSet("init", flag.ContinueOnError),
		config: config.DefaultConfig(),
	}

	// Define flags
	ic.fs.BoolVar(&ic.config.InPlace, "in-place", false, "Initialize in current directory")
	ic.fs.BoolVar(&ic.config.NoOverwrite, "no-overwrite", false, "Don't overwrite existing files")
	ic.fs.StringVar(&ic.template, "template", "base", "Project template (base, typescript, python, rust, php, go)")
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
		"maajise init my-project --template=typescript",
		"maajise init my-service --template=python",
		"maajise init --in-place --template=rust",
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

	// Create standard files from template
	if err := ic.createFiles(repoPath); err != nil {
		return err
	}

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
			ui.Info(fmt.Sprintf("Using current directory: %s", cwd))
		}
		return cwd, nil
	}

	// Standard mode: create nested structure
	ui.Info("Creating directory structure...")
	innerPath := filepath.Join(ic.config.ProjectName, ic.config.ProjectName)

	if ic.fileExists(ic.config.ProjectName) {
		return "", fmt.Errorf("directory '%s' already exists", ic.config.ProjectName)
	}

	if err := os.MkdirAll(innerPath, 0755); err != nil {
		return "", err
	}

	ui.Success(fmt.Sprintf("Created %s/", innerPath))
	return innerPath, nil
}

func (ic *InitCommand) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (ic *InitCommand) initGit(repoDir string) error {
	if ic.config.SkipGit {
		if ic.config.Verbose {
			ui.Info("Skipping Git (--skip-git)")
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
			ui.Info("Skipping Git user configuration")
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
			ui.Info(fmt.Sprintf("Using provided Git config: %s <%s>", userName, userEmail))
		}
	} else {
		// Interactive prompts
		ui.Info("Configuring Git user...")
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

	ui.Success(fmt.Sprintf("Git user configured: %s <%s>", userName, userEmail))
	fmt.Println()
	return nil
}

func (ic *InitCommand) initBeads(repoDir string) {
	if ic.config.SkipBeads {
		if ic.config.Verbose {
			ui.Info("Skipping Beads (--skip-beads)")
		}
		return
	}

	if err := beads.Init(repoDir, ic.config.Verbose); err != nil {
		ui.Warn("Beads init failed (run 'bd init' manually)")
		return
	}
}

func (ic *InitCommand) createFiles(repoDir string) error {
	tmpl, ok := templates.Get(ic.template)
	if !ok {
		return fmt.Errorf("unknown template: %s (use 'maajise templates' to list available)", ic.template)
	}

	files := tmpl.Files(ic.config.ProjectName)
	for filename, content := range files {
		path := filepath.Join(repoDir, filename)
		if err := ic.writeFileIfNotExists(path, content); err != nil {
			return err
		}
	}
	return nil
}

func (ic *InitCommand) writeFileIfNotExists(path, content string) error {
	if ic.fileExists(path) {
		if ic.config.NoOverwrite {
			ui.Warn(fmt.Sprintf("Skipped %s (exists, --no-overwrite)", filepath.Base(path)))
			return nil
		}
		ui.Warn(fmt.Sprintf("Overwriting %s", filepath.Base(path)))
	}

	// Create parent directories if needed
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	ui.Success(fmt.Sprintf("Created %s", path))
	return nil
}


func (ic *InitCommand) createInitialCommit(repoDir string) error {
	if ic.config.SkipGit || ic.config.SkipCommit {
		if ic.config.Verbose {
			ui.Info("Skipping initial commit")
		}
		return nil
	}

	ui.Info("Creating initial commit...")

	// Check if there's anything to commit
	hasChanges, err := git.HasChanges(repoDir)
	if err != nil {
		return err
	}

	if !hasChanges {
		ui.Warn("Nothing to commit")
		return nil
	}

	// Get files from template
	tmpl, _ := templates.Get(ic.template)
	files := tmpl.Files(ic.config.ProjectName)
	fileList := make([]string, 0, len(files))
	for filename := range files {
		fileList = append(fileList, filename)
	}

	// Add files
	if err := git.AddFiles(repoDir, fileList, ic.config.Verbose); err != nil {
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

	ui.Success("Initial commit created")
	return nil
}

func (ic *InitCommand) setupGitRemote(repoDir string) {
	if ic.config.SkipGit || ic.config.SkipRemote {
		if ic.config.Verbose {
			ui.Info("Skipping remote setup")
		}
		return
	}

	// Check if remote already exists
	existingURL, err := git.GetRemote(repoDir, "origin")
	if err == nil {
		// Remote already exists
		ui.Warn(fmt.Sprintf("Remote 'origin' already exists: %s", existingURL))
		return
	}

	fmt.Println()
	ui.Info("Git remote setup (optional)")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("? Add git remote? (y/N): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	response = strings.TrimSpace(response)

	if strings.ToLower(response) != "y" {
		ui.Info("Skipped remote setup")
		return
	}

	fmt.Println()
	ui.Info(fmt.Sprintf("Enter remote URL (e.g., https://github.com/username/%s.git)", ic.config.ProjectName))
	fmt.Print("→ Remote URL: ")
	remoteURL, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	remoteURL = strings.TrimSpace(remoteURL)

	if remoteURL == "" {
		ui.Info("Skipped remote setup")
		return
	}

	// Validate the git remote URL
	if err := validate.ValidateGitURL(remoteURL); err != nil {
		ui.Error(fmt.Sprintf("Invalid git remote URL: %v", err))
		return
	}

	if err := git.AddRemote(repoDir, "origin", remoteURL, ic.config.Verbose); err != nil {
		ui.Warn("Failed to add remote (may already exist)")
		return
	}

	ui.Success(fmt.Sprintf("Added remote: origin → %s", remoteURL))
	fmt.Println()
	ui.Info("Push to remote with:")
	if ic.config.InPlace {
		fmt.Println("  git push -u origin main")
	} else {
		fmt.Printf("  cd %s/%s\n", ic.config.ProjectName, ic.config.ProjectName)
		fmt.Println("  git push -u origin main")
	}
	fmt.Println()
}

func (ic *InitCommand) showSummary(repoPath string) {
	ui.Summary(
		"✓ Repository initialized successfully!",
		"",
		fmt.Sprintf("Project: %s", ic.config.ProjectName),
		fmt.Sprintf("Location: %s", repoPath),
	)

	if !ic.config.InPlace {
		ui.Info("Next steps:")
		fmt.Printf("  1. cd %s/%s\n", ic.config.ProjectName, ic.config.ProjectName)
		fmt.Println("  2. Create your project files")
	} else {
		ui.Info("Next steps:")
		fmt.Println("  1. Create your project files")
	}

	fmt.Println("  2. Run 'ubs .' to scan for issues")
	fmt.Println("  3. Run 'bd list' to manage tasks")
	fmt.Println()
	ui.Info("Quick commands:")
	fmt.Println("  bd create --title \"Task name\"    # Create new task")
	fmt.Println("  bd list                           # View all tasks")
	fmt.Println("  ubs .                             # Scan for bugs")
	fmt.Println()
}

func init() {
	Register(NewInitCommand())
}

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
	"maajise/internal/fsutil"
	"maajise/internal/git"
	"maajise/internal/ui"
	"maajise/internal/validate"
	"maajise/templates"
)

type InitCommand struct {
	fs         *flag.FlagSet
	config     config.Config
	template   string
	fileConfig *config.FileConfig
	dryRun     bool
	interactive bool
}

func NewInitCommand() *InitCommand {
	ic := &InitCommand{
		fs:     flag.NewFlagSet("init", flag.ContinueOnError),
		config: config.DefaultConfig(),
	}

	// Load file config and merge defaults BEFORE defining flags
	// This way, flag defaults can reflect config file values
	if fc, err := config.LoadFileConfig(); err == nil && fc != nil {
		ic.config.MergeFileConfig(fc)
		ic.fileConfig = fc // Store for template variables later
	}

	// Define flags (will use merged defaults)
	ic.fs.BoolVar(&ic.config.InPlace, "in-place", false, "Initialize in current directory")
	ic.fs.BoolVar(&ic.config.NoOverwrite, "no-overwrite", false, "Don't overwrite existing files")
	ic.fs.StringVar(&ic.template, "template", ic.config.Template, "Project template (base, typescript, python, rust, php, go)")
	ic.fs.BoolVar(&ic.config.SkipGit, "skip-git", false, "Skip Git initialization")
	ic.fs.BoolVar(&ic.config.SkipBeads, "skip-beads", false, "Skip Beads initialization")
	ic.fs.BoolVar(&ic.config.SkipCommit, "skip-commit", false, "Skip initial commit")
	ic.fs.BoolVar(&ic.config.SkipRemote, "skip-remote", false, "Skip remote setup")
	ic.fs.BoolVar(&ic.config.SkipGitUser, "skip-git-user", false, "Skip Git user configuration")
	ic.fs.StringVar(&ic.config.GitName, "git-name", "", "Git user.name (non-interactive)")
	ic.fs.StringVar(&ic.config.GitEmail, "git-email", "", "Git user.email (non-interactive)")
	ic.fs.BoolVar(&ic.dryRun, "dry-run", false, "Preview what would be created without making changes")
	ic.fs.BoolVar(&ic.interactive, "interactive", false, "Interactive mode with prompts")
	ic.fs.BoolVar(&ic.interactive, "i", false, "Interactive mode with prompts")
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
		"maajise init my-project --dry-run",
		"maajise init --interactive",
		"maajise init --in-place --template=rust",
	}
}

func (ic *InitCommand) Execute(args []string) error {
	// Parse flags
	if err := ic.fs.Parse(args); err != nil {
		return err
	}

	remainingArgs := ic.fs.Args()

	// Auto-enter interactive mode if no project name and not in-place
	if !ic.config.InPlace && len(remainingArgs) == 0 && !ic.interactive {
		// Offer interactive mode
		fmt.Println("No project name specified.")
		fmt.Print("Enter interactive mode? (Y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response == "" || response == "y" || response == "yes" {
			ic.interactive = true
		} else {
			return fmt.Errorf("project name required (or use --interactive)")
		}
	}

	// Handle interactive mode
	if ic.interactive {
		return ic.runInteractive()
	}

	// Get project name from remaining args
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

	// Template from flag takes precedence
	if ic.template == "" {
		ic.template = "base" // Ultimate fallback
	}

	// Execute initialization or dry-run
	if ic.dryRun {
		return ic.runDryRun()
	}
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

func (ic *InitCommand) runDryRun() error {
	ui.Info("[dry-run] Preview of initialization:")
	fmt.Println()

	// Show directory structure
	var targetDir string
	if ic.config.InPlace {
		cwd, _ := os.Getwd()
		targetDir = cwd
		ui.Info(fmt.Sprintf("[dry-run] Would initialize in: %s", targetDir))
	} else {
		targetDir = filepath.Join(ic.config.ProjectName, ic.config.ProjectName)
		ui.Info(fmt.Sprintf("[dry-run] Would create directory: %s", targetDir))
	}
	fmt.Println()

	// Show Git initialization
	if !ic.config.SkipGit {
		ui.Info("[dry-run] Would initialize Git repository")
		if ic.config.GitName != "" && ic.config.GitEmail != "" {
			fmt.Printf("         user.name: %s\n", ic.config.GitName)
			fmt.Printf("         user.email: %s\n", ic.config.GitEmail)
		} else {
			fmt.Println("         (would prompt for user.name and user.email)")
		}
	}

	// Show Beads initialization
	if !ic.config.SkipBeads {
		ui.Info("[dry-run] Would initialize Beads issue tracking")
	}

	// Show files from template
	tmpl, ok := templates.Get(ic.template)
	if !ok {
		return fmt.Errorf("unknown template: %s", ic.template)
	}

	files := tmpl.Files(ic.config.ProjectName)

	fmt.Println()
	ui.Info(fmt.Sprintf("[dry-run] Would create files (template: %s):", ic.template))
	for filename := range files {
		fullPath := filepath.Join(targetDir, filename)
		if fsutil.FileExists(fullPath) {
			if ic.config.NoOverwrite {
				fmt.Printf("         skip: %s (exists, --no-overwrite)\n", filename)
			} else {
				fmt.Printf("         overwrite: %s\n", filename)
			}
		} else {
			fmt.Printf("         create: %s\n", filename)
		}
	}

	// Show commit
	if !ic.config.SkipGit && !ic.config.SkipCommit {
		fmt.Println()
		ui.Info("[dry-run] Would create initial commit")
	}

	fmt.Println()
	ui.Info("[dry-run] No changes made. Remove --dry-run to execute.")

	return nil
}

func (ic *InitCommand) runInteractive() error {
	reader := bufio.NewReader(os.Stdin)

	ui.Info("Maajise Interactive Setup")
	fmt.Println()

	// Project name
	fmt.Print("Project name: ")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("project name required")
	}
	if err := ic.validateProjectName(projectName); err != nil {
		return err
	}
	ic.config.ProjectName = projectName

	// Template selection
	fmt.Println()
	fmt.Println("Available templates:")
	allTemplates := templates.All()
	for i, tmpl := range allTemplates {
		fmt.Printf("  %d. %s - %s\n", i+1, tmpl.Name(), tmpl.Description())
	}
	fmt.Printf("Select template [1-%d] (default: base): ", len(allTemplates))
	templateChoice, _ := reader.ReadString('\n')
	templateChoice = strings.TrimSpace(templateChoice)
	if templateChoice == "" {
		ic.template = "base"
	} else {
		// Parse number or name
		var found bool
		for i, tmpl := range allTemplates {
			if templateChoice == fmt.Sprintf("%d", i+1) || templateChoice == tmpl.Name() {
				ic.template = tmpl.Name()
				found = true
				break
			}
		}
		if !found {
			ic.template = "base"
			ui.Warn(fmt.Sprintf("Unknown template %q, using base", templateChoice))
		}
	}

	// Git configuration
	fmt.Println()
	fmt.Print("Initialize Git? (Y/n): ")
	gitChoice, _ := reader.ReadString('\n')
	gitChoice = strings.TrimSpace(strings.ToLower(gitChoice))
	ic.config.SkipGit = (gitChoice == "n" || gitChoice == "no")

	if !ic.config.SkipGit {
		// Git name
		defaultName := ""
		if ic.fileConfig != nil && ic.fileConfig.Defaults.GitName != "" {
			defaultName = ic.fileConfig.Defaults.GitName
		}
		if defaultName != "" {
			fmt.Printf("Git user.name [%s]: ", defaultName)
		} else {
			fmt.Print("Git user.name: ")
		}
		gitName, _ := reader.ReadString('\n')
		gitName = strings.TrimSpace(gitName)
		if gitName == "" {
			gitName = defaultName
		}
		ic.config.GitName = gitName

		// Git email
		defaultEmail := ""
		if ic.fileConfig != nil && ic.fileConfig.Defaults.GitEmail != "" {
			defaultEmail = ic.fileConfig.Defaults.GitEmail
		}
		if defaultEmail != "" {
			fmt.Printf("Git user.email [%s]: ", defaultEmail)
		} else {
			fmt.Print("Git user.email: ")
		}
		gitEmail, _ := reader.ReadString('\n')
		gitEmail = strings.TrimSpace(gitEmail)
		if gitEmail == "" {
			gitEmail = defaultEmail
		}
		ic.config.GitEmail = gitEmail
	}

	// Beads
	fmt.Println()
	fmt.Print("Initialize Beads issue tracking? (Y/n): ")
	beadsChoice, _ := reader.ReadString('\n')
	beadsChoice = strings.TrimSpace(strings.ToLower(beadsChoice))
	ic.config.SkipBeads = (beadsChoice == "n" || beadsChoice == "no")

	// Summary
	fmt.Println()
	ui.Info("Summary:")
	fmt.Printf("  Project: %s\n", ic.config.ProjectName)
	fmt.Printf("  Template: %s\n", ic.template)
	fmt.Printf("  Git: %v\n", !ic.config.SkipGit)
	fmt.Printf("  Beads: %v\n", !ic.config.SkipBeads)
	fmt.Println()

	fmt.Print("Proceed? (Y/n): ")
	proceed, _ := reader.ReadString('\n')
	proceed = strings.TrimSpace(strings.ToLower(proceed))
	if proceed == "n" || proceed == "no" {
		ui.Info("Cancelled")
		return nil
	}

	// Execute
	fmt.Println()
	return ic.runInit()
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

	if fsutil.PathExists(ic.config.ProjectName) {
		return "", fmt.Errorf("directory '%s' already exists", ic.config.ProjectName)
	}

	if err := os.MkdirAll(innerPath, 0755); err != nil {
		return "", err
	}

	ui.Success(fmt.Sprintf("Created %s/", innerPath))
	return innerPath, nil
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
	if fsutil.FileExists(path) {
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

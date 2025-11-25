package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"maajise/internal/beads"
	"maajise/internal/git"
	"maajise/internal/ui"
	"maajise/templates"
)

type AddCommand struct {
	fs       *flag.FlagSet
	force    bool
	template string
	verbose  bool
	dryRun   bool
}

// Tooling items that can be added
var toolingItems = map[string]string{
	"git":   "Initialize Git repository",
	"beads": "Initialize Beads issue tracking",
	"ubs":   "Add .ubsignore file",
}

func NewAddCommand() *AddCommand {
	ac := &AddCommand{
		fs: flag.NewFlagSet("add", flag.ContinueOnError),
	}

	ac.fs.BoolVar(&ac.force, "force", false, "Overwrite existing files")
	ac.fs.StringVar(&ac.template, "template", "", "Template for file content (auto-detects if not specified)")
	ac.fs.BoolVar(&ac.verbose, "v", false, "Verbose output")
	ac.fs.BoolVar(&ac.verbose, "verbose", false, "Verbose output")
	ac.fs.BoolVar(&ac.dryRun, "dry-run", false, "Preview without making changes")

	return ac
}

func (ac *AddCommand) Name() string {
	return "add"
}

func (ac *AddCommand) Description() string {
	return "Add files or tooling to an existing project"
}

func (ac *AddCommand) Usage() string {
	return "maajise add <item> [flags]"
}

func (ac *AddCommand) Examples() []string {
	return []string{
		"maajise add git",
		"maajise add beads",
		"maajise add ubs",
		"maajise add .gitignore",
		"maajise add .gitignore --template=typescript",
		"maajise add readme",
		"maajise add --dry-run git",
	}
}

func (ac *AddCommand) Execute(args []string) error {
	if err := ac.fs.Parse(args); err != nil {
		return err
	}

	items := ac.fs.Args()
	if len(items) == 0 {
		return ac.showHelp()
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	projectName := filepath.Base(cwd)

	// Detect template if not specified
	if ac.template == "" {
		ac.template = ac.detectTemplate(cwd)
		if ac.verbose {
			ui.Info(fmt.Sprintf("Detected template: %s", ac.template))
		}
	}

	for _, item := range items {
		if err := ac.addItem(cwd, projectName, item); err != nil {
			ui.Error(fmt.Sprintf("Failed to add %s: %v", item, err))
			// Continue with other items
		}
	}

	return nil
}

func (ac *AddCommand) showHelp() error {
	fmt.Println("Usage: maajise add <item> [flags]")
	fmt.Println()
	fmt.Println("Tooling:")
	for name, desc := range toolingItems {
		fmt.Printf("  %-12s  %s\n", name, desc)
	}
	fmt.Println()
	fmt.Println("Files:")
	fmt.Println("  .gitignore   Add .gitignore from template")
	fmt.Println("  .ubsignore   Add .ubsignore from template")
	fmt.Println("  readme       Add README.md from template")
	fmt.Println()
	return nil
}

func (ac *AddCommand) addItem(dir, projectName, item string) error {
	// Normalize item name
	itemLower := strings.ToLower(item)

	// Check if it's a tooling item
	switch itemLower {
	case "git":
		return ac.addGit(dir)
	case "beads":
		return ac.addBeads(dir)
	case "ubs", ".ubsignore", "ubsignore":
		return ac.addFile(dir, projectName, ".ubsignore")
	case ".gitignore", "gitignore":
		return ac.addFile(dir, projectName, ".gitignore")
	case "readme", "readme.md":
		return ac.addFile(dir, projectName, "README.md")
	default:
		// Try as template file
		return ac.addFile(dir, projectName, item)
	}
}

func (ac *AddCommand) addGit(dir string) error {
	gitDir := filepath.Join(dir, ".git")
	if ac.dirExists(gitDir) {
		ui.Warn("Git already initialized")
		return nil
	}

	if ac.dryRun {
		ui.Info("[dry-run] Would initialize Git repository")
		return nil
	}

	if err := git.Init(dir, ac.verbose); err != nil {
		return err
	}

	ui.Success("Initialized Git repository")
	return nil
}

func (ac *AddCommand) addBeads(dir string) error {
	beadsDir := filepath.Join(dir, ".beads")
	if ac.dirExists(beadsDir) {
		ui.Warn("Beads already initialized")
		return nil
	}

	if ac.dryRun {
		ui.Info("[dry-run] Would initialize Beads issue tracking")
		return nil
	}

	if err := beads.Init(dir, ac.verbose); err != nil {
		return err
	}

	ui.Success("Initialized Beads issue tracking")
	return nil
}

func (ac *AddCommand) addFile(dir, projectName, filename string) error {
	tmpl, ok := templates.Get(ac.template)
	if !ok {
		return fmt.Errorf("unknown template: %s", ac.template)
	}

	files := tmpl.Files(projectName)
	content, ok := files[filename]
	if !ok {
		return fmt.Errorf("file %s not found in template %s", filename, ac.template)
	}

	path := filepath.Join(dir, filename)

	if ac.fileExists(path) && !ac.force {
		ui.Warn(fmt.Sprintf("Skipped %s (exists, use --force to overwrite)", filename))
		return nil
	}

	if ac.dryRun {
		if ac.fileExists(path) {
			ui.Info(fmt.Sprintf("[dry-run] Would overwrite: %s", filename))
		} else {
			ui.Info(fmt.Sprintf("[dry-run] Would create: %s", filename))
		}
		return nil
	}

	// Create parent directories if needed
	parentDir := filepath.Dir(path)
	if parentDir != "." && parentDir != "" {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	if ac.fileExists(path) && ac.force {
		ui.Success(fmt.Sprintf("Updated %s", filename))
	} else {
		ui.Success(fmt.Sprintf("Created %s", filename))
	}

	return nil
}

func (ac *AddCommand) fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func (ac *AddCommand) dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func (ac *AddCommand) detectTemplate(dir string) string {
	checks := []struct {
		file     string
		template string
	}{
		{"package.json", "typescript"},
		{"Cargo.toml", "rust"},
		{"pyproject.toml", "python"},
		{"composer.json", "php"},
		{"go.mod", "go"},
	}

	for _, check := range checks {
		if ac.fileExists(filepath.Join(dir, check.file)) {
			return check.template
		}
	}

	return "base"
}

func init() {
	Register(NewAddCommand())
}

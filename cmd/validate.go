package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"maajise/internal/detect"
	"maajise/internal/fsutil"
	"maajise/internal/ui"
	"maajise/templates"
)

type ValidateCommand struct {
	fs      *flag.FlagSet
	verbose bool
	strict  bool
}

type ValidationResult struct {
	Check   string
	Status  string // "pass", "warn", "fail"
	Message string
}

func NewValidateCommand() *ValidateCommand {
	vc := &ValidateCommand{
		fs: flag.NewFlagSet("validate", flag.ContinueOnError),
	}

	vc.fs.BoolVar(&vc.verbose, "v", false, "Verbose output")
	vc.fs.BoolVar(&vc.verbose, "verbose", false, "Verbose output")
	vc.fs.BoolVar(&vc.strict, "strict", false, "Treat warnings as failures")

	return vc
}

func (vc *ValidateCommand) Name() string {
	return "validate"
}

func (vc *ValidateCommand) Description() string {
	return "Validate project setup and configuration"
}

func (vc *ValidateCommand) LongDescription() string {
	return `Validate the current project's setup and configuration.

Checks for Git initialization, Beads setup, required configuration files, and
template-specific requirements. Reports any issues and warnings found.`
}

func (vc *ValidateCommand) Usage() string {
	return "maajise validate [flags]"
}

func (vc *ValidateCommand) Examples() string {
	return `  # Validate current project
  maajise validate

  # Strict mode - treat warnings as failures
  maajise validate --strict

  # Verbose output with detailed information
  maajise validate --verbose`
}

func (vc *ValidateCommand) Run(args []string) error {
	if err := vc.fs.Parse(args); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	projectName := filepath.Base(cwd)
	ui.Info(fmt.Sprintf("Validating project: %s", projectName))
	fmt.Println()

	results := []ValidationResult{}

	// Check Git initialization
	results = append(results, vc.checkGit(cwd))

	// Check Beads initialization
	results = append(results, vc.checkBeads(cwd))

	// Check required files
	results = append(results, vc.checkRequiredFiles(cwd)...)

	// Detect and validate template-specific files
	template := detect.Template(cwd)
	if vc.verbose {
		ui.Info(fmt.Sprintf("Detected template: %s", template))
	}
	results = append(results, vc.checkTemplateFiles(cwd, template)...)

	// Display results
	passes := 0
	warns := 0
	fails := 0

	for _, r := range results {
		switch r.Status {
		case "pass":
			ui.Success(fmt.Sprintf("✓ %s: %s", r.Check, r.Message))
			passes++
		case "warn":
			ui.Warn(fmt.Sprintf("⚠ %s: %s", r.Check, r.Message))
			warns++
		case "fail":
			ui.Error(fmt.Sprintf("✗ %s: %s", r.Check, r.Message))
			fails++
		}
	}

	fmt.Println()
	ui.Info(fmt.Sprintf("Results: %d passed, %d warnings, %d failed", passes, warns, fails))

	// Return error if failures (or warnings in strict mode)
	if fails > 0 {
		return fmt.Errorf("validation failed with %d errors", fails)
	}
	if vc.strict && warns > 0 {
		return fmt.Errorf("validation failed with %d warnings (strict mode)", warns)
	}

	return nil
}

func (vc *ValidateCommand) checkGit(dir string) ValidationResult {
	gitDir := filepath.Join(dir, ".git")
	if fsutil.DirExists(gitDir) {
		return ValidationResult{"Git", "pass", "Repository initialized"}
	}
	return ValidationResult{"Git", "fail", "Not a git repository (run 'git init')"}
}

func (vc *ValidateCommand) checkBeads(dir string) ValidationResult {
	beadsDir := filepath.Join(dir, ".beads")
	if fsutil.DirExists(beadsDir) {
		return ValidationResult{"Beads", "pass", "Issue tracking initialized"}
	}
	return ValidationResult{"Beads", "warn", "Not initialized (run 'bd init')"}
}

func (vc *ValidateCommand) checkRequiredFiles(dir string) []ValidationResult {
	results := []ValidationResult{}

	required := []string{".gitignore", "README.md"}
	for _, f := range required {
		if fsutil.FileExists(filepath.Join(dir, f)) {
			results = append(results, ValidationResult{f, "pass", "Present"})
		} else {
			results = append(results, ValidationResult{f, "warn", "Missing"})
		}
	}

	// .ubsignore is recommended but not required
	if fsutil.FileExists(filepath.Join(dir, ".ubsignore")) {
		results = append(results, ValidationResult{".ubsignore", "pass", "Present"})
	} else if vc.verbose {
		results = append(results, ValidationResult{".ubsignore", "warn", "Missing (recommended)"})
	}

	return results
}

func (vc *ValidateCommand) checkTemplateFiles(dir, template string) []ValidationResult {
	results := []ValidationResult{}

	tmpl, ok := templates.Get(template)
	if !ok {
		return results
	}

	// Get expected files from template
	files := tmpl.Files(filepath.Base(dir))

	for filename := range files {
		// Skip common files already checked
		if filename == ".gitignore" || filename == "README.md" || filename == ".ubsignore" {
			continue
		}

		if fsutil.FileExists(filepath.Join(dir, filename)) {
			if vc.verbose {
				results = append(results, ValidationResult{filename, "pass", "Present"})
			}
		} else {
			results = append(results, ValidationResult{filename, "warn", fmt.Sprintf("Missing (expected for %s template)", template)})
		}
	}

	return results
}

func init() {
	Register(NewValidateCommand())
}

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

type UpdateCommand struct {
	fs        *flag.FlagSet
	force     bool
	template  string
	verbose   bool
	dryRun    bool
	filesOnly []string
}

func NewUpdateCommand() *UpdateCommand {
	uc := &UpdateCommand{
		fs: flag.NewFlagSet("update", flag.ContinueOnError),
	}

	uc.fs.BoolVar(&uc.force, "force", false, "Overwrite existing files (default: skip existing)")
	uc.fs.StringVar(&uc.template, "template", "", "Template to use (auto-detects if not specified)")
	uc.fs.BoolVar(&uc.verbose, "v", false, "Verbose output")
	uc.fs.BoolVar(&uc.verbose, "verbose", false, "Verbose output")
	uc.fs.BoolVar(&uc.dryRun, "dry-run", false, "Show what would be updated without making changes")

	return uc
}

func (uc *UpdateCommand) Name() string {
	return "update"
}

func (uc *UpdateCommand) Description() string {
	return "Update configuration files in an existing project"
}

func (uc *UpdateCommand) LongDescription() string {
	return `Update configuration files in an existing project.

Updates files like .gitignore, .ubsignore, and README.md based on the project's template.
If no specific files are provided, all standard configuration files are updated. Use --dry-run
to preview changes before applying them.`
}

func (uc *UpdateCommand) Usage() string {
	return "maajise update [flags] [files...]"
}

func (uc *UpdateCommand) Examples() string {
	return `  # Update all configuration files
  maajise update

  # Force overwrite existing files
  maajise update --force

  # Update with specific template
  maajise update --template=typescript

  # Update specific files only
  maajise update .gitignore .ubsignore

  # Preview changes without applying
  maajise update --dry-run

  # Verbose output
  maajise update --verbose
      Shows detailed information about each file updated`
}

func (uc *UpdateCommand) Run(args []string) error {
	if err := uc.fs.Parse(args); err != nil {
		return err
	}

	// Get specific files to update (if any)
	uc.filesOnly = uc.fs.Args()

	// Determine working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get project name from directory
	projectName := filepath.Base(cwd)

	// Detect or use specified template
	templateName := uc.template
	if templateName == "" {
		templateName = detect.Template(cwd)
		if uc.verbose {
			ui.Info(fmt.Sprintf("Detected template: %s", templateName))
		}
	}

	tmpl, ok := templates.Get(templateName)
	if !ok {
		return ui.UsageError("update", fmt.Sprintf("unknown template: %s (available: base, typescript, python, rust, php, go)", templateName))
	}

	// Get files from template
	files := tmpl.Files(projectName)

	// Filter to specific files if requested
	if len(uc.filesOnly) > 0 {
		filtered := make(map[string]string)
		for _, f := range uc.filesOnly {
			if content, ok := files[f]; ok {
				filtered[f] = content
			} else {
				ui.Warn(fmt.Sprintf("File %s not in template %s", f, templateName))
			}
		}
		files = filtered
	}

	if len(files) == 0 {
		return ui.UsageError("update", "no files to update")
	}

	// Update files
	updated := 0
	skipped := 0
	for filename, content := range files {
		path := filepath.Join(cwd, filename)

		// Check if file exists BEFORE writing
		existed := fsutil.PathExists(path)

		if uc.dryRun {
			if existed {
				if uc.force {
					ui.Info(fmt.Sprintf("[dry-run] Would overwrite: %s", filename))
				} else {
					ui.Info(fmt.Sprintf("[dry-run] Would skip (exists): %s", filename))
				}
			} else {
				ui.Info(fmt.Sprintf("[dry-run] Would create: %s", filename))
			}
			continue
		}

		if existed && !uc.force {
			if uc.verbose {
				ui.Warn(fmt.Sprintf("Skipped %s (exists, use --force to overwrite)", filename))
			}
			skipped++
			continue
		}

		// Create parent directories if needed
		dir := filepath.Dir(path)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}

		// Report based on whether file existed
		if existed && uc.force {
			ui.Success(fmt.Sprintf("Updated %s", filename))
		} else {
			ui.Success(fmt.Sprintf("Created %s", filename))
		}
		updated++
	}

	if !uc.dryRun {
		fmt.Println()
		ui.Info(fmt.Sprintf("Updated: %d, Skipped: %d", updated, skipped))
	}

	return nil
}

func init() {
	Register(NewUpdateCommand())
}

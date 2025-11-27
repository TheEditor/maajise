package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"maajise/internal/detect"
	"maajise/internal/fsutil"
	"maajise/internal/ui"
)

type StatusCommand struct {
	fs *flag.FlagSet
}

func NewStatusCommand() *StatusCommand {
	return &StatusCommand{
		fs: flag.NewFlagSet("status", flag.ContinueOnError),
	}
}

func (sc *StatusCommand) Name() string {
	return "status"
}

func (sc *StatusCommand) Description() string {
	return "Show quick project status"
}

func (sc *StatusCommand) Usage() string {
	return "maajise status"
}

func (sc *StatusCommand) Examples() []string {
	return []string{
		"maajise status",
	}
}

func (sc *StatusCommand) Execute(args []string) error {
	if err := sc.fs.Parse(args); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	projectName := filepath.Base(cwd)

	fmt.Printf("Project: %s\n", projectName)
	fmt.Printf("Path:    %s\n", cwd)
	fmt.Println()

	// Git status
	gitDir := filepath.Join(cwd, ".git")
	if fsutil.DirExists(gitDir) {
		ui.Success("Git:     initialized")
	} else {
		ui.Warn("Git:     not initialized")
	}

	// Beads status
	beadsDir := filepath.Join(cwd, ".beads")
	if fsutil.DirExists(beadsDir) {
		ui.Success("Beads:   initialized")
	} else {
		ui.Warn("Beads:   not initialized")
	}

	// Template detection
	template := detect.Template(cwd)
	fmt.Printf("Template: %s\n", template)

	// Key files
	fmt.Println()
	fmt.Println("Files:")
	keyFiles := []string{".gitignore", ".ubsignore", "README.md"}
	for _, f := range keyFiles {
		if fsutil.FileExists(filepath.Join(cwd, f)) {
			fmt.Printf("  ✓ %s\n", f)
		} else {
			fmt.Printf("  ✗ %s\n", f)
		}
	}

	return nil
}

func init() {
	Register(NewStatusCommand())
}

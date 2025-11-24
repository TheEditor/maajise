package beads

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"maajise/internal/ui"
)

// Init initializes Beads issue tracking
func Init(repoDir string, verbose bool) error {
	// Check if already initialized
	beadsDir := filepath.Join(repoDir, ".beads")
	if _, err := os.Stat(beadsDir); err == nil {
		ui.Warn("Beads already initialized")
		return nil
	}

	ui.Info("Initializing Beads...")

	cmd := exec.Command("bd", "init")
	cmd.Dir = repoDir
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		ui.Warn("Beads init failed (run 'bd init' manually)")
		return nil // Don't fail whole operation
	}

	ui.Success("Beads initialized")
	return nil
}

// CheckAvailable checks if bd (beads) is available
func CheckAvailable() error {
	if _, err := exec.LookPath("bd"); err != nil {
		return fmt.Errorf("beads (bd) not found")
	}
	return nil
}

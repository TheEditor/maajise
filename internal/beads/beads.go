package beads

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"maajise/internal/ui"
)

// Init initializes Beads issue tracking via beads_rust (br)
func Init(repoDir string, verbose bool) error {
	// Check if already initialized
	beadsDir := filepath.Join(repoDir, ".beads")
	if _, err := os.Stat(beadsDir); err == nil {
		ui.Warn("Beads already initialized")
		return nil
	}

	ui.Info("Initializing Beads (beads_rust)...")

	cmd := exec.Command("br", "init")
	cmd.Dir = repoDir
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		ui.Warn("Beads init failed (run 'br init' manually)")
		return nil // Don't fail whole operation
	}

	ui.Success("Beads initialized")
	return nil
}

// CheckAvailable checks if br (beads_rust) is available
func CheckAvailable() error {
	if _, err := exec.LookPath("br"); err != nil {
		return fmt.Errorf("beads_rust (br) not found")
	}
	return nil
}

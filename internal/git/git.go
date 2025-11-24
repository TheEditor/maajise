package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"maajise/internal/ui"
)

// Init initializes a Git repository
func Init(repoDir string, verbose bool) error {
	// Check if already a git repo
	gitDir := filepath.Join(repoDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		ui.Warn("Git already initialized")
		return nil
	}

	ui.Info("Initializing Git...")

	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		ui.Error("Git init failed")
		return err
	}

	ui.Success("Git initialized")
	return nil
}

// CreateCommit creates a commit with the given files and message
func CreateCommit(repoDir string, files []string, message string, verbose bool) error {
	ui.Info("Creating initial commit...")

	// Check if there's anything to commit
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = repoDir
	output, err := statusCmd.Output()
	if err != nil {
		ui.Error("Failed to check git status")
		return err
	}

	if len(output) == 0 {
		ui.Warn("Nothing to commit")
		return nil
	}

	// Add files
	args := append([]string{"add"}, files...)
	addCmd := exec.Command("git", args...)
	addCmd.Dir = repoDir
	if verbose {
		addCmd.Stdout = os.Stdout
		addCmd.Stderr = os.Stderr
	}
	if err := addCmd.Run(); err != nil {
		ui.Error("Failed to add files")
		return err
	}

	// Commit
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = repoDir
	if verbose {
		commitCmd.Stdout = os.Stdout
		commitCmd.Stderr = os.Stderr
	}

	if err := commitCmd.Run(); err != nil {
		ui.Error("Failed to create commit")
		return err
	}

	ui.Success("Initial commit created")
	return nil
}

// CheckAvailable checks if git is available
func CheckAvailable() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found")
	}
	return nil
}

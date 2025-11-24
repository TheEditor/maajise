package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// SetConfig sets a git configuration value in a specific repository
func SetConfig(repoDir string, key string, value string, verbose bool) error {
	cmd := exec.Command("git", "config", key, value)
	cmd.Dir = repoDir
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set git config %s: %w", key, err)
	}

	return nil
}

// GetConfig retrieves a git configuration value from a specific repository
func GetConfig(repoDir string, key string) (string, error) {
	cmd := exec.Command("git", "config", key)
	cmd.Dir = repoDir

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git config %s: %w", key, err)
	}

	return strings.TrimSpace(string(output)), nil
}

// HasChanges checks if there are any uncommitted changes
func HasChanges(repoDir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoDir

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	return len(output) > 0, nil
}

// AddFiles stages files for commit
func AddFiles(repoDir string, files []string, verbose bool) error {
	if len(files) == 0 {
		return fmt.Errorf("no files to add")
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	return nil
}

// CreateCommit creates a commit with the given message
func CreateCommit(repoDir string, message string, verbose bool) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoDir
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}

// GetRemote retrieves the URL for a named remote
func GetRemote(repoDir string, remoteName string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", remoteName)
	cmd.Dir = repoDir

	output, err := cmd.Output()
	if err != nil {
		// Remote doesn't exist
		return "", fmt.Errorf("remote %s not found", remoteName)
	}

	return strings.TrimSpace(string(output)), nil
}

// AddRemote adds a new remote to the repository
func AddRemote(repoDir string, remoteName string, url string, verbose bool) error {
	cmd := exec.Command("git", "remote", "add", remoteName, url)
	cmd.Dir = repoDir
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add remote %s: %w", remoteName, err)
	}

	return nil
}

// CheckAvailable checks if git is available
func CheckAvailable() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found")
	}
	return nil
}

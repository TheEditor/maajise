package beads

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Helper to check if bd is available
func bdAvailable(t *testing.T) {
	if err := CheckAvailable(); err != nil {
		t.Skip("bd not installed, skipping test")
	}
}

// Helper to create a temporary directory
func createTempDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "beads-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	return tmpDir
}

// Helper to clean up a directory
func cleanup(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("failed to clean up temp directory: %v", err)
	}
}

func TestCheckAvailable(t *testing.T) {
	err := CheckAvailable()
	if err != nil {
		t.Skipf("bd not available: %v", err)
	}
}

func TestInit(t *testing.T) {
	bdAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize git first (required for bd init)
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	err := Init(tmpDir, false)
	if err != nil {
		t.Errorf("Init() failed: %v", err)
	}

	// Verify .beads directory exists
	beadsDir := filepath.Join(tmpDir, ".beads")
	if _, err := os.Stat(beadsDir); err != nil {
		t.Errorf("beads directory not created: %v", err)
	}
}

func TestInit_AlreadyInitialized(t *testing.T) {
	bdAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize git first
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// Initialize beads first time
	if err := Init(tmpDir, false); err != nil {
		t.Fatalf("Initial Init() failed: %v", err)
	}

	// Initialize again - should not error
	err := Init(tmpDir, false)
	if err != nil {
		t.Errorf("Init() on existing beads should not error: %v", err)
	}
}

func TestInit_Verbose(t *testing.T) {
	bdAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize git first
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// Init with verbose output
	err := Init(tmpDir, true)
	if err != nil {
		t.Errorf("Init() with verbose failed: %v", err)
	}

	// Verify .beads directory exists
	beadsDir := filepath.Join(tmpDir, ".beads")
	if _, err := os.Stat(beadsDir); err != nil {
		t.Errorf("beads directory not created: %v", err)
	}
}

func TestIntegration_FullBeadsWorkflow(t *testing.T) {
	bdAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// 1. Initialize git
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = tmpDir
	if err := gitCmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// 2. Check if bd is available
	if err := CheckAvailable(); err != nil {
		t.Fatalf("CheckAvailable() failed: %v", err)
	}

	// 3. Initialize Beads
	if err := Init(tmpDir, false); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 4. Verify .beads directory exists
	beadsDir := filepath.Join(tmpDir, ".beads")
	if _, err := os.Stat(beadsDir); err != nil {
		t.Errorf("beads directory not created: %v", err)
	}

	// 5. Verify we can check for initialization
	// (skipping 'bd list' as it may have file locks on Windows)
	if _, err := os.Stat(filepath.Join(beadsDir, "beads.jsonl")); err != nil {
		// It's ok if the file doesn't exist yet, the directory is what matters
		if !os.IsNotExist(err) {
			t.Errorf("failed to check beads.jsonl: %v", err)
		}
	}
}

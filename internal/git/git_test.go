package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Helper to check if git is available
func gitAvailable(t *testing.T) {
	if err := CheckAvailable(); err != nil {
		t.Skip("git not installed, skipping test")
	}
}

// Helper to create a temporary directory
func createTempDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "git-test-*")
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

// Helper to initialize a git repo with git command directly
func initGitRepo(t *testing.T, repoDir string) {
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to initialize git repo: %v", err)
	}

	// Configure user for commits
	_ = exec.Command("git", "config", "user.name", "Test User").Run()
	_ = exec.Command("git", "config", "user.email", "test@example.com").Run()
}

func TestCheckAvailable(t *testing.T) {
	err := CheckAvailable()
	if err != nil {
		t.Skipf("git not available: %v", err)
	}
}

func TestInit(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	err := Init(tmpDir, false)
	if err != nil {
		t.Errorf("Init() failed: %v", err)
	}

	// Verify .git directory exists
	gitDir := filepath.Join(tmpDir, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		t.Errorf("git directory not created: %v", err)
	}
}

func TestInit_AlreadyInitialized(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize first time
	if err := Init(tmpDir, false); err != nil {
		t.Fatalf("Initial Init() failed: %v", err)
	}

	// Initialize again - should not error
	err := Init(tmpDir, false)
	if err != nil {
		t.Errorf("Init() on existing repo should not error: %v", err)
	}
}

func TestSetConfig(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo first
	if err := Init(tmpDir, false); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Set config
	err := SetConfig(tmpDir, "user.name", "Test User", false)
	if err != nil {
		t.Errorf("SetConfig() failed: %v", err)
	}

	// Verify it was set
	result, err := GetConfig(tmpDir, "user.name")
	if err != nil {
		t.Errorf("GetConfig() failed: %v", err)
	}
	if result != "Test User" {
		t.Errorf("Expected 'Test User', got '%s'", result)
	}
}

func TestGetConfig(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo first
	if err := Init(tmpDir, false); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Set and retrieve
	if err := SetConfig(tmpDir, "test.key", "test.value", false); err != nil {
		t.Fatalf("SetConfig() failed: %v", err)
	}

	result, err := GetConfig(tmpDir, "test.key")
	if err != nil {
		t.Errorf("GetConfig() failed: %v", err)
	}
	if result != "test.value" {
		t.Errorf("Expected 'test.value', got '%s'", result)
	}
}

func TestGetConfig_NotFound(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo first
	if err := Init(tmpDir, false); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Try to get non-existent key
	_, err := GetConfig(tmpDir, "nonexistent.key")
	if err == nil {
		t.Error("GetConfig() should error for non-existent key")
	}
}

func TestHasChanges(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo
	initGitRepo(t, tmpDir)

	// Should have no changes initially
	hasChanges, err := HasChanges(tmpDir)
	if err != nil {
		t.Errorf("HasChanges() failed: %v", err)
	}
	if hasChanges {
		t.Error("Expected no changes initially")
	}

	// Create a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Should now have changes
	hasChanges, err = HasChanges(tmpDir)
	if err != nil {
		t.Errorf("HasChanges() failed: %v", err)
	}
	if !hasChanges {
		t.Error("Expected changes after creating file")
	}
}

func TestAddFiles(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo
	initGitRepo(t, tmpDir)

	// Create a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Add file
	err := AddFiles(tmpDir, []string{"test.txt"}, false)
	if err != nil {
		t.Errorf("AddFiles() failed: %v", err)
	}

	// Verify file is staged
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to check git diff: %v", err)
	}

	if len(output) == 0 {
		t.Error("File was not staged")
	}
}

func TestAddFiles_NoFiles(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo
	initGitRepo(t, tmpDir)

	// Try to add empty list
	err := AddFiles(tmpDir, []string{}, false)
	if err == nil {
		t.Error("AddFiles() with empty list should error")
	}
}

func TestCreateCommit(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo
	initGitRepo(t, tmpDir)

	// Create and add a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := AddFiles(tmpDir, []string{"test.txt"}, false); err != nil {
		t.Fatalf("AddFiles() failed: %v", err)
	}

	// Create commit
	err := CreateCommit(tmpDir, "Test commit", false)
	if err != nil {
		t.Errorf("CreateCommit() failed: %v", err)
	}

	// Verify commit exists
	cmd := exec.Command("git", "log", "--oneline", "-1")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to check git log: %v", err)
	}

	if len(output) == 0 {
		t.Error("Commit was not created")
	}
}

func TestGetRemote(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo
	initGitRepo(t, tmpDir)

	// Try to get non-existent remote
	_, err := GetRemote(tmpDir, "origin")
	if err == nil {
		t.Error("GetRemote() should error for non-existent remote")
	}
}

func TestAddRemote(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// Initialize repo
	initGitRepo(t, tmpDir)

	// Add remote
	testURL := "https://github.com/test/repo.git"
	err := AddRemote(tmpDir, "origin", testURL, false)
	if err != nil {
		t.Errorf("AddRemote() failed: %v", err)
	}

	// Verify remote was added
	result, err := GetRemote(tmpDir, "origin")
	if err != nil {
		t.Errorf("GetRemote() failed: %v", err)
	}
	if result != testURL {
		t.Errorf("Expected '%s', got '%s'", testURL, result)
	}
}

func TestIntegration_FullWorkflow(t *testing.T) {
	gitAvailable(t)
	tmpDir := createTempDir(t)
	defer cleanup(t, tmpDir)

	// 1. Initialize repo
	if err := Init(tmpDir, false); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 2. Configure user
	if err := SetConfig(tmpDir, "user.name", "Test User", false); err != nil {
		t.Fatalf("SetConfig() failed: %v", err)
	}
	if err := SetConfig(tmpDir, "user.email", "test@example.com", false); err != nil {
		t.Fatalf("SetConfig() failed: %v", err)
	}

	// 3. Create a file
	testFile := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test\n"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// 4. Check for changes
	hasChanges, err := HasChanges(tmpDir)
	if err != nil {
		t.Fatalf("HasChanges() failed: %v", err)
	}
	if !hasChanges {
		t.Fatal("Expected changes after creating file")
	}

	// 5. Add file
	if err := AddFiles(tmpDir, []string{"README.md"}, false); err != nil {
		t.Fatalf("AddFiles() failed: %v", err)
	}

	// 6. Create commit
	if err := CreateCommit(tmpDir, "Initial commit", false); err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	// 7. Add remote
	if err := AddRemote(tmpDir, "origin", "https://github.com/test/repo.git", false); err != nil {
		t.Fatalf("AddRemote() failed: %v", err)
	}

	// 8. Verify remote
	remoteURL, err := GetRemote(tmpDir, "origin")
	if err != nil {
		t.Fatalf("GetRemote() failed: %v", err)
	}
	if remoteURL != "https://github.com/test/repo.git" {
		t.Errorf("Remote URL mismatch: expected 'https://github.com/test/repo.git', got '%s'", remoteURL)
	}

	// 9. Verify no more changes (after commit)
	hasChanges, err = HasChanges(tmpDir)
	if err != nil {
		t.Fatalf("HasChanges() failed: %v", err)
	}
	if hasChanges {
		t.Error("Expected no changes after committing")
	}
}

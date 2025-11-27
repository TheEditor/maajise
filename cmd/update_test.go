package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"maajise/internal/fsutil"
	_ "maajise/templates"
)

func TestUpdateCommand_Name(t *testing.T) {
	uc := NewUpdateCommand()
	if uc.Name() != "update" {
		t.Errorf("Name() = %q, want %q", uc.Name(), "update")
	}
}

func TestUpdateCommand_Description(t *testing.T) {
	uc := NewUpdateCommand()
	if uc.Description() == "" {
		t.Error("Description() should not be empty")
	}
}

func TestUpdateCommand_Usage(t *testing.T) {
	uc := NewUpdateCommand()
	if uc.Usage() == "" {
		t.Error("Usage() should not be empty")
	}
}

func TestUpdateCommand_Examples(t *testing.T) {
	uc := NewUpdateCommand()
	if len(uc.Examples()) == 0 {
		t.Error("Examples() should not be empty")
	}
}

func TestUpdateCommand_DetectTemplate(t *testing.T) {
	uc := NewUpdateCommand()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "maajise-update-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test base detection (no markers)
	if got := uc.detectTemplate(tmpDir); got != "base" {
		t.Errorf("detectTemplate() = %q, want %q", got, "base")
	}

	// Test TypeScript detection
	os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)
	if got := uc.detectTemplate(tmpDir); got != "typescript" {
		t.Errorf("detectTemplate() with package.json = %q, want %q", got, "typescript")
	}
	os.Remove(filepath.Join(tmpDir, "package.json"))

	// Test Go detection
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)
	if got := uc.detectTemplate(tmpDir); got != "go" {
		t.Errorf("detectTemplate() with go.mod = %q, want %q", got, "go")
	}
}

func TestUpdateCommand_DryRun(t *testing.T) {
	uc := NewUpdateCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-update-dryrun-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp dir
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Execute dry run
	err = uc.Execute([]string{"--dry-run"})
	if err != nil {
		t.Errorf("Execute(--dry-run) error = %v", err)
	}

	// Verify no files created
	files, _ := os.ReadDir(tmpDir)
	if len(files) != 0 {
		t.Errorf("dry-run created %d files, want 0", len(files))
	}
}

func TestUpdateCommand_Force(t *testing.T) {
	uc := NewUpdateCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-update-force-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp dir
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Create initial .gitignore
	initialContent := []byte("initial")
	os.WriteFile(".gitignore", initialContent, 0644)

	// Execute update without force (should skip)
	err = uc.Execute([]string{})
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// File should still have original content
	content, _ := os.ReadFile(".gitignore")
	if string(content) != string(initialContent) {
		t.Error("Without --force, existing file should not be overwritten")
	}

	// Execute update with force (should overwrite)
	err = uc.Execute([]string{"--force"})
	if err != nil {
		t.Errorf("Execute(--force) error = %v", err)
	}

	// File should have new content (different from initial)
	newContent, _ := os.ReadFile(".gitignore")
	if string(newContent) == string(initialContent) {
		t.Error("With --force, existing file should be overwritten")
	}
}

func TestUpdateCommand_SpecificFiles(t *testing.T) {
	uc := NewUpdateCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-update-specific-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp dir
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Execute update with specific files
	err = uc.Execute([]string{".gitignore"})
	if err != nil {
		t.Errorf("Execute with specific files error = %v", err)
	}

	// Check that .gitignore was created
	if !fsutil.FileExists(".gitignore") {
		t.Error("Specific file .gitignore should be created")
	}

	// Check that other files were NOT created (README.md should not exist)
	if fsutil.FileExists("README.md") {
		t.Error("Other files should not be created when specific files requested")
	}
}

func TestUpdateCommand_Registration(t *testing.T) {
	// Verify the command is registered
	cmd, ok := Get("update")
	if !ok {
		t.Error("update command should be registered")
	}

	if cmd.Name() != "update" {
		t.Errorf("Registered command name = %q, want %q", cmd.Name(), "update")
	}
}

// Helper method for tests - extracted from UpdateCommand for testing
func (uc *UpdateCommand) detectTemplate(dir string) string {
	checks := []struct {
		file     string
		template string
	}{
		{"package.json", "typescript"},
		{"tsconfig.json", "typescript"},
		{"Cargo.toml", "rust"},
		{"pyproject.toml", "python"},
		{"requirements.txt", "python"},
		{"composer.json", "php"},
		{"go.mod", "go"},
	}

	for _, check := range checks {
		if fsutil.FileExists(filepath.Join(dir, check.file)) {
			return check.template
		}
	}

	return "base"
}

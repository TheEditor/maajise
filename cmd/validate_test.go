package cmd

import (
	"os"
	"path/filepath"
	"testing"

	_ "maajise/templates"
)

func TestValidateCommand_Name(t *testing.T) {
	vc := NewValidateCommand()
	if vc.Name() != "validate" {
		t.Errorf("Name() = %q, want %q", vc.Name(), "validate")
	}
}

func TestValidateCommand_Description(t *testing.T) {
	vc := NewValidateCommand()
	if vc.Description() == "" {
		t.Error("Description() should not be empty")
	}
}

func TestValidateCommand_Usage(t *testing.T) {
	vc := NewValidateCommand()
	if vc.Usage() == "" {
		t.Error("Usage() should not be empty")
	}
}

func TestValidateCommand_Examples(t *testing.T) {
	vc := NewValidateCommand()
	if len(vc.Examples()) == 0 {
		t.Error("Examples() should not be empty")
	}
}

func TestValidateCommand_CheckGit(t *testing.T) {
	vc := NewValidateCommand()

	// Create temp directory without .git
	tmpDir, err := os.MkdirTemp("", "maajise-validate-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Should fail without .git
	result := vc.checkGit(tmpDir)
	if result.Status != "fail" {
		t.Errorf("checkGit() without .git = %q, want %q", result.Status, "fail")
	}

	// Create .git directory
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0755)

	// Should pass with .git
	result = vc.checkGit(tmpDir)
	if result.Status != "pass" {
		t.Errorf("checkGit() with .git = %q, want %q", result.Status, "pass")
	}
}

func TestValidateCommand_CheckBeads(t *testing.T) {
	vc := NewValidateCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-validate-beads-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Should warn without .beads
	result := vc.checkBeads(tmpDir)
	if result.Status != "warn" {
		t.Errorf("checkBeads() without .beads = %q, want %q", result.Status, "warn")
	}

	// Create .beads directory
	os.MkdirAll(filepath.Join(tmpDir, ".beads"), 0755)

	// Should pass with .beads
	result = vc.checkBeads(tmpDir)
	if result.Status != "pass" {
		t.Errorf("checkBeads() with .beads = %q, want %q", result.Status, "pass")
	}
}

func TestValidateCommand_CheckRequiredFiles(t *testing.T) {
	vc := NewValidateCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-validate-files-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Without required files
	results := vc.checkRequiredFiles(tmpDir)
	hasGitignore := false
	hasReadme := false
	for _, r := range results {
		if r.Check == ".gitignore" && r.Status == "warn" {
			hasGitignore = true
		}
		if r.Check == "README.md" && r.Status == "warn" {
			hasReadme = true
		}
	}
	if !hasGitignore || !hasReadme {
		t.Error("Missing required files should be warnings")
	}

	// Create required files
	os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0644)

	// With required files
	results = vc.checkRequiredFiles(tmpDir)
	for _, r := range results {
		if (r.Check == ".gitignore" || r.Check == "README.md") && r.Status != "pass" {
			t.Errorf("Present required file %s should be pass, got %q", r.Check, r.Status)
		}
	}
}

func TestValidateCommand_StrictMode(t *testing.T) {
	vc := NewValidateCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-validate-strict-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create .git but not .beads (will generate warning)
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0755)
	os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0644)

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Normal mode should pass (warnings OK)
	err = vc.Execute([]string{})
	if err != nil {
		t.Errorf("Execute() normal mode error = %v, want nil", err)
	}

	// Strict mode should fail (warnings = errors)
	vc = NewValidateCommand() // Reset command
	err = vc.Execute([]string{"--strict"})
	if err == nil {
		t.Error("Execute(--strict) should fail with warnings")
	}
}

func TestValidateCommand_VerboseMode(t *testing.T) {
	vc := NewValidateCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-validate-verbose-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Set up a minimal valid project
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".beads"), 0755)
	os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0644)

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Execute with verbose flag
	vc = NewValidateCommand()
	err = vc.Execute([]string{"--verbose"})
	if err != nil {
		t.Errorf("Execute(--verbose) error = %v, want nil", err)
	}
}

func TestValidateCommand_Registration(t *testing.T) {
	// Verify the command is registered
	cmd, ok := Get("validate")
	if !ok {
		t.Error("validate command should be registered")
	}

	if cmd.Name() != "validate" {
		t.Errorf("Registered command name = %q, want %q", cmd.Name(), "validate")
	}
}

func TestValidateCommand_ValidationResult(t *testing.T) {
	result := ValidationResult{
		Check:   "test",
		Status:  "pass",
		Message: "test message",
	}

	if result.Check != "test" {
		t.Error("ValidationResult.Check not set correctly")
	}
	if result.Status != "pass" {
		t.Error("ValidationResult.Status not set correctly")
	}
	if result.Message != "test message" {
		t.Error("ValidationResult.Message not set correctly")
	}
}

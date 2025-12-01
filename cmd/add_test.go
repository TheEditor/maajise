package cmd

import (
	"os"
	"path/filepath"
	"testing"

	_ "maajise/templates"
)

func TestAddCommand_Name(t *testing.T) {
	ac := NewAddCommand()
	if ac.Name() != "add" {
		t.Errorf("Name() = %q, want %q", ac.Name(), "add")
	}
}

func TestAddCommand_Description(t *testing.T) {
	ac := NewAddCommand()
	if ac.Description() == "" {
		t.Error("Description() returned empty string")
	}
}

func TestAddCommand_Usage(t *testing.T) {
	ac := NewAddCommand()
	if ac.Usage() == "" {
		t.Error("Usage() returned empty string")
	}
}

func TestAddCommand_Examples(t *testing.T) {
	ac := NewAddCommand()
	examples := ac.Examples()
	if examples == "" {
		t.Error("Examples() returned empty slice")
	}
}

func TestAddCommand_InterfaceImplementation(t *testing.T) {
	var _ Command = NewAddCommand()
}

func TestAddCommand_AddFile_DryRun(t *testing.T) {
	ac := NewAddCommand()
	ac.dryRun = true
	ac.template = "base"

	tmpDir, err := os.MkdirTemp("", "maajise-add-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = ac.addFile(tmpDir, "test-project", ".gitignore")
	if err != nil {
		t.Errorf("addFile() error = %v", err)
	}

	// File should NOT exist (dry-run)
	if ac.fileExists(filepath.Join(tmpDir, ".gitignore")) {
		t.Error("dry-run created file")
	}
}

func TestAddCommand_AddFile_Create(t *testing.T) {
	ac := NewAddCommand()
	ac.template = "base"

	tmpDir, err := os.MkdirTemp("", "maajise-add-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = ac.addFile(tmpDir, "test-project", ".gitignore")
	if err != nil {
		t.Errorf("addFile() error = %v", err)
	}

	// File should exist
	if !ac.fileExists(filepath.Join(tmpDir, ".gitignore")) {
		t.Error("addFile() did not create file")
	}
}

func TestAddCommand_AddFile_SkipExisting(t *testing.T) {
	ac := NewAddCommand()
	ac.template = "base"
	ac.force = false

	tmpDir, err := os.MkdirTemp("", "maajise-add-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create existing file with known content
	existingPath := filepath.Join(tmpDir, ".gitignore")
	os.WriteFile(existingPath, []byte("existing content"), 0644)

	// Add should skip
	err = ac.addFile(tmpDir, "test-project", ".gitignore")
	if err != nil {
		t.Errorf("addFile() error = %v", err)
	}

	// Content should be unchanged
	content, _ := os.ReadFile(existingPath)
	if string(content) != "existing content" {
		t.Error("addFile() modified existing file without --force")
	}
}

func TestAddCommand_DetectTemplate_PackageJson(t *testing.T) {
	ac := NewAddCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-detect-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create package.json
	os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)

	detected := ac.detectTemplate(tmpDir)
	if detected != "typescript" {
		t.Errorf("detectTemplate() = %q, want %q", detected, "typescript")
	}
}

func TestAddCommand_DetectTemplate_CargoToml(t *testing.T) {
	ac := NewAddCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-detect-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Cargo.toml
	os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte("[package]"), 0644)

	detected := ac.detectTemplate(tmpDir)
	if detected != "rust" {
		t.Errorf("detectTemplate() = %q, want %q", detected, "rust")
	}
}

func TestAddCommand_DetectTemplate_Default(t *testing.T) {
	ac := NewAddCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-detect-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// No markers, should default to base
	detected := ac.detectTemplate(tmpDir)
	if detected != "base" {
		t.Errorf("detectTemplate() = %q, want %q", detected, "base")
	}
}

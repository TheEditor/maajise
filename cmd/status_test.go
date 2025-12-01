package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"maajise/internal/fsutil"
	_ "maajise/templates"
)

func TestStatusCommand_Name(t *testing.T) {
	sc := NewStatusCommand()
	if sc.Name() != "status" {
		t.Errorf("Name() = %q, want %q", sc.Name(), "status")
	}
}

func TestStatusCommand_Description(t *testing.T) {
	sc := NewStatusCommand()
	if sc.Description() == "" {
		t.Error("Description() should not be empty")
	}
}

func TestStatusCommand_Usage(t *testing.T) {
	sc := NewStatusCommand()
	if sc.Usage() == "" {
		t.Error("Usage() should not be empty")
	}
}

func TestStatusCommand_Examples(t *testing.T) {
	sc := NewStatusCommand()
	if len(sc.Examples()) == 0 {
		t.Error("Examples() should not be empty")
	}
}

func TestStatusCommand_Execute(t *testing.T) {
	sc := NewStatusCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-status-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Should not error even in empty directory
	err = sc.Run([]string{})
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestStatusCommand_WithGitInitialized(t *testing.T) {
	sc := NewStatusCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-status-git-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create .git directory
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0755)

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	err = sc.Run([]string{})
	if err != nil {
		t.Errorf("Execute() with .git error = %v", err)
	}
}

func TestStatusCommand_WithBeadsInitialized(t *testing.T) {
	sc := NewStatusCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-status-beads-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create .beads directory
	os.MkdirAll(filepath.Join(tmpDir, ".beads"), 0755)

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	err = sc.Run([]string{})
	if err != nil {
		t.Errorf("Execute() with .beads error = %v", err)
	}
}

func TestStatusCommand_WithTemplateMarkers(t *testing.T) {
	sc := NewStatusCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-status-template-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create TypeScript marker
	os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	err = sc.Run([]string{})
	if err != nil {
		t.Errorf("Execute() with template marker error = %v", err)
	}
}

func TestStatusCommand_FilePresenceChecks(t *testing.T) {
	sc := NewStatusCommand()

	tmpDir, err := os.MkdirTemp("", "maajise-status-files-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some key files
	os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0644)
	// Note: not creating .ubsignore to test both present and missing files

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	err = sc.Run([]string{})
	if err != nil {
		t.Errorf("Execute() with files error = %v", err)
	}

	// Verify files exist check works
	if !fsutil.FileExists(filepath.Join(tmpDir, ".gitignore")) {
		t.Error("FileExists should detect .gitignore")
	}
	if !fsutil.FileExists(filepath.Join(tmpDir, "README.md")) {
		t.Error("FileExists should detect README.md")
	}
	if fsutil.FileExists(filepath.Join(tmpDir, ".ubsignore")) {
		t.Error("FileExists should not detect non-existent .ubsignore")
	}
}

func TestStatusCommand_DirectoryCheck(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "maajise-status-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)

	// DirExists should work correctly
	if !fsutil.DirExists(gitDir) {
		t.Error("DirExists should detect .git directory")
	}

	// FileExists should NOT match directories
	if fsutil.FileExists(gitDir) {
		t.Error("FileExists should not match directories")
	}
}

func TestStatusCommand_Registration(t *testing.T) {
	// Verify the command is registered
	cmd, ok := Get("status")
	if !ok {
		t.Error("status command should be registered")
	}

	if cmd.Name() != "status" {
		t.Errorf("Registered command name = %q, want %q", cmd.Name(), "status")
	}
}

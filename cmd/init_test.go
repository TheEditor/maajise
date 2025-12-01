package cmd

import (
	"strings"
	"testing"

	"maajise/internal/config"
)

// contains is a helper function for substring checking in tests
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestNewInitCommand(t *testing.T) {
	ic := NewInitCommand()

	if ic.Name() != "init" {
		t.Errorf("Expected name 'init', got '%s'", ic.Name())
	}

	if ic.Description() == "" {
		t.Error("Description should not be empty")
	}

	if ic.Usage() == "" {
		t.Error("Usage should not be empty")
	}

	examples := ic.Examples()
	if examples == "" {
		t.Error("Examples should not be empty")
	}
}

func TestInitCommand_Execute_NoProjectName(t *testing.T) {
	ic := NewInitCommand()

	err := ic.Run([]string{})
	if err == nil {
		t.Error("Expected error when no project name provided")
	}

	// Check that error contains the key message and help hint
	errMsg := err.Error()
	if !contains(errMsg, "project name required") {
		t.Errorf("Expected 'project name required' in error, got: %v", err)
	}
	if !contains(errMsg, "maajise init --help") {
		t.Errorf("Expected help hint in error, got: %v", err)
	}
}

func TestInitCommand_Execute_InvalidProjectName(t *testing.T) {
	ic := NewInitCommand()

	// Invalid characters should fail
	err := ic.Run([]string{"my project!"})
	if err == nil {
		t.Error("Expected error for invalid project name")
	}
}

func TestInitCommand_Execute_ValidProjectName(t *testing.T) {
	ic := NewInitCommand()

	// Test that the command properly implements the interface
	// We can't test the actual execution without side effects,
	// but we can verify the command structure
	if ic.Name() != "init" {
		t.Errorf("Expected name 'init', got '%s'", ic.Name())
	}
}

func TestInitCommand_ValidateProjectName(t *testing.T) {
	ic := NewInitCommand()

	tests := []struct {
		name    string
		project string
		valid   bool
	}{
		{"valid hyphenated", "my-project", true},
		{"valid underscored", "my_project", true},
		{"valid alphanumeric", "myproject123", true},
		{"invalid space", "my project", false},
		{"invalid special char", "my!project", false},
		{"invalid empty", "", false},
	}

	for _, tt := range tests {
		err := ic.validateProjectName(tt.project)
		if (err == nil) != tt.valid {
			t.Errorf("validateProjectName(%q) = %v, want valid=%v", tt.project, err, tt.valid)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.InPlace != false {
		t.Error("Expected InPlace to default to false")
	}

	if cfg.SkipGit != false {
		t.Error("Expected SkipGit to default to false")
	}

	if cfg.SkipBeads != false {
		t.Error("Expected SkipBeads to default to false")
	}

	if cfg.NoOverwrite != false {
		t.Error("Expected NoOverwrite to default to false")
	}

	if cfg.MainBranch != "main" {
		t.Errorf("Expected MainBranch 'main', got '%s'", cfg.MainBranch)
	}
}

func TestCommandRegistry(t *testing.T) {
	// Get init command from registry
	cmd, ok := Get("init")
	if !ok {
		t.Error("Expected to find 'init' command in registry")
	}

	if cmd == nil {
		t.Error("Command should not be nil")
	}

	if cmd.Name() != "init" {
		t.Errorf("Expected command name 'init', got '%s'", cmd.Name())
	}
}

func TestAll(t *testing.T) {
	// Get all commands
	all := All()

	if len(all) == 0 {
		t.Error("Expected at least one command registered")
	}

	// Should have init command
	found := false
	for _, c := range all {
		if c.Name() == "init" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'init' command in All()")
	}
}

func TestInitCommand_InterfaceImplementation(t *testing.T) {
	// Verify that InitCommand properly implements the Command interface
	var _ Command = (*InitCommand)(nil)

	ic := NewInitCommand()

	// All interface methods should be callable
	_ = ic.Name()
	_ = ic.Description()
	_ = ic.Usage()
	_ = ic.Examples()
	// Execute is tested separately due to side effects
}

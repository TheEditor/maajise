package cmd

import (
	"testing"
)

func TestDoctorCommand_Name(t *testing.T) {
	dc := NewDoctorCommand()
	if dc.Name() != "doctor" {
		t.Errorf("Name() = %q, want %q", dc.Name(), "doctor")
	}
}

func TestDoctorCommand_Description(t *testing.T) {
	dc := NewDoctorCommand()
	if dc.Description() == "" {
		t.Error("Description() returned empty string")
	}
}

func TestDoctorCommand_Usage(t *testing.T) {
	dc := NewDoctorCommand()
	if dc.Usage() == "" {
		t.Error("Usage() returned empty string")
	}
}

func TestDoctorCommand_Examples(t *testing.T) {
	dc := NewDoctorCommand()
	examples := dc.Examples()
	if examples == "" {
		t.Error("Examples() returned empty slice")
	}
}

func TestDoctorCommand_InterfaceImplementation(t *testing.T) {
	var _ Command = NewDoctorCommand()
}

func TestDoctorCommand_RunCheck_Git(t *testing.T) {
	dc := NewDoctorCommand()

	check := DependencyCheck{
		Name:    "git",
		Command: "git",
		Args:    []string{"--version"},
	}

	result := dc.runCheck(check)

	// Git should be available in most environments
	if !result.Found {
		t.Skip("git not installed, skipping")
	}

	if result.Version == "" {
		t.Error("runCheck() returned empty version for git")
	}
}

func TestDoctorCommand_RunCheck_Missing(t *testing.T) {
	dc := NewDoctorCommand()

	check := DependencyCheck{
		Name:    "nonexistent",
		Command: "definitely-not-a-real-command-12345",
		Args:    []string{"--version"},
	}

	result := dc.runCheck(check)

	if result.Found {
		t.Error("runCheck() found nonexistent command")
	}
	if result.Error == "" {
		t.Error("runCheck() should set Error for missing command")
	}
}

func TestDoctorCommand_Execute(t *testing.T) {
	dc := NewDoctorCommand()
	// Execute without args should work
	err := dc.Run([]string{})
	if err == nil {
		// Expected: may error if required dependencies missing, but that's ok for test
		return
	}
	// Check if it's the expected error about required dependencies
	if err.Error() == "required dependencies missing" {
		return
	}
	t.Errorf("Execute() returned unexpected error: %v", err)
}

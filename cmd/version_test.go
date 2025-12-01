package cmd

import (
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	vc := NewVersionCommand()

	if vc.Name() != "version" {
		t.Errorf("Expected name 'version', got '%s'", vc.Name())
	}

	if vc.Description() == "" {
		t.Error("Description should not be empty")
	}

	if vc.Usage() == "" {
		t.Error("Usage should not be empty")
	}

	examples := vc.Examples()
	if examples == "" {
		t.Error("Examples should not be empty")
	}
}

func TestVersionCommand_InterfaceImplementation(t *testing.T) {
	// Verify VersionCommand implements Command interface
	var _ Command = (*VersionCommand)(nil)
}

func TestVersionCommand_Run(t *testing.T) {
	vc := NewVersionCommand()

	err := vc.Run([]string{})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestVersionCommand_LongDescription(t *testing.T) {
	vc := NewVersionCommand()

	if vc.LongDescription() == "" {
		t.Error("LongDescription should not be empty")
	}
}

func TestVersionConstant(t *testing.T) {
	if Version == "" {
		t.Error("Version constant should not be empty")
	}

	if Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%s'", Version)
	}
}

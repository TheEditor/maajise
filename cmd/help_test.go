package cmd

import (
	"testing"
)

func TestNewHelpCommand(t *testing.T) {
	hc := NewHelpCommand()

	if hc.Name() != "help" {
		t.Errorf("Expected name 'help', got '%s'", hc.Name())
	}

	if hc.Description() == "" {
		t.Error("Description should not be empty")
	}

	if hc.Usage() == "" {
		t.Error("Usage should not be empty")
	}

	examples := hc.Examples()
	if len(examples) == 0 {
		t.Error("Examples should not be empty")
	}
}

func TestHelpCommand_InterfaceImplementation(t *testing.T) {
	// Verify HelpCommand implements Command interface
	var _ Command = (*HelpCommand)(nil)
}

func TestHelpCommand_Execute_GeneralHelp(t *testing.T) {
	hc := NewHelpCommand()

	// Execute with no args should show general help (no error)
	err := hc.Execute([]string{})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestHelpCommand_Execute_CommandHelp(t *testing.T) {
	hc := NewHelpCommand()

	// Execute with valid command name
	err := hc.Execute([]string{"init"})
	if err != nil {
		t.Errorf("Expected no error for valid command, got: %v", err)
	}
}

func TestHelpCommand_Execute_InvalidCommand(t *testing.T) {
	hc := NewHelpCommand()

	// Execute with invalid command name should error
	err := hc.Execute([]string{"nonexistent"})
	if err == nil {
		t.Error("Expected error for invalid command")
	}

	expectedMsg := "unknown command"
	if err.Error() != "unknown command: nonexistent" {
		t.Errorf("Expected error containing '%s', got: %v", expectedMsg, err)
	}
}

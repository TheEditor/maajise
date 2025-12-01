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
	if examples == "" {
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
	err := hc.Run([]string{})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestHelpCommand_Execute_CommandHelp(t *testing.T) {
	hc := NewHelpCommand()

	// Execute with valid command name
	err := hc.Run([]string{"init"})
	if err != nil {
		t.Errorf("Expected no error for valid command, got: %v", err)
	}
}

func TestHelpCommand_Execute_InvalidCommand(t *testing.T) {
	hc := NewHelpCommand()

	// Execute with invalid command name should error
	err := hc.Run([]string{"nonexistent"})
	if err == nil {
		t.Error("Expected error for invalid command")
	}

	expectedMsg := "unknown command"
	if err.Error() != "unknown command: nonexistent" {
		t.Errorf("Expected error containing '%s', got: %v", expectedMsg, err)
	}
}

// TestCommandLongDescription verifies all commands have LongDescription
func TestCommandLongDescription(t *testing.T) {
	commands := All()

	for _, cmd := range commands {
		t.Run(cmd.Name(), func(t *testing.T) {
			// Verify LongDescription is not empty
			longDesc := cmd.LongDescription()
			if longDesc == "" {
				t.Errorf("%s: LongDescription() returned empty string", cmd.Name())
			}

			// Verify it's different from short description
			desc := cmd.Description()
			if longDesc == desc {
				t.Errorf("%s: LongDescription() should be more detailed than Description()", cmd.Name())
			}

			// Verify minimum length (should be more substantial than Description)
			if len(longDesc) <= len(desc) {
				t.Errorf("%s: LongDescription() should be longer than Description()", cmd.Name())
			}
		})
	}
}

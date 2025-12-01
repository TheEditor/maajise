package ui

import (
	"strings"
	"testing"
)

func TestUsageError(t *testing.T) {
	err := UsageError("init", "project name required")
	errMsg := err.Error()

	// Check error contains user message
	if !strings.Contains(errMsg, "project name required") {
		t.Error("Error message missing user message")
	}

	// Check error contains help hint
	if !strings.Contains(errMsg, "maajise init --help") {
		t.Error("Error message missing help hint")
	}

	// Check format (blank line separator)
	if !strings.Contains(errMsg, "\n\n") {
		t.Error("Error message should have blank line between message and hint")
	}
}

func TestUsageErrorDifferentCommands(t *testing.T) {
	tests := []struct {
		cmd      string
		msg      string
		expected string
	}{
		{"init", "missing name", "maajise init --help"},
		{"add", "invalid item", "maajise add --help"},
		{"validate", "no project", "maajise validate --help"},
		{"update", "bad template", "maajise update --help"},
		{"help", "unknown command", "maajise help --help"},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			err := UsageError(tt.cmd, tt.msg)
			if !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("Expected %q in error, got: %s", tt.expected, err.Error())
			}
		})
	}
}

func TestUsageErrorFormat(t *testing.T) {
	err := UsageError("init", "test message")
	errMsg := err.Error()

	parts := strings.Split(errMsg, "\n\n")
	if len(parts) != 2 {
		t.Errorf("Expected 2 parts separated by blank line, got %d parts", len(parts))
	}

	// First part should be the user message
	if parts[0] != "test message" {
		t.Errorf("First part should be user message, got: %s", parts[0])
	}

	// Second part should contain help hint
	if !strings.Contains(parts[1], "maajise init --help") {
		t.Errorf("Second part should contain help hint, got: %s", parts[1])
	}
}

package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	_ "maajise/templates" // Import to register templates
)

func TestTemplatesCommand(t *testing.T) {
	tc := NewTemplatesCommand()

	if tc.Name() != "templates" {
		t.Errorf("Name() = %q, want %q", tc.Name(), "templates")
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := tc.Execute([]string{})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check output contains expected templates
	if !strings.Contains(output, "base") {
		t.Error("output should contain 'base' template")
	}

	if !strings.Contains(output, "Available templates") {
		t.Error("output should contain header")
	}
}

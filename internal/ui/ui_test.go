package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// Helper to capture stdout
func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Helper to capture stderr
func captureStderr(fn func()) string {
	r, w, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestError(t *testing.T) {
	output := captureStderr(func() {
		Error("test error message")
	})

	if !strings.Contains(output, "test error message") {
		t.Errorf("Error() output missing message: %q", output)
	}
	if !strings.Contains(output, "✗") {
		t.Errorf("Error() output missing checkmark: %q", output)
	}
}

func TestSuccess(t *testing.T) {
	output := captureStdout(func() {
		Success("test success message")
	})

	if !strings.Contains(output, "test success message") {
		t.Errorf("Success() output missing message: %q", output)
	}
	if !strings.Contains(output, "✓") {
		t.Errorf("Success() output missing checkmark: %q", output)
	}
}

func TestInfo(t *testing.T) {
	output := captureStdout(func() {
		Info("test info message")
	})

	if !strings.Contains(output, "test info message") {
		t.Errorf("Info() output missing message: %q", output)
	}
	if !strings.Contains(output, "→") {
		t.Errorf("Info() output missing arrow: %q", output)
	}
}

func TestWarn(t *testing.T) {
	output := captureStdout(func() {
		Warn("test warning message")
	})

	if !strings.Contains(output, "test warning message") {
		t.Errorf("Warn() output missing message: %q", output)
	}
	if !strings.Contains(output, "⚠") {
		t.Errorf("Warn() output missing warning symbol: %q", output)
	}
}

func TestHeader(t *testing.T) {
	output := captureStdout(func() {
		Header("Test Header")
	})

	if !strings.Contains(output, "Test Header") {
		t.Errorf("Header() output missing title: %q", output)
	}
	if !strings.Contains(output, "╔") {
		t.Errorf("Header() output missing box character: %q", output)
	}
}

func TestSummary(t *testing.T) {
	output := captureStdout(func() {
		Summary("Line 1", "Line 2", "Line 3")
	})

	if !strings.Contains(output, "Line 1") {
		t.Errorf("Summary() output missing line 1: %q", output)
	}
	if !strings.Contains(output, "Line 2") {
		t.Errorf("Summary() output missing line 2: %q", output)
	}
	if !strings.Contains(output, "Line 3") {
		t.Errorf("Summary() output missing line 3: %q", output)
	}
	if !strings.Contains(output, "╔") {
		t.Errorf("Summary() output missing box character: %q", output)
	}
}

func TestSummary_EmptyLines(t *testing.T) {
	output := captureStdout(func() {
		Summary()
	})

	if !strings.Contains(output, "╔") {
		t.Errorf("Summary() with no lines output missing box: %q", output)
	}
}

func TestAnsiColors_Defined(t *testing.T) {
	// Verify ANSI color constants are defined
	if Red == "" {
		t.Error("Red color constant is empty")
	}
	if Green == "" {
		t.Error("Green color constant is empty")
	}
	if Yellow == "" {
		t.Error("Yellow color constant is empty")
	}
	if Blue == "" {
		t.Error("Blue color constant is empty")
	}
	if Reset == "" {
		t.Error("Reset color constant is empty")
	}
}

func TestColorConstants_Format(t *testing.T) {
	// Verify color constants have expected ANSI format
	tests := []struct {
		name  string
		color string
	}{
		{"Red", Red},
		{"Green", Green},
		{"Yellow", Yellow},
		{"Blue", Blue},
		{"Reset", Reset},
	}

	for _, tt := range tests {
		if !strings.HasPrefix(tt.color, "\033[") {
			t.Errorf("%s color doesn't start with ANSI escape sequence", tt.name)
		}
		if !strings.HasSuffix(tt.color, "m") {
			t.Errorf("%s color doesn't end with 'm'", tt.name)
		}
	}
}

func TestIntegration_MultipleOutputs(t *testing.T) {
	output := captureStdout(func() {
		Header("Test Session")
		Info("Starting process")
		Success("Step 1 complete")
		Warn("Be careful")
		Summary("Final summary", "All done")
	})

	// Verify all messages appear
	messages := []string{"Test Session", "Starting process", "Step 1 complete", "Be careful", "Final summary", "All done"}
	for _, msg := range messages {
		if !strings.Contains(output, msg) {
			t.Errorf("Integration output missing message: %q", msg)
		}
	}

	// Verify multiple box characters appear
	boxCount := strings.Count(output, "╔")
	if boxCount < 2 {
		t.Errorf("Expected at least 2 box characters, got %d", boxCount)
	}
}

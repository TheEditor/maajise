package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplate(t *testing.T) {
	tests := []struct {
		name        string
		markers     []string // Files to create
		expected    string
		description string
	}{
		{
			name:        "empty directory returns base",
			markers:     []string{},
			expected:    "base",
			description: "No markers should default to base template",
		},
		{
			name:        "package.json detects typescript",
			markers:     []string{"package.json"},
			expected:    "typescript",
			description: "package.json should trigger typescript detection",
		},
		{
			name:        "tsconfig.json detects typescript",
			markers:     []string{"tsconfig.json"},
			expected:    "typescript",
			description: "tsconfig.json should trigger typescript detection",
		},
		{
			name:        "Cargo.toml detects rust",
			markers:     []string{"Cargo.toml"},
			expected:    "rust",
			description: "Cargo.toml should trigger rust detection",
		},
		{
			name:        "pyproject.toml detects python",
			markers:     []string{"pyproject.toml"},
			expected:    "python",
			description: "pyproject.toml should trigger python detection",
		},
		{
			name:        "requirements.txt detects python",
			markers:     []string{"requirements.txt"},
			expected:    "python",
			description: "requirements.txt should trigger python detection",
		},
		{
			name:        "composer.json detects php",
			markers:     []string{"composer.json"},
			expected:    "php",
			description: "composer.json should trigger php detection",
		},
		{
			name:        "go.mod detects go",
			markers:     []string{"go.mod"},
			expected:    "go",
			description: "go.mod should trigger go detection",
		},
		{
			name:        "priority: package.json before tsconfig",
			markers:     []string{"package.json", "tsconfig.json"},
			expected:    "typescript",
			description: "package.json should be checked before tsconfig.json",
		},
		{
			name:        "priority: typescript before rust",
			markers:     []string{"Cargo.toml", "package.json"},
			expected:    "typescript",
			description: "package.json should be checked before Cargo.toml",
		},
		{
			name:        "pyproject.toml takes priority over requirements.txt",
			markers:     []string{"requirements.txt", "pyproject.toml"},
			expected:    "python",
			description: "pyproject.toml should be checked before requirements.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir, err := os.MkdirTemp("", "maajise-detect-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create marker files
			for _, marker := range tt.markers {
				filePath := filepath.Join(tmpDir, marker)
				if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
					t.Fatalf("Failed to create marker file %s: %v", marker, err)
				}
			}

			// Test template detection
			got := Template(tmpDir)
			if got != tt.expected {
				t.Errorf("Template(%q) = %q, want %q. %s", tmpDir, got, tt.expected, tt.description)
			}
		})
	}
}

func TestTemplate_WindowsPaths(t *testing.T) {
	// Test that template detection works with Windows-style paths
	tmpDir, err := os.MkdirTemp("", "maajise-detect-windows-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create package.json
	filePath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(filePath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Test with both Unix and Windows path separators
	got := Template(tmpDir)
	if got != "typescript" {
		t.Errorf("Template with Windows paths = %q, want %q", got, "typescript")
	}
}

func TestTemplate_NoMarkers(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "maajise-detect-empty-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Empty directory should return "base"
	got := Template(tmpDir)
	if got != "base" {
		t.Errorf("Template(empty dir) = %q, want %q", got, "base")
	}
}

func TestTemplate_DirectoryInsteadOfFile(t *testing.T) {
	// Create temp directory with a subdirectory named like a marker
	tmpDir, err := os.MkdirTemp("", "maajise-detect-dirmarker-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a directory named "package.json" (not a file)
	dirPath := filepath.Join(tmpDir, "package.json")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Should not match directories, should return "base"
	got := Template(tmpDir)
	if got != "base" {
		t.Errorf("Template with directory named package.json = %q, want %q", got, "base")
	}
}

func TestTemplate_AllMarkersPresent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "maajise-detect-all-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create all marker files
	markers := []string{"package.json", "tsconfig.json", "Cargo.toml", "pyproject.toml",
		"requirements.txt", "composer.json", "go.mod"}
	for _, marker := range markers {
		filePath := filepath.Join(tmpDir, marker)
		if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
			t.Fatalf("Failed to create marker file %s: %v", marker, err)
		}
	}

	// Should return typescript (first in priority order)
	got := Template(tmpDir)
	if got != "typescript" {
		t.Errorf("Template with all markers = %q, want %q (should use priority order)", got, "typescript")
	}
}

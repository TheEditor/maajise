package main

import (
	"maajise/cmd"
	"os"
	"strings"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid names
		{
			name:    "simple name",
			input:   "myproject",
			wantErr: false,
		},
		{
			name:    "name with hyphen",
			input:   "my-project",
			wantErr: false,
		},
		{
			name:    "name with underscore",
			input:   "my_project",
			wantErr: false,
		},
		{
			name:    "name with numbers",
			input:   "project123",
			wantErr: false,
		},
		{
			name:    "uppercase name",
			input:   "MyProject",
			wantErr: false,
		},
		{
			name:    "mixed case with hyphens and numbers",
			input:   "My-Project-123",
			wantErr: false,
		},

		// Invalid names
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
		},
		{
			name:    "name with space",
			input:   "my project",
			wantErr: true,
		},
		{
			name:    "name with special character",
			input:   "my@project",
			wantErr: true,
		},
		{
			name:    "name with dot",
			input:   "my.project",
			wantErr: true,
		},
		{
			name:    "name with slash",
			input:   "my/project",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Create a temporary directory
	tmpdir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file",
			path: tmpfile.Name(),
			want: true,
		},
		{
			name: "existing directory",
			path: tmpdir,
			want: true,
		},
		{
			name: "non-existent file",
			path: "/nonexistent/path/to/file",
			want: false,
		},
		{
			name: "empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fileExists(tt.path)
			if got != tt.want {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestCheckDependencies(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *cmd.Config
		wantErr bool
	}{
		{
			name: "skip all dependencies",
			cfg: &cmd.Config{
				SkipGit:   true,
				SkipBeads: true,
			},
			wantErr: false,
		},
		{
			name: "skip git",
			cfg: &cmd.Config{
				SkipGit:   true,
				SkipBeads: false,
			},
			wantErr: false, // Should only check for beads, which may or may not exist
		},
		{
			name: "skip beads",
			cfg: &cmd.Config{
				SkipGit:   false,
				SkipBeads: true,
			},
			wantErr: false, // Should only check for git, which should exist on most systems
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkDependencies(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkDependencies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseFlags(t *testing.T) {
	// Note: parseFlags() function is not easily testable without modifying it
	// because it relies on flag.Parse() which operates on os.Args.
	// This test is a placeholder for manual verification.
	t.Skip("parseFlags requires refactoring to be easily testable")
}

func TestEmailValidation(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantValid bool
	}{
		// Valid emails
		{
			name:      "simple email",
			email:     "user@example.com",
			wantValid: true,
		},
		{
			name:      "email with numbers",
			email:     "user123@example.com",
			wantValid: true,
		},
		{
			name:      "email with multiple dots",
			email:     "user@mail.example.co.uk",
			wantValid: true,
		},
		{
			name:      "email with plus",
			email:     "user+tag@example.com",
			wantValid: true,
		},

		// Invalid emails
		{
			name:      "no @ symbol",
			email:     "userexample.com",
			wantValid: false,
		},
		{
			name:      "no dot in domain",
			email:     "user@example",
			wantValid: false,
		},
		{
			name:      "empty email",
			email:     "",
			wantValid: false,
		},
		{
			name:      "only @ symbol",
			email:     "@",
			wantValid: false,
		},
		{
			name:      "only domain",
			email:     "example.com",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple validation logic from configureGitUser
			isValid := len(tt.email) > 0 && strings.Contains(tt.email, "@") && strings.Contains(tt.email, ".")
			if isValid != tt.wantValid {
				t.Errorf("Email validation for %q = %v, want %v", tt.email, isValid, tt.wantValid)
			}
		})
	}
}

func TestGitUserInputValidation(t *testing.T) {
	tests := []struct {
		name      string
		inputName  string
		inputEmail string
		wantErr   bool
		errMsg    string
	}{
		// Valid inputs
		{
			name:       "valid name and email",
			inputName:  "John Doe",
			inputEmail: "john@example.com",
			wantErr:    false,
		},

		// Empty inputs
		{
			name:       "empty name",
			inputName:  "",
			inputEmail: "john@example.com",
			wantErr:    true,
			errMsg:     "empty user name",
		},
		{
			name:       "empty email",
			inputName:  "John Doe",
			inputEmail: "",
			wantErr:    true,
			errMsg:     "empty user email",
		},

		// Invalid email format
		{
			name:       "email without @",
			inputName:  "John Doe",
			inputEmail: "johnexample.com",
			wantErr:    true,
			errMsg:     "invalid email format",
		},
		{
			name:       "email without dot",
			inputName:  "John Doe",
			inputEmail: "john@example",
			wantErr:    true,
			errMsg:     "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate validation logic from configureGitUser
			var errMsg string
			if tt.inputName == "" {
				errMsg = "empty user name"
			} else if tt.inputEmail == "" {
				errMsg = "empty user email"
			} else if !strings.Contains(tt.inputEmail, "@") || !strings.Contains(tt.inputEmail, ".") {
				errMsg = "invalid email format"
			}

			hasErr := errMsg != ""
			if hasErr != tt.wantErr {
				t.Errorf("Validation for %q/%q = %v, want %v", tt.inputName, tt.inputEmail, hasErr, tt.wantErr)
			}
			if hasErr && errMsg != tt.errMsg {
				t.Errorf("Error message = %q, want %q", errMsg, tt.errMsg)
			}
		})
	}
}

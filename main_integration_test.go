package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestSetupGitRemote(t *testing.T) {
	tests := []struct {
		name        string
		remoteURL   string
		shouldError bool
		checkRemote bool
	}{
		{
			name:        "valid https github url",
			remoteURL:   "https://github.com/test/repo.git\n",
			shouldError: false,
			checkRemote: true,
		},
		{
			name:        "valid https gitlab url",
			remoteURL:   "https://gitlab.com/test/repo.git\n",
			shouldError: false,
			checkRemote: true,
		},
		{
			name:        "valid ssh url",
			remoteURL:   "git@github.com:test/repo.git\n",
			shouldError: false,
			checkRemote: true,
		},
		{
			name:        "invalid url with semicolon",
			remoteURL:   "https://github.com/test/repo.git; echo hacked\n",
			shouldError: true,
			checkRemote: false,
		},
		{
			name:        "invalid url with backtick",
			remoteURL:   "https://github.com/test/`whoami`.git\n",
			shouldError: true,
			checkRemote: false,
		},
		{
			name:        "invalid url with pipe",
			remoteURL:   "https://github.com/test/repo.git | cat /etc/passwd\n",
			shouldError: true,
			checkRemote: false,
		},
		{
			name:        "invalid format",
			remoteURL:   "not-a-url\n",
			shouldError: true,
			checkRemote: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary git repository
			tmpDir, err := os.MkdirTemp("", "test-repo")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Initialize git repo
			cmd := exec.Command("git", "init")
			cmd.Dir = tmpDir
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to init git repo: %v", err)
			}

			// Mock stdin with the remote URL
			cfg := &Config{
				ProjectName: "test-project",
				SkipGit:     false,
				SkipRemote:  false,
				Verbose:     false,
			}

			// Redirect stdin to provide the remote URL
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			defer r.Close()

			// Write test input
			if _, err := w.WriteString("y\n" + tt.remoteURL); err != nil {
				t.Fatalf("Failed to write to pipe: %v", err)
			}
			w.Close()

			os.Stdin = r

			// Call setupGitRemote
			err = setupGitRemote(tmpDir, cfg)

			// Check error status
			if (err != nil) != tt.shouldError {
				t.Errorf("setupGitRemote() error = %v, shouldError %v", err, tt.shouldError)
			}

			// If setup succeeded, verify the remote was actually added
			if !tt.shouldError && tt.checkRemote {
				cmd := exec.Command("git", "remote", "get-url", "origin")
				cmd.Dir = tmpDir
				output, err := cmd.Output()
				if err != nil {
					t.Errorf("Failed to get remote: %v", err)
					return
				}

				actual := string(bytes.TrimSpace(output))
				expected := string(bytes.TrimSpace([]byte(tt.remoteURL)))
				if actual != expected {
					t.Errorf("Remote URL = %q, want %q", actual, expected)
				}
			}

			// If setup failed, verify no remote was added
			if tt.shouldError {
				cmd := exec.Command("git", "remote", "get-url", "origin")
				cmd.Dir = tmpDir
				_, err := cmd.Output()
				// Error is expected when remote doesn't exist
				if err == nil {
					t.Errorf("Expected remote to not exist, but it was added")
				}
			}
		})
	}
}

func TestSetupGitRemoteSkipFlags(t *testing.T) {
	tests := []struct {
		name      string
		skipGit   bool
		skipRemote bool
	}{
		{
			name:      "skip git flag",
			skipGit:   true,
			skipRemote: false,
		},
		{
			name:      "skip remote flag",
			skipGit:   false,
			skipRemote: true,
		},
		{
			name:      "both skip flags",
			skipGit:   true,
			skipRemote: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary git repository
			tmpDir, err := os.MkdirTemp("", "test-repo")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Initialize git repo
			cmd := exec.Command("git", "init")
			cmd.Dir = tmpDir
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to init git repo: %v", err)
			}

			cfg := &Config{
				ProjectName: "test-project",
				SkipGit:     tt.skipGit,
				SkipRemote:  tt.skipRemote,
				Verbose:     false,
			}

			// setupGitRemote should return nil without prompting
			err = setupGitRemote(tmpDir, cfg)
			if err != nil {
				t.Errorf("setupGitRemote() error = %v, want nil", err)
			}

			// Verify no remote was added
			cmd = exec.Command("git", "remote", "get-url", "origin")
			cmd.Dir = tmpDir
			_, err = cmd.Output()
			if err == nil {
				t.Errorf("Expected no remote to be added when flags are set")
			}
		})
	}
}

func TestSetupGitRemoteExistingRemote(t *testing.T) {
	// Create a temporary git repository
	tmpDir, err := os.MkdirTemp("", "test-repo")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Add an existing remote
	cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/existing/repo.git")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add remote: %v", err)
	}

	cfg := &Config{
		ProjectName: "test-project",
		SkipGit:     false,
		SkipRemote:  false,
		Verbose:     false,
	}

	// setupGitRemote should return nil without prompting
	err = setupGitRemote(tmpDir, cfg)
	if err != nil {
		t.Errorf("setupGitRemote() error = %v, want nil", err)
	}

	// Verify the original remote is still there
	cmd = exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Failed to get remote: %v", err)
		return
	}

	actual := string(bytes.TrimSpace(output))
	expected := "https://github.com/existing/repo.git"
	if actual != expected {
		t.Errorf("Remote URL = %q, want %q", actual, expected)
	}
}

func TestSetupGitRemoteEmptyURL(t *testing.T) {
	// Create a temporary git repository
	tmpDir, err := os.MkdirTemp("", "test-repo")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	cfg := &Config{
		ProjectName: "test-project",
		SkipGit:     false,
		SkipRemote:  false,
		Verbose:     false,
	}

	// Mock stdin - answer 'y' then provide empty URL
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()

	// Write test input: yes, then empty line
	if _, err := w.WriteString("y\n\n"); err != nil {
		t.Fatalf("Failed to write to pipe: %v", err)
	}
	w.Close()

	os.Stdin = r

	// setupGitRemote should return nil without error
	err = setupGitRemote(tmpDir, cfg)
	if err != nil {
		t.Errorf("setupGitRemote() error = %v, want nil", err)
	}

	// Verify no remote was added
	cmd = exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = tmpDir
	_, err = cmd.Output()
	if err == nil {
		t.Errorf("Expected no remote to be added when empty URL provided")
	}
}

func TestSetupGitRemoteNoAnswer(t *testing.T) {
	// Create a temporary git repository
	tmpDir, err := os.MkdirTemp("", "test-repo")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	cfg := &Config{
		ProjectName: "test-project",
		SkipGit:     false,
		SkipRemote:  false,
		Verbose:     false,
	}

	// Mock stdin - answer 'n' to skip remote setup
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()

	// Write test input: no
	if _, err := w.WriteString("n\n"); err != nil {
		t.Fatalf("Failed to write to pipe: %v", err)
	}
	w.Close()

	os.Stdin = r

	// setupGitRemote should return nil without error
	err = setupGitRemote(tmpDir, cfg)
	if err != nil {
		t.Errorf("setupGitRemote() error = %v, want nil", err)
	}

	// Verify no remote was added
	cmd = exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = tmpDir
	_, err = cmd.Output()
	if err == nil {
		t.Errorf("Expected no remote to be added when user declines")
	}
}

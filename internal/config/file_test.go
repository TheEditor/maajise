package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadFileConfig_NoFile(t *testing.T) {
	// Should return empty config, not error
	cfg, err := LoadFileConfig()
	if err != nil {
		t.Errorf("LoadFileConfig() error = %v", err)
	}
	if cfg == nil {
		t.Error("LoadFileConfig() returned nil config")
	}
}

func TestFileConfig_RoundTrip(t *testing.T) {
	// Create temp config file
	tmpDir, err := os.MkdirTemp("", "maajise-config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test config
	testCfg := &FileConfig{}
	testCfg.Defaults.Template = "typescript"
	testCfg.Defaults.GitName = "Test User"
	testCfg.Variables.Author = "Test Author"

	// Write to temp file
	path := filepath.Join(tmpDir, ".maajiserc")
	data := []byte(`defaults:
  template: typescript
  git_name: Test User
variables:
  author: Test Author
`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Read back - note: this tests parsing, not LoadFileConfig path logic
	cfg := &FileConfig{}
	fileData, _ := os.ReadFile(path)
	if err := yaml.Unmarshal(fileData, cfg); err != nil {
		t.Errorf("yaml.Unmarshal error = %v", err)
	}

	if cfg.Defaults.Template != "typescript" {
		t.Errorf("Template = %q, want %q", cfg.Defaults.Template, "typescript")
	}
	if cfg.Variables.Author != "Test Author" {
		t.Errorf("Author = %q, want %q", cfg.Variables.Author, "Test Author")
	}
}

func TestMergeFileConfig(t *testing.T) {
	fc := &FileConfig{}
	fc.Defaults.Template = "rust"
	fc.Defaults.GitName = "Config User"
	fc.Variables.Author = "Config Author"

	c := DefaultConfig()
	c.Template = "" // Empty, should be filled
	c.GitName = ""
	c.GitEmail = "user@example.com" // Already set, should not be overwritten

	c.MergeFileConfig(fc)

	if c.Template != "rust" {
		t.Errorf("Template after merge = %q, want %q", c.Template, "rust")
	}
	if c.GitName != "Config User" {
		t.Errorf("GitName after merge = %q, want %q", c.GitName, "Config User")
	}
	// Email was already set, should not change
	if c.GitEmail != "user@example.com" {
		t.Errorf("GitEmail after merge = %q, want %q", c.GitEmail, "user@example.com")
	}
}

func TestDefaultFileConfig(t *testing.T) {
	cfg := DefaultFileConfig()
	if cfg.Defaults.Template != "base" {
		t.Errorf("Default template = %q, want %q", cfg.Defaults.Template, "base")
	}
	if cfg.Defaults.MainBranch != "main" {
		t.Errorf("Default MainBranch = %q, want %q", cfg.Defaults.MainBranch, "main")
	}
	if cfg.Variables.Year == "" {
		t.Error("Default Year should not be empty")
	}
	if cfg.Variables.License != "MIT" {
		t.Errorf("Default License = %q, want %q", cfg.Variables.License, "MIT")
	}
}

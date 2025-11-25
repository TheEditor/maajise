package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfig represents the ~/.maajiserc configuration file
type FileConfig struct {
	// Default values for init command
	Defaults struct {
		Template    string `yaml:"template"`
		GitName     string `yaml:"git_name"`
		GitEmail    string `yaml:"git_email"`
		SkipRemote  bool   `yaml:"skip_remote"`
		SkipBeads   bool   `yaml:"skip_beads"`
		MainBranch  string `yaml:"main_branch"`
	} `yaml:"defaults"`

	// Template variables for substitution
	Variables struct {
		Author  string `yaml:"author"`
		Email   string `yaml:"email"`
		Year    string `yaml:"year"`
		License string `yaml:"license"`
		GitHub  string `yaml:"github"`
	} `yaml:"variables"`

	// Custom template directory
	TemplatesDir string `yaml:"templates_dir"`
}

// ConfigPath returns the path to the config file
func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".maajiserc")
}

// LoadFileConfig loads configuration from ~/.maajiserc
// Returns empty config (not error) if file doesn't exist
func LoadFileConfig() (*FileConfig, error) {
	cfg := &FileConfig{}

	path := ConfigPath()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // No config file is fine
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveFileConfig saves configuration to ~/.maajiserc
func SaveFileConfig(cfg *FileConfig) error {
	path := ConfigPath()
	if path == "" {
		return os.ErrNotExist
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ConfigExists returns true if config file exists
func ConfigExists() bool {
	path := ConfigPath()
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

// DefaultFileConfig returns a FileConfig with example values (for init)
func DefaultFileConfig() *FileConfig {
	cfg := &FileConfig{}
	cfg.Defaults.Template = "base"
	cfg.Defaults.MainBranch = "main"
	cfg.Variables.Year = "2025"
	cfg.Variables.License = "MIT"
	return cfg
}

// MergeFileConfig applies FileConfig values to Config where Config has zero values
func (c *Config) MergeFileConfig(fc *FileConfig) {
	if c.Template == "" && fc.Defaults.Template != "" {
		c.Template = fc.Defaults.Template
	}
	if c.GitName == "" && fc.Defaults.GitName != "" {
		c.GitName = fc.Defaults.GitName
	}
	if c.GitEmail == "" && fc.Defaults.GitEmail != "" {
		c.GitEmail = fc.Defaults.GitEmail
	}
	if c.MainBranch == "" && fc.Defaults.MainBranch != "" {
		c.MainBranch = fc.Defaults.MainBranch
	}
	// Bool flags: only apply if explicitly set in file AND not overridden by CLI
	// This requires more sophisticated flag tracking - defer to P4-007 (interactive mode)
}

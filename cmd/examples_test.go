package cmd

import (
	"strings"
	"testing"
)

// TestInitExamplesIncludeAllFlags verifies InitCommand documents all flags
func TestInitExamplesIncludeAllFlags(t *testing.T) {
	cmd := NewInitCommand()
	examples := cmd.Examples()

	// All flags that should be documented
	requiredFlags := []string{
		"--in-place", "--no-overwrite", "--template",
		"--skip-git", "--skip-beads", "--skip-commit",
		"--skip-remote", "--skip-git-user",
		"--git-name", "--git-email",
		"--dry-run", "--interactive", "--verbose",
	}

	for _, flag := range requiredFlags {
		if !strings.Contains(examples, flag) {
			t.Errorf("InitCommand examples missing flag: %s", flag)
		}
	}

	// Check that all template options are listed
	templates := []string{"base", "typescript", "python", "rust", "php", "go"}
	for _, tmpl := range templates {
		if !strings.Contains(examples, tmpl) {
			t.Errorf("InitCommand examples missing template: %s", tmpl)
		}
	}
}

// TestAddExamplesIncludeAllFlags verifies AddCommand documents all flags
func TestAddExamplesIncludeAllFlags(t *testing.T) {
	cmd := NewAddCommand()
	examples := cmd.Examples()

	requiredFlags := []string{"--force", "--template", "--dry-run", "--verbose"}

	for _, flag := range requiredFlags {
		if !strings.Contains(examples, flag) {
			t.Errorf("AddCommand examples missing flag: %s", flag)
		}
	}
}

// TestUpdateExamplesIncludeAllFlags verifies UpdateCommand documents all flags
func TestUpdateExamplesIncludeAllFlags(t *testing.T) {
	cmd := NewUpdateCommand()
	examples := cmd.Examples()

	requiredFlags := []string{"--force", "--template", "--dry-run", "--verbose"}

	for _, flag := range requiredFlags {
		if !strings.Contains(examples, flag) {
			t.Errorf("UpdateCommand examples missing flag: %s", flag)
		}
	}
}

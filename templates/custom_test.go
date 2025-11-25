package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCustomTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "maajise-custom-template-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test template
	templateContent := `name: test-custom
description: A test custom template
dependencies:
  - git
files:
  README.md: |
    # {{.ProjectName}}
    Custom template test
`
	templatePath := filepath.Join(tmpDir, "test-custom.yaml")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Load it
	err = loadCustomTemplate(templatePath)
	if err != nil {
		t.Errorf("loadCustomTemplate() error = %v", err)
	}

	// Check it was registered
	tmpl, ok := Get("test-custom")
	if !ok {
		t.Error("custom template not registered")
	}

	if tmpl.Description() != "A test custom template" {
		t.Errorf("Description = %q", tmpl.Description())
	}

	files := tmpl.Files("my-project")
	if !strings.Contains(files["README.md"], "my-project") {
		t.Error("ProjectName not substituted")
	}
}

func TestLoadCustomTemplates_MissingDir(t *testing.T) {
	// Should not error if directory doesn't exist
	err := LoadCustomTemplates("/nonexistent/path")
	if err != nil {
		t.Errorf("LoadCustomTemplates() with missing dir error = %v", err)
	}
}

func TestLoadCustomTemplates_EmptyDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "maajise-custom-templates-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Should not error if directory is empty
	err = LoadCustomTemplates(tmpDir)
	if err != nil {
		t.Errorf("LoadCustomTemplates() with empty dir error = %v", err)
	}
}

func TestDefaultCustomTemplatesDir(t *testing.T) {
	dir := DefaultCustomTemplatesDir()
	if !strings.Contains(dir, ".maajise") {
		t.Errorf("DefaultCustomTemplatesDir() = %q, should contain .maajise", dir)
	}
	if !strings.HasSuffix(dir, "templates") {
		t.Errorf("DefaultCustomTemplatesDir() = %q, should end with templates", dir)
	}
}

func TestCustomTemplate_Implementation(t *testing.T) {
	var _ Template = &CustomTemplate{}
}

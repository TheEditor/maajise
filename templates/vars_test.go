package templates

import (
	"strings"
	"testing"
)

func TestProcessTemplate(t *testing.T) {
	vars := TemplateVars{
		ProjectName: "test-project",
		Author:      "John Doe",
		Email:       "john@example.com",
		Year:        "2025",
		License:     "MIT",
	}

	content := "Author: {{.Author}} <{{.Email}}>\nYear: {{.Year}}"
	result, err := ProcessTemplate(content, vars)
	if err != nil {
		t.Errorf("ProcessTemplate() error = %v", err)
	}

	if !strings.Contains(result, "John Doe") {
		t.Error("result should contain author name")
	}
	if !strings.Contains(result, "john@example.com") {
		t.Error("result should contain email")
	}
	if !strings.Contains(result, "2025") {
		t.Error("result should contain year")
	}
}

func TestProcessTemplate_NoVars(t *testing.T) {
	vars := TemplateVars{ProjectName: "test"}
	content := "No variables here, just plain text"

	result, err := ProcessTemplate(content, vars)
	if err != nil {
		t.Errorf("ProcessTemplate() error = %v", err)
	}

	if result != content {
		t.Errorf("ProcessTemplate() modified content without variables")
	}
}

func TestDefaultVars(t *testing.T) {
	vars := DefaultVars("my-project")

	if vars.ProjectName != "my-project" {
		t.Errorf("ProjectName = %q, want %q", vars.ProjectName, "my-project")
	}
	if vars.Year == "" {
		t.Error("Year should not be empty")
	}
	if vars.License != "MIT" {
		t.Errorf("License = %q, want %q", vars.License, "MIT")
	}
}

func TestFilesWithVars(t *testing.T) {
	vars := TemplateVars{
		ProjectName: "my-app",
		Author:      "Jane Smith",
		Year:        "2024",
	}

	// Use base template
	tmpl, ok := Get("base")
	if !ok {
		t.Skip("base template not available")
	}

	files := FilesWithVars(tmpl, vars)

	// README should have been processed
	if readme, ok := files["README.md"]; ok {
		if !strings.Contains(readme, "my-app") {
			t.Error("ProjectName not substituted in README")
		}
	}
}

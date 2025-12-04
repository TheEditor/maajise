package templates

import (
	"strings"
	"testing"
)

func TestGoTemplate(t *testing.T) {
	tmpl, ok := Get("go")
	if !ok {
		t.Fatal("go template not registered")
	}

	if tmpl.Name() != "go" {
		t.Errorf("Name() = %q, want %q", tmpl.Name(), "go")
	}

	files := tmpl.Files("test-project")

	expectedFiles := []string{".gitignore", ".ubsignore", "README.md", "go.mod"}
	for _, f := range expectedFiles {
		if _, ok := files[f]; !ok {
			t.Errorf("missing file: %s", f)
		}
	}

	if !strings.Contains(files["go.mod"], "test-project") {
		t.Error("go.mod should contain project name")
	}

	if !strings.Contains(files["go.mod"], "go 1.23") {
		t.Error("go.mod should specify go version")
	}
}

func TestGoTemplateFiles(t *testing.T) {
	tmpl := &GoTemplate{}
	files := tmpl.Files("myapp")

	// Verify required files exist
	requiredFiles := []string{
		".gitignore",
		".ubsignore",
		"README.md",
		"go.mod",
		"cmd/myapp/main.go",
		"api/.gitkeep",
		"build/.gitkeep",
		"configs/.gitkeep",
		"docs/.gitkeep",
		"internal/.gitkeep",
		"pkg/.gitkeep",
		"scripts/.gitkeep",
		"test/.gitkeep",
	}

	for _, f := range requiredFiles {
		if _, exists := files[f]; !exists {
			t.Errorf("Missing required file: %s", f)
		}
	}

	// Verify old flat main.go doesn't exist
	if _, exists := files["main.go"]; exists {
		t.Error("main.go should not exist at root (should be in cmd/<project>/)")
	}

	// Verify we have the expected number of files (4 original + 8 .gitkeep + 1 cmd/*/main.go)
	expectedCount := 13
	if len(files) != expectedCount {
		t.Errorf("Files() returned %d files, expected %d", len(files), expectedCount)
	}
}

func TestGoTemplateDescription(t *testing.T) {
	tmpl := &GoTemplate{}
	desc := tmpl.Description()

	if !strings.Contains(desc, "Standard Go Project Layout") {
		t.Error("Description should mention Standard Go Project Layout")
	}
}

func TestGoTemplateREADME(t *testing.T) {
	tmpl := &GoTemplate{}
	readme := tmpl.readme("test-app")

	requiredSections := []string{
		"## Project Structure",
		"Standard Go Project Layout",
		"cmd/test-app/",
		"internal/",
		"pkg/",
		"api/",
		"## Getting Started",
	}

	for _, section := range requiredSections {
		if !strings.Contains(readme, section) {
			t.Errorf("README missing section/content: %s", section)
		}
	}
}

func TestGoTemplateMainGo(t *testing.T) {
	tmpl := &GoTemplate{}
	main := tmpl.mainGo()

	if !strings.Contains(main, "TODO") {
		t.Error("main.go should contain TODO comment")
	}

	if !strings.Contains(main, "internal/") {
		t.Error("main.go should mention internal/ imports")
	}
}

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

	expectedFiles := []string{".gitignore", ".ubsignore", "README.md", "go.mod", "main.go"}
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

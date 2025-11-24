package templates

import (
	"strings"
	"testing"
)

func TestPythonTemplate(t *testing.T) {
	tmpl, ok := Get("python")
	if !ok {
		t.Fatal("python template not registered")
	}

	files := tmpl.Files("test-project")

	expectedFiles := []string{".gitignore", ".ubsignore", "README.md", "pyproject.toml", "requirements.txt", "src/__init__.py", "src/main.py"}
	for _, f := range expectedFiles {
		if _, ok := files[f]; !ok {
			t.Errorf("missing file: %s", f)
		}
	}

	if !strings.Contains(files["pyproject.toml"], "test-project") {
		t.Error("pyproject.toml should contain project name")
	}

	if !strings.Contains(files[".gitignore"], "__pycache__") {
		t.Error(".gitignore should contain __pycache__")
	}

	if !strings.Contains(files[".gitignore"], "venv") {
		t.Error(".gitignore should contain venv")
	}
}

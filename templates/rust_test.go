package templates

import (
	"strings"
	"testing"
)

func TestRustTemplate(t *testing.T) {
	tmpl, ok := Get("rust")
	if !ok {
		t.Fatal("rust template not registered")
	}

	files := tmpl.Files("test-project")

	expectedFiles := []string{".gitignore", ".ubsignore", "README.md", "Cargo.toml", "src/main.rs"}
	for _, f := range expectedFiles {
		if _, ok := files[f]; !ok {
			t.Errorf("missing file: %s", f)
		}
	}

	if !strings.Contains(files["Cargo.toml"], "test-project") {
		t.Error("Cargo.toml should contain project name")
	}

	if !strings.Contains(files[".gitignore"], "/target/") {
		t.Error(".gitignore should contain /target/")
	}
}

package templates

import (
	"strings"
	"testing"
)

func TestPHPTemplate(t *testing.T) {
	tmpl, ok := Get("php")
	if !ok {
		t.Fatal("php template not registered")
	}

	files := tmpl.Files("test-project")

	expectedFiles := []string{".gitignore", ".ubsignore", "README.md", "composer.json", "src/index.php"}
	for _, f := range expectedFiles {
		if _, ok := files[f]; !ok {
			t.Errorf("missing file: %s", f)
		}
	}

	if !strings.Contains(files["composer.json"], "test-project") {
		t.Error("composer.json should contain project name")
	}

	if !strings.Contains(files[".gitignore"], "/vendor/") {
		t.Error(".gitignore should contain /vendor/")
	}
}

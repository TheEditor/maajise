package templates

import (
	"strings"
	"testing"
)

func TestTypeScriptTemplate(t *testing.T) {
	tmpl, ok := Get("typescript")
	if !ok {
		t.Fatal("typescript template not registered")
	}

	if tmpl.Name() != "typescript" {
		t.Errorf("Name() = %q, want %q", tmpl.Name(), "typescript")
	}

	files := tmpl.Files("test-project")

	expectedFiles := []string{".gitignore", ".ubsignore", "README.md", "package.json", "tsconfig.json", "src/index.ts"}
	for _, f := range expectedFiles {
		if _, ok := files[f]; !ok {
			t.Errorf("missing file: %s", f)
		}
	}

	// Check package.json contains project name
	if !strings.Contains(files["package.json"], "test-project") {
		t.Error("package.json should contain project name")
	}

	// Check node_modules in gitignore
	if !strings.Contains(files[".gitignore"], "node_modules") {
		t.Error(".gitignore should contain node_modules")
	}
}

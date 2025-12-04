package templates

import (
	"encoding/json"
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

func TestTypeScriptTemplateDirectoryStructure(t *testing.T) {
	tmpl := &TypeScriptTemplate{}
	files := tmpl.Files("test-app")

	// Verify all .gitkeep files exist with empty content
	expectedDirs := []string{
		"src/config/.gitkeep",
		"src/controllers/.gitkeep",
		"src/middleware/.gitkeep",
		"src/models/.gitkeep",
		"src/routes/.gitkeep",
		"src/services/.gitkeep",
		"src/types/.gitkeep",
		"src/utils/.gitkeep",
		"tests/.gitkeep",
	}

	for _, path := range expectedDirs {
		content, exists := files[path]
		if !exists {
			t.Errorf("Missing .gitkeep file: %s", path)
		}
		if content != "" {
			t.Errorf(".gitkeep %s should be empty, got: %q", path, content)
		}
	}

	// Verify we have the expected number of files (6 original + 9 .gitkeep)
	expectedCount := 15
	if len(files) != expectedCount {
		t.Errorf("Files() returned %d files, expected %d", len(files), expectedCount)
	}
}

func TestTypeScriptTemplatePathAliases(t *testing.T) {
	tmpl := &TypeScriptTemplate{}
	tsconfig := tmpl.tsconfig()

	// Verify baseUrl
	if !strings.Contains(tsconfig, `"baseUrl": "./src"`) {
		t.Error("tsconfig.json should have baseUrl set to ./src")
	}

	// Verify all path aliases exist
	requiredAliases := []string{
		`"@config/*"`,
		`"@controllers/*"`,
		`"@services/*"`,
		`"@models/*"`,
		`"@middleware/*"`,
		`"@routes/*"`,
		`"@types/*"`,
		`"@utils/*"`,
	}

	for _, alias := range requiredAliases {
		if !strings.Contains(tsconfig, alias) {
			t.Errorf("tsconfig.json missing path alias: %s", alias)
		}
	}

	// Verify resolveJsonModule is present
	if !strings.Contains(tsconfig, `"resolveJsonModule": true`) {
		t.Error("tsconfig.json should have resolveJsonModule: true")
	}
}

func TestTypeScriptTemplateDescription(t *testing.T) {
	tmpl := &TypeScriptTemplate{}
	desc := tmpl.Description()

	if !strings.Contains(desc, "layered architecture") {
		t.Error("Description should mention layered architecture")
	}
}

func TestTypeScriptTemplateREADMESections(t *testing.T) {
	tmpl := &TypeScriptTemplate{}
	readme := tmpl.readme("test-app")

	requiredSections := []string{
		"## Project Structure",
		"## Architecture",
		"## Development",
		"## Path Aliases",
		"src/config/",
		"src/controllers/",
		"src/services/",
		"src/models/",
		"@config/*",
		"Routes →",
		"Controllers →",
	}

	for _, section := range requiredSections {
		if !strings.Contains(readme, section) {
			t.Errorf("README missing section/content: %s", section)
		}
	}
}

func TestTypeScriptTemplateIndexTODOs(t *testing.T) {
	tmpl := &TypeScriptTemplate{}
	index := tmpl.indexTS()

	if !strings.Contains(index, "TODO") {
		t.Error("index.ts should contain TODO comments")
	}

	if !strings.Contains(index, "Example imports") {
		t.Error("index.ts should contain Example imports")
	}

	if !strings.Contains(index, "Hello, TypeScript") {
		t.Error("index.ts should contain Hello message")
	}
}

func TestTypeScriptTemplateJSONValidity(t *testing.T) {
	tmpl := &TypeScriptTemplate{}
	files := tmpl.Files("test-app")

	// Verify package.json is valid JSON
	var pkg map[string]interface{}
	if err := json.Unmarshal([]byte(files["package.json"]), &pkg); err != nil {
		t.Errorf("package.json is not valid JSON: %v", err)
	}

	// Verify tsconfig.json is valid JSON
	var tsconfig map[string]interface{}
	if err := json.Unmarshal([]byte(files["tsconfig.json"]), &tsconfig); err != nil {
		t.Errorf("tsconfig.json is not valid JSON: %v", err)
	}
}

func TestTypeScriptTemplateCleanScript(t *testing.T) {
	tmpl := &TypeScriptTemplate{}
	pkg := tmpl.packageJSON("test-app")

	if !strings.Contains(pkg, `"clean"`) {
		t.Error("package.json should contain clean script")
	}

	// Verify it uses Node.js (cross-platform)
	if !strings.Contains(pkg, "node") || !strings.Contains(pkg, "rmSync") {
		t.Error("clean script should use Node.js fs.rmSync for cross-platform support")
	}
}

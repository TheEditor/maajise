package templates

import (
	"strings"
	"testing"
)

func TestPHPTemplateFiles(t *testing.T) {
	tmpl := &PHPTemplate{}
	files := tmpl.Files("myapp")

	// Verify required files exist
	requiredFiles := []string{
		".gitignore",
		".ubsignore",
		"README.md",
		"composer.json",
		"bin/.gitkeep",
		"config/.gitkeep",
		"public/index.php",
		"public/assets/.gitkeep",
		"src/Controllers/.gitkeep",
		"src/Models/.gitkeep",
		"src/Services/.gitkeep",
		"src/Middleware/.gitkeep",
		"tests/.gitkeep",
		"var/cache/.gitkeep",
		"var/logs/.gitkeep",
		"vendor/.gitkeep",
	}

	for _, f := range requiredFiles {
		if _, exists := files[f]; !exists {
			t.Errorf("Missing required file: %s", f)
		}
	}

	// Verify old src/index.php doesn't exist
	if _, exists := files["src/index.php"]; exists {
		t.Error("src/index.php should not exist (should be public/index.php)")
	}
}

func TestPHPComposerHasPSR4(t *testing.T) {
	tmpl := &PHPTemplate{}
	composer := tmpl.composerJSON("myapp")

	if !strings.Contains(composer, `"psr-4"`) {
		t.Error("composer.json should contain PSR-4 autoload configuration")
	}
	if !strings.Contains(composer, `"App\\"`) {
		t.Error("composer.json should map App\\ namespace")
	}
	if !strings.Contains(composer, `"Tests\\"`) {
		t.Error("composer.json should map Tests\\ namespace in autoload-dev")
	}
	if !strings.Contains(composer, `"autoload-dev"`) {
		t.Error("composer.json should have autoload-dev section")
	}
}

func TestPHPIndexUsesAutoloader(t *testing.T) {
	tmpl := &PHPTemplate{}
	index := tmpl.indexPHP("myapp")

	if !strings.Contains(index, "vendor/autoload.php") {
		t.Error("index.php should require Composer autoloader")
	}
	if !strings.Contains(index, "../vendor/autoload.php") {
		t.Error("index.php should use correct relative path to autoloader")
	}
	if !strings.Contains(index, "declare(strict_types=1)") {
		t.Error("index.php should have strict types declaration")
	}
}

func TestPHPTemplateDescription(t *testing.T) {
	tmpl := &PHPTemplate{}
	desc := tmpl.Description()

	if !strings.Contains(desc, "PSR-4") {
		t.Error("Description should mention PSR-4")
	}
	if !strings.Contains(desc, "public/") {
		t.Error("Description should mention public/ directory")
	}
}

func TestPHPTemplate(t *testing.T) {
	tmpl, ok := Get("php")
	if !ok {
		t.Fatal("php template not registered")
	}

	files := tmpl.Files("test-project")

	expectedFiles := []string{".gitignore", ".ubsignore", "README.md", "composer.json", "public/index.php"}
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

	if !strings.Contains(files[".gitignore"], "/var/cache/") {
		t.Error(".gitignore should contain /var/cache/ pattern")
	}
}

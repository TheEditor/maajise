package templates

import (
	"strings"
	"testing"
)

func TestSwiftTemplate(t *testing.T) {
	tmpl, ok := Get("swift")
	if !ok {
		t.Fatal("swift template not registered")
	}

	if got := tmpl.Description(); got != "Swift package with SwiftPM configuration (Sources/, Tests/)" {
		t.Errorf("Description() = %q, want %q", got, "Swift package with SwiftPM configuration (Sources/, Tests/)")
	}

	deps := tmpl.Dependencies()
	expectedDeps := []string{"git", "br", "swift"}
	if len(deps) != len(expectedDeps) {
		t.Fatalf("Dependencies() length = %d, want %d", len(deps), len(expectedDeps))
	}
	for i := range expectedDeps {
		if deps[i] != expectedDeps[i] {
			t.Errorf("Dependencies()[%d] = %q, want %q", i, deps[i], expectedDeps[i])
		}
	}

	files := tmpl.Files("myapp")

	expectedFiles := []string{
		".gitignore",
		".ubsignore",
		"README.md",
		"Package.swift",
		"Sources/myapp/main.swift",
		"Tests/myappTests/myappTests.swift",
	}

	if len(files) != len(expectedFiles) {
		t.Fatalf("Files() returned %d files, want %d", len(files), len(expectedFiles))
	}

	for _, f := range expectedFiles {
		if _, ok := files[f]; !ok {
			t.Errorf("missing file: %s", f)
		}
	}

	if !strings.Contains(files["Package.swift"], "name: \"myapp\"") {
		t.Error("Package.swift should contain project name")
	}

	if !strings.Contains(files["Tests/myappTests/myappTests.swift"], "@testable import myapp") {
		t.Error("test file should import module name")
	}

	if !strings.Contains(files["README.md"], "git add -f Package.resolved") {
		t.Error("README should document how to commit Package.resolved when needed")
	}
}

func TestSwiftTemplate_HyphenProjectName(t *testing.T) {
	tmpl, ok := Get("swift")
	if !ok {
		t.Fatal("swift template not registered")
	}

	files := tmpl.Files("my-app")
	testFile := files["Tests/my-appTests/my-appTests.swift"]

	if !strings.Contains(testFile, "@testable import my_app") {
		t.Error("hyphenated project names should normalize import module with underscores")
	}
}

package detect

import (
	"path/filepath"

	"maajise/internal/fsutil"
)

// Template detects the project template based on marker files found in the directory.
// It checks for language-specific files in a specific order and returns the corresponding
// template name. If no markers are found, it returns "base".
func Template(dir string) string {
	checks := []struct {
		file     string
		template string
	}{
		{"package.json", "typescript"},
		{"tsconfig.json", "typescript"},
		{"Cargo.toml", "rust"},
		{"pyproject.toml", "python"},
		{"requirements.txt", "python"},
		{"composer.json", "php"},
		{"go.mod", "go"},
	}

	for _, check := range checks {
		if fsutil.FileExists(filepath.Join(dir, check.file)) {
			return check.template
		}
	}

	return "base"
}

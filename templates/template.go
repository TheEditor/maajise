package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// Template defines the interface for project templates
type Template interface {
	Name() string
	Description() string
	Files(projectName string) map[string]string
	Dependencies() []string
}

// TemplateVars holds variables for template substitution
type TemplateVars struct {
	ProjectName string
	Author      string
	Email       string
	Year        string
	License     string
	GitHub      string
}

// DefaultVars returns TemplateVars with sensible defaults
func DefaultVars(projectName string) TemplateVars {
	return TemplateVars{
		ProjectName: projectName,
		Year:        time.Now().Format("2006"),
		License:     "MIT",
	}
}

// ProcessTemplate applies variable substitution to content
func ProcessTemplate(content string, vars TemplateVars) (string, error) {
	tmpl, err := template.New("content").Parse(content)
	if err != nil {
		// If template parsing fails, return original content
		// This handles files that aren't meant to have variables
		return content, nil
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return content, nil
	}

	return buf.String(), nil
}

// FilesWithVars returns template files with variable substitution applied
func FilesWithVars(tmpl Template, vars TemplateVars) map[string]string {
	files := tmpl.Files(vars.ProjectName)
	result := make(map[string]string, len(files))

	for name, content := range files {
		processed, _ := ProcessTemplate(content, vars)
		result[name] = processed
	}

	return result
}

// Registry holds all available templates
var registry = make(map[string]Template)

// Register adds a template to the registry
func Register(tmpl Template) {
	registry[tmpl.Name()] = tmpl
}

// Get retrieves a template by name
func Get(name string) (Template, bool) {
	tmpl, ok := registry[name]
	return tmpl, ok
}

// All returns all registered templates
func All() []Template {
	tmpls := make([]Template, 0, len(registry))
	for _, tmpl := range registry {
		tmpls = append(tmpls, tmpl)
	}
	return tmpls
}

// List returns template names
func List() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// CustomTemplate represents a user-defined template loaded from YAML
type CustomTemplate struct {
	name        string
	description string
	deps        []string
	files       map[string]string
}

func (t *CustomTemplate) Name() string        { return t.name }
func (t *CustomTemplate) Description() string { return t.description }
func (t *CustomTemplate) Dependencies() []string { return t.deps }
func (t *CustomTemplate) Files(projectName string) map[string]string {
	// Replace {{.ProjectName}} in file contents
	result := make(map[string]string, len(t.files))
	for name, content := range t.files {
		// Simple replacement for project name
		processed := strings.ReplaceAll(content, "{{.ProjectName}}", projectName)
		result[name] = processed
	}
	return result
}

// CustomTemplateFile represents the YAML structure for custom templates
type CustomTemplateFile struct {
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Dependencies []string          `yaml:"dependencies"`
	Files        map[string]string `yaml:"files"`
}

// LoadCustomTemplates loads templates from a directory
func LoadCustomTemplates(dir string) error {
	if dir == "" {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist, that's OK
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		path := filepath.Join(dir, name)
		if err := loadCustomTemplate(path); err != nil {
			// Log warning but continue
			// Can't import ui here, just skip
		}
	}

	return nil
}

func loadCustomTemplate(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var ctf CustomTemplateFile
	if err := yaml.Unmarshal(data, &ctf); err != nil {
		return err
	}

	if ctf.Name == "" {
		return nil // Skip templates without a name
	}

	tmpl := &CustomTemplate{
		name:        ctf.Name,
		description: ctf.Description,
		deps:        ctf.Dependencies,
		files:       ctf.Files,
	}

	Register(tmpl)
	return nil
}

// DefaultCustomTemplatesDir returns ~/.maajise/templates/
func DefaultCustomTemplatesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".maajise", "templates")
}

func init() {
	// Load custom templates from default directory
	LoadCustomTemplates(DefaultCustomTemplatesDir())
}

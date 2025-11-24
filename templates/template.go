package templates

// Template defines the interface for project templates
type Template interface {
	Name() string
	Description() string
	Files(projectName string) map[string]string
	Dependencies() []string
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

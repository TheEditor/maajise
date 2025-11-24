package templates

import "fmt"

func init() {
	Register(&RustTemplate{})
}

// RustTemplate is a Rust project template
type RustTemplate struct{}

func (t *RustTemplate) Name() string {
	return "rust"
}

func (t *RustTemplate) Description() string {
	return "Rust project with Cargo configuration"
}

func (t *RustTemplate) Dependencies() []string {
	return []string{"git", "bd", "cargo", "rustc"}
}

func (t *RustTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore":  t.gitignore(),
		".ubsignore":  t.ubsignore(),
		"README.md":   t.readme(projectName),
		"Cargo.toml":  t.cargoToml(projectName),
		"src/main.rs": t.mainRs(),
	}
}

func (t *RustTemplate) gitignore() string {
	return `# Rust
/target/
Cargo.lock

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
`
}

func (t *RustTemplate) ubsignore() string {
	return `# UBS Scanner Ignore File
target/
.git/
.vscode/
.idea/
.beads/
.claude/
*.md
*.toml
*.lock
`
}

func (t *RustTemplate) readme(projectName string) string {
	return fmt.Sprintf(`# %s

A Rust project.

## Setup

` + "```bash" + `
# Build the project
cargo build

# Run the project
cargo run
` + "```" + `

## Development

` + "```bash" + `
# Run in release mode
cargo run --release

# Run tests
cargo test

# Check code without building
cargo check

# Format code
cargo fmt

# Lint
cargo clippy
` + "```" + `

## Issue Tracking

` + "```bash" + `
bd list           # View issues
bd create --title "Task"  # Create issue
` + "```" + `

## Code Quality

` + "```bash" + `
ubs .             # Scan for bugs
` + "```" + `
`, projectName)
}

func (t *RustTemplate) cargoToml(projectName string) string {
	return fmt.Sprintf(`[package]
name = "%s"
version = "0.1.0"
edition = "2021"

[dependencies]
`, projectName)
}

func (t *RustTemplate) mainRs() string {
	return `fn main() {
    println!("Hello, Rust!");
}
`
}

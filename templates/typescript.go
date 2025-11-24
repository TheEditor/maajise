package templates

import "fmt"

func init() {
	Register(&TypeScriptTemplate{})
}

// TypeScriptTemplate is a TypeScript project template
type TypeScriptTemplate struct{}

func (t *TypeScriptTemplate) Name() string {
	return "typescript"
}

func (t *TypeScriptTemplate) Description() string {
	return "TypeScript project with npm configuration"
}

func (t *TypeScriptTemplate) Dependencies() []string {
	return []string{"git", "bd", "node", "npm"}
}

func (t *TypeScriptTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore":    t.gitignore(),
		".ubsignore":    t.ubsignore(),
		"README.md":     t.readme(projectName),
		"package.json":  t.packageJSON(projectName),
		"tsconfig.json": t.tsconfig(),
		"src/index.ts":  t.indexTS(),
	}
}

func (t *TypeScriptTemplate) gitignore() string {
	return `# Dependencies
node_modules/

# Build outputs
dist/
build/
*.js
*.js.map
*.d.ts
!*.config.js

# Logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Environment
.env
.env.local
.env*.local

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Testing
coverage/
.nyc_output/

# Cache
.cache/
*.tsbuildinfo
`
}

func (t *TypeScriptTemplate) ubsignore() string {
	return `# UBS Scanner Ignore File
node_modules/
dist/
build/
coverage/
.git/
.vscode/
.idea/
.beads/
.claude/
*.md
*.json
*.lock
*.log
`
}

func (t *TypeScriptTemplate) readme(projectName string) string {
	return fmt.Sprintf(`# %s

A TypeScript project.

## Setup

` + "```bash" + `
npm install
` + "```" + `

## Development

` + "```bash" + `
# Run in development mode
npm run dev

# Build for production
npm run build

# Run tests
npm test
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

func (t *TypeScriptTemplate) packageJSON(projectName string) string {
	return fmt.Sprintf(`{
  "name": "%s",
  "version": "0.1.0",
  "description": "",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "dev": "tsc --watch",
    "start": "node dist/index.js",
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}
`, projectName)
}

func (t *TypeScriptTemplate) tsconfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
`
}

func (t *TypeScriptTemplate) indexTS() string {
	return `// Entry point
console.log("Hello, TypeScript!");
`
}

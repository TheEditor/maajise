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
	return "TypeScript project with layered architecture (controllers, services, models, routes)"
}

func (t *TypeScriptTemplate) Dependencies() []string {
	return []string{"git", "bd", "node", "npm"}
}

func (t *TypeScriptTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore":               t.gitignore(),
		".ubsignore":               t.ubsignore(),
		"README.md":                t.readme(projectName),
		"package.json":             t.packageJSON(projectName),
		"tsconfig.json":            t.tsconfig(),
		"src/index.ts":             t.indexTS(),
		"src/config/.gitkeep":      "",
		"src/controllers/.gitkeep": "",
		"src/middleware/.gitkeep":  "",
		"src/models/.gitkeep":      "",
		"src/routes/.gitkeep":      "",
		"src/services/.gitkeep":    "",
		"src/types/.gitkeep":       "",
		"src/utils/.gitkeep":       "",
		"tests/.gitkeep":           "",
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

## Project Structure

This project uses a layered architecture for scalability and maintainability.

| Directory | Purpose |
|-----------|---------|
| src/config/ | Application configuration and environment variables |
| src/controllers/ | Request handlers and HTTP logic |
| src/middleware/ | Express/HTTP middleware functions |
| src/models/ | Data models and schemas |
| src/routes/ | API route definitions |
| src/services/ | Business logic layer |
| src/types/ | TypeScript type definitions and interfaces |
| src/utils/ | Utility functions and helpers |
| tests/ | Unit and integration tests |
| dist/ | Compiled JavaScript output |

## Architecture

Request flow: Routes → Controllers → Services → Models

- **Routes**: Define endpoints, validate input
- **Controllers**: Handle HTTP request/response
- **Services**: Business logic, reusable across controllers
- **Models**: Data structures, database schemas

## Development

### Install dependencies
` + "```bash" + `
npm install
` + "```" + `

### Build
` + "```bash" + `
npm run build
` + "```" + `

### Run
` + "```bash" + `
npm start
` + "```" + `

### Development mode (watch)
` + "```bash" + `
npm run dev
` + "```" + `

### Clean build artifacts
` + "```bash" + `
npm run clean
` + "```" + `

## Path Aliases

This project uses path aliases for cleaner imports:

- @config/* → src/config/*
- @controllers/* → src/controllers/*
- @services/* → src/services/*
- @models/* → src/models/*
- @middleware/* → src/middleware/*
- @routes/* → src/routes/*
- @types/* → src/types/*
- @utils/* → src/utils/*

Example: import { logger } from '@utils/logger';

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
    "clean": "node -e \"require('fs').rmSync('dist',{recursive:true,force:true})\"",
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
    "sourceMap": true,
    "resolveJsonModule": true,
    "baseUrl": "./src",
    "paths": {
      "@config/*": ["config/*"],
      "@controllers/*": ["controllers/*"],
      "@services/*": ["services/*"],
      "@models/*": ["models/*"],
      "@middleware/*": ["middleware/*"],
      "@routes/*": ["routes/*"],
      "@types/*": ["types/*"],
      "@utils/*": ["utils/*"]
    }
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "tests"]
}
`
}

func (t *TypeScriptTemplate) indexTS() string {
	return `// Entry point
console.log("Hello, TypeScript!");

// TODO: Initialize your application here
// Example structure:
// 1. Load configuration from ./config
// 2. Set up middleware from ./middleware
// 3. Register routes from ./routes
// 4. Start server or run application logic

// Example imports (uncomment as you build):
// import { config } from './config';
// import { setupRoutes } from './routes';
// import { logger } from './utils';
`
}

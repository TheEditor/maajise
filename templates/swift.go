package templates

import (
	"fmt"
	"strings"
)

func init() {
	Register(&SwiftTemplate{})
}

// SwiftTemplate scaffolds a minimal SwiftPM executable project.
type SwiftTemplate struct{}

func (t *SwiftTemplate) Name() string {
	return "swift"
}

func (t *SwiftTemplate) Description() string {
	return "Swift package with SwiftPM configuration (Sources/, Tests/)"
}

func (t *SwiftTemplate) Dependencies() []string {
	return []string{"git", "br", "swift"}
}

func (t *SwiftTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		".gitignore": t.gitignore(),
		".ubsignore": t.ubsignore(),
		"README.md": t.readme(projectName),
		"Package.swift": t.packageSwift(projectName),
		"Sources/" + projectName + "/main.swift": t.mainSwift(),
		"Tests/" + projectName + "Tests/" + projectName + "Tests.swift": t.testSwift(projectName),
	}
}

func (t *SwiftTemplate) gitignore() string {
	return `.DS_Store
/.build
/Packages
xcuserdata/
DerivedData/
.swiftpm/configuration/registries.json
.swiftpm/xcode/package.xcworkspace/contents.xcworkspacedata
.netrc
Package.resolved
`
}

func (t *SwiftTemplate) ubsignore() string {
	return `# UBS Scanner Ignore File
.build/
DerivedData/
Packages/
.swiftpm/
xcuserdata/
.git/
.beads/
*.md
`
}

func (t *SwiftTemplate) readme(projectName string) string {
	return fmt.Sprintf(`# %s

## Prerequisites

- Swift 6.0+ (`+"`swift --version`"+`)
- Xcode 16+ (optional, for IDE workflow on macOS)

## Build

`+"```bash"+`
swift build
`+"```"+`

## Run

`+"```bash"+`
swift run %s
`+"```"+`

## Test

`+"```bash"+`
swift test
`+"```"+`

## Project Structure

`+"```text"+`
%s/
├── Package.swift
├── Sources/
│   └── %s/
│       └── main.swift
└── Tests/
    └── %sTests/
        └── %sTests.swift
`+"```"+`

## Package.resolved

`+"`Package.resolved`"+` is gitignored by default during early development.
If you want reproducible dependency resolution for an app/release workflow, commit it explicitly:

`+"```bash"+`
git add -f Package.resolved
`+"```"+`
`, projectName, projectName, projectName, projectName, projectName, projectName)
}

func (t *SwiftTemplate) packageSwift(projectName string) string {
	return fmt.Sprintf(`// swift-tools-version: 6.0
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "%s",
    targets: [
        .executableTarget(
            name: "%s"
        ),
        .testTarget(
            name: "%sTests",
            dependencies: ["%s"]
        ),
    ]
)
`, projectName, projectName, projectName, projectName)
}

func (t *SwiftTemplate) mainSwift() string {
	return `import Foundation

print("Hello, Swift!")
`
}

func (t *SwiftTemplate) testSwift(projectName string) string {
	moduleName := strings.ReplaceAll(projectName, "-", "_")

	return fmt.Sprintf(`import Testing
@testable import %s

@Test func example() {
    #expect(1 + 1 == 2)
}
`, moduleName)
}

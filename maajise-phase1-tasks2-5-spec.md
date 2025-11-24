# Maajise Phase 1 Tasks 2-5 - Complete Architecture Migration

## Project Context
**Project:** Maajise CLI Tool  
**Phase:** 1 - Complete Architecture Migration  
**Tasks:** 2-5 (help, version, main.go refactor, internal packages)  
**Date:** 2025-11-23  
**Prerequisites:** Phase 1 Task 1 completed (cmd/init.go exists)

---

## ⚠️ IMPLEMENTATION WORKFLOW

**CRITICAL FIRST STEP:** Create all Beads issues before implementing

### Execution Steps

**STEP 1:** Verify Phase 1 Task 1 is complete
```bash
bd list
# All Task 1 issues should be closed
ls cmd/
# Should see: command.go, config.go, init.go
```

**STEP 2:** Read entire specification to understand scope

**STEP 3:** Create all Beads issues (see "CREATE ISSUES NOW" section below)

**STEP 4:** Verify issues created
```bash
bd list
# Should show all 12 new issues (ISSUE-006 through ISSUE-017)
```

**STEP 5:** Work through issues following "Implementation Order" section

**STEP 6:** Mark issues complete as you go
```bash
bd close <issue-id>
```

**STEP 7:** Verify completion using Success Criteria checklist

---

## Issues

### TASK 2: Create cmd/help.go

#### ISSUE-006: Create HelpCommand Structure

**Severity:** CRITICAL  
**File:** New file: `cmd/help.go`  
**Description:** Create help command that auto-discovers and displays all registered commands. Provides user-friendly documentation for the CLI.

**Current State:** No help command exists. Users must read source or guess commands.

**Required Changes:**

1. Create `cmd/help.go`:

```go
package cmd

import (
    "flag"
    "fmt"
    "sort"
)

type HelpCommand struct {
    fs *flag.FlagSet
}

func NewHelpCommand() *HelpCommand {
    return &HelpCommand{
        fs: flag.NewFlagSet("help", flag.ExitOnError),
    }
}

func (hc *HelpCommand) Name() string {
    return "help"
}

func (hc *HelpCommand) Description() string {
    return "Display help information about commands"
}

func (hc *HelpCommand) FlagSet() *flag.FlagSet {
    return hc.fs
}

func (hc *HelpCommand) Execute(args []string) error {
    if err := hc.fs.Parse(args); err != nil {
        return err
    }
    
    remainingArgs := hc.fs.Args()
    
    // If specific command requested, show detailed help
    if len(remainingArgs) > 0 {
        return hc.showCommandHelp(remainingArgs[0])
    }
    
    // Otherwise show general help
    hc.showGeneralHelp()
    return nil
}

func (hc *HelpCommand) showGeneralHelp() {
    fmt.Println("Maajise - Project Initialization Tool")
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println("  maajise <command> [arguments]")
    fmt.Println("  maajise <project-name>         (shorthand for 'init <project-name>')")
    fmt.Println()
    fmt.Println("Available commands:")
    
    // Get all commands and sort alphabetically
    commands := All()
    names := make([]string, 0, len(commands))
    for name := range commands {
        names = append(names, name)
    }
    sort.Strings(names)
    
    // Display each command with description
    for _, name := range names {
        cmd := commands[name]
        fmt.Printf("  %-12s %s\n", name, cmd.Description())
    }
    
    fmt.Println()
    fmt.Println("Use 'maajise help <command>' for more information about a command.")
}

func (hc *HelpCommand) showCommandHelp(cmdName string) error {
    cmd, ok := Get(cmdName)
    if !ok {
        return fmt.Errorf("unknown command: %s", cmdName)
    }
    
    fmt.Printf("Command: %s\n", cmd.Name())
    fmt.Printf("Description: %s\n", cmd.Description())
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Printf("  maajise %s [options]\n", cmd.Name())
    fmt.Println()
    
    // Show available flags
    fs := cmd.FlagSet()
    if fs != nil {
        fmt.Println("Options:")
        fs.PrintDefaults()
    }
    
    return nil
}

func init() {
    Register(NewHelpCommand())
}
```

**Test Requirements:**
- Test HelpCommand implements Command interface
- Test showGeneralHelp() displays all commands
- Test showCommandHelp() with valid command
- Test showCommandHelp() with invalid command returns error
- Test Execute() with no args calls showGeneralHelp()
- Test Execute() with command name calls showCommandHelp()

**Dependencies:** Task 1 complete (Command interface exists)  
**Blocks:** ISSUE-008

---

#### ISSUE-007: Add Tests for HelpCommand

**Severity:** WARNING  
**File:** New file: `cmd/help_test.go`  
**Description:** Comprehensive tests for help command functionality.

**Current State:** No tests for help command

**Required Changes:**

1. Create `cmd/help_test.go`:

```go
package cmd

import (
    "testing"
)

func TestNewHelpCommand(t *testing.T) {
    hc := NewHelpCommand()
    
    if hc.Name() != "help" {
        t.Errorf("Expected name 'help', got '%s'", hc.Name())
    }
    
    if hc.Description() == "" {
        t.Error("Description should not be empty")
    }
    
    if hc.FlagSet() == nil {
        t.Error("FlagSet should not be nil")
    }
}

func TestHelpCommand_Execute_GeneralHelp(t *testing.T) {
    // Clear and setup test registry
    Registry = make(map[string]Command)
    Register(NewHelpCommand())
    Register(NewInitCommand())
    
    hc := NewHelpCommand()
    
    // Execute with no args should show general help
    err := hc.Execute([]string{})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
}

func TestHelpCommand_Execute_CommandHelp(t *testing.T) {
    // Setup test registry
    Registry = make(map[string]Command)
    Register(NewHelpCommand())
    Register(NewInitCommand())
    
    hc := NewHelpCommand()
    
    // Execute with command name
    err := hc.Execute([]string{"init"})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
}

func TestHelpCommand_Execute_InvalidCommand(t *testing.T) {
    Registry = make(map[string]Command)
    Register(NewHelpCommand())
    
    hc := NewHelpCommand()
    
    err := hc.Execute([]string{"nonexistent"})
    if err == nil {
        t.Error("Expected error for invalid command")
    }
}
```

**Test Requirements:**
- All tests pass
- Coverage >70% for help.go

**Dependencies:** ISSUE-006  
**Blocks:** None

---

### TASK 3: Create cmd/version.go

#### ISSUE-008: Create VersionCommand

**Severity:** INFO  
**File:** New file: `cmd/version.go`  
**Description:** Simple version display command. Shows current Maajise version and build information.

**Current State:** No version command exists

**Required Changes:**

1. Create `cmd/version.go`:

```go
package cmd

import (
    "flag"
    "fmt"
)

const Version = "2.0.0-alpha"

type VersionCommand struct {
    fs *flag.FlagSet
}

func NewVersionCommand() *VersionCommand {
    return &VersionCommand{
        fs: flag.NewFlagSet("version", flag.ExitOnError),
    }
}

func (vc *VersionCommand) Name() string {
    return "version"
}

func (vc *VersionCommand) Description() string {
    return "Display version information"
}

func (vc *VersionCommand) FlagSet() *flag.FlagSet {
    return vc.fs
}

func (vc *VersionCommand) Execute(args []string) error {
    if err := vc.fs.Parse(args); err != nil {
        return err
    }
    
    fmt.Printf("Maajise version %s\n", Version)
    return nil
}

func init() {
    Register(NewVersionCommand())
}
```

**Test Requirements:**
- Test VersionCommand implements Command interface
- Test Execute() displays version
- Test version string is valid

**Dependencies:** Task 1 complete  
**Blocks:** None

---

#### ISSUE-009: Add Tests for VersionCommand

**Severity:** WARNING  
**File:** New file: `cmd/version_test.go`  
**Description:** Tests for version command.

**Current State:** No tests for version command

**Required Changes:**

1. Create `cmd/version_test.go`:

```go
package cmd

import (
    "testing"
)

func TestNewVersionCommand(t *testing.T) {
    vc := NewVersionCommand()
    
    if vc.Name() != "version" {
        t.Errorf("Expected name 'version', got '%s'", vc.Name())
    }
    
    if vc.Description() == "" {
        t.Error("Description should not be empty")
    }
    
    if vc.FlagSet() == nil {
        t.Error("FlagSet should not be nil")
    }
}

func TestVersionCommand_Execute(t *testing.T) {
    vc := NewVersionCommand()
    
    err := vc.Execute([]string{})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
}

func TestVersionConstant(t *testing.T) {
    if Version == "" {
        t.Error("Version constant should not be empty")
    }
}
```

**Test Requirements:**
- All tests pass
- Coverage >80% for version.go

**Dependencies:** ISSUE-008  
**Blocks:** None

---

### TASK 4: Refactor main.go

#### ISSUE-010: Slim Down main.go to Dispatcher

**Severity:** CRITICAL  
**File:** `main.go` (modify)  
**Description:** Complete the main.go refactoring to create a thin dispatcher that routes commands. Target: <150 lines.

**Current State:** main.go still contains helper functions and logic (from before Task 1 migration)

**Required Changes:**

1. Update main.go to final form:

```go
package main

import (
    "fmt"
    "os"
    
    "maajise/cmd"
)

func main() {
    // Handle no arguments
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }
    
    // Get command name
    cmdName := os.Args[1]
    cmdArgs := os.Args[2:]
    
    // Special case: --help or -h
    if cmdName == "--help" || cmdName == "-h" {
        cmdName = "help"
        cmdArgs = []string{}
    }
    
    // Special case: --version or -v
    if cmdName == "--version" || cmdName == "-v" {
        cmdName = "version"
        cmdArgs = []string{}
    }
    
    // Convenience: if first arg doesn't match a command, assume "init"
    if _, ok := cmd.Get(cmdName); !ok {
        // Treat as project name for init command
        cmdArgs = append([]string{cmdName}, cmdArgs...)
        cmdName = "init"
    }
    
    // Get command
    command, ok := cmd.Get(cmdName)
    if !ok {
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmdName)
        fmt.Fprintln(os.Stderr)
        printUsage()
        os.Exit(1)
    }
    
    // Execute command
    if err := command.Execute(cmdArgs); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func printUsage() {
    if helpCmd, ok := cmd.Get("help"); ok {
        helpCmd.Execute([]string{})
    } else {
        // Fallback if help command not registered
        fmt.Println("Maajise - Project Initialization Tool")
        fmt.Println()
        fmt.Println("Usage:")
        fmt.Println("  maajise <command> [arguments]")
        fmt.Println("  maajise <project-name>  (shorthand for 'init')")
    }
}
```

2. Remove all remaining functions that were moved to cmd/init.go

3. Verify no dead code remains

4. Verify imports are minimal

**Test Requirements:**
- main.go compiles with no errors
- main.go < 150 lines
- `go build` succeeds
- Binary runs all commands correctly

**Dependencies:** ISSUE-006, ISSUE-008 (help and version commands exist)  
**Blocks:** None

---

#### ISSUE-011: Add Integration Tests for main.go

**Severity:** WARNING  
**File:** New file: `main_test.go`  
**Description:** Integration tests to verify main.go dispatcher works correctly with all commands.

**Current State:** No integration tests for main dispatcher

**Required Changes:**

1. Create `main_test.go`:

```go
package main

import (
    "os"
    "testing"
    
    "maajise/cmd"
)

func TestMain(m *testing.M) {
    // Setup: ensure commands are registered
    cmd.Register(cmd.NewInitCommand())
    cmd.Register(cmd.NewHelpCommand())
    cmd.Register(cmd.NewVersionCommand())
    
    // Run tests
    code := m.Run()
    
    // Cleanup
    os.Exit(code)
}

func TestPrintUsage(t *testing.T) {
    // Just verify it doesn't panic
    printUsage()
}

// Note: Testing main() directly is difficult
// These tests verify the supporting functions work
// Full integration testing done manually
```

**Test Requirements:**
- Tests pass
- No panics during execution

**Dependencies:** ISSUE-010  
**Blocks:** None

---

### TASK 5: Wire Up Internal Packages

#### ISSUE-012: Audit internal/ Package Usage

**Severity:** INFO  
**File:** Documentation  
**Description:** Survey current internal/ packages and determine what exists vs what needs to be created or integrated.

**Current State:** Unknown which internal packages are implemented and which are stubs

**Required Changes:**

1. List contents of internal/ directory:

```bash
ls -la internal/
```

2. For each subdirectory, check:
   - Does it have actual implementation?
   - Does it have tests?
   - Is it currently used anywhere?

3. Create audit document: `internal/AUDIT.md`:

```markdown
# Internal Packages Audit

## internal/git
Status: [EXISTS/STUB/MISSING]
Functions: [list main functions]
Tests: [YES/NO/PARTIAL]
Current usage: [list files using it]

## internal/beads
Status: [EXISTS/STUB/MISSING]
Functions: [list main functions]
Tests: [YES/NO/PARTIAL]
Current usage: [list files using it]

## internal/ui
Status: [EXISTS/STUB/MISSING]
Functions: [list main functions]
Tests: [YES/NO/PARTIAL]
Current usage: [list files using it]

## internal/validate
Status: [EXISTS/STUB/MISSING]
Functions: [list main functions]
Tests: [YES/NO/PARTIAL]
Current usage: [list files using it]

## Recommendations
[What needs to be built, what can be integrated now]
```

**Test Requirements:**
- Document is accurate
- All internal/ subdirectories surveyed

**Dependencies:** None  
**Blocks:** ISSUE-013, ISSUE-014, ISSUE-015

---

#### ISSUE-013: Integrate internal/git Package

**Severity:** CRITICAL  
**File:** `cmd/init.go` (modify), `internal/git/` (verify/implement)  
**Description:** Replace direct exec.Command("git"...) calls in cmd/init.go with internal/git package functions. Centralizes Git operations with proper error handling.

**Current State:** cmd/init.go uses exec.Command directly for Git operations

**Required Changes:**

1. Verify `internal/git/git.go` exists and has these functions (create if missing):

```go
package git

import (
    "fmt"
    "os/exec"
    "strings"
)

// Init initializes a new Git repository with the specified branch
func Init(branch string) error {
    if err := exec.Command("git", "init").Run(); err != nil {
        return fmt.Errorf("git init failed: %w", err)
    }
    
    if err := exec.Command("git", "branch", "-M", branch).Run(); err != nil {
        return fmt.Errorf("failed to set branch to %s: %w", branch, err)
    }
    
    return nil
}

// IsRepo checks if current directory is a Git repository
func IsRepo() bool {
    err := exec.Command("git", "rev-parse", "--git-dir").Run()
    return err == nil
}

// AddRemote adds a remote with the given name and URL
func AddRemote(name, url string) error {
    if err := exec.Command("git", "remote", "add", name, url).Run(); err != nil {
        return fmt.Errorf("failed to add remote %s: %w", name, err)
    }
    return nil
}

// GetRemotes returns list of configured remotes
func GetRemotes() ([]string, error) {
    out, err := exec.Command("git", "remote").Output()
    if err != nil {
        return nil, fmt.Errorf("failed to list remotes: %w", err)
    }
    
    remotes := strings.Split(strings.TrimSpace(string(out)), "\n")
    if len(remotes) == 1 && remotes[0] == "" {
        return []string{}, nil
    }
    
    return remotes, nil
}
```

2. Update `cmd/init.go` to use internal/git:

```go
import (
    "maajise/internal/git"
)

func initializeGit(branch string) error {
    if git.IsRepo() {
        return fmt.Errorf("directory is already a git repository")
    }
    
    if err := git.Init(branch); err != nil {
        return err
    }
    
    return nil
}

func setupGitRemote(remoteURL string) error {
    // Validate URL first (still using exec for now)
    // TODO: Move to internal/validate in future
    
    if err := git.AddRemote("origin", remoteURL); err != nil {
        return err
    }
    
    return nil
}
```

3. Remove direct exec.Command("git"...) calls from init.go

**Test Requirements:**
- Create `internal/git/git_test.go` with tests for each function
- Test coverage >70%
- Integration test: full init creates valid Git repo

**Dependencies:** ISSUE-012 (audit complete)  
**Blocks:** None

---

#### ISSUE-014: Integrate internal/beads Package

**Severity:** CRITICAL  
**File:** `cmd/init.go` (modify), `internal/beads/` (verify/implement)  
**Description:** Replace direct exec.Command("bd"...) calls with internal/beads package functions.

**Current State:** cmd/init.go uses exec.Command directly for Beads operations

**Required Changes:**

1. Verify `internal/beads/beads.go` exists and has these functions (create if missing):

```go
package beads

import (
    "fmt"
    "os/exec"
)

// Init initializes a new Beads issue tracker
func Init() error {
    if err := exec.Command("bd", "init").Run(); err != nil {
        return fmt.Errorf("bd init failed: %w", err)
    }
    return nil
}

// IsInitialized checks if Beads is initialized in current directory
func IsInitialized() bool {
    err := exec.Command("bd", "list").Run()
    return err == nil
}

// IsInstalled checks if bd command is available
func IsInstalled() bool {
    _, err := exec.LookPath("bd")
    return err == nil
}
```

2. Update `cmd/init.go` to use internal/beads:

```go
import (
    "maajise/internal/beads"
)

func initializeBeads() error {
    if !beads.IsInstalled() {
        return fmt.Errorf("bd command not found - install Beads first")
    }
    
    if beads.IsInitialized() {
        return fmt.Errorf("directory already has Beads initialized")
    }
    
    if err := beads.Init(); err != nil {
        return err
    }
    
    return nil
}
```

3. Remove direct exec.Command("bd"...) calls from init.go

**Test Requirements:**
- Create `internal/beads/beads_test.go`
- Test coverage >70%
- Mock exec calls or skip if bd not installed

**Dependencies:** ISSUE-012 (audit complete)  
**Blocks:** None

---

#### ISSUE-015: Integrate internal/ui Package

**Severity:** WARNING  
**File:** `cmd/init.go`, `cmd/help.go`, `main.go` (modify), `internal/ui/` (verify/implement)  
**Description:** Centralize output formatting with internal/ui package. Enables consistent colored output, progress indicators, etc.

**Current State:** Commands use fmt.Printf directly

**Required Changes:**

1. Verify `internal/ui/ui.go` exists with these functions (create if missing):

```go
package ui

import (
    "fmt"
)

// Success prints a success message with checkmark
func Success(msg string) {
    fmt.Printf("✓ %s\n", msg)
}

// Error prints an error message
func Error(msg string) {
    fmt.Printf("✗ %s\n", msg)
}

// Info prints an informational message
func Info(msg string) {
    fmt.Printf("ℹ %s\n", msg)
}

// Header prints a section header
func Header(msg string) {
    fmt.Printf("\n=== %s ===\n\n", msg)
}

// TODO: Add color support in future
// For now, keep it simple
```

2. Update cmd/init.go to use internal/ui:

```go
import (
    "maajise/internal/ui"
)

func (ic *InitCommand) runInit() error {
    ui.Info(fmt.Sprintf("Initializing project: %s", ic.config.ProjectName))
    
    if err := createProjectDirectory(ic.config.ProjectName); err != nil {
        return err
    }
    ui.Success("Project directory created")
    
    // ... other steps with ui.Success() calls
    
    ui.Success(fmt.Sprintf("Project '%s' initialized successfully", ic.config.ProjectName))
    return nil
}
```

3. Update cmd/help.go to use internal/ui for headers

4. Update main.go error output to use ui.Error()

**Test Requirements:**
- Create `internal/ui/ui_test.go`
- Test each function produces output
- Coverage >60%

**Dependencies:** ISSUE-012 (audit complete)  
**Blocks:** None

---

#### ISSUE-016: Create internal/config Package

**Severity:** INFO  
**File:** New: `internal/config/config.go`, Update: `cmd/config.go`  
**Description:** Move Config to internal/config for better organization and future extensibility.

**Current State:** Config lives in cmd/config.go

**Required Changes:**

1. Create `internal/config/config.go`:

```go
package config

// Config holds initialization configuration for project setup
type Config struct {
    ProjectName    string
    UseBeads       bool
    NoGit          bool
    NoOverwrite    bool
    SkipUBS        bool
    MainBranch     string
    DefaultRemote  string
}

// New returns a Config with sensible defaults
func New() Config {
    return Config{
        UseBeads:      true,
        NoGit:         false,
        NoOverwrite:   false,
        SkipUBS:       false,
        MainBranch:    "main",
        DefaultRemote: "",
    }
}
```

2. Update `cmd/config.go` to re-export for backwards compatibility:

```go
package cmd

import "maajise/internal/config"

// Config is re-exported from internal/config for compatibility
type Config = config.Config

// DefaultConfig returns a new config with defaults
func DefaultConfig() Config {
    return config.New()
}
```

3. Update all imports if needed

**Test Requirements:**
- Verify all code still compiles
- No functional changes

**Dependencies:** None  
**Blocks:** None

---

#### ISSUE-017: Final Integration Testing

**Severity:** WARNING  
**File:** Documentation  
**Description:** Comprehensive manual testing of entire Phase 1 to verify all components work together.

**Current State:** Individual issues tested, but not end-to-end integration

**Required Changes:**

1. Create test script: `test-phase1.sh`:

```bash
#!/bin/bash
set -e

echo "=== Phase 1 Integration Tests ==="

# Clean up previous test runs
rm -rf test-project-* 2>/dev/null || true

echo ""
echo "Test 1: Basic init"
./maajise init test-project-1
cd test-project-1
if [ ! -d .git ]; then
    echo "FAIL: .git directory not created"
    exit 1
fi
cd ..
echo "✓ Test 1 passed"

echo ""
echo "Test 2: Convenience shorthand"
./maajise test-project-2
if [ ! -d test-project-2/.git ]; then
    echo "FAIL: convenience shorthand failed"
    exit 1
fi
echo "✓ Test 2 passed"

echo ""
echo "Test 3: Help command"
./maajise help
./maajise help init
echo "✓ Test 3 passed"

echo ""
echo "Test 4: Version command"
./maajise version
echo "✓ Test 4 passed"

echo ""
echo "Test 5: Flags work"
./maajise init test-project-3 --no-git
if [ -d test-project-3/.git ]; then
    echo "FAIL: --no-git didn't work"
    exit 1
fi
echo "✓ Test 5 passed"

echo ""
echo "Test 6: Global --help"
./maajise --help
echo "✓ Test 6 passed"

echo ""
echo "Test 7: Global --version"
./maajise --version
echo "✓ Test 7 passed"

echo ""
echo "=== All Phase 1 Integration Tests Passed ==="
```

2. Run tests:

```bash
chmod +x test-phase1.sh
go build
./test-phase1.sh
```

3. Document any failures and fix before closing issue

**Test Requirements:**
- All 7 integration tests pass
- No crashes or panics
- Error messages are clear

**Dependencies:** All previous issues (ISSUE-006 through ISSUE-016)  
**Blocks:** None

---

## Implementation Order

Follow this sequence to respect dependencies:

### Task 2: Help Command (30-40 minutes)
1. **ISSUE-006** - Create HelpCommand (25 min)
2. **ISSUE-007** - Add tests for HelpCommand (15 min)

### Task 3: Version Command (15-20 minutes)
3. **ISSUE-008** - Create VersionCommand (10 min)
4. **ISSUE-009** - Add tests for VersionCommand (10 min)

### Task 4: Main.go Refactor (30-40 minutes)
5. **ISSUE-010** - Slim down main.go (25 min)
6. **ISSUE-011** - Add integration tests for main.go (15 min)

### Task 5: Internal Packages (60-90 minutes)
7. **ISSUE-012** - Audit internal/ packages (15 min)
8. **ISSUE-013** - Integrate internal/git (25 min)
9. **ISSUE-014** - Integrate internal/beads (20 min)
10. **ISSUE-015** - Integrate internal/ui (20 min)
11. **ISSUE-016** - Create internal/config (15 min)
12. **ISSUE-017** - Final integration testing (20 min)

**Total Estimated Time:** 2.5-3.5 hours

---

## CREATE ISSUES NOW ⚠️

**Run these commands BEFORE starting implementation:**

```bash
bd create --title "ISSUE-006: Create HelpCommand Structure" --description "Create help command that displays all registered commands. See ISSUE-006 in spec."

bd create --title "ISSUE-007: Add Tests for HelpCommand" --description "Comprehensive tests for help command functionality. See ISSUE-007 in spec."

bd create --title "ISSUE-008: Create VersionCommand" --description "Simple version display command with build info. See ISSUE-008 in spec."

bd create --title "ISSUE-009: Add Tests for VersionCommand" --description "Tests for version command. See ISSUE-009 in spec."

bd create --title "ISSUE-010: Slim Down main.go to Dispatcher" --description "Complete main.go refactor to thin dispatcher (<150 lines). See ISSUE-010 in spec."

bd create --title "ISSUE-011: Add Integration Tests for main.go" --description "Integration tests for main dispatcher. See ISSUE-011 in spec."

bd create --title "ISSUE-012: Audit internal/ Package Usage" --description "Survey internal packages to determine what exists and what needs work. See ISSUE-012 in spec."

bd create --title "ISSUE-013: Integrate internal/git Package" --description "Replace exec.Command git calls with internal/git functions. See ISSUE-013 in spec."

bd create --title "ISSUE-014: Integrate internal/beads Package" --description "Replace exec.Command bd calls with internal/beads functions. See ISSUE-014 in spec."

bd create --title "ISSUE-015: Integrate internal/ui Package" --description "Centralize output formatting with internal/ui. See ISSUE-015 in spec."

bd create --title "ISSUE-016: Create internal/config Package" --description "Move Config to internal/config for better organization. See ISSUE-016 in spec."

bd create --title "ISSUE-017: Final Integration Testing" --description "End-to-end testing of all Phase 1 components. See ISSUE-017 in spec."
```

**Verify all issues created:**

```bash
bd list
# Should show 12 new issues (plus 5 from Task 1 = 17 total)
```

---

## Testing Strategy

### Unit Tests

After each issue:

```bash
go test ./cmd/...
go test ./internal/...
```

### Coverage Targets

- cmd/help.go: >70%
- cmd/version.go: >80%
- internal/git: >70%
- internal/beads: >70%
- internal/ui: >60%
- Overall: >65%

### Integration Tests

After ISSUE-017:

```bash
./test-phase1.sh
```

### Manual Verification

Test each command:

```bash
maajise help
maajise help init
maajise version
maajise init test-proj
maajise test-proj-2
maajise init --help
```

---

## Success Criteria

### Issue Tracking
- [ ] All 12 Beads issues created (verify with `bd list`)
- [ ] Each issue marked complete as work finishes
- [ ] All issues closed at end

### Code Structure
- [ ] cmd/help.go exists and works
- [ ] cmd/version.go exists and works
- [ ] main.go < 150 lines (target < 100)
- [ ] internal/git package integrated
- [ ] internal/beads package integrated
- [ ] internal/ui package integrated
- [ ] internal/config package created

### Functionality
- [ ] `maajise help` lists all commands
- [ ] `maajise help init` shows init details
- [ ] `maajise version` displays version
- [ ] `maajise --help` works
- [ ] `maajise --version` works
- [ ] All init functionality still works
- [ ] Convenience shorthand still works

### Testing
- [ ] All unit tests pass: `go test ./...`
- [ ] Test coverage >65%: `go test -cover ./...`
- [ ] Integration test script passes
- [ ] No regressions

### Build & Quality
- [ ] Clean build: `go build`
- [ ] No compiler warnings
- [ ] UBS scan passes (0 critical issues)
- [ ] Binary size reasonable (<10MB)

### Phase 1 Complete ✅
- [ ] Architecture fully modular
- [ ] main.go is thin dispatcher
- [ ] All subcommands follow Command interface
- [ ] Internal packages in use
- [ ] Ready for Phase 2 (templates)

---

## Notes for Implementation

### Command Pattern

All commands follow same structure:
1. Implement Command interface
2. Define flags in constructor
3. Parse flags in Execute()
4. Call worker functions
5. Return errors clearly

### Internal Package Design

Keep internal packages focused:
- internal/git: Git operations only
- internal/beads: Beads operations only
- internal/ui: Output formatting only
- internal/config: Configuration structs only

### Error Handling

Chain errors with context:

```go
if err := git.Init("main"); err != nil {
    return fmt.Errorf("failed to initialize git: %w", err)
}
```

### Testing Strategy

For packages wrapping exec.Command:
- Test IsInstalled() and similar checks
- Mock exec calls or skip tests if tool missing
- Focus on error handling and edge cases

### Integration Testing

The test-phase1.sh script is crucial:
- Catches integration issues
- Verifies user-facing behavior
- Quick smoke test before release

### Code Organization

After Phase 1:
```
maajise/
├── main.go                 (~80 lines)
├── cmd/
│   ├── command.go         (interface + registry)
│   ├── config.go          (re-export)
│   ├── init.go            (init command)
│   ├── help.go            (help command)
│   ├── version.go         (version command)
│   └── *_test.go          (tests)
├── internal/
│   ├── git/               (git operations)
│   ├── beads/             (beads operations)
│   ├── ui/                (output formatting)
│   └── config/            (config struct)
└── test-phase1.sh         (integration tests)
```

---

## Checkpoint: After Each Task

### After Task 2 (Help Command)
- [ ] `maajise help` works
- [ ] `maajise help init` works
- [ ] Tests pass

### After Task 3 (Version Command)
- [ ] `maajise version` works
- [ ] `maajise --version` works
- [ ] Tests pass

### After Task 4 (Main Refactor)
- [ ] main.go < 150 lines
- [ ] All commands still work
- [ ] Convenience shorthand works
- [ ] Error handling clean

### After Task 5 (Internal Packages)
- [ ] No direct exec calls in cmd/
- [ ] All internal packages tested
- [ ] Integration tests pass
- [ ] Ready for Phase 2

---

## Common Issues & Solutions

### Issue: main.go still >150 lines
**Solution:** Move printUsage to cmd/help and call it

### Issue: Tests fail when bd not installed
**Solution:** Skip tests with:
```go
if !beads.IsInstalled() {
    t.Skip("bd not installed")
}
```

### Issue: Import cycle
**Solution:** Keep cmd/ → internal/ direction only, never reverse

### Issue: Integration tests fail
**Solution:** Check build is fresh: `go build` before running tests

---

## Future Work (Out of Scope)

Phase 2 and beyond:
- Template system
- Multi-language support
- Additional commands (update, validate)
- Configuration files
- Custom templates

Phase 1 establishes the architecture that makes all of this possible.

---

## Version History

- v1.0.0 (2025-11-23): Initial specification for Phase 1 Tasks 2-5

---

## References

- Maajise High-Level Roadmap: Phase 1 details
- Task Specification Best Practices: Structure and format guidance
- Phase 1 Task 1 Spec: Foundation this builds on

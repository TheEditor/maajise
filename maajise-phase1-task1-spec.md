# Maajise Phase 1 Task 1 - Create cmd/init.go

## Project Context
**Project:** Maajise CLI Tool  
**Phase:** 1 - Complete Architecture Migration  
**Task:** 1 - Create cmd/init.go  
**Date:** 2025-11-23  
**Current State:** Monolithic main.go (~700 lines), architecture 50% migrated

---

## ⚠️ IMPLEMENTATION WORKFLOW

**CRITICAL FIRST STEP:** Create all Beads issues before implementing

### Execution Steps

**STEP 1:** Read entire specification to understand scope

**STEP 2:** Create all Beads issues (see "CREATE ISSUES NOW" section below)

**STEP 3:** Verify issues created
```bash
bd list
# Should show all 5 issues
```

**STEP 4:** Work through issues following "Implementation Order" section

**STEP 5:** Mark issues complete as you go
```bash
bd close <issue-id>
```

**STEP 6:** Verify completion using Success Criteria checklist

---

## Issues

### ISSUE-001: Create Command Interface

**Severity:** CRITICAL  
**File:** New file: `cmd/command.go`  
**Description:** Define the Command interface that all subcommands will implement. This establishes the contract for command registration and execution in the dispatcher pattern.

**Current State:** No command interface exists. main.go handles everything directly.

**Required Changes:**

1. Create directory: `cmd/`

2. Create `cmd/command.go` with this interface:

```go
package cmd

import "flag"

// Command represents a subcommand that can be executed
type Command interface {
    // Name returns the command name (e.g., "init", "help")
    Name() string
    
    // Description returns a brief description for help text
    Description() string
    
    // FlagSet returns the command's flag set for parsing args
    FlagSet() *flag.FlagSet
    
    // Execute runs the command with parsed flags and remaining args
    Execute(args []string) error
}

// Registry holds all registered commands
var Registry = make(map[string]Command)

// Register adds a command to the registry
func Register(cmd Command) {
    Registry[cmd.Name()] = cmd
}

// Get retrieves a command by name
func Get(name string) (Command, bool) {
    cmd, ok := Registry[name]
    return cmd, ok
}

// All returns all registered commands
func All() map[string]Command {
    return Registry
}
```

**Test Requirements:**
- Test Register() adds commands to Registry
- Test Get() retrieves existing commands
- Test Get() returns false for non-existent commands
- Test All() returns all registered commands

**Dependencies:** None  
**Blocks:** ISSUE-002, ISSUE-003, ISSUE-004

---

### ISSUE-002: Extract Config Struct to cmd Package

**Severity:** CRITICAL  
**File:** `cmd/config.go` (new) and `main.go` (modify)  
**Description:** Move the Config struct from main.go to cmd package so it can be shared by all subcommands. This enables consistent configuration across the modular architecture.

**Current State:** Config struct defined in main.go, not accessible to cmd/ packages

**Required Changes:**

1. Create `cmd/config.go`:

```go
package cmd

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

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
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

2. Update main.go imports to use `"maajise/cmd"`

3. Replace main.go Config references with `cmd.Config`

4. Remove old Config struct definition from main.go

**Test Requirements:**
- Test DefaultConfig() returns expected values
- Test Config fields are correctly typed
- Verify main.go compiles with new import

**Dependencies:** ISSUE-001  
**Blocks:** ISSUE-003

---

### ISSUE-003: Create InitCommand Implementation

**Severity:** CRITICAL  
**File:** New file: `cmd/init.go`  
**Description:** Migrate all initialization logic from main.go to a proper InitCommand that implements the Command interface. This is the core of the Phase 1 Task 1 objective.

**Current State:** All init logic embedded in main() function in main.go

**Required Changes:**

1. Create `cmd/init.go`:

```go
package cmd

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
)

type InitCommand struct {
    fs     *flag.FlagSet
    config Config
}

func NewInitCommand() *InitCommand {
    ic := &InitCommand{
        fs:     flag.NewFlagSet("init", flag.ExitOnError),
        config: DefaultConfig(),
    }
    
    // Define flags
    ic.fs.BoolVar(&ic.config.UseBeads, "use-beads", true, "Initialize with Beads issue tracker")
    ic.fs.BoolVar(&ic.config.UseBeads, "ub", true, "Shorthand for --use-beads")
    ic.fs.BoolVar(&ic.config.NoGit, "no-git", false, "Skip Git initialization")
    ic.fs.BoolVar(&ic.config.NoGit, "ng", false, "Shorthand for --no-git")
    ic.fs.BoolVar(&ic.config.NoOverwrite, "no-overwrite", false, "Don't overwrite existing files")
    ic.fs.BoolVar(&ic.config.NoOverwrite, "no", false, "Shorthand for --no-overwrite")
    ic.fs.BoolVar(&ic.config.SkipUBS, "skip-ubs", false, "Skip UBS configuration")
    ic.fs.BoolVar(&ic.config.SkipUBS, "su", false, "Shorthand for --skip-ubs")
    ic.fs.StringVar(&ic.config.MainBranch, "branch", "main", "Main branch name")
    ic.fs.StringVar(&ic.config.MainBranch, "b", "main", "Shorthand for --branch")
    ic.fs.StringVar(&ic.config.DefaultRemote, "remote", "", "Default remote URL")
    ic.fs.StringVar(&ic.config.DefaultRemote, "r", "", "Shorthand for --remote")
    
    return ic
}

func (ic *InitCommand) Name() string {
    return "init"
}

func (ic *InitCommand) Description() string {
    return "Initialize a new project with Git, Beads, and configuration files"
}

func (ic *InitCommand) FlagSet() *flag.FlagSet {
    return ic.fs
}

func (ic *InitCommand) Execute(args []string) error {
    // Parse flags
    if err := ic.fs.Parse(args); err != nil {
        return err
    }
    
    // Get project name from remaining args
    remainingArgs := ic.fs.Args()
    if len(remainingArgs) == 0 {
        return fmt.Errorf("project name required")
    }
    ic.config.ProjectName = remainingArgs[0]
    
    // Validate project name
    if err := validateProjectName(ic.config.ProjectName); err != nil {
        return err
    }
    
    // Execute initialization steps
    return ic.runInit()
}

func (ic *InitCommand) runInit() error {
    // Create project directory
    if err := createProjectDirectory(ic.config.ProjectName); err != nil {
        return err
    }
    
    // Change to project directory
    if err := os.Chdir(ic.config.ProjectName); err != nil {
        return fmt.Errorf("failed to change directory: %w", err)
    }
    
    // Initialize Git (unless --no-git)
    if !ic.config.NoGit {
        if err := initializeGit(ic.config.MainBranch); err != nil {
            return err
        }
    }
    
    // Initialize Beads (unless disabled)
    if ic.config.UseBeads {
        if err := initializeBeads(); err != nil {
            return err
        }
    }
    
    // Create configuration files
    if err := createConfigFiles(ic.config); err != nil {
        return err
    }
    
    // Setup remote if specified
    if ic.config.DefaultRemote != "" {
        if err := setupGitRemote(ic.config.DefaultRemote); err != nil {
            return err
        }
    }
    
    fmt.Printf("✓ Project '%s' initialized successfully\n", ic.config.ProjectName)
    return nil
}

// Helper functions (to be migrated from main.go)
func validateProjectName(name string) error {
    // TODO: Move validation logic from main.go
    return nil
}

func createProjectDirectory(name string) error {
    // TODO: Move directory creation logic from main.go
    return nil
}

func initializeGit(branch string) error {
    // TODO: Move Git init logic from main.go
    // Should use internal/git package once wired up
    return nil
}

func initializeBeads() error {
    // TODO: Move Beads init logic from main.go
    // Should use internal/beads package once wired up
    return nil
}

func createConfigFiles(config Config) error {
    // TODO: Move file creation logic from main.go
    return nil
}

func setupGitRemote(remoteURL string) error {
    // TODO: Move remote setup logic from main.go
    // Should use internal/git package once wired up
    return nil
}

func init() {
    Register(NewInitCommand())
}
```

2. Identify all functions in main.go used by init logic

3. Copy those functions to cmd/init.go (mark as TODO for refactoring)

4. Update function signatures if needed for new Config location

**Test Requirements:**
- Test InitCommand implements Command interface
- Test flag parsing works correctly
- Test Execute() with valid project name succeeds
- Test Execute() without project name returns error
- Test NewInitCommand() registers itself
- Integration test: full init flow creates expected files

**Dependencies:** ISSUE-001, ISSUE-002  
**Blocks:** ISSUE-004

---

### ISSUE-004: Update main.go to Use InitCommand

**Severity:** CRITICAL  
**File:** `main.go` (modify)  
**Description:** Refactor main.go to dispatch to InitCommand instead of containing init logic. This completes the migration of init functionality out of main.go.

**Current State:** main.go contains all init logic directly in main() function

**Required Changes:**

1. Update main.go structure:

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
    
    // Convenience: if first arg doesn't match a command, assume "init"
    if _, ok := cmd.Get(cmdName); !ok {
        // Treat as project name for init command
        cmdName = "init"
        os.Args = append([]string{os.Args[0], "init"}, os.Args[1:]...)
    }
    
    // Get and execute command
    command, ok := cmd.Get(cmdName)
    if !ok {
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmdName)
        printUsage()
        os.Exit(1)
    }
    
    // Execute with remaining args
    if err := command.Execute(os.Args[2:]); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func printUsage() {
    fmt.Println("Maajise - Project Initialization Tool")
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println("  maajise <command> [arguments]")
    fmt.Println("  maajise <project-name>  (convenience shorthand for 'init')")
    fmt.Println()
    fmt.Println("Available commands:")
    for name, cmd := range cmd.All() {
        fmt.Printf("  %-10s %s\n", name, cmd.Description())
    }
}
```

2. Remove all init-related functions from main.go (now in cmd/init.go)

3. Remove old Config struct (now in cmd/config.go)

4. Verify imports are correct

**Test Requirements:**
- Test `maajise init my-project` works
- Test `maajise my-project` (convenience) works
- Test `maajise unknown-command` shows error and usage
- Test `maajise` (no args) shows usage
- Verify main.go is < 100 lines

**Dependencies:** ISSUE-003  
**Blocks:** None

---

### ISSUE-005: Add Tests for cmd Package

**Severity:** WARNING  
**File:** New file: `cmd/init_test.go`  
**Description:** Create comprehensive tests for the new cmd package to ensure Command interface works correctly and InitCommand behaves as expected.

**Current State:** No tests for cmd package

**Required Changes:**

1. Create `cmd/init_test.go`:

```go
package cmd

import (
    "testing"
)

func TestNewInitCommand(t *testing.T) {
    ic := NewInitCommand()
    
    if ic.Name() != "init" {
        t.Errorf("Expected name 'init', got '%s'", ic.Name())
    }
    
    if ic.Description() == "" {
        t.Error("Description should not be empty")
    }
    
    if ic.FlagSet() == nil {
        t.Error("FlagSet should not be nil")
    }
}

func TestInitCommand_Execute_NoProjectName(t *testing.T) {
    ic := NewInitCommand()
    
    err := ic.Execute([]string{})
    if err == nil {
        t.Error("Expected error when no project name provided")
    }
    
    if err.Error() != "project name required" {
        t.Errorf("Expected 'project name required' error, got: %v", err)
    }
}

func TestInitCommand_Execute_WithProjectName(t *testing.T) {
    // TODO: This needs to be an integration test
    // For now, just verify parsing works
    ic := NewInitCommand()
    
    // Don't actually execute, just test flag parsing
    err := ic.FlagSet().Parse([]string{"--no-git", "test-project"})
    if err != nil {
        t.Errorf("Flag parsing failed: %v", err)
    }
    
    if !ic.config.NoGit {
        t.Error("Expected NoGit to be true")
    }
}

func TestCommandRegistry(t *testing.T) {
    // Clear registry for test
    Registry = make(map[string]Command)
    
    ic := NewInitCommand()
    Register(ic)
    
    cmd, ok := Get("init")
    if !ok {
        t.Error("Expected to find 'init' command in registry")
    }
    
    if cmd.Name() != "init" {
        t.Errorf("Expected retrieved command name 'init', got '%s'", cmd.Name())
    }
    
    _, ok = Get("nonexistent")
    if ok {
        t.Error("Expected false when getting nonexistent command")
    }
}

func TestAll(t *testing.T) {
    // Clear registry for test
    Registry = make(map[string]Command)
    
    ic := NewInitCommand()
    Register(ic)
    
    all := All()
    if len(all) != 1 {
        t.Errorf("Expected 1 command in registry, got %d", len(all))
    }
    
    if _, ok := all["init"]; !ok {
        t.Error("Expected 'init' command in All()")
    }
}

func TestDefaultConfig(t *testing.T) {
    config := DefaultConfig()
    
    if config.UseBeads != true {
        t.Error("Expected UseBeads to default to true")
    }
    
    if config.NoGit != false {
        t.Error("Expected NoGit to default to false")
    }
    
    if config.MainBranch != "main" {
        t.Errorf("Expected MainBranch 'main', got '%s'", config.MainBranch)
    }
}
```

2. Run tests and verify coverage:

```bash
go test ./cmd/...
go test -cover ./cmd/...
```

**Test Requirements:**
- All tests pass
- Test coverage > 60% for cmd package
- No test panics or hangs

**Dependencies:** ISSUE-003, ISSUE-004  
**Blocks:** None

---

## Implementation Order

Follow this sequence to respect dependencies:

1. **ISSUE-001** - Create Command Interface
   - No dependencies
   - Foundation for everything else
   - Quick win (~15 minutes)

2. **ISSUE-002** - Extract Config Struct
   - Depends on cmd/ directory existing
   - Needed before InitCommand
   - Quick win (~10 minutes)

3. **ISSUE-003** - Create InitCommand
   - Depends on Command interface and Config
   - Core of the work (~30-45 minutes)
   - Most complex issue

4. **ISSUE-004** - Update main.go
   - Depends on InitCommand existing
   - Simpler than it looks (~20 minutes)
   - Satisfying to see main.go shrink

5. **ISSUE-005** - Add Tests
   - Depends on everything else
   - Can be done incrementally (~20 minutes)
   - Ensures quality

**Total Estimated Time:** 1-2 hours

---

## CREATE ISSUES NOW ⚠️

**Run these commands BEFORE starting implementation:**

```bash
bd create --title "ISSUE-001: Create Command Interface" --description "Define Command interface and registry for subcommand architecture. See ISSUE-001 in spec."

bd create --title "ISSUE-002: Extract Config Struct to cmd Package" --description "Move Config from main.go to cmd/config.go for shared access. See ISSUE-002 in spec."

bd create --title "ISSUE-003: Create InitCommand Implementation" --description "Migrate initialization logic from main.go to cmd/init.go implementing Command interface. See ISSUE-003 in spec."

bd create --title "ISSUE-004: Update main.go to Use InitCommand" --description "Refactor main.go to dispatch to InitCommand instead of containing init logic. See ISSUE-004 in spec."

bd create --title "ISSUE-005: Add Tests for cmd Package" --description "Create comprehensive tests for Command interface and InitCommand. See ISSUE-005 in spec."
```

**Verify all issues created:**

```bash
bd list
# Should show all 5 issues
```

---

## Testing Strategy

### Unit Tests

Run after each issue completion:

```bash
go test ./cmd/...
```

### Coverage Targets

- cmd/command.go: >80% (simple interface)
- cmd/config.go: >80% (simple struct)
- cmd/init.go: >60% (complex logic, some integration)
- Overall cmd package: >65%

### Integration Tests

After ISSUE-004 completion, test full workflows:

```bash
# Test basic init
maajise init test-project-1

# Test convenience shorthand
maajise test-project-2

# Test with flags
maajise init test-project-3 --no-git --skip-ubs

# Test help text
maajise init --help
```

### Manual Verification

After completion:
1. Build binary: `go build`
2. Run binary: `./maajise --help`
3. Initialize test project
4. Verify files created
5. Check UBS scan passes

---

## Success Criteria

### Issue Tracking
- [x] All 5 Beads issues created (verify with `bd list`)
- [ ] Each issue marked complete as work finishes
- [ ] All issues closed at end

### Code Structure
- [ ] `cmd/command.go` exists with Command interface
- [ ] `cmd/config.go` exists with Config struct
- [ ] `cmd/init.go` exists with InitCommand
- [ ] main.go < 100 lines
- [ ] main.go only imports and dispatches

### Functionality
- [ ] `maajise init my-project` works
- [ ] `maajise my-project` (convenience) works
- [ ] All original flags work (--no-git, --use-beads, etc.)
- [ ] Project initialization completes successfully
- [ ] Files created correctly

### Testing
- [ ] All unit tests pass: `go test ./cmd/...`
- [ ] Test coverage >65%: `go test -cover ./cmd/...`
- [ ] Integration tests pass manually
- [ ] No regressions in existing functionality

### Build & Scan
- [ ] Clean build: `go build`
- [ ] No compiler warnings
- [ ] UBS scan passes (if available)
- [ ] Binary runs correctly

---

## Notes for Implementation

### Code Style

Follow existing Maajise conventions:
- Use fmt.Errorf for error wrapping
- Descriptive variable names
- Comments on exported functions
- Keep functions focused and small

### Error Handling

Be explicit with errors:

```go
// Good
if err := initializeGit(); err != nil {
    return fmt.Errorf("git initialization failed: %w", err)
}

// Avoid
initializeGit()  // Ignoring error
```

### Flag Design

Maintain existing flag patterns:
- Long form: `--no-git`
- Short form: `-ng`
- Boolean defaults
- String flags for names/URLs

### Migration Strategy

For ISSUE-003, you can:
1. Copy functions from main.go verbatim initially
2. Mark them with `// TODO: refactor to use internal/git`
3. They'll be cleaned up in Phase 1 Task 5

This allows splitting the work: structure first, internal package integration later.

### Testing Approach

Start with simple tests:
- Test interface implementation
- Test flag parsing
- Test error cases

Save complex integration tests for after ISSUE-004 when everything is wired up.

---

## Future Work (Out of Scope)

These are Phase 1 Task 5 or later:

- Wire up internal/git package
- Wire up internal/beads package
- Wire up internal/ui package
- Move Config to internal/config
- Extract file creation to templates package

For now, keep the existing exec.Command calls and logic. We're just moving code, not refactoring it yet.

---

## Version History

- v1.0.0 (2025-11-23): Initial task specification for Phase 1 Task 1

---

## References

- Maajise High-Level Roadmap: Phase 1 details
- Task Specification Best Practices: Structure and format guidance
- Go Command Pattern: Standard subcommand architecture

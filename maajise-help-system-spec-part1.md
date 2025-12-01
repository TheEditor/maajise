# Maajise Help System Enhancement - Task Specification (Part 1)

## ⚠️ READ FIRST

**BEFORE STARTING ANY WORK:**
1. Read `AGENTS.md` in the project root for critical project context and Beads workflow
2. Read `/mnt/skills/user/bd-issue-tracking/SKILL.md` for issue tracking best practices
3. Then read this entire specification

---

## Project Context

**Project:** Maajise v2.0 (Go CLI Tool)
**Current State:** Global `maajise help` exists, no per-command help
**Date:** 2025-11-30
**Goal:** Comprehensive help system for all commands and flags

---

## ⚠️ IMPLEMENTATION WORKFLOW

**CRITICAL FIRST STEP:** Create all issues before implementing

**STEP 1:** Read this entire spec to understand scope

**STEP 2:** Create all Beads issues (see "CREATE ISSUES NOW" section below)

**STEP 3:** Verify issues created
```bash
bd ready --json
# Should show at least 1 ready issue (ISSUE-000 equivalent)
```

**STEP 4:** Implement in dependency order (see "Implementation Order" section below)

**STEP 5:** Mark complete as you go
```bash
# Use actual ID returned from bd create, not "ISSUE-000"
bd close <actual-issue-id> --reason "Implemented and tested"
```

---

## Help Text Style Guide

### Writing Conventions

**Imperative mood:** "Initialize a repository" (not "Initializes a repository")

**Present tense:** For what the command does

**Structure:**
- First paragraph = one-line summary (shown in command list)
- Subsequent paragraphs = details (2-3 sentences)
- Important notes/warnings if applicable

### Examples with Output

Show representative output for commands that create visible artifacts:

```
  # Basic initialization
  maajise init my-project
  
  Output:
    Creating project: my-project
    ✓ Created directory structure
    ✓ Initialized Git repository
    ✓ Initialized Beads issue tracker
```

### Flag Documentation

Document flags through examples (not separate flag list). Each major flag should appear in at least one example with contextual explanation.

**Rationale:** Learning by example is more effective than dry flag lists. Matches conventions of popular CLI tools (git, docker).

---

## Issues

### ISSUE-000: Rename Execute() to Run() in Command Interface

**Severity:** INFO (refactoring for convention)
**File:** `cmd/command.go`, all command files, `main.go`
**Description:** Align with Go CLI conventions (Cobra, mitchellh/cli) by renaming Execute() method to Run(). This is a simple refactor to match standard Go patterns.

**Current Code:**
```go
type Command interface {
    Name() string
    Description() string
    Execute(args []string) error
}
```

**Required Changes:**

1. Update interface in cmd/command.go:
```go
type Command interface {
    Name() string
    Description() string
    Run(args []string) error
}
```

2. Rename method in all command files:
   - cmd/init.go: `func (c *InitCommand) Execute(args []string) error` → `Run`
   - cmd/help.go: `func (h *HelpCommand) Execute(args []string) error` → `Run`
   - cmd/version.go: `func (v *VersionCommand) Execute(args []string) error` → `Run`
   - Any other existing command files

3. Update main.go (line ~57):
```go
// Change:
if err := command.Execute(os.Args[2:]); err != nil {

// To:
if err := command.Run(os.Args[2:]); err != nil {
```

4. Update all help command calls:
```go
// Find all instances of:
helpCmd.Execute(args)

// Change to:
helpCmd.Run(args)
```

**Test Requirements:**
- All existing tests still pass
- Build succeeds with no errors
- All commands still work: `maajise init`, `maajise help`, `maajise version`

**Dependencies:** None
**Blocks:** ISSUE-001 (interface must be stable first)

---

### ISSUE-001: Add Command Registry and Help Methods

**Severity:** INFO (foundational)
**File:** `cmd/command.go` (existing)
**Description:** Extend Command interface to support detailed help output and add registry functions for command discovery. Each command needs to provide usage, description, and examples.

**Current Code:**
```go
type Command interface {
    Name() string
    Description() string
    Run(args []string) error
}
```

**Required Changes:**

1. Add registry to cmd/command.go (after interface definition):
```go
var registry = make(map[string]Command)

// Register adds a command to the registry
func Register(cmd Command) {
    registry[cmd.Name()] = cmd
}

// Get retrieves a command by name
func Get(name string) (Command, bool) {
    cmd, exists := registry[name]
    return cmd, exists
}

// All returns all registered commands
func All() []Command {
    cmds := make([]Command, 0, len(registry))
    for _, cmd := range registry {
        cmds = append(cmds, cmd)
    }
    return cmds
}
```

2. Add help methods to Command interface:
```go
type Command interface {
    Name() string
    Description() string
    Usage() string           // Single-line syntax
    LongDescription() string // Multi-line detailed explanation
    Examples() string        // Usage examples with output
    Run(args []string) error
}
```

3. Update all existing commands to register in init():
```go
// cmd/init.go
func init() {
    Register(&InitCommand{})
}

// cmd/help.go
func init() {
    Register(&HelpCommand{})
}

// cmd/version.go
func init() {
    Register(&VersionCommand{})
}
```

4. Implement help methods for each command following style guide.

**Example Implementation (VersionCommand):**
```go
func (v *VersionCommand) Usage() string {
    return "maajise version"
}

func (v *VersionCommand) LongDescription() string {
    return `Display version information for Maajise.

Shows the current version number and build information.`
}

func (v *VersionCommand) Examples() string {
    return `  # Show version
  maajise version
  
  Output:
    maajise v2.0.0`
}
```

**Example Implementation (InitCommand - partial):**
```go
func (c *InitCommand) Usage() string {
    return "maajise init [--in-place] [--no-overwrite] [--template=<name>] [<project-name>]"
}

func (c *InitCommand) LongDescription() string {
    return `Initialize a new repository with Git, Beads, and UBS.

Creates project structure, initializes version control, sets up
issue tracking, and configures security scanning.

Available templates: base, typescript, python, rust, php
Use 'maajise templates' to see detailed template descriptions.`
}

func (c *InitCommand) Examples() string {
    return `  # Basic initialization
  maajise init my-project
  
  Output:
    Creating project: my-project
    ✓ Created directory structure
    ✓ Initialized Git repository
    ✓ Initialized Beads issue tracker
    ✓ Created .gitignore
    ✓ Created README.md

  # Initialize in current directory
  maajise init --in-place
      Use current directory instead of creating new one
  
  Output:
    Initializing in current directory
    ✓ Initialized Git repository
    ✓ Initialized Beads issue tracker

  # Use TypeScript template
  maajise init my-ts-app --template=typescript
      Creates TypeScript project with tsconfig.json, package.json
      Available: base, typescript, python, rust, php
  
  # Safe mode - don't overwrite existing files
  maajise init my-project --no-overwrite
      Skips any files that already exist
  
  # Combine flags
  maajise init --in-place --template=python --no-overwrite
      Initialize Python project in current directory safely`
}
```

**Test Requirements:**
- Verify all registered commands implement full interface
- Test that help methods return non-empty strings
- Test registry functions (Register, Get, All)
- Test that init() functions register commands
- Verify examples follow style guide (imperative mood, output shown)

**Dependencies:** ISSUE-000 (Run method must exist)
**Blocks:** ISSUE-002, ISSUE-003, ISSUE-004

---

### ISSUE-002: Implement Command-Specific Help

**Severity:** INFO
**File:** `cmd/help.go` (existing)
**Description:** Support `maajise help <command>` to show detailed help for specific commands. Include proper error handling for invalid commands and argument counts.

**Current Behavior:** Only `maajise help` (no args) works

**Required Changes:**

1. Modify help.go Run() method:
```go
func (h *HelpCommand) Run(args []string) error {
    if len(args) == 0 {
        // Existing: show all commands
        return h.showAllCommands()
    }
    
    if len(args) > 1 {
        return fmt.Errorf("help command takes only one argument\n\nUsage: maajise help [command]")
    }
    
    // New: show specific command help
    cmdName := args[0]
    return h.showCommandHelp(cmdName)
}
```

2. Implement `showCommandHelp()`:
```go
func (h *HelpCommand) showCommandHelp(cmdName string) error {
    cmd, exists := Get(cmdName)
    if !exists {
        return fmt.Errorf("unknown command: %s\n\nRun 'maajise help' to see available commands", cmdName)
    }
    
    fmt.Printf("Usage: %s\n\n", cmd.Usage())
    fmt.Printf("%s\n\n", cmd.LongDescription())
    fmt.Printf("Examples:\n%s\n", cmd.Examples())
    return nil
}
```

3. Update HelpCommand's own help methods:
```go
func (h *HelpCommand) Usage() string {
    return "maajise help [command]"
}

func (h *HelpCommand) LongDescription() string {
    return `Show help information for Maajise commands.

Without arguments, displays a list of all available commands.
With a command name, displays detailed help for that specific command.`
}

func (h *HelpCommand) Examples() string {
    return `  # Show all commands
  maajise help
  
  # Show help for specific command
  maajise help init
  
  Output:
    Usage: maajise init [--in-place] [--template=<name>] [<project-name>]
    
    Initialize a new repository with Git, Beads, and UBS.
    ...`
}
```

**Error Handling Rules:**
- Zero args: Show all commands (existing behavior)
- One arg: Show specific command help
- Two+ args: Error "help takes only one argument. Usage: maajise help [command]"
- Invalid command: Error "unknown command: X. Run 'maajise help' to see available commands"

**Test Requirements:**
- Test `maajise help init` shows init help
- Test `maajise help invalid` shows error with hint
- Test `maajise help help` shows help about help
- Test `maajise help init extra` shows error about too many args
- Test output formatting is clean and readable

**Dependencies:** ISSUE-001 (needs help methods and registry)
**Blocks:** None

---

### ISSUE-003: Add --help Flag Support

**Severity:** INFO
**File:** `main.go`
**Description:** Support `--help` and `-h` flags for each command. Should work as `maajise init --help` or `maajise --help`. Help flag takes precedence over all other flags.

**Current State:** Global --help works (line 23-26), but not per-command

**Required Changes:**

Add flag detection in main.go after command lookup (around line 50, before command execution):

```go
// After: command, ok := cmd.Get(strings.ToLower(cmdName))
// After: handling fallback to init
// BEFORE: if err := command.Run(os.Args[2:]); err != nil

// Check for --help flag in remaining args
for _, arg := range os.Args[2:] {
    if arg == "--help" || arg == "-h" {
        if helpCmd, ok := cmd.Get("help"); ok {
            helpCmd.Run([]string{cmdName})
            os.Exit(0)
        }
        break
    }
}
```

**Exact insertion point in main.go:**
```go
// Line ~48-52: After command lookup and init fallback
command, ok := cmd.Get(strings.ToLower(cmdName))
if !ok {
    // ... existing fallback to init ...
}

// NEW CODE HERE (before command execution):
// Check for --help flag in remaining args
for _, arg := range os.Args[2:] {
    if arg == "--help" || arg == "-h" {
        if helpCmd, ok := cmd.Get("help"); ok {
            helpCmd.Run([]string{cmdName})
            os.Exit(0)
        }
        break
    }
}

// Execute the command with remaining args (existing line ~57)
if err := command.Run(os.Args[2:]); err != nil {
```

**Help Flag Behavior Rules:**
- **Precedence:** --help/-h takes priority, ignore all other flags/args
- **Position:** Scan all args for --help/-h, not just first position
- **Unknown command:** Show "unknown command" error (can't show help for non-existent command)
- **Multiple flags:** `maajise init --in-place --help` shows help (ignores --in-place)

**Test Requirements:**
- Test `maajise --help` works (existing functionality)
- Test `maajise init --help` shows init help
- Test `maajise init --in-place --help` shows help (not init execution)
- Test `maajise init --help --template=ts` shows help (position doesn't matter)
- Test `-h` short form works
- Test `maajise invalidcmd --help` shows "unknown command" error

**Dependencies:** ISSUE-001 (needs help methods), ISSUE-002 (uses help command)
**Blocks:** None

---

## Implementation Order

Follow this sequence to respect dependencies:

**PHASE 1: Foundation (REQUIRED FIRST)**

1. **ISSUE-000** - Rename Execute() to Run()
   - Rationale: Interface change must happen first
   - Time: 0.1 session (simple find-replace)
   - Verification: All tests pass, build succeeds

**PHASE 2: Core Help Infrastructure**

2. **ISSUE-001** - Add registry and help methods
   - Rationale: Foundation for all help functionality
   - Time: 0.5 session
   - Verification: Commands register, help methods compile

3. **ISSUE-002** - Command-specific help
   - Rationale: Uses registry from ISSUE-001
   - Time: 0.25 session
   - Verification: `maajise help init` works

4. **ISSUE-003** - --help flag support
   - Rationale: Builds on ISSUE-001 and ISSUE-002
   - Time: 0.25 session
   - Verification: `maajise init --help` works

**Total Estimated Time:** 1.1 sessions

---

## Testing Strategy

### Automated Unit Tests

Create `cmd/help_test.go`:

```go
package cmd

import (
    "strings"
    "testing"
)

func TestCommandHelpMethods(t *testing.T) {
    commands := []Command{
        &InitCommand{},
        &HelpCommand{},
        &VersionCommand{},
    }
    
    for _, cmd := range commands {
        t.Run(cmd.Name(), func(t *testing.T) {
            // Test all help methods return non-empty
            if cmd.Usage() == "" {
                t.Error("Usage() returned empty string")
            }
            if cmd.LongDescription() == "" {
                t.Error("LongDescription() returned empty string")
            }
            if cmd.Examples() == "" {
                t.Error("Examples() returned empty string")
            }
            
            // Test usage line starts with "maajise"
            if !strings.HasPrefix(cmd.Usage(), "maajise ") {
                t.Errorf("Usage() should start with 'maajise ', got: %s", cmd.Usage())
            }
        })
    }
}

func TestCommandRegistry(t *testing.T) {
    // Test Get function
    if _, exists := Get("init"); !exists {
        t.Error("init command not registered")
    }
    
    if _, exists := Get("help"); !exists {
        t.Error("help command not registered")
    }
    
    if _, exists := Get("version"); !exists {
        t.Error("version command not registered")
    }
    
    // Test invalid command
    if _, exists := Get("notacommand"); exists {
        t.Error("Get() returned true for non-existent command")
    }
    
    // Test All function
    all := All()
    if len(all) < 3 {
        t.Errorf("Expected at least 3 commands, got %d", len(all))
    }
}

func TestHelpCommandWithArgs(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
        errMsg  string
    }{
        {"no args", []string{}, false, ""},
        {"valid command", []string{"init"}, false, ""},
        {"invalid command", []string{"notacommand"}, true, "unknown command"},
        {"too many args", []string{"init", "extra"}, true, "takes only one argument"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h := &HelpCommand{}
            err := h.Run(tt.args)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
                t.Errorf("Error message missing expected text '%s', got: %v", tt.errMsg, err)
            }
        })
    }
}
```

### Integration Tests (Manual)

Run these commands to verify help system:

```bash
# Test global help
maajise help
maajise --help
maajise -h

# Test command-specific help
maajise help init
maajise help help
maajise help version
maajise init --help
maajise version --help

# Test help with invalid commands
maajise help invalidcmd      # Should show error
maajise help init extra      # Should show "too many args" error

# Test help flag in various positions
maajise init --help
maajise init --in-place --help
maajise init --help --template=ts
```

### Coverage Target
- **Help command:** >90%
- **Command help methods:** 100%
- **Registry functions:** 100%
- **Overall:** >70%

---

## CREATE ISSUES NOW ⚠️

**CRITICAL**: Issues must be created BEFORE implementation begins.

### Beads Workflow Reminders

**Before starting:** Read `/mnt/skills/user/bd-issue-tracking/SKILL.md` if you haven't already.

**Key concepts:**
- Beads assigns IDs automatically (you don't choose them)
- Each `bd create` returns JSON like `{"id":"maajise-XXX",...}`
- Use the actual returned IDs in your work, not the ISSUE-000 labels in this spec
- ISSUE-000, ISSUE-001, etc. are just labels for humans to reference tasks in this doc

### Step 1: Create All Issues

Run these commands **in order**. Each command returns JSON with an `id` field - you don't choose the ID, Beads assigns it automatically.

```bash
# Create ISSUE-000
bd create --title "ISSUE-000: Rename Execute() to Run() in Command Interface" \
  --description "Align with Go CLI conventions by renaming Execute() method to Run(). Simple refactor. See ISSUE-000 in spec for details." \
  --json

# Capture the returned ID (looks like: {"id":"maajise-042",...})
# You'll use this actual ID (not "ISSUE-000") in commits and close commands

# Create ISSUE-001
bd create --title "ISSUE-001: Add Command Registry and Help Methods" \
  --description "Extend Command interface with Usage(), LongDescription(), Examples() methods. Add registry functions for command discovery. See ISSUE-001 in spec for implementation." \
  --json

# Create ISSUE-002
bd create --title "ISSUE-002: Implement Command-Specific Help" \
  --description "Support 'maajise help <command>' for detailed command help. Includes error handling for invalid commands. See ISSUE-002 in spec for details." \
  --json

# Create ISSUE-003
bd create --title "ISSUE-003: Add --help Flag Support" \
  --description "Support --help and -h flags globally and per-command. Help flag takes precedence. See ISSUE-003 in spec for main.go changes." \
  --json
```

### Step 2: Verify Issues Created

```bash
# Check all issues were created
bd list --status open

# Verify what's ready to work on (should show ISSUE-000 equivalent first)
bd ready --json
```

**Expected Output:** 
- `bd list` should show 4 open issues
- `bd ready` should show at least 1 issue (ISSUE-000 equivalent has no dependencies)

If `bd ready` shows 0 issues, something went wrong with issue creation.

### Step 3: Begin Implementation

Work through issues in dependency order (see "Implementation Order" section below).

As you complete each issue:
```bash
# Mark complete (use actual ID returned from creation)
bd close <actual-issue-id> --reason "Implemented and tested"

# Verify it's gone from ready list
bd ready --json
```

**Important**: Include the issue ID in your commit messages:
```bash
git commit -m "feat: rename Execute to Run (bd:<actual-issue-id>)"
```

---

## Success Criteria

### Issue Tracking
- [ ] All 4 issues created (actual IDs captured from bd create)
- [ ] Issues tracked with `bd ready` during implementation
- [ ] All issues closed with `bd close <id> --reason "..."`
- [ ] Verify closure: `bd list --status open` shows 0

### Code Changes
- [ ] Execute() renamed to Run() throughout codebase
- [ ] Command interface extended with help methods
- [ ] All commands implement new methods
- [ ] Registry functions working (Register, Get, All)
- [ ] Help command supports args
- [ ] --help flags work everywhere

### Testing
- [ ] All unit tests pass (`go test ./...`)
- [ ] Coverage targets met (>70% overall, >90% help)
- [ ] Manual integration tests complete
- [ ] No regressions in existing commands

### User Experience
- [ ] `maajise help` shows clear command list
- [ ] `maajise help <command>` shows detailed help
- [ ] `maajise <command> --help` works consistently
- [ ] Help output is readable and well-formatted

### Verification
- [ ] Clean build: `go build`
- [ ] All tests pass: `go test ./...`
- [ ] All Beads issues closed: `bd ready` shows 0 ready issues
- [ ] Ready for Part 2 implementation

---

## Next Steps

After completing Part 1:
1. Verify all success criteria met
2. Run full test suite
3. Proceed to **Part 2** for flag documentation and error handling improvements

---

## End of Part 1

**Completed:** Core help system infrastructure
**Remaining:** Flag documentation and usage error improvements (see Part 2)


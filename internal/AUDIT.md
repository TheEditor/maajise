# Internal Packages Audit

## Summary

This audit surveys all internal packages to determine current status, implementation completeness, test coverage, and usage patterns.

**Date:** 2025-11-24
**Status:** All packages exist with implementations; audit complete

---

## Package Status Overview

| Package | Status | Functions | Tests | Current Usage |
|---------|--------|-----------|-------|----------------|
| internal/git | EXISTS | 3 | None | cmd/init.go (direct exec calls) |
| internal/beads | EXISTS | 2 | None | cmd/init.go (direct exec calls) |
| internal/ui | EXISTS | 5 | None | git, beads, help, version |
| internal/validate | EXISTS | 3 | YES | cmd/init.go (URL validation) |
| internal/config | EXISTS | 1 | None | cmd package (partial use) |

---

## Detailed Package Analysis

### internal/git

**Status:** EXISTS - Full implementation

**Location:** `internal/git/git.go` (94 lines)

**Functions:**
1. `Init(repoDir string, verbose bool) error` - Initializes a Git repository
2. `CreateCommit(repoDir string, files []string, message string, verbose bool) error` - Creates initial commit with files
3. `CheckAvailable() error` - Verifies git command is available

**Implementation Notes:**
- Uses `os/exec` to execute git commands
- Supports verbose output (stdout/stderr passthrough)
- Checks for existing .git directory before init
- Provides user-friendly error messages via ui package
- Properly sets working directory for commands

**Tests:** None yet (needed for ISSUE-013)

**Current Usage:**
- cmd/init.go uses `exec.Command("git", ...)` directly in 6 locations
- Does NOT currently import or use internal/git package
- Lines with direct exec: 201, 278, 289, 564, 576, 593, 617, 665

**Integration Complexity:** Medium
- init.go has multiple git operations that vary in arguments
- Will need wrapper functions or better function signatures to cover all cases

---

### internal/beads

**Status:** EXISTS - Basic implementation

**Location:** `internal/beads/beads.go` (46 lines)

**Functions:**
1. `Init(repoDir string, verbose bool) error` - Initializes Beads issue tracker
2. `CheckAvailable() error` - Verifies bd command is available

**Implementation Notes:**
- Uses `os/exec` to run `bd init` command
- Gracefully handles missing bd command (warnings instead of errors)
- Checks for existing .beads directory before init
- Integrates with ui package for user feedback

**Tests:** None yet (needed for ISSUE-014)

**Current Usage:**
- cmd/init.go uses `exec.Command("bd", "init")` directly (line 325)
- Does NOT currently import or use internal/beads package
- Only one location of direct bd usage

**Integration Complexity:** Low
- Straightforward replacement; only one exec.Command call to replace

---

### internal/ui

**Status:** EXISTS - Full implementation with color support

**Location:** `internal/ui/ui.go` (56 lines)

**Functions:**
1. `Error(msg string)` - Red error message with ✗ symbol
2. `Success(msg string)` - Green success message with ✓ symbol
3. `Info(msg string)` - Blue info message with → symbol
4. `Warn(msg string)` - Yellow warning message with ⚠ symbol
5. `Header(title string)` - Formatted header box with blue border
6. `Summary(lines ...string)` - Success summary box with green border

**Implementation Notes:**
- Uses ANSI color codes for terminal output
- Writes errors to stderr, others to stdout
- Provides formatted boxes for headers and summaries
- Already integrated into git and beads packages

**Tests:** None yet (recommended for ISSUE-015)

**Current Usage:**
- internal/git: Calls `ui.Info()`, `ui.Success()`, `ui.Error()`, `ui.Warn()`
- internal/beads: Calls `ui.Info()`, `ui.Success()`, `ui.Warn()`
- help.go: Uses fmt.Printf directly (should use ui.Header)
- version.go: Uses fmt.Printf directly
- main.go: Uses fmt.Fprintf directly (should use ui.Error)
- cmd/init.go: Uses fmt.Printf/Fprintf directly

**Integration Complexity:** Low
- Already partially integrated
- Just need to replace fmt calls with ui calls in cmd files

---

### internal/validate

**Status:** EXISTS - Full implementation with comprehensive tests

**Location:** `internal/validate/validate.go` (77 lines)

**Tests:** `internal/validate/validate_test.go` (309 lines)

**Functions:**
1. `ValidateGitURL(url string) error` - Validates git remote URLs with allowlist approach
2. `SanitizeInput(input string, maxLen int) (string, error)` - Trims and validates input
3. `ContainsShellMetachars(s string) bool` - Detects shell injection attempts

**Implementation Notes:**
- Uses regex pattern matching for git URLs
- Implements shell metacharacter detection (13 characters checked)
- Maximum URL length: 2048 characters
- Security-focused with allowlist approach (https://, git@, ssh:// only)
- Comprehensive test coverage (3 test functions, 19+ test cases per function)

**Tests:** YES - Comprehensive
- TestValidateGitURL: 14 test cases covering valid HTTPS, SSH, and invalid formats
- TestSanitizeInput: 11 test cases covering edge cases
- TestContainsShellMetachars: 14 test cases for safety verification
- **Test Coverage:** Estimated 90%+

**Current Usage:**
- cmd/init.go: Imports validate package (line 13)
- Used for ValidateGitURL() calls in remote URL setup

**Integration Status:** Already integrated ✓

---

### internal/config

**Status:** EXISTS - Minimal implementation

**Location:** `internal/config/config.go` (22 lines)

**Struct:**
```go
type Config struct {
    ProjectName  string
    InPlace      bool
    NoOverwrite  bool
    SkipGit      bool
    SkipBeads    bool
    SkipCommit   bool
    SkipRemote   bool
    Template     string
    Verbose      bool
}
```

**Functions:**
1. `New() *Config` - Creates new Config with defaults

**Implementation Notes:**
- 9 fields covering project initialization options
- Defaults: Template="base", all Skip* flags=false, InPlace=false
- Simple, minimal implementation
- **CONFLICT NOTE:** cmd/config.go defines a different Config struct with 18 fields including GitName, GitEmail, etc.

**Tests:** None

**Current Usage:**
- cmd/config.go does NOT import or use internal/config
- cmd/config.go defines its own Config struct
- cmd/init.go uses cmd/config.Config (not internal/config.Config)

**Integration Status:** BLOCKED - Struct Duplication Issue
- Two different Config structs exist
- Needs resolution in ISSUE-018 (Reconcile Config Struct Duplication)

**Relationship to cmd/config.go:**
```
cmd/config.go (18 fields):
  - ProjectName, Template, Verbose (shared with internal/config)
  - GitName, GitEmail, MainBranch (unique to cmd/config)
  - DefaultRemote, AutoCommit, etc. (unique to cmd/config)

internal/config.go (9 fields):
  - ProjectName, Verbose (shared with cmd/config)
  - InPlace, NoOverwrite, SkipGit, SkipBeads, SkipCommit, SkipRemote, Template
```

---

## Integration Status Summary

### Already Integrated ✓
- **internal/validate:** Used by cmd/init.go for URL validation
- **internal/ui:** Used by internal/git, internal/beads

### Partially Integrated (in use but not fully)
- **internal/ui:** Available in git/beads but not used in cmd files

### Needs Integration (direct exec calls remain)
- **internal/git:** cmd/init.go has 8 direct `exec.Command("git", ...)` calls
- **internal/beads:** cmd/init.go has 1 direct `exec.Command("bd", "init")` call

### Needs Reconciliation
- **internal/config:** Conflicts with cmd/config.Config struct (18-field version)

---

## Recommendations

### Priority 1: Integrate internal/git and internal/beads
**Issue:** ISSUE-013, ISSUE-014

1. **ISSUE-013: Integrate internal/git Package**
   - Replace 8 direct git exec.Command calls in cmd/init.go
   - Potentially expand internal/git functions to cover all use cases
   - Create tests for internal/git package
   - May need wrapper functions for user.name/user.email configuration

2. **ISSUE-014: Integrate internal/beads Package**
   - Replace 1 direct bd exec.Command call in cmd/init.go
   - Create tests for internal/beads package

### Priority 2: Expand internal/ui usage
**Issue:** ISSUE-015

Replace fmt.Printf/fprintf calls with ui package functions in:
- cmd/init.go (multiple ui.Info/Success calls)
- cmd/help.go (use ui.Header for section headers)
- cmd/version.go (use ui.Info or ui.Success)
- main.go (use ui.Error for error messages)

### Priority 3: Reconcile Config struct duplication
**Issue:** ISSUE-018

**Problem:** Two different Config structs in codebase:
- `internal/config.Config`: 9 fields (newer, simpler)
- `cmd/config.Config`: 18 fields (existing, more comprehensive)

**Solution Options:**
1. Merge configs: Create unified Config with all 18 fields
2. Separate concerns: Keep internal for runtime, cmd for command-specific
3. Re-export approach: Have cmd/config re-export internal/config with extensions

**Recommended:** Option 1 - Unified Config
- Move comprehensive Config to internal/config
- Re-export from cmd/config for backwards compatibility
- Update internal/git and internal/beads to use Config if needed

---

## Test Coverage Assessment

### Current Test Status
| Package | Unit Tests | Lines | Estimated Coverage |
|---------|-----------|-------|-------------------|
| internal/validate | 24+ cases | 77 | ~90% |
| internal/git | None | 94 | ~0% |
| internal/beads | None | 46 | ~0% |
| internal/ui | None | 56 | ~0% |
| internal/config | None | 22 | ~0% |

### Recommended Test Targets
- **internal/git:** >70% (3 functions, moderate complexity)
- **internal/beads:** >70% (2 functions, low complexity)
- **internal/ui:** >60% (5 functions, mostly output - skip tests ok)
- **internal/config:** >80% (1 function, trivial)

---

## Blockers and Dependencies

### ISSUE-013 Blocker Analysis (internal/git)
- **CheckAvailable()** can be tested without git installed (mock exec.LookPath)
- **Init()** needs either:
  - Real git command available
  - Mock exec.Command calls
  - Skip test if git not available
- **CreateCommit()** same constraints as Init()

### ISSUE-014 Blocker Analysis (internal/beads)
- **CheckAvailable()** can be tested without bd installed (mock exec.LookPath)
- **Init()** needs either:
  - Real bd command available
  - Skip test if bd not installed

### ISSUE-018 Blocker Analysis (Config Reconciliation)
- Can be done independently
- Does not block other issues
- Recommended: Solve after ISSUE-013/014/015 to understand full usage patterns

---

## Next Steps

1. **Proceed with ISSUE-012 Complete** ✓ (This audit)
2. **ISSUE-013:** Integrate internal/git
   - Create internal/git/git_test.go
   - Refactor cmd/init.go to use internal/git functions
3. **ISSUE-014:** Integrate internal/beads
   - Create internal/beads/beads_test.go
   - Refactor cmd/init.go to use internal/beads functions
4. **ISSUE-015:** Integrate internal/ui
   - Create internal/ui/ui_test.go (optional - output functions)
   - Update cmd files to use ui.* functions
5. **ISSUE-016:** Create internal/config Package (already exists)
   - Verify structure vs spec requirements
6. **ISSUE-018:** Reconcile Config Struct
   - Decide merge/separate strategy
   - Update imports and usage

---

## Files Summary

```
internal/
├── git/
│   ├── git.go          (94 lines, 3 functions, NO TESTS)
│   └── [git_test.go]   (needed: ISSUE-013)
├── beads/
│   ├── beads.go        (46 lines, 2 functions, NO TESTS)
│   └── [beads_test.go] (needed: ISSUE-014)
├── ui/
│   ├── ui.go           (56 lines, 6 functions, NO TESTS)
│   └── [ui_test.go]    (recommended: ISSUE-015)
├── validate/
│   ├── validate.go     (77 lines, 3 functions)
│   └── validate_test.go (309 lines, 3 test functions, 40+ test cases) ✓
├── config/
│   └── config.go       (22 lines, 1 function, NO TESTS)
└── AUDIT.md            (this file)
```

---

**Audit Completed:** 2025-11-24
**Auditor:** Claude Code Assistant
**Status:** Ready for ISSUE-013 (Integrate internal/git Package)

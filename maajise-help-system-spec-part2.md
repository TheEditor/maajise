# Maajise Help System Enhancement - Task Specification (Part 2)

## ⚠️ READ FIRST

**BEFORE STARTING ANY WORK:**
1. Read `AGENTS.md` in the project root for critical project context and Beads workflow
2. Read `/mnt/skills/user/bd-issue-tracking/SKILL.md` for issue tracking best practices
3. Then read this entire specification

---

## ⚠️ PREREQUISITE: Complete Part 1 First

This specification continues the help system work from Part 1. Before starting Part 2:
- Part 1 must be fully complete (all issues closed)
- All Part 1 tests must pass
- Core help infrastructure must be working

---

## Project Context

**Continuation of:** Part 1 (ISSUE-000 through ISSUE-003)
**This Part Contains:** ISSUE-004 through ISSUE-005
**Focus:** Flag documentation and user-friendly error messages

---

## ⚠️ IMPLEMENTATION WORKFLOW

**STEP 1:** Verify Part 1 is complete
```bash
bd list --status closed
# Should show 4 closed issues from Part 1
```

**STEP 2:** Read this spec to understand remaining work

**STEP 3:** Create Beads issues for Part 2 (see "CREATE ISSUES NOW" section below)

**STEP 4:** Verify new issues created
```bash
bd ready --json
# Should show at least 1 ready issue (ISSUE-004 or ISSUE-005 equivalent)
```

**STEP 5:** Implement in order (see "Implementation Order" section below)

**STEP 6:** Mark complete as you go
```bash
# Use actual ID returned from bd create
bd close <actual-issue-id> --reason "Implemented and tested"
```

---

## Issues (Continued)

### ISSUE-004: Document All Flags in Examples

**Severity:** INFO
**File:** All command files (cmd/init.go, cmd/help.go, etc.)
**Description:** Ensure all command flags are documented through examples with contextual explanations. Each flag should appear in at least one example showing realistic usage.

**Current State:** Init command has --in-place, --no-overwrite, --template flags but they may not all be documented in examples

**Required Changes:**

1. Audit all commands for flags (use grep or code review):
```bash
# Find all flag definitions
grep -r "flag.Bool\|flag.String\|flag.Int" cmd/
```

2. For each flag found, ensure Examples() includes:
   - At least one example using the flag
   - Brief explanation of what the flag does (indented comment)
   - Realistic use case

3. Example pattern to follow:
```go
  # Use TypeScript template
  maajise init my-ts-app --template=typescript
      Creates TypeScript project with tsconfig.json, package.json
      Available templates: base, typescript, python, rust, php
```

4. Flags to document for InitCommand (verify current flags in code):
   - --in-place: Initialize in current directory
   - --no-overwrite: Skip existing files (safe mode)
   - --template: Specify project template
   - Any other flags present

5. Update InitCommand examples to show all template options:
```go
  # Use TypeScript template
  maajise init my-ts-app --template=typescript
      Available templates: base, typescript, python, rust, php
      TypeScript: Creates tsconfig.json, package.json, .gitignore
```

**Test Requirements:**
- Manually review Examples() for each command
- Verify all flags appear in at least one example
- Verify flag examples are accurate (can be copy-pasted and run)
- Check that template options are all listed

**Dependencies:** ISSUE-001 (Examples() method must exist)
**Blocks:** None

---

### ISSUE-005: Add Usage Errors with Help Hints

**Severity:** INFO
**File:** `internal/ui/errors.go` (new file), all command files
**Description:** When commands fail due to incorrect usage, show helpful error message with hint to use --help. Provides better user experience for validation failures.

**Required Changes:**

1. Create error helper in internal/ui/errors.go:
```go
package ui

import "fmt"

// UsageError returns formatted usage error with help hint
func UsageError(cmdName string, msg string) error {
    return fmt.Errorf("%s\n\nRun 'maajise %s --help' for usage information", msg, cmdName)
}
```

2. Update all commands to use UsageError for validation failures.

**Example: InitCommand validation**
```go
// In InitCommand.Run()

// Validate project name requirement
if projectName == "" && !inPlace {
    return ui.UsageError("init", "project name is required when not using --in-place")
}

// Validate template selection
if template != "" && !isValidTemplate(template) {
    validTemplates := []string{"base", "typescript", "python", "rust", "php"}
    return ui.UsageError("init", 
        fmt.Sprintf("unknown template: %s\nAvailable templates: %s", 
            template, strings.Join(validTemplates, ", ")))
}

// Validate conflicting flags
if inPlace && projectName != "" {
    return ui.UsageError("init", "--in-place cannot be used with a project name")
}
```

3. Add validation for common error cases:
   - Missing required arguments
   - Invalid flag values
   - Conflicting flags
   - Invalid file paths

**Test Requirements:**
- Test missing required args show usage error
- Test invalid flag values show usage error
- Test conflicting flags show usage error
- Verify help hint points to correct command
- Test error messages are clear and actionable

**Dependencies:** None (independent utility)
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

**PHASE 3: User-Facing Features (CAN BE DONE IN EITHER ORDER)**

4. **ISSUE-004** - Document flags in examples
   - Rationale: Improves help text completeness
   - Time: 0.25 session
   - Verification: All flags appear in examples

5. **ISSUE-005** - Usage errors with hints
   - Rationale: Better error messages for users
   - Time: 0.25 session
   - Verification: Validation failures show helpful errors

**Total Estimated Time:** 0.5 sessions

---

## Testing Strategy

### Automated Unit Tests

Create `cmd/examples_test.go`:

```go
package cmd

import (
    "strings"
    "testing"
)

func TestInitExamplesIncludeAllFlags(t *testing.T) {
    cmd := &InitCommand{}
    examples := cmd.Examples()
    
    // Check that all known flags are documented
    requiredFlags := []string{"--in-place", "--no-overwrite", "--template"}
    for _, flag := range requiredFlags {
        if !strings.Contains(examples, flag) {
            t.Errorf("Examples missing flag: %s", flag)
        }
    }
    
    // Check that template options are listed
    templates := []string{"typescript", "python", "rust", "php"}
    for _, tmpl := range templates {
        if !strings.Contains(examples, tmpl) {
            t.Errorf("Examples missing template option: %s", tmpl)
        }
    }
}
```

Create `internal/ui/errors_test.go`:

```go
package ui

import (
    "strings"
    "testing"
)

func TestUsageError(t *testing.T) {
    err := UsageError("init", "project name required")
    errMsg := err.Error()
    
    // Check error contains user message
    if !strings.Contains(errMsg, "project name required") {
        t.Error("Error message missing user message")
    }
    
    // Check error contains help hint
    if !strings.Contains(errMsg, "maajise init --help") {
        t.Error("Error message missing help hint")
    }
    
    // Check format
    if !strings.Contains(errMsg, "\n\n") {
        t.Error("Error message should have blank line between message and hint")
    }
}
```

### Integration Tests (Manual)

Run these commands to verify Part 2 additions:

```bash
# Test flag documentation completeness
maajise help init | grep "template"    # Should show template options
maajise help init | grep "in-place"    # Should show in-place flag
maajise help init | grep "no-overwrite" # Should show no-overwrite flag

# Test usage error messages
maajise init                           # Should show error with help hint
maajise init --template=invalid        # Should show error with valid templates
maajise init --in-place myproject      # Should show conflicting flags error

# Verify error format
maajise init 2>&1 | grep "maajise init --help"  # Help hint present
```

### Manual Testing Checklist

- [ ] `maajise help` lists all commands
- [ ] `maajise help init` shows init details
- [ ] `maajise init --help` shows same as `help init`
- [ ] `maajise --help` shows global help
- [ ] Error messages include help hints
- [ ] All flags documented in examples
- [ ] Examples are accurate and copy-pasteable
- [ ] Help output is readable (80 column wrap)
- [ ] Output uses colors from internal/ui
- [ ] Template options clearly listed

### Coverage Targets

- **Help command:** >90% (critical for UX)
- **Command help methods:** 100% (trivial getters)
- **Error helpers:** 100%
- **Registry functions:** 100%
- **Overall project:** >70%

Run coverage:
```bash
# All tests with coverage
go test ./... -cover

# Detailed coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## CREATE ISSUES NOW ⚠️

**PREREQUISITE**: Part 1 must be complete (4 issues closed).

### Beads Workflow Reminders

**Before starting:** Read `/mnt/skills/user/bd-issue-tracking/SKILL.md` if you haven't already.

**Key concepts:**
- Beads assigns IDs automatically (you don't choose them)
- Each `bd create` returns JSON like `{"id":"maajise-XXX",...}`
- Use the actual returned IDs in your work, not the ISSUE-004 labels in this spec
- ISSUE-004, ISSUE-005 are just labels for humans to reference tasks in this doc

### Step 1: Verify Part 1 Complete

```bash
# Check that Part 1 issues are closed
bd list --status closed

# Should see 4 closed issues (the actual IDs from Part 1)
```

### Step 2: Create Part 2 Issues

Run these commands **in order**. Each returns JSON with an assigned `id` field.

```bash
# Create ISSUE-004
bd create --title "ISSUE-004: Document All Flags in Examples" \
  --description "Ensure all command flags are documented in Examples() with contextual explanations. See ISSUE-004 in spec for pattern." \
  --json

# Capture the returned ID (looks like: {"id":"maajise-046",...})
# You'll use this actual ID in commits and close commands

# Create ISSUE-005
bd create --title "ISSUE-005: Add Usage Errors with Help Hints" \
  --description "Show helpful errors with --help hints for validation failures. Create UsageError helper. See ISSUE-005 in spec for implementation." \
  --json
```

### Step 3: Verify Issues Created

```bash
# Check all issues were created
bd list --status open

# Verify what's ready to work on
bd ready --json
```

**Expected Output:** 
- `bd list --status open` should show 2 open issues
- `bd ready` should show both issues (neither has dependencies)

If `bd ready` shows 0 issues, something went wrong with issue creation.

### Step 4: Begin Implementation

Work through issues in either order (see "Implementation Order" section below - no dependencies between them).

As you complete each issue:
```bash
# Mark complete (use actual ID returned from creation)
bd close <actual-issue-id> --reason "Implemented and tested"

# Verify it's gone from ready list
bd ready --json
```

**Important**: Include the issue ID in your commit messages:
```bash
git commit -m "feat: document flags in examples (bd:<actual-issue-id>)"
```

---

## Success Criteria

### Issue Tracking
- [ ] 2 new issues created (actual IDs captured from bd create)
- [ ] Previous 4 issues from Part 1 remain closed
- [ ] All issues closed with `bd close <id> --reason "..."`

### Code Changes
- [ ] All flags documented in command examples
- [ ] Template options clearly listed in init command
- [ ] UsageError helper created in internal/ui
- [ ] All commands use UsageError for validation failures
- [ ] Error messages include help hints

### Testing
- [ ] All new unit tests pass
- [ ] Examples test verifies all flags present
- [ ] UsageError test passes (100% coverage)
- [ ] No regressions in Part 1 functionality

### Documentation
- [ ] All flags appear in examples with explanations
- [ ] Template options listed (base, typescript, python, rust, php)
- [ ] Examples are copy-pasteable

### User Experience
- [ ] Error messages guide users to --help
- [ ] Invalid inputs show clear, actionable errors
- [ ] Help text is complete and polished

### Verification
- [ ] Clean build: `go build`
- [ ] All tests pass: `go test ./...`
- [ ] All Part 2 issues closed: `bd list --status open` shows 0
- [ ] Help system fully complete

---

## Notes for Implementation

### Flag Documentation

Ensure every flag appears in at least one example:
- Show realistic use cases
- Include contextual explanations (indented 6 spaces)
- List available options where applicable (templates, etc.)

### Error Messages

Good error pattern:
```
<specific problem>

Run 'maajise <command> --help' for usage information
```

Bad error:
```
Invalid input
```

### Validation Examples

```go
// Good - specific and helpful
if !isValidTemplate(template) {
    return ui.UsageError("init", fmt.Sprintf(
        "unknown template: %s\nAvailable: base, typescript, python, rust, php",
        template))
}

// Bad - vague
if !isValidTemplate(template) {
    return fmt.Errorf("invalid template")
}
```

---

## End of Part 2

**Completed:** Full help system implementation
**Status:** All 6 issues complete
**Next:** Help system is production-ready

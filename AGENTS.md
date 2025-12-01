# Agent Instructions for Maajise

## BEFORE ANYTHING ELSE

**CRITICAL**: Run `bd onboard` and follow the instructions.

After onboarding, read this entire file before starting any work.

---

## Project Overview

**Maajise** is a CLI tool for standardized repository initialization that sets up Git repositories, Beads issue tracking, and UBS (Ultimate Bug Scanner) configurations.

**Primary Goal**: Transform Maajise from a functional but monolithic tool into a mature, professional-grade CLI that can initialize projects across multiple programming languages (TypeScript, Python, Rust, PHP, Go).

**Technology Stack**: Go, Git integration, Beads integration, UBS configuration

**Current State**: Phase 1 (architectural migration to modular command-based structure) partially complete. Tool successfully initializes repositories but needs refactoring from monolithic main.go to clean subcommand architecture.

---

## Working with Beads Issue Tracking

### Step 0: Read the bd-issue-tracking Skill

**MANDATORY FIRST STEP**: Before doing ANYTHING else in this project, execute:

```
file_read /mnt/skills/user/bd-issue-tracking/SKILL.md
```

This skill explains:
- How bd assigns issue IDs automatically (you cannot choose them)
- How to track which issue you're currently working on
- How to properly close issues when work is complete
- How `bd ready` shows what's available to work on
- Common pitfalls and how to avoid them

**DO NOT PROCEED** until you have read and understood this skill.

---

### Understanding Beads ID Assignment

**CRITICAL CONCEPT**: Beads assigns issue IDs automatically. You CANNOT specify them.

#### When Creating Issues

```bash
# You run this command:
bd create "Add SQLite tracking database" -t task -p 0 -d "Description..." --json

# Beads returns JSON like this:
{"id":"maajise-008","title":"Add SQLite tracking database",...}

# The ID "maajise-008" was ASSIGNED by the system
# You must capture it and use it in subsequent commands
```

#### Using Captured IDs

When task specs show commands like:
```bash
bd create "Task" -t task -p 0 --parent maajise-007 -d "..." --json
```

The `maajise-007` is a **placeholder**. Replace it with the **actual ID** returned from the parent epic creation.

**Example Workflow**:
1. Create parent epic → returns `{"id":"maajise-042",...}`
2. Use `maajise-042` (not `maajise-007`) in all child task commands
3. Capture each child task ID as it's returned
4. Use those actual IDs in dependency commands

---

### Daily Workflow

#### Starting Work

1. **Check ready work**:
   ```bash
   bd ready --json
   ```
   This shows issues with no blockers that are ready to work on.

2. **Pick an issue**: Choose based on priority (P0 = highest, P3 = lowest)

3. **Update status** (optional but recommended):
   ```bash
   bd update <issue-id> --status in_progress
   ```

4. **Do the work**: Implement according to the task spec

#### During Work

- **Discovered new work?** File an issue immediately:
  ```bash
  bd create "Found bug in parser" -t bug -p 1 --json
  bd dep add <new-issue-id> <current-issue-id> --type discovered-from
  ```

- **Need to check dependencies?**
  ```bash
  bd dep tree <issue-id>
  ```

- **Stuck on a blocker?** Update the issue:
  ```bash
  bd update <issue-id> --status blocked
  ```

#### Completing Work

1. **Verify acceptance criteria**: Check the issue's acceptance criteria are met

2. **Commit your work**:
   ```bash
   git add .
   git commit -m "feat: implement feature (bd:<issue-id>)"
   git push
   ```
   
   **IMPORTANT**: Include `(bd:<issue-id>)` in commit message for traceability

3. **Close the issue**:
   ```bash
   bd close <issue-id> --reason "Implemented and tested"
   ```

4. **Verify closure**:
   ```bash
   bd ready --json
   ```
   Confirm the issue is no longer in the ready list

---

### Session Ending Protocol

**CRITICAL**: Before ending ANY session, complete this checklist:

#### 1. Issue Tracker Hygiene

- [ ] File issues for any discovered bugs, TODOs, or follow-up work
- [ ] Close all completed issues with `bd close <issue-id>`
- [ ] Update status for any in-progress work
- [ ] Run `bd ready` to confirm no orphaned issues

#### 2. Quality Gates (if code changed)

- [ ] Tests pass: `make test`
- [ ] Code builds: `make build`
- [ ] Coverage meets threshold: `make test-coverage`
- [ ] If builds broken: File P0 issue immediately

#### 3. Sync Issue Tracker

```bash
# Sync Beads database
bd sync

# Add and commit changes
git add .beads/issues.jsonl
git commit -m "chore: sync issue tracker"

# Push to remote
git push
```

**Handle conflicts carefully**:
- If merge conflict in `.beads/issues.jsonl`, resolve thoughtfully
- Sometimes accepting remote and re-importing is safest
- Goal: ensure no issues are lost

#### 4. Verify Clean State

- [ ] All code changes committed and pushed
- [ ] No untracked files: `git status`
- [ ] Beads database synced: `bd info`
- [ ] Ready work reviewed: `bd ready`

#### 5. Next Session Context

Provide a brief summary for the next session:
- What was completed
- What issues remain
- What's ready to work on next
- Any blockers or concerns

**Example**:
```
Session Summary:
- Completed: maajise-008 (SQLite DB), maajise-009 (Diff detection)
- Ready next: maajise-010 (Sync command) - no blockers
- Note: maajise-015 blocked on maajise-010 completion
```

---

## Task Specification Workflow

### Reading Task Specs

Task specifications are located in the project root:
- `maajise-phase1-task-spec.md` - Phase 1 implementation
- `maajise-phase2-task-spec.md` - Phase 2 implementation
- etc.

Each spec contains:
1. **Beads Issue Setup** - Commands to create issues (you execute these)
2. **Task Details** - Implementation instructions for each issue
3. **Acceptance Criteria** - How to verify completion
4. **Time Estimates** - Expected duration for planning

### Executing Task Specs

**Step 1: Create Issues**
- Read the "Beads Issue Setup" section
- Execute commands **in order** (STEP 1, STEP 2, STEP 3, etc.)
- **Capture each returned ID** and use in subsequent commands
- Do NOT use the example IDs (maajise-001, maajise-007, etc.) - use actual returned IDs

**Step 2: Verify Setup**
```bash
bd ready --json  # Should show tasks with no blockers
bd dep tree <any-task-id>  # Verify dependency graph
```

**Step 3: Work Through Tasks**
- Pick ready tasks from `bd ready`
- Implement according to task section in spec
- Verify acceptance criteria
- Close when complete

### Common Pitfalls

❌ **DON'T**: Use example IDs from spec (maajise-007, etc.)
✅ **DO**: Use actual IDs returned by bd create

❌ **DON'T**: Create issues out of order
✅ **DO**: Follow STEP 1, 2, 3, 4 exactly as written

❌ **DON'T**: Forget to close issues when done
✅ **DO**: Run `bd close <issue-id>` after each task

❌ **DON'T**: Skip the bd-issue-tracking skill
✅ **DO**: Read it first, every time

---

## Git Commit Conventions

Follow Conventional Commits format with Beads issue reference:

```
<type>(<scope>): <description> (bd:<issue-id>)

<body>

<footer>
```

**Examples**:
```bash
git commit -m "feat: add SQLite tracking database (bd:maajise-008)

- Create schema with findings table
- Implement CRUD operations
- Add fingerprint computation
- Include comprehensive tests"

git commit -m "fix: correct fingerprint algorithm (bd:maajise-008)

Use code context instead of line numbers for stability"

git commit -m "docs: add Phase 2 sync documentation (bd:maajise-014)"
```

**Types**: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `perf`

---

## Project Structure

```
maajise/
├── cmd/                  # Command implementations
│   ├── command.go       # Command interface & registry
│   ├── init.go          # Init command
│   ├── help.go          # Help command
│   └── version.go       # Version command
├── internal/
│   ├── git/             # Git operations
│   ├── beads/           # Beads integration
│   ├── config/          # Configuration handling
│   ├── ui/              # UI utilities (colors, output)
│   └── validate/        # Input validation
├── templates/           # Project templates
│   ├── base.go          # Base template
│   ├── typescript.go    # TypeScript template
│   ├── python.go        # Python template
│   ├── rust.go          # Rust template
│   └── php.go           # PHP template
├── .beads/              # Beads issue tracker database
├── main.go              # CLI entry point
├── go.mod               # Go dependencies
├── README.md            # User documentation
└── AGENTS.md            # This file
```

---

## Testing Requirements

**All code must have tests**. Acceptance criteria often specify coverage thresholds.

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

**Coverage threshold**: 80% minimum (enforced by `make test-coverage`)

### Test Organization

- Unit tests: `*_test.go` alongside implementation
- Integration tests: `cmd/maajise/*_test.go`
- Test fixtures: `testdata/` directory

---

## Build & Run

```bash
# Build binary (Windows)
.\build.bat

# Build binary (Linux/Mac)
./build.sh

# Run (after building)
./maajise init my-project
./maajise init --in-place
./maajise init my-project --template=typescript

# Install globally (Windows PowerShell as Admin)
.\install.ps1

# After global install, use from anywhere
maajise init my-new-project
```

---

## Getting Help

### Beads Commands Reference

```bash
bd help                    # General help
bd help create            # Help for specific command
bd ready                  # Show ready work
bd show <issue-id>        # Show issue details
bd list --status open     # List open issues
bd dep tree <issue-id>    # Show dependency tree
bd info                   # Database info
```

### When Stuck

1. **Read the task spec** - Implementation details are there
2. **Check issue dependencies** - `bd dep tree <issue-id>`
3. **Review acceptance criteria** - What does "done" look like?
4. **Check bd-issue-tracking skill** - Common patterns and solutions
5. **Ask the human** - When genuinely blocked

---

## Critical Reminders

1. **Read bd-issue-tracking skill FIRST** - Every time you start work
2. **Capture actual IDs** - Don't use example IDs from specs
3. **Close issues when done** - Verify with `bd ready`
4. **Sync before ending session** - `bd sync` then `git push`
5. **Follow session ending protocol** - Complete the entire checklist
6. **Include issue ID in commits** - `(bd:<issue-id>)` format

---

## Questions?

If anything in this document is unclear, ask the human before proceeding. It's better to clarify expectations upfront than to waste time on incorrect implementation.

---

**Last Updated**: 2025-11-30
**Project Status**: Phase 1 (architecture migration) in progress, help system enhancement specified

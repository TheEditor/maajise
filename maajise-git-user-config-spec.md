# Maajise - Git User Configuration Feature

## Project Context
Repository: Maajise CLI tool
Current Version: 2.0.0
Task: Add explicit Git user.name and user.email configuration during repo initialization
Date: 2025-11-23

## IMPLEMENTATION WORKFLOW ⚠️

**STEP 1:** Create all Beads issues
Run the commands in the "CREATE ISSUES NOW" section below

**STEP 2:** Verify issues created
```bash
bd list
# Should show all 3 issues
```

**STEP 3:** Begin implementation
Work through issues in the order specified in "Implementation Order" section

**STEP 4:** Mark issues complete
As each issue is resolved:
```bash
bd close <issue-id>
```

---

## Issues

### ISSUE-001: Add Git User Configuration Function

**Severity:** INFO
**File:** main.go (new function after initGit)
**Description:** Create a new function to prompt for and configure Git user.name and user.email at the repository level. This ensures commits have correct author information and prevents Claude/AI attribution issues.

**Current State:** Git is initialized without setting user configuration. User relies on global config which may be incorrect or missing.

**Required Changes:**

1. Add new function `configureGitUser()` after `initGit()`:

```go
func configureGitUser(repoDir string, cfg *Config) error {
	if cfg.SkipGit {
		return nil
	}

	info("Configuring Git user...")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Get user.name
	fmt.Printf("%s→%s Git user.name (full name): ", blue, reset)
	userName, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user name: %w", err)
	}
	userName = strings.TrimSpace(userName)

	if userName == "" {
		errorMsg("Git user.name cannot be empty")
		return fmt.Errorf("empty user name")
	}

	// Get user.email
	fmt.Printf("%s→%s Git user.email (email address): ", blue, reset)
	userEmail, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user email: %w", err)
	}
	userEmail = strings.TrimSpace(userEmail)

	if userEmail == "" {
		errorMsg("Git user.email cannot be empty")
		return fmt.Errorf("empty user email")
	}

	// Validate email format (basic check)
	if !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") {
		errorMsg("Invalid email format")
		return fmt.Errorf("invalid email format")
	}

	// Set user.name
	nameCmd := exec.Command("git", "config", "user.name", userName)
	nameCmd.Dir = repoDir
	if cfg.Verbose {
		nameCmd.Stdout = os.Stdout
		nameCmd.Stderr = os.Stderr
	}
	if err := nameCmd.Run(); err != nil {
		errorMsg("Failed to set user.name")
		return err
	}

	// Set user.email
	emailCmd := exec.Command("git", "config", "user.email", userEmail)
	emailCmd.Dir = repoDir
	if cfg.Verbose {
		emailCmd.Stdout = os.Stdout
		emailCmd.Stderr = os.Stderr
	}
	if err := emailCmd.Run(); err != nil {
		errorMsg("Failed to set user.email")
		return err
	}

	success(fmt.Sprintf("Git user configured: %s <%s>", userName, userEmail))
	fmt.Println()
	return nil
}
```

2. Call `configureGitUser()` in `main()` after `initGit()`:

```go
// In main() function, after:
if err := initGit(repoPath, cfg); err != nil {
	os.Exit(1)
}

// Add:
if err := configureGitUser(repoPath, cfg); err != nil {
	os.Exit(1)
}
```

**Test Requirements:**
- Unit tests for email validation logic
- Test empty name rejection
- Test empty email rejection
- Test invalid email format rejection
- Test valid email formats accepted
- Integration test: verify git config commands execute
- Integration test: verify user.name and user.email set in repo

**Dependencies:** None
**Blocks:** ISSUE-002

---

### ISSUE-002: Add Skip Flag for Git User Configuration

**Severity:** INFO
**File:** main.go (parseFlags and Config struct)
**Description:** Add `--skip-git-user` flag to allow bypassing interactive Git user configuration. Useful for automated workflows or when user wants to configure manually.

**Required Changes:**

1. Add field to `Config` struct:
```go
type Config struct {
	ProjectName  string
	InPlace      bool
	NoOverwrite  bool
	SkipGit      bool
	SkipBeads    bool
	SkipCommit   bool
	SkipRemote   bool
	SkipGitUser  bool  // NEW
	Verbose      bool
}
```

2. Add flag in `parseFlags()`:
```go
flag.BoolVar(&cfg.SkipGitUser, "skip-git-user", false, "Skip Git user configuration")
```

3. Update `showUsage()` to document new flag:
```go
Flags:
  --in-place          Initialize in current directory (no nested folders)
  --no-overwrite      Skip files that already exist
  --skip-git          Don't initialize Git
  --skip-beads        Don't initialize Beads
  --skip-commit       Don't create initial commit
  --skip-remote       Don't prompt for remote setup
  --skip-git-user     Don't prompt for Git user config  // NEW
  -v, --verbose       Verbose output
  -h, --help          Show this help
```

4. Update `configureGitUser()` to check flag:
```go
func configureGitUser(repoDir string, cfg *Config) error {
	if cfg.SkipGit || cfg.SkipGitUser {
		if cfg.Verbose {
			info("Skipping Git user configuration")
		}
		return nil
	}
	// ... rest of function
}
```

**Test Requirements:**
- Test `--skip-git-user` flag prevents prompting
- Test flag appears in help output
- Test `--skip-git` also skips user config
- Test combination of flags (e.g., `--skip-git-user --skip-remote`)

**Dependencies:** ISSUE-001 (needs configureGitUser function)
**Blocks:** None

---

### ISSUE-003: Add Non-Interactive Mode for Git User Config

**Severity:** INFO
**File:** main.go (parseFlags and configureGitUser)
**Description:** Add optional flags `--git-name` and `--git-email` to provide Git user configuration non-interactively. Useful for automation and CI/CD pipelines.

**Required Changes:**

1. Add fields to `Config` struct:
```go
type Config struct {
	ProjectName  string
	InPlace      bool
	NoOverwrite  bool
	SkipGit      bool
	SkipBeads    bool
	SkipCommit   bool
	SkipRemote   bool
	SkipGitUser  bool
	GitName      string  // NEW
	GitEmail     string  // NEW
	Verbose      bool
}
```

2. Add flags in `parseFlags()`:
```go
flag.StringVar(&cfg.GitName, "git-name", "", "Git user.name (non-interactive)")
flag.StringVar(&cfg.GitEmail, "git-email", "", "Git user.email (non-interactive)")
```

3. Update `showUsage()`:
```go
Flags:
  --in-place          Initialize in current directory (no nested folders)
  --no-overwrite      Skip files that already exist
  --skip-git          Don't initialize Git
  --skip-beads        Don't initialize Beads
  --skip-commit       Don't create initial commit
  --skip-remote       Don't prompt for remote setup
  --skip-git-user     Don't prompt for Git user config
  --git-name NAME     Set Git user.name (non-interactive)      // NEW
  --git-email EMAIL   Set Git user.email (non-interactive)    // NEW
  -v, --verbose       Verbose output
  -h, --help          Show this help
```

4. Update `configureGitUser()` to use provided values:
```go
func configureGitUser(repoDir string, cfg *Config) error {
	if cfg.SkipGit || cfg.SkipGitUser {
		if cfg.Verbose {
			info("Skipping Git user configuration")
		}
		return nil
	}

	var userName, userEmail string

	// Check if values provided via flags
	if cfg.GitName != "" && cfg.GitEmail != "" {
		userName = cfg.GitName
		userEmail = cfg.GitEmail
		
		// Validate email format
		if !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") {
			errorMsg("Invalid email format")
			return fmt.Errorf("invalid email format")
		}
		
		if cfg.Verbose {
			info(fmt.Sprintf("Using provided Git config: %s <%s>", userName, userEmail))
		}
	} else {
		// Interactive prompts (existing code)
		info("Configuring Git user...")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)
		// ... existing interactive code ...
	}

	// Set git config (existing code)
	// ...
}
```

**Test Requirements:**
- Test `--git-name` and `--git-email` flags work together
- Test invalid email format rejected
- Test missing one flag (name without email) still prompts interactively
- Test non-interactive mode doesn't prompt
- Test verbose mode shows provided values

**Dependencies:** ISSUE-001, ISSUE-002
**Blocks:** None

---

## Implementation Order

Work through issues in this order:

1. **ISSUE-001** - Add Git user configuration function
   - Core functionality
   - No dependencies
   - Blocks other issues

2. **ISSUE-002** - Add skip flag
   - Simple addition
   - Depends on ISSUE-001

3. **ISSUE-003** - Add non-interactive mode
   - Enhancement
   - Depends on ISSUE-001 and ISSUE-002

---

## CREATE ISSUES NOW ⚠️

**BEFORE IMPLEMENTING:** Run these commands to create all issues:

```bash
bd create --title "ISSUE-001: Add Git user configuration function" --description "Prompt for and set user.name and user.email during git init. See ISSUE-001 in spec."

bd create --title "ISSUE-002: Add skip flag for Git user config" --description "Add --skip-git-user flag to bypass interactive prompts. See ISSUE-002 in spec."

bd create --title "ISSUE-003: Add non-interactive mode for Git user config" --description "Add --git-name and --git-email flags for automation. See ISSUE-003 in spec."
```

**VERIFY:** Run `bd list` to confirm all 3 issues created.

---

## Testing Strategy

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific test
go test -v -run TestConfigureGitUser
```

### Test Coverage Goals

- New `configureGitUser()` function: > 80%
- Email validation logic: 100%
- Flag parsing: > 70%

### Manual Testing

```bash
# Test interactive mode
./maajise.exe test-repo

# Test skip flag
./maajise.exe --skip-git-user test-repo2

# Test non-interactive mode
./maajise.exe --git-name "John Doe" --git-email "john@example.com" test-repo3

# Verify config was set
cd test-repo/test-repo
git config user.name
git config user.email
```

---

## Success Criteria

Task completion checklist:

### Issue Tracking
- [ ] All 3 Beads issues created (verify with `bd list`)
- [ ] Issues have correct dependencies mapped
- [ ] Issues closed as work completes

### Code Changes
- [ ] `configureGitUser()` function implemented
- [ ] Function called in main() after initGit()
- [ ] `--skip-git-user` flag added and working
- [ ] `--git-name` and `--git-email` flags added and working
- [ ] Help text updated with all new flags
- [ ] Email validation implemented

### Testing
- [ ] Unit tests pass for email validation
- [ ] Unit tests pass for empty input rejection
- [ ] Integration tests verify git config commands work
- [ ] Manual testing confirms user.name and user.email set correctly
- [ ] Test all flag combinations work

### Verification
- [ ] Build succeeds: `go build -o maajise.exe`
- [ ] Help shows new flags: `./maajise.exe --help`
- [ ] Interactive mode prompts for name/email
- [ ] Non-interactive mode works without prompts
- [ ] Skip flag prevents prompting
- [ ] Git config properly set in test repos

---

## Notes for Implementation

**Code Style:**
- Follow existing error handling patterns
- Use existing color constants (blue, reset) for prompts
- Use existing output functions (info, success, errorMsg)
- Keep function signatures consistent with other functions

**Email Validation:**
- Use simple validation (contains @ and .)
- Don't over-engineer with complex regex
- Clear error messages for invalid format

**User Experience:**
- Prompts should be clear and concise
- Show configured values after setting
- Provide helpful error messages
- Support verbose mode for debugging

**Git Commands:**
- Use `git config user.name` (local, not --global)
- Use `git config user.email` (local, not --global)
- Set working directory (cmd.Dir) to repoDir
- Handle command failures gracefully

**Testing Notes:**
- Mock stdin for testing interactive prompts
- Use temporary directories for integration tests
- Clean up test repos in defer statements
- Test both Windows and Unix line endings in input

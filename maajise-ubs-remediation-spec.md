# Maajise - UBS Scan Remediation Tasks

## Project Context
Repository initialization CLI tool written in Go. Current version: 2.0.0
Scan Date: 2025-11-23
Scan Tool: UBS (Ultimate Bug Scanner) v5.0.0

## Critical Issues

### ISSUE-001: Command Injection Risk in Git Remote Setup
**Severity:** CRITICAL
**File:** main.go:612 (setupGitRemote function)
**Description:** User input from stdin flows directly to exec.Command without validation or sanitization. The remote URL is read from user input and passed to `git remote add origin <url>` without checking for shell metacharacters or validating URL format.

**Current Code Flow:**
```
bufio.NewReader(os.Stdin) → reader.ReadString('\n') → remoteURL → exec.Command("git", "remote", "add", "origin", remoteURL)
```

**Required Changes:**
1. Add URL validation function that:
   - Validates against git URL patterns (https://, git@, ssh://)
   - Rejects URLs containing shell metacharacters: `;`, `&`, `|`, `$`, `` ` ``, `\n`, etc.
   - Validates URL scheme and format
   - Returns error for invalid URLs

2. Sanitize input by:
   - Trimming whitespace
   - Checking maximum length (e.g., 2048 chars)
   - Rejecting multi-line input

3. Apply validation before exec.Command call

**Test Requirements:**
- Test valid HTTPS URLs (github.com, gitlab.com patterns)
- Test valid SSH URLs (git@github.com patterns)
- Test rejection of URLs with semicolons
- Test rejection of URLs with backticks
- Test rejection of URLs with pipe characters
- Test rejection of multiline input
- Test rejection of URLs exceeding max length

**Dependencies:** None
**Blocks:** ISSUE-005

---

## Warnings

### ISSUE-002: Update Go Version Directive
**Severity:** WARNING
**File:** go.mod
**Description:** go.mod specifies Go 1.22, but Go 1.23+ is recommended for improved language features and standard library updates.

**Required Changes:**
1. Update go.mod: change `go 1.22` to `go 1.23`
2. Run `go mod tidy` to ensure compatibility
3. Test build process with updated version

**Test Requirements:**
- Verify successful build with Go 1.23
- Run existing functionality tests (if any)

**Dependencies:** None

---

## Code Quality Improvements

### ISSUE-003: Add Input Validation Helper Functions
**Severity:** INFO
**File:** New file: internal/validate/validate.go
**Description:** Create reusable validation utilities to support ISSUE-001 and future input validation needs.

**Required Changes:**
1. Create `internal/validate/` package
2. Implement `ValidateGitURL(url string) error` function
3. Implement `SanitizeInput(input string) string` function
4. Implement `ContainsShellMetachars(s string) bool` function

**Function Specifications:**

```go
// ValidateGitURL checks if a string is a valid git remote URL
// Accepts: https://, git@, ssh:// formats
// Rejects: shell metacharacters, invalid schemes
func ValidateGitURL(url string) error

// SanitizeInput trims whitespace and checks basic constraints
func SanitizeInput(input string, maxLen int) (string, error)

// ContainsShellMetachars returns true if string contains shell metacharacters
func ContainsShellMetachars(s string) bool
```

**Test Requirements:**
- Unit tests for each validation function
- Test coverage > 80%
- Test edge cases (empty strings, very long strings, unicode)

**Dependencies:** None
**Blocks:** ISSUE-001

---

### ISSUE-004: Add Unit Tests for Main Package
**Severity:** INFO
**File:** New file: main_test.go
**Description:** Create comprehensive unit tests for maajise core functionality. Currently no test files exist.

**Required Test Coverage:**
1. `parseFlags()` - flag parsing logic
   - Test all flag combinations
   - Test --in-place behavior
   - Test project name extraction
   - Test help flags

2. `validateProjectName()` - project name validation
   - Test valid names
   - Test names with invalid characters
   - Test empty names
   - Test edge cases

3. `fileExists()` - file existence check
   - Test existing files
   - Test non-existent files
   - Test directories

4. `checkDependencies()` - dependency verification
   - Mock exec.LookPath for testing
   - Test with all dependencies present
   - Test with missing dependencies
   - Test with --skip flags

5. Input/Output helper functions
   - Test error message formatting
   - Test success message formatting
   - Test info message formatting
   - Test warn message formatting

**Test Requirements:**
- Use Go's standard testing package
- Use table-driven tests where appropriate
- Mock external dependencies (exec.Command, os.Stdin)
- Target > 70% code coverage for main package

**Dependencies:** None

---

### ISSUE-005: Add Integration Tests for setupGitRemote
**Severity:** INFO
**File:** New file: main_integration_test.go
**Description:** Create integration tests specifically for the git remote setup flow, including the security fixes from ISSUE-001.

**Required Test Coverage:**
1. Test complete remote setup flow with valid URLs
2. Test rejection of malicious URLs
3. Test behavior when remote already exists
4. Test --skip-remote flag
5. Test interactive prompts (mocked stdin)

**Test Requirements:**
- Use temporary directories for test repos
- Clean up test artifacts
- Mock stdin for automated testing
- Verify actual git commands execute correctly

**Dependencies:** ISSUE-001, ISSUE-003
**Blocks:** None

---

### ISSUE-006: Address Ignored Error Return Values
**Severity:** INFO
**File:** main.go (location TBD by grep search)
**Description:** UBS scan found 1 instance of assignment discarding secondary return values (likely an error). Review and handle appropriately.

**Required Changes:**
1. Search codebase for `_` used to discard error returns
2. For each instance, either:
   - Handle the error appropriately
   - Document why the error can be safely ignored (with comment)
   - Use explicit `//nolint` comment if intentional

**Example patterns to search:**
```go
result, _ := someFunc()  // Bad: ignoring error
_, err := someFunc()     // OK: intentionally ignoring non-error value
```

**Test Requirements:**
- No new tests required
- Ensure existing error handling tests still pass

**Dependencies:** None

---

## Implementation Order

Recommended implementation sequence to handle dependencies:

1. **ISSUE-002** - Update Go version (quick, no dependencies)
2. **ISSUE-003** - Add validation helpers (needed by ISSUE-001)
3. **ISSUE-001** - Fix command injection (critical, uses ISSUE-003)
4. **ISSUE-006** - Address ignored errors (quick review task)
5. **ISSUE-004** - Add unit tests (foundation for testing)
6. **ISSUE-005** - Add integration tests (validates ISSUE-001 fix)

---

## Testing Strategy

### Test Execution
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Coverage Targets
- Main package: > 70%
- Validation package: > 80%
- Overall project: > 65%

### CI Integration
After tests are implemented, add to pre-commit:
```bash
go test ./... || exit 1
ubs . --fail-on-warning || exit 1
```

---

## Beads Issue Creation Commands

```bash
# Critical
bd create --title "Fix command injection in git remote setup" --description "Validate and sanitize user input for git remote URL. See ISSUE-001 in spec."

# Warning
bd create --title "Update go.mod to Go 1.23" --description "Update go directive to 1.23+. See ISSUE-002 in spec."

# Quality
bd create --title "Create input validation helper package" --description "Add internal/validate package with URL and input validation. See ISSUE-003 in spec."

bd create --title "Add unit tests for main package" --description "Create main_test.go with >70% coverage. See ISSUE-004 in spec."

bd create --title "Add integration tests for git remote setup" --description "Test git remote flow including security fixes. See ISSUE-005 in spec."

bd create --title "Review and handle ignored error values" --description "Find and properly handle discarded error returns. See ISSUE-006 in spec."
```

---

## Success Criteria

Task completion checklist:
- [ ] All CRITICAL issues resolved
- [ ] All WARNING issues resolved
- [ ] Input validation package created and tested
- [ ] Unit test coverage > 70% for main package
- [ ] Integration tests pass
- [ ] All ignored errors reviewed and documented
- [ ] UBS scan shows 0 critical, 0 warnings
- [ ] `go test ./...` passes
- [ ] Build succeeds with Go 1.23

---

## Notes for AI Assistant

**Code Style:**
- Follow existing code conventions (already established in main.go)
- Use table-driven tests for Go tests
- Add comments for non-obvious security decisions
- Keep functions small and focused

**Validation Approach:**
- Prefer allowlist over blocklist for URL validation
- Use regex patterns for git URL formats
- Be strict: reject ambiguous input rather than trying to fix it
- Return clear error messages for rejected input

**Testing Approach:**
- Use `testing.T` for unit tests
- Use temporary directories for integration tests
- Mock external dependencies (stdin, exec.Command)
- Clean up test artifacts in defer statements
- Use subtests (t.Run) to organize test cases

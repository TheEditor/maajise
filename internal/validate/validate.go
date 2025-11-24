package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// MaxURLLength is the maximum allowed length for a git remote URL
const MaxURLLength = 2048

// shellMetachars contains characters that have special meaning in shells
var shellMetachars = []string{";", "&", "|", "$", "`", "\n", "\r", "<", ">", "(", ")", "{", "}"}

// gitURLPattern matches valid git remote URLs:
// - https://example.com/user/repo.git
// - git@example.com:user/repo.git
// - ssh://git@example.com/user/repo.git
var gitURLPattern = regexp.MustCompile(`^(https://|git@|ssh://).+`)

// ValidateGitURL checks if a string is a valid git remote URL.
// It uses an allowlist approach, accepting only:
// - HTTPS URLs (https://)
// - SSH URLs (git@hostname or ssh://)
// It rejects URLs containing shell metacharacters or exceeding max length.
func ValidateGitURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if len(url) > MaxURLLength {
		return fmt.Errorf("URL exceeds maximum length of %d characters", MaxURLLength)
	}

	if ContainsShellMetachars(url) {
		return fmt.Errorf("URL contains shell metacharacters")
	}

	if !gitURLPattern.MatchString(url) {
		return fmt.Errorf("invalid git URL format; must start with https://, git@, or ssh://")
	}

	return nil
}

// SanitizeInput trims whitespace and validates basic constraints.
// Returns the sanitized input or an error if constraints are violated.
func SanitizeInput(input string, maxLen int) (string, error) {
	sanitized := strings.TrimSpace(input)

	if sanitized == "" {
		return "", fmt.Errorf("input cannot be empty or whitespace only")
	}

	if len(sanitized) > maxLen {
		return "", fmt.Errorf("input exceeds maximum length of %d characters", maxLen)
	}

	// Reject multiline input
	if strings.Contains(sanitized, "\n") {
		return "", fmt.Errorf("multiline input not allowed")
	}

	return sanitized, nil
}

// ContainsShellMetachars returns true if the string contains any shell metacharacters.
// This is a safety check to prevent shell injection attacks.
func ContainsShellMetachars(s string) bool {
	for _, char := range shellMetachars {
		if strings.Contains(s, char) {
			return true
		}
	}
	return false
}

package validate

import (
	"testing"
)

func TestValidateGitURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		// Valid HTTPS URLs
		{
			name:    "valid github https url",
			url:     "https://github.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "valid gitlab https url",
			url:     "https://gitlab.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "https url without .git suffix",
			url:     "https://github.com/user/repo",
			wantErr: false,
		},

		// Valid SSH URLs
		{
			name:    "valid ssh git@ url",
			url:     "git@github.com:user/repo.git",
			wantErr: false,
		},
		{
			name:    "valid ssh git@ url without .git",
			url:     "git@github.com:user/repo",
			wantErr: false,
		},
		{
			name:    "valid ssh:// protocol url",
			url:     "ssh://git@github.com/user/repo.git",
			wantErr: false,
		},

		// Invalid URLs - shell metacharacters
		{
			name:    "url with semicolon",
			url:     "https://github.com/user/repo.git; echo hacked",
			wantErr: true,
		},
		{
			name:    "url with ampersand",
			url:     "https://github.com/user/repo.git & malicious",
			wantErr: true,
		},
		{
			name:    "url with pipe",
			url:     "https://github.com/user/repo.git | cat /etc/passwd",
			wantErr: true,
		},
		{
			name:    "url with dollar sign",
			url:     "https://github.com/user/$REPO.git",
			wantErr: true,
		},
		{
			name:    "url with backtick",
			url:     "https://github.com/user/`whoami`.git",
			wantErr: true,
		},
		{
			name:    "url with newline",
			url:     "https://github.com/user/repo.git\nmalicious",
			wantErr: true,
		},
		{
			name:    "url with carriage return",
			url:     "https://github.com/user/repo.git\r\nmalicious",
			wantErr: true,
		},
		{
			name:    "url with angle bracket",
			url:     "https://github.com/user/repo.git > /tmp/file",
			wantErr: true,
		},

		// Invalid URLs - format
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
		{
			name:    "url without valid scheme",
			url:     "github.com/user/repo.git",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://github.com/user/repo.git",
			wantErr: true,
		},
		{
			name:    "url exceeding max length",
			url:     "https://github.com/" + string(make([]byte, MaxURLLength)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGitURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGitURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLen    int
		want      string
		wantErr   bool
		wantEmpty bool
	}{
		// Valid inputs
		{
			name:    "simple input",
			input:   "hello",
			maxLen:  100,
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "input with leading whitespace",
			input:   "  hello",
			maxLen:  100,
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "input with trailing whitespace",
			input:   "hello  ",
			maxLen:  100,
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "input with both leading and trailing whitespace",
			input:   "  hello world  ",
			maxLen:  100,
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "input at max length",
			input:   "12345",
			maxLen:  5,
			want:    "12345",
			wantErr: false,
		},

		// Invalid inputs
		{
			name:      "empty input",
			input:     "",
			maxLen:    100,
			wantErr:   true,
			wantEmpty: true,
		},
		{
			name:      "whitespace only input",
			input:     "   ",
			maxLen:    100,
			wantErr:   true,
			wantEmpty: true,
		},
		{
			name:    "input exceeding max length",
			input:   "123456",
			maxLen:  5,
			wantErr: true,
		},
		{
			name:    "multiline input",
			input:   "hello\nworld",
			maxLen:  100,
			wantErr: true,
		},
		{
			name:    "tabs and newlines",
			input:   "hello\t\nworld",
			maxLen:  100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeInput(tt.input, tt.maxLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeInput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("SanitizeInput() = %q, want %q", got, tt.want)
			}
			if tt.wantEmpty && got != "" {
				t.Errorf("SanitizeInput() expected empty string, got %q", got)
			}
		})
	}
}

func TestContainsShellMetachars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		// Safe strings
		{
			name: "simple string",
			s:    "hello",
			want: false,
		},
		{
			name: "url without metacharacters",
			s:    "https://github.com/user/repo.git",
			want: false,
		},
		{
			name: "string with numbers and hyphens",
			s:    "my-project-123",
			want: false,
		},
		{
			name: "string with underscores and dots",
			s:    "my_project.tar.gz",
			want: false,
		},

		// Dangerous strings
		{
			name: "string with semicolon",
			s:    "test;echo",
			want: true,
		},
		{
			name: "string with ampersand",
			s:    "test&echo",
			want: true,
		},
		{
			name: "string with pipe",
			s:    "test|echo",
			want: true,
		},
		{
			name: "string with dollar sign",
			s:    "test$var",
			want: true,
		},
		{
			name: "string with backtick",
			s:    "test`cmd`",
			want: true,
		},
		{
			name: "string with newline",
			s:    "test\necho",
			want: true,
		},
		{
			name: "string with carriage return",
			s:    "test\recho",
			want: true,
		},
		{
			name: "string with redirect",
			s:    "test > output",
			want: true,
		},
		{
			name: "string with parentheses",
			s:    "test(echo)",
			want: true,
		},
		{
			name: "string with braces",
			s:    "test{a,b}",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsShellMetachars(tt.s)
			if got != tt.want {
				t.Errorf("ContainsShellMetachars(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

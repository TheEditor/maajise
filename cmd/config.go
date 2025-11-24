package cmd

// Config holds initialization configuration for project setup
type Config struct {
	ProjectName   string
	InPlace       bool
	NoOverwrite   bool
	SkipGit       bool
	SkipBeads     bool
	SkipCommit    bool
	SkipRemote    bool
	SkipGitUser   bool
	GitName       string
	GitEmail      string
	Verbose       bool
	Template      string
	MainBranch    string
	DefaultRemote string
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		InPlace:       false,
		NoOverwrite:   false,
		SkipGit:       false,
		SkipBeads:     false,
		SkipCommit:    false,
		SkipRemote:    false,
		SkipGitUser:   false,
		GitName:       "",
		GitEmail:      "",
		Verbose:       false,
		Template:      "",
		MainBranch:    "main",
		DefaultRemote: "",
	}
}

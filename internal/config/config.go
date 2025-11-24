package config

// Config holds runtime configuration for commands
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

// New creates a new config with defaults
func New() *Config {
	return &Config{
		Template: "base",
	}
}

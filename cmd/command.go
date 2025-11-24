package cmd

import "flag"

// Command represents a CLI command
type Command interface {
	Name() string
	Description() string
	Usage() string
	Examples() []string
	Execute(args []string) error
}

// Registry holds all registered commands
var registry = make(map[string]Command)

// Register adds a command to the registry
func Register(cmd Command) {
	registry[cmd.Name()] = cmd
}

// Get retrieves a command by name
func Get(name string) (Command, bool) {
	cmd, ok := registry[name]
	return cmd, ok
}

// All returns all registered commands
func All() []Command {
	cmds := make([]Command, 0, len(registry))
	for _, cmd := range registry {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// ParseFlags is a helper for commands that use flags
func ParseFlags(fs *flag.FlagSet, args []string) error {
	return fs.Parse(args)
}

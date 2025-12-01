package cmd

import (
	"flag"
	"fmt"
	"sort"

	"maajise/internal/ui"
)

type HelpCommand struct {
	fs *flag.FlagSet
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{
		fs: flag.NewFlagSet("help", flag.ContinueOnError),
	}
}

func (hc *HelpCommand) Name() string {
	return "help"
}

func (hc *HelpCommand) Description() string {
	return "Display help information about commands"
}

func (hc *HelpCommand) LongDescription() string {
	return `Show help information for Maajise commands.

Without arguments, displays a list of all available commands with brief descriptions.
With a command name, shows detailed help including usage syntax, description, and examples.`
}

func (hc *HelpCommand) Usage() string {
	return "maajise help [command]"
}

func (hc *HelpCommand) Examples() string {
	return `  # Show all commands
  maajise help

  # Show help for specific command
  maajise help init`
}

func (hc *HelpCommand) Run(args []string) error {
	if err := hc.fs.Parse(args); err != nil {
		return err
	}

	remainingArgs := hc.fs.Args()

	// If specific command requested, show detailed help
	if len(remainingArgs) > 0 {
		return hc.showCommandHelp(remainingArgs[0])
	}

	// Otherwise show general help
	hc.showGeneralHelp()
	return nil
}

func (hc *HelpCommand) showGeneralHelp() {
	fmt.Println("Maajise - Project Initialization Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  maajise <command> [arguments]")
	fmt.Println("  maajise <project-name>         (shorthand for 'init <project-name>')")
	fmt.Println()
	fmt.Println("Available commands:")

	// Get all commands and sort alphabetically
	commands := All()
	names := make([]string, 0, len(commands))
	for _, cmd := range commands {
		names = append(names, cmd.Name())
	}
	sort.Strings(names)

	// Display each command with description
	for _, name := range names {
		cmd, _ := Get(name)
		fmt.Printf("  %-12s %s\n", name, cmd.Description())
	}

	fmt.Println()
	fmt.Println("Use 'maajise help <command>' for more information about a command.")
}

func (hc *HelpCommand) showCommandHelp(cmdName string) error {
	cmd, ok := Get(cmdName)
	if !ok {
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	ui.Header(fmt.Sprintf("Command: %s", cmd.Name()))
	fmt.Println(cmd.LongDescription())
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", cmd.Usage())
	fmt.Println()

	// Show examples
	examples := cmd.Examples()
	if examples != "" {
		fmt.Println("Examples:")
		fmt.Println(examples)
	}

	return nil
}

func init() {
	Register(NewHelpCommand())
}

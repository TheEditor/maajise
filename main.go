package main

import (
	"fmt"
	"os"
	"strings"

	"maajise/cmd"
	"maajise/internal/ui"
)

const VERSION = "2.0.0"

func main() {
	// Handle no arguments - show help
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmdName := os.Args[1]

	// Handle help flags
	if cmdName == "-h" || cmdName == "--help" || cmdName == "help" {
		printUsage()
		os.Exit(0)
	}

	// Handle version flags
	if cmdName == "-v" || cmdName == "--version" || cmdName == "version" {
		if vc, ok := cmd.Get("version"); ok {
			vc.Run([]string{})
		} else {
			fmt.Printf("Maajise v%s\n", VERSION)
		}
		os.Exit(0)
	}

	// Try to get the command (case-insensitive lookup)
	command, ok := cmd.Get(strings.ToLower(cmdName))
	if !ok {
		// If not found, assume it's a project name and default to "init"
		// This provides the convenience shorthand: maajise my-project
		command, _ = cmd.Get("init")
		if command == nil {
			ui.Error(fmt.Sprintf("Unknown command: %s", cmdName))
			printUsage()
			os.Exit(1)
		}
		// Reconstruct args to include "init" as the command
		os.Args = append([]string{os.Args[0], "init"}, os.Args[1:]...)
	}

	// Check for --help or -h flag in remaining args
	for _, arg := range os.Args[2:] {
		if arg == "--help" || arg == "-h" {
			if helpCmd, ok := cmd.Get("help"); ok {
				helpCmd.Run([]string{cmdName})
				os.Exit(0)
			}
			break
		}
	}

	// Execute the command with remaining args
	if err := command.Run(os.Args[2:]); err != nil {
		ui.Error(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`Maajise v%s - Repository Initialization Tool

Usage:
  maajise <command> [flags] [arguments]
  maajise <project-name>  (convenience shorthand for 'init')

Commands:
`, VERSION)

	// Group commands logically
	groups := []struct {
		name     string
		commands []string
	}{
		{"Project Setup", []string{"init", "add", "update"}},
		{"Project Info", []string{"status", "validate", "templates"}},
		{"System", []string{"doctor"}},
		{"Help", []string{"help", "version"}},
	}

	for _, g := range groups {
		fmt.Printf("\n  %s:\n", g.name)
		for _, cmdName := range g.commands {
			if c, ok := cmd.Get(cmdName); ok {
				fmt.Printf("    %-12s  %s\n", c.Name(), c.Description())
			}
		}
	}

	fmt.Println()
	fmt.Println("Global Flags:")
	fmt.Println("  -h, --help     Show help")
	fmt.Println("  -v, --version  Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  maajise init my-project --template=typescript")
	fmt.Println("  maajise init --interactive")
	fmt.Println("  maajise add git")
	fmt.Println("  maajise doctor")
	fmt.Println("  maajise validate")
	fmt.Println()
	fmt.Println("Run 'maajise help <command>' for detailed help on a command.")
	fmt.Println()
}

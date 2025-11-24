package main

import (
	"fmt"
	"os"

	"maajise/cmd"
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
			vc.Execute([]string{})
		} else {
			fmt.Printf("Maajise v%s\n", VERSION)
		}
		os.Exit(0)
	}

	// Try to get the command
	command, ok := cmd.Get(cmdName)
	if !ok {
		// If not found, assume it's a project name and default to "init"
		// This provides the convenience shorthand: maajise my-project
		command, _ = cmd.Get("init")
		if command == nil {
			fmt.Fprintf(os.Stderr, "✗ Unknown command: %s\n", cmdName)
			printUsage()
			os.Exit(1)
		}
		// Reconstruct args to include "init" as the command
		os.Args = append([]string{os.Args[0], "init"}, os.Args[1:]...)
	}

	// Execute the command with remaining args
	if err := command.Execute(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
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
	for _, c := range cmd.All() {
		fmt.Printf("  %-12s  %s\n", c.Name(), c.Description())
	}
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -h, --help     Show this help")
	fmt.Println("  -v, --version  Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  maajise init my-project")
	fmt.Println("  maajise my-project  (same as above)")
	fmt.Println("  maajise init --in-place")
	fmt.Println("  maajise init my-project --no-git")
	fmt.Println()
}

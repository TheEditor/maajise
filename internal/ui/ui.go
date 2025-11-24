package ui

import (
	"fmt"
	"os"
)

// ANSI colors
const (
	Red    = "\033[0;31m"
	Green  = "\033[0;32m"
	Yellow = "\033[1;33m"
	Blue   = "\033[0;34m"
	Reset  = "\033[0m"
)

// Error prints an error message
func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s✗%s %s\n", Red, Reset, msg)
}

// Success prints a success message
func Success(msg string) {
	fmt.Printf("%s✓%s %s\n", Green, Reset, msg)
}

// Info prints an info message
func Info(msg string) {
	fmt.Printf("%s→%s %s\n", Blue, Reset, msg)
}

// Warn prints a warning message
func Warn(msg string) {
	fmt.Printf("%s⚠%s %s\n", Yellow, Reset, msg)
}

// Header prints a formatted header box
func Header(title string) {
	fmt.Println()
	fmt.Printf("%s╔═══════════════════════════════════════════════════════════╗%s\n", Blue, Reset)
	fmt.Printf("%s║  %-57s║%s\n", Blue, title, Reset)
	fmt.Printf("%s╚═══════════════════════════════════════════════════════════╝%s\n", Blue, Reset)
	fmt.Println()
}

// Summary prints a success summary box
func Summary(lines ...string) {
	fmt.Println()
	fmt.Printf("%s╔═══════════════════════════════════════════════════════════╗%s\n", Green, Reset)
	for _, line := range lines {
		fmt.Printf("%s║  %-57s║%s\n", Green, line, Reset)
	}
	fmt.Printf("%s╚═══════════════════════════════════════════════════════════╝%s\n", Green, Reset)
	fmt.Println()
}

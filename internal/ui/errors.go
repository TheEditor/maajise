package ui

import "fmt"

// UsageError returns formatted usage error with help hint
func UsageError(cmdName string, msg string) error {
	return fmt.Errorf("%s\n\nRun 'maajise %s --help' for usage information", msg, cmdName)
}

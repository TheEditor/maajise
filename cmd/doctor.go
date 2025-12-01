package cmd

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"maajise/internal/config"
	"maajise/internal/ui"
)

type DoctorCommand struct {
	fs      *flag.FlagSet
	verbose bool
}

type DependencyCheck struct {
	Name     string
	Command  string
	Args     []string
	Required bool
	Version  string
	Found    bool
	Error    string
}

func NewDoctorCommand() *DoctorCommand {
	dc := &DoctorCommand{
		fs: flag.NewFlagSet("doctor", flag.ContinueOnError),
	}
	dc.fs.BoolVar(&dc.verbose, "v", false, "Verbose output")
	dc.fs.BoolVar(&dc.verbose, "verbose", false, "Verbose output")
	return dc
}

func (dc *DoctorCommand) Name() string {
	return "doctor"
}

func (dc *DoctorCommand) Description() string {
	return "Check system dependencies and configuration"
}

func (dc *DoctorCommand) LongDescription() string {
	return `Check system dependencies and configuration for Maajise.

Verifies that required tools (Git, Beads) are installed and accessible, and optionally
checks for optional tools (UBS, Go). Helps diagnose setup issues and verify the system
is properly configured for using Maajise.`
}

func (dc *DoctorCommand) Usage() string {
	return "maajise doctor [flags]"
}

func (dc *DoctorCommand) Examples() string {
	return `  # Check dependencies
  maajise doctor

  # Verbose output with version details
  maajise doctor --verbose`
}

func (dc *DoctorCommand) Run(args []string) error {
	if err := dc.fs.Parse(args); err != nil {
		return err
	}

	ui.Info("Checking maajise dependencies...")
	fmt.Println()

	checks := []DependencyCheck{
		{Name: "git", Command: "git", Args: []string{"--version"}, Required: true},
		{Name: "bd (Beads)", Command: "bd", Args: []string{"--version"}, Required: true},
		{Name: "ubs", Command: "ubs", Args: []string{"--version"}, Required: false},
		{Name: "go", Command: "go", Args: []string{"version"}, Required: false},
	}

	allOK := true
	requiredMissing := false

	for i, check := range checks {
		checks[i] = dc.runCheck(check)
		c := checks[i]

		status := "✓"
		if !c.Found {
			if c.Required {
				status = "✗"
				requiredMissing = true
			} else {
				status = "○"
			}
			allOK = false
		}

		reqStr := ""
		if c.Required {
			reqStr = " (required)"
		} else {
			reqStr = " (optional)"
		}

		if c.Found {
			ui.Success(fmt.Sprintf("%s %s: %s", status, c.Name, c.Version))
		} else {
			if c.Required {
				ui.Error(fmt.Sprintf("%s %s: not found%s", status, c.Name, reqStr))
			} else {
				ui.Warn(fmt.Sprintf("%s %s: not found%s", status, c.Name, reqStr))
			}
		}

		if dc.verbose && c.Error != "" {
			fmt.Printf("    Error: %s\n", c.Error)
		}
	}

	// Check config file
	fmt.Println()
	ui.Info("Configuration:")
	if config.ConfigExists() {
		ui.Success(fmt.Sprintf("✓ Config file: %s", config.ConfigPath()))
		if dc.verbose {
			if fc, err := config.LoadFileConfig(); err == nil {
				if fc.Defaults.Template != "" {
					fmt.Printf("    Default template: %s\n", fc.Defaults.Template)
				}
				if fc.Defaults.GitName != "" {
					fmt.Printf("    Default git name: %s\n", fc.Defaults.GitName)
				}
			}
		}
	} else {
		ui.Warn(fmt.Sprintf("○ Config file: not found (%s)", config.ConfigPath()))
	}

	// Check custom templates directory
	if fc, err := config.LoadFileConfig(); err == nil && fc.TemplatesDir != "" {
		ui.Info(fmt.Sprintf("  Custom templates: %s", fc.TemplatesDir))
	}

	fmt.Println()

	if requiredMissing {
		return fmt.Errorf("required dependencies missing")
	}

	if allOK {
		ui.Success("All checks passed!")
	} else {
		ui.Info("Some optional dependencies missing (maajise will still work)")
	}

	return nil
}

func (dc *DoctorCommand) runCheck(check DependencyCheck) DependencyCheck {
	cmd := exec.Command(check.Command, check.Args...)
	output, err := cmd.Output()

	if err != nil {
		check.Found = false
		check.Error = err.Error()
		return check
	}

	check.Found = true
	check.Version = strings.TrimSpace(string(output))

	// Clean up version strings
	if strings.Contains(check.Version, "\n") {
		check.Version = strings.Split(check.Version, "\n")[0]
	}

	return check
}

func init() {
	Register(NewDoctorCommand())
}

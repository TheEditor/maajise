package cmd

import (
	"flag"
	"fmt"
)

const Version = "2.0.0"

type VersionCommand struct {
	fs *flag.FlagSet
}

func NewVersionCommand() *VersionCommand {
	return &VersionCommand{
		fs: flag.NewFlagSet("version", flag.ContinueOnError),
	}
}

func (vc *VersionCommand) Name() string {
	return "version"
}

func (vc *VersionCommand) Description() string {
	return "Display version information"
}

func (vc *VersionCommand) Usage() string {
	return "maajise version"
}

func (vc *VersionCommand) Examples() []string {
	return []string{
		"maajise version",
	}
}

func (vc *VersionCommand) FlagSet() *flag.FlagSet {
	return vc.fs
}

func (vc *VersionCommand) Execute(args []string) error {
	if err := vc.fs.Parse(args); err != nil {
		return err
	}

	fmt.Printf("Maajise version %s\n", Version)
	return nil
}

func init() {
	Register(NewVersionCommand())
}

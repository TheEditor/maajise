package cmd

import (
	"flag"
	"fmt"
	"sort"

	"maajise/templates"
)

// TemplatesCommand lists available templates
type TemplatesCommand struct {
	fs *flag.FlagSet
}

func NewTemplatesCommand() *TemplatesCommand {
	return &TemplatesCommand{
		fs: flag.NewFlagSet("templates", flag.ContinueOnError),
	}
}

func (tc *TemplatesCommand) Name() string {
	return "templates"
}

func (tc *TemplatesCommand) Description() string {
	return "List available project templates"
}

func (tc *TemplatesCommand) Usage() string {
	return "maajise templates"
}

func (tc *TemplatesCommand) Examples() []string {
	return []string{
		"maajise templates",
	}
}

func (tc *TemplatesCommand) Execute(args []string) error {
	if err := tc.fs.Parse(args); err != nil {
		return err
	}

	allTemplates := templates.All()

	// Sort by name for consistent output
	sort.Slice(allTemplates, func(i, j int) bool {
		return allTemplates[i].Name() < allTemplates[j].Name()
	})

	fmt.Println("Available templates:")
	fmt.Println()

	for _, tmpl := range allTemplates {
		fmt.Printf("  %-12s  %s\n", tmpl.Name(), tmpl.Description())
		if deps := tmpl.Dependencies(); len(deps) > 0 {
			fmt.Printf("               Dependencies: %v\n", deps)
		}
	}

	fmt.Println()
	fmt.Println("Usage: maajise init <project-name> --template=<template>")
	fmt.Println()

	return nil
}

func init() {
	Register(NewTemplatesCommand())
}

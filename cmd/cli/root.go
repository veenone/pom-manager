package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/pom-manager/cmd/cli/commands"
)

var (
	verbose bool
	noColor bool
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "pom-manager",
	Short: "Maven POM Manager - Create and manage Maven POM files",
	Long: `A CLI tool for creating, validating, and managing Maven POM files.

Supports template-based project creation, dependency management,
and POM validation following Maven conventions.`,
	Version: "0.1.0-MVP",
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")

	// Add subcommands
	rootCmd.AddCommand(commands.CreateCmd)
	rootCmd.AddCommand(commands.ValidateCmd)
	rootCmd.AddCommand(commands.AddDepCmd)
	rootCmd.AddCommand(commands.TemplatesCmd)
	rootCmd.AddCommand(commands.InfoCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

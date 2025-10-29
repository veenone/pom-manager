package commands

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/user/pom-manager/internal/core/pom"
)

var (
	jsonOutput bool
)

var InfoCmd = &cobra.Command{
	Use:   "info <file>",
	Short: "Display POM file information",
	Long:  `Display information about a Maven POM file including coordinates, dependencies, and plugins.`,
	Example: `  pom-manager info pom.xml
  pom-manager info --json pom.xml`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	InfoCmd.Flags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
}

func runInfo(cmd *cobra.Command, args []string) error {
	file := args[0]

	// Parse POM
	parser := pom.NewParser()
	project, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parsing POM: %w", err)
	}

	if jsonOutput {
		data, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	// Display info
	color.Cyan("=== POM Information ===\n")

	color.Green("Project:")
	fmt.Printf("  Group ID:    %s\n", project.GroupID)
	fmt.Printf("  Artifact ID: %s\n", project.ArtifactID)
	fmt.Printf("  Version:     %s\n", project.Version)
	fmt.Printf("  Packaging:   %s\n", project.Packaging)

	if project.Name != "" {
		fmt.Printf("  Name:        %s\n", project.Name)
	}

	if len(project.Dependencies) > 0 {
		color.Green("\nDependencies (%d):", len(project.Dependencies))
		for _, dep := range project.Dependencies {
			scope := dep.Scope
			if scope == "" {
				scope = "compile"
			}
			fmt.Printf("  - %s:%s:%s [%s]\n", dep.GroupID, dep.ArtifactID, dep.Version, scope)
		}
	}

	if project.Build != nil && len(project.Build.Plugins) > 0 {
		color.Green("\nPlugins (%d):", len(project.Build.Plugins))
		for _, plugin := range project.Build.Plugins {
			fmt.Printf("  - %s:%s", plugin.GroupID, plugin.ArtifactID)
			if plugin.Version != "" {
				fmt.Printf(":%s", plugin.Version)
			}
			fmt.Println()
		}
	}

	return nil
}

package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/user/pom-manager/internal/core/pom"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a Maven POM file",
	Long:  `Parse and validate a Maven POM file against Maven conventions.`,
	Example: `  pom-manager validate pom.xml
  pom-manager validate --verbose pom.xml`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	file := args[0]

	// Parse POM
	parser := pom.NewParser()
	project, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parsing POM: %w", err)
	}

	color.Cyan("Parsed: %s", project.Coordinates.String())

	// Validate
	validator := pom.NewValidator()
	result := validator.Validate(project)

	if result.Valid {
		color.Green("✓ POM is valid")
		return nil
	}

	// Print errors
	color.Red("✗ Validation failed:\n")

	if len(result.Errors.Coordinates) > 0 {
		color.Yellow("Coordinate Errors:")
		for _, err := range result.Errors.Coordinates {
			color.Red("  - %s", err.Error())
		}
	}

	if len(result.Errors.Dependencies) > 0 {
		color.Yellow("Dependency Errors:")
		for _, err := range result.Errors.Dependencies {
			color.Red("  - %s", err.Error())
		}
	}

	if len(result.Errors.Build) > 0 {
		color.Yellow("Build Errors:")
		for _, err := range result.Errors.Build {
			color.Red("  - %s", err.Error())
		}
	}

	if len(result.Errors.General) > 0 {
		color.Yellow("General Errors:")
		for _, err := range result.Errors.General {
			color.Red("  - %s", err.Error())
		}
	}

	return fmt.Errorf("validation failed")
}

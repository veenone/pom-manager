package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/user/pom-manager/internal/core/pom"
)

var (
	depGroup    string
	depArtifact string
	depVersion  string
	depScope    string
	depFile     string
)

var AddDepCmd = &cobra.Command{
	Use:   "add-dep",
	Short: "Add a dependency to a POM file",
	Long:  `Add a Maven dependency to an existing POM file.`,
	Example: `  pom-manager add-dep --group junit --artifact junit --version 4.13.2 --scope test
  pom-manager add-dep -g org.slf4j -a slf4j-api -v 2.0.0 --file myproject/pom.xml`,
	RunE: runAddDep,
}

func init() {
	AddDepCmd.Flags().StringVarP(&depGroup, "group", "g", "", "dependency groupId (required)")
	AddDepCmd.Flags().StringVarP(&depArtifact, "artifact", "a", "", "dependency artifactId (required)")
	AddDepCmd.Flags().StringVarP(&depVersion, "version", "V", "", "dependency version (required)")
	AddDepCmd.Flags().StringVarP(&depScope, "scope", "s", "compile", "dependency scope")
	AddDepCmd.Flags().StringVarP(&depFile, "file", "f", "pom.xml", "POM file to modify")
	AddDepCmd.MarkFlagRequired("group")
	AddDepCmd.MarkFlagRequired("artifact")
	AddDepCmd.MarkFlagRequired("version")
}

func runAddDep(cmd *cobra.Command, args []string) error {
	// Parse existing POM
	parser := pom.NewParser()
	project, err := parser.ParseFile(depFile)
	if err != nil {
		return fmt.Errorf("parsing POM: %w", err)
	}

	// Add dependency
	dep := pom.Dependency{
		GroupID:    depGroup,
		ArtifactID: depArtifact,
		Version:    depVersion,
		Scope:      depScope,
	}

	// Check if already exists
	exists := false
	for i, existing := range project.Dependencies {
		if existing.GroupID == dep.GroupID && existing.ArtifactID == dep.ArtifactID {
			project.Dependencies[i] = dep // Update version
			exists = true
			color.Yellow("Updated existing dependency")
			break
		}
	}

	if !exists {
		project.Dependencies = append(project.Dependencies, dep)
		color.Green("Added new dependency")
	}

	// Validate
	validator := pom.NewValidator()
	result := validator.Validate(project)
	if !result.Valid {
		color.Red("✗ Validation failed after adding dependency:")
		for _, err := range result.Errors.AllErrors() {
			color.Red("  - %s", err.Error())
		}
		return fmt.Errorf("validation failed")
	}

	// Write back
	generator := pom.NewGenerator()
	if err := generator.GenerateToFile(project, depFile); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	color.Green("✓ Dependency added to %s", depFile)
	fmt.Printf("  %s:%s:%s [%s]\n", depGroup, depArtifact, depVersion, depScope)

	return nil
}

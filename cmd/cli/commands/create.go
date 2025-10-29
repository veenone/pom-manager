package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/user/pom-manager/internal/core/pom"
)

var (
	groupID    string
	artifactID string
	version    string
	template   string
	output     string
	force      bool
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Maven POM file",
	Long:  `Create a new Maven POM file from a template or with custom coordinates.`,
	Example: `  # Interactive mode
  pom-manager create

  # Non-interactive with flags
  pom-manager create --group com.example --artifact my-app --version 1.0.0

  # With template
  pom-manager create --template java-library --group com.example --artifact my-lib --version 1.0.0`,
	RunE: runCreate,
}

func init() {
	CreateCmd.Flags().StringVarP(&groupID, "group", "g", "", "Maven groupId")
	CreateCmd.Flags().StringVarP(&artifactID, "artifact", "a", "", "Maven artifactId")
	CreateCmd.Flags().StringVarP(&version, "version", "V", "", "project version")
	CreateCmd.Flags().StringVarP(&template, "template", "t", "basic-java", "template name")
	CreateCmd.Flags().StringVarP(&output, "output", "o", "pom.xml", "output file path")
	CreateCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing file")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Check if file exists
	if !force {
		if _, err := os.Stat(output); err == nil {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("File %s already exists. Overwrite", output),
				IsConfirm: true,
			}
			result, err := prompt.Run()
			if err != nil || result != "y" {
				color.Yellow("Cancelled")
				return nil
			}
		}
	}

	// Interactive mode if coordinates not provided
	if groupID == "" || artifactID == "" || version == "" {
		if err := interactiveCreate(); err != nil {
			return err
		}
	}

	// Create coordinates
	coords := pom.Coordinates{
		GroupID:    groupID,
		ArtifactID: artifactID,
		Version:    version,
	}

	// Create project from template
	tm := pom.NewTemplateManager()
	project, err := tm.Create(template, coords)
	if err != nil {
		return fmt.Errorf("creating project: %w", err)
	}

	// Validate
	validator := pom.NewValidator()
	result := validator.Validate(project)
	if !result.Valid {
		color.Red("✗ Validation failed:")
		for _, err := range result.Errors.AllErrors() {
			color.Red("  - %s", err.Error())
		}
		return fmt.Errorf("project validation failed")
	}

	// Generate and write
	generator := pom.NewGenerator()
	if err := generator.GenerateToFile(project, output); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	color.Green("✓ Created POM file: %s", output)
	color.Cyan("  Group ID:    %s", project.GroupID)
	color.Cyan("  Artifact ID: %s", project.ArtifactID)
	color.Cyan("  Version:     %s", project.Version)
	color.Cyan("  Template:    %s", template)

	return nil
}

func interactiveCreate() error {
	color.Cyan("=== Create New Maven Project ===\n")

	// Select template
	tm := pom.NewTemplateManager()
	templates := tm.List()
	templateNames := make([]string, len(templates))
	for i, t := range templates {
		templateNames[i] = fmt.Sprintf("%s - %s", t.Name, t.Description)
	}

	templatePrompt := promptui.Select{
		Label: "Select Template",
		Items: templateNames,
	}
	idx, _, err := templatePrompt.Run()
	if err != nil {
		return err
	}
	template = templates[idx].Name

	// Get coordinates
	groupPrompt := promptui.Prompt{
		Label:   "Group ID",
		Default: "com.example",
	}
	groupID, err = groupPrompt.Run()
	if err != nil {
		return err
	}

	artifactPrompt := promptui.Prompt{
		Label:   "Artifact ID",
		Default: "my-app",
	}
	artifactID, err = artifactPrompt.Run()
	if err != nil {
		return err
	}

	versionPrompt := promptui.Prompt{
		Label:   "Version",
		Default: "1.0.0",
	}
	version, err = versionPrompt.Run()
	if err != nil {
		return err
	}

	fmt.Println()
	return nil
}

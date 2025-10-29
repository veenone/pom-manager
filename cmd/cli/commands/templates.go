package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/user/pom-manager/internal/core/pom"
)

var TemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available POM templates",
	Long:  `List all available Maven POM templates with descriptions.`,
	Example: `  pom-manager templates`,
	RunE: runTemplates,
}

func runTemplates(cmd *cobra.Command, args []string) error {
	tm := pom.NewTemplateManager()
	templates := tm.List()

	color.Cyan("Available POM Templates:\n")
	for _, t := range templates {
		color.Green("  %s", t.Name)
		fmt.Printf("    %s\n", t.Description)
	}

	return nil
}

package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// HelpDialog displays Maven concepts and help information
type HelpDialog struct {
	window fyne.Window
}

// NewHelpDialog creates a new help dialog
func NewHelpDialog(window fyne.Window) *HelpDialog {
	return &HelpDialog{
		window: window,
	}
}

// Show displays the Maven basics help dialog
func (h *HelpDialog) Show() {
	// Create sections for different Maven concepts
	sections := []struct {
		title   string
		content string
	}{
		{
			title: "POM Structure",
			content: `The Project Object Model (POM) is an XML file that contains information about the project and configuration details used by Maven.

Key Elements:
• <groupId>: Identifies the project's group or organization
• <artifactId>: Identifies the project's unique name
• <version>: Specifies the project version
• <packaging>: Defines output type (jar, war, pom)
• <name>: Human-readable project name
• <description>: Project description`,
		},
		{
			title: "Dependencies",
			content: `Dependencies are external libraries your project needs to compile and run.

Dependency Coordinates:
• groupId: Organization that created the library
• artifactId: Name of the library
• version: Library version to use
• scope: When the dependency is needed (compile, test, runtime)

Common Scopes:
• compile: Available in all classpaths (default)
• test: Only for testing
• provided: Provided by runtime environment
• runtime: Not needed for compilation`,
		},
		{
			title: "Plugins",
			content: `Maven plugins provide goals that extend Maven's build functionality.

Common Plugins:
• maven-compiler-plugin: Compiles Java source code
• maven-jar-plugin: Creates JAR files
• maven-war-plugin: Creates WAR files for web apps
• maven-surefire-plugin: Runs unit tests

Plugin Configuration:
Plugins can be configured with <configuration> elements to customize their behavior.`,
		},
		{
			title: "Lifecycle Phases",
			content: `Maven builds follow a standard lifecycle with predefined phases.

Default Lifecycle Phases:
• validate: Validate project structure
• compile: Compile source code
• test: Run unit tests
• package: Package compiled code (JAR/WAR)
• verify: Run integration tests
• install: Install to local repository
• deploy: Deploy to remote repository

Running Phases:
Execute: mvn <phase>
Example: mvn clean install`,
		},
		{
			title: "Properties",
			content: `Properties define reusable values throughout the POM.

Usage:
• Define: <properties><key>value</key></properties>
• Reference: ${key}

Common Properties:
• project.version: Project version
• maven.compiler.source: Java source version
• maven.compiler.target: Java target version
• project.build.sourceEncoding: Source file encoding

Custom Properties:
You can define any custom properties for your project and reference them in dependencies, plugins, or other configuration.`,
		},
	}

	// Create accordion with sections
	var accordionItems []*widget.AccordionItem
	for _, section := range sections {
		// Need to capture section in closure
		sectionCopy := section

		// Create multiline label for content
		contentLabel := widget.NewLabel(sectionCopy.content)
		contentLabel.Wrapping = fyne.TextWrapWord

		// Create scrollable content
		scrollContent := container.NewVScroll(contentLabel)
		scrollContent.SetMinSize(fyne.NewSize(600, 150))

		// Create accordion item
		item := widget.NewAccordionItem(
			sectionCopy.title,
			scrollContent,
		)
		accordionItems = append(accordionItems, item)
	}

	// Create accordion
	accordion := widget.NewAccordion(accordionItems...)

	// Open first section by default
	if len(accordionItems) > 0 {
		accordion.Open(0)
	}

	// Create main content
	content := container.NewBorder(
		widget.NewLabel("Maven Basics - Quick Reference"),
		nil, nil, nil,
		container.NewVScroll(accordion),
	)

	// Create dialog
	customDialog := dialog.NewCustom(
		"Maven Basics",
		"Close",
		content,
		h.window,
	)

	customDialog.Resize(fyne.NewSize(700, 600))
	customDialog.Show()
}

// ShowQuickHelp shows a quick help message with keyboard shortcuts
func (h *HelpDialog) ShowQuickHelp() {
	helpText := `Maven POM Manager - Quick Help

Keyboard Shortcuts:
• Ctrl+N - Create new POM file
• Ctrl+O - Open existing POM file
• Ctrl+S - Save current file
• Ctrl+Shift+S - Save as new file
• Ctrl+W / Ctrl+Q - Quit application
• F1 - Show this help
• F5 - Refresh and validate

Getting Started:
1. Create a new POM using File → New or Ctrl+N
2. Fill in project coordinates (Group ID, Artifact ID, Version)
3. Choose a template (basic-java, java-library, or web-app)
4. Add dependencies and plugins as needed
5. Save your POM file using Ctrl+S

For more information about Maven concepts, see Help → Maven Basics.`

	dialog.ShowInformation("Quick Help", helpText, h.window)
}

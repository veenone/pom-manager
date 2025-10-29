package panels

import (
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
	"github.com/user/pom-manager/internal/gui/widgets"
)

// LifecyclePanel displays plugin executions organized by Maven lifecycle phase
type LifecyclePanel struct {
	// UI components
	accordion       *widget.Accordion
	addButton       *widgets.ButtonWithTooltip
	mainContainer   *fyne.Container

	// State
	organizer pom.Organizer
	project   *pom.Project
	phaseMap  map[string][]pom.PluginExecution

	// Callbacks
	onAddExecution    func(pluginIndex int, execution pom.PluginExecution)
	onRemoveExecution func(pluginIndex int, executionID string)
}

// NewLifecyclePanel creates a new LifecyclePanel
func NewLifecyclePanel() *LifecyclePanel {
	panel := &LifecyclePanel{
		organizer: pom.NewOrganizer(),
		phaseMap:  make(map[string][]pom.PluginExecution),
	}

	panel.createUI()
	return panel
}

// createUI creates the panel layout
func (p *LifecyclePanel) createUI() {
	// Create accordion for phases
	p.accordion = widget.NewAccordion()

	// Add Execution button
	p.addButton = widgets.NewButtonWithTooltip("Add Execution",
		"Add a new plugin execution and bind it to a lifecycle phase",
		func() {
			if p.onAddExecution != nil && p.project != nil && p.project.Build != nil && len(p.project.Build.Plugins) > 0 {
				// Callback will be triggered from parent
			}
		})

	buttonBar := container.NewHBox(p.addButton)

	// Create main container with title and buttons
	p.mainContainer = container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Maven Lifecycle Phases"),
			widget.NewSeparator(),
			widget.NewLabel("Plugin executions grouped by lifecycle phase"),
		),
		buttonBar,
		nil, nil,
		container.NewScroll(p.accordion),
	)
}

// LoadProject updates the panel with project data
func (p *LifecyclePanel) LoadProject(project *pom.Project) {
	p.project = project
	p.refresh()
}

// refresh rebuilds the accordion with current project data
func (p *LifecyclePanel) refresh() {
	if p.project == nil {
		// UI updates must be called on UI thread
		fyne.Do(func() {
			p.accordion.Items = nil
			p.accordion.Refresh()
		})
		return
	}

	// Get plugin executions organized by phase
	p.phaseMap = p.organizer.ByPhase(p.project)

	// Get phase order
	phaseOrder := p.organizer.GetPhaseOrder()

	// Build accordion items in phase order
	var items []*widget.AccordionItem

	for _, phase := range phaseOrder {
		executions, hasExecutions := p.phaseMap[phase]
		if !hasExecutions || len(executions) == 0 {
			continue
		}

		// Create content for this phase
		phaseContent := p.createPhaseContent(phase, executions)

		// Create accordion item
		title := fmt.Sprintf("%s (%d)", phase, len(executions))
		item := widget.NewAccordionItem(title, phaseContent)
		items = append(items, item)
	}

	// Check for executions with no phase (should be rare)
	if len(p.phaseMap) == 0 && p.project.Build != nil && len(p.project.Build.Plugins) > 0 {
		// Show message that no phase-bound executions exist
		noPhaseLabel := widget.NewLabel("No plugin executions bound to lifecycle phases.\nPlugin executions need a <phase> element to appear here.")
		noPhaseLabel.Wrapping = fyne.TextWrapWord
		emptyItem := widget.NewAccordionItem("ℹ️ Information", noPhaseLabel)
		items = append(items, emptyItem)
	}

	if len(items) == 0 {
		// Show empty state
		emptyLabel := widget.NewLabel("No plugin executions configured.\nAdd plugins with executions to see them organized by phase.")
		emptyLabel.Wrapping = fyne.TextWrapWord
		emptyItem := widget.NewAccordionItem("ℹ️ Getting Started", emptyLabel)
		items = append(items, emptyItem)
	}

	// UI updates must be called on UI thread
	// Completely recreate accordion to ensure proper expansion behavior
	fyne.Do(func() {
		// Clear all existing items
		for len(p.accordion.Items) > 0 {
			p.accordion.Remove(p.accordion.Items[0])
		}

		// Add new items
		for _, item := range items {
			p.accordion.Append(item)
		}

		p.accordion.Refresh()
	})
}

// createPhaseContent creates the content for a single phase section
func (p *LifecyclePanel) createPhaseContent(phase string, executions []pom.PluginExecution) fyne.CanvasObject {
	// Sort executions by plugin name for consistent display
	sortedExecs := make([]pom.PluginExecution, len(executions))
	copy(sortedExecs, executions)
	sort.Slice(sortedExecs, func(i, j int) bool {
		return sortedExecs[i].ID < sortedExecs[j].ID
	})

	var widgets []fyne.CanvasObject

	// Phase description
	phaseDesc := widget.NewLabel(getPhaseDescription(phase))
	phaseDesc.Wrapping = fyne.TextWrapWord
	phaseDesc.TextStyle = fyne.TextStyle{Italic: true}
	widgets = append(widgets, phaseDesc, widget.NewSeparator())

	// List each execution
	for i, exec := range sortedExecs {
		execCard := p.createExecutionCard(i+1, exec)
		widgets = append(widgets, execCard)
	}

	return container.NewVBox(widgets...)
}

// createExecutionCard creates a card displaying a single plugin execution
func (p *LifecyclePanel) createExecutionCard(index int, exec pom.PluginExecution) fyne.CanvasObject {
	// Execution ID
	idLabel := widget.NewLabel(fmt.Sprintf("%d. Execution ID: %s", index, exec.ID))
	idLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Goals
	goalsText := "Goals: " + formatGoals(exec.Goals)
	goalsLabel := widget.NewLabel(goalsText)

	// Build card content
	cardContent := container.NewVBox(
		idLabel,
		goalsLabel,
	)

	// Configuration info (if present)
	if exec.Configuration != nil {
		configLabel := widget.NewLabel("✓ Has configuration")
		configLabel.TextStyle = fyne.TextStyle{Italic: true}
		cardContent.Add(configLabel)
	}

	// Create card with padding
	card := container.NewPadded(cardContent)

	return card
}

// formatGoals formats a list of goals into a comma-separated string
func formatGoals(goals []string) string {
	if len(goals) == 0 {
		return "(none)"
	}
	result := ""
	for i, goal := range goals {
		if i > 0 {
			result += ", "
		}
		result += goal
	}
	return result
}

// getPhaseDescription returns a description for each Maven lifecycle phase
func getPhaseDescription(phase string) string {
	descriptions := map[string]string{
		pom.PhaseValidate:              "Validate project structure and configuration",
		pom.PhaseInitialize:            "Initialize build state",
		pom.PhaseGenerateSources:       "Generate source code",
		pom.PhaseProcessSources:        "Process source code",
		pom.PhaseGenerateResources:     "Generate resources",
		pom.PhaseProcessResources:      "Copy and process resources to output directory",
		pom.PhaseCompile:               "Compile source code",
		pom.PhaseProcessClasses:        "Post-process compiled classes",
		pom.PhaseGenerateTestSources:   "Generate test source code",
		pom.PhaseProcessTestSources:    "Process test source code",
		pom.PhaseGenerateTestResources: "Generate test resources",
		pom.PhaseProcessTestResources:  "Copy and process test resources",
		pom.PhaseTestCompile:           "Compile test source code",
		pom.PhaseProcessTestClasses:    "Post-process compiled test classes",
		pom.PhaseTest:                  "Run unit tests",
		pom.PhasePreparePackage:        "Prepare for packaging",
		pom.PhasePackage:               "Package compiled code (JAR, WAR, etc.)",
		pom.PhasePreIntegrationTest:    "Prepare integration test environment",
		pom.PhaseIntegrationTest:       "Run integration tests",
		pom.PhasePostIntegrationTest:   "Clean up integration test environment",
		pom.PhaseVerify:                "Verify package is valid",
		pom.PhaseInstall:               "Install package to local repository",
		pom.PhaseDeploy:                "Deploy package to remote repository",
	}

	if desc, found := descriptions[phase]; found {
		return desc
	}
	return "Custom phase"
}

// GetContainer returns the main container for embedding
func (p *LifecyclePanel) GetContainer() *fyne.Container {
	return p.mainContainer
}

// Clear clears the panel
func (p *LifecyclePanel) Clear() {
	p.project = nil
	p.phaseMap = make(map[string][]pom.PluginExecution)
	p.accordion.Items = nil
	p.accordion.Refresh()
}

// OnAddExecution sets the callback for adding an execution
func (p *LifecyclePanel) OnAddExecution(callback func(pluginIndex int, execution pom.PluginExecution)) {
	p.onAddExecution = callback
	// Update button callback
	p.addButton.OnTapped = func() {
		if callback != nil && p.project != nil && p.project.Build != nil && len(p.project.Build.Plugins) > 0 {
			callback(-1, pom.PluginExecution{}) // -1 means show add dialog
		}
	}
}

// OnRemoveExecution sets the callback for removing an execution
func (p *LifecyclePanel) OnRemoveExecution(callback func(pluginIndex int, executionID string)) {
	p.onRemoveExecution = callback
}

// GetProject returns the current project
func (p *LifecyclePanel) GetProject() *pom.Project {
	return p.project
}

package dialogs

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
)

// ExecutionDialog manages plugin execution creation and editing
type ExecutionDialog struct {
	window  fyne.Window
	plugins []pom.Plugin

	// Form widgets
	pluginSelect *widget.Select
	executionID  *widget.Entry
	phaseSelect  *widget.Select
	goalsEntry   *widget.Entry
}

// NewExecutionDialog creates a new ExecutionDialog
func NewExecutionDialog(window fyne.Window, plugins []pom.Plugin) *ExecutionDialog {
	return &ExecutionDialog{
		window:  window,
		plugins: plugins,
	}
}

// ShowAdd displays the dialog for adding a new execution
func (d *ExecutionDialog) ShowAdd(callback func(pluginIndex int, execution pom.PluginExecution)) {
	d.show("Add Plugin Execution", "", pom.PluginExecution{}, func(pluginIndex int, exec pom.PluginExecution) {
		if callback != nil {
			callback(pluginIndex, exec)
		}
	})
}

// ShowEdit displays the dialog for editing an existing execution
func (d *ExecutionDialog) ShowEdit(pluginIndex int, execution pom.PluginExecution, callback func(pluginIndex int, execution pom.PluginExecution)) {
	pluginKey := d.getPluginKey(pluginIndex)
	d.show("Edit Plugin Execution", pluginKey, execution, callback)
}

// show displays the execution dialog
func (d *ExecutionDialog) show(title string, preselectedPlugin string, existing pom.PluginExecution, callback func(int, pom.PluginExecution)) {
	// Plugin selection
	pluginOptions := d.getPluginOptions()
	if len(pluginOptions) == 0 {
		dialog.ShowError(nil, d.window)
		return
	}

	d.pluginSelect = widget.NewSelect(pluginOptions, nil)
	if preselectedPlugin != "" {
		d.pluginSelect.SetSelected(preselectedPlugin)
	} else {
		d.pluginSelect.SetSelected(pluginOptions[0])
	}

	// Execution ID
	d.executionID = widget.NewEntry()
	d.executionID.SetPlaceHolder("default-execution")
	if existing.ID != "" {
		d.executionID.SetText(existing.ID)
	}

	// Phase selection
	phases := []string{
		pom.PhaseValidate,
		pom.PhaseInitialize,
		pom.PhaseGenerateSources,
		pom.PhaseProcessSources,
		pom.PhaseGenerateResources,
		pom.PhaseProcessResources,
		pom.PhaseCompile,
		pom.PhaseProcessClasses,
		pom.PhaseGenerateTestSources,
		pom.PhaseProcessTestSources,
		pom.PhaseGenerateTestResources,
		pom.PhaseProcessTestResources,
		pom.PhaseTestCompile,
		pom.PhaseProcessTestClasses,
		pom.PhaseTest,
		pom.PhasePreparePackage,
		pom.PhasePackage,
		pom.PhasePreIntegrationTest,
		pom.PhaseIntegrationTest,
		pom.PhasePostIntegrationTest,
		pom.PhaseVerify,
		pom.PhaseInstall,
		pom.PhaseDeploy,
	}

	d.phaseSelect = widget.NewSelect(phases, nil)
	if existing.Phase != "" {
		d.phaseSelect.SetSelected(existing.Phase)
	} else {
		d.phaseSelect.SetSelected(pom.PhaseCompile)
	}

	// Goals entry (comma-separated)
	d.goalsEntry = widget.NewEntry()
	d.goalsEntry.SetPlaceHolder("compile, testCompile")
	if len(existing.Goals) > 0 {
		d.goalsEntry.SetText(strings.Join(existing.Goals, ", "))
	}

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Plugin", Widget: d.pluginSelect},
			{Text: "Execution ID", Widget: d.executionID},
			{Text: "Phase", Widget: d.phaseSelect},
			{Text: "Goals (comma-separated)", Widget: d.goalsEntry},
		},
	}

	// Info label
	infoLabel := widget.NewLabel("Select a plugin and configure when its goals should execute during the build lifecycle.")
	infoLabel.Wrapping = fyne.TextWrapWord
	infoLabel.TextStyle = fyne.TextStyle{Italic: true}

	content := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		form,
	)

	// Create dialog
	customDialog := dialog.NewCustomConfirm(
		title,
		"Save",
		"Cancel",
		content,
		func(save bool) {
			if save {
				// Validate and create execution
				pluginIndex := d.getSelectedPluginIndex()
				if pluginIndex < 0 {
					return
				}

				exec := pom.PluginExecution{
					ID:    d.executionID.Text,
					Phase: d.phaseSelect.Selected,
					Goals: d.parseGoals(d.goalsEntry.Text),
				}

				// Default execution ID if empty
				if exec.ID == "" {
					exec.ID = "default"
				}

				if callback != nil {
					callback(pluginIndex, exec)
				}
			}
		},
		d.window,
	)

	customDialog.Resize(fyne.NewSize(500, 400))
	customDialog.Show()
}

// getPluginOptions returns a list of plugin display strings
func (d *ExecutionDialog) getPluginOptions() []string {
	var options []string
	for _, plugin := range d.plugins {
		options = append(options, d.formatPlugin(plugin))
	}
	return options
}

// formatPlugin formats a plugin for display
func (d *ExecutionDialog) formatPlugin(plugin pom.Plugin) string {
	if plugin.Version != "" {
		return plugin.GroupID + ":" + plugin.ArtifactID + ":" + plugin.Version
	}
	return plugin.GroupID + ":" + plugin.ArtifactID
}

// getPluginKey returns the key for a specific plugin index
func (d *ExecutionDialog) getPluginKey(index int) string {
	if index >= 0 && index < len(d.plugins) {
		return d.formatPlugin(d.plugins[index])
	}
	return ""
}

// getSelectedPluginIndex returns the index of the selected plugin
func (d *ExecutionDialog) getSelectedPluginIndex() int {
	selected := d.pluginSelect.Selected
	for i, plugin := range d.plugins {
		if d.formatPlugin(plugin) == selected {
			return i
		}
	}
	return -1
}

// parseGoals parses a comma-separated string of goals
func (d *ExecutionDialog) parseGoals(goalsStr string) []string {
	if goalsStr == "" {
		return []string{}
	}

	parts := strings.Split(goalsStr, ",")
	var goals []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			goals = append(goals, trimmed)
		}
	}
	return goals
}

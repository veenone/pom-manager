package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
)

// PluginDialog is a modal dialog for adding or editing plugins
type PluginDialog struct {
	window fyne.Window

	// Form fields
	commonPluginSelect *widget.Select
	groupIDEntry       *widget.Entry
	artifactIDEntry    *widget.Entry
	versionEntry       *widget.Entry

	// Callbacks
	onSave func(pom.Plugin)
}

// Common Maven plugins
var commonPlugins = map[string]struct {
	GroupID    string
	ArtifactID string
	Version    string
}{
	"Maven Compiler Plugin": {"org.apache.maven.plugins", "maven-compiler-plugin", "3.11.0"},
	"Maven JAR Plugin":      {"org.apache.maven.plugins", "maven-jar-plugin", "3.3.0"},
	"Maven WAR Plugin":      {"org.apache.maven.plugins", "maven-war-plugin", "3.4.0"},
	"Maven Surefire Plugin": {"org.apache.maven.plugins", "maven-surefire-plugin", "3.1.2"},
	"Maven Assembly Plugin": {"org.apache.maven.plugins", "maven-assembly-plugin", "3.6.0"},
}

// NewPluginDialog creates a new plugin dialog
func NewPluginDialog(window fyne.Window) *PluginDialog {
	return &PluginDialog{
		window: window,
	}
}

// ShowAdd displays the dialog for adding a new plugin
func (d *PluginDialog) ShowAdd(callback func(pom.Plugin)) {
	d.onSave = callback
	d.show("Add Plugin", nil)
}

// ShowEdit displays the dialog for editing an existing plugin
func (d *PluginDialog) ShowEdit(plugin pom.Plugin, callback func(pom.Plugin)) {
	d.onSave = callback
	d.show("Edit Plugin", &plugin)
}

// show creates and displays the dialog
func (d *PluginDialog) show(title string, existingPlugin *pom.Plugin) {
	// Create common plugins dropdown
	pluginNames := []string{"(Custom)", "Maven Compiler Plugin", "Maven JAR Plugin", "Maven WAR Plugin", "Maven Surefire Plugin", "Maven Assembly Plugin"}
	d.commonPluginSelect = widget.NewSelect(pluginNames, func(selected string) {
		if selected != "(Custom)" {
			if plugin, ok := commonPlugins[selected]; ok {
				d.groupIDEntry.SetText(plugin.GroupID)
				d.artifactIDEntry.SetText(plugin.ArtifactID)
				d.versionEntry.SetText(plugin.Version)
			}
		}
	})
	d.commonPluginSelect.SetSelected("(Custom)")

	// Create form fields
	d.groupIDEntry = widget.NewEntry()
	d.groupIDEntry.SetPlaceHolder("org.apache.maven.plugins")

	d.artifactIDEntry = widget.NewEntry()
	d.artifactIDEntry.SetPlaceHolder("maven-compiler-plugin")

	d.versionEntry = widget.NewEntry()
	d.versionEntry.SetPlaceHolder("3.11.0")

	// Populate fields if editing
	if existingPlugin != nil {
		d.groupIDEntry.SetText(existingPlugin.GroupID)
		d.artifactIDEntry.SetText(existingPlugin.ArtifactID)
		d.versionEntry.SetText(existingPlugin.Version)
	}

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Common Plugins", Widget: d.commonPluginSelect},
			{Text: "Group ID", Widget: d.groupIDEntry},
			{Text: "Artifact ID", Widget: d.artifactIDEntry},
			{Text: "Version", Widget: d.versionEntry},
		},
	}

	// Create dialog
	content := container.NewVBox(form)

	customDialog := dialog.NewCustomConfirm(
		title,
		"Save",
		"Cancel",
		content,
		func(save bool) {
			if save && d.onSave != nil {
				plugin := pom.Plugin{
					GroupID:    d.groupIDEntry.Text,
					ArtifactID: d.artifactIDEntry.Text,
					Version:    d.versionEntry.Text,
				}
				d.onSave(plugin)
			}
		},
		d.window,
	)

	customDialog.Resize(fyne.NewSize(450, 280))
	customDialog.Show()
}

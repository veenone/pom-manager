package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
)

// DependencyDialog is a modal dialog for adding or editing dependencies
type DependencyDialog struct {
	window fyne.Window

	// Form fields
	groupIDEntry    *widget.Entry
	artifactIDEntry *widget.Entry
	versionEntry    *widget.Entry
	scopeSelect     *widget.Select

	// Callbacks
	onSave   func(pom.Dependency)
	onCancel func()
}

// NewDependencyDialog creates a new dependency dialog
func NewDependencyDialog(window fyne.Window) *DependencyDialog {
	return &DependencyDialog{
		window: window,
	}
}

// ShowAdd displays the dialog for adding a new dependency
func (d *DependencyDialog) ShowAdd(callback func(pom.Dependency)) {
	d.onSave = callback
	d.show("Add Dependency", nil)
}

// ShowEdit displays the dialog for editing an existing dependency
func (d *DependencyDialog) ShowEdit(dep pom.Dependency, callback func(pom.Dependency)) {
	d.onSave = callback
	d.show("Edit Dependency", &dep)
}

// show creates and displays the dialog
func (d *DependencyDialog) show(title string, existingDep *pom.Dependency) {
	// Create form fields
	d.groupIDEntry = widget.NewEntry()
	d.groupIDEntry.SetPlaceHolder("org.example")

	d.artifactIDEntry = widget.NewEntry()
	d.artifactIDEntry.SetPlaceHolder("library-name")

	d.versionEntry = widget.NewEntry()
	d.versionEntry.SetPlaceHolder("1.0.0")

	d.scopeSelect = widget.NewSelect(
		[]string{"compile", "test", "provided", "runtime", "system"},
		nil,
	)
	d.scopeSelect.SetSelected("compile")

	// Populate fields if editing
	if existingDep != nil {
		d.groupIDEntry.SetText(existingDep.GroupID)
		d.artifactIDEntry.SetText(existingDep.ArtifactID)
		d.versionEntry.SetText(existingDep.Version)
		if existingDep.Scope != "" {
			d.scopeSelect.SetSelected(existingDep.Scope)
		}
	}

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Group ID", Widget: d.groupIDEntry},
			{Text: "Artifact ID", Widget: d.artifactIDEntry},
			{Text: "Version", Widget: d.versionEntry},
			{Text: "Scope", Widget: d.scopeSelect},
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
				dep := pom.Dependency{
					GroupID:    d.groupIDEntry.Text,
					ArtifactID: d.artifactIDEntry.Text,
					Version:    d.versionEntry.Text,
					Scope:      d.scopeSelect.Selected,
				}
				d.onSave(dep)
			}
		},
		d.window,
	)

	customDialog.Resize(fyne.NewSize(400, 250))
	customDialog.Show()
}

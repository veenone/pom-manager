package wizard

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
	"github.com/user/pom-manager/internal/gui/widgets"
)

// CreateWizard is a multi-step wizard for creating new POM files
type CreateWizard struct {
	window fyne.Window

	// Step 1: Coordinates
	groupIDEntry    *widget.Entry
	artifactIDEntry *widget.Entry
	versionEntry    *widget.Entry
	packagingSelect *widget.Select

	// Step 2: Template selection
	templateSelect *widget.RadioGroup
	templateDesc   *widget.Label

	// Wizard state
	currentStep int
	maxSteps    int

	// Callbacks
	onComplete func(coords pom.Coordinates, template string)
	onCancel   func()
}

// NewCreateWizard creates a new project creation wizard
func NewCreateWizard(window fyne.Window) *CreateWizard {
	return &CreateWizard{
		window:      window,
		currentStep: 1,
		maxSteps:    2,
	}
}

// Show displays the wizard
func (w *CreateWizard) Show(onComplete func(pom.Coordinates, string)) {
	w.onComplete = onComplete
	w.showStep1()
}

// showStep1 displays Step 1: Project Coordinates
func (w *CreateWizard) showStep1() {
	// Create form fields
	w.groupIDEntry = widget.NewEntry()
	w.groupIDEntry.SetPlaceHolder("com.example")

	w.artifactIDEntry = widget.NewEntry()
	w.artifactIDEntry.SetPlaceHolder("my-app")

	w.versionEntry = widget.NewEntry()
	w.versionEntry.SetPlaceHolder("1.0.0")
	w.versionEntry.SetText("1.0.0") // Default

	w.packagingSelect = widget.NewSelect(
		[]string{"jar", "war", "pom", "maven-plugin"},
		nil,
	)
	w.packagingSelect.SetSelected("jar")

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Group ID *", Widget: w.groupIDEntry},
			{Text: "Artifact ID *", Widget: w.artifactIDEntry},
			{Text: "Version *", Widget: w.versionEntry},
			{Text: "Packaging", Widget: w.packagingSelect},
		},
	}

	content := container.NewVBox(
		widget.NewLabel("Step 1 of 2: Project Coordinates"),
		widget.NewSeparator(),
		form,
	)

	customDialog := dialog.NewCustomConfirm(
		"New POM Project",
		"Next",
		"Cancel",
		content,
		func(next bool) {
			if next {
				// Validate coordinates
				if w.groupIDEntry.Text == "" || w.artifactIDEntry.Text == "" || w.versionEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("all required fields must be filled"), w.window)
					w.showStep1() // Show again
					return
				}
				w.showStep2()
			}
		},
		w.window,
	)

	customDialog.Resize(fyne.NewSize(450, 300))
	customDialog.Show()
}

// showStep2 displays Step 2: Template Selection
func (w *CreateWizard) showStep2() {
	// Template options
	templates := []string{
		"basic-java",
		"java-library",
		"web-app",
		"javacard",
	}

	descriptions := map[string]string{
		"basic-java":   "Basic Java JAR project with compiler plugin",
		"java-library": "Java library project with compiler and JAR plugins",
		"web-app":      "Java web application (WAR) project",
		"javacard":     "JavaCard applet project for smart cards (CAP packaging)",
	}

	w.templateSelect = widget.NewRadioGroup(templates, func(selected string) {
		if desc, ok := descriptions[selected]; ok {
			w.templateDesc.SetText(desc)
		}
	})
	w.templateSelect.SetSelected("basic-java")

	w.templateDesc = widget.NewLabel(descriptions["basic-java"])
	w.templateDesc.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		widget.NewLabel("Step 2 of 2: Choose Template"),
		widget.NewSeparator(),
		widget.NewLabel("Select a project template:"),
		w.templateSelect,
		widget.NewSeparator(),
		widget.NewLabel("Description:"),
		w.templateDesc,
	)

	// Create dialog variable to reference in button callbacks
	var customDialog dialog.Dialog

	// Create custom buttons with tooltips that will reference the dialog
	backButton := widgets.NewButtonWithTooltip("Back",
		"Go back to project coordinates step",
		func() {
			if customDialog != nil {
				customDialog.Hide()
				w.showStep1()
			}
		})

	finishButton := widgets.NewButtonWithTooltip("Finish",
		"Create the project with the selected template",
		func() {
			if customDialog != nil {
				customDialog.Hide()
				if w.onComplete != nil {
					coords := pom.Coordinates{
						GroupID:    w.groupIDEntry.Text,
						ArtifactID: w.artifactIDEntry.Text,
						Version:    w.versionEntry.Text,
					}
					w.onComplete(coords, w.templateSelect.Selected)
				}
			}
		})

	buttonBar := container.NewHBox(
		backButton,
		finishButton,
	)

	// Build the complete content with buttons BEFORE creating dialog
	finalContent := container.NewBorder(
		nil, buttonBar, nil, nil,
		content,
	)

	// Now create the dialog with the complete content
	customDialog = dialog.NewCustom(
		"New POM Project",
		"Cancel",
		finalContent,
		w.window,
	)

	customDialog.Resize(fyne.NewSize(450, 350))
	customDialog.Show()
}

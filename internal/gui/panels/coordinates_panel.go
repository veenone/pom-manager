package panels

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
)

// CoordinatesPanel provides a form for editing project coordinates
type CoordinatesPanel struct {
	// Form fields
	groupIDEntry    *widget.Entry
	artifactIDEntry *widget.Entry
	versionEntry    *widget.Entry
	packagingSelect *widget.Select
	nameEntry       *widget.Entry
	descriptionEntry *widget.Entry

	// Main container
	mainContainer *fyne.Container

	// Callbacks
	onChange func(pom.Coordinates)

	// State
	loading bool // Flag to prevent onChange during programmatic updates
}

// NewCoordinatesPanel creates a new CoordinatesPanel
func NewCoordinatesPanel() *CoordinatesPanel {
	panel := &CoordinatesPanel{}
	panel.createUI()
	panel.setupCallbacks()
	return panel
}

// createUI creates the form layout
func (p *CoordinatesPanel) createUI() {
	// Create entry fields
	p.groupIDEntry = widget.NewEntry()
	p.groupIDEntry.SetPlaceHolder("com.example")

	p.artifactIDEntry = widget.NewEntry()
	p.artifactIDEntry.SetPlaceHolder("my-app")

	p.versionEntry = widget.NewEntry()
	p.versionEntry.SetPlaceHolder("1.0.0")

	// Packaging type selector
	p.packagingSelect = widget.NewSelect(
		[]string{"jar", "war", "pom", "maven-plugin"},
		nil,
	)
	p.packagingSelect.SetSelected("jar") // Default

	p.nameEntry = widget.NewEntry()
	p.nameEntry.SetPlaceHolder("My Application")

	p.descriptionEntry = widget.NewEntry()
	p.descriptionEntry.SetPlaceHolder("Project description")

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Group ID *", Widget: p.groupIDEntry},
			{Text: "Artifact ID *", Widget: p.artifactIDEntry},
			{Text: "Version *", Widget: p.versionEntry},
			{Text: "Packaging", Widget: p.packagingSelect},
			{Text: "Name", Widget: p.nameEntry},
			{Text: "Description", Widget: p.descriptionEntry},
		},
	}

	p.mainContainer = container.NewVBox(
		widget.NewLabel("Project Coordinates"),
		widget.NewSeparator(),
		form,
	)
}

// setupCallbacks sets up change callbacks for all fields
func (p *CoordinatesPanel) setupCallbacks() {
	p.groupIDEntry.OnChanged = func(s string) {
		p.notifyChange()
	}
	p.artifactIDEntry.OnChanged = func(s string) {
		p.notifyChange()
	}
	p.versionEntry.OnChanged = func(s string) {
		p.notifyChange()
	}
	p.packagingSelect.OnChanged = func(s string) {
		p.notifyChange()
	}
	p.nameEntry.OnChanged = func(s string) {
		p.notifyChange()
	}
	p.descriptionEntry.OnChanged = func(s string) {
		p.notifyChange()
	}
}

// notifyChange triggers the onChange callback with current coordinates
func (p *CoordinatesPanel) notifyChange() {
	// Don't notify if we're loading data programmatically
	if p.loading {
		return
	}

	if p.onChange != nil {
		coords := p.GetCoordinates()
		p.onChange(coords)
	}
}

// LoadCoordinates populates the form with coordinates
func (p *CoordinatesPanel) LoadCoordinates(coords pom.Coordinates) {
	// UI updates must be called on UI thread
	fyne.Do(func() {
		p.loading = true
		p.groupIDEntry.SetText(coords.GroupID)
		p.artifactIDEntry.SetText(coords.ArtifactID)
		p.versionEntry.SetText(coords.Version)
		p.loading = false
	})
}

// LoadProject populates the form from a full project
func (p *CoordinatesPanel) LoadProject(project *pom.Project) {
	if project == nil {
		return
	}

	// UI updates must be called on UI thread
	fyne.Do(func() {
		p.loading = true
		p.groupIDEntry.SetText(project.GroupID)
		p.artifactIDEntry.SetText(project.ArtifactID)
		p.versionEntry.SetText(project.Version)

		if project.Packaging != "" {
			p.packagingSelect.SetSelected(project.Packaging)
		} else {
			p.packagingSelect.SetSelected("jar")
		}

		p.nameEntry.SetText(project.Name)
		p.descriptionEntry.SetText(project.Description)
		p.loading = false
	})
}

// GetCoordinates returns the current coordinates from the form
func (p *CoordinatesPanel) GetCoordinates() pom.Coordinates {
	return pom.Coordinates{
		GroupID:    p.groupIDEntry.Text,
		ArtifactID: p.artifactIDEntry.Text,
		Version:    p.versionEntry.Text,
	}
}

// GetPackaging returns the selected packaging type
func (p *CoordinatesPanel) GetPackaging() string {
	return p.packagingSelect.Selected
}

// GetName returns the project name
func (p *CoordinatesPanel) GetName() string {
	return p.nameEntry.Text
}

// GetDescription returns the project description
func (p *CoordinatesPanel) GetDescription() string {
	return p.descriptionEntry.Text
}

// OnChange sets the callback for when coordinates change
func (p *CoordinatesPanel) OnChange(callback func(pom.Coordinates)) {
	p.onChange = callback
}

// GetContainer returns the main container for embedding
func (p *CoordinatesPanel) GetContainer() *fyne.Container {
	return p.mainContainer
}

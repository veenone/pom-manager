package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
)

// TemplateGallery displays available templates in a gallery view
type TemplateGallery struct {
	window fyne.Window

	// UI components
	templateList *widget.List
	previewText  *widget.Entry

	// State
	templates       []pom.TemplateInfo
	selectedIndex   int
	templateManager pom.TemplateManager

	// Callbacks
	onSelect func(templateName string)
}

// NewTemplateGallery creates a new template gallery
func NewTemplateGallery(window fyne.Window, templateManager pom.TemplateManager) *TemplateGallery {
	return &TemplateGallery{
		window:          window,
		templateManager: templateManager,
		selectedIndex:   -1,
	}
}

// Show displays the template gallery
func (g *TemplateGallery) Show(callback func(templateName string)) {
	g.onSelect = callback
	g.templates = g.templateManager.List()

	// Create template list
	g.templateList = widget.NewList(
		func() int {
			return len(g.templates)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Template Name"),
				widget.NewLabel("Description"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			nameLabel := box.Objects[0].(*widget.Label)
			descLabel := box.Objects[1].(*widget.Label)

			template := g.templates[id]
			nameLabel.SetText("📄 " + template.Name)
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			descLabel.SetText(template.Description)
		},
	)

	g.templateList.OnSelected = func(id widget.ListItemID) {
		g.selectedIndex = int(id)
		g.updatePreview()
	}

	// Create preview pane
	g.previewText = widget.NewMultiLineEntry()
	g.previewText.Disable()
	g.previewText.SetPlaceHolder("Select a template to preview...")

	// Create layout
	content := container.NewHSplit(
		g.templateList,
		container.NewBorder(
			widget.NewLabel("Preview:"),
			nil, nil, nil,
			container.NewScroll(g.previewText),
		),
	)
	content.SetOffset(0.4) // 40% for list, 60% for preview

	// Create dialog
	customDialog := dialog.NewCustomConfirm(
		"Template Gallery",
		"Use Template",
		"Cancel",
		content,
		func(useTemplate bool) {
			if useTemplate && g.selectedIndex >= 0 && g.onSelect != nil {
				template := g.templates[g.selectedIndex]
				g.onSelect(template.Name)
			}
		},
		g.window,
	)

	customDialog.Resize(fyne.NewSize(700, 500))
	customDialog.Show()
}

// updatePreview updates the preview pane with selected template
func (g *TemplateGallery) updatePreview() {
	if g.selectedIndex < 0 || g.selectedIndex >= len(g.templates) {
		return
	}

	template := g.templates[g.selectedIndex]

	// Generate a sample POM from the template
	sampleCoords := pom.Coordinates{
		GroupID:    "com.example",
		ArtifactID: "sample-project",
		Version:    "1.0.0",
	}

	project, err := g.templateManager.Create(template.Name, sampleCoords)
	if err != nil {
		g.previewText.SetText("Error generating preview: " + err.Error())
		return
	}

	// Generate XML for preview
	generator := pom.NewGenerator()
	xmlData, err := generator.Generate(project)
	if err != nil {
		g.previewText.SetText("Error generating XML: " + err.Error())
		return
	}

	g.previewText.SetText(string(xmlData))
}

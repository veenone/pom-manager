package panels

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
	"github.com/user/pom-manager/internal/gui/widgets"
)

// DependenciesPanel provides interface for managing project dependencies
type DependenciesPanel struct {
	// UI components
	dependenciesList *widget.List
	addButton        *widgets.ButtonWithTooltip
	editButton       *widgets.ButtonWithTooltip
	removeButton     *widgets.ButtonWithTooltip
	mainContainer    *fyne.Container

	// State
	dependencies     []pom.Dependency
	selectedIndex    int

	// Callbacks
	onAdd    func()
	onEdit   func(pom.Dependency)
	onRemove func(pom.Dependency)
}

// NewDependenciesPanel creates a new DependenciesPanel
func NewDependenciesPanel() *DependenciesPanel {
	panel := &DependenciesPanel{
		dependencies:  make([]pom.Dependency, 0),
		selectedIndex: -1,
	}

	panel.createUI()
	return panel
}

// createUI creates the panel layout
func (p *DependenciesPanel) createUI() {
	// Create list
	p.dependenciesList = widget.NewList(
		func() int {
			return len(p.dependencies)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			dep := p.dependencies[id]
			scope := dep.Scope
			if scope == "" {
				scope = "compile"
			}
			label.SetText(fmt.Sprintf("%s:%s:%s [%s]",
				dep.GroupID, dep.ArtifactID, dep.Version, scope))
		},
	)

	p.dependenciesList.OnSelected = func(id widget.ListItemID) {
		p.selectedIndex = int(id)
		p.updateButtonStates()
	}

	p.dependenciesList.OnUnselected = func(id widget.ListItemID) {
		p.selectedIndex = -1
		p.updateButtonStates()
	}

	// Create buttons with tooltips
	p.addButton = widgets.NewButtonWithTooltip("Add Dependency",
		"Add a new Maven dependency to the project",
		func() {
			if p.onAdd != nil {
				p.onAdd()
			}
		})

	p.editButton = widgets.NewButtonWithTooltip("Edit",
		"Edit the selected dependency",
		func() {
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.dependencies) && p.onEdit != nil {
				p.onEdit(p.dependencies[p.selectedIndex])
			}
		})
	p.editButton.Disable()

	p.removeButton = widgets.NewButtonWithTooltip("Remove",
		"Remove the selected dependency from the project",
		func() {
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.dependencies) && p.onRemove != nil {
				p.onRemove(p.dependencies[p.selectedIndex])
			}
		})
	p.removeButton.Disable()

	// Create layout
	buttonBar := container.NewHBox(
		p.addButton,
		p.editButton,
		p.removeButton,
	)

	p.mainContainer = container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Dependencies"),
			widget.NewSeparator(),
		),
		buttonBar,
		nil, nil,
		p.dependenciesList,
	)
}

// LoadDependencies updates the list with dependencies
func (p *DependenciesPanel) LoadDependencies(deps []pom.Dependency) {
	p.dependencies = deps
	// UI updates must be called on UI thread
	fyne.Do(func() {
		p.dependenciesList.Refresh()
		p.selectedIndex = -1
		p.updateButtonStates()
	})
}

// updateButtonStates enables/disables buttons based on selection
func (p *DependenciesPanel) updateButtonStates() {
	hasSelection := p.selectedIndex >= 0 && p.selectedIndex < len(p.dependencies)
	if hasSelection {
		p.editButton.Enable()
		p.removeButton.Enable()
	} else {
		p.editButton.Disable()
		p.removeButton.Disable()
	}
}

// OnAdd sets the callback for adding a dependency
func (p *DependenciesPanel) OnAdd(callback func()) {
	p.onAdd = callback
}

// OnEdit sets the callback for editing a dependency
func (p *DependenciesPanel) OnEdit(callback func(pom.Dependency)) {
	p.onEdit = callback
}

// OnRemove sets the callback for removing a dependency
func (p *DependenciesPanel) OnRemove(callback func(pom.Dependency)) {
	p.onRemove = callback
}

// GetContainer returns the main container for embedding
func (p *DependenciesPanel) GetContainer() *fyne.Container {
	return p.mainContainer
}

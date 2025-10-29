package panels

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
	"github.com/user/pom-manager/internal/gui/widgets"
)

// PluginsPanel provides interface for managing build plugins
type PluginsPanel struct {
	// UI components
	pluginsList   *widget.List
	addButton     *widgets.ButtonWithTooltip
	editButton    *widgets.ButtonWithTooltip
	removeButton  *widgets.ButtonWithTooltip
	mainContainer *fyne.Container

	// State
	plugins       []pom.Plugin
	selectedIndex int

	// Callbacks
	onAdd    func()
	onEdit   func(pom.Plugin)
	onRemove func(pom.Plugin)
}

// NewPluginsPanel creates a new PluginsPanel
func NewPluginsPanel() *PluginsPanel {
	panel := &PluginsPanel{
		plugins:       make([]pom.Plugin, 0),
		selectedIndex: -1,
	}

	panel.createUI()
	return panel
}

// createUI creates the panel layout
func (p *PluginsPanel) createUI() {
	// Create list
	p.pluginsList = widget.NewList(
		func() int {
			return len(p.plugins)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			plugin := p.plugins[id]
			if plugin.Version != "" {
				label.SetText(fmt.Sprintf("%s:%s:%s",
					plugin.GroupID, plugin.ArtifactID, plugin.Version))
			} else {
				label.SetText(fmt.Sprintf("%s:%s",
					plugin.GroupID, plugin.ArtifactID))
			}
		},
	)

	p.pluginsList.OnSelected = func(id widget.ListItemID) {
		p.selectedIndex = int(id)
		p.updateButtonStates()
	}

	p.pluginsList.OnUnselected = func(id widget.ListItemID) {
		p.selectedIndex = -1
		p.updateButtonStates()
	}

	// Create buttons with tooltips
	p.addButton = widgets.NewButtonWithTooltip("Add Plugin",
		"Add a new Maven build plugin to the project",
		func() {
			if p.onAdd != nil {
				p.onAdd()
			}
		})

	p.editButton = widgets.NewButtonWithTooltip("Edit",
		"Edit the selected build plugin",
		func() {
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.plugins) && p.onEdit != nil {
				p.onEdit(p.plugins[p.selectedIndex])
			}
		})
	p.editButton.Disable()

	p.removeButton = widgets.NewButtonWithTooltip("Remove",
		"Remove the selected build plugin from the project",
		func() {
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.plugins) && p.onRemove != nil {
				p.onRemove(p.plugins[p.selectedIndex])
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
			widget.NewLabel("Build Plugins"),
			widget.NewSeparator(),
		),
		buttonBar,
		nil, nil,
		p.pluginsList,
	)
}

// LoadPlugins updates the list with plugins
func (p *PluginsPanel) LoadPlugins(plugins []pom.Plugin) {
	p.plugins = plugins
	// UI updates must be called on UI thread
	fyne.Do(func() {
		p.pluginsList.Refresh()
		p.selectedIndex = -1
		p.updateButtonStates()
	})
}

// updateButtonStates enables/disables buttons based on selection
func (p *PluginsPanel) updateButtonStates() {
	hasSelection := p.selectedIndex >= 0 && p.selectedIndex < len(p.plugins)
	if hasSelection {
		p.editButton.Enable()
		p.removeButton.Enable()
	} else {
		p.editButton.Disable()
		p.removeButton.Disable()
	}
}

// OnAdd sets the callback for adding a plugin
func (p *PluginsPanel) OnAdd(callback func()) {
	p.onAdd = callback
}

// OnEdit sets the callback for editing a plugin
func (p *PluginsPanel) OnEdit(callback func(pom.Plugin)) {
	p.onEdit = callback
}

// OnRemove sets the callback for removing a plugin
func (p *PluginsPanel) OnRemove(callback func(pom.Plugin)) {
	p.onRemove = callback
}

// GetContainer returns the main container for embedding
func (p *PluginsPanel) GetContainer() *fyne.Container {
	return p.mainContainer
}

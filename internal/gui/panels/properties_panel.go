package panels

import (
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/gui/widgets"
)

// PropertiesPanel provides interface for managing Maven properties
type PropertiesPanel struct {
	// UI components
	propertiesList *widget.List
	addButton      *widgets.ButtonWithTooltip
	editButton     *widgets.ButtonWithTooltip
	removeButton   *widgets.ButtonWithTooltip
	mainContainer  *fyne.Container

	// State
	properties    map[string]string
	propertyKeys  []string
	selectedIndex int

	// Parent window for dialogs
	window fyne.Window

	// Callbacks
	onChange func(map[string]string)
}

// NewPropertiesPanel creates a new PropertiesPanel
func NewPropertiesPanel(window fyne.Window) *PropertiesPanel {
	panel := &PropertiesPanel{
		properties:    make(map[string]string),
		propertyKeys:  make([]string, 0),
		selectedIndex: -1,
		window:        window,
	}

	panel.createUI()
	return panel
}

// createUI creates the panel layout
func (p *PropertiesPanel) createUI() {
	// Create list
	p.propertiesList = widget.NewList(
		func() int {
			return len(p.propertyKeys)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			key := p.propertyKeys[id]
			value := p.properties[key]
			label.SetText(fmt.Sprintf("%s = %s", key, value))
		},
	)

	p.propertiesList.OnSelected = func(id widget.ListItemID) {
		p.selectedIndex = int(id)
		p.updateButtonStates()
	}

	p.propertiesList.OnUnselected = func(id widget.ListItemID) {
		p.selectedIndex = -1
		p.updateButtonStates()
	}

	// Create buttons with tooltips
	p.addButton = widgets.NewButtonWithTooltip("Add Property",
		"Add a new Maven property (key-value pair)",
		func() {
			p.showPropertyDialog("", "")
		})

	p.editButton = widgets.NewButtonWithTooltip("Edit",
		"Edit the selected property",
		func() {
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.propertyKeys) {
				key := p.propertyKeys[p.selectedIndex]
				value := p.properties[key]
				p.showPropertyDialog(key, value)
			}
		})
	p.editButton.Disable()

	p.removeButton = widgets.NewButtonWithTooltip("Remove",
		"Remove the selected property from the project",
		func() {
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.propertyKeys) {
				key := p.propertyKeys[p.selectedIndex]
				delete(p.properties, key)
				p.rebuildKeys()
				p.propertiesList.Refresh()
				p.selectedIndex = -1
				p.updateButtonStates()
				p.notifyChange()
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
			widget.NewLabel("Maven Properties"),
			widget.NewSeparator(),
		),
		buttonBar,
		nil, nil,
		p.propertiesList,
	)
}

// showPropertyDialog shows a dialog for adding or editing a property
func (p *PropertiesPanel) showPropertyDialog(existingKey, existingValue string) {
	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("property.name")
	if existingKey != "" {
		keyEntry.SetText(existingKey)
		keyEntry.Disable() // Don't allow key editing
	}

	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("value")
	valueEntry.SetText(existingValue)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Property Name", Widget: keyEntry},
			{Text: "Value", Widget: valueEntry},
		},
	}

	title := "Add Property"
	if existingKey != "" {
		title = "Edit Property"
	}

	customDialog := dialog.NewCustomConfirm(
		title,
		"Save",
		"Cancel",
		form,
		func(save bool) {
			if save {
				key := keyEntry.Text
				value := valueEntry.Text

				if key != "" {
					// Remove old key if renamed (not applicable since we disable key editing)
					if existingKey != "" && existingKey != key {
						delete(p.properties, existingKey)
					}

					p.properties[key] = value
					p.rebuildKeys()
					p.propertiesList.Refresh()
					p.notifyChange()
				}
			}
		},
		p.window,
	)

	customDialog.Resize(fyne.NewSize(400, 200))
	customDialog.Show()
}

// LoadProperties updates the panel with properties
func (p *PropertiesPanel) LoadProperties(props map[string]string) {
	p.properties = make(map[string]string)
	for k, v := range props {
		p.properties[k] = v
	}
	p.rebuildKeys()
	// UI updates must be called on UI thread
	fyne.Do(func() {
		p.propertiesList.Refresh()
		p.selectedIndex = -1
		p.updateButtonStates()
	})
}

// GetProperties returns the current properties
func (p *PropertiesPanel) GetProperties() map[string]string {
	result := make(map[string]string)
	for k, v := range p.properties {
		result[k] = v
	}
	return result
}

// rebuildKeys rebuilds the propertyKeys slice from the map
// Keys are sorted alphabetically to maintain consistent display order
func (p *PropertiesPanel) rebuildKeys() {
	p.propertyKeys = make([]string, 0, len(p.properties))
	for k := range p.properties {
		p.propertyKeys = append(p.propertyKeys, k)
	}
	// Sort keys alphabetically for consistent order
	sort.Strings(p.propertyKeys)
}

// updateButtonStates enables/disables buttons based on selection
func (p *PropertiesPanel) updateButtonStates() {
	hasSelection := p.selectedIndex >= 0 && p.selectedIndex < len(p.propertyKeys)
	if hasSelection {
		p.editButton.Enable()
		p.removeButton.Enable()
	} else {
		p.editButton.Disable()
		p.removeButton.Disable()
	}
}

// OnChange sets the callback for when properties change
func (p *PropertiesPanel) OnChange(callback func(map[string]string)) {
	p.onChange = callback
}

// notifyChange triggers the onChange callback
func (p *PropertiesPanel) notifyChange() {
	if p.onChange != nil {
		p.onChange(p.GetProperties())
	}
}

// GetContainer returns the main container for embedding
func (p *PropertiesPanel) GetContainer() *fyne.Container {
	return p.mainContainer
}

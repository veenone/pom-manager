package panels

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
)

// ErrorsPanel displays validation errors grouped by category
type ErrorsPanel struct {
	// UI components
	errorsList    *widget.List
	mainContainer *fyne.Container

	// State
	errors        []errorItem
	visible       bool

	// Callbacks
	onErrorClick func(errorType string, index int)
}

// errorItem represents a single error with category
type errorItem struct {
	category string
	message  string
	index    int
}

// NewErrorsPanel creates a new ErrorsPanel
func NewErrorsPanel() *ErrorsPanel {
	panel := &ErrorsPanel{
		errors:  make([]errorItem, 0),
		visible: false,
	}

	panel.createUI()
	return panel
}

// createUI creates the panel layout
func (p *ErrorsPanel) createUI() {
	// Create error list
	p.errorsList = widget.NewList(
		func() int {
			return len(p.errors)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ErrorIcon()),
				widget.NewLabel("template"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[1].(*widget.Label)
			err := p.errors[id]
			label.SetText(fmt.Sprintf("[%s] %s", err.category, err.message))
		},
	)

	p.errorsList.OnSelected = func(id widget.ListItemID) {
		if p.onErrorClick != nil && int(id) < len(p.errors) {
			err := p.errors[id]
			p.onErrorClick(err.category, err.index)
		}
	}

	// Wrap errors list in a container with minimum height
	scrolledList := container.NewScroll(p.errorsList)
	scrolledList.SetMinSize(fyne.NewSize(0, 150)) // Minimum 150px height

	p.mainContainer = container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Validation Errors"),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		scrolledList,
	)
}

// SetErrors updates the panel with validation errors
func (p *ErrorsPanel) SetErrors(result pom.ValidationResult) {
	p.errors = make([]errorItem, 0)

	if result.Valid {
		p.visible = false
		// UI updates must be called on UI thread
		fyne.Do(func() {
			p.errorsList.Refresh()
		})
		return
	}

	// Add coordinate errors
	for i, err := range result.Errors.Coordinates {
		p.errors = append(p.errors, errorItem{
			category: "Coordinates",
			message:  err.Error(),
			index:    i,
		})
	}

	// Add dependency errors
	for i, err := range result.Errors.Dependencies {
		p.errors = append(p.errors, errorItem{
			category: "Dependencies",
			message:  err.Error(),
			index:    i,
		})
	}

	// Add build errors
	for i, err := range result.Errors.Build {
		p.errors = append(p.errors, errorItem{
			category: "Build",
			message:  err.Error(),
			index:    i,
		})
	}

	// Add general errors
	for i, err := range result.Errors.General {
		p.errors = append(p.errors, errorItem{
			category: "General",
			message:  err.Error(),
			index:    i,
		})
	}

	p.visible = len(p.errors) > 0
	// UI updates must be called on UI thread
	fyne.Do(func() {
		p.errorsList.Refresh()
	})
}

// Clear clears all errors
func (p *ErrorsPanel) Clear() {
	p.errors = make([]errorItem, 0)
	p.visible = false
	p.errorsList.Refresh()
}

// IsVisible returns whether the panel should be visible
func (p *ErrorsPanel) IsVisible() bool {
	return p.visible
}

// OnErrorClick sets the callback for error clicks
func (p *ErrorsPanel) OnErrorClick(callback func(errorType string, index int)) {
	p.onErrorClick = callback
}

// GetContainer returns the main container for embedding
func (p *ErrorsPanel) GetContainer() *fyne.Container {
	return p.mainContainer
}

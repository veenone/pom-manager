package widgets

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ValidationBadge is a widget that displays validation status
// Shows green checkmark for valid, red X with error count for invalid
type ValidationBadge struct {
	widget.BaseWidget
	valid      bool
	errorCount int
	onClick    func()

	// Internal widgets
	statusIcon  *canvas.Text
	statusLabel *widget.Label
	container   *fyne.Container
}

// NewValidationBadge creates a new ValidationBadge
func NewValidationBadge() *ValidationBadge {
	badge := &ValidationBadge{
		valid:      true,
		errorCount: 0,
	}

	badge.ExtendBaseWidget(badge)
	badge.createUI()
	badge.updateDisplay()

	return badge
}

// createUI creates the internal UI components
func (v *ValidationBadge) createUI() {
	v.statusIcon = canvas.NewText("✓", theme.SuccessColor())
	v.statusIcon.TextSize = 16
	v.statusIcon.TextStyle = fyne.TextStyle{Bold: true}

	v.statusLabel = widget.NewLabel("Valid")

	v.container = container.NewHBox(
		v.statusIcon,
		v.statusLabel,
	)
}

// SetValid updates the validation status
func (v *ValidationBadge) SetValid(valid bool) {
	v.valid = valid
	v.updateDisplay()
	v.Refresh()
}

// SetErrorCount sets the number of validation errors
func (v *ValidationBadge) SetErrorCount(count int) {
	v.errorCount = count
	v.updateDisplay()
	v.Refresh()
}

// OnClick sets the callback for when the badge is clicked
func (v *ValidationBadge) OnClick(callback func()) {
	v.onClick = callback
}

// updateDisplay updates the visual appearance based on status
func (v *ValidationBadge) updateDisplay() {
	if v.valid {
		v.statusIcon.Text = "✓"
		v.statusIcon.Color = theme.SuccessColor()
		v.statusLabel.SetText("Valid")
	} else {
		v.statusIcon.Text = "✗"
		v.statusIcon.Color = theme.ErrorColor()
		if v.errorCount > 0 {
			v.statusLabel.SetText(fmt.Sprintf("Invalid (%d errors)", v.errorCount))
		} else {
			v.statusLabel.SetText("Invalid")
		}
	}
}

// CreateRenderer creates the renderer for this widget
func (v *ValidationBadge) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.container)
}

// Tapped handles tap/click events
func (v *ValidationBadge) Tapped(*fyne.PointEvent) {
	if v.onClick != nil {
		v.onClick()
	}
}

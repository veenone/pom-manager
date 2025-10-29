package widgets

import (
	"fyne.io/fyne/v2/widget"
)

// ButtonWithTooltip creates a button with tooltip support
type ButtonWithTooltip struct {
	widget.Button
	tooltip string
}

// NewButtonWithTooltip creates a new button with a tooltip
func NewButtonWithTooltip(label string, tooltip string, tapped func()) *ButtonWithTooltip {
	btn := &ButtonWithTooltip{
		tooltip: tooltip,
	}
	btn.Text = label
	btn.OnTapped = tapped
	btn.ExtendBaseWidget(btn)
	return btn
}

// TooltipText returns the tooltip text (implements fyne.Tooltipable)
func (b *ButtonWithTooltip) TooltipText() string {
	return b.tooltip
}

// SetTooltip updates the tooltip text
func (b *ButtonWithTooltip) SetTooltip(tooltip string) {
	b.tooltip = tooltip
	b.Refresh()
}

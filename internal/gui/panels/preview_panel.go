package panels

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/user/pom-manager/internal/gui/widgets"
)

// PreviewPane displays the generated POM XML with validation status
type PreviewPane struct {
	// UI components
	validationBadge *widgets.ValidationBadge
	xmlViewer       *widgets.XMLViewer
	copyButton      *widgets.ButtonWithTooltip
	refreshButton   *widgets.ButtonWithTooltip
	toolbar         *fyne.Container
	mainContainer   *fyne.Container

	// State
	livePreview bool
	currentXML  string // Store current XML for clipboard
}

// NewPreviewPane creates a new PreviewPane
func NewPreviewPane() *PreviewPane {
	pane := &PreviewPane{
		livePreview: true,
	}

	pane.createUI()
	return pane
}

// createUI creates the internal UI components
func (p *PreviewPane) createUI() {
	// Validation badge
	p.validationBadge = widgets.NewValidationBadge()

	// XML viewer with syntax highlighting
	p.xmlViewer = widgets.NewXMLViewer()

	// Copy to clipboard button with tooltip
	p.copyButton = widgets.NewButtonWithTooltip("Copy to Clipboard",
		"Copy the generated POM XML to clipboard",
		func() {
			p.CopyToClipboard()
		})

	// Refresh button with tooltip
	p.refreshButton = widgets.NewButtonWithTooltip("Refresh",
		"Manually refresh the XML preview (only enabled when live preview is off)",
		func() {
			// Refresh callback will be set by parent
		})

	// Toolbar with validation badge and buttons
	p.toolbar = container.NewBorder(
		nil, nil,
		p.validationBadge, // Left
		container.NewHBox(p.refreshButton, p.copyButton), // Right
	)

	// Main container with toolbar and XML display
	p.mainContainer = container.NewBorder(
		p.toolbar, // Top
		nil,       // Bottom
		nil,       // Left
		nil,       // Right
		p.xmlViewer, // Center (XMLViewer handles its own scrolling)
	)
}

// SetXML updates the displayed XML content
func (p *PreviewPane) SetXML(xml string) {
	p.currentXML = xml
	// UI updates must be called on UI thread
	fyne.Do(func() {
		p.xmlViewer.SetXML(xml)
	})
}

// SetValidationStatus updates the validation status indicator
func (p *PreviewPane) SetValidationStatus(valid bool, errorCount int) {
	// UI updates must be called on UI thread
	fyne.Do(func() {
		p.validationBadge.SetValid(valid)
		p.validationBadge.SetErrorCount(errorCount)
	})
}

// CopyToClipboard copies the XML content to the system clipboard
func (p *PreviewPane) CopyToClipboard() {
	if p.currentXML != "" {
		fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(p.currentXML)
	}
}

// SetLivePreview enables or disables live preview mode
func (p *PreviewPane) SetLivePreview(enabled bool) {
	p.livePreview = enabled
	// UI updates must be called on UI thread
	fyne.Do(func() {
		if enabled {
			p.refreshButton.Disable()
		} else {
			p.refreshButton.Enable()
		}
	})
}

// OnValidationBadgeClick sets the callback for when validation badge is clicked
func (p *PreviewPane) OnValidationBadgeClick(callback func()) {
	p.validationBadge.OnClick(callback)
}

// OnRefresh sets the callback for the refresh button
func (p *PreviewPane) OnRefresh(callback func()) {
	p.refreshButton.OnTapped = callback
}

// GetContainer returns the main container for embedding in parent layouts
func (p *PreviewPane) GetContainer() *fyne.Container {
	return p.mainContainer
}

package dialogs

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/gui/state"
	"github.com/user/pom-manager/internal/gui/widgets"
)

// SettingsDialog manages application settings through a tabbed interface
type SettingsDialog struct {
	window fyne.Window

	// Settings copy for editing
	tempSettings *state.Settings

	// General tab widgets
	themeSelect        *widget.Select
	autoSaveEntry      *widget.Entry
	restoreSessionCheck *widget.Check

	// Editor tab widgets
	fontSizeSlider     *widget.Slider
	fontSizeLabel      *widget.Label
	livePreviewCheck   *widget.Check
	validationDelayEntry *widget.Entry
	syntaxHighlightCheck *widget.Check

	// Templates tab widgets
	defaultTemplateSelect *widget.Select
	customTemplateDirEntry *widget.Entry

	// Advanced tab widgets
	mavenTimeoutEntry   *widget.Entry
	debugLogCheck       *widget.Check
	cacheDirEntry       *widget.Entry

	// Callbacks
	onSave func(*state.Settings)
}

// NewSettingsDialog creates a new settings dialog
func NewSettingsDialog(window fyne.Window, currentSettings *state.Settings) *SettingsDialog {
	// Create a copy of current settings for editing
	tempSettings := &state.Settings{}
	*tempSettings = *currentSettings

	return &SettingsDialog{
		window:       window,
		tempSettings: tempSettings,
	}
}

// Show displays the settings dialog
func (d *SettingsDialog) Show(onSave func(*state.Settings)) {
	d.onSave = onSave

	// Create tabs
	generalTab := d.createGeneralTab()
	editorTab := d.createEditorTab()
	templatesTab := d.createTemplatesTab()
	advancedTab := d.createAdvancedTab()

	// Create tabbed container
	tabs := container.NewAppTabs(
		container.NewTabItem("General", generalTab),
		container.NewTabItem("Editor", editorTab),
		container.NewTabItem("Templates", templatesTab),
		container.NewTabItem("Advanced", advancedTab),
	)

	// Create dialog variable for button callbacks
	var customDialog dialog.Dialog

	// Create buttons with tooltips
	resetButton := widgets.NewButtonWithTooltip("Reset to Defaults",
		"Reset all settings to their default values",
		func() {
			d.resetToDefaults()
		})

	cancelButton := widgets.NewButtonWithTooltip("Cancel",
		"Cancel and discard any changes",
		func() {
			if customDialog != nil {
				customDialog.Hide()
			}
		})

	okButton := widgets.NewButtonWithTooltip("OK",
		"Save settings and close the dialog",
		func() {
			if d.validateSettings() {
				if customDialog != nil {
					customDialog.Hide()
				}
				if d.onSave != nil {
					d.onSave(d.tempSettings)
				}
			}
		})
	okButton.Importance = widget.HighImportance

	// Button bar with proper spacing
	buttonBar := container.NewBorder(
		nil, nil,
		resetButton,
		container.NewHBox(cancelButton, okButton),
		widget.NewLabel(""), // Spacer
	)

	// Build complete content with padding
	content := container.NewBorder(
		nil,
		container.NewPadded(buttonBar),
		nil, nil,
		tabs,
	)

	// Create dialog without Close button (OK and Cancel are sufficient)
	customDialog = dialog.NewCustom(
		"Settings",
		"",
		content,
		d.window,
	)

	customDialog.Resize(fyne.NewSize(600, 500))
	customDialog.Show()
}

// createGeneralTab creates the General settings tab
func (d *SettingsDialog) createGeneralTab() fyne.CanvasObject {
	// Theme selection with immediate application
	d.themeSelect = widget.NewSelect(
		[]string{"light", "dark"},
		func(value string) {
			d.tempSettings.Theme = value
			// Apply theme immediately for preview
			d.applyThemePreview(value)
		},
	)
	d.themeSelect.SetSelected(d.tempSettings.Theme)

	// Auto-save interval
	d.autoSaveEntry = widget.NewEntry()
	d.autoSaveEntry.SetText(fmt.Sprintf("%d", d.tempSettings.AutoSaveInterval))
	d.autoSaveEntry.SetPlaceHolder("Minutes (0 = disabled)")

	// Restore session checkbox
	d.restoreSessionCheck = widget.NewCheck("Restore last opened file on startup", func(checked bool) {
		d.tempSettings.RestoreSession = checked
	})
	d.restoreSessionCheck.SetChecked(d.tempSettings.RestoreSession)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Theme", Widget: d.themeSelect},
			{Text: "Auto-save Interval (min)", Widget: d.autoSaveEntry},
			{Text: "Session Restore", Widget: d.restoreSessionCheck},
		},
	}

	return container.NewVBox(
		widget.NewLabel("General Settings"),
		widget.NewSeparator(),
		form,
	)
}

// createEditorTab creates the Editor settings tab
func (d *SettingsDialog) createEditorTab() fyne.CanvasObject {
	// Font size slider
	fontSizeBind := binding.NewFloat()
	fontSizeBind.Set(float64(d.tempSettings.FontSize))

	d.fontSizeSlider = widget.NewSliderWithData(10, 18, fontSizeBind)
	d.fontSizeLabel = widget.NewLabel(fmt.Sprintf("%d pt", d.tempSettings.FontSize))

	fontSizeBind.AddListener(binding.NewDataListener(func() {
		val, _ := fontSizeBind.Get()
		d.tempSettings.FontSize = int(val)
		d.fontSizeLabel.SetText(fmt.Sprintf("%d pt", int(val)))
	}))

	fontSizeContainer := container.NewHBox(
		d.fontSizeSlider,
		d.fontSizeLabel,
	)

	// Live preview checkbox
	d.livePreviewCheck = widget.NewCheck("Enable real-time XML preview", func(checked bool) {
		d.tempSettings.LivePreview = checked
	})
	d.livePreviewCheck.SetChecked(d.tempSettings.LivePreview)

	// Validation delay
	d.validationDelayEntry = widget.NewEntry()
	d.validationDelayEntry.SetText(fmt.Sprintf("%d", d.tempSettings.ValidationDelay))
	d.validationDelayEntry.SetPlaceHolder("Milliseconds")

	// Syntax highlighting checkbox
	d.syntaxHighlightCheck = widget.NewCheck("Enable XML syntax highlighting", func(checked bool) {
		d.tempSettings.SyntaxHighlight = checked
	})
	d.syntaxHighlightCheck.SetChecked(d.tempSettings.SyntaxHighlight)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Font Size", Widget: fontSizeContainer},
			{Text: "Live Preview", Widget: d.livePreviewCheck},
			{Text: "Validation Delay (ms)", Widget: d.validationDelayEntry},
			{Text: "Syntax Highlighting", Widget: d.syntaxHighlightCheck},
		},
	}

	return container.NewVBox(
		widget.NewLabel("Editor Settings"),
		widget.NewSeparator(),
		form,
	)
}

// createTemplatesTab creates the Templates settings tab
func (d *SettingsDialog) createTemplatesTab() fyne.CanvasObject {
	// Default template selection
	d.defaultTemplateSelect = widget.NewSelect(
		[]string{"basic-java", "java-library", "web-app"},
		func(value string) {
			d.tempSettings.DefaultTemplate = value
		},
	)
	d.defaultTemplateSelect.SetSelected(d.tempSettings.DefaultTemplate)

	// Custom template directory
	d.customTemplateDirEntry = widget.NewEntry()
	d.customTemplateDirEntry.SetText(d.tempSettings.CustomTemplateDir)
	d.customTemplateDirEntry.SetPlaceHolder("Path to custom templates directory")

	browseButton := widget.NewButton("Browse...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				d.customTemplateDirEntry.SetText(uri.Path())
				d.tempSettings.CustomTemplateDir = uri.Path()
			}
		}, d.window)
	})

	customDirContainer := container.NewBorder(
		nil, nil, nil, browseButton,
		d.customTemplateDirEntry,
	)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Default Template", Widget: d.defaultTemplateSelect},
			{Text: "Custom Templates Dir", Widget: customDirContainer},
		},
	}

	return container.NewVBox(
		widget.NewLabel("Template Settings"),
		widget.NewSeparator(),
		form,
		widget.NewLabel("Custom templates must follow the template structure."),
	)
}

// createAdvancedTab creates the Advanced settings tab
func (d *SettingsDialog) createAdvancedTab() fyne.CanvasObject {
	// Maven Central timeout
	d.mavenTimeoutEntry = widget.NewEntry()
	d.mavenTimeoutEntry.SetText(fmt.Sprintf("%d", d.tempSettings.MavenCentralTimeout))
	d.mavenTimeoutEntry.SetPlaceHolder("Seconds")

	// Debug log checkbox
	d.debugLogCheck = widget.NewCheck("Enable debug logging", func(checked bool) {
		d.tempSettings.EnableDebugLog = checked
	})
	d.debugLogCheck.SetChecked(d.tempSettings.EnableDebugLog)

	// Cache directory
	d.cacheDirEntry = widget.NewEntry()
	d.cacheDirEntry.SetText(d.tempSettings.CacheDir)
	d.cacheDirEntry.SetPlaceHolder("Default: ~/.pom-manager/cache")

	browseCacheButton := widget.NewButton("Browse...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				d.cacheDirEntry.SetText(uri.Path())
				d.tempSettings.CacheDir = uri.Path()
			}
		}, d.window)
	})

	cacheDirContainer := container.NewBorder(
		nil, nil, nil, browseCacheButton,
		d.cacheDirEntry,
	)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Maven Central Timeout (s)", Widget: d.mavenTimeoutEntry},
			{Text: "Debug Logging", Widget: d.debugLogCheck},
			{Text: "Cache Directory", Widget: cacheDirContainer},
		},
	}

	return container.NewVBox(
		widget.NewLabel("Advanced Settings"),
		widget.NewSeparator(),
		form,
		widget.NewLabel("⚠️ Advanced settings - change with caution"),
	)
}

// validateSettings validates all settings before saving
func (d *SettingsDialog) validateSettings() bool {
	// Validate auto-save interval
	autoSave, err := strconv.Atoi(d.autoSaveEntry.Text)
	if err != nil || autoSave < 0 {
		dialog.ShowError(fmt.Errorf("auto-save interval must be a non-negative number"), d.window)
		return false
	}
	d.tempSettings.AutoSaveInterval = autoSave

	// Validate validation delay
	validationDelay, err := strconv.Atoi(d.validationDelayEntry.Text)
	if err != nil || validationDelay < 0 || validationDelay > 5000 {
		dialog.ShowError(fmt.Errorf("validation delay must be between 0 and 5000 ms"), d.window)
		return false
	}
	d.tempSettings.ValidationDelay = validationDelay

	// Validate Maven timeout
	mavenTimeout, err := strconv.Atoi(d.mavenTimeoutEntry.Text)
	if err != nil || mavenTimeout < 1 || mavenTimeout > 300 {
		dialog.ShowError(fmt.Errorf("maven Central timeout must be between 1 and 300 seconds"), d.window)
		return false
	}
	d.tempSettings.MavenCentralTimeout = mavenTimeout

	return true
}

// resetToDefaults resets all settings to default values
func (d *SettingsDialog) resetToDefaults() {
	defaults := state.NewSettings()
	d.tempSettings = defaults

	// Update UI widgets
	d.themeSelect.SetSelected(defaults.Theme)
	d.autoSaveEntry.SetText(fmt.Sprintf("%d", defaults.AutoSaveInterval))
	d.restoreSessionCheck.SetChecked(defaults.RestoreSession)

	d.fontSizeSlider.SetValue(float64(defaults.FontSize))
	d.fontSizeLabel.SetText(fmt.Sprintf("%d pt", defaults.FontSize))
	d.livePreviewCheck.SetChecked(defaults.LivePreview)
	d.validationDelayEntry.SetText(fmt.Sprintf("%d", defaults.ValidationDelay))
	d.syntaxHighlightCheck.SetChecked(defaults.SyntaxHighlight)

	d.defaultTemplateSelect.SetSelected(defaults.DefaultTemplate)
	d.customTemplateDirEntry.SetText(defaults.CustomTemplateDir)

	d.mavenTimeoutEntry.SetText(fmt.Sprintf("%d", defaults.MavenCentralTimeout))
	d.debugLogCheck.SetChecked(defaults.EnableDebugLog)
	d.cacheDirEntry.SetText(defaults.CacheDir)

	// Apply default theme
	d.applyThemePreview(defaults.Theme)

	dialog.ShowInformation("Reset", "Settings have been reset to defaults", d.window)
}

// applyThemePreview applies the selected theme immediately for preview
func (d *SettingsDialog) applyThemePreview(themeName string) {
	app := fyne.CurrentApp()
	switch themeName {
	case "dark":
		app.Settings().SetTheme(theme.DarkTheme())
	case "light":
		app.Settings().SetTheme(theme.LightTheme())
	default:
		app.Settings().SetTheme(theme.DefaultTheme())
	}
}

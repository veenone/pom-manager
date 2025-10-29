package main

// GUI Entry Point for Maven POM Manager
//
// Build Requirements:
// - Fyne v2 requires CGO to be enabled
// - Windows: Install TDM-GCC or MinGW-w64 and set CGO_ENABLED=1
// - Build command: go build -o pom-manager-gui.exe ./cmd/gui
//
// For development without CGO, the CLI interface is available at cmd/cli/

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"github.com/user/pom-manager/internal/core/pom"
	"github.com/user/pom-manager/internal/gui/presenters"
	"github.com/user/pom-manager/internal/gui/state"
	"github.com/user/pom-manager/internal/gui/windows"
)

const (
	AppName    = "Maven POM Manager"
	AppID      = "com.pom-manager.gui"
	AppVersion = "0.1.0-MVP"
)

func main() {
	// Create Fyne application
	myApp := app.NewWithID(AppID)
	myApp.SetIcon(nil) // TODO: Add application icon later

	// Load settings from disk
	settings, err := state.LoadSettings()
	if err != nil {
		// Use defaults if loading fails
		settings = state.NewSettings()
	}

	// Apply theme based on settings
	applyTheme(myApp, settings.Theme)

	// Create main window
	window := myApp.NewWindow(AppName + " - " + AppVersion)

	// Restore window size from settings or use defaults
	windowSize := fyne.NewSize(
		float32(settings.WindowWidth),
		float32(settings.WindowHeight),
	)
	window.Resize(windowSize)

	// Initialize core engine components
	parser := pom.NewParser()
	generator := pom.NewGenerator()
	validator := pom.NewValidator()
	repository := pom.NewRepository()
	templateManager := pom.NewTemplateManager()

	// Initialize state with loaded settings
	appState := state.NewAppState()
	appState.SetSettings(settings)

	// Initialize presenter
	presenter := presenters.NewMainPresenter(
		parser,
		generator,
		validator,
		repository,
		templateManager,
		appState,
	)

	// Create main window
	mainWin := windows.NewMainWindow(window, presenter, appState)

	// Setup window close handler to save settings
	window.SetOnClosed(func() {
		// Update window size in settings
		currentSettings := appState.GetSettings()
		size := window.Content().Size()
		currentSettings.WindowWidth = int(size.Width)
		currentSettings.WindowHeight = int(size.Height)

		// Save current file path for session restore
		currentSettings.LastOpenedFile = appState.GetFilePath()

		// Save settings to disk
		_ = state.SaveSettings(currentSettings)
	})

	// Restore last file if RestoreSession is enabled
	if settings.RestoreSession && settings.LastOpenedFile != "" {
		_ = presenter.LoadPOM(settings.LastOpenedFile)
	}

	// Show main window
	mainWin.Show()
}

// applyTheme applies the specified theme to the application
func applyTheme(app fyne.App, themeName string) {
	switch themeName {
	case "dark":
		app.Settings().SetTheme(theme.DarkTheme())
	case "light":
		app.Settings().SetTheme(theme.LightTheme())
	default:
		// Use system default
		app.Settings().SetTheme(theme.DefaultTheme())
	}
}

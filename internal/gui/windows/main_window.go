package windows

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
	"github.com/user/pom-manager/internal/gui/dialogs"
	"github.com/user/pom-manager/internal/gui/dialogs/wizard"
	"github.com/user/pom-manager/internal/gui/panels"
	"github.com/user/pom-manager/internal/gui/presenters"
	"github.com/user/pom-manager/internal/gui/state"
)

// MainWindow is the main application window
type MainWindow struct {
	window    fyne.Window
	presenter presenters.MainPresenter
	appState  *state.AppState

	// Panels
	treePanel         *panels.TreePanel
	coordsPanel       *panels.CoordinatesPanel
	depsPanel         *panels.DependenciesPanel
	pluginsPanel      *panels.PluginsPanel
	propsPanel        *panels.PropertiesPanel
	profilesPanel     *panels.ProfilesPanel
	lifecyclePanel    *panels.LifecyclePanel
	previewPane       *panels.PreviewPane
	errorsPanel       *panels.ErrorsPanel

	// UI components
	tabContainer *container.AppTabs
	statusLabel  *widget.Label
	mainContent  *fyne.Container

	// Debouncing for preview updates
	refreshTimer    *time.Timer
	refreshPending  bool
	refreshDebounce time.Duration
}

// NewMainWindow creates a new main window
func NewMainWindow(
	window fyne.Window,
	presenter presenters.MainPresenter,
	appState *state.AppState,
) *MainWindow {
	mw := &MainWindow{
		window:    window,
		presenter: presenter,
		appState:  appState,
	}

	// Initialize debouncing from settings
	settings := appState.GetSettings()
	mw.refreshDebounce = time.Duration(settings.ValidationDelay) * time.Millisecond

	mw.createPanels()
	mw.createMenu()
	mw.createLayout()
	mw.setupCallbacks()

	return mw
}

// createPanels initializes all panels
func (mw *MainWindow) createPanels() {
	mw.treePanel = panels.NewTreePanel()
	mw.coordsPanel = panels.NewCoordinatesPanel()
	mw.depsPanel = panels.NewDependenciesPanel()
	mw.pluginsPanel = panels.NewPluginsPanel()
	mw.propsPanel = panels.NewPropertiesPanel(mw.window)
	mw.profilesPanel = panels.NewProfilesPanel()
	mw.lifecyclePanel = panels.NewLifecyclePanel()
	mw.previewPane = panels.NewPreviewPane()
	mw.errorsPanel = panels.NewErrorsPanel()
}

// createMenu creates the menu bar
func (mw *MainWindow) createMenu() {
	// File menu
	newItem := fyne.NewMenuItem("New", mw.handleNew)
	openItem := fyne.NewMenuItem("Open", mw.handleOpen)

	// Open Recent submenu
	recentMenu := fyne.NewMenu("Open Recent")
	mw.updateRecentFilesMenu(recentMenu)
	recentItem := fyne.NewMenuItem("Open Recent", nil)
	recentItem.ChildMenu = recentMenu

	saveItem := fyne.NewMenuItem("Save", mw.handleSave)
	saveAsItem := fyne.NewMenuItem("Save As...", mw.handleSaveAs)
	exitItem := fyne.NewMenuItem("Exit", func() {
		mw.window.Close()
	})

	fileMenu := fyne.NewMenu("File", newItem, openItem, recentItem, fyne.NewMenuItemSeparator(), saveItem, saveAsItem, fyne.NewMenuItemSeparator(), exitItem)

	// Edit menu
	settingsItem := fyne.NewMenuItem("Settings...", mw.handleSettings)
	editMenu := fyne.NewMenu("Edit", settingsItem)

	// Help menu
	quickHelpItem := fyne.NewMenuItem("Quick Help", mw.handleQuickHelp)
	mavenBasicsItem := fyne.NewMenuItem("Maven Basics", mw.handleMavenBasics)
	aboutItem := fyne.NewMenuItem("About", mw.handleAbout)
	helpMenu := fyne.NewMenu("Help", quickHelpItem, mavenBasicsItem, fyne.NewMenuItemSeparator(), aboutItem)

	mainMenu := fyne.NewMainMenu(fileMenu, editMenu, helpMenu)
	mw.window.SetMainMenu(mainMenu)
}

// createLayout creates the main window layout
func (mw *MainWindow) createLayout() {
	// Create tabs for editor panels
	mw.tabContainer = container.NewAppTabs(
		container.NewTabItem("Coordinates", mw.coordsPanel.GetContainer()),
		container.NewTabItem("Dependencies", mw.depsPanel.GetContainer()),
		container.NewTabItem("Plugins", mw.pluginsPanel.GetContainer()),
		container.NewTabItem("Properties", mw.propsPanel.GetContainer()),
		container.NewTabItem("Profiles", mw.profilesPanel.GetContainer()),
		container.NewTabItem("Lifecycle Phases", mw.lifecyclePanel.GetContainer()),
	)

	// Create center panel with tabs and errors
	centerPanel := container.NewBorder(
		nil,
		mw.errorsPanel.GetContainer(),
		nil, nil,
		mw.tabContainer,
	)

	// Create three-panel layout
	splitLeft := container.NewHSplit(
		mw.treePanel.GetContainer(),
		centerPanel,
	)
	splitLeft.SetOffset(0.2) // 20% for tree

	splitMain := container.NewHSplit(
		splitLeft,
		mw.previewPane.GetContainer(),
	)
	splitMain.SetOffset(0.65) // 65% for left (tree + editor), 35% for preview

	// Status bar
	mw.statusLabel = widget.NewLabel("Ready")
	statusBar := container.NewHBox(mw.statusLabel)

	// Main content
	mw.mainContent = container.NewBorder(
		nil,        // Top (menu is separate)
		statusBar,  // Bottom
		nil, nil,   // Left, Right
		splitMain,  // Center
	)

	mw.window.SetContent(mw.mainContent)
}

// setupCallbacks sets up event handlers
func (mw *MainWindow) setupCallbacks() {
	// Subscribe to state changes with debouncing
	mw.presenter.SubscribeToChanges(func() {
		mw.debouncedRefreshUI()
	})

	// Coordinates panel
	mw.coordsPanel.OnChange(func(coords pom.Coordinates) {
		mw.presenter.UpdateCoordinates(coords)
	})

	// Dependencies panel
	mw.depsPanel.OnAdd(func() {
		depDialog := dialogs.NewDependencyDialog(mw.window)
		depDialog.ShowAdd(func(dep pom.Dependency) {
			mw.presenter.AddDependency(dep)
		})
	})

	mw.depsPanel.OnEdit(func(dep pom.Dependency) {
		depDialog := dialogs.NewDependencyDialog(mw.window)
		depDialog.ShowEdit(dep, func(updated pom.Dependency) {
			mw.presenter.AddDependency(updated) // Add/update logic
		})
	})

	mw.depsPanel.OnRemove(func(dep pom.Dependency) {
		mw.presenter.RemoveDependency(dep.GroupID, dep.ArtifactID)
	})

	// Plugins panel
	mw.pluginsPanel.OnAdd(func() {
		pluginDialog := dialogs.NewPluginDialog(mw.window)
		pluginDialog.ShowAdd(func(plugin pom.Plugin) {
			mw.presenter.AddPlugin(plugin)
		})
	})

	mw.pluginsPanel.OnEdit(func(plugin pom.Plugin) {
		pluginDialog := dialogs.NewPluginDialog(mw.window)
		pluginDialog.ShowEdit(plugin, func(updated pom.Plugin) {
			mw.presenter.AddPlugin(updated)
		})
	})

	mw.pluginsPanel.OnRemove(func(plugin pom.Plugin) {
		mw.presenter.RemovePlugin(plugin.GroupID, plugin.ArtifactID)
	})

	// Properties panel
	mw.propsPanel.OnChange(func(props map[string]string) {
		mw.presenter.UpdateProperties(props)
	})

	// Lifecycle panel
	mw.lifecyclePanel.OnAddExecution(func(pluginIndex int, execution pom.PluginExecution) {
		mw.handleAddExecution(pluginIndex, execution)
	})

	mw.lifecyclePanel.OnRemoveExecution(func(pluginIndex int, executionID string) {
		mw.handleRemoveExecution(pluginIndex, executionID)
	})

	// Tree panel - navigate to corresponding tab when node selected
	mw.treePanel.OnNodeSelected(func(nodeType string, id string) {
		fyne.Do(func() {
			switch nodeType {
			case "coordinates":
				mw.tabContainer.SelectIndex(0) // Coordinates tab
			case "dependencies", "dep":
				mw.tabContainer.SelectIndex(1) // Dependencies tab
			case "plugins", "plugin":
				mw.tabContainer.SelectIndex(2) // Plugins tab
			case "properties", "prop":
				mw.tabContainer.SelectIndex(3) // Properties tab
			case "profiles", "profile":
				mw.tabContainer.SelectIndex(4) // Profiles tab
			}
		})
	})

	// Setup keyboard shortcuts
	mw.setupKeyboardShortcuts()
}

// setupKeyboardShortcuts configures keyboard shortcuts for the main window
func (mw *MainWindow) setupKeyboardShortcuts() {
	// Ctrl+N: New POM
	mw.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		mw.handleNew()
	})

	// Ctrl+O: Open
	mw.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		mw.handleOpen()
	})

	// Ctrl+S: Save
	mw.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		mw.handleSave()
	})

	// Ctrl+Shift+S: Save As
	mw.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		mw.handleSaveAs()
	})

	// Ctrl+W: Close
	mw.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyW,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		mw.window.Close()
	})

	// Ctrl+Q: Quit
	mw.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		mw.window.Close()
	})

	// F1: Help
	mw.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyF1,
	}, func(shortcut fyne.Shortcut) {
		mw.handleHelp()
	})

	// F5: Refresh/Validate
	mw.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyF5,
	}, func(shortcut fyne.Shortcut) {
		mw.handleRefresh()
	})
}

// debouncedRefreshUI debounces refreshUI calls to prevent excessive updates
func (mw *MainWindow) debouncedRefreshUI() {
	// If debouncing is disabled (0ms), refresh immediately
	if mw.refreshDebounce == 0 {
		mw.refreshUI()
		return
	}

	// Cancel existing timer if pending
	if mw.refreshTimer != nil {
		mw.refreshTimer.Stop()
	}

	// Schedule a new refresh
	mw.refreshTimer = time.AfterFunc(mw.refreshDebounce, func() {
		mw.refreshUI()
		mw.refreshPending = false
	})
	mw.refreshPending = true
}

// refreshUI updates all UI components from current state
func (mw *MainWindow) refreshUI() {
	project := mw.presenter.GetCurrentProject()
	if project == nil {
		return
	}

	// Update panels
	mw.coordsPanel.LoadProject(project)
	mw.depsPanel.LoadDependencies(project.Dependencies)

	if project.Build != nil {
		mw.pluginsPanel.LoadPlugins(project.Build.Plugins)
	}

	mw.propsPanel.LoadProperties(project.Properties)
	mw.profilesPanel.LoadProfiles(project.Profiles)
	mw.lifecyclePanel.LoadProject(project)
	mw.treePanel.LoadProject(project)

	// Validate and update preview
	result, _ := mw.presenter.ValidateCurrent()

	// Update errors panel
	mw.errorsPanel.SetErrors(result)

	// Update preview pane
	generator := pom.NewGenerator()
	xmlData, err := generator.Generate(project)
	if err == nil {
		mw.previewPane.SetXML(string(xmlData))
	}

	errorCount := len(result.Errors.AllErrors())
	mw.previewPane.SetValidationStatus(result.Valid, errorCount)

	// Update status bar (must be on UI thread)
	filePath := mw.appState.GetFilePath()
	statusText := ""
	if filePath != "" {
		statusText = fmt.Sprintf("File: %s | %s", filePath, mw.getValidationStatus(result))
	} else {
		statusText = fmt.Sprintf("Unsaved | %s", mw.getValidationStatus(result))
	}

	fyne.Do(func() {
		mw.statusLabel.SetText(statusText)
	})
}

// getValidationStatus returns validation status string
func (mw *MainWindow) getValidationStatus(result pom.ValidationResult) string {
	if result.Valid {
		return "✓ Valid"
	}
	errorCount := len(result.Errors.AllErrors())
	return fmt.Sprintf("✗ Invalid (%d errors)", errorCount)
}

// Menu handlers
func (mw *MainWindow) handleNew() {
	wiz := wizard.NewCreateWizard(mw.window)
	wiz.Show(func(coords pom.Coordinates, template string) {
		err := mw.presenter.CreateNewPOM(coords, template)
		if err != nil {
			dialog.ShowError(err, mw.window)
		}
	})
}

func (mw *MainWindow) handleOpen() {
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		path := reader.URI().Path()
		err = mw.presenter.LoadPOM(path)
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}

		// Add to recent files
		settings := mw.appState.GetSettings()
		settings.AddRecentFile(path)
		mw.appState.SetSettings(settings)
		state.SaveSettings(settings) // Save to persist
	}, mw.window)

	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".xml"}))
	fileDialog.Show()
}

// updateRecentFilesMenu updates the Open Recent submenu
func (mw *MainWindow) updateRecentFilesMenu(menu *fyne.Menu) {
	menu.Items = nil // Clear existing items

	settings := mw.appState.GetSettings()
	recentFiles := settings.GetRecentFiles()

	if len(recentFiles) == 0 {
		menu.Items = append(menu.Items, fyne.NewMenuItem("(No recent files)", nil))
		menu.Items[0].Disabled = true
		return
	}

	for _, filePath := range recentFiles {
		// Create copy for closure
		path := filePath
		// Get file name only for display
		fileName := filePath
		if idx := len(filePath) - 1; idx >= 0 {
			for i := idx; i >= 0; i-- {
				if filePath[i] == '/' || filePath[i] == '\\' {
					fileName = filePath[i+1:]
					break
				}
			}
		}

		item := fyne.NewMenuItem(fileName, func() {
			err := mw.presenter.LoadPOM(path)
			if err != nil {
				dialog.ShowError(err, mw.window)
			}
		})
		menu.Items = append(menu.Items, item)
	}

	// Add "Clear Recent" option
	menu.Items = append(menu.Items, fyne.NewMenuItemSeparator())
	clearItem := fyne.NewMenuItem("Clear Recent Files", func() {
		settings := mw.appState.GetSettings()
		settings.RecentFiles = []string{}
		mw.appState.SetSettings(settings)
		state.SaveSettings(settings)
		// Refresh menu
		mw.createMenu()
	})
	menu.Items = append(menu.Items, clearItem)
}

func (mw *MainWindow) handleSave() {
	filePath := mw.appState.GetFilePath()
	if filePath == "" {
		mw.handleSaveAs()
		return
	}

	err := mw.presenter.SavePOM(filePath)
	if err != nil {
		dialog.ShowError(err, mw.window)
	} else {
		dialog.ShowInformation("Saved", "POM file saved successfully", mw.window)
	}
}

func (mw *MainWindow) handleSaveAs() {
	fileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		path := writer.URI().Path()
		err = mw.presenter.SavePOM(path)
		if err != nil {
			dialog.ShowError(err, mw.window)
		} else {
			dialog.ShowInformation("Saved", "POM file saved successfully", mw.window)
		}
	}, mw.window)

	fileDialog.SetFileName("pom.xml")
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".xml"}))
	fileDialog.Show()
}

func (mw *MainWindow) handleSettings() {
	currentSettings := mw.appState.GetSettings()
	settingsDialog := dialogs.NewSettingsDialog(mw.window, currentSettings)
	settingsDialog.Show(func(updatedSettings *state.Settings) {
		// Update app state
		mw.appState.SetSettings(updatedSettings)

		// Save to disk
		if err := state.SaveSettings(updatedSettings); err != nil {
			dialog.ShowError(err, mw.window)
			return
		}

		dialog.ShowInformation("Settings Saved", "Settings have been saved successfully", mw.window)
	})
}

func (mw *MainWindow) handleHelp() {
	helpDialog := dialogs.NewHelpDialog(mw.window)
	helpDialog.ShowQuickHelp()
}

func (mw *MainWindow) handleQuickHelp() {
	helpDialog := dialogs.NewHelpDialog(mw.window)
	helpDialog.ShowQuickHelp()
}

func (mw *MainWindow) handleMavenBasics() {
	helpDialog := dialogs.NewHelpDialog(mw.window)
	helpDialog.Show()
}

func (mw *MainWindow) handleRefresh() {
	// Force re-validation and UI refresh
	result, err := mw.presenter.ValidateCurrent()
	if err != nil {
		dialog.ShowError(err, mw.window)
		return
	}

	// Update errors panel
	mw.errorsPanel.SetErrors(result)

	// Update preview pane
	project := mw.presenter.GetCurrentProject()
	if project != nil {
		generator := pom.NewGenerator()
		xmlData, err := generator.Generate(project)
		if err == nil {
			mw.previewPane.SetXML(string(xmlData))
		}
		errorCount := len(result.Errors.AllErrors())
		mw.previewPane.SetValidationStatus(result.Valid, errorCount)
	}

	// Update status bar
	mw.statusLabel.SetText(fmt.Sprintf("Refreshed | %s", mw.getValidationStatus(result)))
}

func (mw *MainWindow) handleAbout() {
	about := dialog.NewInformation("About",
		"Maven POM Manager v0.1.0-MVP\n\nA desktop application for creating and managing Maven POM files.",
		mw.window)
	about.Show()
}

// handleAddExecution handles adding a new plugin execution
func (mw *MainWindow) handleAddExecution(pluginIndex int, execution pom.PluginExecution) {
	project := mw.presenter.GetCurrentProject()
	if project == nil || project.Build == nil {
		return
	}

	// Show execution dialog
	execDialog := dialogs.NewExecutionDialog(mw.window, project.Build.Plugins)
	execDialog.ShowAdd(func(selectedPluginIndex int, newExecution pom.PluginExecution) {
		if selectedPluginIndex >= 0 && selectedPluginIndex < len(project.Build.Plugins) {
			// Add execution to the selected plugin
			project.Build.Plugins[selectedPluginIndex].Executions = append(
				project.Build.Plugins[selectedPluginIndex].Executions,
				newExecution,
			)
			// Notify presenter of change
			mw.presenter.UpdateProject(project)
		}
	})
}

// handleRemoveExecution handles removing a plugin execution
func (mw *MainWindow) handleRemoveExecution(pluginIndex int, executionID string) {
	project := mw.presenter.GetCurrentProject()
	if project == nil || project.Build == nil || pluginIndex < 0 || pluginIndex >= len(project.Build.Plugins) {
		return
	}

	plugin := &project.Build.Plugins[pluginIndex]

	// Find and remove the execution
	for i, exec := range plugin.Executions {
		if exec.ID == executionID {
			plugin.Executions = append(plugin.Executions[:i], plugin.Executions[i+1:]...)
			mw.presenter.UpdateProject(project)
			break
		}
	}
}

// Show displays the window
func (mw *MainWindow) Show() {
	mw.window.ShowAndRun()
}

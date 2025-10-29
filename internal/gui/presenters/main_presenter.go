package presenters

import (
	"fmt"

	"github.com/user/pom-manager/internal/core/pom"
	"github.com/user/pom-manager/internal/gui/state"
)

// MainPresenter orchestrates the main window logic and coordinates
// between UI components and the core POM engine
type MainPresenter interface {
	// File operations
	LoadPOM(path string) error
	SavePOM(path string) error
	CreateNewPOM(coords pom.Coordinates, template string) error

	// POM operations
	ValidateCurrent() (pom.ValidationResult, error)
	UpdateCoordinates(coords pom.Coordinates) error
	AddDependency(dep pom.Dependency) error
	RemoveDependency(groupID, artifactID string) error
	AddPlugin(plugin pom.Plugin) error
	RemovePlugin(groupID, artifactID string) error
	UpdateProperties(props map[string]string) error
	UpdateProject(project *pom.Project) error

	// State access
	GetCurrentProject() *pom.Project
	SubscribeToChanges(callback func())
}

// mainPresenter is the concrete implementation of MainPresenter
type mainPresenter struct {
	parser          pom.Parser
	generator       pom.Generator
	validator       pom.Validator
	repository      pom.Repository
	templateManager pom.TemplateManager
	appState        *state.AppState
}

// NewMainPresenter creates a new MainPresenter with injected dependencies
func NewMainPresenter(
	parser pom.Parser,
	generator pom.Generator,
	validator pom.Validator,
	repository pom.Repository,
	templateManager pom.TemplateManager,
	appState *state.AppState,
) MainPresenter {
	return &mainPresenter{
		parser:          parser,
		generator:       generator,
		validator:       validator,
		repository:      repository,
		templateManager: templateManager,
		appState:        appState,
	}
}

// LoadPOM loads a POM file from the specified path
func (p *mainPresenter) LoadPOM(path string) error {
	// Parse the file
	project, err := p.parser.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to load POM: %w", err)
	}

	// Update app state
	p.appState.SetCurrentProject(project)
	p.appState.SetFilePath(path)
	p.appState.SetDirty(false)

	return nil
}

// SavePOM saves the current POM to the specified path
func (p *mainPresenter) SavePOM(path string) error {
	project := p.appState.GetCurrentProject()
	if project == nil {
		return fmt.Errorf("no project loaded")
	}

	// Generate XML
	xmlData, err := p.generator.Generate(project)
	if err != nil {
		return fmt.Errorf("failed to generate POM XML: %w", err)
	}

	// Write to file
	if err := p.repository.Write(path, xmlData); err != nil {
		return fmt.Errorf("failed to save POM: %w", err)
	}

	// Update app state
	p.appState.SetFilePath(path)
	p.appState.SetDirty(false)

	return nil
}

// CreateNewPOM creates a new POM from a template with the given coordinates
func (p *mainPresenter) CreateNewPOM(coords pom.Coordinates, template string) error {
	// Create project from template
	project, err := p.templateManager.Create(template, coords)
	if err != nil {
		return fmt.Errorf("failed to create POM from template: %w", err)
	}

	// Update app state
	p.appState.SetCurrentProject(project)
	p.appState.SetFilePath("") // New file, not saved yet
	p.appState.SetDirty(true)

	return nil
}

// ValidateCurrent validates the current project
func (p *mainPresenter) ValidateCurrent() (pom.ValidationResult, error) {
	project := p.appState.GetCurrentProject()
	if project == nil {
		return pom.ValidationResult{}, fmt.Errorf("no project loaded")
	}

	result := p.validator.Validate(project)
	return result, nil
}

// UpdateCoordinates updates the project coordinates
func (p *mainPresenter) UpdateCoordinates(coords pom.Coordinates) error {
	project := p.appState.GetCurrentProject()
	if project == nil {
		return fmt.Errorf("no project loaded")
	}

	// Update coordinates
	project.GroupID = coords.GroupID
	project.ArtifactID = coords.ArtifactID
	project.Version = coords.Version
	project.Coordinates = coords

	// Mark as dirty and notify
	p.appState.SetDirty(true)
	p.appState.SetCurrentProject(project) // This triggers notification

	return nil
}

// AddDependency adds a new dependency to the project
func (p *mainPresenter) AddDependency(dep pom.Dependency) error {
	project := p.appState.GetCurrentProject()
	if project == nil {
		return fmt.Errorf("no project loaded")
	}

	// Check for duplicates
	for i, existing := range project.Dependencies {
		if existing.GroupID == dep.GroupID && existing.ArtifactID == dep.ArtifactID {
			// Update existing dependency
			project.Dependencies[i] = dep
			p.appState.SetDirty(true)
			p.appState.SetCurrentProject(project)
			return nil
		}
	}

	// Add new dependency
	project.Dependencies = append(project.Dependencies, dep)
	p.appState.SetDirty(true)
	p.appState.SetCurrentProject(project)

	return nil
}

// RemoveDependency removes a dependency from the project
func (p *mainPresenter) RemoveDependency(groupID, artifactID string) error {
	project := p.appState.GetCurrentProject()
	if project == nil {
		return fmt.Errorf("no project loaded")
	}

	// Find and remove dependency
	for i, dep := range project.Dependencies {
		if dep.GroupID == groupID && dep.ArtifactID == artifactID {
			project.Dependencies = append(project.Dependencies[:i], project.Dependencies[i+1:]...)
			p.appState.SetDirty(true)
			p.appState.SetCurrentProject(project)
			return nil
		}
	}

	return fmt.Errorf("dependency not found: %s:%s", groupID, artifactID)
}

// AddPlugin adds a new plugin to the project's build configuration
func (p *mainPresenter) AddPlugin(plugin pom.Plugin) error {
	project := p.appState.GetCurrentProject()
	if project == nil {
		return fmt.Errorf("no project loaded")
	}

	// Ensure Build section exists
	if project.Build == nil {
		project.Build = &pom.Build{
			Plugins: make([]pom.Plugin, 0),
		}
	}

	// Check for duplicates
	for i, existing := range project.Build.Plugins {
		if existing.GroupID == plugin.GroupID && existing.ArtifactID == plugin.ArtifactID {
			// Update existing plugin
			project.Build.Plugins[i] = plugin
			p.appState.SetDirty(true)
			p.appState.SetCurrentProject(project)
			return nil
		}
	}

	// Add new plugin
	project.Build.Plugins = append(project.Build.Plugins, plugin)
	p.appState.SetDirty(true)
	p.appState.SetCurrentProject(project)

	return nil
}

// RemovePlugin removes a plugin from the project's build configuration
func (p *mainPresenter) RemovePlugin(groupID, artifactID string) error {
	project := p.appState.GetCurrentProject()
	if project == nil {
		return fmt.Errorf("no project loaded")
	}

	if project.Build == nil {
		return fmt.Errorf("no build configuration")
	}

	// Find and remove plugin
	for i, plugin := range project.Build.Plugins {
		if plugin.GroupID == groupID && plugin.ArtifactID == artifactID {
			project.Build.Plugins = append(project.Build.Plugins[:i], project.Build.Plugins[i+1:]...)
			p.appState.SetDirty(true)
			p.appState.SetCurrentProject(project)
			return nil
		}
	}

	return fmt.Errorf("plugin not found: %s:%s", groupID, artifactID)
}

// UpdateProperties updates the project properties
func (p *mainPresenter) UpdateProperties(props map[string]string) error {
	project := p.appState.GetCurrentProject()
	if project == nil {
		return fmt.Errorf("no project loaded")
	}

	// Update properties
	project.Properties = props
	p.appState.SetDirty(true)
	p.appState.SetCurrentProject(project)

	return nil
}

// UpdateProject updates the entire project
func (p *mainPresenter) UpdateProject(project *pom.Project) error {
	if project == nil {
		return fmt.Errorf("project cannot be nil")
	}

	p.appState.SetDirty(true)
	p.appState.SetCurrentProject(project)

	return nil
}

// GetCurrentProject returns the current project from app state
func (p *mainPresenter) GetCurrentProject() *pom.Project {
	return p.appState.GetCurrentProject()
}

// SubscribeToChanges registers a callback for state changes
func (p *mainPresenter) SubscribeToChanges(callback func()) {
	p.appState.Subscribe(callback)
}

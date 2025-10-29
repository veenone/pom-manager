package presenters

import (
	"testing"

	"github.com/user/pom-manager/internal/core/pom"
	"github.com/user/pom-manager/internal/gui/state"
)

func TestNewMainPresenter(t *testing.T) {
	parser := pom.NewParser()
	generator := pom.NewGenerator()
	validator := pom.NewValidator()
	repository := pom.NewRepository()
	templateManager := pom.NewTemplateManager()
	appState := state.NewAppState()

	presenter := NewMainPresenter(
		parser,
		generator,
		validator,
		repository,
		templateManager,
		appState,
	)

	if presenter == nil {
		t.Fatal("NewMainPresenter returned nil")
	}
}

func TestCreateNewPOM(t *testing.T) {
	parser := pom.NewParser()
	generator := pom.NewGenerator()
	validator := pom.NewValidator()
	repository := pom.NewRepository()
	templateManager := pom.NewTemplateManager()
	appState := state.NewAppState()

	presenter := NewMainPresenter(
		parser,
		generator,
		validator,
		repository,
		templateManager,
		appState,
	)

	// Create coordinates
	coords := pom.Coordinates{
		GroupID:    "com.example",
		ArtifactID: "test-app",
		Version:    "1.0.0",
	}

	// Create new POM with basic-java template
	err := presenter.CreateNewPOM(coords, "basic-java")
	if err != nil {
		t.Fatalf("Failed to create new POM: %v", err)
	}

	// Verify project was created
	project := presenter.GetCurrentProject()
	if project == nil {
		t.Fatal("Expected project to be created")
	}

	if project.GroupID != "com.example" {
		t.Errorf("Expected GroupID 'com.example', got '%s'", project.GroupID)
	}

	if project.ArtifactID != "test-app" {
		t.Errorf("Expected ArtifactID 'test-app', got '%s'", project.ArtifactID)
	}

	if project.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", project.Version)
	}
}

func TestUpdateCoordinates(t *testing.T) {
	parser := pom.NewParser()
	generator := pom.NewGenerator()
	validator := pom.NewValidator()
	repository := pom.NewRepository()
	templateManager := pom.NewTemplateManager()
	appState := state.NewAppState()

	presenter := NewMainPresenter(
		parser,
		generator,
		validator,
		repository,
		templateManager,
		appState,
	)

	// Create initial project
	initialCoords := pom.Coordinates{
		GroupID:    "com.initial",
		ArtifactID: "initial-app",
		Version:    "0.1.0",
	}
	_ = presenter.CreateNewPOM(initialCoords, "basic-java")

	// Update coordinates
	updatedCoords := pom.Coordinates{
		GroupID:    "com.updated",
		ArtifactID: "updated-app",
		Version:    "1.0.0",
	}
	presenter.UpdateCoordinates(updatedCoords)

	// Verify update
	project := presenter.GetCurrentProject()
	if project.GroupID != "com.updated" {
		t.Errorf("Expected updated GroupID 'com.updated', got '%s'", project.GroupID)
	}

	if project.ArtifactID != "updated-app" {
		t.Errorf("Expected updated ArtifactID 'updated-app', got '%s'", project.ArtifactID)
	}

	if project.Version != "1.0.0" {
		t.Errorf("Expected updated Version '1.0.0', got '%s'", project.Version)
	}
}

func TestAddDependency(t *testing.T) {
	parser := pom.NewParser()
	generator := pom.NewGenerator()
	validator := pom.NewValidator()
	repository := pom.NewRepository()
	templateManager := pom.NewTemplateManager()
	appState := state.NewAppState()

	presenter := NewMainPresenter(
		parser,
		generator,
		validator,
		repository,
		templateManager,
		appState,
	)

	// Create project
	coords := pom.Coordinates{
		GroupID:    "com.example",
		ArtifactID: "test-app",
		Version:    "1.0.0",
	}
	_ = presenter.CreateNewPOM(coords, "basic-java")

	// Add dependency
	dep := pom.Dependency{
		GroupID:    "junit",
		ArtifactID: "junit",
		Version:    "4.13.2",
		Scope:      "test",
	}
	presenter.AddDependency(dep)

	// Verify dependency was added
	project := presenter.GetCurrentProject()
	if len(project.Dependencies) == 0 {
		t.Fatal("Expected at least one dependency")
	}

	found := false
	for _, d := range project.Dependencies {
		if d.GroupID == "junit" && d.ArtifactID == "junit" {
			found = true
			if d.Version != "4.13.2" {
				t.Errorf("Expected version '4.13.2', got '%s'", d.Version)
			}
			if d.Scope != "test" {
				t.Errorf("Expected scope 'test', got '%s'", d.Scope)
			}
			break
		}
	}

	if !found {
		t.Error("Dependency not found in project")
	}
}

func TestRemoveDependency(t *testing.T) {
	parser := pom.NewParser()
	generator := pom.NewGenerator()
	validator := pom.NewValidator()
	repository := pom.NewRepository()
	templateManager := pom.NewTemplateManager()
	appState := state.NewAppState()

	presenter := NewMainPresenter(
		parser,
		generator,
		validator,
		repository,
		templateManager,
		appState,
	)

	// Create project and add dependency
	coords := pom.Coordinates{
		GroupID:    "com.example",
		ArtifactID: "test-app",
		Version:    "1.0.0",
	}
	_ = presenter.CreateNewPOM(coords, "basic-java")

	dep := pom.Dependency{
		GroupID:    "junit",
		ArtifactID: "junit",
		Version:    "4.13.2",
		Scope:      "test",
	}
	presenter.AddDependency(dep)

	// Verify it was added
	project := presenter.GetCurrentProject()
	initialCount := len(project.Dependencies)

	// Remove dependency
	presenter.RemoveDependency("junit", "junit")

	// Verify it was removed
	project = presenter.GetCurrentProject()
	if len(project.Dependencies) >= initialCount {
		t.Error("Expected dependency count to decrease")
	}

	for _, d := range project.Dependencies {
		if d.GroupID == "junit" && d.ArtifactID == "junit" {
			t.Error("Dependency should have been removed")
		}
	}
}

func TestValidateCurrent(t *testing.T) {
	parser := pom.NewParser()
	generator := pom.NewGenerator()
	validator := pom.NewValidator()
	repository := pom.NewRepository()
	templateManager := pom.NewTemplateManager()
	appState := state.NewAppState()

	presenter := NewMainPresenter(
		parser,
		generator,
		validator,
		repository,
		templateManager,
		appState,
	)

	// Create valid project
	coords := pom.Coordinates{
		GroupID:    "com.example",
		ArtifactID: "test-app",
		Version:    "1.0.0",
	}
	_ = presenter.CreateNewPOM(coords, "basic-java")

	// Validate
	result, err := presenter.ValidateCurrent()
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.Valid {
		t.Error("Expected project to be valid")
	}
}

func TestSubscribeToChanges(t *testing.T) {
	parser := pom.NewParser()
	generator := pom.NewGenerator()
	validator := pom.NewValidator()
	repository := pom.NewRepository()
	templateManager := pom.NewTemplateManager()
	appState := state.NewAppState()

	presenter := NewMainPresenter(
		parser,
		generator,
		validator,
		repository,
		templateManager,
		appState,
	)

	// Track if callback was called
	callbackCalled := false

	// Subscribe to changes
	presenter.SubscribeToChanges(func() {
		callbackCalled = true
	})

	// Make a change that triggers notification
	coords := pom.Coordinates{
		GroupID:    "com.example",
		ArtifactID: "test-app",
		Version:    "1.0.0",
	}
	_ = presenter.CreateNewPOM(coords, "basic-java")

	// The callback should be called automatically via CreateNewPOM
	// which calls appState.SetCurrentProject internally

	if !callbackCalled {
		t.Error("Expected callback to be called after changes")
	}
}

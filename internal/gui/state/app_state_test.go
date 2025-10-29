package state

import (
	"testing"

	"github.com/user/pom-manager/internal/core/pom"
)

func TestNewAppState(t *testing.T) {
	state := NewAppState()

	if state == nil {
		t.Fatal("NewAppState returned nil")
	}

	if state.GetSettings() == nil {
		t.Error("Expected default settings to be initialized")
	}

	if state.GetCurrentProject() != nil {
		t.Error("Expected no project initially")
	}

	if state.GetFilePath() != "" {
		t.Error("Expected empty file path initially")
	}
}

func TestSetAndGetCurrentProject(t *testing.T) {
	state := NewAppState()

	// Create test project
	project := &pom.Project{
		GroupID:    "com.test",
		ArtifactID: "test-app",
		Version:    "1.0.0",
		Packaging:  "jar",
	}

	// Set project
	state.SetCurrentProject(project)

	// Get project
	retrieved := state.GetCurrentProject()
	if retrieved == nil {
		t.Fatal("Expected to retrieve project")
	}

	if retrieved.GroupID != "com.test" {
		t.Errorf("Expected GroupID 'com.test', got '%s'", retrieved.GroupID)
	}

	if retrieved.ArtifactID != "test-app" {
		t.Errorf("Expected ArtifactID 'test-app', got '%s'", retrieved.ArtifactID)
	}

	if retrieved.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", retrieved.Version)
	}
}

func TestSetAndGetFilePath(t *testing.T) {
	state := NewAppState()

	testPath := "/path/to/pom.xml"
	state.SetFilePath(testPath)

	retrieved := state.GetFilePath()
	if retrieved != testPath {
		t.Errorf("Expected file path '%s', got '%s'", testPath, retrieved)
	}
}

func TestSetAndGetSettings(t *testing.T) {
	state := NewAppState()

	// Create custom settings
	customSettings := &Settings{
		Theme:    "dark",
		FontSize: 16,
	}

	state.SetSettings(customSettings)

	retrieved := state.GetSettings()
	if retrieved == nil {
		t.Fatal("Expected to retrieve settings")
	}

	if retrieved.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got '%s'", retrieved.Theme)
	}

	if retrieved.FontSize != 16 {
		t.Errorf("Expected font size 16, got %d", retrieved.FontSize)
	}
}

func TestNotifyObservers(t *testing.T) {
	state := NewAppState()

	// Track if observer was called
	observerCalled := false
	var callCount int

	// Subscribe observer
	state.Subscribe(func() {
		observerCalled = true
		callCount++
	})

	// Trigger notification
	state.Notify()

	if !observerCalled {
		t.Error("Expected observer to be called")
	}

	if callCount != 1 {
		t.Errorf("Expected observer to be called once, was called %d times", callCount)
	}

	// Trigger again
	state.Notify()

	if callCount != 2 {
		t.Errorf("Expected observer to be called twice, was called %d times", callCount)
	}
}

func TestMultipleObservers(t *testing.T) {
	state := NewAppState()

	// Track calls for multiple observers
	observer1Called := false
	observer2Called := false

	// Subscribe multiple observers
	state.Subscribe(func() {
		observer1Called = true
	})

	state.Subscribe(func() {
		observer2Called = true
	})

	// Trigger notification
	state.Notify()

	if !observer1Called {
		t.Error("Expected first observer to be called")
	}

	if !observer2Called {
		t.Error("Expected second observer to be called")
	}
}

func TestConcurrentAccess(t *testing.T) {
	state := NewAppState()

	// Test concurrent reads and writes
	done := make(chan bool)

	// Goroutine 1: Write project
	go func() {
		for i := 0; i < 100; i++ {
			project := &pom.Project{
				GroupID:    "com.test",
				ArtifactID: "test-app",
				Version:    "1.0.0",
			}
			state.SetCurrentProject(project)
		}
		done <- true
	}()

	// Goroutine 2: Read project
	go func() {
		for i := 0; i < 100; i++ {
			_ = state.GetCurrentProject()
		}
		done <- true
	}()

	// Goroutine 3: Write settings
	go func() {
		for i := 0; i < 100; i++ {
			settings := NewSettings()
			state.SetSettings(settings)
		}
		done <- true
	}()

	// Wait for all goroutines
	<-done
	<-done
	<-done

	// If we reach here without race conditions, test passes
}

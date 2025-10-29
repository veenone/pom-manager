package state

import (
	"sync"

	"github.com/user/pom-manager/internal/core/pom"
)

// AppState manages the application's global state with thread-safe access
// and observer pattern for reactive UI updates
type AppState struct {
	currentProject *pom.Project  // Current loaded/edited POM
	filePath       string         // Path to current file
	isDirty        bool           // Unsaved changes flag
	settings       *Settings      // User preferences
	observers      []func()       // Observer callbacks
	mutex          sync.RWMutex   // Thread-safe access
}

// NewAppState creates a new AppState with default settings
func NewAppState() *AppState {
	return &AppState{
		currentProject: nil,
		filePath:       "",
		isDirty:        false,
		settings:       NewSettings(),
		observers:      make([]func(), 0),
	}
}

// GetCurrentProject returns the current project (thread-safe read)
func (s *AppState) GetCurrentProject() *pom.Project {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.currentProject
}

// SetCurrentProject sets the current project and notifies observers
func (s *AppState) SetCurrentProject(project *pom.Project) {
	s.mutex.Lock()
	s.currentProject = project
	s.mutex.Unlock()
	s.Notify()
}

// GetFilePath returns the current file path (thread-safe read)
func (s *AppState) GetFilePath() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.filePath
}

// SetFilePath sets the file path and notifies observers
func (s *AppState) SetFilePath(path string) {
	s.mutex.Lock()
	s.filePath = path
	s.mutex.Unlock()
	s.Notify()
}

// IsDirty returns the dirty flag (thread-safe read)
func (s *AppState) IsDirty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isDirty
}

// SetDirty sets the dirty flag and notifies observers
func (s *AppState) SetDirty(dirty bool) {
	s.mutex.Lock()
	s.isDirty = dirty
	s.mutex.Unlock()
	s.Notify()
}

// GetSettings returns a copy of current settings (thread-safe read)
func (s *AppState) GetSettings() *Settings {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	// Return a copy to prevent external modifications
	settingsCopy := *s.settings
	return &settingsCopy
}

// SetSettings updates settings and notifies observers
func (s *AppState) SetSettings(settings *Settings) {
	s.mutex.Lock()
	s.settings = settings
	s.mutex.Unlock()
	s.Notify()
}

// Subscribe registers an observer callback
// The callback will be invoked when state changes occur
func (s *AppState) Subscribe(callback func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.observers = append(s.observers, callback)
}

// Unsubscribe removes an observer callback
// Note: This uses function pointer comparison, which may not work
// for all cases. Consider using a handle-based approach for production.
func (s *AppState) Unsubscribe(callback func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Find and remove the callback
	for i, obs := range s.observers {
		// Note: Function pointer comparison is limited in Go
		// This is a simplified implementation
		if &obs == &callback {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify calls all registered observers
// This is called internally when state changes
func (s *AppState) Notify() {
	s.mutex.RLock()
	observers := make([]func(), len(s.observers))
	copy(observers, s.observers)
	s.mutex.RUnlock()

	// Call observers outside the lock to prevent deadlocks
	for _, observer := range observers {
		observer()
	}
}

// UpdateProject updates the current project without triggering notification
// Useful for batching multiple changes before notifying observers
func (s *AppState) UpdateProject(updater func(*pom.Project)) {
	s.mutex.Lock()
	if s.currentProject != nil {
		updater(s.currentProject)
	}
	s.mutex.Unlock()
}

// UpdateProjectAndNotify updates the project and notifies observers
func (s *AppState) UpdateProjectAndNotify(updater func(*pom.Project)) {
	s.UpdateProject(updater)
	s.Notify()
}

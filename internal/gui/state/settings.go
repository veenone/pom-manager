package state

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Settings holds user preferences for the GUI application
type Settings struct {
	// General settings
	Theme            string `yaml:"theme"`             // "light" | "dark"
	AutoSaveInterval int    `yaml:"auto_save_interval"` // Minutes (0 = disabled)
	RestoreSession   bool   `yaml:"restore_session"`   // Restore last file on startup

	// Editor settings
	FontSize         int  `yaml:"font_size"`         // 10-18 pt
	LivePreview      bool `yaml:"live_preview"`      // Enable real-time preview
	ValidationDelay  int  `yaml:"validation_delay"`  // Milliseconds
	SyntaxHighlight  bool `yaml:"syntax_highlight"`  // Enable XML syntax highlighting

	// Templates settings
	DefaultTemplate   string `yaml:"default_template"`    // Default template name
	CustomTemplateDir string `yaml:"custom_template_dir"` // Path to custom templates

	// Advanced settings
	MavenCentralTimeout int    `yaml:"maven_central_timeout"` // Seconds
	EnableDebugLog      bool   `yaml:"enable_debug_log"`      // Debug logging
	CacheDir            string `yaml:"cache_dir"`             // Cache directory path

	// Window settings
	WindowWidth  int `yaml:"window_width"`  // Last window width
	WindowHeight int `yaml:"window_height"` // Last window height
	WindowX      int `yaml:"window_x"`      // Last window X position
	WindowY      int `yaml:"window_y"`      // Last window Y position

	// Session restore
	LastOpenedFile string   `yaml:"last_opened_file"` // Last opened file path
	RecentFiles    []string `yaml:"recent_files"`     // List of recently opened files
}

// NewSettings creates Settings with default values
func NewSettings() *Settings {
	return &Settings{
		// General defaults
		Theme:            "light",
		AutoSaveInterval: 5, // 5 minutes
		RestoreSession:   true,

		// Editor defaults
		FontSize:         12,
		LivePreview:      true,
		ValidationDelay:  100, // 100ms
		SyntaxHighlight:  true,

		// Templates defaults
		DefaultTemplate:   "basic-java",
		CustomTemplateDir: "",

		// Advanced defaults
		MavenCentralTimeout: 10, // 10 seconds
		EnableDebugLog:      false,
		CacheDir:            "", // Will use default ~/.pom-manager/cache

		// Window defaults
		WindowWidth:  1024,
		WindowHeight: 768,
		WindowX:      0,
		WindowY:      0,

		// Session defaults
		LastOpenedFile: "",
		RecentFiles:    []string{},
	}
}

// AddRecentFile adds a file to the recent files list (max 10)
func (s *Settings) AddRecentFile(filePath string) {
	// Remove if already exists
	for i, path := range s.RecentFiles {
		if path == filePath {
			s.RecentFiles = append(s.RecentFiles[:i], s.RecentFiles[i+1:]...)
			break
		}
	}

	// Add to front
	s.RecentFiles = append([]string{filePath}, s.RecentFiles...)

	// Keep only last 10
	if len(s.RecentFiles) > 10 {
		s.RecentFiles = s.RecentFiles[:10]
	}
}

// GetRecentFiles returns the list of recent files
func (s *Settings) GetRecentFiles() []string {
	// Filter out files that don't exist anymore
	validFiles := []string{}
	for _, path := range s.RecentFiles {
		if _, err := os.Stat(path); err == nil {
			validFiles = append(validFiles, path)
		}
	}
	s.RecentFiles = validFiles
	return validFiles
}

// GetConfigDir returns the config directory path (~/.pom-manager)
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".pom-manager")
	return configDir, nil
}

// GetConfigFilePath returns the full path to the GUI config file
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "gui-config.yaml"), nil
}

// LoadSettings loads settings from the config file
// If the file doesn't exist or can't be read, returns default settings
func LoadSettings() (*Settings, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return NewSettings(), fmt.Errorf("failed to get config path: %w", err)
	}

	// If config file doesn't exist, return defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return NewSettings(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return NewSettings(), fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal YAML
	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return NewSettings(), fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate loaded settings
	if err := validateSettings(&settings); err != nil {
		return NewSettings(), fmt.Errorf("invalid settings in config file: %w", err)
	}

	return &settings, nil
}

// SaveSettings saves settings to the config file
func SaveSettings(settings *Settings) error {
	// Validate settings before saving
	if err := validateSettings(settings); err != nil {
		return fmt.Errorf("invalid settings: %w", err)
	}

	configPath, err := GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal settings to YAML
	data, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// validateSettings validates that settings values are within acceptable ranges
func validateSettings(s *Settings) error {
	if s.FontSize < 10 || s.FontSize > 18 {
		return fmt.Errorf("font size must be between 10 and 18")
	}
	if s.AutoSaveInterval < 0 {
		return fmt.Errorf("auto-save interval must be non-negative")
	}
	if s.ValidationDelay < 0 || s.ValidationDelay > 5000 {
		return fmt.Errorf("validation delay must be between 0 and 5000 ms")
	}
	if s.MavenCentralTimeout < 1 || s.MavenCentralTimeout > 300 {
		return fmt.Errorf("Maven Central timeout must be between 1 and 300 seconds")
	}
	if s.Theme != "light" && s.Theme != "dark" {
		return fmt.Errorf("theme must be 'light' or 'dark'")
	}
	return nil
}

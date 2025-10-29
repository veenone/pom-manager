package pom

import (
	"fmt"
	"os"
	"path/filepath"
)

// Repository interface for file I/O operations
type Repository interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
	Exists(path string) bool
}

// fileRepository implements Repository using the file system
type fileRepository struct{}

// NewRepository creates a new file system repository
func NewRepository() Repository {
	return &fileRepository{}
}

// Read reads file contents
func (r *fileRepository) Read(path string) ([]byte, error) {
	// Check file size
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, path)
		}
		return nil, fmt.Errorf("stat file %s: %w", path, err)
	}

	if info.Size() > MaxFileSizeBytes {
		return nil, fmt.Errorf("%w: file %s size %d exceeds maximum %d bytes",
			ErrFileTooBig, path, info.Size(), MaxFileSizeBytes)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, path)
		}
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	return data, nil
}

// Write writes data to file, creating directories if needed
func (r *fileRepository) Write(path string, data []byte) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("%w: %s", ErrPermissionDenied, dir)
		}
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(path, data, 0644); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("%w: %s", ErrPermissionDenied, path)
		}
		return fmt.Errorf("writing file %s: %w", path, err)
	}

	return nil
}

// Exists checks if a file exists
func (r *fileRepository) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

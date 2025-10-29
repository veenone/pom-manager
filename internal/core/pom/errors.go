package pom

import "errors"

// Parsing errors
var (
	// ErrInvalidXML indicates malformed XML structure
	ErrInvalidXML = errors.New("invalid XML structure")

	// ErrMissingRequired indicates missing required fields
	ErrMissingRequired = errors.New("missing required fields")

	// ErrInvalidFormat indicates invalid format for a field
	ErrInvalidFormat = errors.New("invalid format")
)

// File operation errors
var (
	// ErrFileTooBig indicates file size exceeds limit
	ErrFileTooBig = errors.New("file size exceeds limit")

	// ErrFileNotFound indicates file does not exist
	ErrFileNotFound = errors.New("file not found")

	// ErrPermissionDenied indicates insufficient permissions
	ErrPermissionDenied = errors.New("permission denied")
)

// Validation errors
var (
	// ErrCircularDependency indicates a circular dependency was detected
	ErrCircularDependency = errors.New("circular dependency detected")

	// ErrInvalidScope indicates an invalid dependency scope
	ErrInvalidScope = errors.New("invalid dependency scope")

	// ErrInvalidPackaging indicates an invalid packaging type
	ErrInvalidPackaging = errors.New("invalid packaging type")

	// ErrInvalidPhase indicates an invalid Maven lifecycle phase
	ErrInvalidPhase = errors.New("invalid Maven lifecycle phase")

	// ErrInvalidProject indicates the project struct failed validation
	ErrInvalidProject = errors.New("invalid project structure")
)

// Generation errors
var (
	// ErrGenerationFailed indicates XML generation failed
	ErrGenerationFailed = errors.New("XML generation failed")
)

// Template errors
var (
	// ErrTemplateNotFound indicates unknown template name
	ErrTemplateNotFound = errors.New("template not found")
)

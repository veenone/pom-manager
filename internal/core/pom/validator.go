package pom

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// Validator interface for validating Project structs
type Validator interface {
	Validate(project *Project) ValidationResult
}

// ValidationRule interface for individual validation rules
type ValidationRule interface {
	Validate(project *Project) []ValidationError
}

// defaultValidator implements Validator
type defaultValidator struct {
	rules []ValidationRule
}

// NewValidator creates a new Validator with all validation rules
func NewValidator() Validator {
	return &defaultValidator{
		rules: []ValidationRule{
			&coordinatesRule{},
			&dependenciesRule{},
			&buildRule{},
		},
	}
}

// Validate runs all validation rules and returns grouped errors
func (v *defaultValidator) Validate(project *Project) ValidationResult {
	result := ValidationResult{
		Valid: true,
		Errors: ValidationErrors{
			Coordinates:  []ValidationError{},
			Dependencies: []ValidationError{},
			Build:        []ValidationError{},
			General:      []ValidationError{},
		},
	}

	if project == nil {
		result.Valid = false
		result.Errors.General = append(result.Errors.General, ValidationError{
			Field:   "project",
			Value:   "nil",
			Message: "project cannot be nil",
		})
		return result
	}

	// Run all validation rules
	for _, rule := range v.rules {
		errors := rule.Validate(project)
		for _, err := range errors {
			result.Valid = false
			// Categorize errors based on field
			if strings.HasPrefix(err.Field, "groupId") || strings.HasPrefix(err.Field, "artifactId") || strings.HasPrefix(err.Field, "version") || strings.HasPrefix(err.Field, "packaging") {
				result.Errors.Coordinates = append(result.Errors.Coordinates, err)
			} else if strings.Contains(err.Field, "dependency") || strings.Contains(err.Field, "scope") {
				result.Errors.Dependencies = append(result.Errors.Dependencies, err)
			} else if strings.Contains(err.Field, "plugin") || strings.Contains(err.Field, "phase") || strings.Contains(err.Field, "build") {
				result.Errors.Build = append(result.Errors.Build, err)
			} else {
				result.Errors.General = append(result.Errors.General, err)
			}
		}
	}

	return result
}

// coordinatesRule validates project coordinates
type coordinatesRule struct{}

func (r *coordinatesRule) Validate(project *Project) []ValidationError {
	var errors []ValidationError

	// Validate groupId
	if project.GroupID == "" {
		errors = append(errors, ValidationError{
			Field:   "groupId",
			Value:   "",
			Message: "groupId is required",
		})
	} else if !isValidGroupID(project.GroupID) {
		errors = append(errors, ValidationError{
			Field:   "groupId",
			Value:   project.GroupID,
			Message: "groupId must be lowercase with dot separators (e.g., 'com.example')",
		})
	}

	// Validate artifactId
	if project.ArtifactID == "" {
		errors = append(errors, ValidationError{
			Field:   "artifactId",
			Value:   "",
			Message: "artifactId is required",
		})
	} else if !isValidArtifactID(project.ArtifactID) {
		errors = append(errors, ValidationError{
			Field:   "artifactId",
			Value:   project.ArtifactID,
			Message: "artifactId must be lowercase with hyphens (e.g., 'my-app')",
		})
	}

	// Validate version
	if project.Version == "" {
		errors = append(errors, ValidationError{
			Field:   "version",
			Value:   "",
			Message: "version is required",
		})
	} else if !isValidVersion(project.Version) {
		errors = append(errors, ValidationError{
			Field:   "version",
			Value:   project.Version,
			Message: "version must follow semantic versioning or Maven snapshot conventions",
		})
	}

	// Validate packaging
	if project.Packaging != "" && !isValidPackaging(project.Packaging) {
		errors = append(errors, ValidationError{
			Field:   "packaging",
			Value:   project.Packaging,
			Message: fmt.Sprintf("packaging must be one of: %s", strings.Join(ValidPackagingTypes, ", ")),
		})
	}

	return errors
}

// dependenciesRule validates dependencies
type dependenciesRule struct{}

func (r *dependenciesRule) Validate(project *Project) []ValidationError {
	var errors []ValidationError

	// Check for circular dependencies (simplified - checks direct duplicates)
	seen := make(map[string]bool)
	for i, dep := range project.Dependencies {
		// Validate required fields
		if dep.GroupID == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("dependencies[%d].groupId", i),
				Value:   "",
				Message: "dependency groupId is required",
			})
		}
		if dep.ArtifactID == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("dependencies[%d].artifactId", i),
				Value:   "",
				Message: "dependency artifactId is required",
			})
		}
		if dep.Version == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("dependencies[%d].version", i),
				Value:   "",
				Message: "dependency version is required",
			})
		}

		// Validate scope
		if dep.Scope != "" && !isValidScope(dep.Scope) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("dependencies[%d].scope", i),
				Value:   dep.Scope,
				Message: fmt.Sprintf("scope must be one of: %s", strings.Join(ValidDependencyScopes, ", ")),
			})
		}

		// Check for duplicates (simple circular dependency detection)
		key := fmt.Sprintf("%s:%s", dep.GroupID, dep.ArtifactID)
		if seen[key] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("dependencies[%d]", i),
				Value:   key,
				Message: "duplicate dependency detected",
			})
		}
		seen[key] = true
	}

	return errors
}

// buildRule validates build configuration
type buildRule struct{}

func (r *buildRule) Validate(project *Project) []ValidationError {
	var errors []ValidationError

	if project.Build == nil {
		return errors
	}

	// Validate plugins
	for i, plugin := range project.Build.Plugins {
		if plugin.GroupID == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("build.plugins[%d].groupId", i),
				Value:   "",
				Message: "plugin groupId is required",
			})
		}
		if plugin.ArtifactID == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("build.plugins[%d].artifactId", i),
				Value:   "",
				Message: "plugin artifactId is required",
			})
		}

		// Validate executions
		for j, exec := range plugin.Executions {
			if exec.Phase != "" && !isValidPhase(exec.Phase) {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("build.plugins[%d].executions[%d].phase", i, j),
					Value:   exec.Phase,
					Message: "phase must be a valid Maven lifecycle phase",
				})
			}
		}
	}

	return errors
}

// isValidGroupID checks if groupId follows Maven conventions
func isValidGroupID(groupID string) bool {
	// Allow lowercase letters, numbers, dots, and hyphens
	// Must not start or end with dot
	pattern := `^[a-z0-9][a-z0-9\-.]*[a-z0-9]$`
	matched, _ := regexp.MatchString(pattern, groupID)
	return matched && strings.Contains(groupID, ".")
}

// isValidArtifactID checks if artifactId follows Maven conventions
func isValidArtifactID(artifactID string) bool {
	// Allow lowercase letters, numbers, and hyphens
	pattern := `^[a-z0-9][a-z0-9\-]*[a-z0-9]$`
	matched, _ := regexp.MatchString(pattern, artifactID)
	return matched
}

// isValidVersion checks if version follows semver or Maven snapshot conventions
func isValidVersion(version string) bool {
	// Try semver first
	if _, err := semver.NewVersion(version); err == nil {
		return true
	}

	// Check for Maven snapshot version
	if strings.HasSuffix(version, "-SNAPSHOT") {
		base := strings.TrimSuffix(version, "-SNAPSHOT")
		if _, err := semver.NewVersion(base); err == nil {
			return true
		}
	}

	// Allow simple version patterns like "1.0", "1.0.0", etc.
	pattern := `^\d+(\.\d+)*(-SNAPSHOT)?$`
	matched, _ := regexp.MatchString(pattern, version)
	return matched
}

// isValidPackaging checks if packaging type is valid
func isValidPackaging(packaging string) bool {
	for _, valid := range ValidPackagingTypes {
		if packaging == valid {
			return true
		}
	}
	return false
}

// isValidScope checks if dependency scope is valid
func isValidScope(scope string) bool {
	for _, valid := range ValidDependencyScopes {
		if scope == valid {
			return true
		}
	}
	return false
}

// isValidPhase checks if Maven lifecycle phase is valid
func isValidPhase(phase string) bool {
	for _, valid := range MavenLifecyclePhases {
		if phase == valid {
			return true
		}
	}
	return false
}

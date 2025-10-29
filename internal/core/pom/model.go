package pom

import "fmt"

// Project represents a complete Maven POM
type Project struct {
	XMLName      string                 `xml:"project"`
	XMLNS        string                 `xml:"xmlns,attr"`
	XSI          string                 `xml:"xmlns:xsi,attr"`
	SchemaLocation string               `xml:"xsi:schemaLocation,attr"`
	ModelVersion string                 `xml:"modelVersion" validate:"required"`
	Coordinates  Coordinates            `xml:"-" validate:"required"`
	GroupID      string                 `xml:"groupId" validate:"required"`
	ArtifactID   string                 `xml:"artifactId" validate:"required"`
	Version      string                 `xml:"version" validate:"required"`
	Packaging    string                 `xml:"packaging,omitempty"`
	Name         string                 `xml:"name,omitempty"`
	Description  string                 `xml:"description,omitempty"`
	Properties   map[string]string      `xml:"-"`
	PropertiesXML *Properties           `xml:"properties,omitempty"`
	Dependencies []Dependency           `xml:"dependencies>dependency,omitempty"`
	Build        *Build                 `xml:"build,omitempty"`
	Modules      []string               `xml:"modules>module,omitempty"`
	Parent       *Parent                `xml:"parent,omitempty"`
	Profiles     []Profile              `xml:"profiles>profile,omitempty"`
}

// Properties represents Maven properties as a map
type Properties struct {
	Entries map[string]string
}

// Coordinates uniquely identify a Maven artifact
type Coordinates struct {
	GroupID    string `validate:"required"`
	ArtifactID string `validate:"required"`
	Version    string `validate:"required"`
}

// String returns coordinates in standard Maven format "groupId:artifactId:version"
func (c Coordinates) String() string {
	return fmt.Sprintf("%s:%s:%s", c.GroupID, c.ArtifactID, c.Version)
}

// Dependency represents a Maven dependency
type Dependency struct {
	GroupID    string      `xml:"groupId" validate:"required"`
	ArtifactID string      `xml:"artifactId" validate:"required"`
	Version    string      `xml:"version" validate:"required"`
	Scope      string      `xml:"scope,omitempty"`
	Optional   bool        `xml:"optional,omitempty"`
	Exclusions []Exclusion `xml:"exclusions>exclusion,omitempty"`
}

// Exclusion represents an excluded transitive dependency
type Exclusion struct {
	GroupID    string `xml:"groupId" validate:"required"`
	ArtifactID string `xml:"artifactId" validate:"required"`
}

// Build represents Maven build configuration
type Build struct {
	SourceDirectory     string   `xml:"sourceDirectory,omitempty"`
	TestSourceDirectory string   `xml:"testSourceDirectory,omitempty"`
	OutputDirectory     string   `xml:"outputDirectory,omitempty"`
	Plugins             []Plugin `xml:"plugins>plugin,omitempty"`
}

// Plugin represents a Maven plugin
type Plugin struct {
	GroupID       string            `xml:"groupId" validate:"required"`
	ArtifactID    string            `xml:"artifactId" validate:"required"`
	Version       string            `xml:"version,omitempty"`
	Configuration *Configuration    `xml:"configuration,omitempty"`
	Executions    []PluginExecution `xml:"executions>execution,omitempty"`
}

// PluginExecution represents a plugin execution
type PluginExecution struct {
	ID            string         `xml:"id,omitempty"`
	Phase         string         `xml:"phase,omitempty"`
	Goals         []string       `xml:"goals>goal,omitempty"`
	Configuration *Configuration `xml:"configuration,omitempty"`
}

// Configuration represents generic plugin or execution configuration
// This is a simplified representation - real Maven configs can be complex nested XML
type Configuration struct {
	Data map[string]interface{}
}

// Parent represents a parent POM reference
type Parent struct {
	GroupID      string `xml:"groupId" validate:"required"`
	ArtifactID   string `xml:"artifactId" validate:"required"`
	Version      string `xml:"version" validate:"required"`
	RelativePath string `xml:"relativePath,omitempty"`
}

// Profile represents a Maven build profile
type Profile struct {
	ID           string            `xml:"id" validate:"required"`
	Activation   *Activation       `xml:"activation,omitempty"`
	Properties   map[string]string `xml:"-"`
	PropertiesXML *Properties      `xml:"properties,omitempty"`
	Dependencies []Dependency      `xml:"dependencies>dependency,omitempty"`
	Build        *Build            `xml:"build,omitempty"`
	Modules      []string          `xml:"modules>module,omitempty"`
}

// Activation defines when a profile should be active
type Activation struct {
	ActiveByDefault bool   `xml:"activeByDefault,omitempty"`
	JDK             string `xml:"jdk,omitempty"`
	Property        *ActivationProperty `xml:"property,omitempty"`
	OS              *ActivationOS `xml:"os,omitempty"`
	File            *ActivationFile `xml:"file,omitempty"`
}

// ActivationProperty represents property-based activation
type ActivationProperty struct {
	Name  string `xml:"name"`
	Value string `xml:"value,omitempty"`
}

// ActivationOS represents OS-based activation
type ActivationOS struct {
	Name    string `xml:"name,omitempty"`
	Family  string `xml:"family,omitempty"`
	Arch    string `xml:"arch,omitempty"`
	Version string `xml:"version,omitempty"`
}

// ActivationFile represents file-based activation
type ActivationFile struct {
	Exists  string `xml:"exists,omitempty"`
	Missing string `xml:"missing,omitempty"`
}

// ValidationResult contains validation errors grouped by category
type ValidationResult struct {
	Valid  bool
	Errors ValidationErrors
}

// ValidationErrors groups errors by category
type ValidationErrors struct {
	Coordinates  []ValidationError
	Dependencies []ValidationError
	Build        []ValidationError
	General      []ValidationError
}

// ValidationError represents a single validation failure
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

// Error returns formatted error message
func (v ValidationError) Error() string {
	return fmt.Sprintf("field '%s' with value '%s': %s", v.Field, v.Value, v.Message)
}

// HasErrors returns true if there are any validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve.Coordinates) > 0 || len(ve.Dependencies) > 0 || len(ve.Build) > 0 || len(ve.General) > 0
}

// AllErrors returns all validation errors as a flat slice
func (ve ValidationErrors) AllErrors() []ValidationError {
	var all []ValidationError
	all = append(all, ve.Coordinates...)
	all = append(all, ve.Dependencies...)
	all = append(all, ve.Build...)
	all = append(all, ve.General...)
	return all
}

// TemplateInfo provides information about a POM template
type TemplateInfo struct {
	Name        string
	Description string
}

package pom

import (
	"fmt"
	"sort"

	"github.com/beevik/etree"
)

// Generator interface for generating Maven POM XML
type Generator interface {
	Generate(project *Project) ([]byte, error)
	GenerateToFile(project *Project, path string) error
}

// defaultGenerator implements Generator using etree
type defaultGenerator struct {
	repo Repository
}

// NewGenerator creates a new Generator instance
func NewGenerator() Generator {
	return &defaultGenerator{
		repo: NewRepository(),
	}
}

// NewGeneratorWithRepo creates a new Generator with custom repository (for testing)
func NewGeneratorWithRepo(repo Repository) Generator {
	return &defaultGenerator{
		repo: repo,
	}
}

// Generate generates XML bytes from a Project struct
func (g *defaultGenerator) Generate(project *Project) ([]byte, error) {
	if project == nil {
		return nil, fmt.Errorf("%w: project is nil", ErrInvalidProject)
	}

	// Basic validation - check required fields
	if project.GroupID == "" || project.ArtifactID == "" || project.Version == "" {
		return nil, fmt.Errorf("%w: missing required fields (groupId, artifactId, or version)", ErrMissingRequired)
	}

	// Create XML document
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	// Create project root element
	root := doc.CreateElement("project")
	root.CreateAttr("xmlns", MavenXMLNamespace)
	root.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	root.CreateAttr("xsi:schemaLocation", MavenXMLSchemaLocation)

	// Add model version
	modelVersion := root.CreateElement("modelVersion")
	if project.ModelVersion != "" {
		modelVersion.SetText(project.ModelVersion)
	} else {
		modelVersion.SetText(DefaultModelVersion)
	}

	// Add parent if present
	if project.Parent != nil {
		g.addParent(root, project.Parent)
	}

	// Add coordinates
	groupID := root.CreateElement("groupId")
	groupID.SetText(project.GroupID)

	artifactID := root.CreateElement("artifactId")
	artifactID.SetText(project.ArtifactID)

	version := root.CreateElement("version")
	version.SetText(project.Version)

	// Add packaging
	if project.Packaging != "" && project.Packaging != DefaultPackaging {
		packaging := root.CreateElement("packaging")
		packaging.SetText(project.Packaging)
	}

	// Add optional name and description
	if project.Name != "" {
		name := root.CreateElement("name")
		name.SetText(project.Name)
	}

	if project.Description != "" {
		desc := root.CreateElement("description")
		desc.SetText(project.Description)
	}

	// Add modules if present
	if len(project.Modules) > 0 {
		modules := root.CreateElement("modules")
		for _, mod := range project.Modules {
			module := modules.CreateElement("module")
			module.SetText(mod)
		}
	}

	// Add properties (sorted alphabetically for consistent output)
	if len(project.Properties) > 0 {
		properties := root.CreateElement("properties")

		// Sort property keys alphabetically
		keys := make([]string, 0, len(project.Properties))
		for key := range project.Properties {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// Add properties in sorted order
		for _, key := range keys {
			prop := properties.CreateElement(key)
			prop.SetText(project.Properties[key])
		}
	}

	// Add dependencies
	if len(project.Dependencies) > 0 {
		dependencies := root.CreateElement("dependencies")
		for _, dep := range project.Dependencies {
			g.addDependency(dependencies, dep)
		}
	}

	// Add build
	if project.Build != nil {
		g.addBuild(root, project.Build)
	}

	// Set indentation for pretty-print (4 spaces)
	doc.Indent(4)

	// Generate XML
	xmlBytes, err := doc.WriteToBytes()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGenerationFailed, err)
	}

	return xmlBytes, nil
}

// GenerateToFile generates XML and writes to file
func (g *defaultGenerator) GenerateToFile(project *Project, path string) error {
	xmlBytes, err := g.Generate(project)
	if err != nil {
		return err
	}

	if err := g.repo.Write(path, xmlBytes); err != nil {
		return fmt.Errorf("writing file %s: %w", path, err)
	}

	return nil
}

// addDependency adds a dependency element
func (g *defaultGenerator) addDependency(parent *etree.Element, dep Dependency) {
	dependency := parent.CreateElement("dependency")

	groupID := dependency.CreateElement("groupId")
	groupID.SetText(dep.GroupID)

	artifactID := dependency.CreateElement("artifactId")
	artifactID.SetText(dep.ArtifactID)

	version := dependency.CreateElement("version")
	version.SetText(dep.Version)

	if dep.Scope != "" && dep.Scope != DefaultScope {
		scope := dependency.CreateElement("scope")
		scope.SetText(dep.Scope)
	}

	if dep.Optional {
		optional := dependency.CreateElement("optional")
		optional.SetText("true")
	}

	// Add exclusions
	if len(dep.Exclusions) > 0 {
		exclusions := dependency.CreateElement("exclusions")
		for _, excl := range dep.Exclusions {
			exclusion := exclusions.CreateElement("exclusion")
			exclGroupID := exclusion.CreateElement("groupId")
			exclGroupID.SetText(excl.GroupID)
			exclArtifactID := exclusion.CreateElement("artifactId")
			exclArtifactID.SetText(excl.ArtifactID)
		}
	}
}

// addBuild adds a build element
func (g *defaultGenerator) addBuild(parent *etree.Element, build *Build) {
	buildElem := parent.CreateElement("build")

	if build.SourceDirectory != "" {
		sourceDir := buildElem.CreateElement("sourceDirectory")
		sourceDir.SetText(build.SourceDirectory)
	}

	if build.TestSourceDirectory != "" {
		testSourceDir := buildElem.CreateElement("testSourceDirectory")
		testSourceDir.SetText(build.TestSourceDirectory)
	}

	if build.OutputDirectory != "" {
		outputDir := buildElem.CreateElement("outputDirectory")
		outputDir.SetText(build.OutputDirectory)
	}

	// Add plugins
	if len(build.Plugins) > 0 {
		plugins := buildElem.CreateElement("plugins")
		for _, plugin := range build.Plugins {
			g.addPlugin(plugins, plugin)
		}
	}
}

// addPlugin adds a plugin element
func (g *defaultGenerator) addPlugin(parent *etree.Element, plugin Plugin) {
	pluginElem := parent.CreateElement("plugin")

	groupID := pluginElem.CreateElement("groupId")
	groupID.SetText(plugin.GroupID)

	artifactID := pluginElem.CreateElement("artifactId")
	artifactID.SetText(plugin.ArtifactID)

	if plugin.Version != "" {
		version := pluginElem.CreateElement("version")
		version.SetText(plugin.Version)
	}

	// Add executions
	if len(plugin.Executions) > 0 {
		executions := pluginElem.CreateElement("executions")
		for _, exec := range plugin.Executions {
			g.addExecution(executions, exec)
		}
	}
}

// addExecution adds an execution element
func (g *defaultGenerator) addExecution(parent *etree.Element, exec PluginExecution) {
	execElem := parent.CreateElement("execution")

	if exec.ID != "" {
		id := execElem.CreateElement("id")
		id.SetText(exec.ID)
	}

	if exec.Phase != "" {
		phase := execElem.CreateElement("phase")
		phase.SetText(exec.Phase)
	}

	// Add goals
	if len(exec.Goals) > 0 {
		goals := execElem.CreateElement("goals")
		for _, goal := range exec.Goals {
			goalElem := goals.CreateElement("goal")
			goalElem.SetText(goal)
		}
	}
}

// addParent adds a parent element
func (g *defaultGenerator) addParent(parent *etree.Element, p *Parent) {
	parentElem := parent.CreateElement("parent")

	groupID := parentElem.CreateElement("groupId")
	groupID.SetText(p.GroupID)

	artifactID := parentElem.CreateElement("artifactId")
	artifactID.SetText(p.ArtifactID)

	version := parentElem.CreateElement("version")
	version.SetText(p.Version)

	if p.RelativePath != "" {
		relativePath := parentElem.CreateElement("relativePath")
		relativePath.SetText(p.RelativePath)
	}
}

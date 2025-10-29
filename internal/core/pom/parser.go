package pom

import (
	"fmt"
	"os"

	"github.com/beevik/etree"
)

// Parser interface for parsing Maven POM files
type Parser interface {
	Parse(xmlData []byte) (*Project, error)
	ParseFile(path string) (*Project, error)
}

// defaultParser implements Parser interface using etree
type defaultParser struct {
	repo Repository
}

// NewParser creates a new Parser instance
func NewParser() Parser {
	return &defaultParser{
		repo: NewRepository(),
	}
}

// NewParserWithRepo creates a new Parser with custom repository (for testing)
func NewParserWithRepo(repo Repository) Parser {
	return &defaultParser{
		repo: repo,
	}
}

// Parse parses XML bytes into a Project struct
func (p *defaultParser) Parse(xmlData []byte) (*Project, error) {
	// Check file size limit
	if len(xmlData) > MaxFileSizeBytes {
		return nil, fmt.Errorf("%w: size %d exceeds maximum %d bytes", ErrFileTooBig, len(xmlData), MaxFileSizeBytes)
	}

	// Parse XML
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidXML, err)
	}

	// Get root project element
	root := doc.SelectElement("project")
	if root == nil {
		return nil, fmt.Errorf("%w: missing <project> root element", ErrInvalidXML)
	}

	project := &Project{
		XMLNS:          MavenXMLNamespace,
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: MavenXMLSchemaLocation,
		ModelVersion:   DefaultModelVersion,
	}

	// Parse model version
	if modelVersion := root.SelectElement("modelVersion"); modelVersion != nil {
		project.ModelVersion = modelVersion.Text()
	}

	// Parse coordinates
	groupID := root.SelectElement("groupId")
	artifactID := root.SelectElement("artifactId")
	version := root.SelectElement("version")

	if groupID == nil || artifactID == nil || version == nil {
		return nil, fmt.Errorf("%w: missing required fields (groupId, artifactId, or version)", ErrMissingRequired)
	}

	project.GroupID = groupID.Text()
	project.ArtifactID = artifactID.Text()
	project.Version = version.Text()
	project.Coordinates = Coordinates{
		GroupID:    project.GroupID,
		ArtifactID: project.ArtifactID,
		Version:    project.Version,
	}

	// Parse optional fields
	if packaging := root.SelectElement("packaging"); packaging != nil {
		project.Packaging = packaging.Text()
	} else {
		project.Packaging = DefaultPackaging
	}

	if name := root.SelectElement("name"); name != nil {
		project.Name = name.Text()
	}

	if description := root.SelectElement("description"); description != nil {
		project.Description = description.Text()
	}

	// Parse properties
	if props := root.SelectElement("properties"); props != nil {
		project.Properties = make(map[string]string)
		for _, child := range props.ChildElements() {
			project.Properties[child.Tag] = child.Text()
		}
	}

	// Parse dependencies
	if dependencies := root.SelectElement("dependencies"); dependencies != nil {
		for _, dep := range dependencies.SelectElements("dependency") {
			dependency, err := p.parseDependency(dep)
			if err != nil {
				return nil, fmt.Errorf("parsing dependency: %w", err)
			}
			project.Dependencies = append(project.Dependencies, dependency)
		}
	}

	// Parse build
	if buildElem := root.SelectElement("build"); buildElem != nil {
		build, err := p.parseBuild(buildElem)
		if err != nil {
			return nil, fmt.Errorf("parsing build: %w", err)
		}
		project.Build = build
	}

	// Parse parent
	if parentElem := root.SelectElement("parent"); parentElem != nil {
		parent, err := p.parseParent(parentElem)
		if err != nil {
			return nil, fmt.Errorf("parsing parent: %w", err)
		}
		project.Parent = parent
	}

	// Parse modules
	if modulesElem := root.SelectElement("modules"); modulesElem != nil {
		for _, module := range modulesElem.SelectElements("module") {
			project.Modules = append(project.Modules, module.Text())
		}
	}

	// Parse profiles
	if profilesElem := root.SelectElement("profiles"); profilesElem != nil {
		for _, profileElem := range profilesElem.SelectElements("profile") {
			profile, err := p.parseProfile(profileElem)
			if err != nil {
				return nil, fmt.Errorf("parsing profile: %w", err)
			}
			project.Profiles = append(project.Profiles, profile)
		}
	}

	return project, nil
}

// ParseFile reads and parses a POM file
func (p *defaultParser) ParseFile(path string) (*Project, error) {
	// Check file size
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, path)
		}
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	if info.Size() > MaxFileSizeBytes {
		return nil, fmt.Errorf("%w: file %s size %d exceeds maximum %d bytes", ErrFileTooBig, path, info.Size(), MaxFileSizeBytes)
	}

	// Read file
	data, err := p.repo.Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	// Parse
	return p.Parse(data)
}

// parseDependency parses a dependency element
func (p *defaultParser) parseDependency(elem *etree.Element) (Dependency, error) {
	dep := Dependency{
		Scope: DefaultScope, // Default scope
	}

	groupID := elem.SelectElement("groupId")
	artifactID := elem.SelectElement("artifactId")
	version := elem.SelectElement("version")

	if groupID == nil || artifactID == nil || version == nil {
		return dep, fmt.Errorf("%w: dependency missing required fields", ErrMissingRequired)
	}

	dep.GroupID = groupID.Text()
	dep.ArtifactID = artifactID.Text()
	dep.Version = version.Text()

	if scope := elem.SelectElement("scope"); scope != nil {
		dep.Scope = scope.Text()
	}

	if optional := elem.SelectElement("optional"); optional != nil {
		dep.Optional = optional.Text() == "true"
	}

	// Parse exclusions
	if exclusions := elem.SelectElement("exclusions"); exclusions != nil {
		for _, excl := range exclusions.SelectElements("exclusion") {
			exclusion, err := p.parseExclusion(excl)
			if err != nil {
				return dep, fmt.Errorf("parsing exclusion: %w", err)
			}
			dep.Exclusions = append(dep.Exclusions, exclusion)
		}
	}

	return dep, nil
}

// parseExclusion parses an exclusion element
func (p *defaultParser) parseExclusion(elem *etree.Element) (Exclusion, error) {
	excl := Exclusion{}

	groupID := elem.SelectElement("groupId")
	artifactID := elem.SelectElement("artifactId")

	if groupID == nil || artifactID == nil {
		return excl, fmt.Errorf("%w: exclusion missing required fields", ErrMissingRequired)
	}

	excl.GroupID = groupID.Text()
	excl.ArtifactID = artifactID.Text()

	return excl, nil
}

// parseBuild parses a build element
func (p *defaultParser) parseBuild(elem *etree.Element) (*Build, error) {
	build := &Build{}

	if sourceDir := elem.SelectElement("sourceDirectory"); sourceDir != nil {
		build.SourceDirectory = sourceDir.Text()
	}

	if testSourceDir := elem.SelectElement("testSourceDirectory"); testSourceDir != nil {
		build.TestSourceDirectory = testSourceDir.Text()
	}

	if outputDir := elem.SelectElement("outputDirectory"); outputDir != nil {
		build.OutputDirectory = outputDir.Text()
	}

	// Parse plugins
	if plugins := elem.SelectElement("plugins"); plugins != nil {
		for _, pluginElem := range plugins.SelectElements("plugin") {
			plugin, err := p.parsePlugin(pluginElem)
			if err != nil {
				return nil, fmt.Errorf("parsing plugin: %w", err)
			}
			build.Plugins = append(build.Plugins, plugin)
		}
	}

	return build, nil
}

// parsePlugin parses a plugin element
func (p *defaultParser) parsePlugin(elem *etree.Element) (Plugin, error) {
	plugin := Plugin{}

	groupID := elem.SelectElement("groupId")
	artifactID := elem.SelectElement("artifactId")

	if groupID == nil || artifactID == nil {
		return plugin, fmt.Errorf("%w: plugin missing required fields", ErrMissingRequired)
	}

	plugin.GroupID = groupID.Text()
	plugin.ArtifactID = artifactID.Text()

	if version := elem.SelectElement("version"); version != nil {
		plugin.Version = version.Text()
	}

	// Parse executions
	if executions := elem.SelectElement("executions"); executions != nil {
		for _, exec := range executions.SelectElements("execution") {
			execution, err := p.parseExecution(exec)
			if err != nil {
				return plugin, fmt.Errorf("parsing execution: %w", err)
			}
			plugin.Executions = append(plugin.Executions, execution)
		}
	}

	return plugin, nil
}

// parseExecution parses an execution element
func (p *defaultParser) parseExecution(elem *etree.Element) (PluginExecution, error) {
	exec := PluginExecution{}

	if id := elem.SelectElement("id"); id != nil {
		exec.ID = id.Text()
	}

	if phase := elem.SelectElement("phase"); phase != nil {
		exec.Phase = phase.Text()
	}

	// Parse goals
	if goals := elem.SelectElement("goals"); goals != nil {
		for _, goal := range goals.SelectElements("goal") {
			exec.Goals = append(exec.Goals, goal.Text())
		}
	}

	return exec, nil
}

// parseParent parses a parent element
func (p *defaultParser) parseParent(elem *etree.Element) (*Parent, error) {
	parent := &Parent{}

	groupID := elem.SelectElement("groupId")
	artifactID := elem.SelectElement("artifactId")
	version := elem.SelectElement("version")

	if groupID == nil || artifactID == nil || version == nil {
		return nil, fmt.Errorf("%w: parent missing required fields", ErrMissingRequired)
	}

	parent.GroupID = groupID.Text()
	parent.ArtifactID = artifactID.Text()
	parent.Version = version.Text()

	if relativePath := elem.SelectElement("relativePath"); relativePath != nil {
		parent.RelativePath = relativePath.Text()
	}

	return parent, nil
}

// parseProfile parses a profile element
func (p *defaultParser) parseProfile(elem *etree.Element) (Profile, error) {
	profile := Profile{}

	// Parse ID (required)
	idElem := elem.SelectElement("id")
	if idElem == nil {
		return profile, fmt.Errorf("%w: profile missing required id", ErrMissingRequired)
	}
	profile.ID = idElem.Text()

	// Parse activation
	if activationElem := elem.SelectElement("activation"); activationElem != nil {
		activation := &Activation{}

		if activeByDefault := activationElem.SelectElement("activeByDefault"); activeByDefault != nil {
			activation.ActiveByDefault = activeByDefault.Text() == "true"
		}

		if jdk := activationElem.SelectElement("jdk"); jdk != nil {
			activation.JDK = jdk.Text()
		}

		if propertyElem := activationElem.SelectElement("property"); propertyElem != nil {
			prop := &ActivationProperty{}
			if name := propertyElem.SelectElement("name"); name != nil {
				prop.Name = name.Text()
			}
			if value := propertyElem.SelectElement("value"); value != nil {
				prop.Value = value.Text()
			}
			activation.Property = prop
		}

		profile.Activation = activation
	}

	// Parse properties
	if props := elem.SelectElement("properties"); props != nil {
		profile.Properties = make(map[string]string)
		for _, child := range props.ChildElements() {
			profile.Properties[child.Tag] = child.Text()
		}
	}

	// Parse dependencies
	if dependencies := elem.SelectElement("dependencies"); dependencies != nil {
		for _, dep := range dependencies.SelectElements("dependency") {
			dependency, err := p.parseDependency(dep)
			if err != nil {
				return profile, fmt.Errorf("parsing profile dependency: %w", err)
			}
			profile.Dependencies = append(profile.Dependencies, dependency)
		}
	}

	// Parse build
	if buildElem := elem.SelectElement("build"); buildElem != nil {
		build, err := p.parseBuild(buildElem)
		if err != nil {
			return profile, fmt.Errorf("parsing profile build: %w", err)
		}
		profile.Build = build
	}

	// Parse modules
	if modulesElem := elem.SelectElement("modules"); modulesElem != nil {
		for _, module := range modulesElem.SelectElements("module") {
			profile.Modules = append(profile.Modules, module.Text())
		}
	}

	return profile, nil
}

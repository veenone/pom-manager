package pom

import "fmt"

// TemplateManager interface for creating Projects from templates
type TemplateManager interface {
	Create(templateName string, coords Coordinates) (*Project, error)
	List() []TemplateInfo
}

// templateManager implements TemplateManager
type templateManager struct{}

// NewTemplateManager creates a new TemplateManager
func NewTemplateManager() TemplateManager {
	return &templateManager{}
}

// Create creates a new Project from a template
func (tm *templateManager) Create(templateName string, coords Coordinates) (*Project, error) {
	switch templateName {
	case "basic-java":
		return tm.createBasicJava(coords), nil
	case "java-library":
		return tm.createJavaLibrary(coords), nil
	case "web-app":
		return tm.createWebApp(coords), nil
	case "javacard":
		return tm.createJavaCard(coords), nil
	default:
		return nil, fmt.Errorf("%w: unknown template '%s', available templates: basic-java, java-library, web-app, javacard", ErrTemplateNotFound, templateName)
	}
}

// List returns all available templates
func (tm *templateManager) List() []TemplateInfo {
	return []TemplateInfo{
		{
			Name:        "basic-java",
			Description: "Basic Java JAR project with compiler plugin",
		},
		{
			Name:        "java-library",
			Description: "Java library project with compiler and JAR plugins",
		},
		{
			Name:        "web-app",
			Description: "Java web application (WAR) project",
		},
		{
			Name:        "javacard",
			Description: "JavaCard applet project for smart cards (CAP packaging)",
		},
	}
}

// createBasicJava creates a basic Java project template
func (tm *templateManager) createBasicJava(coords Coordinates) *Project {
	return &Project{
		XMLNS:          MavenXMLNamespace,
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: MavenXMLSchemaLocation,
		ModelVersion:   DefaultModelVersion,
		GroupID:        coords.GroupID,
		ArtifactID:     coords.ArtifactID,
		Version:        coords.Version,
		Coordinates:    coords,
		Packaging:      PackagingJar,
		Properties: map[string]string{
			"project.build.sourceEncoding": "UTF-8",
			"maven.compiler.source":        "11",
			"maven.compiler.target":        "11",
		},
		Build: &Build{
			Plugins: []Plugin{
				{
					GroupID:    "org.apache.maven.plugins",
					ArtifactID: "maven-compiler-plugin",
					Version:    "3.11.0",
				},
			},
		},
	}
}

// createJavaLibrary creates a Java library template
func (tm *templateManager) createJavaLibrary(coords Coordinates) *Project {
	return &Project{
		XMLNS:          MavenXMLNamespace,
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: MavenXMLSchemaLocation,
		ModelVersion:   DefaultModelVersion,
		GroupID:        coords.GroupID,
		ArtifactID:     coords.ArtifactID,
		Version:        coords.Version,
		Coordinates:    coords,
		Packaging:      PackagingJar,
		Properties: map[string]string{
			"project.build.sourceEncoding": "UTF-8",
			"maven.compiler.source":        "11",
			"maven.compiler.target":        "11",
		},
		Dependencies: []Dependency{
			{
				GroupID:    "junit",
				ArtifactID: "junit",
				Version:    "4.13.2",
				Scope:      ScopeTest,
			},
		},
		Build: &Build{
			Plugins: []Plugin{
				{
					GroupID:    "org.apache.maven.plugins",
					ArtifactID: "maven-compiler-plugin",
					Version:    "3.11.0",
				},
				{
					GroupID:    "org.apache.maven.plugins",
					ArtifactID: "maven-jar-plugin",
					Version:    "3.3.0",
				},
			},
		},
	}
}

// createWebApp creates a web application template
func (tm *templateManager) createWebApp(coords Coordinates) *Project {
	return &Project{
		XMLNS:          MavenXMLNamespace,
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: MavenXMLSchemaLocation,
		ModelVersion:   DefaultModelVersion,
		GroupID:        coords.GroupID,
		ArtifactID:     coords.ArtifactID,
		Version:        coords.Version,
		Coordinates:    coords,
		Packaging:      PackagingWar,
		Properties: map[string]string{
			"project.build.sourceEncoding": "UTF-8",
			"maven.compiler.source":        "11",
			"maven.compiler.target":        "11",
		},
		Dependencies: []Dependency{
			{
				GroupID:    "javax.servlet",
				ArtifactID: "javax.servlet-api",
				Version:    "4.0.1",
				Scope:      ScopeProvided,
			},
			{
				GroupID:    "junit",
				ArtifactID: "junit",
				Version:    "4.13.2",
				Scope:      ScopeTest,
			},
		},
		Build: &Build{
			Plugins: []Plugin{
				{
					GroupID:    "org.apache.maven.plugins",
					ArtifactID: "maven-compiler-plugin",
					Version:    "3.11.0",
				},
				{
					GroupID:    "org.apache.maven.plugins",
					ArtifactID: "maven-war-plugin",
					Version:    "3.3.2",
				},
			},
		},
	}
}

// createJavaCard creates a JavaCard applet template
func (tm *templateManager) createJavaCard(coords Coordinates) *Project {
	return &Project{
		XMLNS:          MavenXMLNamespace,
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: MavenXMLSchemaLocation,
		ModelVersion:   DefaultModelVersion,
		GroupID:        coords.GroupID,
		ArtifactID:     coords.ArtifactID,
		Version:        coords.Version,
		Coordinates:    coords,
		Packaging:      "jar", // CAP packaging typically builds from JAR
		Properties: map[string]string{
			"project.build.sourceEncoding": "UTF-8",
			"maven.compiler.source":        "1.8",
			"maven.compiler.target":        "1.8",
			"javacard.version":             "3.0.5",
			"globalplatform.version":       "1.7.0",
		},
		Dependencies: []Dependency{
			{
				GroupID:    "com.github.martinpaljak",
				ArtifactID: "globalplatform",
				Version:    "1.7.0",
				Scope:      ScopeProvided,
			},
			{
				GroupID:    "com.github.martinpaljak",
				ArtifactID: "javacard-api",
				Version:    "3.0.5u3",
				Scope:      ScopeProvided,
			},
			{
				GroupID:    "junit",
				ArtifactID: "junit",
				Version:    "4.13.2",
				Scope:      ScopeTest,
			},
		},
		Build: &Build{
			Plugins: []Plugin{
				{
					GroupID:    "org.apache.maven.plugins",
					ArtifactID: "maven-compiler-plugin",
					Version:    "3.11.0",
				},
				{
					GroupID:    "com.github.martinpaljak",
					ArtifactID: "ant-javacard",
					Version:    "23.08.08",
					Executions: []PluginExecution{
						{
							ID:    "build-cap",
							Phase: PhasePackage,
							Goals: []string{"cap"},
						},
					},
				},
			},
		},
	}
}

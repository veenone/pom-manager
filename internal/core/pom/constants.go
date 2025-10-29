// Package pom provides Maven POM file parsing, generation, and validation.
//
// The package supports Maven POM schema 4.0.0 and provides functionality for:
//   - Parsing XML POM files into Go structures
//   - Generating valid POM XML from Go structures
//   - Validating POM files against Maven schema rules
//   - Organizing plugin executions by phase, goal, or plugin
//
// Example usage:
//
//     parser := pom.NewParser()
//     project, err := parser.ParseFile("pom.xml")
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Println(project.Coordinates.String())
package pom

// Maven POM model version
const (
	DefaultModelVersion = "4.0.0"
)

// Maven lifecycle phases in execution order
const (
	PhaseValidate            = "validate"
	PhaseInitialize          = "initialize"
	PhaseGenerateSources     = "generate-sources"
	PhaseProcessSources      = "process-sources"
	PhaseGenerateResources   = "generate-resources"
	PhaseProcessResources    = "process-resources"
	PhaseCompile             = "compile"
	PhaseProcessClasses      = "process-classes"
	PhaseGenerateTestSources = "generate-test-sources"
	PhaseProcessTestSources  = "process-test-sources"
	PhaseGenerateTestResources = "generate-test-resources"
	PhaseProcessTestResources = "process-test-resources"
	PhaseTestCompile         = "test-compile"
	PhaseProcessTestClasses  = "process-test-classes"
	PhaseTest                = "test"
	PhasePreparePackage      = "prepare-package"
	PhasePackage             = "package"
	PhasePreIntegrationTest  = "pre-integration-test"
	PhaseIntegrationTest     = "integration-test"
	PhasePostIntegrationTest = "post-integration-test"
	PhaseVerify              = "verify"
	PhaseInstall             = "install"
	PhaseDeploy              = "deploy"
)

// MavenLifecyclePhases returns all Maven lifecycle phases in execution order
var MavenLifecyclePhases = []string{
	PhaseValidate,
	PhaseInitialize,
	PhaseGenerateSources,
	PhaseProcessSources,
	PhaseGenerateResources,
	PhaseProcessResources,
	PhaseCompile,
	PhaseProcessClasses,
	PhaseGenerateTestSources,
	PhaseProcessTestSources,
	PhaseGenerateTestResources,
	PhaseProcessTestResources,
	PhaseTestCompile,
	PhaseProcessTestClasses,
	PhaseTest,
	PhasePreparePackage,
	PhasePackage,
	PhasePreIntegrationTest,
	PhaseIntegrationTest,
	PhasePostIntegrationTest,
	PhaseVerify,
	PhaseInstall,
	PhaseDeploy,
}

// Maven dependency scopes
const (
	ScopeCompile  = "compile"
	ScopeProvided = "provided"
	ScopeRuntime  = "runtime"
	ScopeTest     = "test"
	ScopeSystem   = "system"
	ScopeImport   = "import"
)

// ValidDependencyScopes contains all valid Maven dependency scopes
var ValidDependencyScopes = []string{
	ScopeCompile,
	ScopeProvided,
	ScopeRuntime,
	ScopeTest,
	ScopeSystem,
	ScopeImport,
}

// Maven packaging types
const (
	PackagingJar        = "jar"
	PackagingWar        = "war"
	PackagingEar        = "ear"
	PackagingPom        = "pom"
	PackagingMavenPlugin = "maven-plugin"
	PackagingRar        = "rar"
	PackagingPar        = "par"
)

// ValidPackagingTypes contains all valid Maven packaging types
var ValidPackagingTypes = []string{
	PackagingJar,
	PackagingWar,
	PackagingEar,
	PackagingPom,
	PackagingMavenPlugin,
	PackagingRar,
	PackagingPar,
}

// Maven XML namespace
const (
	MavenXMLNamespace = "http://maven.apache.org/POM/4.0.0"
	MavenXMLSchemaLocation = "http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd"
)

// File size limits
const (
	MaxFileSizeBytes = 10 * 1024 * 1024 // 10MB
)

// Default values
const (
	DefaultPackaging = PackagingJar
	DefaultScope     = ScopeCompile
)

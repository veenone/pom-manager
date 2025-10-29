package main

import (
	"fmt"
	"os"

	"github.com/user/pom-manager/internal/core/pom"
)

func main() {
	fmt.Println("=== Maven POM Manager MVP Test ===\n")

	// 1. Create a new project from template
	fmt.Println("1. Creating project from 'basic-java' template...")
	tm := pom.NewTemplateManager()

	coords := pom.Coordinates{
		GroupID:    "com.example",
		ArtifactID: "my-app",
		Version:    "1.0.0",
	}

	project, err := tm.Create("basic-java", coords)
	if err != nil {
		fmt.Printf("❌ Error creating project: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Created project: %s\n\n", project.Coordinates.String())

	// 2. List available templates
	fmt.Println("2. Available templates:")
	for _, tmpl := range tm.List() {
		fmt.Printf("   - %s: %s\n", tmpl.Name, tmpl.Description)
	}
	fmt.Println()

	// 3. Validate the project
	fmt.Println("3. Validating project...")
	validator := pom.NewValidator()
	result := validator.Validate(project)

	if result.Valid {
		fmt.Println("✓ Project is valid\n")
	} else {
		fmt.Println("❌ Project has validation errors:")
		for _, err := range result.Errors.AllErrors() {
			fmt.Printf("   - %s\n", err.Error())
		}
		os.Exit(1)
	}

	// 4. Generate XML
	fmt.Println("4. Generating POM XML...")
	generator := pom.NewGenerator()
	xmlBytes, err := generator.Generate(project)
	if err != nil {
		fmt.Printf("❌ Error generating XML: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Generated %d bytes of XML\n\n", len(xmlBytes))

	// 5. Write to file
	outputFile := "test-pom.xml"
	fmt.Printf("5. Writing to %s...\n", outputFile)
	err = generator.GenerateToFile(project, outputFile)
	if err != nil {
		fmt.Printf("❌ Error writing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Wrote POM to %s\n\n", outputFile)

	// 6. Parse it back
	fmt.Printf("6. Parsing %s back...\n", outputFile)
	parser := pom.NewParser()
	parsedProject, err := parser.ParseFile(outputFile)
	if err != nil {
		fmt.Printf("❌ Error parsing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Parsed project: %s\n\n", parsedProject.Coordinates.String())

	// 7. Add a dependency
	fmt.Println("7. Adding junit dependency...")
	parsedProject.Dependencies = append(parsedProject.Dependencies, pom.Dependency{
		GroupID:    "junit",
		ArtifactID: "junit",
		Version:    "4.13.2",
		Scope:      pom.ScopeTest,
	})
	fmt.Printf("✓ Added dependency, now has %d dependencies\n\n", len(parsedProject.Dependencies))

	// 8. Validate again
	fmt.Println("8. Validating modified project...")
	result = validator.Validate(parsedProject)
	if result.Valid {
		fmt.Println("✓ Modified project is valid\n")
	} else {
		fmt.Println("❌ Modified project has validation errors:")
		for _, err := range result.Errors.AllErrors() {
			fmt.Printf("   - %s\n", err.Error())
		}
	}

	// 9. Generate final XML
	fmt.Println("9. Generating final POM...")
	finalXML, err := generator.Generate(parsedProject)
	if err != nil {
		fmt.Printf("❌ Error generating final XML: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== Generated POM.XML ===")
	fmt.Println(string(finalXML))

	fmt.Println("\n=== MVP Test Complete ===")
	fmt.Println("✓ All operations successful!")
	fmt.Printf("\nCore capabilities demonstrated:\n")
	fmt.Println("  - Template-based project creation")
	fmt.Println("  - Project validation")
	fmt.Println("  - XML generation")
	fmt.Println("  - File I/O (write and read)")
	fmt.Println("  - POM parsing")
	fmt.Println("  - Dependency management")
	fmt.Println("\nMVP is ready! 🎉")
}

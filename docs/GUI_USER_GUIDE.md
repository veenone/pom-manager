# Maven POM Manager - GUI User Guide

## Table of Contents
1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [User Interface Overview](#user-interface-overview)
4. [Creating a New Project](#creating-a-new-project)
5. [Opening and Saving Projects](#opening-and-saving-projects)
6. [Editing Project Coordinates](#editing-project-coordinates)
7. [Managing Dependencies](#managing-dependencies)
8. [Configuring Build Plugins](#configuring-build-plugins)
9. [Managing Properties](#managing-properties)
10. [Working with Profiles](#working-with-profiles)
11. [Lifecycle Phase Management](#lifecycle-phase-management)
12. [XML Preview and Validation](#xml-preview-and-validation)
13. [Application Settings](#application-settings)
14. [Keyboard Shortcuts](#keyboard-shortcuts)
15. [Tips and Best Practices](#tips-and-best-practices)
16. [Troubleshooting](#troubleshooting)

---

## Introduction

The Maven POM Manager is a desktop application that provides a visual, user-friendly interface for creating and editing Maven Project Object Model (POM) files. Instead of manually editing XML, you can use intuitive forms and dialogs to manage all aspects of your Maven project configuration.

### What is a Maven POM?

A POM (Project Object Model) is the fundamental unit of work in Maven. It's an XML file (`pom.xml`) that contains information about the project and configuration details used by Maven to build the project. The POM includes:

- Project coordinates (Group ID, Artifact ID, Version)
- Dependencies on other libraries
- Build plugins and their configurations
- Project properties
- Build profiles for different environments

### Why Use This Application?

- **Visual Editing**: No need to remember XML syntax or Maven structure
- **Real-time Validation**: Instant feedback on errors and missing fields
- **Template-Based**: Quick start with pre-configured project templates
- **Safe Editing**: Validation prevents creating invalid POMs
- **Live Preview**: See the generated XML as you make changes
- **Intuitive Navigation**: Tree view and tabs organize the POM structure

---

## Getting Started

### Installation

1. **Download** the appropriate build for your operating system
2. **Extract** the archive to your preferred location
3. **Run** the executable:
   - Windows: `pom-manager-gui.exe`
   - Linux/macOS: `./pom-manager-gui`

### First Launch

When you first launch the application:

1. The main window opens with an empty project
2. The **File → New** menu option is available to create your first project
3. Settings are saved to `~/.pom-manager/gui-config.yaml`

### System Requirements

- **Operating System**: Windows 10+, Linux (recent distributions), macOS 10.13+
- **Memory**: 256 MB RAM minimum, 512 MB recommended
- **Disk Space**: 50 MB for application and configuration files

---

## User Interface Overview

The application window is divided into several key areas:

### 1. Menu Bar (Top)

- **File**: New, Open, Open Recent, Save, Save As, Exit
- **Edit**: Settings
- **Help**: Quick Help, Maven Basics, About

### 2. Tree Navigation Panel (Left, ~20%)

- Hierarchical view of your POM structure
- Click on nodes to navigate to corresponding tabs
- Sections:
  - **Project** (root): Project coordinates
  - **Properties**: Maven properties
  - **Dependencies**: Project dependencies
  - **Plugins**: Build plugins
  - **Profiles**: Build profiles (if any)

### 3. Editor Tabs (Center, ~45%)

Six tabs for editing different aspects of your POM:

- **Coordinates**: Project metadata
- **Dependencies**: Dependency management
- **Plugins**: Build plugin configuration
- **Properties**: Key-value properties
- **Profiles**: Build profile details
- **Lifecycle Phases**: Plugin execution phases

### 4. XML Preview Panel (Right, ~35%)

- **Validation Badge**: Shows validation status (green ✓ or red ✗)
- **XML Display**: Color-coded XML with proper indentation
- **Copy to Clipboard**: Button to copy the generated XML
- **Refresh**: Manual refresh (when live preview is off)

### 5. Status Bar (Bottom)

- Shows current file status
- Displays validation summary
- Indicates save state (modified/saved)

---

## Creating a New Project

### Using the Project Creation Wizard

1. **Open the Wizard**
   - Click **File → New** or press **Ctrl+N**
   - The wizard appears with Step 1

2. **Step 1: Project Coordinates**
   - **Group ID**: Your organization's reverse domain (e.g., `com.example`)
   - **Artifact ID**: Your project name (e.g., `my-app`)
   - **Version**: Initial version (default: `1.0.0`)
   - **Packaging**: Select JAR, WAR, POM, or maven-plugin
   - Click **Next**

3. **Step 2: Choose Template**
   - **Basic Java**: Standard JAR project with compiler plugin
   - **Java Library**: Library project with JUnit and JAR plugin
   - **Web App**: WAR-based web application with servlet dependencies
   - **JavaCard**: Smart card applet project with JavaCard APIs
   - Click **Finish**

4. **Result**
   - The project is created and loaded
   - All panels update with template defaults
   - XML preview shows the generated POM

### Template Details

#### Basic Java Template
- **Packaging**: JAR
- **Java Version**: 11 (source and target)
- **Plugins**: maven-compiler-plugin 3.11.0
- **Properties**: UTF-8 encoding, Java 11

#### Java Library Template
- **Packaging**: JAR
- **Java Version**: 11
- **Dependencies**: JUnit 4.13.2 (test scope)
- **Plugins**: maven-compiler-plugin, maven-jar-plugin
- **Use Case**: Creating reusable libraries

#### Web App Template
- **Packaging**: WAR
- **Java Version**: 11
- **Dependencies**: javax.servlet-api 4.0.1 (provided), JUnit (test)
- **Plugins**: maven-compiler-plugin, maven-war-plugin
- **Use Case**: Java web applications

#### JavaCard Template
- **Packaging**: JAR (builds to CAP)
- **Java Version**: 1.8 (JavaCard requirement)
- **Dependencies**:
  - GlobalPlatform API 1.7.0 (provided)
  - JavaCard API 3.0.5u3 (provided)
  - JUnit 4.13.2 (test)
- **Plugins**:
  - maven-compiler-plugin
  - ant-javacard 23.08.08 with CAP goal bound to package phase
- **Use Case**: Smart card applet development

---

## Opening and Saving Projects

### Opening Existing POMs

1. **From Menu**
   - Click **File → Open** or press **Ctrl+O**
   - Browse to your `pom.xml` file
   - Select and click **Open**

2. **From Recent Files**
   - Click **File → Open Recent**
   - Select from the last 10 opened files
   - Files are automatically filtered (deleted files don't appear)

3. **Drag and Drop**
   - Drag a `pom.xml` file onto the application window
   - The file opens automatically

### Saving Projects

1. **Save (Ctrl+S)**
   - Saves to the current file path
   - If new project, prompts for location (acts as Save As)

2. **Save As (Ctrl+Shift+S)**
   - Opens file dialog to choose new location
   - Creates a new file without modifying the original
   - Updates the current file path to the new location

### Auto-Save and Recovery

*Note: Auto-save feature is planned but not yet implemented (Task 20)*

---

## Editing Project Coordinates

The **Coordinates** tab contains the core project metadata.

### Required Fields (marked with *)

1. **Group ID***
   - Format: Reverse domain notation
   - Example: `com.example`, `org.apache.maven`
   - Validation: Required, no spaces

2. **Artifact ID***
   - Format: Lowercase with hyphens
   - Example: `my-app`, `spring-core`
   - Validation: Required, no spaces

3. **Version***
   - Format: Semantic versioning recommended
   - Example: `1.0.0`, `2.1.3-SNAPSHOT`
   - Validation: Required

### Optional Fields

4. **Packaging**
   - Values: `jar` (default), `war`, `pom`, `maven-plugin`
   - When to use:
     - `jar`: Standard applications and libraries
     - `war`: Web applications
     - `pom`: Parent POMs or aggregator projects
     - `maven-plugin`: Custom Maven plugins

5. **Name**
   - Human-readable project name
   - Example: "My Application"
   - Appears in build output and documentation

6. **Description**
   - Brief description of the project
   - Used in generated site documentation

### Real-time Validation

- Fields turn **green** with a ✓ when valid
- Fields turn **red** with an ✗ when invalid
- Validation errors appear in the validation badge and errors panel

---

## Managing Dependencies

The **Dependencies** tab lets you manage project dependencies.

### Viewing Dependencies

- Dependencies appear in a list
- Format: `groupId:artifactId:version [scope]`
- Example: `org.springframework:spring-core:5.3.30 [compile]`
- Click a dependency to select it
- Selection enables Edit and Remove buttons

### Adding a Dependency

1. Click **Add Dependency**
2. Fill in the dialog:
   - **Group ID**: Provider organization (e.g., `org.springframework`)
   - **Artifact ID**: Library name (e.g., `spring-core`)
   - **Version**: Version number (e.g., `5.3.30`)
   - **Scope**: When the dependency is needed
3. Click **OK**

### Dependency Scopes

- **compile** (default): Available in all phases
- **provided**: Required for compilation, provided by runtime (e.g., servlet-api)
- **runtime**: Not needed for compilation, only for execution
- **test**: Only for testing (e.g., JUnit)
- **system**: Like provided, but you specify the JAR path
- **import**: Used in dependencyManagement for POM imports

### Editing a Dependency

1. **Select** the dependency from the list
2. Click **Edit**
3. Modify fields in the dialog
4. Click **OK**

### Removing a Dependency

1. **Select** the dependency
2. Click **Remove**
3. Confirm if prompted
4. Dependency is removed immediately

### Exclusions

*Note: Exclusion management in the dialog is a planned enhancement*

### Common Dependencies

- **JUnit 4**: `junit:junit:4.13.2` (test)
- **JUnit 5**: `org.junit.jupiter:junit-jupiter:5.9.3` (test)
- **SLF4J API**: `org.slf4j:slf4j-api:2.0.7`
- **Logback**: `ch.qos.logback:logback-classic:1.4.7`
- **Spring Core**: `org.springframework:spring-core:5.3.30`
- **Jackson**: `com.fasterxml.jackson.core:jackson-databind:2.15.2`

---

## Configuring Build Plugins

The **Plugins** tab manages Maven build plugins.

### Viewing Plugins

- Plugins appear in a list
- Format: `groupId:artifactId:version`
- Example: `org.apache.maven.plugins:maven-compiler-plugin:3.11.0`

### Adding a Plugin

1. Click **Add Plugin**
2. Fill in the dialog:
   - **Group ID**: Plugin provider
   - **Artifact ID**: Plugin name
   - **Version**: Plugin version
3. Click **OK**

### Common Build Plugins

1. **maven-compiler-plugin** (`org.apache.maven.plugins:maven-compiler-plugin:3.11.0`)
   - Compiles Java source code
   - Configure Java version via properties

2. **maven-jar-plugin** (`org.apache.maven.plugins:maven-jar-plugin:3.3.0`)
   - Creates JAR archives
   - Configure manifest entries

3. **maven-war-plugin** (`org.apache.maven.plugins:maven-war-plugin:3.3.2`)
   - Creates WAR files for web applications

4. **maven-surefire-plugin** (`org.apache.maven.plugins:maven-surefire-plugin:3.1.2`)
   - Runs unit tests during the test phase

5. **maven-shade-plugin** (`org.apache.maven.plugins:maven-shade-plugin:3.5.0`)
   - Creates uber/fat JARs with dependencies

### Plugin Executions

Plugins can have multiple executions bound to different lifecycle phases. See the **Lifecycle Phases** tab for execution management.

---

## Managing Properties

The **Properties** tab manages Maven properties (key-value pairs).

### Viewing Properties

- Properties display as `name = value`
- Sorted alphabetically by property name
- Common properties appear by default in templates

### Adding a Property

1. Click **Add Property**
2. Enter:
   - **Property Name**: Key (e.g., `maven.compiler.source`)
   - **Value**: Value (e.g., `11`)
3. Click **Save**

### Editing a Property

1. **Select** the property
2. Click **Edit**
3. Modify the **Value** (name is read-only)
4. Click **Save**

### Removing a Property

1. **Select** the property
2. Click **Remove**
3. Property is deleted immediately

### Common Properties

1. **Java Version**
   ```
   maven.compiler.source = 11
   maven.compiler.target = 11
   ```

2. **Encoding**
   ```
   project.build.sourceEncoding = UTF-8
   ```

3. **Dependency Versions** (for version management)
   ```
   spring.version = 5.3.30
   junit.version = 4.13.2
   ```

4. **JavaCard Versions**
   ```
   javacard.version = 3.0.5
   globalplatform.version = 1.7.0
   ```

### Using Properties

Once defined, properties can be referenced in other parts of the POM using `${property.name}` syntax. However, the current GUI version doesn't support variable substitution display - this happens during Maven build.

---

## Working with Profiles

The **Profiles** tab displays Maven build profiles (view-only in current version).

### What are Profiles?

Profiles allow you to customize builds for different environments (development, staging, production) or different platforms.

### Viewing Profile Details

1. Navigate to the **Profiles** tab
2. **Left panel**: List of profile IDs
   - Active profiles shown with a checkmark ✓
3. **Right panel**: Detailed information
   - Activation conditions
   - Profile-specific properties
   - Profile-specific dependencies
   - Profile-specific plugins
   - Profile-specific modules

### Profile Activation

Profiles can be activated by:
- **activeByDefault**: `true/false`
- **JDK version**: e.g., `1.8`, `11`
- **Operating System**: name, family, arch, version
- **Property**: presence or value of a system property
- **File**: existence or absence of a file

### Example Profile Use Cases

1. **Development Profile**
   - Activates by default
   - Includes debug dependencies
   - Disables minification

2. **Production Profile**
   - Activated explicitly (`-Pproduction`)
   - Enables optimization plugins
   - Uses production database properties

3. **Platform Profiles**
   - Windows-specific plugin configurations
   - Linux-specific native libraries

*Note: Profile editing (add/remove/edit) is a planned enhancement*

---

## Lifecycle Phase Management

The **Lifecycle Phases** tab manages plugin execution bindings.

### Understanding Maven Lifecycle

Maven builds follow a predefined lifecycle with phases executed in order:

1. **validate**: Validate project structure
2. **compile**: Compile source code
3. **test**: Run unit tests
4. **package**: Package compiled code (JAR/WAR)
5. **verify**: Run integration tests
6. **install**: Install to local repository
7. **deploy**: Deploy to remote repository

Additional phases: initialize, generate-sources, process-sources, generate-resources, process-resources, process-classes, generate-test-sources, process-test-sources, generate-test-resources, process-test-resources, test-compile, process-test-classes, prepare-package, pre-integration-test, integration-test, post-integration-test

### Viewing Phase-Bound Executions

- Executions are organized by phase in an accordion
- Each phase shows the number of bound executions
- Expand a phase to see execution details:
  - Execution ID
  - Plugin goals
  - Configuration (if present)

### Adding an Execution

1. Ensure you have plugins added in the **Plugins** tab
2. Click **Add Execution**
3. Fill in the dialog:
   - **Plugin**: Select from existing plugins
   - **Execution ID**: Unique identifier (e.g., `build-cap`, `generate-sources`)
   - **Phase**: Maven lifecycle phase to bind to
   - **Goals**: Comma-separated plugin goals (e.g., `compile, testCompile`)
4. Click **Save**

### Example: JavaCard CAP Building

For a JavaCard project:
- **Plugin**: `com.github.martinpaljak:ant-javacard:23.08.08`
- **Execution ID**: `build-cap`
- **Phase**: `package`
- **Goals**: `cap`
- **Result**: CAP file is built during the package phase

### Removing an Execution

*Note: Execution removal is a planned enhancement*

---

## XML Preview and Validation

The **XML Preview** panel (right side) shows the generated POM XML.

### Validation Status Badge

The badge at the top shows:
- **Green ✓**: POM is valid
- **Red ✗ with count**: Number of validation errors

Click the badge to jump to the errors panel.

### XML Display Features

1. **Syntax Highlighting**
   - Tags: Green
   - Attribute names: Purple
   - Attribute values: Orange
   - XML declaration: Standard color

2. **Indentation**
   - 4-space indentation for readability
   - Proper nesting structure

3. **Auto-Update**
   - Updates within 100ms of changes
   - Configurable debounce delay in settings

### Copy to Clipboard

Click **Copy to Clipboard** to copy the entire XML to your clipboard for:
- Pasting into other editors
- Sharing with team members
- Manual review

### Refresh Button

- Enabled only when live preview is disabled
- Click to manually regenerate XML
- Useful for performance with very large POMs

### Validation Errors Panel

The **Errors** panel (bottom of center area) shows:
- Errors grouped by category:
  - Coordinates errors
  - Dependencies errors
  - Build errors
  - General errors
- Click an error to navigate to the relevant field
- Errors clear automatically when fixed

---

## Application Settings

Access settings via **Edit → Settings** or the settings button.

### General Tab

1. **Theme**
   - Light (default)
   - Dark
   - Changes apply immediately

2. **Auto-save Interval**
   - Minutes between auto-saves
   - 0 = disabled
   - *Feature planned but not yet active*

3. **Restore Session**
   - Checkbox: Reopen last file on startup
   - Restores window position

### Editor Tab

1. **Font Size**
   - Slider: 10-18 points
   - Default: 12
   - Affects all text in editor panels

2. **Live Preview**
   - Checkbox: Enable real-time XML updates
   - Disable for large POMs if performance is slow

3. **Validation Delay**
   - Milliseconds to wait before validating
   - Default: 100ms
   - Higher values reduce CPU usage during rapid typing

4. **Syntax Highlighting**
   - Checkbox: Enable XML syntax highlighting
   - Disable if colors are distracting

### Templates Tab

1. **Default Template**
   - Select which template to use by default
   - Options: basic-java, java-library, web-app, javacard

2. **Custom Template Directory**
   - Path to folder with custom templates
   - *Feature planned*

### Advanced Tab

1. **Maven Central Timeout**
   - Seconds to wait for Maven Central
   - *Used by planned search feature*

2. **Enable Debug Logging**
   - Checkbox: Write debug logs
   - Logs to `~/.pom-manager/debug.log`

3. **Cache Directory**
   - Location for cached data
   - Leave empty for default

### Buttons

- **OK**: Save settings and close
- **Cancel**: Discard changes
- **Reset to Defaults**: Restore default values

---

## Keyboard Shortcuts

### File Operations
- **Ctrl+N**: New project (wizard)
- **Ctrl+O**: Open file
- **Ctrl+S**: Save
- **Ctrl+Shift+S**: Save As
- **Ctrl+W** or **Ctrl+Q**: Close/Quit

### Application
- **F1**: Help
- **F5**: Refresh/Validate
- **Ctrl+,**: Settings (platform-specific)

### Navigation
- Click tree nodes to switch tabs
- Click on validation errors to jump to fields

---

## Tips and Best Practices

### Maven Best Practices

1. **Use Properties for Versions**
   - Define `<spring.version>5.3.30</spring.version>`
   - Reference as `${spring.version}` in dependencies
   - Makes version upgrades easier

2. **Scope Dependencies Correctly**
   - Use `provided` for server-provided libraries
   - Use `test` for testing frameworks
   - Keep `compile` dependencies minimal

3. **Semantic Versioning**
   - Format: `MAJOR.MINOR.PATCH`
   - Example: `1.2.3`
   - Use `-SNAPSHOT` for development versions

4. **Group ID Convention**
   - Use reverse domain: `com.company.project`
   - Matches Java package structure
   - Ensures uniqueness

### Application Tips

1. **Use Templates**
   - Start with a template matching your project type
   - Saves time setting up common plugins

2. **Watch the Validation Badge**
   - Green ✓ means your POM is valid
   - Fix errors (red ✗) before saving

3. **Use Recent Files**
   - Quick access to frequently edited POMs
   - Saves navigation time

4. **Learn Keyboard Shortcuts**
   - Faster workflow
   - Ctrl+S to save frequently

5. **Copy XML for Review**
   - Use "Copy to Clipboard" for code reviews
   - Paste into IDE for comparison

6. **Check the Tree View**
   - Quick overview of POM structure
   - See what's defined at a glance

### Performance Tips

1. **For Large POMs**
   - Disable live preview in settings
   - Use manual refresh (Refresh button)
   - Increase validation delay

2. **Many Dependencies**
   - Group related dependencies in your mental model
   - Use search (Ctrl+F) in XML preview if needed

---

## Troubleshooting

### Application Won't Start

**Symptoms**: Application crashes immediately or won't launch

**Solutions**:
1. Check that you have the required runtime libraries
2. On Windows, ensure Visual C++ Redistributable is installed
3. Check `~/.pom-manager/gui-config.yaml` isn't corrupted
   - Delete the file to reset to defaults
4. Check available disk space in home directory

### POM Won't Validate

**Symptoms**: Validation badge shows red ✗

**Solutions**:
1. Click the validation badge to see error details
2. Check required fields (Group ID, Artifact ID, Version)
3. Verify dependency coordinates are correct
4. Check for duplicate dependencies or plugins
5. Review errors panel for specific issues

### Changes Not Saving

**Symptoms**: File saves but changes don't persist

**Solutions**:
1. Check file permissions on the pom.xml
2. Verify the file path is correct (shown in status bar)
3. Try Save As to a new location
4. Check disk space

### XML Preview Not Updating

**Symptoms**: Preview doesn't show recent changes

**Solutions**:
1. Check if live preview is enabled (Settings → Editor)
2. Click the **Refresh** button manually
3. Increase validation delay if typing rapidly
4. Check for validation errors blocking generation

### Can't Add Dependency

**Symptoms**: Add Dependency button disabled or dialog doesn't work

**Solutions**:
1. Verify you've entered all required fields (Group ID, Artifact ID, Version)
2. Check for exact duplicate (same coordinates)
3. Try a different scope
4. Restart the application

### Application Runs Slowly

**Symptoms**: UI lags, slow response

**Solutions**:
1. Disable live preview for large POMs
2. Increase validation delay in settings
3. Close other resource-intensive applications
4. Check if POM has hundreds of dependencies
5. Reduce font size if using large fonts

### Lost Unsaved Changes

**Symptoms**: Application closed without saving

**Solutions**:
1. *Auto-save feature is planned (Task 20)*
2. Currently: Save frequently with Ctrl+S
3. Enable "Restore Session" to reopen last file

### Theme Issues

**Symptoms**: Colors hard to read, contrast problems

**Solutions**:
1. Switch theme (Light/Dark) in Settings
2. Adjust system display settings for brightness
3. Disable syntax highlighting if colors are distracting

### Recent Files Not Appearing

**Symptoms**: Open Recent menu is empty

**Solutions**:
1. Recent files are filtered - deleted files don't appear
2. Maximum 10 recent files stored
3. Check `~/.pom-manager/gui-config.yaml` for recent_files list

---

## Getting Help

### In-Application Help

- **F1**: Quick help dialog
- **Help → Maven Basics**: Introduction to Maven concepts
- **Help → About**: Version information

### External Resources

- **Maven Official Documentation**: https://maven.apache.org/guides/
- **Maven Central Repository**: https://search.maven.org/
- **POM Reference**: https://maven.apache.org/pom.html

### Reporting Issues

If you encounter bugs or have feature requests:

1. Check the troubleshooting section above
2. Verify you're using the latest version
3. Report issues on the project repository
4. Include:
   - Application version (Help → About)
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior

---

## Glossary

- **Artifact**: A file (JAR, WAR, etc.) produced by a Maven build
- **Artifact ID**: The name of the artifact (e.g., `my-app`)
- **CAP**: Converted Applet file format for JavaCard
- **Dependency**: An external library your project uses
- **Execution**: A specific invocation of a plugin goal
- **Goal**: An action that a Maven plugin can perform
- **Group ID**: Organization identifier (e.g., `com.example`)
- **Lifecycle**: The sequence of build phases Maven follows
- **Phase**: A step in the Maven build lifecycle (compile, test, package, etc.)
- **Plugin**: A tool that extends Maven's functionality
- **POM**: Project Object Model - Maven's project configuration file
- **Profile**: Alternative build configuration for different environments
- **Property**: A key-value pair for configuration (like a variable)
- **Scope**: When a dependency is needed (compile, runtime, test, etc.)
- **Version**: The version number of your project or a dependency

---

**Version**: 0.1.0-MVP
**Last Updated**: 2025-10-29
**Author**: Maven POM Manager Team

# Maven POM Manager

A modern, user-friendly GUI application for creating, editing, and managing Maven POM (Project Object Model) files. Built with Go and Fyne, this tool provides an intuitive interface for working with Maven projects without manually editing XML.

## Features

### Core Functionality
- **Template-Based Project Creation**: Quick start with pre-configured templates (Basic Java, Java Library, Web App, JavaCard)
- **Visual POM Editor**: Edit Maven coordinates, dependencies, plugins, properties, and profiles through an intuitive GUI
- **XML Syntax Highlighting**: Color-coded XML preview with proper indentation (4 spaces)
- **Real-Time Validation**: Instant feedback on POM structure and required fields
- **Dependency Management**: Add, edit, and remove dependencies with scope and exclusions support
- **Plugin Configuration**: Manage build plugins and their executions
- **Lifecycle Phase Management**: Add plugin executions and bind them to specific Maven lifecycle phases
- **Properties Editor**: Manage project properties with alphabetically sorted output
- **Profiles Support**: View and manage Maven build profiles with activation conditions
- **Tree Navigation**: Hierarchical view of POM structure for easy navigation

### GUI Features
- **Multi-Tab Interface**: Separate tabs for Coordinates, Dependencies, Plugins, Properties, Profiles, and Lifecycle Phases
- **Recent Files Menu**: Quick access to recently opened POM files (max 10)
- **Settings Dialog**: Customize theme, font size, auto-save, and more
- **XML Preview**: Live preview with syntax highlighting and copy-to-clipboard
- **Validation Badge**: Visual indicator of POM validation status
- **Button Tooltips**: All buttons have helpful tooltips (500ms hover delay)

### Recent Enhancements
- XML pretty-print with 4-space indentation
- Improved syntax highlighting (tags: green, attributes: purple, values: orange)
- Non-breaking spaces for proper indentation display
- Removed close button from settings modal for cleaner UI
- Properties alphabetically sorted in XML output
- Button tooltips for improved discoverability (500ms hover delay)
- Comprehensive GUI User Guide with tutorials and troubleshooting

### Documentation

📖 **[Complete GUI User Guide](docs/GUI_USER_GUIDE.md)** - Comprehensive guide covering:
- Getting started tutorials
- Detailed feature explanations
- Maven concepts for beginners
- Keyboard shortcuts
- Troubleshooting guide
- Best practices and tips

## Installation

### Prerequisites
- **Go 1.19+** (for building from source)
- **CGO Compiler**:
  - Windows: TDM-GCC or MinGW-w64
  - Linux: gcc
  - macOS: Xcode Command Line Tools
- **Git** (for cloning the repository)

### Building from Source

#### Windows (with TDM-GCC)
```bash
# Clone the repository
git clone <repository-url>
cd pom_manager

# Set GCC path (adjust path as needed)
export PATH="/c/TDM-GCC-64/bin:$PATH"

# Build the GUI application
make gui

# Run the application
./build/pom-manager-gui.exe
```

#### Linux/macOS
```bash
# Clone the repository
git clone <repository-url>
cd pom_manager

# Build the GUI application
make gui

# Run the application
./build/pom-manager-gui
```

## Usage

### Quick Start

For a comprehensive guide to all features, see the **[GUI User Guide](docs/GUI_USER_GUIDE.md)**.

### Starting the Application
```bash
# Using the built executable
./build/pom-manager-gui.exe   # Windows
./build/pom-manager-gui        # Linux/macOS

# Or using make
make run-gui
```

### Creating a New POM

1. **Launch the application**
2. **File → New** or use the wizard
3. **Select a template**:
   - **Basic Java**: Standard JAR project with compiler plugin
   - **Java Library**: Library project with JAR and compiler plugins
   - **Web App**: WAR-based web application
   - **JavaCard**: JavaCard applet project for smart cards (CAP packaging)
4. **Fill in coordinates**:
   - Group ID: `com.example`
   - Artifact ID: `my-app`
   - Version: `1.0.0`
5. **Click Create**

### Editing an Existing POM

1. **File → Open** or drag-and-drop a `pom.xml` file
2. **Navigate using the tree view** on the left
3. **Edit in the appropriate tab**:
   - **Coordinates**: Project metadata
   - **Dependencies**: Manage dependencies
   - **Plugins**: Configure build plugins
   - **Properties**: Project properties
   - **Profiles**: Build profiles
4. **Preview XML** in real-time
5. **File → Save** when done

### Adding Dependencies

1. Navigate to the **Dependencies** tab
2. Click **Add Dependency**
3. Enter:
   - Group ID: `org.springframework`
   - Artifact ID: `spring-core`
   - Version: `5.3.30`
   - Scope: `compile` (default)
4. Click **OK**
5. Preview updates automatically

### Managing Lifecycle Phases

1. Navigate to the **Lifecycle Phases** tab
2. Click **Add Execution** to bind a plugin goal to a lifecycle phase
3. Select:
   - **Plugin**: Choose from existing plugins in your project
   - **Execution ID**: Unique identifier for this execution
   - **Phase**: Maven lifecycle phase (compile, test, package, etc.)
   - **Goals**: Comma-separated list of plugin goals to execute
4. Click **Save**
5. View executions organized by phase in the accordion view

### Managing Settings

1. **Edit → Settings**
2. Configure:
   - **General**: Theme (light/dark), auto-save interval
   - **Editor**: Font size, live preview, syntax highlighting
   - **Templates**: Default template, custom template directory
   - **Advanced**: Maven Central timeout, debug logging, cache directory
3. Click **OK** to apply

## Project Structure

```
pom_manager/
├── cmd/
│   ├── cli/            # CLI application entry point
│   └── gui/            # GUI application entry point
├── internal/
│   ├── core/
│   │   └── pom/        # Core POM logic (model, parser, generator, validator)
│   ├── gui/
│   │   ├── dialogs/    # Dialog windows (settings, wizard)
│   │   ├── panels/     # Main UI panels (dependencies, plugins, etc.)
│   │   ├── widgets/    # Custom widgets (XML viewer, validation badge)
│   │   ├── windows/    # Main window
│   │   ├── state/      # Application state management
│   │   └── presenters/ # Presenter layer (MVC pattern)
│   └── templates/      # Built-in POM templates
├── build/              # Built executables
├── test_mvp.go         # MVP test script
├── Makefile            # Build automation
└── README.md           # This file
```

## Development

### Running Tests
```bash
# Run core tests
make test

# Run MVP integration test
go run test_mvp.go

# The test demonstrates:
# - Template-based project creation
# - Project validation
# - XML generation
# - File I/O operations
# - POM parsing
# - Dependency management
```

### Building

```bash
# Build GUI (requires CGO)
make gui

# Build CLI
make cli

# Build both
make build

# Clean build artifacts
make clean
```

### Code Organization

The project follows a clean architecture pattern:

- **Core Layer** (`internal/core/pom`): Domain logic, no external dependencies
  - `model.go`: Data structures
  - `parser.go`: XML → Model
  - `generator.go`: Model → XML
  - `validator.go`: Validation logic
  - `template.go`: Template management

- **GUI Layer** (`internal/gui`): Fyne-based user interface
  - **MVC Pattern**: Presenters mediate between views and model
  - **State Management**: Centralized application state
  - **Custom Widgets**: Reusable UI components

### Dependencies

**Core**:
- `github.com/beevik/etree`: XML processing
- `gopkg.in/yaml.v3`: Settings file format

**GUI**:
- `fyne.io/fyne/v2`: Cross-platform GUI framework
- Requires CGO for native rendering

## Configuration

Settings are stored in `~/.pom-manager/gui-config.yaml`:

```yaml
# General
theme: light
auto_save_interval: 5
restore_session: true

# Editor
font_size: 12
live_preview: true
validation_delay: 100
syntax_highlight: true

# Templates
default_template: basic-java
custom_template_dir: ""

# Advanced
maven_central_timeout: 10
enable_debug_log: false
cache_dir: ""

# Window
window_width: 1024
window_height: 768
window_x: 0
window_y: 0

# Recent files
recent_files:
  - /path/to/project1/pom.xml
  - /path/to/project2/pom.xml
```

## Troubleshooting

### Build Issues

**CGO Compiler Not Found**:
```bash
# Windows: Install TDM-GCC
# Download from: https://jmeubank.github.io/tdm-gcc/

# Set PATH
export PATH="/c/TDM-GCC-64/bin:$PATH"
```

**Fyne Dependencies Missing**:
```bash
# Linux: Install required packages
sudo apt-get install gcc libgl1-mesa-dev xorg-dev

# macOS: Install Xcode Command Line Tools
xcode-select --install
```

### Runtime Issues

**Application Won't Start**:
- Ensure CGO runtime libraries are available
- Check settings file is not corrupted: `~/.pom-manager/gui-config.yaml`
- Run with debug logging enabled in settings

**XML Validation Errors**:
- Ensure required fields are filled (groupId, artifactId, version)
- Check XML structure in Preview tab
- Validation errors show in status badge (click for details)

## Contributing

Contributions are welcome! Please ensure:

1. Code follows Go conventions (`gofmt`, `golint`)
2. Tests pass (`make test`)
3. GUI builds successfully (`make gui`)
4. MVP test succeeds (`go run test_mvp.go`)

## License

MIT License

Copyright (c) 2025 Achmad Fienan Rahardianto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Authors

Achmad Fienan Rahardianto (veenone@gmail.com)

## Acknowledgments

- Built with [Fyne](https://fyne.io/) GUI toolkit
- XML processing by [etree](https://github.com/beevik/etree)
- Inspired by Maven's POM structure

---

**Status**: MVP Complete ✓

All core features implemented and tested. Ready for production use!

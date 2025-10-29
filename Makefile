# Maven POM Manager - Makefile
# Build system for CLI and GUI applications

# Variables
BINARY_NAME_GUI=pom-manager-gui
BINARY_NAME_CLI=pom-manager-cli
VERSION=1.0.0
BUILD_DIR=build
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w -X main.AppVersion=$(VERSION)"

# Detect OS
ifeq ($(OS),Windows_NT)
	BINARY_EXT=.exe
	RM=del /Q
	RMDIR=rmdir /S /Q
	MKDIR=mkdir
else
	BINARY_EXT=
	RM=rm -f
	RMDIR=rm -rf
	MKDIR=mkdir -p
endif

# Default target
.PHONY: all
all: help

# Help target
.PHONY: help
help:
	@echo "Maven POM Manager - Build System"
	@echo "================================="
	@echo ""
	@echo "Available targets:"
	@echo "  make cli          - Build CLI application (no CGO required)"
	@echo "  make gui          - Build GUI application (requires CGO + MinGW/TDM-GCC)"
	@echo "  make build        - Build both CLI and GUI"
	@echo "  make test         - Run all tests"
	@echo "  make test-cli     - Run CLI tests only"
	@echo "  make test-gui     - Run GUI tests only"
	@echo "  make test-core    - Run core engine tests"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make fmt          - Format code"
	@echo "  make vet          - Run go vet"
	@echo "  make install-deps - Install/update dependencies"
	@echo "  make run-cli      - Build and run CLI"
	@echo "  make run-gui      - Build and run GUI"
	@echo ""
	@echo "GUI Build Requirements:"
	@echo "  - CGO must be enabled (set CGO_ENABLED=1)"
	@echo "  - Windows: Install TDM-GCC or MinGW-w64"
	@echo "  - Add gcc to PATH"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make cli          - Try the CLI (works without CGO)"
	@echo "  2. make test         - Run all tests"
	@echo "  3. make gui          - Build GUI (after installing GCC)"
	@echo ""

# Build CLI application (no CGO required)
.PHONY: cli
cli:
	@echo "Building CLI application..."
	@$(MKDIR) $(BUILD_DIR) 2>nul || echo Directory exists
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME_CLI)$(BINARY_EXT) ./cmd/cli
	@echo "CLI built successfully: $(BUILD_DIR)/$(BINARY_NAME_CLI)$(BINARY_EXT)"

# Build GUI application (requires CGO + GCC)
.PHONY: gui
gui:
	@echo "Building GUI application..."
	@echo "Note: This requires CGO and a C compiler (TDM-GCC or MinGW-w64)"
	@$(MKDIR) $(BUILD_DIR) 2>nul || echo Directory exists
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME_GUI)$(BINARY_EXT) ./cmd/gui
	@echo "GUI built successfully: $(BUILD_DIR)/$(BINARY_NAME_GUI)$(BINARY_EXT)"

# Build both applications
.PHONY: build
build: cli gui

# Run CLI application
.PHONY: run-cli
run-cli: cli
	@echo "Running CLI application..."
	@$(BUILD_DIR)/$(BINARY_NAME_CLI)$(BINARY_EXT)

# Run GUI application
.PHONY: run-gui
run-gui: gui
	@echo "Running GUI application..."
	@$(BUILD_DIR)/$(BINARY_NAME_GUI)$(BINARY_EXT)

# Run all tests
.PHONY: test
test:
	@echo "Running all tests..."
	$(GO) test ./... -v -cover

# Run CLI tests only
.PHONY: test-cli
test-cli:
	@echo "Running CLI tests..."
	$(GO) test ./internal/cli/... -v

# Run GUI tests only
.PHONY: test-gui
test-gui:
	@echo "Running GUI tests..."
	$(GO) test ./internal/gui/... -v

# Run core engine tests
.PHONY: test-core
test-core:
	@echo "Running core engine tests..."
	$(GO) test ./internal/core/... -v

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Install/update dependencies
.PHONY: install-deps
install-deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies installed successfully"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@if exist $(BUILD_DIR) $(RMDIR) $(BUILD_DIR) 2>nul
	@if exist coverage.out $(RM) coverage.out 2>nul
	@if exist coverage.html $(RM) coverage.html 2>nul
	@if exist test-pom.xml $(RM) test-pom.xml 2>nul
	@echo "Clean complete"

# Check if GCC is available
.PHONY: check-gcc
check-gcc:
	@echo "Checking for GCC..."
	@gcc --version || echo "ERROR: GCC not found. Install TDM-GCC or MinGW-w64 for GUI builds."

# Verify environment for GUI build
.PHONY: check-env
check-env: check-gcc
	@echo "Checking Go environment..."
	@$(GO) version
	@echo ""
	@echo "CGO Status:"
	@$(GO) env CGO_ENABLED
	@echo ""
	@echo "Environment check complete"

# Development build (faster, with debug info)
.PHONY: dev-cli
dev-cli:
	@echo "Building CLI (development mode)..."
	@$(MKDIR) $(BUILD_DIR) 2>nul || echo Directory exists
	CGO_ENABLED=0 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME_CLI)$(BINARY_EXT) ./cmd/cli

.PHONY: dev-gui
dev-gui:
	@echo "Building GUI (development mode)..."
	@$(MKDIR) $(BUILD_DIR) 2>nul || echo Directory exists
	CGO_ENABLED=1 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME_GUI)$(BINARY_EXT) ./cmd/gui

# Quick test for MVP functionality
.PHONY: test-mvp
test-mvp:
	@echo "Running MVP functionality test..."
	$(GO) run test_mvp.go

# Install the application (copy to GOPATH/bin or specified location)
.PHONY: install
install: build
	@echo "Installing binaries..."
	@copy $(BUILD_DIR)\$(BINARY_NAME_CLI)$(BINARY_EXT) $(GOPATH)\bin\ 2>nul || echo "Install CLI manually"
	@copy $(BUILD_DIR)\$(BINARY_NAME_GUI)$(BINARY_EXT) $(GOPATH)\bin\ 2>nul || echo "Install GUI manually"
	@echo "Installation complete (if GOPATH is set)"

# Show build information
.PHONY: info
info:
	@echo "Build Information"
	@echo "================="
	@echo "Version:     $(VERSION)"
	@echo "CLI Binary:  $(BUILD_DIR)/$(BINARY_NAME_CLI)$(BINARY_EXT)"
	@echo "GUI Binary:  $(BUILD_DIR)/$(BINARY_NAME_GUI)$(BINARY_EXT)"
	@echo "Go Version:  " && $(GO) version
	@echo "Build Dir:   $(BUILD_DIR)"

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	@$(GO) doc ./...

# Benchmark tests
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

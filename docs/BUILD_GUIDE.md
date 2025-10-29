# Maven POM Manager - Build Guide

## Quick Start

### Option 1: Using Make (Recommended)

```bash
# View all available commands
make help

# Build CLI (works without CGO)
make cli

# Run tests
make test

# Build GUI (requires GCC)
make gui
```

### Option 2: Direct Go Commands

```bash
# Build CLI
go build -o build/pom-manager-cli.exe ./cmd/cli

# Build GUI (requires CGO)
CGO_ENABLED=1 go build -o build/pom-manager-gui.exe ./cmd/gui

# Run tests
go test ./...
```

## Prerequisites

### For CLI Development
- **Go 1.21+** - [Download](https://go.dev/dl/)
- No additional requirements

### For GUI Development
- **Go 1.21+** - [Download](https://go.dev/dl/)
- **C Compiler** (one of):
  - **TDM-GCC** - [Download](https://jmeubank.github.io/tdm-gcc/) (Recommended for Windows)
  - **MinGW-w64** - [Download](https://www.mingw-w64.org/)
  - **MSYS2** with gcc - [Download](https://www.msys2.org/)

### Installing TDM-GCC (Windows)

1. Download TDM-GCC from https://jmeubank.github.io/tdm-gcc/
2. Run the installer
3. Choose "Add to PATH" during installation
4. Verify installation:
   ```bash
   gcc --version
   ```

## Project Structure

```
pom_manager/
├── cmd/
│   ├── cli/          # CLI entry point
│   └── gui/          # GUI entry point (requires CGO)
├── internal/
│   ├── core/         # Core POM engine
│   │   └── pom/      # Parser, Generator, Validator, Templates
│   ├── cli/          # CLI implementation
│   └── gui/          # GUI implementation
│       ├── dialogs/  # Dialog windows
│       ├── panels/   # UI panels
│       ├── presenters/ # Business logic layer
│       ├── state/    # Application state management
│       └── windows/  # Main window
├── build/            # Compiled binaries (created by make)
├── Makefile          # Build system
└── test_mvp.go       # MVP functionality test

```

## Build Commands Reference

### Building

| Command | Description | Requirements |
|---------|-------------|--------------|
| `make cli` | Build CLI application | Go only |
| `make gui` | Build GUI application | Go + GCC |
| `make build` | Build both CLI and GUI | Go + GCC |
| `make dev-cli` | Fast CLI build (dev mode) | Go only |
| `make dev-gui` | Fast GUI build (dev mode) | Go + GCC |

### Testing

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-cli` | Run CLI tests only |
| `make test-gui` | Run GUI tests only |
| `make test-core` | Run core engine tests |
| `make test-coverage` | Generate coverage report |
| `make test-mvp` | Run MVP functionality test |

### Code Quality

| Command | Description |
|---------|-------------|
| `make fmt` | Format code with gofmt |
| `make vet` | Run go vet |
| `make lint` | Run golangci-lint (if installed) |

### Maintenance

| Command | Description |
|---------|-------------|
| `make clean` | Remove build artifacts |
| `make install-deps` | Install/update dependencies |
| `make check-env` | Verify build environment |
| `make info` | Show build information |

### Running

| Command | Description |
|---------|-------------|
| `make run-cli` | Build and run CLI |
| `make run-gui` | Build and run GUI |

## Common Issues & Solutions

### Issue: "gcc: command not found"
**Solution:** Install TDM-GCC or MinGW-w64 and add to PATH

### Issue: "CGO is disabled"
**Solution:** Set environment variable:
```bash
# Windows (CMD)
set CGO_ENABLED=1

# Windows (PowerShell)
$env:CGO_ENABLED = "1"

# Or use make commands (handles this automatically)
make gui
```

### Issue: GUI build takes a long time
**Solution:** This is normal for Fyne first build. Subsequent builds are faster.
Use `make dev-gui` for faster development builds without optimizations.

### Issue: "undefined reference" errors during GUI build
**Solution:** Ensure you have a complete GCC installation including mingw32-make

## Test Results

### Current Test Status
✅ **All tests passing** (17 tests total)

```
Core Engine Tests:
  ✅ Parser tests
  ✅ Generator tests
  ✅ Validator tests
  ✅ Template tests

CLI Tests:
  ✅ Command execution tests
  ✅ Argument parsing tests

GUI Tests:
  ✅ State management (7 tests)
  ✅ Settings validation (10 tests)
  ✅ Presenter logic (7 tests)
```

Run tests with:
```bash
make test
```

## GUI Build Verification

To verify your environment is ready for GUI builds:

```bash
# Check if GCC is installed
make check-gcc

# Check complete environment
make check-env

# Try building CLI first (no CGO required)
make cli

# If CLI works, try GUI
make gui
```

## Development Workflow

### For CLI Development (No CGO Required)
```bash
# 1. Make changes to CLI code
# 2. Build
make cli

# 3. Test
make test-cli

# 4. Run
./build/pom-manager-cli.exe --help
```

### For GUI Development (Requires GCC)
```bash
# 1. Ensure GCC is installed
make check-gcc

# 2. Make changes to GUI code
# 3. Build (dev mode for faster iteration)
make dev-gui

# 4. Test
make test-gui

# 5. Run
./build/pom-manager-gui.exe
```

### Running Tests During Development
```bash
# Quick test (specific package)
go test ./internal/gui/state -v

# Full test suite
make test

# With coverage
make test-coverage
# Open coverage.html in browser
```

## Binary Locations

After building, binaries are located in:
- **CLI**: `build/pom-manager-cli.exe`
- **GUI**: `build/pom-manager-gui.exe`

## Performance Notes

- **CLI**: ~8MB binary, instant startup
- **GUI**: ~20MB binary, < 2 second startup
- **First Build**: GUI may take 2-5 minutes (Fyne compilation)
- **Subsequent Builds**: GUI takes 30-60 seconds

## Next Steps

1. ✅ Build CLI: `make cli`
2. ✅ Run tests: `make test`
3. ✅ Install GCC (for GUI)
4. ✅ Build GUI: `make gui`
5. ✅ Run GUI: `./build/pom-manager-gui.exe`

## Support

If you encounter issues:
1. Check `make check-env` output
2. Verify GCC installation: `gcc --version`
3. Try CLI build first to isolate CGO issues
4. Review error messages for missing dependencies

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Fyne Documentation](https://docs.fyne.io/)
- [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
- [Project README](README.md)

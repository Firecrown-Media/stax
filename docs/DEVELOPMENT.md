# Stax Development Guide

Guide for contributing to Stax development and adding new features.

## Table of Contents
- [Architecture Overview](#architecture-overview)
- [Development Setup](#development-setup)
- [Code Structure](#code-structure)
- [Adding Features](#adding-features)
- [Testing](#testing)
- [Build and Release](#build-and-release)
- [Contributing Guidelines](#contributing-guidelines)

## Architecture Overview

### How Stax Works

```
┌─────────────┐
│   CLI User  │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────┐
│     Stax CLI (Cobra)        │
│  ┌───────┬──────┬────────┐  │
│  │  cmd/ │ pkg/ │ config │  │
│  └───────┴──────┴────────┘  │
└──────────┬──────────────────┘
           │
           ▼
┌──────────────────┐      ┌──────────────────┐
│      DDEV        │      │    WP Engine     │
│   (Container     │      │      API         │
│    Management)   │      │   (Remote Sync)  │
└──────────────────┘      └──────────────────┘
           │
           ▼
┌──────────────────┐
│     Docker       │
│   (Containers)   │
└──────────────────┘
```

### Technology Stack

- **Language:** Go 1.19+
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra)
- **Configuration:** [Viper](https://github.com/spf13/viper)
- **Container Management:** DDEV
- **Build System:** Make
- **Testing:** Go testing package

## Development Setup

### Prerequisites

1. **Install Go:**
```bash
# Mac
brew install go

# Linux
sudo apt install golang-go

# Verify
go version  # Should be 1.19+
```

2. **Clone Repository:**
```bash
git clone https://github.com/Firecrown-Media/stax.git
cd stax
```

3. **Install Dependencies:**
```bash
# Go dependencies
go mod download
go mod tidy

# Development tools
brew install golangci-lint  # Linter
```

### Development Workflow

1. **Create feature branch:**
```bash
git checkout -b feature/my-feature
```

2. **Make changes and test:**
```bash
# Build locally
make build

# Test your changes
./build/stax --version
./build/stax [command]

# Run tests
make test

# Run linter
make lint
```

3. **Install for testing:**
```bash
# Quick install to Go bin
make update

# Test installed version
stax --version
```

## Code Structure

### Directory Layout

```
stax/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and flags
│   ├── init.go            # stax init command
│   ├── start.go           # stax start command
│   ├── wpe.go             # WP Engine commands
│   ├── swap.go            # Hot swap functionality
│   └── swap_presets.go    # Swap preset definitions
├── pkg/                    # Reusable packages
│   ├── config/            # Configuration management
│   ├── ddev/              # DDEV integration
│   ├── wordpress/         # WordPress operations
│   ├── wpengine/          # WP Engine client
│   ├── ui/                # User interface (spinner, etc.)
│   └── errors/            # Custom error types
├── main.go                # Entry point
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
└── Makefile              # Build automation
```

### Package Descriptions

#### cmd/ - Commands
Each file implements a Cobra command:
```go
// cmd/init.go
var initCmd = &cobra.Command{
    Use:   "init [project-name]",
    Short: "Initialize a new project",
    RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
    // Command implementation
}
```

#### pkg/config - Configuration
Handles global and project configurations:
```go
// pkg/config/config.go
type Config struct {
    PHPVersion    string `yaml:"php_version"`
    MySQLVersion  string `yaml:"mysql_version"`
    WordPress     WordPressConfig
}
```

#### pkg/ddev - DDEV Integration
Wraps DDEV CLI commands:
```go
// pkg/ddev/ddev.go
func Start(projectPath string) error {
    cmd := exec.Command("ddev", "start")
    cmd.Dir = projectPath
    return cmd.Run()
}
```

#### pkg/wpengine - WP Engine Client
Handles WP Engine API and SSH operations:
```go
// pkg/wpengine/client.go
type Client struct {
    Username string
    Password string
    SSHKey   string
}

func (c *Client) Sync(install string, opts SyncOptions) error {
    // Sync implementation
}
```

## Adding Features

### Adding a New Command

1. **Create command file:**
```go
// cmd/mycommand.go
package cmd

import (
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand [args]",
    Short: "Short description",
    Long:  `Longer description...`,
    RunE:  runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCmd)

    // Add flags
    myCmd.Flags().StringP("option", "o", "", "Option description")
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
```

2. **Add to root command:**
```go
// cmd/root.go
func init() {
    // ... existing code
    rootCmd.AddCommand(myCmd)
}
```

### Adding a Subcommand

```go
// cmd/parent.go
var parentCmd = &cobra.Command{
    Use:   "parent",
    Short: "Parent command",
}

var subCmd = &cobra.Command{
    Use:   "sub",
    Short: "Subcommand",
    RunE:  runSub,
}

func init() {
    rootCmd.AddCommand(parentCmd)
    parentCmd.AddCommand(subCmd)
}
```

### Adding a New Package

1. **Create package directory:**
```bash
mkdir pkg/mypackage
```

2. **Implement package:**
```go
// pkg/mypackage/mypackage.go
package mypackage

type Client struct {
    // Fields
}

func New(opts ...Option) *Client {
    // Constructor
}

func (c *Client) DoSomething() error {
    // Method implementation
}
```

3. **Write tests:**
```go
// pkg/mypackage/mypackage_test.go
package mypackage

import "testing"

func TestDoSomething(t *testing.T) {
    client := New()
    err := client.DoSomething()
    if err != nil {
        t.Errorf("Expected nil, got %v", err)
    }
}
```

### Adding a New Hosting Provider

Example: Adding Pantheon support

1. **Create provider package:**
```go
// pkg/pantheon/client.go
package pantheon

type Client struct {
    Email    string
    Password string
}

func New(email, password string) *Client {
    return &Client{
        Email:    email,
        Password: password,
    }
}

func (c *Client) ListSites() ([]Site, error) {
    // Implementation
}

func (c *Client) Sync(siteName string, opts SyncOptions) error {
    // Implementation
}
```

2. **Add CLI commands:**
```go
// cmd/pantheon.go
package cmd

var pantheonCmd = &cobra.Command{
    Use:   "pantheon",
    Short: "Pantheon hosting commands",
}

var pantheonSyncCmd = &cobra.Command{
    Use:   "sync [site]",
    Short: "Sync from Pantheon",
    RunE:  runPantheonSync,
}

func init() {
    rootCmd.AddCommand(pantheonCmd)
    pantheonCmd.AddCommand(pantheonSyncCmd)
}
```

### Adding Configuration Options

1. **Update config struct:**
```go
// pkg/config/config.go
type Config struct {
    // Existing fields...
    MyNewOption string `yaml:"my_new_option"`
}
```

2. **Add to Viper:**
```go
// cmd/root.go
func initConfig() {
    // Existing code...
    viper.SetDefault("my_new_option", "default_value")
}
```

## Testing

### Unit Tests

```go
// pkg/wordpress/wordpress_test.go
package wordpress

import "testing"

func TestHasWordPress(t *testing.T) {
    tests := []struct {
        name     string
        path     string
        expected bool
    }{
        {"WordPress exists", "/path/to/wp", true},
        {"No WordPress", "/empty/path", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := HasWordPress(tt.path)
            if result != tt.expected {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

### Integration Tests

```go
// test/integration/ddev_test.go
// +build integration

package integration

import (
    "testing"
    "github.com/Firecrown-Media/stax/pkg/ddev"
)

func TestDDEVIntegration(t *testing.T) {
    if !ddev.IsInstalled() {
        t.Skip("DDEV not installed")
    }

    // Test actual DDEV operations
    err := ddev.Init("/tmp/test-project", ddev.Config{
        ProjectName: "test",
    })

    if err != nil {
        t.Fatalf("Failed to init: %v", err)
    }

    // Cleanup
    defer ddev.Delete("/tmp/test-project", "test", true, true)
}
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./pkg/wordpress/...

# Run with verbose output
go test -v ./...

# Run integration tests
go test -tags=integration ./test/integration/...
```

## Build and Release

### Local Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Output in build/ directory:
# - stax-darwin-amd64 (Intel Mac)
# - stax-darwin-arm64 (Apple Silicon)
# - stax-linux-amd64
```

### Version Management

Version is set in the build:
```go
// cmd/version.go
var Version = "dev"  // Set by build

// Makefile
LDFLAGS=-ldflags "-X github.com/Firecrown-Media/stax/cmd.Version=$(VERSION)"
```

### Release Process

1. **Update version:**
```bash
git tag v1.2.3
git push origin v1.2.3
```

2. **GitHub Actions automatically:**
   - Builds binaries for all platforms
   - Creates GitHub release
   - Updates Homebrew tap

3. **Manual release:**
```bash
# Build release artifacts
make release

# Creates in build/release/:
# - stax-darwin-amd64.tar.gz
# - stax-darwin-arm64.tar.gz
# - stax-linux-amd64.tar.gz
```

## Contributing Guidelines

### Code Style

1. **Follow Go conventions:**
```bash
# Format code
go fmt ./...
make fmt

# Vet code
go vet ./...
make vet

# Lint
golangci-lint run
make lint
```

2. **Naming conventions:**
- Exported functions: `CapitalCase`
- Unexported functions: `camelCase`
- Constants: `UPPER_SNAKE_CASE`
- Packages: lowercase

### Commit Messages

Follow conventional commits:
```
feat: Add Pantheon integration
fix: Resolve SSH connection timeout
docs: Update README with examples
chore: Update dependencies
test: Add wpengine client tests
```

### Pull Request Process

1. **Before submitting:**
```bash
# Ensure tests pass
make test

# Ensure linting passes
make lint

# Ensure it builds
make build
```

2. **PR Description should include:**
- What changes were made
- Why changes were needed
- How to test changes
- Breaking changes (if any)

### Error Handling

```go
// Good - Return errors for handling
func DoSomething() error {
    if err := operation(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}

// Bad - Panic on error
func DoSomething() {
    if err := operation(); err != nil {
        panic(err)  // Don't do this
    }
}
```

### Logging

```go
// Use structured logging
import "github.com/Firecrown-Media/stax/pkg/logger"

logger.Info("Starting operation",
    "project", projectName,
    "path", projectPath)

// Verbose mode
if verbose {
    logger.Debug("Detailed information...")
}
```

## Development Tips

### Debugging

1. **Add debug output:**
```go
if verbose {
    fmt.Fprintf(os.Stderr, "Debug: %v\n", variable)
}
```

2. **Use delve debugger:**
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug main.go -- init test-project
```

### Performance

1. **Profile code:**
```go
import "runtime/pprof"

f, _ := os.Create("cpu.prof")
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()
```

2. **Benchmark:**
```go
func BenchmarkOperation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Operation()
    }
}
```

## Resources

- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [DDEV Documentation](https://ddev.readthedocs.io/)

## Getting Help

- Review existing code for patterns
- Check [GitHub Issues](https://github.com/Firecrown-Media/stax/issues)
- Contact maintainers at dev@firecrown.com
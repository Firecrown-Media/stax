# Testing Guide

This document describes the testing strategy, infrastructure, and best practices for the Stax CLI tool.

## Table of Contents

- [Overview](#overview)
- [Test Organization](#test-organization)
- [Running Tests](#running-tests)
- [Writing Tests](#writing-tests)
- [Test Coverage](#test-coverage)
- [Continuous Integration](#continuous-integration)
- [Debugging Tests](#debugging-tests)

## Overview

The Stax project maintains a comprehensive test suite to ensure reliability and correctness. We aim for >70% code coverage across the codebase.

### Test Types

1. **Unit Tests**: Test individual functions and packages in isolation
2. **Integration Tests**: Test interaction between multiple components
3. **End-to-End Tests**: Test complete user workflows
4. **Security Tests**: Test security-critical functionality

## Test Organization

```
stax/
├── pkg/                        # Package code with unit tests
│   ├── config/
│   │   ├── config.go
│   │   ├── config_test.go     # Unit tests for config
│   │   ├── loader.go
│   │   └── loader_test.go     # Unit tests for loader
│   ├── ddev/
│   │   ├── manager.go
│   │   └── manager_test.go    # Unit tests for DDEV manager
│   └── ...
├── test/                      # Integration and E2E tests
│   ├── helpers/               # Test helper functions
│   │   └── helpers.go
│   ├── fixtures/              # Test data and configurations
│   │   ├── config.yml
│   │   └── database.sql
│   ├── mocks/                 # Mock implementations
│   │   ├── wpengine.go
│   │   ├── ssh.go
│   │   └── ddev.go
│   ├── integration/           # Integration tests
│   │   ├── init_test.go
│   │   └── db_test.go
│   └── e2e/                   # End-to-end tests
│       └── e2e_test.go
└── pkg/testutil/              # Test utilities
    └── testutil.go
```

## Running Tests

### Quick Start

```bash
# Run all unit tests (fast)
make test

# Run all tests including integration and security
make test-all

# Generate coverage report
make test-coverage
```

### Specific Test Types

```bash
# Unit tests only
make test-unit

# Integration tests
make test-integration

# End-to-end tests
make test-e2e

# Security tests
make test-security

# Verbose output
make test-verbose
```

### Test Tags

Integration and E2E tests use build tags to avoid running by default:

```bash
# Run integration tests manually
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/...

# Run E2E tests manually
RUN_E2E_TESTS=true go test -v -tags=e2e ./test/e2e/...
```

### Individual Package Tests

```bash
# Test specific package
go test -v ./pkg/config/

# Test with coverage
go test -v -cover ./pkg/config/

# Test with race detection
go test -v -race ./pkg/config/
```

## Writing Tests

### Unit Test Structure

Unit tests should be placed in the same package as the code being tested, with a `_test.go` suffix.

```go
package config

import (
    "testing"
    "github.com/firecrown-media/stax/pkg/testutil"
)

func TestConfigLoad(t *testing.T) {
    // Create test directory
    dir := testutil.TempDir(t)

    // Setup test data
    cfg := Defaults()
    cfg.Project.Name = "test-project"

    // Test the functionality
    err := Save(cfg, filepath.Join(dir, ".stax.yml"))
    testutil.AssertNoError(t, err, "save config")

    // Verify results
    loaded, err := Load("", dir)
    testutil.AssertNoError(t, err, "load config")
    testutil.AssertEqual(t, loaded.Project.Name, "test-project")
}
```

### Table-Driven Tests

Use table-driven tests for testing multiple scenarios:

```go
func TestConfigMerge(t *testing.T) {
    tests := []struct {
        name     string
        base     *Config
        override *Config
        want     string
    }{
        {
            name:     "override project name",
            base:     &Config{Project: ProjectConfig{Name: "base"}},
            override: &Config{Project: ProjectConfig{Name: "override"}},
            want:     "override",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := mergeConfigs(tt.base, tt.override)
            if result.Project.Name != tt.want {
                t.Errorf("got %q, want %q", result.Project.Name, tt.want)
            }
        })
    }
}
```

### Integration Tests

Integration tests should use the `integration` build tag:

```go
// +build integration

package integration

import (
    "os"
    "testing"
)

func TestInitWorkflow(t *testing.T) {
    if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
        t.Skip("Skipping integration test")
    }

    // Test complete workflow
    // ...
}
```

### Using Test Helpers

The `testutil` package provides common test helpers:

```go
import "github.com/firecrown-media/stax/pkg/testutil"

// Create temporary directory (auto-cleaned)
dir := testutil.TempDir(t)

// Write test file
testutil.WriteTestFile(t, filepath.Join(dir, "test.txt"), "content")

// Assert file exists
testutil.AssertFileExists(t, filepath.Join(dir, "test.txt"))

// Assert file contains text
testutil.AssertFileContains(t, filepath.Join(dir, "test.txt"), "content")

// Assert no error
testutil.AssertNoError(t, err, "operation description")

// Assert values equal
testutil.AssertEqual(t, got, want)
```

### Using Mocks

The `test/mocks` package provides mock implementations:

```go
import "github.com/firecrown-media/stax/test/mocks"

// Mock WPEngine client
client := mocks.NewMockWPEngineClient()
client.GetInstallFunc = func(install string) (*wpengine.Install, error) {
    return &wpengine.Install{Name: install}, nil
}

// Mock SSH client
ssh := mocks.NewMockSSHClient()
ssh.AddCommandResponse("wp --version", "WP-CLI 2.8.0")

// Mock DDEV manager
ddev := mocks.NewMockDDEVManager()
ddev.WithRunningState(true)
```

## Test Coverage

### Viewing Coverage

```bash
# Generate HTML coverage report
make test-coverage

# Open in browser (macOS)
open coverage/coverage.html

# View function-level coverage in terminal
make test-coverage-func
```

### Coverage Goals

- **Overall**: >70% coverage
- **Core packages** (config, ddev, wordpress): >80% coverage
- **Security package**: >90% coverage
- **Commands**: >60% coverage (due to CLI interaction)

### Improving Coverage

1. Run coverage report to identify gaps:
   ```bash
   make test-coverage-func | grep -v "100.0%"
   ```

2. Focus on:
   - Critical business logic
   - Error handling paths
   - Security-sensitive code
   - Edge cases

3. Don't aim for 100% coverage of:
   - Generated code
   - Simple getters/setters
   - Trivial helper functions

## Continuous Integration

Tests run automatically on GitHub Actions for:

- All pushes to `main`, `develop`, and `feature/**` branches
- All pull requests

### CI Pipeline

1. **Unit Tests**: Run on Go 1.22 and 1.23
2. **Integration Tests**: Run on Go 1.22
3. **Coverage**: Generate and upload to Codecov
4. **Code Quality**: fmt, vet, golangci-lint
5. **Build**: Verify binary builds successfully

### Local CI Simulation

Run the full CI pipeline locally:

```bash
make ci
```

This runs:
- Code formatting
- Go vet
- Linting
- All tests
- Coverage report

## Debugging Tests

### Verbose Output

```bash
# Run with verbose output
go test -v ./pkg/config/

# Show all test names
go test -v ./pkg/config/ | grep -E "^(=== RUN|--- PASS|--- FAIL)"
```

### Running Single Test

```bash
# Run specific test
go test -v -run TestConfigLoad ./pkg/config/

# Run tests matching pattern
go test -v -run TestConfig.* ./pkg/config/
```

### Debug with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug specific test
dlv test ./pkg/config/ -- -test.run TestConfigLoad
```

### Skip Slow Tests

```bash
# Skip tests marked with testing.Short()
go test -short ./...
```

### Test Timeouts

```bash
# Set custom timeout
go test -timeout 30s ./pkg/config/

# Default is 10 minutes
```

## Best Practices

### DO

- ✅ Use descriptive test names
- ✅ Test one thing per test function
- ✅ Use table-driven tests for multiple scenarios
- ✅ Clean up test resources (use `t.Cleanup()`)
- ✅ Use `testutil` helpers for common operations
- ✅ Test error cases thoroughly
- ✅ Mock external dependencies
- ✅ Keep tests fast and independent
- ✅ Use meaningful assertion messages

### DON'T

- ❌ Test implementation details
- ❌ Create interdependent tests
- ❌ Use hardcoded paths
- ❌ Leave test data files in the repo
- ❌ Skip error checking in tests
- ❌ Use `t.Fatal()` in goroutines (use `t.Error()`)
- ❌ Ignore race conditions
- ❌ Write flaky tests

## Environment Variables

Tests can be controlled via environment variables:

- `RUN_INTEGRATION_TESTS=true`: Enable integration tests
- `RUN_E2E_TESTS=true`: Enable end-to-end tests
- `RUN_DDEV_TESTS=true`: Enable DDEV integration tests
- `RUN_WP_TESTS=true`: Enable WordPress CLI tests

Example:
```bash
RUN_INTEGRATION_TESTS=true make test-integration
```

## Troubleshooting

### Tests Fail Locally But Pass in CI

- Check Go version matches CI
- Ensure dependencies are up to date: `go mod download`
- Check for OS-specific issues
- Verify environment variables

### Flaky Tests

- Use `go test -count=100` to reproduce
- Check for race conditions: `go test -race`
- Review test isolation
- Check for timing issues

### Coverage Doesn't Update

- Clean coverage directory: `make clean`
- Regenerate: `make test-coverage`
- Check that new tests actually run

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Go Test Coverage](https://blog.golang.org/cover)
- [Testing Best Practices](https://github.com/golang/go/wiki/CodeReviewComments#tests)

.PHONY: build install test test-unit test-integration test-e2e test-security test-coverage test-all clean dev help fmt vet lint tidy deps version version-build release-snapshot release-dry-run release-check release man man-install man-uninstall man-preview

# Variables
BINARY_NAME=stax
VERSION?=dev
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/firecrown-media/stax/cmd.Version=$(VERSION) \
                  -X github.com/firecrown-media/stax/cmd.GitCommit=$(GIT_COMMIT) \
                  -X github.com/firecrown-media/stax/cmd.BuildDate=$(BUILD_DATE)"

# Test configuration
COVERAGE_DIR=coverage
COVERAGE_PROFILE=$(COVERAGE_DIR)/coverage.out
COVERAGE_HTML=$(COVERAGE_DIR)/coverage.html

# Default target
all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Built $(BINARY_NAME) successfully!"

## install: Install to /usr/local/bin
install: build man
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installing man page..."
	@sudo mkdir -p /usr/local/share/man/man1
	@sudo cp dist/man/stax.1 /usr/local/share/man/man1/
	@sudo mandb 2>/dev/null || sudo makewhatis 2>/dev/null || true
	@echo "Installed successfully! Run 'stax --version' to verify."
	@echo "Run 'man stax' to view the manual."

## test: Run unit tests only (fast)
test: test-unit

## test-unit: Run unit tests with race detection
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -short ./pkg/... ./cmd/...
	@echo "✓ Unit tests complete"

## test-integration: Run integration tests
test-integration:
	@echo "Running integration tests..."
	@RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/...
	@echo "✓ Integration tests complete"

## test-e2e: Run end-to-end tests
test-e2e:
	@echo "Running end-to-end tests..."
	@RUN_E2E_TESTS=true go test -v -tags=e2e ./test/e2e/...
	@echo "✓ End-to-end tests complete"

## test-security: Run security tests
test-security:
	@echo "Running security tests..."
	@go test -v ./pkg/security/...
	@echo "✓ Security tests complete"

## test-coverage: Generate test coverage report
test-coverage:
	@echo "Generating test coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@go tool cover -func=$(COVERAGE_PROFILE) | grep total | awk '{print "Total coverage: " $$3}'
	@echo "✓ Coverage report generated at $(COVERAGE_HTML)"

## test-coverage-func: Show function-level test coverage
test-coverage-func:
	@echo "Showing function-level coverage..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./...
	@go tool cover -func=$(COVERAGE_PROFILE)

## test-all: Run all tests (unit, integration, security)
test-all: test-unit test-integration test-security
	@echo "✓ All tests complete"

## test-verbose: Run all tests with verbose output
test-verbose:
	@echo "Running all tests (verbose)..."
	@go test -v -race ./...

## clean: Clean build artifacts and test coverage
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(COVERAGE_DIR)
	@go clean
	@echo "✓ Clean complete"

## dev: Build and run in development mode
dev: build
	@./$(BINARY_NAME)

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Go vet complete"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed, skipping"

## tidy: Tidy go modules
tidy:
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "✓ Modules tidied"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "✓ Dependencies downloaded"

## verify: Run verification checks (fmt, vet, lint)
verify: fmt vet lint
	@echo "✓ Verification complete"

## ci: Run CI checks (verify + test-all + test-coverage)
ci: verify test-all test-coverage
	@echo "✓ CI checks complete"

## version: Show current version
version:
	@git describe --tags --abbrev=0 2>/dev/null || echo "No version tags yet (use v0.0.0)"

## version-build: Build and show binary version
version-build: build
	@./$(BINARY_NAME) --version

## release-snapshot: Build release snapshot (test release locally)
release-snapshot:
	@echo "Building release snapshot..."
	@goreleaser build --snapshot --clean
	@echo "✓ Snapshot built in ./dist/"

## release-dry-run: Dry run release (test without publishing)
release-dry-run:
	@echo "Running release dry-run..."
	@goreleaser release --snapshot --skip=publish --clean
	@echo "✓ Dry-run complete. Check ./dist/ for artifacts."

## release-check: Validate GoReleaser configuration
release-check:
	@echo "Validating GoReleaser configuration..."
	@goreleaser check
	@echo "✓ Configuration is valid"

## release: Show release instructions
release:
	@echo "Release Process:"
	@echo ""
	@echo "  Option 1 (Recommended): Use GitHub Actions 'Version Bump' workflow"
	@echo "    1. Go to: https://github.com/firecrown-media/stax/actions"
	@echo "    2. Select 'Version Bump' workflow"
	@echo "    3. Click 'Run workflow'"
	@echo "    4. Choose version type (patch/minor/major)"
	@echo ""
	@echo "  Option 2: Manual tag"
	@echo "    git tag -a vX.Y.Z -m 'Release vX.Y.Z'"
	@echo "    git push origin vX.Y.Z"
	@echo ""
	@echo "See docs/RELEASE_PROCESS.md for detailed instructions"

## man: Generate man page
man:
	@echo "Generating man page..."
	@./stax man -o dist/man/ || bash scripts/generate-man.sh

## man-preview: Preview man page
man-preview: man
	@man dist/man/stax.1

## man-install: Install man page
man-install: man
	@echo "Installing man page..."
	@sudo mkdir -p /usr/local/share/man/man1
	@sudo cp dist/man/stax.1 /usr/local/share/man/man1/
	@sudo mandb 2>/dev/null || sudo makewhatis 2>/dev/null || true
	@echo "Man page installed. Run 'man stax' to view."

## man-uninstall: Uninstall man page
man-uninstall:
	@echo "Uninstalling man page..."
	@sudo rm -f /usr/local/share/man/man1/stax.1
	@sudo mandb 2>/dev/null || sudo makewhatis 2>/dev/null || true
	@echo "Man page uninstalled."

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build targets:"
	@echo "  build              - Build the binary"
	@echo "  install            - Install to /usr/local/bin"
	@echo "  clean              - Clean build artifacts"
	@echo ""
	@echo "Test targets:"
	@echo "  test               - Run unit tests (default)"
	@echo "  test-unit          - Run unit tests with race detection"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-e2e           - Run end-to-end tests"
	@echo "  test-security      - Run security tests"
	@echo "  test-coverage      - Generate coverage report"
	@echo "  test-coverage-func - Show function-level coverage"
	@echo "  test-all           - Run all tests"
	@echo "  test-verbose       - Run all tests with verbose output"
	@echo ""
	@echo "Quality targets:"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  lint               - Run golangci-lint"
	@echo "  verify             - Run all verification checks"
	@echo ""
	@echo "Documentation targets:"
	@echo "  man                - Generate man page"
	@echo "  man-preview        - Preview man page"
	@echo "  man-install        - Install man page"
	@echo "  man-uninstall      - Uninstall man page"
	@echo ""
	@echo "Release targets:"
	@echo "  version            - Show current git version tag"
	@echo "  version-build      - Build and show binary version"
	@echo "  release-snapshot   - Build release snapshot locally"
	@echo "  release-dry-run    - Test release without publishing"
	@echo "  release-check      - Validate GoReleaser config"
	@echo "  release            - Show release instructions"
	@echo ""
	@echo "Other targets:"
	@echo "  deps               - Download dependencies"
	@echo "  tidy               - Tidy go modules"
	@echo "  ci                 - Run CI checks"
	@echo "  help               - Show this help message"

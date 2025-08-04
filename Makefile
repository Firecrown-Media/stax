# Stax CLI Makefile

# Build variables
BINARY_NAME=stax
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=build
MAN_DIR=$(BUILD_DIR)/man
LDFLAGS=-ldflags "-X github.com/Firecrown-Media/stax/cmd.Version=$(VERSION)"
MAN_PAGE_NAME=stax.1

# Go variables
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: help build clean test lint install dev run docs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go

build-all: ## Build for all supported platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(MAN_DIR)
	@rm -f $(BINARY_NAME)

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

install: build docs ## Install the binary and man page to Go bin directory
	@echo "Installing $(BINARY_NAME)..."
	@if [ -n "$$GOPATH" ]; then \
		echo "Installing to $$GOPATH/bin/"; \
		mkdir -p $$GOPATH/bin && cp $(BUILD_DIR)/$(BINARY_NAME) $$GOPATH/bin/; \
	else \
		echo "Installing to $$HOME/go/bin/ (creating directory)"; \
		mkdir -p $$HOME/go/bin && cp $(BUILD_DIR)/$(BINARY_NAME) $$HOME/go/bin/; \
		echo "⚠️  Make sure $$HOME/go/bin is in your PATH"; \
	fi
	@echo "Installing man page..."
	@if [ -w "/opt/homebrew/share/man/man1" ]; then \
		echo "Installing man page to /opt/homebrew/share/man/man1/"; \
		mkdir -p /opt/homebrew/share/man/man1 && cp $(MAN_DIR)/$(MAN_PAGE_NAME) /opt/homebrew/share/man/man1/; \
	elif [ -n "$$GOPATH" ] && [ -w "$$GOPATH/share/man/man1" ]; then \
		echo "Installing man page to $$GOPATH/share/man/man1/"; \
		mkdir -p $$GOPATH/share/man/man1 && cp $(MAN_DIR)/$(MAN_PAGE_NAME) $$GOPATH/share/man/man1/; \
	elif [ -w "/usr/local/share/man/man1" ]; then \
		echo "Installing man page to /usr/local/share/man/man1/"; \
		mkdir -p /usr/local/share/man/man1 && sudo cp $(MAN_DIR)/$(MAN_PAGE_NAME) /usr/local/share/man/man1/; \
	else \
		echo "Could not find a writable man page directory."; \
	fi
	@echo "✅ $(BINARY_NAME) installed successfully!"

install-local: build docs ## Install the binary and man page to /usr/local (requires sudo)
	@echo "Installing $(BINARY_NAME) to /usr/local/bin (requires sudo)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installing man page to /usr/local/share/man/man1..."
	@sudo mkdir -p /usr/local/share/man/man1
	@sudo cp $(MAN_DIR)/$(MAN_PAGE_NAME) /usr/local/share/man/man1/
	@echo "✅ $(BINARY_NAME) installed to /usr/local successfully!"

update: ## Quick rebuild and install like 'go install' for development testing
	@echo "🔄 Updating $(BINARY_NAME) installation (Go-style)..."
	@$(MAKE) clean
	@$(MAKE) build
	@echo "Installing to Go bin directory..."
	@if [ -n "$$GOPATH" ]; then \
		echo "Installing to $$GOPATH/bin/$(BINARY_NAME)"; \
		mkdir -p $$GOPATH/bin && cp $(BUILD_DIR)/$(BINARY_NAME) $$GOPATH/bin/; \
	else \
		echo "Installing to $$HOME/go/bin/$(BINARY_NAME)"; \
		mkdir -p $$HOME/go/bin && cp $(BUILD_DIR)/$(BINARY_NAME) $$HOME/go/bin/; \
	fi
	@echo "✅ $(BINARY_NAME) updated successfully!"
	@echo "📋 Test with: $(BINARY_NAME) --version"
	@echo "💡 Make sure your Go bin directory is in PATH"

update-dev: build ## Quick update for development (copies to ./build and shows how to test)
	@echo "🔄 Development build ready!"
	@echo "📁 Binary location: ./$(BUILD_DIR)/$(BINARY_NAME)"
	@echo "📋 Test with: ./$(BUILD_DIR)/$(BINARY_NAME) --version"
	@echo "🔧 To test commands: ./$(BUILD_DIR)/$(BINARY_NAME) [command]"

go-install: ## Install directly like 'go install' (no build artifacts, straight to bin)
	@echo "🚀 Installing $(BINARY_NAME) directly to Go bin (like 'go install')..."
	@if [ -n "$$GOPATH" ]; then \
		echo "Installing to $$GOPATH/bin/$(BINARY_NAME)"; \
		go build $(LDFLAGS) -o $$GOPATH/bin/$(BINARY_NAME) main.go; \
	else \
		echo "Installing to $$HOME/go/bin/$(BINARY_NAME)"; \
		mkdir -p $$HOME/go/bin; \
		go build $(LDFLAGS) -o $$HOME/go/bin/$(BINARY_NAME) main.go; \
	fi
	@echo "✅ $(BINARY_NAME) installed successfully!"
	@echo "📋 Test with: $(BINARY_NAME) --version"
	@echo "💡 Make sure your Go bin directory is in PATH"

uninstall: ## Uninstall the binary from all common locations
	@echo "Uninstalling $(BINARY_NAME)..."
	@removed=false; \
	if [ -f "/opt/homebrew/bin/$(BINARY_NAME)" ]; then \
		echo "Removing from /opt/homebrew/bin/"; \
		rm "/opt/homebrew/bin/$(BINARY_NAME)" && removed=true; \
	fi; \
	if [ -f "$$GOPATH/bin/$(BINARY_NAME)" ]; then \
		echo "Removing from $$GOPATH/bin/"; \
		rm "$$GOPATH/bin/$(BINARY_NAME)" && removed=true; \
	fi; \
	if [ -f "$$HOME/go/bin/$(BINARY_NAME)" ]; then \
		echo "Removing from $$HOME/go/bin/"; \
		rm "$$HOME/go/bin/$(BINARY_NAME)" && removed=true; \
	fi; \
	if [ -f "$$HOME/bin/$(BINARY_NAME)" ]; then \
		echo "Removing from $$HOME/bin/"; \
		rm "$$HOME/bin/$(BINARY_NAME)" && removed=true; \
	fi; \
	if [ -f "/opt/homebrew/share/man/man1/$(MAN_PAGE_NAME)" ]; then \
		echo "Removing man page from /opt/homebrew/share/man/man1/"; \
		rm "/opt/homebrew/share/man/man1/$(MAN_PAGE_NAME)" && removed=true; \
	fi; \
	if [ -n "$$GOPATH" ] && [ -f "$$GOPATH/share/man/man1/$(MAN_PAGE_NAME)" ]; then \
		echo "Removing man page from $$GOPATH/share/man/man1/"; \
		rm "$$GOPATH/share/man/man1/$(MAN_PAGE_NAME)" && removed=true; \
	fi; \
	if [ -f "/usr/local/share/man/man1/$(MAN_PAGE_NAME)" ]; then \
		echo "Removing man page from /usr/local/share/man/man1/"; \
		rm "/usr/local/share/man/man1/$(MAN_PAGE_NAME)" && removed=true; \
	fi; \
	if [ "$$removed" = "true" ]; then \
		echo "✅ $(BINARY_NAME) uninstalled successfully"; \
	else \
		echo "$(BINARY_NAME) not found in any common locations"; \
		echo "Try: which $(BINARY_NAME) to find its location and remove manually"; \
	fi

uninstall-local: ## Uninstall the binary and man page from /usr/local (requires sudo)
	@echo "Uninstalling $(BINARY_NAME) from /usr/local/bin (requires sudo)..."
	@if [ -f "/usr/local/bin/$(BINARY_NAME)" ]; then \
		sudo rm "/usr/local/bin/$(BINARY_NAME)"; \
		echo "✅ $(BINARY_NAME) uninstalled from /usr/local/bin"; \
	else \
		echo "$(BINARY_NAME) not found in /usr/local/bin"; \
	fi
	@echo "Uninstalling man page from /usr/local/share/man/man1 (requires sudo)..."
	@if [ -f "/usr/local/share/man/man1/$(MAN_PAGE_NAME)" ]; then \
		sudo rm "/usr/local/share/man/man1/$(MAN_PAGE_NAME)"; \
		echo "✅ Man page uninstalled from /usr/local/share/man/man1"; \
	else \
		echo "Man page not found in /usr/local/share/man/man1"; \
	fi

uninstall-all: uninstall ## Uninstall binary and remove configuration files
	@echo "Removing configuration files..."
	@if [ -f "$$HOME/.stax.yaml" ]; then \
		echo "Removing global config: $$HOME/.stax.yaml"; \
		rm "$$HOME/.stax.yaml"; \
	fi
	@echo "✅ Complete uninstall finished"
	@echo "Note: Project-specific stax.yaml files are not removed"

dev: ## Build and run in development mode
	@echo "Building and running in development mode..."
	@go run main.go

run: build ## Build and run the binary
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

docs: ## Generate man pages
	@echo "Generating man pages..."
	@go run main.go docs man --dir=$(MAN_DIR)


install-deps: ## Install external dependencies (wp-cli, etc.)
	@echo "Installing external dependencies..."
	@if ! command -v wp >/dev/null 2>&1; then \
		echo "Installing WP-CLI..."; \
		if command -v brew >/dev/null 2>&1; then \
			brew install wp-cli; \
		elif command -v curl >/dev/null 2>&1; then \
			echo "Installing WP-CLI via curl..."; \
			curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar; \
			chmod +x wp-cli.phar; \
			if [ -w "/opt/homebrew/bin" ]; then \
				mv wp-cli.phar /opt/homebrew/bin/wp; \
			elif [ -w "/usr/local/bin" ]; then \
				sudo mv wp-cli.phar /usr/local/bin/wp; \
			else \
				mkdir -p $$HOME/bin; \
				mv wp-cli.phar $$HOME/bin/wp; \
				echo "⚠️  Make sure $$HOME/bin is in your PATH"; \
			fi; \
		else \
			echo "❌ Unable to install WP-CLI: neither brew nor curl found"; \
			echo "Please install WP-CLI manually: https://wp-cli.org/#installing"; \
		fi; \
	else \
		echo "✅ WP-CLI already installed"; \
	fi
	@if ! command -v ddev >/dev/null 2>&1; then \
		echo "⚠️  DDEV not found. Please install DDEV manually:"; \
		echo "   https://ddev.readthedocs.io/en/stable/#installation"; \
	else \
		echo "✅ DDEV already installed"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

check: fmt vet lint test ## Run all checks (format, vet, lint, test)

release: check build-all ## Prepare a release
	@echo "Preparing release $(VERSION)..."
	@mkdir -p $(BUILD_DIR)/release
	@for binary in $(BUILD_DIR)/$(BINARY_NAME)-*; do \
		if [ -f "$$binary" ]; then \
			basename=$$(basename $$binary); \
			tar -czf $(BUILD_DIR)/release/$$basename.tar.gz -C $(BUILD_DIR) $$basename; \
		fi \
	done
	@echo "Release artifacts created in $(BUILD_DIR)/release/"

demo-swap: build ## Demo hot swap functionality
	@echo "Demo: Hot swap functionality"
	@echo "Available commands:"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) swap list     # List available versions and presets"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) swap status   # Show current configuration"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) swap php 8.3  # Switch to PHP 8.3"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) swap preset modern  # Apply modern preset"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) swap --rollback      # Rollback to previous config"
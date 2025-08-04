# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Stax is a Go CLI tool for web development that streamlines local environment setup using DDEV. It offers robust support for content management systems like WordPress and can integrate with hosting providers like WP Engine. It is designed to be versatile for managing multiple projects, from individual developers to large teams. For example, an organization like Firecrown can use Stax to manage a large portfolio of brand websites.

## Architecture

**Core Structure:**
- `main.go` - Entry point that calls `cmd.Execute()`
- `cmd/` - Cobra CLI commands (root, init, setup, status, start, stop, etc.)
- `pkg/` - Reusable packages:
  - `config/` - Configuration management
  - `ddev/` - DDEV integration
  - `wordpress/` - WordPress-specific operations
  - `wpengine/` - WP Engine API client (for hosting integration)
  - `ui/` - User interface components (spinner)
  - `errors/` - Error handling

**Key Dependencies:**
- Cobra (CLI framework)
- Viper (configuration)
- DDEV (local development environments)
- WP-CLI integration

## Build and Development Commands

**Primary commands:**
```bash
make build          # Build binary to build/stax
make install        # Build and install to PATH
make test           # Run all tests
make lint           # Run golangci-lint
make check          # Run fmt, vet, lint, test
make dev            # Run with go run main.go
```

**Testing:**
```bash
go test -v ./...                    # Run tests
go test -v -coverprofile=coverage.out ./...  # With coverage
```

**Linting:**
```bash
golangci-lint run   # Must pass before commits
```

## Configuration

**Global config:** `~/.stax.yaml` for user or team defaults
**Project config:** Local `stax.yaml` for project-specific settings

Configuration uses Viper with:
- YAML files
- Environment variables with `STAX_` prefix
- Command-line flags

## Key Workflows

**Project Environment Setup:**
1. `stax init [project-name]` - Initialize new environment
2. `stax setup [project-name] --install-wp` - Example: Full WordPress setup
3. `stax wpe sync [project]` - Example: Sync a WordPress site from WP Engine

**Hot Swap Environment Management:**
1. `stax swap php 8.3` - Switch PHP version without data loss
2. `stax swap mysql 8.4` - Change database version
3. `stax swap preset modern` - Apply predefined configuration preset
4. `stax swap --rollback` - Revert to previous configuration
5. `stax swap status` - View current environment configuration
6. `stax swap list` - List available versions and presets

**DDEV Integration:**
- All environments use DDEV containers
- Standard PHP 8.2, nginx-fpm, MySQL 8.0
- URLs: `https://[project-name].ddev.site`
- Hot swapping preserves data and container state

**Example Hosting Integration: WP Engine**
- Sync databases and files from production/staging
- Deploy to WP Engine environments
- SSH key authentication

Stax can be extended to support other hosting providers. The built-in WP Engine integration serves as a model for connecting to other platforms.

## Development Guidelines

**Code Style:**
- Follow standard Go conventions
- Use gofmt for formatting
- Pass golangci-lint checks

**Error Handling:**
- Use custom error types in `pkg/errors/`
- Provide helpful error messages
- Log appropriately with verbose flag

**Commands:**
- Use Cobra command structure
- Consistent flag naming
- Support --verbose for debugging

**Testing:**
- Unit tests for all packages
- Integration tests for DDEV/WP Engine operations
- Mock external dependencies

## Hot Swap Environment Management

**Overview:**
The hot swap feature allows rapid switching between different versions of core components without rebuilding environments or losing data.

**Supported Components:**
- PHP versions: 7.4, 8.0, 8.1, 8.2, 8.3, 8.4
- MySQL versions: 5.7, 8.0, 8.4
- Web servers: nginx-fpm, apache-fpm, nginx-fpm-arm64

**Built-in Presets:**
- `legacy` - PHP 7.4, MySQL 5.7 for older WordPress sites
- `stable` - PHP 8.2, MySQL 8.0 (default)
- `modern` - PHP 8.3, MySQL 8.4 for latest WordPress features
- `bleeding-edge` - PHP 8.4, MySQL 8.4 for testing new features
- `performance` - PHP 8.3, MySQL 8.4, optimized for speed
- `compatibility` - PHP 8.1, MySQL 8.0 for maximum plugin compatibility
- `development` - PHP 8.3, MySQL 8.0 with enhanced debugging

**Hot Swap Process:**
1. Automatic backup of current configuration
2. Environment stop with data preservation
3. DDEV configuration update
4. Environment restart with new versions
5. Rollback capability to previous state

**Key Files:**
- `cmd/swap.go` - Main swap command implementation
- `cmd/swap_presets.go` - Preset definitions and backup management
- `.stax.backup.yaml` - Automatic backup of previous configuration
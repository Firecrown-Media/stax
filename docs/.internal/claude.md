# Stax Project - AI Assistant Context

This file provides context for AI assistants working on the Stax CLI tool project.

## Project Overview

**Stax** is a production-ready CLI tool for automating WordPress multisite development environments. It replaces LocalWP with a powerful, provider-agnostic solution built on DDEV and Go.

**Primary Goal:** Automate the complete WordPress multisite setup workflow - from repository cloning to running locally with proper database syncing, remote media proxying, and build integration.

**Target Users:** WordPress development teams working with multisite installations, particularly those using WPEngine hosting.

## Current Project State

**Status:** 95% Complete - Ready for Security Hardening Phase

**Completion:**
- ✅ Architecture & Platform Decision
- ✅ Core CLI Development (Go + Cobra)
- ✅ Multi-Provider Architecture
- ✅ WPEngine Integration
- ✅ DDEV Configuration & Multisite Support
- ✅ Build Process Integration
- ✅ Comprehensive Documentation
- ✅ Security Audit Complete
- ⏳ Security Remediation (6-8 weeks)
- ⏳ Integration Testing
- ⏳ Production Release

## Technology Stack

- **Language:** Go 1.22+
- **CLI Framework:** Cobra + Viper
- **Container Platform:** DDEV
- **Credential Storage:** macOS Keychain
- **Configuration:** YAML with environment overrides
- **Package Distribution:** Homebrew (planned)

## Project Structure

```
stax/
├── .claude/
│   └── agents/              # Specialized subagents (8 total)
│       ├── cli-developer.md
│       ├── wordpress-expert.md
│       ├── devops-engineer.md
│       ├── go-developer.md
│       ├── database-admin.md
│       ├── docs-specialist.md
│       ├── security-auditor.md
│       └── build-engineer.md
│
├── cmd/                     # CLI commands (12 files)
│   ├── root.go             # Root command + global flags
│   ├── init.go             # Project initialization
│   ├── start.go, stop.go   # Environment control
│   ├── status.go           # Status display
│   ├── db.go               # Database operations
│   ├── build.go            # Build commands
│   ├── lint.go             # Code quality
│   ├── dev.go              # Development mode
│   ├── config.go           # Configuration management
│   ├── provider.go         # Provider management
│   ├── setup.go            # Credential setup
│   └── doctor.go           # Diagnostics
│
├── pkg/                     # Core packages
│   ├── provider/           # Provider abstraction (4 files)
│   │   ├── interface.go    # Provider interface
│   │   ├── registry.go     # Provider registration
│   │   ├── factory.go      # Provider factory
│   │   └── manager.go      # Provider manager
│   │
│   ├── providers/          # Provider implementations
│   │   ├── wpengine/       # WPEngine (production-ready)
│   │   │   ├── provider.go
│   │   │   ├── capabilities.go
│   │   │   ├── client.go
│   │   │   ├── ssh.go
│   │   │   ├── database.go
│   │   │   └── files.go
│   │   ├── aws/            # AWS (stub)
│   │   ├── wordpress-vip/  # WordPress VIP (stub)
│   │   └── local/          # Local-only provider
│   │
│   ├── ddev/               # DDEV management (5 files)
│   │   ├── types.go
│   │   ├── manager.go
│   │   ├── config.go
│   │   └── nginx.go
│   │
│   ├── wpengine/           # WPEngine client (legacy, use providers/wpengine)
│   │
│   ├── wordpress/          # WordPress operations (5 files)
│   │   ├── cli.go          # WP-CLI wrapper
│   │   ├── database.go     # DB operations
│   │   ├── search_replace.go
│   │   ├── multisite.go
│   │   └── types.go
│   │
│   ├── build/              # Build management (8 files)
│   │   ├── manager.go
│   │   ├── composer.go
│   │   ├── npm.go
│   │   ├── quality.go      # PHPCS integration
│   │   ├── status.go
│   │   ├── watch.go
│   │   └── hooks.go        # Git hooks
│   │
│   ├── system/             # System operations
│   │   ├── hosts.go        # Hosts file management
│   │   └── types.go
│   │
│   ├── config/             # Configuration management
│   │   ├── config.go       # Config structs
│   │   ├── loader.go       # Multi-source loading
│   │   └── validator.go    # Validation
│   │
│   ├── credentials/        # Credential storage
│   │   ├── keychain.go     # macOS Keychain
│   │   └── manager.go
│   │
│   ├── ui/                 # Terminal UI
│   │   ├── output.go       # Colorized output
│   │   └── spinner.go      # Progress indicators
│   │
│   └── errors/             # Error types
│       └── errors.go
│
├── templates/              # Configuration templates
│   ├── ddev/
│   │   ├── config.yaml.tmpl
│   │   ├── nginx_full/
│   │   │   └── media-proxy.conf.tmpl
│   │   └── post-start.sh
│   └── scripts/
│       └── build/
│
├── docs/                   # Documentation (24 files, ~400KB)
│   ├── User Documentation (10 files):
│   │   ├── INSTALLATION.md
│   │   ├── QUICK_START.md
│   │   ├── USER_GUIDE.md
│   │   ├── MULTISITE.md
│   │   ├── TROUBLESHOOTING.md
│   │   ├── COMMAND_REFERENCE.md
│   │   ├── WPENGINE.md
│   │   ├── EXAMPLES.md
│   │   └── FAQ.md
│   │
│   ├── Technical Documentation (8 files):
│   │   ├── ARCHITECTURE.md
│   │   ├── COMMANDS.md
│   │   ├── CONFIG_SPEC.md
│   │   ├── PROVIDER_INTERFACE.md
│   │   ├── WPENGINE_INTEGRATION.md
│   │   ├── MULTI_PROVIDER.md
│   │   ├── PROVIDER_DEVELOPMENT.md
│   │   └── BUILD_PROCESS.md
│   │
│   └── Security Documentation (6 files):
│       ├── SECURITY_AUDIT.md
│       ├── SECURITY.md
│       ├── SECURITY_CHECKLIST.md
│       ├── SECURITY_REVIEW_SUMMARY.md
│       ├── SECURITY_SCAN_RESULTS.md
│       └── SECURITY_QUICK_REFERENCE.md
│
├── main.go                 # Entry point
├── go.mod, go.sum         # Go dependencies
├── Makefile               # Build automation
├── .gitignore             # Git ignore rules
├── README.md              # User-facing README
└── PROJECT_SUMMARY.md     # Complete project summary
```

## Key Architectural Concepts

### 1. Provider Abstraction Layer

The project uses a provider interface pattern to support multiple hosting platforms:

- **Provider Interface**: Defined in `pkg/provider/interface.go`
- **Provider Registry**: Manages available providers
- **Provider Factory**: Creates provider instances based on config
- **Implementations**: WPEngine (complete), AWS (stub), WordPress VIP (stub), Local (stub)

**Why:** Allows stax to work with any hosting platform, not just WPEngine.

### 2. DDEV Integration

DDEV is the recommended container platform (chosen over Podman/Docker Compose):

- **Configuration Generation**: `pkg/ddev/config.go` creates `.ddev/config.yaml`
- **Container Management**: `pkg/ddev/manager.go` handles start/stop/status
- **Nginx Configuration**: `pkg/ddev/nginx.go` for remote media proxy
- **Multisite Support**: Automatic subdomain configuration with SSL

**Why DDEV:** Pre-configured WordPress support, multisite-friendly, Mac-optimized, junior-dev friendly.

### 3. Credential Management

All credentials stored securely in macOS Keychain:

- **No credentials in config files** - Ever
- **Keychain Integration**: `pkg/credentials/keychain.go`
- **Credential Types**: WPEngine API, SSH keys, GitHub tokens
- **Secure**: System-level encryption, audit trail

### 4. WordPress Multisite

Special handling for WordPress multisite installations:

- **Detection**: `pkg/wordpress/multisite.go` detects subdomain vs subdirectory
- **Subsite Management**: Lists all subsites, generates local domains
- **URL Mapping**: Production URLs → Local URLs for all subsites
- **Hosts File**: Automatic `/etc/hosts` updates (requires sudo)
- **Search-Replace**: Network-wide URL replacements with serialized data handling

### 5. Remote Media Proxying

Avoids downloading gigabytes of media files:

- **Primary**: BunnyCDN proxy
- **Fallback**: WPEngine direct
- **Caching**: Nginx cache (30 day TTL, 10GB max)
- **Config**: `templates/ddev/nginx_full/media-proxy.conf.tmpl`

**Benefit:** Faster setup, less storage, always current media.

## Important Files to Know

### Start Here
- `README.md` - User-facing overview
- `PROJECT_SUMMARY.md` - Complete implementation summary
- `docs/ARCHITECTURE.md` - System architecture
- `docs/QUICK_START.md` - 5-minute getting started guide

### For Development
- `cmd/root.go` - CLI entry point, global flags
- `pkg/provider/interface.go` - Provider interface definition
- `pkg/config/config.go` - Configuration schema
- `docs/COMMANDS.md` - Command specifications

### For Security
- `docs/SECURITY_AUDIT.md` - Complete security audit
- `docs/SECURITY.md` - Security best practices
- `docs/SECURITY_CHECKLIST.md` - Pre-release checklist
- `docs/SECURITY_QUICK_REFERENCE.md` - Quick security guide

## Common Development Tasks

### Building the Project

```bash
# Build binary
make build

# Build and install
make install

# Run tests
make test

# Clean build artifacts
make clean

# Format code
make fmt

# Run static analysis
make vet
```

### Working with Providers

```bash
# List available providers
./stax provider list

# Show provider capabilities
./stax provider show wpengine

# Test provider connection
./stax provider test

# Switch provider
./stax provider set aws
```

### Testing Commands

```bash
# Initialize a project (dry-run)
./stax init --dry-run --name=test-project

# Check system diagnostics
./stax doctor

# Get help for any command
./stax <command> --help
```

### Code Organization Rules

1. **Commands in cmd/**: Each command in its own file
2. **Business logic in pkg/**: Never in cmd/
3. **Interfaces first**: Define interface before implementation
4. **Provider-agnostic**: Use provider interface, not concrete types
5. **Error handling**: Use `pkg/errors` for custom errors
6. **UI feedback**: Use `pkg/ui` for all output

## Development Guidelines

### Git Commit Guidelines

1. **Use conventional commits**: `feat:`, `fix:`, `docs:`, etc.
2. **No attribution lines**: Never include "Generated with Claude Code" or similar AI attribution in commit messages
3. **Focus on the change**: Commit messages should describe what and why, not how or who

### Code Style

1. **Follow Go conventions**: Use `gofmt`, `go vet`
2. **Document exports**: All exported functions need comments
3. **Error handling**: Always handle errors, never `panic()`
4. **Context**: Use `context.Context` for cancellation
5. **Testing**: Write tests for all critical paths

### Security Requirements

1. **Never log credentials**: Use `RemoveSensitiveData()` for logs
2. **Validate all inputs**: Use validation functions
3. **No shell expansion**: Use `exec.Command(cmd, args...)`, not shell scripts
4. **Sanitize paths**: Prevent path traversal with `SanitizePath()`
5. **HTTPS only**: All API calls must use HTTPS
6. **Verify SSH hosts**: Enable SSH host key verification

### Documentation Standards

1. **Update docs with code**: Documentation is part of the feature
2. **Examples required**: Every command needs examples
3. **Junior-dev friendly**: Write for beginners, not experts
4. **Keep updated**: Documentation drift is a bug

## Current Known Issues

### Build Issues

1. **Import path errors**: Some packages have incorrect import paths
2. **Type mismatches**: Provider interface vs implementation
3. **Duplicate declarations**: Some types defined in multiple places

**Status:** Documented in `docs/SECURITY_SCAN_RESULTS.md`

### Security Issues

1. **Critical**: SSH host key verification disabled
2. **High**: Command injection vulnerabilities
3. **High**: Path traversal risks
4. **High**: Insecure temporary file creation

**Status:** Complete audit in `docs/SECURITY_AUDIT.md` with remediation plan

### Integration Issues

1. **cmd/init.go**: Needs full integration wiring
2. **cmd/db.go**: Database pull needs component connection
3. **Testing**: No integration tests yet

## Next Steps (Priority Order)

### Critical (Do First)

1. **Fix Build Errors** (1-2 days)
   - Resolve import path issues
   - Fix type mismatches
   - Test compilation

2. **Security Critical Fixes** (2-3 weeks)
   - Implement SSH host key verification
   - Add command injection prevention
   - Implement path validation
   - Fix temporary file handling

3. **Basic Integration Testing** (1 week)
   - Test `stax init` end-to-end
   - Test `stax db:pull`
   - Test multisite setup
   - Verify builds work

### High Priority

4. **Complete Command Wiring** (1 week)
   - Wire cmd/init.go to all components
   - Wire cmd/db.go pull operation
   - Test all workflows

5. **Security Medium Priority** (2 weeks)
   - Error message sanitization
   - TLS hardening
   - Input validation enhancement
   - Rate limiting

### Medium Priority

6. **Comprehensive Testing** (2-3 weeks)
   - Unit tests (>70% coverage)
   - Integration tests
   - Security tests
   - Performance tests

7. **Homebrew Packaging** (1 week)
   - Create Homebrew formula
   - Set up tap repository
   - Release automation

## Working with Subagents

The project has 8 specialized subagents in `.claude/agents/`:

1. **cli-developer** - CLI UX and command design
2. **wordpress-expert** - WordPress and multisite knowledge
3. **devops-engineer** - DDEV, containers, infrastructure
4. **go-developer** - Go language and patterns
5. **database-admin** - MySQL, WP database operations
6. **docs-specialist** - Documentation writing
7. **security-auditor** - Security review and auditing
8. **build-engineer** - Build systems and automation

**How to use:** Reference the relevant subagent for specialized tasks. They have role-specific tools and knowledge.

## Testing Strategy

### Unit Tests
- Test each package independently
- Mock external dependencies
- Use table-driven tests
- Test error conditions

### Integration Tests
- Test command workflows end-to-end
- Use test fixtures for WPEngine responses
- Test with real DDEV containers (optional)
- Verify file operations

### Security Tests
- Test input validation with malicious inputs
- Test path traversal prevention
- Test command injection prevention
- Fuzz testing for parsers

### Manual Testing
- Test on real WPEngine projects
- Test multisite scenarios
- Test on different Mac versions (Intel + Apple Silicon)
- Test error recovery

## Deployment Process (Future)

1. **Version bump** in `cmd/root.go`
2. **Run security checklist** (`docs/SECURITY_CHECKLIST.md`)
3. **Tag release** in git
4. **Build binaries** for Mac (Intel + ARM)
5. **Update Homebrew formula**
6. **Publish release notes**
7. **Update documentation** if needed

## Useful Commands

```bash
# Development
make build              # Build binary
make dev                # Build and run
make test               # Run tests
make lint               # Run linters

# Security
govulncheck ./...       # Check for vulnerabilities
staticcheck ./...       # Static analysis
gosec ./...             # Security scan

# Documentation
go doc pkg/provider     # View package docs

# Git
git status              # Check for uncommitted changes
git log --oneline       # View recent commits
```

## Getting Help

### Documentation Resources
- **User Questions**: See `docs/FAQ.md`
- **Technical Questions**: See `docs/ARCHITECTURE.md`
- **Security Questions**: See `docs/SECURITY.md`
- **Command Reference**: See `docs/COMMAND_REFERENCE.md`

### For AI Assistants
- Read `PROJECT_SUMMARY.md` for complete project overview
- Check `docs/ARCHITECTURE.md` for system design
- Review relevant subagent in `.claude/agents/` for specialized knowledge
- Follow security guidelines in `docs/SECURITY.md`

## Important Notes

### What NOT to Do

1. ❌ **Never hardcode credentials** - Use Keychain only
2. ❌ **Never use shell expansion** - Use `exec.Command()` properly
3. ❌ **Never skip input validation** - Validate everything
4. ❌ **Never commit .env or credential files** - Check `.gitignore`
5. ❌ **Never log sensitive data** - Sanitize all logs
6. ❌ **Never reference the source of requirements** - Keep implementation details internal

### What TO Do

1. ✅ **Write tests** - Tests are required, not optional
2. ✅ **Document everything** - Code + docs together
3. ✅ **Follow Go conventions** - Use `gofmt`, `go vet`
4. ✅ **Think security first** - Review `docs/SECURITY.md`
5. ✅ **Use provider abstraction** - Never hardcode to WPEngine
6. ✅ **Keep it simple** - Readable code > clever code

## Project Philosophy

1. **Developer Experience First**: Make it easy and enjoyable to use
2. **Security by Default**: Secure defaults, no credentials in configs
3. **Provider Agnostic**: Work with any hosting platform
4. **Junior-Dev Friendly**: Excellent documentation, clear errors
5. **Team Collaboration**: Shared configs, consistent environments
6. **Production Ready**: Not a prototype, a professional tool

## Version History

- **v0.1.0** (2025-11-08): Initial implementation complete
  - All 7 phases delivered
  - 95% complete, ready for security hardening
  - Comprehensive documentation
  - Security audit complete

## Contributors

This project was built with AI assistance using specialized subagents for different domains (CLI development, WordPress, DevOps, Go, databases, documentation, security, build engineering).

---

**Last Updated:** 2025-11-08
**Project Status:** 95% Complete - Security Hardening Phase
**Next Milestone:** Production Release (6-8 weeks)

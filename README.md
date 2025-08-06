# Stax

**A specialized CLI tool for Firecrown's web development workflow, built to streamline local development and hosting provider integration.**

Stax was developed as a custom solution for Firecrown's development team to automate local environment setup using DDEV while providing seamless integration with hosting platforms. The tool is purpose-built for WordPress development but designed with extensibility in mind for future needs.

## Table of Contents

- [About This Project](#about-this-project)
- [Key Capabilities](#key-capabilities)
- [Getting Started](#getting-started)
- [Installation](#installation)
  - [System Requirements](#system-requirements)
  - [Install via Homebrew (Recommended)](#install-via-homebrew-recommended)
  - [Install from Source](#install-from-source)
  - [Development Installation](#development-installation)
  - [Uninstalling Stax](#uninstalling-stax)
- [Hosting Provider Integration](#hosting-provider-integration)
- [WP Engine Integration](#wp-engine-integration)
  - [Setup Requirements](#setup-requirements)
  - [Configuration](#configuration)
- [Common Firecrown Workflows](#common-firecrown-workflows)
- [WP Engine Commands Reference](#wp-engine-commands-reference)
- [Stax Commands Reference](#stax-commands-reference)
- [Hot Swap Environment Management](#hot-swap-environment-management)
- [Release Management & Updates](#release-management--updates)
- [Extending Stax](#extending-stax)
- [Troubleshooting](#troubleshooting)
- [Bug Reports & Feature Requests](#bug-reports--feature-requests)
- [Support & Maintenance](#support--maintenance)
- [Project Background](#project-background)

## About This Project

This CLI was created by an independent software engineer under contract with Firecrown to address specific workflow challenges in managing multiple brand projects. The tool focuses on:

- Rapid local development environment setup
- Streamlined WordPress project initialization
- Direct integration with WP Engine hosting
- Hot-swappable environment configurations
- Consistent development practices across team members

## Getting Started

**New to Stax?** Start with our comprehensive [Developer Onboarding Guide](ONBOARDING.md) which walks through the complete development lifecycle from initial setup through production deployment.

The onboarding guide covers:
- Complete installation and setup process
- WP Engine integration configuration
- Real-world development workflow examples
- Git integration and deployment pipelines
- Troubleshooting common issues

For quick reference, continue reading below, or jump directly to [Installation](#installation) if you're already familiar with the workflow.

## Key Capabilities

- **🚀 Rapid Project Setup**: Initialize new WordPress development environments in under 2 minutes
- **🔄 Production Sync**: One-command sync from WP Engine production/staging to local development
- **⚡ Hot Swap Management**: Switch PHP/MySQL versions instantly without losing data or configuration
- **🌐 Multi-Project Workflows**: Manage dozens of brand projects simultaneously
- **📦 Team Standardization**: Enforce consistent development environments across Firecrown's team
- **🔌 Extensible Architecture**: Built to accommodate future hosting provider integrations

## Installation

### System Requirements

Before installing Stax, ensure these dependencies are available:

**Required:**

- **Docker Provider**: Container runtime for DDEV environments

  - [Docker Desktop](https://www.docker.com/products/docker-desktop/) for Mac/Windows
  - [Docker Engine](https://docs.docker.com/engine/install/) for Linux
  - [OrbStack](https://orbstack.dev/): Recommended, easiest to install, most performant, commercial, not open-source
  - [Lima](https://lima-vm.io/): Free, open-source
  - [Rancher Desktop](https://rancherdesktop.io/): Free, open-source, simple installation, slower startup
  - [Colima](https://github.com/abiosoft/colima): Free, open-source. Depends on separate Lima installation (managed by Homebrew)

- **DDEV**: Local development environment manager
  - [Install DDEV](https://ddev.readthedocs.io/en/stable/#installation)

**Recommended:**

- **WP-CLI**: WordPress command-line interface
  - [Install WP-CLI](https://wp-cli.org/#installing) or use `make install-deps`

### Install via Homebrew (Recommended)

```bash
# Add Firecrown's Stax tap
brew tap firecrown-media/stax

# Install Stax
brew install stax

# Verify installation
stax --version
```

### Install from Source

For development or customization:

```bash
# Clone the repository
git clone https://github.com/Firecrown-Media/stax.git
cd stax

# Install dependencies (optional)
make install-deps

# Build and install to Go bin directory
make install

# Verify installation
stax --version
```

### Development Installation

For contributors working on Stax improvements:

```bash
# Quick rebuild and install (like 'go install')
make update

# Build locally for testing without system installation
make update-dev
./build/stax --version

# Direct install to Go bin directory
make go-install
```

### Uninstalling Stax

#### Uninstall via Homebrew

```bash
# Remove Stax installation
brew uninstall stax

# Remove the tap (optional)
brew untap firecrown-media/stax
```

#### Uninstall from Source Installation

```bash
# Basic uninstall (removes binary and man pages)
make uninstall

# Complete uninstall (removes binary, man pages, and global config)
make uninstall-all

# If installed to /usr/local (requires sudo)
make uninstall-local
```

#### Manual Cleanup

If needed, you can manually remove configuration files:

```bash
# Remove global configuration
rm ~/.stax.yaml

# Project-specific configurations are left intact in project directories
# Remove manually if desired: rm [project-directory]/stax.yaml
```

## Hosting Provider Integration

Stax was architected with extensible hosting provider support. Currently implements full WP Engine integration with the framework designed to accommodate additional providers.

### Current Platform Support

- **✅ WP Engine** - Complete sync, deploy, and management integration
- **🔄 Future Platforms** - Architecture ready for other hosting solutions

### Integration Architecture

The hosting provider system was designed as a pluggable architecture:

```go
type HostingProvider interface {
    Connect(credentials) error
    Sync(project, options) error
    Deploy(project, environment) error
    GetInfo(project) (*ProjectInfo, error)
}
```

This allows Firecrown to easily add new hosting providers as brand needs dictate, without major architectural changes.

## WP Engine Integration

### Setup Requirements

The WP Engine integration requires several prerequisite steps. Work with the admin on the Firecrown team to ensure proper access:

#### 1. WP Engine Account Access

- **Contact**: Wayne/Nick
- **Required Permissions**:
  - Access to target site installations
  - API access permissions for programmatic operations
  - SSH access permissions for secure file/database operations

#### 2. WP Engine API Configuration

- **Purpose**: Enables automated access to site installations and environments
- **Setup Guide**: [WP Engine API Documentation](https://wpengine.com/support/enabling-wp-engine-api/)
- **Deliverable**: API credentials for secure authentication

#### 3. SSH Key Setup

- **Purpose**: Secure authentication for database and file synchronization operations
- **Documentation**:
  - [SSH Keys for Shell Access](https://wpengine.com/support/ssh-keys-for-shell-access/)
  - [SSH Gateway Key Configuration](https://wpengine.com/support/ssh-gateway/#Add_SSH_Key)
- **Required**: Add your development machine's SSH public key to your WP Engine account

### Configuration

#### Environment Variables

Store WP Engine credentials securely (recommend using a password manager):

```bash
# WP Engine API credentials (from Step 2 above)
export WPE_USERNAME=your-wpe-username
export WPE_PASSWORD=your-secure-api-password

# Optional: Set default environment preferences
export WPE_DEFAULT_ENV=production
```

#### Global Team Configuration

Create a team-wide configuration file (`~/.stax.yaml`):

```yaml
# Standard Firecrown development environment settings
php_version: "8.2"
webserver: "nginx-fpm"
database: "mysql:8.0"

# Default WordPress configuration for new projects
wordpress_defaults:
  admin_user: "fcadmin"
  admin_email: "dev@firecrown.com"

# WP Engine integration settings
hosting:
  wpengine:
    username: "your-wpe-username"
    sync_defaults:
      skip_media: true # Use production CDN, faster sync
      exclude_dirs:
        - "wp-content/cache/"
        - "wp-content/uploads/backup-*"
        - "wp-content/ai1wm-backups/"
    ssh_key_path: "~/.ssh/id_rsa"

# Standard Firecrown WordPress plugin stack
default_plugins:
  - "advanced-custom-fields-pro"
  - "yoast-seo"
  - "wp-rocket"
```

#### Project-Specific Configuration

Each site project can have custom settings (`project-directory/stax.yaml`):

```yaml
name: "site-project-local"
type: "wordpress"

# Environment overrides for specific site requirements
php_version: "8.1" # site legacy requirement
database: "mysql:8.0"

# WordPress configuration
wordpress:
  url: "https://site-project-local.ddev.site"
  title: "site Project - Firecrown Development"

# WP Engine site installation mapping
hosting:
  wpengine:
    install_name: "siteprojectprod"
    environment: "production"

# site-specific plugin requirements
plugins:
  - "woocommerce"
  - "gravityforms"
  - "site-custom-plugin"
```

## Common Firecrown Workflows

### New site Project Onboarding

When taking on a new client with an existing WordPress installation:

```bash
# Create local development environment for site assessment
mkdir site-name-assessment && cd site-name-assessment
stax init site-name-local --php-version=8.2

# Mirror their existing WP Engine site for analysis
stax wpe sync site-name-install

# Document current setup for project planning
stax wp plugin list --status=active > site-tech-stack.txt
stax wp theme list --status=active >> site-tech-stack.txt

# Site ready for development at: https://site-name-local.ddev.site
```

### Feature Development Workflow

Developing new features for an existing site:

```bash
# Create feature-specific environment
mkdir site-feature-development && cd site-feature-development
stax init site-feature --php-version=8.2

# Sync latest production data for realistic testing
stax wpe sync client-install --skip-files  # Database only for speed

# Install development and testing tools
stax wp plugin install query-monitor --activate
stax wp plugin install log-deprecated-notices --activate

# Ready for feature development
```

### Environment Testing with Hot Swap

Testing plugin compatibility across different PHP versions:

```bash
# Setup compatibility testing environment
mkdir plugin-compatibility-test && cd plugin-compatibility-test
stax init compat-test --php-version=8.2

# Test with stable environment (current Firecrown standard)
stax swap preset stable
stax wp plugin install new-client-plugin --activate
stax wp eval "echo 'Testing on PHP ' . PHP_VERSION;"

# Test compatibility with modern environment
stax swap preset modern
stax wp plugin activate new-client-plugin
stax wp eval "echo 'Testing on PHP ' . PHP_VERSION;"

# Rollback if issues discovered
stax swap --rollback

# Document results for client recommendation
echo "Plugin compatible with stable (PHP 8.2) and modern (PHP 8.3)" > compatibility-report.txt
```

## WP Engine Commands Reference

| Command                                         | Purpose                                    | Usage Example                                    |
| ----------------------------------------------- | ------------------------------------------ | ------------------------------------------------ |
| `stax wpe list`                                 | List accessible WP Engine installations    | `stax wpe list`                                  |
| `stax wpe info [install]`                       | Show installation details and environments | `stax wpe info clientsite`                       |
| `stax wpe sync [install]`                       | Sync from WP Engine to local development   | `stax wpe sync clientsite`                       |
| `stax wpe sync [install] --environment=staging` | Sync from specific environment             | `stax wpe sync clientsite --environment=staging` |
| `stax wpe sync [install] --skip-files`          | Database-only sync (faster)                | `stax wpe sync clientsite --skip-files`          |
| `stax wpe sync [install] --skip-database`       | Files-only sync                            | `stax wpe sync clientsite --skip-database`       |

## Stax Commands Reference

### Core Commands

| Command | Description | Required Flags | Optional Flags |
|---------|-------------|----------------|----------------|
| `stax init [project-name]` | Initialize new WordPress development environment | None | `--path`, `--php-version`, `--webserver`, `--database` |
| `stax start [project]` | Start DDEV environment | None | `--path` |
| `stax stop [project]` | Stop DDEV environment | None | `--path` |
| `stax status` | Show status of all environments | None | None |
| `stax poweroff` | Stop all DDEV environments | None | None |
| `stax delete [project]` | Delete DDEV environment and data | None | `--path`, `--force` |

### WP Engine Integration Commands

| Command | Description | Required Flags | Optional Flags |
|---------|-------------|----------------|----------------|
| `stax wpe list` | List accessible WP Engine installations | None | `--username`, `--api-key` |
| `stax wpe connect [install]` | Test connection to WP Engine install | `[install-name]` | `--username`, `--environment` |
| `stax wpe sync [install]` | Sync database and files from WP Engine | `[install-name]` | `--skip-files`, `--skip-database`, `--skip-media`, `--delete-local`, `--suppress-debug`, `--create-upload-redirect` |

### WP Engine Database Commands

| Command | Description | Required Flags | Optional Flags |
|---------|-------------|----------------|----------------|
| `stax wpe db download [install]` | Download database from WP Engine | `[install-name]` | `--output`, `--username`, `--environment` |
| `stax wpe db import [file]` | Import WP Engine database to local | `[database-file]` | `--suppress-debug`, `--create-upload-redirect` |
| `stax wpe db analyze` | Analyze media URLs in WordPress database | None | `--path` |
| `stax wpe db rewrite [install]` | Rewrite URLs for local development | `[install-name]` | `--suppress-debug`, `--create-upload-redirect` |
| `stax wpe db diagnose` | Diagnose media URL routing behavior | None | `--path` |

### Hot Swap Commands

| Command | Description | Required Flags | Optional Flags |
|---------|-------------|----------------|----------------|
| `stax swap <component> <version>` | Swap component version | `<component>`, `<version>` | `--force` |
| `stax swap preset <preset-name>` | Apply predefined preset | `<preset-name>` | `--force` |
| `stax swap --rollback` | Rollback to previous configuration | None | `--force` |
| `stax swap list` | List available versions and presets | None | None |
| `stax swap status` | Show current environment configuration | None | None |

### WordPress Commands

| Command | Description | Required Flags | Optional Flags |
|---------|-------------|----------------|----------------|
| `stax wp <command>` | Execute WP-CLI commands in DDEV environment | `<command>` | `--path` |

### Global Flags

Available on all commands:

| Flag | Description | Default |
|------|-------------|---------|
| `--verbose`, `-v` | Enable verbose output | `false` |
| `--config` | Specify config file path | `$HOME/.stax.yaml` |
| `--path`, `-p` | Specify project path | Current directory |

### Flag Details

#### Init Command Flags
- `--path`, `-p`: Directory to initialize project (default: current directory)
- `--php-version`: PHP version to use (default: "8.2", options: 7.4, 8.0, 8.1, 8.2, 8.3, 8.4)
- `--webserver`: Web server type (default: "nginx-fpm", options: "nginx-fpm", "apache-fpm")
- `--database`: Database type and version (default: "mysql:8.0", options: "mysql:5.7", "mysql:8.0", "mysql:8.4")

#### WP Engine Global Flags
- `--username`: WP Engine username (or set `WPE_USERNAME` env var)
- `--api-key`: WP Engine API key (or set `WPE_API_KEY` env var)
- `--environment`: WP Engine environment (default: "production", options: "production", "staging", "development")
- `--install`: WP Engine install name (alternative to positional argument)

#### WP Engine Sync Flags
- `--skip-files`: Skip syncing WordPress files (default: false)
- `--skip-database`: Skip syncing database (default: false)
- `--skip-media`: Skip syncing media files (default: true)
- `--delete-local`: Delete local files not present on remote (WARNING: dangerous)
- `--suppress-debug`: Suppress WordPress debug notices for cleaner output
- `--create-upload-redirect`: Create must-use plugin to redirect upload URLs to remote

#### Hot Swap Components
- `php`: PHP versions (7.4, 8.0, 8.1, 8.2, 8.3, 8.4)
- `mysql`: MySQL versions (5.7, 8.0, 8.4)
- `webserver`: Web server types (nginx-fpm, apache-fpm)
- `preset`: Predefined combinations (legacy, stable, modern, bleeding-edge, performance, compatibility, development)

### Examples

```bash
# Initialize new project with PHP 8.3
stax init my-project --php-version=8.3

# Sync database only from WP Engine
stax wpe sync clientsite --skip-files

# Switch to modern preset (PHP 8.3, MySQL 8.4)
stax swap preset modern

# Execute WP-CLI command
stax wp plugin list --status=active

# Check environment status
stax status
```

## Hot Swap Environment Management

One of Stax's specialized features is the ability to quickly switch environment components without losing data:

```bash
# List available versions and presets
stax swap list

# Switch PHP version for testing
stax swap php 8.3

# Apply Firecrown's standard modern preset
stax swap preset modern

# Rollback if compatibility issues arise
stax swap --rollback

# Check current environment configuration
stax swap status
```

### Firecrown Environment Presets

- **`stable`** - PHP 8.2, MySQL 8.0 (current Firecrown standard)
- **`modern`** - PHP 8.3, MySQL 8.4 (forward-looking development)
- **`legacy`** - PHP 7.4, MySQL 5.7 (older client maintenance)
- **`performance`** - PHP 8.3, MySQL 8.4, optimized for high-traffic sites
- **`compatibility`** - PHP 8.1, MySQL 8.0, maximum plugin compatibility
- **`bleeding-edge`** - PHP 8.4, MySQL 8.4 (experimental testing)

## Release Management & Updates

Stax uses automated release management designed for Firecrown's development workflow:

### Automated Releases

New versions are released through GitHub Actions triggered by version tags:

```bash
# Create new release (maintainer only)
git tag v1.2.3
git push origin v1.2.3

# This automatically:
# 1. Builds cross-platform binaries (macOS Intel/ARM, Linux)
# 2. Creates GitHub release with downloadable assets
# 3. Updates Homebrew tap for team installation
# 4. Synchronizes `stax --version` with the git tag
```

### Team Updates

Firecrown team members receive updates through Homebrew:

```bash
# Update to latest version
brew upgrade stax

# Verify new version
stax --version
```

## Extending Stax

### Adding New Hosting Providers

As Firecrown's client base grows, new hosting providers can be added to Stax:

1. **Implement the interface** in `pkg/[provider]/`
2. **Add CLI commands** in `cmd/[provider].go`
3. **Update configuration structures** in `pkg/config/`
4. **Add documentation** and usage examples
5. **Write integration tests**

Example structure for a new provider (e.g., Pantheon):

```
pkg/
  pantheon/
    client.go      # API client implementation
    database.go    # Sync operations
    types.go       # Provider-specific types
cmd/
  pantheon.go      # CLI commands (pantheon sync, pantheon deploy, etc.)
```

### Development Contributions

For team members contributing to Stax development:

```bash
# Clone and setup development environment
git clone https://github.com/Firecrown-Media/stax.git
cd stax

# Install development dependencies
make install-deps

# Run quality checks
make check

# Development workflow
make update-dev        # Build for local testing
./build/stax --version # Test changes

make update           # Install system-wide when ready
```

## Troubleshooting

### Common Issues

**WP Engine Connection Problems:**

```bash
# Verify SSH key configuration
ssh your-username@client-install.ssh.wpengine.net

# Test API credentials
stax wpe list

# Check specific installation access
stax wpe info client-install
```

**⚠️ Known Issue: Warp Terminal SSH Compatibility**

If you're using [Warp terminal](https://www.warp.dev/) and experiencing SSH connection issues with WP Engine sidecars, this is a known compatibility bug. 

**Symptoms:**
- SSH connections to WP Engine hang or fail
- `stax wpe sync` operations timeout during SSH operations
- Manual SSH commands to `*.ssh.wpengine.net` don't work properly

**Workaround:**
Use an alternative terminal for WP Engine SSH operations:
- **macOS**: Terminal.app, iTerm2, or Alacritty
- **Linux**: GNOME Terminal, Konsole, or Alacritty  
- **Cross-platform**: VS Code integrated terminal

**Example:**
```bash
# In alternative terminal (not Warp)
stax wpe sync client-install

# Or run SSH commands directly
ssh your-username@client-install.ssh.wpengine.net
```

This issue is specific to Warp's SSH handling and doesn't affect other Stax functionality.

**Environment Port Conflicts:**

```bash
# Check for conflicts across multiple projects
stax status

# Stop unused environments to free resources
stax stop unused-project-name

# Nuclear option: stop all environments
stax poweroff
```

**Performance Issues with Multiple Projects:**

```bash
# Clean up environments older than one week
stax cleanup --older-than=7d

# Free up Docker resources
docker system prune -f

# Check Docker resource usage
docker stats
```

### Getting Help

- **Command Documentation**: `stax --help` or `stax [command] --help`
- **Firecrown Team**: Check with your project lead or technical team members
- **GitHub Issues**: [Report bugs or request features](https://github.com/Firecrown-Media/stax/issues)
- **Development Questions**: Contact the original developer or current maintainer

## Bug Reports & Feature Requests

### Reporting Issues

Found a bug or have a feature request? Help improve Stax by reporting issues on GitHub:

**🐛 Bug Reports**

When reporting bugs, please include:

1. **Stax version**: Run `stax --version`
2. **Operating System**: macOS version, Linux distribution, etc.
3. **DDEV version**: Run `ddev version`
4. **Steps to reproduce**: Clear, numbered steps
5. **Expected behavior**: What should happen
6. **Actual behavior**: What actually happens
7. **Error messages**: Full error output with `--verbose` flag
8. **Environment details**: PHP version, project type, etc.

**📝 Template for Bug Reports:**

```
## Bug Description
Brief description of the issue

## Environment
- Stax version: [run `stax --version`]
- OS: [macOS 14.5, Ubuntu 22.04, etc.]
- DDEV version: [run `ddev version`]
- PHP version: [if applicable]
- Project type: [WordPress, etc.]

## Steps to Reproduce
1. Run `stax init project-name`
2. Execute `stax wpe sync install-name`
3. Error occurs...

## Expected Behavior
What should happen...

## Actual Behavior
What actually happens...

## Error Output
```
[Include full error output with --verbose flag]
```

## Additional Context
Any other relevant information...
```

**✨ Feature Requests**

For new features or improvements:

1. **Check existing issues** first to avoid duplicates
2. **Describe the use case** - what problem does this solve?
3. **Provide examples** of how the feature would work
4. **Consider impact** on existing workflows
5. **Suggest implementation** if you have technical insights

### Creating GitHub Issues

**Quick Links:**
- [Report a Bug](https://github.com/Firecrown-Media/stax/issues/new?labels=bug&template=bug_report.md)
- [Request a Feature](https://github.com/Firecrown-Media/stax/issues/new?labels=enhancement&template=feature_request.md)
- [View All Issues](https://github.com/Firecrown-Media/stax/issues)

**Step-by-Step Process:**

1. **Visit the GitHub repository**: https://github.com/Firecrown-Media/stax
2. **Click "Issues"** tab
3. **Click "New Issue"** button
4. **Choose appropriate template**:
   - 🐛 Bug Report
   - ✨ Feature Request  
   - 📚 Documentation Issue
   - ❓ Question/Discussion
5. **Fill out the template** with detailed information
6. **Add relevant labels** (bug, enhancement, documentation, etc.)
7. **Submit the issue**

### Issue Labels

Common labels used in the repository:

| Label | Description |
|-------|-------------|
| `bug` | Something isn't working correctly |
| `enhancement` | New feature or improvement request |
| `documentation` | Documentation needs improvement |
| `good first issue` | Good for newcomers to contribute |
| `help wanted` | Extra attention or help needed |
| `question` | General questions about usage |
| `wontfix` | Issue that won't be addressed |
| `duplicate` | Issue already exists |
| `priority-high` | Critical issues affecting functionality |
| `priority-low` | Nice-to-have improvements |

### Contributing to Issue Resolution

**For Firecrown Team Members:**
- Reproduce issues locally when possible
- Provide additional context or workarounds
- Test proposed solutions
- Help triage and label issues appropriately

**For External Contributors:**
- Pull requests welcome for bug fixes
- Discuss major changes in issues first
- Follow existing code style and patterns
- Include tests when appropriate

### Security Issues

**🔒 For security-related issues:**

Do NOT create public GitHub issues for security vulnerabilities. Instead:

1. **Email directly**: security@firecrown.com
2. **Include**: Detailed description and steps to reproduce
3. **Provide**: Your contact information for follow-up
4. **Allow time**: For investigation and resolution before disclosure

### Getting Help Before Creating Issues

**Try these resources first:**

1. **Built-in help**: `stax --help` or `stax [command] --help`
2. **Onboarding guide**: [ONBOARDING.md](ONBOARDING.md)
3. **Search existing issues**: Someone may have already reported it
4. **Ask your team lead**: For Firecrown team members
5. **Check documentation**: This README and related docs

### Issue Response Timeline

**Expected response times:**

- **Critical bugs**: Within 24-48 hours
- **Standard bugs**: Within 1 week
- **Feature requests**: Within 2 weeks
- **Documentation issues**: Within 1 week
- **Questions**: Within 3-5 business days

*Note: Response times may vary based on issue complexity and maintainer availability.*

## Support & Maintenance

### For Firecrown Team Members

- **User Documentation**: This README and built-in help (`stax --help`)
- **Feature Requests**: Submit via GitHub Issues with business justification
- **Bug Reports**: Include environment details and reproduction steps
- **Training**: New team member onboarding should include Stax workflow training

### For External Contributors

- **Open Source**: Stax is available for modification and extension
- **Pull Requests**: Welcome for bug fixes and improvements
- **Custom Integrations**: Hosting provider additions encouraged

---

## Project Background

Stax was developed in 2025 as a custom solution for Firecrown's specific WordPress development workflow challenges. The tool addresses the complexity of managing multiple site projects, each with different hosting environments, WordPress configurations, and development requirements.

The CLI consolidates what previously required multiple tools and manual processes into a single, consistent interface that enforces Firecrown's development standards while providing the flexibility needed for diverse site requirements.

---

_This tool was built specifically for Firecrown's development workflow but designed with extensibility in mind. The hosting provider architecture and configuration system can be adapted for other agencies with similar multi-client WordPress development needs._

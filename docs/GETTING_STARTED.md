# Getting Started with Stax

Welcome to Stax! This guide will help you set up your WordPress development environment in under 10 minutes.

## What is Stax?

Stax is a powerful CLI tool that streamlines WordPress development workflows. It integrates with DDEV for local development and WPEngine for database synchronization, making it easy to work with WordPress sites (single site or multisite).

### What You'll Accomplish

By the end of this guide, you'll have:
- Stax installed and configured
- A local WordPress development environment running
- Database synced from WPEngine (optional)
- Understanding of daily development workflow

**Estimated time:** 10-15 minutes

---

## Prerequisites

Before starting, ensure you have:

1. **Docker Desktop** - Installed and running
   - Download: https://www.docker.com/products/docker-desktop
   - Must be running before using Stax

2. **DDEV** - Local development environment
   ```bash
   brew install ddev
   ```

3. **WPEngine Account** (Optional but recommended)
   - SSH access to your WPEngine install
   - Useful for database synchronization

4. **GitHub Repository Access** (Optional)
   - If working with an existing codebase

---

## Installation

### Install via Homebrew (Recommended)

```bash
# Add the Stax tap
brew install firecrown-media/tap/stax

# Verify installation
stax --version
```

You should see output like:
```
stax version 2.0.0
Git Commit: ...
Build Date: ...
```

### Build from Source (Advanced)

```bash
git clone https://github.com/firecrown-media/stax.git
cd stax
go build -o stax main.go
sudo mv stax /usr/local/bin/
```

---

## Quick Start Workflow

### Step 1: Set Up Credentials (Optional)

If you plan to sync databases from WPEngine, configure your credentials:

```bash
stax setup
```

**What this does:**
- Prompts for WPEngine username and password
- Stores credentials securely in macOS Keychain
- Can be skipped if doing local-only development

**Example interaction:**
```
? WPEngine Username: your-username
? WPEngine Password: ********
✓ Credentials saved to macOS Keychain
```

---

### Step 2: List Available Sites (Optional)

See which WPEngine installations you have access to:

```bash
stax list
```

**Example output:**
```
Available WPEngine Installations:
- mysite-prod (Production)
- mysite-staging (Staging)
- anothersite-prod (Production)
```

This helps you choose which install to sync with your local environment.

---

### Step 3: Initialize Your Project

Navigate to your project directory (or create a new one) and run:

```bash
# Create and navigate to project directory
mkdir my-wordpress-site && cd my-wordpress-site

# Run interactive initialization
stax init
```

**What happens during `stax init`:**
1. Interactive prompts guide you through configuration
2. Creates `.stax.yml` configuration file
3. Sets up DDEV environment (`.ddev/config.yaml`)
4. Optionally clones your GitHub repository
5. Optionally pulls database from WPEngine
6. Starts the environment automatically

**Example prompts and answers:**

```
? Project name: my-wordpress-site
? Project type: wordpress (or wordpress-multisite)
? PHP version: 8.1
? MySQL version: 8.0
? GitHub repository URL: https://github.com/myorg/mysite.git (or leave blank)
? Branch: main
? WPEngine install: mysite-prod (select from list)
? WPEngine environment: production
? Start environment now? Yes
? Pull database now? Yes
```

**After initialization:**
```
✓ Created .stax.yml
✓ Created DDEV configuration
✓ Repository cloned
✓ Environment started
✓ Database pulled and imported

Your site is ready at: https://my-wordpress-site.ddev.site
```

---

### Step 4: Use Your Environment

Now that your environment is running, you can use these commands:

#### Check Status

```bash
stax status
```

Shows comprehensive information:
- Project name and type
- Primary URL
- Container status (web, database, router)
- PHP and database versions
- Xdebug status
- Mailhog URL for testing emails

#### Stop Environment

```bash
# Stop current project
stax stop

# Stop all DDEV projects
stax stop --all
```

#### Start Again

```bash
# Basic start
stax start

# Start with Xdebug enabled
stax start --xdebug

# Start and run build
stax start --build
```

#### Check System Health

```bash
stax doctor
```

Diagnoses issues and provides solutions:
- Docker availability
- DDEV installation
- Configuration validity
- Port conflicts
- Disk space
- Credentials

---

## Common Workflows

### Workflow 1: New Project from WPEngine

Start with a WPEngine site and create local environment:

```bash
# Create project directory
mkdir my-site && cd my-site

# Initialize with WPEngine integration
stax init

# Answer prompts:
# - Project name: my-site
# - WPEngine install: mysite-prod
# - Pull database: Yes
# - Start now: Yes

# Your environment is now running with WPEngine database!
```

---

### Workflow 2: Existing Git Repository

Clone an existing WordPress repository and set up local environment:

```bash
# Create project directory
mkdir my-site && cd my-site

# Initialize with repository
stax init

# Answer prompts:
# - Repository: https://github.com/org/repo.git
# - WPEngine install: mysite-prod
# - Pull database: Yes

# Repository cloned and environment started!
```

---

### Workflow 3: Local Only (No WPEngine)

Set up a local WordPress environment without WPEngine integration:

```bash
# Create project directory
mkdir my-site && cd my-site

# Initialize without WPEngine
stax init

# Answer prompts:
# - Skip WPEngine setup (leave install name blank)
# - Use defaults for DDEV configuration

# Local-only environment created!
```

---

### Workflow 4: Import Existing DDEV Project

Already have a DDEV project? Add Stax features:

```bash
# Navigate to existing DDEV project
cd /path/to/existing/ddev/project

# Import into Stax
stax init --from-ddev

# Optionally add WPEngine integration
# Creates .stax.yml from your DDEV settings

# Now you can use all Stax commands!
```

---

## Daily Development

### Starting Your Day

```bash
# Navigate to your project
cd ~/projects/my-site

# Start the environment
stax start

# Your site opens automatically in browser
# → https://my-site.ddev.site
```

---

### Syncing Database

Pull the latest database from WPEngine:

```bash
# Pull database (creates backup first)
stax db pull

# Pull without backup (faster)
stax db pull --skip-backup
```

**What happens:**
1. Creates local snapshot (unless skipped)
2. Downloads database from WPEngine
3. Imports into local DDEV database
4. Runs search-replace for local URLs
5. Shows completion message

---

### Running Builds

Build your theme and plugin assets:

```bash
# One-time build
stax build

# Development mode with file watching
stax dev

# Development mode for specific theme
stax dev --theme=my-theme
```

---

### Working with WordPress

Run WP-CLI commands:

```bash
# List plugins
stax wp -- plugin list

# Update WordPress
stax wp -- core update

# Run search-replace
stax wp -- search-replace 'oldurl.com' 'my-site.ddev.site'

# Any WP-CLI command
stax wp -- <command>
```

---

### Stopping Environment

When you're done for the day:

```bash
# Stop current project
stax stop

# Stop all DDEV projects
stax stop --all

# Stop and remove data (destructive)
stax stop --remove-data
```

---

## Troubleshooting

### "stax: command not found"

**Problem:** Stax is not installed or not in PATH

**Solution:**
```bash
# Install via Homebrew
brew install firecrown-media/tap/stax

# Or add to PATH if built from source
export PATH="/usr/local/bin:$PATH"
```

---

### "Docker is not running"

**Problem:** Docker Desktop is not started

**Solution:**
1. Open Docker Desktop application
2. Wait for it to fully start (whale icon stops animating)
3. Run `stax doctor` to verify
4. Try your command again

---

### "No project configuration found"

**Problem:** Neither `.stax.yml` nor `.ddev/config.yaml` exists

**Solution:**
```bash
# Initialize the project
stax init

# Or if you just want DDEV
ddev config --project-type=wordpress
stax start
```

---

### "Port already in use"

**Problem:** Required ports (80, 443, 3306, etc.) are in use

**Solution:**
```bash
# Run diagnostics
stax doctor

# Shows which ports are in use
# Stop conflicting services:
# - Apache: sudo apachectl stop
# - MySQL: brew services stop mysql
# - Other DDEV projects: stax stop --all

# Or change ports in .ddev/config.yaml
```

---

### "Failed to pull database"

**Problem:** WPEngine credentials invalid or network issue

**Solution:**
```bash
# Reconfigure credentials
stax setup

# Verify credentials are correct
stax list

# Check WPEngine install name
# Make sure you have SSH access to the install
```

---

### "DDEV start failed"

**Problem:** Various DDEV issues

**Solution:**
```bash
# Run diagnostics
stax doctor

# Common fixes:
# 1. Restart Docker Desktop
# 2. Clean up DDEV
ddev poweroff
ddev clean

# 3. Check disk space
df -h

# 4. Try starting again
stax start
```

---

### Configuration Issues

**Problem:** Something wrong with `.stax.yml` or DDEV config

**Solution:**
```bash
# Validate your configuration
stax validate

# Shows what's wrong and how to fix it

# If config is broken, regenerate:
stax init
```

---

## Next Steps

Now that you have Stax set up, explore these resources:

### Documentation
- [User Guide](USER_GUIDE.md) - Comprehensive feature documentation
- [Multisite Guide](MULTISITE.md) - WordPress multisite setup
- [Development Guide](DEVELOPMENT.md) - Contributing to Stax
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues and solutions

### Advanced Features
- Database snapshots and restoration
- Remote media proxying (BunnyCDN + WPEngine)
- Custom build configurations
- Team configuration sharing
- Multiple environment support

### Command Reference

Get help anytime:
```bash
# General help
stax --help

# Command-specific help
stax <command> --help

# Examples:
stax init --help
stax start --help
stax db --help
```

---

## Getting Help

### Run Diagnostics

The `doctor` command solves most issues:
```bash
stax doctor
```

Provides:
- System health checks
- Configuration validation
- Actionable solutions
- Links to documentation

### Check Documentation

All commands have built-in help:
```bash
stax <command> --help
```

### Report Issues

Found a bug or have a suggestion?
- GitHub Issues: https://github.com/firecrown-media/stax/issues
- Include output from `stax doctor`
- Describe what you expected vs. what happened

---

## Tips for Success

1. **Run `stax doctor` regularly** - Catches issues early
2. **Use `stax validate`** - After manual config changes
3. **Keep Docker running** - Required for all operations
4. **Back up before major changes** - `stax db snapshot`
5. **Use `--verbose` for debugging** - `stax start --verbose`
6. **Stop projects when not in use** - `stax stop --all` (frees resources)

---

## Quick Reference

```bash
# Setup
stax setup              # Configure credentials
stax init               # Initialize project
stax init --from-ddev   # Import existing DDEV project

# Daily Use
stax start              # Start environment
stax stop               # Stop environment
stax status             # Check status
stax db pull            # Sync database

# Development
stax build              # Build assets
stax dev                # Dev mode with watching
stax wp -- <command>    # Run WP-CLI

# Troubleshooting
stax doctor             # Diagnose issues
stax validate           # Validate config
stax --help             # Get help
```

---

**You're all set!** 🎉

Your WordPress development environment is ready. Happy coding!

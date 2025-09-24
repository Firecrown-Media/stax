# Stax User Guide

A comprehensive guide for using Stax in your daily WordPress development workflow.

## Table of Contents
- [Core Concepts](#core-concepts)
- [Project Workflows](#project-workflows)
- [Command Reference](#command-reference)
- [Real-World Scenarios](#real-world-scenarios)
- [Best Practices](#best-practices)
- [Tips and Tricks](#tips-and-tricks)

## Core Concepts

### Understanding the Stack

Before diving into commands, let's understand what Stax manages:

```
Your Computer
    └── Stax (CLI Tool)
         └── DDEV (Container Manager)
              └── Docker (Container Runtime)
                   ├── Web Container (nginx/apache + PHP)
                   ├── Database Container (MySQL/MariaDB)
                   └── Additional Services (mailhog, phpmyadmin, etc.)
```

### Key Terms

| Term | What It Means | Example |
|------|---------------|---------|
| **Project** | A local WordPress site | `client-website` |
| **Environment** | The container setup for a project | PHP 8.2 + MySQL 8.0 |
| **Install** | A WP Engine site name | `clientprod` |
| **Sync** | Copy data from remote to local | Database + files |
| **Hot Swap** | Change versions without data loss | PHP 8.1 → 8.3 |

## Project Workflows

### Creating Your First Project

Let's create a fresh WordPress installation:

```bash
# 1. Initialize the project
stax init my-blog
# Output: ✓ Project 'my-blog' initialized successfully!

# 2. Set up WordPress
stax setup my-blog --install-wp
# Output: ✓ Starting DDEV environment...
#         ✓ Downloading WordPress...
#         ✓ Creating database...
#         ✓ WordPress installed successfully!

# 3. Start the environment
stax start my-blog
# Output: ✓ Starting containers...
#         ✓ Site ready at: https://my-blog.ddev.site

# 4. Access WordPress admin
open https://my-blog.ddev.site/wp-admin
# Default credentials: admin / admin
```

### Working with Existing Projects

When you return to work on a project:

```bash
# Navigate to project directory
cd ~/projects/client-site

# Start the environment
stax start
# Stax detects you're in a project directory

# Check status
stax status
# Shows: client-site (running) - https://client-site.ddev.site

# When done for the day
stax stop
```

### Cloning from WP Engine

Sync a production site for local development:

```bash
# 1. Create local project with matching name
stax init acme-corp

# 2. Sync from WP Engine (database + files)
stax wpe sync acmecorpprod
# Output: 📥 Downloading database...
#         📥 Downloading files...
#         ✓ Site synced successfully!

# 3. For faster sync (database only)
stax wpe sync acmecorpprod --skip-files

# 4. Start working
stax start acme-corp
```

### Hot Swapping Environments

Test compatibility across PHP versions:

```bash
# Check current configuration
stax swap status
# Shows: PHP 8.2, MySQL 8.0

# Switch to newer PHP
stax swap php 8.3
# Output: 📸 Backing up configuration...
#         🔄 Updating PHP to 8.3...
#         ✓ Environment updated successfully!

# Test with legacy environment
stax swap preset legacy
# Switches to: PHP 7.4, MySQL 5.7

# Something broke? Rollback!
stax swap --rollback
# Returns to previous configuration
```

## Command Reference

### Essential Daily Commands

These are the commands you'll use most often:

```bash
# Project Management
stax init [name]              # Create new project
stax start [name]             # Start development
stax stop [name]              # Stop when done
stax status                   # What's running?
stax delete [name] --yes      # Remove project

# WordPress Operations
stax wp [command]             # Run WP-CLI commands
stax wp plugin list           # List plugins
stax wp user create           # Create user
stax wp db export backup.sql  # Export database

# Environment Management
stax swap list                # Show available versions
stax swap php 8.3            # Change PHP version
stax swap preset modern      # Apply preset configuration
stax swap --rollback         # Undo last change
```

### WP Engine Commands

```bash
# Listing and Info
stax wpe list                 # Show all WP Engine installs
stax wpe info [install]       # Show install details

# Syncing
stax wpe sync [install]                    # Full sync (db + files)
stax wpe sync [install] --skip-files       # Database only
stax wpe sync [install] --skip-database    # Files only
stax wpe sync [install] --environment=staging  # From staging

# Database Operations
stax wpe db download [install]             # Download database
stax wpe db import database.sql            # Import database
stax wpe db analyze                        # Check media URLs
```

### Advanced Commands

```bash
# Multiple environments
stax start project1 project2   # Start multiple projects

# Custom configurations
stax init mysite --php-version=8.1 --database=mysql:5.7

# Debugging
stax --verbose [command]       # Show detailed output
stax describe [project]        # Show project details
```

## Real-World Scenarios

### Scenario 1: Monday Morning Bug Fix

Your client reports a critical bug on their live site:

```bash
# 1. Create fresh local copy
stax init clientsite-bugfix

# 2. Sync latest production data
stax wpe sync clientprod --skip-files
# Files aren't needed for debugging

# 3. Start and reproduce the issue
stax start clientsite-bugfix
open https://clientsite-bugfix.ddev.site

# 4. Enable debugging
stax wp config set WP_DEBUG true --raw
stax wp plugin install query-monitor --activate

# 5. Fix the issue, test thoroughly

# 6. When done, clean up
stax stop clientsite-bugfix
stax delete clientsite-bugfix --yes
```

### Scenario 2: Plugin Compatibility Testing

Testing if a plugin works with newer PHP:

```bash
# 1. Set up test environment
stax init plugin-test
stax wpe sync clientsite --skip-files

# 2. Document current state
stax swap status > compatibility-report.txt
stax wp plugin list --status=active >> compatibility-report.txt

# 3. Test with current PHP
stax start plugin-test
stax wp plugin install new-plugin --activate
# Test functionality...

# 4. Test with newer PHP
stax swap php 8.3
stax start plugin-test
# Test again...

# 5. Test with bleeding edge
stax swap preset bleeding-edge
# Test once more...

# 6. Report findings
echo "Plugin compatible with PHP 8.2, 8.3, and 8.4" >> compatibility-report.txt
```

### Scenario 3: New Developer Onboarding

Setting up a new team member:

```bash
# 1. They install dependencies
brew install docker ddev/ddev/ddev
brew tap firecrown-media/stax
brew install stax

# 2. Clone team configuration
git clone git@github.com:company/wp-configs.git
cp wp-configs/.stax.yaml ~/

# 3. Set up main project
stax init main-project
stax wpe sync mainprod

# 4. Verify setup
stax status
stax wp user create newdev newdev@company.com --role=administrator

# New developer ready to work!
```

### Scenario 4: Performance Optimization

Testing site performance with different configurations:

```bash
# 1. Baseline test
stax swap preset stable
stax start performance-test
# Run performance tests, note results

# 2. Test with performance preset
stax swap preset performance
stax start performance-test
# Run same tests, compare

# 3. Test with specific optimizations
stax swap php 8.3
stax wp plugin deactivate --all
stax wp theme activate twentytwentythree
# Test minimal setup

# 4. Document optimal configuration
stax swap status > optimal-config.txt
```

## Best Practices

### 1. Project Organization

```bash
# Recommended directory structure
~/projects/
  ├── client-a/
  │   ├── project1/
  │   └── project2/
  ├── client-b/
  │   └── main-site/
  └── personal/
      └── my-blog/
```

### 2. Naming Conventions

```bash
# Use consistent, descriptive names
stax init clientname-purpose
# Examples:
stax init acme-redesign
stax init acme-staging-copy
stax init acme-feature-test
```

### 3. Resource Management

```bash
# Stop unused projects to free resources
stax status  # Check what's running
stax stop project1 project2  # Stop specific projects
stax poweroff  # Stop everything

# Clean up old projects
stax delete old-project --yes
docker system prune -a  # Clean Docker
```

### 4. Configuration Management

```bash
# Use global config for defaults
cat > ~/.stax.yaml << EOF
default_php_version: "8.2"
wordpress:
  admin_email: "dev@company.com"
EOF

# Override per-project as needed
cat > ./stax.yaml << EOF
php_version: "8.1"  # This project needs older PHP
EOF
```

## Tips and Tricks

### Speed Optimizations

```bash
# Skip files for faster syncing
stax wpe sync install --skip-files

# Use media redirect instead of downloading
stax wpe sync install --skip-files --create-upload-redirect

# Run multiple operations in parallel
stax start project1 & stax start project2 &
```

### Debugging Helpers

```bash
# Enable verbose output
stax --verbose wpe sync install

# Check DDEV directly
ddev describe
ddev logs

# Access containers
ddev ssh  # Web container
ddev mysql  # Database CLI
```

### Shortcuts and Aliases

Add to your `~/.zshrc` or `~/.bashrc`:

```bash
# Stax shortcuts
alias si="stax init"
alias ss="stax start"
alias st="stax stop"
alias sstat="stax status"
alias swp="stax wp"

# Common workflows
alias sync-prod="stax wpe sync $1 --skip-files"
alias quick-test="stax init test-$RANDOM && stax setup test-$RANDOM --install-wp"
```

### Database Management

```bash
# Export database before major changes
stax wp db export backup-$(date +%Y%m%d).sql

# Quick database reset
stax wp db reset --yes
stax wpe db import production.sql

# Search and replace
stax wp search-replace "old-domain.com" "new-domain.com"
```

## Common Workflows Summary

| Task | Commands |
|------|----------|
| **New project** | `stax init` → `stax setup --install-wp` → `stax start` |
| **Clone from production** | `stax init` → `stax wpe sync` → `stax start` |
| **Daily development** | `stax start` → work → `stax stop` |
| **Test compatibility** | `stax swap php 8.3` → test → `stax swap --rollback` |
| **Clean up** | `stax stop` → `stax delete --yes` |

## Next Steps

- Learn about [WP Engine Integration](WPENGINE.md)
- Explore [Troubleshooting](TROUBLESHOOTING.md)
- Set up [WordPress Multisite](MULTISITE.md)
- Contribute to [Development](DEVELOPMENT.md)

## Getting Help

- Run `stax --help` or `stax [command] --help`
- Check [Troubleshooting](TROUBLESHOOTING.md)
- Report issues at [GitHub](https://github.com/Firecrown-Media/stax/issues)
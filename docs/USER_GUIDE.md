# Stax User Guide

A comprehensive guide to using Stax for WordPress development.

---

## Table of Contents

- [Introduction](#introduction)
- [Single Site or Multisite?](#single-site-or-multisite)
- [Discovering WPEngine Installs](#discovering-wpengine-installs)
- [Daily Workflows](#daily-workflows)
- [Database Management](#database-management)
- [WordPress Operations](#wordpress-operations)
- [Build and Development](#build-and-development)
- [Configuration Management](#configuration-management)
- [Team Collaboration](#team-collaboration)
- [Advanced Workflows](#advanced-workflows)
- [Best Practices](#best-practices)
- [Tips and Tricks](#tips-and-tricks)

---

## Introduction

This guide covers everything you need to know to use Stax effectively. If you haven't installed Stax yet, start with the [Installation Guide](./INSTALLATION.md) and [Quick Start](./QUICK_START.md).

### Prerequisites

Before following this guide, make sure you've:
- [x] Installed Stax and prerequisites
- [x] Configured credentials with `stax setup`
- [x] Initialized at least one project with `stax init`

---

## Single Site or Multisite?

Stax works with both standard WordPress installations and multisite networks. Choose what's right for your project:

### Single WordPress Site (Default)

Most WordPress projects are single-site installations. Stax makes it easy to develop them locally:

```yaml
# .stax.yml
project:
  name: my-site
  type: wordpress  # Single site

wordpress:
  domain: mysite.local
```

**Perfect for:**
- Client websites
- Blogs
- Marketing sites
- E-commerce sites
- Most WordPress projects

### WordPress Multisite (Optional)

If you're working with a WordPress multisite network, Stax has first-class support:

```yaml
# .stax.yml
project:
  name: my-network
  type: wordpress-multisite
  mode: subdomain  # or subdirectory

network:
  domain: mynetwork.local
  sites:
    - name: site1
      domain: site1.mynetwork.local
```

**Use multisite when:**
- You have multiple sites sharing code and database
- You're working with subdomain or subdirectory networks
- You need network-level administration

**Not sure?** If you're asking whether you need multisite, you probably don't. Most WordPress projects are single sites.

See [MULTISITE.md](./MULTISITE.md) for detailed multisite documentation.

---

## Discovering WPEngine Installs

Before you can initialize a Stax project, you need to know your WPEngine install name. The `stax list` command helps you discover available installs without needing a stax.yml file first.

### Listing All Installs

**See all WPEngine installs available to your account**:

```bash
stax list
```

**Expected output**:
```
Listing WPEngine Installs

INSTALL NAME           ENVIRONMENT   PRIMARY DOMAIN              PHP   STATUS
myinstall              production    mysite.wpengine.com         8.1   active
myinstall-staging      staging       myinstall-staging.wpe       8.1   active
client-site-prod       production    clientsite.com              8.2   active
client-site-staging    staging       clientsite-staging.wpe      8.2   active

Total: 4 installs
```

This shows:
- **Install Name**: The identifier you'll use in `stax init`
- **Environment**: Whether it's production or staging
- **Primary Domain**: The main domain for this install
- **PHP**: PHP version running on WPEngine
- **Status**: Usually "active"

### Filtering Installs

**Filter by name using regex**:
```bash
# Find all client installs
stax list --filter="client.*"

# Find production installs only
stax list --filter=".*-prod"

# Find installs starting with "fs"
stax list --filter="^fs-"
```

**Filter by environment**:
```bash
# Show only production installs
stax list --environment=production

# Show only staging installs
stax list --environment=staging
```

**Combine filters**:
```bash
# Client staging installs only
stax list --filter="client.*" --environment=staging
```

### Different Output Formats

**JSON output** (for scripting):
```bash
stax list --output=json
```

**YAML output**:
```bash
stax list --output=yaml
```

**Table output** (default, human-readable):
```bash
stax list --output=table
# or just:
stax list
```

### Common Use Cases

**Onboarding a new team member**:
```bash
# List all installs to find the one you need
stax list

# Found it! Now initialize
stax init
# Enter "myinstall" when prompted for install name
```

**Finding the staging environment**:
```bash
# List staging installs only
stax list --environment=staging

# Use the install name in your project
stax config set wpengine.install myinstall-staging
```

**Auditing available installs**:
```bash
# Get complete list in JSON
stax list --output=json > wpengine-installs.json

# Process with jq
stax list --output=json | jq '.[] | select(.environment == "production")'
```

### Troubleshooting List Command

**"WPEngine credentials not found"**:
```bash
# Configure your credentials first
stax setup
```

**"Authentication failed"**:
- Check your WPEngine API credentials
- Verify you have the correct username/password
- Try `stax setup` to reconfigure

**"No installs found"**:
- Verify your WPEngine account has install access
- Check if your API user has the correct permissions
- Contact WPEngine support if needed

**Notes**:
- This command doesn't require a stax.yml file
- No SSH key needed (API only)
- Fast operation (typically <5 seconds)
- Safe to run anytime - read-only operation

### What You'll Learn

- How to use Stax in your daily development workflow
- How to manage databases and snapshots
- How to work with WordPress multisite
- How to configure and customize Stax
- Best practices for team collaboration

---

## Daily Workflows

### Starting Your Day

**Monday morning, you're ready to code**:

```bash
# Navigate to your project
cd ~/Sites/my-project

# Start your environment
stax start
```

**What happens**:
- Docker containers start (web, database, router)
- Takes ~10 seconds if containers exist
- First time may take 1-2 minutes (downloading images)

**Expected output**:
```
🚀 Starting my-project

✓ Starting DDEV containers
  Web: https://my-project.local
  Database: MySQL 8.0 (ready)
  MailHog: http://my-project.local:8025

Environment started successfully!
```

Your site is now accessible at https://my-project.local.

### Checking Environment Status

**See what's running**:

```bash
stax status
```

**Output**:
```
📊 Status: my-project

Environment: Running ✓

Containers:
  ✓ ddev-my-project-web       (healthy)
  ✓ ddev-my-project-db        (healthy)
  ✓ ddev-router               (healthy)

URLs:
  Network:    https://my-project.local
  Site 1:     https://site1.my-project.local
  Site 2:     https://site2.my-project.local
  MailHog:    http://my-project.local:8025

Configuration:
  PHP:        8.1
  MySQL:      8.0
  Webserver:  nginx

Database:
  Size:       245 MB
  Tables:     127
  Sites:      3

Last database pull: 2 hours ago
```

This tells you everything you need to know about your environment.

### Making Code Changes

**Your typical workflow**:

1. **Start the dev watcher** (optional but recommended):
   ```bash
   stax dev
   ```
   This runs your build tools in watch mode - when you save a file, assets rebuild automatically.

2. **Make your changes**:
   - Edit theme files
   - Modify plugins
   - Update configurations
   - Changes are reflected immediately (PHP) or after rebuild (CSS/JS)

3. **View your changes**:
   - Open https://my-project.local in your browser
   - Hard refresh if needed (Cmd+Shift+R)

4. **Test your changes**:
   ```bash
   # Clear WordPress caches
   stax wp cache flush

   # Run linters
   stax lint

   # Run tests (if configured)
   stax ssh
   npm test
   composer test
   exit
   ```

### Working with Git

**Stax and Git work together seamlessly**:

```bash
# Create a feature branch
git checkout -b feature/new-header

# Make changes...
# Edit files, commit as normal

# Rebuild if needed
stax build

# Test your changes
open https://my-project.local

# Push when ready
git push origin feature/new-header
```

**Important**: The `.stax.yml` file should be committed to Git so your team shares the same configuration.

### Ending Your Day

**When you're done for the day**:

```bash
stax stop
```

**What happens**:
- Containers stop gracefully
- Database data is preserved
- Your code remains unchanged
- Takes ~5 seconds

**Expected output**:
```
🛑 Stopping my-project

✓ Stopping DDEV containers

Environment stopped.
```

**Should I stop every night?**
- Yes, if you want to free up system resources
- No, if you want instant access tomorrow (containers can stay running)
- Docker Desktop uses ~2-4GB RAM when running

---

## Database Management

### Pulling from WPEngine

**Get the latest database from production**:

```bash
stax db pull
```

**What happens**:
1. Creates automatic snapshot (named like `auto_2024-11-08_14-30-00`)
2. Connects to WPEngine via SSH
3. Detects table prefix
4. Exports database (excluding logs, transients, spam by default)
5. Transfers to your machine
6. Imports to local database
7. Runs search-replace for all domains
8. Flushes WordPress caches

**Expected output**:
```
🗄️  Pulling database from WPEngine (production)

✓ Creating snapshot: auto_2024-11-08_14-30-00
✓ Connecting to WPEngine SSH Gateway
✓ Detecting table prefix: wp_
✓ Exporting database
  Tables: 127
  Size: 245 MB
✓ Transferring database
  [████████████████████] 245 MB
✓ Importing to local database
  Rows: 1,245,678
✓ Running search-replace
  Network: mysite.wpengine.com → my-project.local (43 replacements)
  site1: site1.com → site1.my-project.local (1,234 replacements)
  site2: site2.com → site2.my-project.local (987 replacements)
✓ Flushing WordPress cache

Database pulled successfully!
Time: 2m 34s
```

**Pull from staging instead**:
```bash
stax db pull --environment=staging
```

**Skip automatic snapshot**:
```bash
stax db pull --snapshot=false
```

**Pull specific tables only**:
```bash
stax db pull --exclude-tables=wp_actionscheduler_logs,wp_wc_admin_notes
```

**Sanitize sensitive data** (anonymize user emails/passwords):
```bash
stax db pull --sanitize
```

### Database Snapshots

Snapshots are local backups of your database. They're fast to create and restore.

**Create a snapshot**:
```bash
# Auto-named snapshot
stax db snapshot

# Named snapshot
stax db snapshot before-testing

# With description
stax db snapshot pre-migration --description="Before user table changes"
```

**List all snapshots**:
```bash
stax db list
```

**Output**:
```
📋 Database Snapshots

Name                         Created        Size    Description
auto_2024-11-08_14-30-00    2 hours ago    245 MB  Auto (before pull)
before-testing              1 day ago      243 MB  Manual snapshot
pre-migration               2 days ago     240 MB  Before user table changes

Total: 3 snapshots (728 MB)
```

**Restore a snapshot**:
```bash
stax db restore before-testing
```

You'll be asked to confirm:
```
♻️  Restoring snapshot: before-testing

⚠️  This will replace your current database. Continue? (yes/no): yes

✓ Creating backup of current database
✓ Restoring snapshot
  Size: 243 MB
  Rows: 1,240,000

Database restored successfully!
```

**Delete old snapshots**:
```bash
stax db delete-snapshot auto_2024-11-01_10-00-00
```

### Exporting and Importing

**Export your local database**:
```bash
# Export to default location
stax db export

# Export to specific file
stax db export ~/backups/my-backup.sql

# Export compressed
stax db export --gzip
```

**Import a SQL file**:
```bash
# Import with automatic snapshot
stax db import ~/Downloads/database.sql

# Import without snapshot (not recommended)
stax db import database.sql --snapshot=false

# Import and run search-replace
stax db import database.sql --replace
```

### Database Queries

**Run a SQL query**:
```bash
# Simple query
stax db query "SELECT * FROM wp_options WHERE option_name = 'siteurl'"

# Query from file
stax db query --file=query.sql

# JSON output
stax db query "SELECT * FROM wp_users LIMIT 5" --format=json
```

**Access MySQL directly**:
```bash
stax ssh
mysql
# You're now in MySQL prompt
SHOW DATABASES;
USE db;
SHOW TABLES;
exit;
exit;
```

Or use phpMyAdmin:
```bash
ddev launch -p
```

---

## WordPress Operations

### Using WP-CLI

Stax provides direct access to WP-CLI:

```bash
# General syntax
stax wp <command> [args] [flags]
```

**Common commands**:

```bash
# List plugins
stax wp plugin list

# Activate a plugin
stax wp plugin activate wordpress-seo

# Deactivate a plugin
stax wp plugin deactivate akismet

# Update all plugins
stax wp plugin update --all

# List themes
stax wp theme list

# Activate a theme
stax wp theme activate twentytwentyfour

# List users
stax wp user list

# Create a user
stax wp user create johndoe john@example.com --role=editor --user_pass=password

# List all sites (multisite)
stax wp site list

# Flush caches
stax wp cache flush

# Rewrite flush
stax wp rewrite flush

# Check core integrity
stax wp core verify-checksums

# Update WordPress core
stax wp core update

# Run cron
stax wp cron event run --due-now
```

### Search and Replace

**Replace URLs across your database**:

```bash
# Dry run (see what would change)
stax wp search-replace 'old.com' 'new.com' --dry-run

# Actually replace
stax wp search-replace 'old.com' 'new.com'

# Network-wide replace
stax wp search-replace 'old.com' 'new.com' --network

# Specific site only
stax wp search-replace 'old.com' 'new.com' --url=site1.my-project.local

# Skip specific columns
stax wp search-replace 'old.com' 'new.com' --skip-columns=guid
```

**Why skip GUID?**
WordPress uses GUID as a permanent identifier. Changing it can break things.

### Managing Plugins

**Install a plugin**:
```bash
# From WordPress.org
stax wp plugin install wordpress-seo --activate

# Specific version
stax wp plugin install wordpress-seo --version=21.0 --activate

# From URL
stax wp plugin install https://example.com/plugin.zip --activate
```

**Network activate plugins**:
```bash
stax wp plugin activate wordpress-seo --network
```

**Deactivate and uninstall**:
```bash
# Deactivate
stax wp plugin deactivate wordpress-seo

# Uninstall (deletes plugin data)
stax wp plugin uninstall wordpress-seo

# Delete (removes files only)
stax wp plugin delete wordpress-seo
```

### Managing Themes

```bash
# Install theme
stax wp theme install twentytwentyfour

# Activate theme
stax wp theme activate twentytwentyfour

# Delete inactive themes
stax wp theme delete twentytwentythree
```

### Managing Multisite

**List all sites**:
```bash
stax wp site list
```

**Output**:
```
+----+-----------+---------------------------+
| ID | url       | state                     |
+----+-----------+---------------------------+
| 1  | my-project.local | 1                         |
| 2  | site1.my-project.local | 1                         |
| 3  | site2.my-project.local | 1                         |
+----+-----------+---------------------------+
```

**Create a new site**:
```bash
stax wp site create --slug=site3 --title="Site 3"
```

**Delete a site**:
```bash
# Careful! This deletes all site content
stax wp site delete 3 --yes
```

**Empty a site** (delete posts/pages but keep site):
```bash
stax wp site empty 3 --yes
```

**More multisite operations**: See [MULTISITE.md](./MULTISITE.md)

---

## Build and Development

### Running Builds

**Run a one-time build**:
```bash
stax build
```

This runs:
1. `composer install`
2. `npm install`
3. Your build script (usually `scripts/build.sh` or `npm run build`)

**What gets built**:
- SCSS → CSS
- JavaScript (transpiling, bundling)
- Image optimization
- Asset copying
- Anything in your build scripts

### Development Mode

**Run build in watch mode**:
```bash
stax dev
```

This:
- Runs your build tools in watch mode
- Rebuilds automatically when files change
- Usually faster than full rebuilds
- Logs output to console

**Typical watch mode output**:
```
🔧 Starting development mode

✓ Starting file watcher
  Watching: wp-content/themes/*/assets/**

  [12:34:56] File changed: theme.scss
  [12:34:57] Compiled: theme.css
  [12:34:57] ✓ Build complete (234ms)
```

Press `Ctrl+C` to stop.

### Linting

**Run code quality checks**:
```bash
stax lint
```

This typically runs:
- **PHP_CodeSniffer** (WordPress coding standards)
- **ESLint** (JavaScript)
- **Stylelint** (CSS/SCSS)
- **PHPStan** (PHP static analysis)

**Output**:
```
🔍 Running linters

✓ PHP_CodeSniffer
  Checked 45 files
  0 errors, 3 warnings

✗ ESLint
  Checked 23 files
  2 errors, 1 warning

  src/js/main.js
    12:5  error  'console' is not defined  no-undef

✓ Stylelint
  Checked 15 files
  0 errors

✗ Linting failed (2 errors)
```

**Fix automatically** (when possible):
```bash
stax lint --fix
```

### SSH Access

**Access the container**:
```bash
stax ssh
```

You're now inside the container at `/var/www/html`.

**Things you can do**:
```bash
# Run WP-CLI
wp plugin list

# Run Composer
composer install
composer update

# Run npm
npm install
npm run build

# Access MySQL
mysql

# View files
ls -la
cat wp-config.php

# Exit container
exit
```

**Run a single command** without entering the container:
```bash
stax ssh "wp plugin list"
stax ssh "composer install"
stax ssh "npm run build"
```

### Viewing Logs

**View container logs**:
```bash
# Last 100 lines
stax logs

# Follow logs (live tail)
stax logs -f

# Last 500 lines
stax logs --tail=500

# Specific service
stax logs --service=web

# With timestamps
stax logs -f --timestamp
```

**Log types you'll see**:
- PHP errors and warnings
- nginx access logs
- MySQL slow query logs
- WordPress debug logs
- Build script output

---

## Configuration Management

### Viewing Configuration

**See your current configuration**:
```bash
stax config list
```

**Output**:
```yaml
Project Configuration (.stax.yml):

project:
  name: my-project
  type: wordpress-multisite
  mode: subdomain

wpengine:
  install: myinstall
  environment: production

network:
  domain: my-project.local
  sites:
    - name: site1
      domain: site1.my-project.local

ddev:
  php_version: "8.1"
  mysql_version: "8.0"
```

**Get a specific value**:
```bash
stax config get wpengine.environment
# Output: production

stax config get ddev.php_version
# Output: 8.1
```

### Updating Configuration

**Set a configuration value**:
```bash
# Switch to staging
stax config set wpengine.environment staging

# Change PHP version
stax config set ddev.php_version 8.2

# Update project mode
stax config set project.mode subdirectory
```

**After changing config**:
```bash
# Restart to apply changes
stax restart
```

Some changes (like PHP version) require regenerating DDEV config.

### Global vs Project Config

**Project config** (`.stax.yml`):
- Lives in your project directory
- Committed to Git
- Shared by whole team
- Project-specific settings

**Global config** (`~/.stax/config.yml`):
- Lives in `~/.stax/config.yml`
- Personal to you
- Not in Git
- Defaults for all projects

**Edit global config**:
```bash
code ~/.stax/config.yml
```

**Example global config**:
```yaml
defaults:
  ddev:
    php_version: "8.1"
    mysql_version: "8.0"

  wpengine:
    environment: staging  # Default to staging for safety

  build:
    auto_build: true
    watch_enabled: false
```

### Validating Configuration

**Check if your config is valid**:
```bash
stax config validate
```

**Output if valid**:
```
✓ Configuration is valid

Warnings:
  - PHP 8.1 is older than WPEngine production (8.2)
    Consider upgrading: stax config set ddev.php_version 8.2
```

**Output if invalid**:
```
✗ Configuration is invalid

Errors:
  1. wpengine.install is required
  2. network.domain must end in .local or .ddev.site
  3. ddev.php_version must be "7.4", "8.0", "8.1", or "8.2"

Fix these errors in .stax.yml
```

---

## Team Collaboration

### Sharing Configuration

**Your `.stax.yml` should be in Git**:

```bash
git add .stax.yml
git commit -m "Add Stax configuration"
git push
```

Now your whole team uses identical settings.

**Team member clones the repo**:
```bash
git clone https://github.com/mycompany/my-project.git
cd my-project

# Configure their credentials (one-time)
stax setup

# Initialize from existing config
stax init
```

Stax reads `.stax.yml` and sets up an identical environment.

### Handling Different Preferences

**Team members can have different preferences**:

Person A likes staging data:
```bash
# In ~/.stax/config.yml
defaults:
  wpengine:
    environment: staging
```

Person B likes production data:
```bash
# In ~/.stax/config.yml
defaults:
  wpengine:
    environment: production
```

The project `.stax.yml` remains the same, but each person's pulls come from their preferred environment.

### Multiple Projects

**Working on multiple projects**:

```bash
# Project 1
cd ~/Sites/project-1
stax start
# Work on project 1

# Switch to Project 2
cd ~/Sites/project-2
stax start
# Work on project 2
```

All projects can run simultaneously. The `ddev-router` container routes requests to the right project based on domain.

**Stop all projects**:
```bash
ddev poweroff
```

**List all DDEV projects**:
```bash
ddev list
```

---

## Advanced Workflows

### Testing with Production Data

```bash
# 1. Create a snapshot
stax db snapshot before-testing

# 2. Pull production database
stax db pull --environment=production

# 3. Test your changes
# ... make changes, test, verify ...

# 4. Something broke? Restore snapshot
stax db restore before-testing

# 5. Or start fresh
stax db pull --environment=staging
```

### Database Migration Testing

```bash
# 1. Pull production data
stax db pull

# 2. Create snapshot
stax db snapshot before-migration

# 3. Run migration
stax ssh
wp db query < migration.sql
exit

# 4. Verify migration
stax wp db query "SELECT COUNT(*) FROM new_table"

# 5. If good, commit migration script
git add migration.sql
git commit -m "Add user migration"

# 6. If bad, rollback
stax db restore before-migration
```

### Environment Switching

**Switch between staging and production**:

```bash
# Currently using production
stax config get wpengine.environment
# Output: production

# Switch to staging
stax config set wpengine.environment staging

# Pull staging database
stax db pull
# Now working with staging data

# Switch back to production
stax config set wpengine.environment production
stax db pull
```

### Remote Media Configuration

**Stax proxies media from production** - you don't need to download gigabytes of images.

**How it works**:
1. WordPress requests `/wp-content/uploads/2024/01/image.jpg`
2. File doesn't exist locally
3. nginx proxies request to production CDN/server
4. Image loads in your browser
5. Optionally cached locally for faster subsequent loads

**Configure remote media**:

Edit `.stax.yml`:
```yaml
media:
  proxy:
    enabled: true
    remote_url: https://cdn.example.com
    cache_locally: true
    cache_duration: 7d
```

Restart to apply:
```bash
stax restart
```

**Disable proxying** (use local files only):
```yaml
media:
  proxy:
    enabled: false
```

### Custom Build Scripts

**Customize the build process**:

Edit `.stax.yml`:
```yaml
build:
  pre_install:
    - echo "Starting build..."

  install:
    - composer install --no-dev
    - npm ci

  post_install:
    - npm run build:production
    - scripts/optimize-images.sh

  watch:
    enabled: true
    command: npm run watch
    paths:
      - wp-content/themes/*/assets/**
```

**Custom build script example** (`scripts/build.sh`):

```bash
#!/bin/bash
set -e

echo "Building theme assets..."

# Compile SCSS
npm run sass

# Bundle JavaScript
npm run webpack

# Optimize images
npm run imagemin

# Copy static assets
npm run copy

echo "Build complete!"
```

Make it executable:
```bash
chmod +x scripts/build.sh
```

---

## Best Practices

### When to Pull from Production

**Good times to pull**:
- Monday morning (start week with fresh data)
- Before testing features that depend on real data
- After major production database changes
- When debugging production-specific issues

**Don't pull too often**:
- Pulls can take 2-5 minutes
- You'll lose local changes to database
- Creates large snapshots

**Recommendation**: Once a day or less for most teams.

### Database Snapshot Strategy

**Always snapshot before**:
- Running migrations
- Testing destructive operations
- Pulling new database
- Major WordPress core/plugin updates

**Snapshot naming**:
- Be descriptive: `before-user-migration`
- Include date for long-term: `2024-11-08-pre-deploy`
- Auto snapshots are fine for routine work

**Clean up old snapshots**:
```bash
# List snapshots
stax db list

# Delete old ones
stax db delete-snapshot old-snapshot-name
```

### Code Quality

**Before committing**:
```bash
# Run linters
stax lint --fix

# Run tests
stax ssh
composer test
npm test
exit

# Check build
stax build
```

### Git Workflow Integration

**Recommended workflow**:

```bash
# 1. Create feature branch
git checkout -b feature/new-feature

# 2. Refresh database (optional)
stax db pull --environment=staging

# 3. Make changes
# ... edit files ...

# 4. Build and test
stax build
stax lint
open https://my-project.local

# 5. Commit
git add .
git commit -m "Add new feature"

# 6. Push
git push origin feature/new-feature

# 7. Create PR
# Use GitHub UI or gh CLI
```

### Keeping Dependencies Updated

**Update Stax**:
```bash
brew update
brew upgrade stax
```

**Update DDEV**:
```bash
brew update
brew upgrade ddev
```

**Update project dependencies**:
```bash
stax ssh
composer update
npm update
exit

# Commit updated lockfiles
git add composer.lock package-lock.json
git commit -m "Update dependencies"
```

---

## Tips and Tricks

### Faster Database Imports

**Skip unnecessary data**:
```bash
stax db pull \
  --skip-logs \
  --skip-transients \
  --skip-spam \
  --exclude-tables=wp_actionscheduler_logs
```

Can reduce import time by 50%+.

### Alias Common Commands

Add to your `~/.zshrc` or `~/.bashrc`:

```bash
# Stax shortcuts
alias ss='stax start'
alias st='stax stop'
alias sr='stax restart'
alias sdp='stax db pull'
alias sds='stax db snapshot'
alias swp='stax wp'
alias ssh-stax='stax ssh'
```

Reload shell:
```bash
source ~/.zshrc
```

Now you can use `ss` instead of `stax start`, etc.

### Multiple Terminal Windows

**Recommended setup**:

1. **Window 1**: Code editor
   ```bash
   code .
   ```

2. **Window 2**: Dev watcher
   ```bash
   stax dev
   ```

3. **Window 3**: Commands
   ```bash
   # Available for running stax commands
   stax wp cache flush
   stax db pull
   etc.
   ```

4. **Window 4**: Logs (optional)
   ```bash
   stax logs -f
   ```

### Quick Testing Cycle

```bash
# Make change
vim wp-content/themes/my-theme/style.css

# Rebuild (if not watching)
stax build

# Clear cache
stax wp cache flush

# View in browser
open https://my-project.local

# Or all in one line:
stax build && stax wp cache flush && open https://my-project.local
```

### Environment Variables

**Override config via environment**:

```bash
# Use staging for this command only
STAX_WPENGINE_ENV=staging stax db pull

# Enable debug mode
STAX_DEBUG=true stax init

# Use different config file
STAX_CONFIG=/path/to/config.yml stax start
```

### VS Code Integration

**Recommended extensions**:
- PHP Intelephense
- WordPress Snippets
- ESLint
- Stylelint

**Configure VS Code to use container PHP**:

`.vscode/settings.json`:
```json
{
  "php.validate.executablePath": "/usr/local/bin/ddev-php",
  "phpcs.executablePath": "/usr/local/bin/ddev-phpcs",
  "intelephense.environment.includePaths": [
    "/path/to/project/vendor"
  ]
}
```

### MailHog Tips

**Access MailHog**:
```
http://my-project.local:8025
```

**Test email sending**:
```bash
stax wp user create testuser test@example.com --send-email
```

Check MailHog - you'll see the email!

**Clear all emails**:
Click "Clear" button in MailHog UI.

---

## Next Steps

You now know how to use Stax effectively! Explore these guides next:

- **[Multisite Guide](./MULTISITE.md)** - Deep dive into multisite features
- **[WPEngine Guide](./WPENGINE.md)** - WPEngine-specific operations
- **[Examples](./EXAMPLES.md)** - Real-world scenarios and workflows
- **[Command Reference](./COMMAND_REFERENCE.md)** - Complete command documentation
- **[Troubleshooting](./TROUBLESHOOTING.md)** - Common problems and solutions

---

**Questions?** Check the [FAQ](./FAQ.md) or contact your team!

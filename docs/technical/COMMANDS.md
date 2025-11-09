# Stax Command Reference

## Command Structure Overview

Stax uses a hierarchical command structure with logical groupings for related operations. All commands follow the pattern:

```
stax [command] [subcommand] [flags] [arguments]
```

## Global Flags

Available on all commands:

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--config` | `-c` | string | Path to config file (default: `.stax.yml`) |
| `--verbose` | `-v` | bool | Enable verbose output |
| `--debug` | `-d` | bool | Enable debug logging |
| `--quiet` | `-q` | bool | Suppress non-error output |
| `--help` | `-h` | bool | Display help for command |
| `--version` | | bool | Display version information |
| `--no-color` | | bool | Disable colored output |

## Command Tree

```
stax
├── init              Initialize new project
├── setup             Configure credentials
├── start             Start environment
├── stop              Stop environment
├── restart           Restart environment
├── delete            Delete environment
├── status            Show environment status
├── ssh               SSH into web container
├── logs              View container logs
├── config            Configuration management
│   ├── set           Set configuration value
│   ├── get           Get configuration value
│   ├── list          List all configuration
│   └── validate      Validate configuration
├── db                Database operations
│   ├── pull          Pull database from WPEngine
│   ├── push          Push database to WPEngine (warning)
│   ├── import        Import SQL file
│   ├── export        Export to SQL file
│   ├── snapshot      Create snapshot
│   ├── restore       Restore snapshot
│   ├── list          List snapshots
│   └── query         Execute SQL query
├── wp                WordPress operations
│   ├── cli           Execute WP-CLI command
│   ├── search-replace Run search-replace
│   ├── plugin        Plugin management
│   ├── theme         Theme management
│   ├── site          Multisite site management
│   └── user          User management
├── wpe               WPEngine operations
│   ├── info          Show WPEngine environment info
│   ├── sync          Sync files from WPEngine
│   ├── backups       List available backups
│   ├── deploy        Deploy to WPEngine
│   └── environments  List environments
└── doctor            Diagnose and fix issues
```

---

## Core Commands

### `stax init`

Initialize a new Stax project in the current directory.

**Usage:**
```bash
stax init [flags]
```

**Flags:**
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | (current dir) | Project name |
| `--type` | string | wordpress-multisite | Project type |
| `--mode` | string | subdomain | Multisite mode (subdomain/subdirectory) |
| `--php-version` | string | 8.1 | PHP version |
| `--mysql-version` | string | 8.0 | MySQL version |
| `--repo` | string | | GitHub repository URL |
| `--branch` | string | main | Repository branch |
| `--wpengine-install` | string | | WPEngine install name |
| `--wpengine-env` | string | production | WPEngine environment |
| `--interactive` | bool | true | Enable interactive prompts |
| `--skip-db` | bool | false | Skip database import |
| `--skip-build` | bool | false | Skip build process |

**Interactive Prompts:**
1. Project name
2. Multisite mode (subdomain/subdirectory)
3. Network domain (e.g., firecrown.local)
4. WPEngine install name
5. WPEngine environment (production/staging)
6. GitHub repository URL
7. Add brand sites (repeatable)
   - Brand name
   - Local domain
   - WPEngine domain

**Examples:**

```bash
# Interactive mode (default)
stax init

# Non-interactive with all flags
stax init \
  --name=firecrown-multisite \
  --mode=subdomain \
  --php-version=8.1 \
  --mysql-version=8.0 \
  --repo=https://github.com/Firecrown-Media/firecrown-multisite.git \
  --wpengine-install=fsmultisite \
  --no-interactive

# Initialize without database import
stax init --skip-db

# Initialize without build process
stax init --skip-build
```

**Output:**
```
🚀 Initializing Stax project: firecrown-multisite

✓ Validating WPEngine credentials
✓ Cloning repository from GitHub
✓ Detecting PHP/MySQL versions from WPEngine
✓ Generating DDEV configuration
✓ Starting DDEV containers
  - Web: https://firecrown.local
  - Database: MySQL 8.0
  - MailHog: http://firecrown.local:8025
✓ Installing Composer dependencies
✓ Installing NPM dependencies
✓ Running build script
✓ Pulling database from WPEngine
✓ Importing database
✓ Running search-replace
  - Network: fsmultisite.wpenginepowered.com → firecrown.local
  - Flying Magazine: flyingmag.com → flyingmag.firecrown.local
  - Plane & Pilot: planeandpilotmag.com → planeandpilot.firecrown.local

✓ Project initialized successfully!

Your sites are ready:
  - Network:         https://firecrown.local
  - Flying Magazine: https://flyingmag.firecrown.local
  - Plane & Pilot:   https://planeandpilot.firecrown.local

Next steps:
  - stax status          Check environment status
  - stax ssh             SSH into web container
  - stax logs -f         View logs
```

**Error Scenarios:**
- Invalid WPEngine credentials → "Failed to authenticate with WPEngine"
- Missing GitHub token → "GitHub repository requires authentication"
- DDEV not installed → "DDEV not found. Please install: brew install ddev/ddev/ddev"
- Port conflicts → "Port 80/443 already in use"

---

### `stax setup`

Configure WPEngine and GitHub credentials.

**Usage:**
```bash
stax setup [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--wpengine-user` | string | WPEngine API username |
| `--wpengine-password` | string | WPEngine API password |
| `--github-token` | string | GitHub personal access token |
| `--ssh-key` | string | Path to SSH private key for WPEngine |
| `--interactive` | bool | Interactive credential setup (default: true) |

**Examples:**

```bash
# Interactive mode
stax setup

# Non-interactive
stax setup \
  --wpengine-user=myuser@example.com \
  --wpengine-password=mypassword \
  --github-token=ghp_xxxxxxxxxxxxx \
  --ssh-key=~/.ssh/wpengine_rsa
```

**Output:**
```
🔐 Setting up Stax credentials

WPEngine API Username: myuser@example.com
WPEngine API Password: ********
GitHub Personal Access Token: ghp_xxxxx...
SSH Key for WPEngine: ~/.ssh/wpengine_rsa

✓ Validating WPEngine credentials
✓ Validating GitHub token
✓ Saving credentials to macOS Keychain

Credentials saved successfully!
```

---

### `stax start`

Start the DDEV environment.

**Usage:**
```bash
stax start [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--build` | bool | Run build process after start |
| `--xdebug` | bool | Enable Xdebug |

**Examples:**

```bash
# Basic start
stax start

# Start with Xdebug enabled
stax start --xdebug

# Start and rebuild
stax start --build
```

**Output:**
```
🚀 Starting firecrown-multisite

✓ Starting DDEV containers
  - Web: https://firecrown.local
  - Database: MySQL 8.0 (ready)
  - MailHog: http://firecrown.local:8025

Environment started successfully!
```

---

### `stax stop`

Stop the DDEV environment.

**Usage:**
```bash
stax stop [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--all` | bool | Stop all DDEV projects |
| `--remove-data` | bool | Remove database data (destructive) |

**Examples:**

```bash
# Stop current project
stax stop

# Stop all DDEV projects
stax stop --all
```

**Output:**
```
🛑 Stopping firecrown-multisite

✓ Stopping DDEV containers

Environment stopped.
```

---

### `stax restart`

Restart the DDEV environment.

**Usage:**
```bash
stax restart [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--build` | bool | Run build process after restart |
| `--xdebug` | bool | Enable Xdebug |

**Examples:**

```bash
stax restart
stax restart --build
```

---

### `stax delete`

Delete the DDEV environment and all data.

**Usage:**
```bash
stax delete [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--force` | bool | Skip confirmation prompt |
| `--keep-config` | bool | Keep .stax.yml and .ddev/ configs |

**Examples:**

```bash
# Delete with confirmation
stax delete

# Force delete
stax delete --force

# Delete but keep configs
stax delete --keep-config
```

**Output:**
```
⚠️  WARNING: This will permanently delete the firecrown-multisite environment.

The following will be deleted:
  - DDEV containers
  - Database data
  - .ddev/ directory

Are you sure? (yes/no): yes

✓ Stopping containers
✓ Removing containers
✓ Removing database volumes
✓ Removing DDEV configuration

Environment deleted successfully.
```

---

### `stax status`

Show detailed environment status.

**Usage:**
```bash
stax status [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--json` | bool | Output as JSON |

**Examples:**

```bash
stax status
stax status --json
```

**Output:**
```
📊 Status: firecrown-multisite

Environment: Running ✓

Containers:
  ✓ ddev-firecrown-web       (healthy)
  ✓ ddev-firecrown-db        (healthy)
  ✓ ddev-router              (healthy)

URLs:
  - Network:         https://firecrown.local
  - Flying Magazine: https://flyingmag.firecrown.local
  - Plane & Pilot:   https://planeandpilot.firecrown.local
  - MailHog:         http://firecrown.local:8025

Configuration:
  - PHP Version:     8.1
  - MySQL Version:   8.0
  - Webserver:       nginx
  - Xdebug:          Disabled

Database:
  - Size:            245 MB
  - Tables:          127
  - Sites:           4

WPEngine:
  - Install:         fsmultisite
  - Environment:     production
  - Last Sync:       2 hours ago
```

**JSON Output:**
```json
{
  "status": "running",
  "containers": [
    {"name": "ddev-firecrown-web", "status": "healthy"},
    {"name": "ddev-firecrown-db", "status": "healthy"},
    {"name": "ddev-router", "status": "healthy"}
  ],
  "urls": {
    "network": "https://firecrown.local",
    "sites": [
      {"name": "Flying Magazine", "url": "https://flyingmag.firecrown.local"},
      {"name": "Plane & Pilot", "url": "https://planeandpilot.firecrown.local"}
    ],
    "mailhog": "http://firecrown.local:8025"
  },
  "config": {
    "php_version": "8.1",
    "mysql_version": "8.0",
    "webserver": "nginx",
    "xdebug": false
  }
}
```

---

### `stax ssh`

SSH into the DDEV web container.

**Usage:**
```bash
stax ssh [flags] [command]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--service` | string | Container to SSH into (default: web) |

**Examples:**

```bash
# Interactive SSH session
stax ssh

# Run single command
stax ssh ls -la wp-content/

# SSH into database container
stax ssh --service=db
```

**Output:**
```
geoff@ddev-firecrown-web:/var/www/html$
```

---

### `stax logs`

View container logs.

**Usage:**
```bash
stax logs [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--follow` / `-f` | bool | Follow log output |
| `--tail` | int | Number of lines to show (default: 100) |
| `--service` | string | Show logs for specific service |
| `--timestamp` | bool | Show timestamps |

**Examples:**

```bash
# View last 100 lines
stax logs

# Follow logs
stax logs -f

# View last 500 lines
stax logs --tail=500

# View web container logs only
stax logs --service=web

# Follow with timestamps
stax logs -f --timestamp
```

---

## Configuration Commands

### `stax config:set`

Set a configuration value.

**Usage:**
```bash
stax config:set <key> <value> [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--global` | bool | Set in global config (~/.stax/config.yml) |

**Examples:**

```bash
# Set project config
stax config:set wpengine.environment staging

# Set global config
stax config:set wpengine.api_user myuser@example.com --global

# Set PHP version
stax config:set ddev.php_version 8.2

# Set multisite mode
stax config:set project.mode subdirectory
```

**Output:**
```
✓ Configuration updated: wpengine.environment = staging
```

---

### `stax config:get`

Get a configuration value.

**Usage:**
```bash
stax config:get <key> [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--global` | bool | Get from global config |

**Examples:**

```bash
# Get project config
stax config:get wpengine.environment

# Get global config
stax config:get wpengine.api_user --global
```

**Output:**
```
production
```

---

### `stax config:list`

List all configuration values.

**Usage:**
```bash
stax config:list [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--global` | bool | List global config |
| `--json` | bool | Output as JSON |

**Examples:**

```bash
stax config:list
stax config:list --global
stax config:list --json
```

**Output:**
```
Project Configuration (.stax.yml):

project:
  name: firecrown-multisite
  type: wordpress-multisite
  mode: subdomain

wpengine:
  environment: production
  install: fsmultisite

ddev:
  php_version: 8.1
  mysql_version: 8.0
  webserver_type: nginx
```

---

### `stax config:validate`

Validate configuration file.

**Usage:**
```bash
stax config:validate [flags]
```

**Examples:**

```bash
stax config:validate
```

**Output:**
```
✓ Configuration is valid

Warnings:
  - PHP version 8.1 is older than WPEngine production (8.2)
```

---

## Database Commands

### `stax db:pull`

Pull database from WPEngine.

**Usage:**
```bash
stax db:pull [flags]
```

**Flags:**
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--environment` | string | (from config) | WPEngine environment |
| `--snapshot` | bool | true | Create snapshot before import |
| `--sanitize` | bool | false | Sanitize user data |
| `--skip-replace` | bool | false | Skip search-replace |
| `--exclude-tables` | string | | Comma-separated tables to exclude |
| `--skip-logs` | bool | true | Skip log tables |
| `--skip-transients` | bool | true | Skip transient tables |
| `--skip-spam` | bool | true | Skip spam/trash |

**Examples:**

```bash
# Basic pull
stax db:pull

# Pull from staging
stax db:pull --environment=staging

# Pull without snapshot
stax db:pull --snapshot=false

# Pull with sanitized data
stax db:pull --sanitize

# Pull specific tables only
stax db:pull --exclude-tables=wp_actionscheduler_logs,wp_wc_admin_notes

# Pull without search-replace
stax db:pull --skip-replace
```

**Output:**
```
🗄️  Pulling database from WPEngine (production)

✓ Creating snapshot: db_2025-11-08_14-30-00
✓ Connecting to WPEngine SSH Gateway
✓ Detecting table prefix: wp_
✓ Exporting database (excluding logs, transients, spam)
  - Tables: 127
  - Estimated size: 245 MB
✓ Transferring database
  [████████████████████████████████] 245 MB / 245 MB
✓ Importing to local database
  - Rows imported: 1,245,678
✓ Running search-replace
  - Network: fsmultisite.wpenginepowered.com → firecrown.local (43 replacements)
  - Flying Magazine: flyingmag.com → flyingmag.firecrown.local (1,234 replacements)
  - Plane & Pilot: planeandpilotmag.com → planeandpilot.firecrown.local (987 replacements)
✓ Flushing WordPress cache

Database pulled successfully!
Time elapsed: 2m 34s
```

---

### `stax db:push`

Push database to WPEngine (with safety warnings).

**Usage:**
```bash
stax db:push [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--environment` | string | WPEngine environment (required) |
| `--force` | bool | Skip confirmation (dangerous) |
| `--dry-run` | bool | Show what would be pushed |

**Examples:**

```bash
# Push with confirmation
stax db:push --environment=staging

# Dry run
stax db:push --environment=staging --dry-run
```

**Output:**
```
⚠️  WARNING: This will OVERWRITE the WPEngine staging database!

You are about to push to:
  - Install: fsmultisite
  - Environment: staging
  - URL: fsmultisite-staging.wpengine.com

This cannot be undone. Are you absolutely sure? (type 'yes' to confirm): yes

✓ Creating backup on WPEngine
✓ Exporting local database
✓ Running search-replace
  - firecrown.local → fsmultisite.wpenginepowered.com
✓ Uploading to WPEngine
✓ Importing on WPEngine

Database pushed successfully.
```

---

### `stax db:import`

Import SQL file into local database.

**Usage:**
```bash
stax db:import <file> [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--snapshot` | bool | Create snapshot before import (default: true) |
| `--replace` | bool | Run search-replace after import |

**Examples:**

```bash
# Import SQL file
stax db:import ~/Downloads/backup.sql

# Import without snapshot
stax db:import backup.sql --snapshot=false

# Import and run search-replace
stax db:import backup.sql --replace
```

**Output:**
```
📥 Importing database from backup.sql

✓ Creating snapshot: db_2025-11-08_14-35-00
✓ Importing SQL file
  - Size: 245 MB
  - Rows imported: 1,245,678

Database imported successfully!
```

---

### `stax db:export`

Export local database to SQL file.

**Usage:**
```bash
stax db:export [file] [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--gzip` | bool | Compress with gzip |
| `--exclude-tables` | string | Comma-separated tables to exclude |
| `--skip-logs` | bool | Skip log tables (default: true) |

**Examples:**

```bash
# Export to default location
stax db:export

# Export to specific file
stax db:export ~/backups/firecrown-$(date +%Y%m%d).sql

# Export compressed
stax db:export --gzip

# Export without logs
stax db:export --skip-logs
```

**Output:**
```
📤 Exporting database

✓ Exporting database
  - Tables: 127
  - Size: 245 MB
✓ Saved to: ./firecrown-multisite-2025-11-08.sql

Database exported successfully!
```

---

### `stax db:snapshot`

Create named database snapshot.

**Usage:**
```bash
stax db:snapshot [name] [flags]
```

**Examples:**

```bash
# Auto-named snapshot
stax db:snapshot

# Named snapshot
stax db:snapshot before-migration

# Snapshot with description
stax db:snapshot before-migration --description="Before user table migration"
```

**Output:**
```
📸 Creating snapshot: before-migration

✓ Snapshot created: before-migration
  - Size: 245 MB
  - Location: ~/.stax/snapshots/firecrown-multisite/before-migration.sql.gz
```

---

### `stax db:restore`

Restore database from snapshot.

**Usage:**
```bash
stax db:restore <snapshot> [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--force` | bool | Skip confirmation |

**Examples:**

```bash
# Restore specific snapshot
stax db:restore before-migration

# Restore latest snapshot
stax db:restore latest

# Force restore
stax db:restore before-migration --force
```

**Output:**
```
♻️  Restoring snapshot: before-migration

⚠️  This will replace your current database. Continue? (yes/no): yes

✓ Creating backup of current database
✓ Restoring snapshot
  - Size: 245 MB
  - Rows imported: 1,245,678

Database restored successfully!
```

---

### `stax db:list`

List all database snapshots.

**Usage:**
```bash
stax db:list [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--json` | bool | Output as JSON |

**Examples:**

```bash
stax db:list
stax db:list --json
```

**Output:**
```
📋 Database Snapshots

Name                    Created              Size     Description
db_2025-11-08_14-30-00  2 hours ago         245 MB   Auto (before pull)
before-migration        1 day ago           243 MB   Before user table migration
db_2025-11-07_09-15-00  1 day ago           240 MB   Auto (before pull)

Total: 3 snapshots (728 MB)
```

---

### `stax db:query`

Execute SQL query.

**Usage:**
```bash
stax db:query <sql> [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--file` | string | Read query from file |
| `--format` | string | Output format (table/json/csv) |

**Examples:**

```bash
# Direct query
stax db:query "SELECT * FROM wp_options WHERE option_name = 'siteurl'"

# Query from file
stax db:query --file=query.sql

# JSON output
stax db:query "SELECT * FROM wp_users LIMIT 5" --format=json
```

---

## WordPress Commands

### `stax wp`

Execute WP-CLI command.

**Usage:**
```bash
stax wp <command> [args...] [flags]
```

**Examples:**

```bash
# Any WP-CLI command
stax wp plugin list
stax wp user list --role=administrator
stax wp site list
stax wp cache flush
stax wp db check
stax wp core version

# Complex commands
stax wp search-replace oldtext newtext --dry-run
stax wp eval "echo WP_CONTENT_DIR;"
```

**Output:**
```
# stax wp plugin list
+---------------------+----------+--------+---------+
| name                | status   | update | version |
+---------------------+----------+--------+---------+
| akismet             | inactive | none   | 4.2.2   |
| hello               | inactive | none   | 1.7.2   |
+---------------------+----------+--------+---------+
```

---

### `stax wp:search-replace`

Run search-replace across all sites.

**Usage:**
```bash
stax wp:search-replace <old> <new> [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--network` | bool | Run across entire network |
| `--site` | string | Run on specific site |
| `--dry-run` | bool | Show what would be replaced |
| `--skip-columns` | string | Skip specific columns |

**Examples:**

```bash
# Replace across network
stax wp:search-replace http://old.com https://new.com --network

# Replace on specific site
stax wp:search-replace old.com new.com --site=flyingmag

# Dry run
stax wp:search-replace old.com new.com --network --dry-run
```

---

### `stax wp:plugin`

Manage WordPress plugins.

**Usage:**
```bash
stax wp:plugin <action> [plugin] [flags]
```

**Actions:** list, activate, deactivate, install, uninstall, update

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--network` | bool | Network-wide activation/deactivation |

**Examples:**

```bash
# List plugins
stax wp:plugin list

# Install plugin
stax wp:plugin install wordpress-seo

# Activate plugin network-wide
stax wp:plugin activate wordpress-seo --network

# Deactivate plugin
stax wp:plugin deactivate akismet

# Update all plugins
stax wp:plugin update --all
```

---

### `stax wp:theme`

Manage WordPress themes.

**Usage:**
```bash
stax wp:theme <action> [theme] [flags]
```

**Actions:** list, activate, install, delete

**Examples:**

```bash
stax wp:theme list
stax wp:theme activate twentytwentyfour
stax wp:theme install twentytwentyfour
```

---

### `stax wp:site`

Manage multisite sites.

**Usage:**
```bash
stax wp:site <action> [args...] [flags]
```

**Actions:** list, create, delete, empty, activate, deactivate

**Examples:**

```bash
# List all sites
stax wp:site list

# Create new site
stax wp:site create --slug=newsite --title="New Site"

# Delete site
stax wp:site delete 5

# Empty site (delete all posts/pages)
stax wp:site empty 5
```

---

### `stax wp:user`

Manage WordPress users.

**Usage:**
```bash
stax wp:user <action> [args...] [flags]
```

**Actions:** list, create, delete, update, reset-password

**Examples:**

```bash
# List users
stax wp:user list

# Create admin user
stax wp:user create admin admin@example.com --role=administrator

# Reset password
stax wp:user reset-password admin --password=newpassword
```

---

## WPEngine Commands

### `stax wpe:info`

Show WPEngine environment information.

**Usage:**
```bash
stax wpe:info [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--environment` | string | Environment to query |
| `--json` | bool | Output as JSON |

**Examples:**

```bash
stax wpe:info
stax wpe:info --environment=staging
stax wpe:info --json
```

**Output:**
```
🌐 WPEngine Environment: fsmultisite (production)

Install Details:
  - Name:            fsmultisite
  - Primary Domain:  fsmultisite.wpengine.com
  - PHP Version:     8.2
  - MySQL Version:   8.0
  - WordPress:       6.4.2
  - Disk Usage:      2.4 GB / 10 GB
  - Account:         Firecrown Media

Environment:
  - Type:            production
  - Git Enabled:     Yes
  - SSH Access:      Yes
  - CDN:             Enabled

Domains:
  - firecrown.local
  - flyingmag.firecrown.local
  - planeandpilot.firecrown.local
  - finescale.firecrown.local
  - avweb.firecrown.local
```

---

### `stax wpe:sync`

Sync files from WPEngine.

**Usage:**
```bash
stax wpe:sync [path] [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--environment` | string | WPEngine environment |
| `--exclude` | string | Exclude patterns (rsync format) |
| `--dry-run` | bool | Show what would be synced |
| `--delete` | bool | Delete local files not on remote |

**Examples:**

```bash
# Sync wp-content/uploads
stax wpe:sync wp-content/uploads

# Sync with dry run
stax wpe:sync wp-content/uploads --dry-run

# Sync from staging
stax wpe:sync wp-content/uploads --environment=staging

# Sync with exclusions
stax wpe:sync wp-content --exclude="*.log,cache/"
```

**Output:**
```
🔄 Syncing from WPEngine (production): wp-content/uploads

✓ Connecting to SSH Gateway
✓ Syncing files
  [████████████████████████████████] 1,234 files (456 MB)

Files synced successfully!
Time elapsed: 1m 23s
```

---

### `stax wpe:backups`

List available WPEngine backups.

**Usage:**
```bash
stax wpe:backups [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--environment` | string | Environment to query |
| `--json` | bool | Output as JSON |

**Examples:**

```bash
stax wpe:backups
stax wpe:backups --environment=staging
```

**Output:**
```
📦 WPEngine Backups: fsmultisite (production)

ID    Type       Created              Size     Status
1234  automatic  2 hours ago         1.2 GB   complete
1233  manual     1 day ago           1.1 GB   complete
1232  automatic  2 days ago          1.1 GB   complete

Total: 3 backups (3.4 GB)
```

---

### `stax wpe:deploy`

Deploy to WPEngine via GitHub Actions.

**Usage:**
```bash
stax wpe:deploy [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--environment` | string | Environment to deploy to |
| `--watch` | bool | Watch deployment progress |
| `--branch` | string | Branch to deploy (default: current) |

**Examples:**

```bash
# Deploy to staging
stax wpe:deploy --environment=staging

# Deploy to production and watch
stax wpe:deploy --environment=production --watch

# Deploy specific branch
stax wpe:deploy --environment=staging --branch=feature/new-feature
```

**Output:**
```
🚀 Deploying to WPEngine (staging)

✓ Triggering GitHub Actions workflow
  - Workflow: Deploy to WPEngine
  - Branch: main
  - Run ID: 1234567890

Watching deployment...
  ⏳ Building...
  ⏳ Running tests...
  ⏳ Deploying to WPEngine...
  ✓ Deployment complete!

View deployment: https://github.com/Firecrown-Media/firecrown-multisite/actions/runs/1234567890
```

---

### `stax wpe:environments`

List available WPEngine environments.

**Usage:**
```bash
stax wpe:environments [flags]
```

**Examples:**

```bash
stax wpe:environments
```

**Output:**
```
🌐 WPEngine Environments: fsmultisite

Name        Type        URL                                  Status
production  production  fsmultisite.wpengine.com            active
staging     staging     fsmultisite-staging.wpengine.com    active

Total: 2 environments
```

---

## Diagnostic Commands

### `stax doctor`

Diagnose and fix common issues.

**Usage:**
```bash
stax doctor [flags]
```

**Flags:**
| Flag | Type | Description |
|------|------|-------------|
| `--fix` | bool | Automatically fix issues |

**Examples:**

```bash
# Diagnose issues
stax doctor

# Diagnose and fix
stax doctor --fix
```

**Output:**
```
🩺 Running diagnostics...

✓ DDEV installed (v1.22.0)
✓ Docker Desktop running
✓ WPEngine credentials valid
✓ GitHub token valid
✗ Port 443 in use by Apache
✗ Database connection slow
✓ SSL certificates valid
✓ WordPress core files intact
⚠️  PHP version mismatch (local: 8.1, WPEngine: 8.2)

Issues found: 2 errors, 1 warning

Errors:
  1. Port 443 in use by Apache
     Fix: sudo apachectl stop

  2. Database connection slow
     Fix: Increase MySQL max_connections

Warnings:
  1. PHP version mismatch
     Recommendation: Update .stax.yml to use PHP 8.2

Run 'stax doctor --fix' to automatically fix issues.
```

---

## Examples and Workflows

### Complete Setup Workflow

```bash
# 1. Initial setup (one-time)
stax setup

# 2. Initialize project
cd ~/Sites/firecrown-multisite
stax init

# 3. Start working
stax start

# 4. Make changes, test
stax wp plugin list
stax logs -f

# 5. Refresh database
stax db:pull

# 6. Stop when done
stax stop
```

### Daily Development Workflow

```bash
# Morning: Start environment
stax start

# Refresh database if needed
stax db:pull --environment=staging

# Work on code...

# Test changes
stax wp cache flush
stax ssh
  # Run tests, etc.

# Evening: Stop environment
stax stop
```

### Database Testing Workflow

```bash
# Create snapshot before testing
stax db:snapshot before-migration

# Test migration
stax wp db query --file=migration.sql

# If something breaks, restore
stax db:restore before-migration
```

### Deployment Workflow

```bash
# Test locally
git checkout feature/new-feature
stax restart --build

# Deploy to staging
git push origin feature/new-feature
stax wpe:deploy --environment=staging --watch

# After QA approval, deploy to production
git checkout main
git merge feature/new-feature
git push origin main
stax wpe:deploy --environment=production --watch
```

---

## Command Aliases

Stax supports common aliases for frequently used commands:

| Alias | Full Command |
|-------|--------------|
| `stax up` | `stax start` |
| `stax down` | `stax stop` |
| `stax rm` | `stax delete` |
| `stax s` | `stax status` |
| `stax i` | `stax init` |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments/flags |
| 3 | Missing dependencies (DDEV, Docker) |
| 4 | Authentication failed |
| 5 | Network error (WPEngine, GitHub) |
| 6 | Container error (DDEV) |
| 10 | User cancelled operation |

---

## Environment Variables

Stax respects the following environment variables:

| Variable | Description |
|----------|-------------|
| `STAX_CONFIG` | Path to config file |
| `STAX_DEBUG` | Enable debug output (true/false) |
| `STAX_NO_COLOR` | Disable colored output |
| `WPENGINE_API_USER` | WPEngine API username |
| `WPENGINE_API_PASSWORD` | WPEngine API password |
| `GITHUB_TOKEN` | GitHub personal access token |

---

## Shell Completion

Stax supports shell completion for Bash, Zsh, and Fish.

**Install completion:**

```bash
# Bash
stax completion bash > /usr/local/etc/bash_completion.d/stax

# Zsh
stax completion zsh > /usr/local/share/zsh/site-functions/_stax

# Fish
stax completion fish > ~/.config/fish/completions/stax.fish
```

**Usage:**

```bash
stax db:<TAB>
# Suggests: pull, push, import, export, snapshot, restore, list, query

stax wp:plugin <TAB>
# Suggests: list, activate, deactivate, install, uninstall, update
```

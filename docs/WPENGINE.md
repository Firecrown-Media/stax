# WP Engine Integration Guide

Complete guide for integrating Stax with WP Engine hosting environments.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [Understanding stax wpe sync](#understanding-stax-wpe-sync)
- [Common Workflows](#common-workflows)
- [Command Reference](#command-reference)
- [Performance Optimization](#performance-optimization)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before using WP Engine features, you need:

### 1. WP Engine Account Access
- Contact your team administrator for account access
- Ensure you have permissions for target installations
- Verify you can log into [my.wpengine.com](https://my.wpengine.com)

### 2. SSH Key Configuration

Generate and add SSH keys for secure connection:

```bash
# Generate SSH key (if you don't have one)
ssh-keygen -t rsa -b 4096 -f ~/.ssh/wpengine_rsa -C "your.email@company.com"

# View your public key
cat ~/.ssh/wpengine_rsa.pub
```

**Add to WP Engine:**
1. Log into [my.wpengine.com](https://my.wpengine.com)
2. Navigate to your name (top right) → "SSH keys"
3. Click "Add SSH key"
4. Paste your public key
5. Name it (e.g., "MacBook Pro - Development")

### 3. SSH Config Setup

Configure SSH for easier connections:

```bash
# Edit SSH config
nano ~/.ssh/config

# Add this configuration
Host *.ssh.wpengine.net
  IdentityFile ~/.ssh/wpengine_rsa
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
```

### 4. Test Connection

Verify SSH access works:

```bash
# Test connection (replace 'installname' with actual install)
ssh installname@installname.ssh.wpengine.net "echo 'Connected successfully!'"
```

### 5. API Credentials (Optional)

For API access, set environment variables:

```bash
# Add to ~/.zshrc or ~/.bashrc
export WPE_USERNAME="your-username"
export WPE_PASSWORD="your-password"

# Reload shell
source ~/.zshrc
```

## Initial Setup

### Finding Install Names

WP Engine install names are the internal identifiers for sites:

```bash
# List all available installs
stax wpe list

# Get details about specific install
stax wpe info installname
```

Common install naming patterns:
- Production: `clientprod` or `clientname`
- Staging: `clientstg` or `clientstage`
- Development: `clientdev`

### First Sync

Set up your first WP Engine sync:

```bash
# 1. Create local project
stax init client-site

# 2. Full sync (database + files)
stax wpe sync clientprod
# This will:
# - Connect via SSH
# - Export database from WP Engine
# - Download database
# - Import to local DDEV
# - Sync WordPress files
# - Update URLs for local development

# 3. Start development
stax start client-site
```

## Understanding stax wpe sync

### How Sync Works

The sync process follows these steps:

1. **SSH Connection** - Connects to WP Engine via SSH
2. **Database Export** - Creates database dump on WP Engine
3. **Database Download** - Transfers SQL file via SCP
4. **Database Import** - Imports into local DDEV MySQL
5. **URL Replacement** - Updates URLs from production to local
6. **File Sync** - Uses rsync to copy WordPress files (optional)
7. **Cleanup** - Removes temporary files

### Sync Flags Explained

| Flag | Purpose | When to Use |
|------|---------|-------------|
| `--skip-files` | Only sync database | Daily development (faster) |
| `--skip-database` | Only sync files | Updating media/plugins |
| `--skip-media` | Exclude wp-content/uploads | Default on, saves bandwidth |
| `--environment` | Choose environment | Sync from staging/dev |
| `--delete-local` | Mirror exactly | ⚠️ Dangerous - deletes local files |
| `--create-upload-redirect` | Redirect media to production | Avoid downloading large media |
| `--suppress-debug` | Hide PHP notices | Cleaner output |

### What Gets Synced

**Database sync includes:**
- All WordPress tables
- Users and permissions
- Posts, pages, and custom content
- Settings and configurations

**File sync includes:**
- WordPress core files
- Themes (`wp-content/themes/`)
- Plugins (`wp-content/plugins/`)
- Must-use plugins (`wp-content/mu-plugins/`)
- Uploads (unless `--skip-media`)

**Automatically excluded:**
- Cache directories
- Backup files
- `.git` directories
- `node_modules`
- Temporary files

## Common Workflows

### Daily Development Workflow

Start your day with fresh data:

```bash
# Morning sync - database only for speed
stax wpe sync clientprod --skip-files

# Start development
stax start client-site

# Work on features...

# End of day
stax stop client-site
```

### Initial Project Setup

First time working on a client site:

```bash
# 1. Get install information
stax wpe list | grep client
stax wpe info clientprod

# 2. Create and sync
stax init client-local
stax wpe sync clientprod --environment=production

# 3. Create development user
stax start client-local
stax wp user create yourusername you@company.com --role=administrator
stax wp user update yourusername --user_pass=secure-password
```

### Staging Environment Sync

Work with staging instead of production:

```bash
# Sync from staging
stax wpe sync clientstg --environment=staging

# Or specify in the command
stax wpe sync clientstaging --skip-files
```

### Media Optimization Workflow

Handle large media libraries efficiently:

```bash
# Option 1: Skip media entirely
stax wpe sync clientprod --skip-files

# Option 2: Redirect media to production
stax wpe sync clientprod --skip-files --create-upload-redirect
# Creates mu-plugin that redirects missing images to production

# Option 3: Selective media sync
stax wpe sync clientprod --skip-database
# Then manually sync specific directories
rsync -avz clientprod@clientprod.ssh.wpengine.net:sites/clientprod/wp-content/uploads/2024/ \
  ./wp-content/uploads/2024/
```

### Database-Only Updates

Quick database refresh without files:

```bash
# Download latest database
stax wpe db download clientprod --output=latest.sql

# Import to existing project
stax wpe db import latest.sql

# Update URLs
stax wpe db rewrite clientprod
```

## Command Reference

### List and Information

```bash
# List all installs
stax wpe list

# Detailed install information
stax wpe info [install]

# Test connection
stax wpe connect [install]
```

### Sync Commands

```bash
# Full sync (default)
stax wpe sync [install]

# Database only
stax wpe sync [install] --skip-files

# Files only
stax wpe sync [install] --skip-database

# From staging
stax wpe sync [install] --environment=staging

# With media redirect
stax wpe sync [install] --skip-files --create-upload-redirect
```

### Database Commands

```bash
# Download database
stax wpe db download [install] --output=backup.sql

# Import database
stax wpe db import backup.sql

# Analyze media URLs
stax wpe db analyze

# Rewrite URLs for local
stax wpe db rewrite [install]

# Diagnose URL issues
stax wpe db diagnose
```

### Advanced Options

```bash
# Verbose output for debugging
stax --verbose wpe sync [install]

# Custom SSH username
stax wpe sync [install] --username=customuser

# Force delete local files (careful!)
stax wpe sync [install] --delete-local
```

## Performance Optimization

### Faster Syncing

**1. Skip unnecessary files:**
```bash
# Most development only needs database
stax wpe sync install --skip-files
```

**2. Use media redirects:**
```bash
# Avoid downloading gigabytes of images
stax wpe sync install --skip-files --create-upload-redirect
```

**3. Sync from staging:**
```bash
# Staging databases are often smaller
stax wpe sync install --environment=staging
```

**4. Exclude large directories:**
```bash
# Configure exclusions in ~/.stax.yaml
hosting:
  wpengine:
    sync_defaults:
      exclude_dirs:
        - "wp-content/uploads/backups/"
        - "wp-content/cache/"
        - "wp-content/ai1wm-backups/"
```

### Managing Large Sites

For sites over 1GB:

```bash
# 1. Initial sync - database only
stax wpe sync largesite --skip-files

# 2. Selective file sync
# Get only current year's uploads
rsync -avz largesite@largesite.ssh.wpengine.net:sites/largesite/wp-content/uploads/2024/ \
  ./wp-content/uploads/2024/

# 3. Use CDN for remaining media
stax wpe sync largesite --create-upload-redirect
```

## Troubleshooting

### SSH Connection Issues

**Problem: "Permission denied (publickey)"**
```bash
# Check SSH key is loaded
ssh-add -l

# Add key if missing
ssh-add ~/.ssh/wpengine_rsa

# Test direct SSH
ssh -v installname@installname.ssh.wpengine.net
```

**Problem: "Could not resolve hostname"**
```bash
# Verify install name
stax wpe list

# Check DNS
nslookup installname.ssh.wpengine.net
```

### Sync Failures

**Problem: "Database export failed"**
```bash
# Check WP Engine status
curl https://wpengine.com/support/status/

# Try manual export
ssh install@install.ssh.wpengine.net
cd sites/install
wp db export backup.sql
exit

# Download manually
scp install@install.ssh.wpengine.net:sites/install/backup.sql ./
```

**Problem: "Import failed - database too large"**
```bash
# Check available space
df -h

# Clear Docker space
docker system prune -a

# Import with progress
pv database.sql | ddev mysql
```

### Performance Issues

**Problem: "Sync is very slow"**
```bash
# Check connection speed
ssh install@install.ssh.wpengine.net "dd if=/dev/zero bs=1M count=10" | dd of=/dev/null

# Use compression
stax wpe sync install --verbose  # Check if compression is enabled

# Skip large directories
stax wpe sync install --skip-media
```

### Warp Terminal Compatibility

**Known Issue:** Warp terminal has SSH issues with WP Engine.

**Workaround:** Use alternative terminals:
- Terminal.app (Mac default)
- iTerm2
- VS Code integrated terminal

```bash
# In alternative terminal
stax wpe sync install
```

## Security Best Practices

1. **Never commit credentials:**
```bash
# Use environment variables
export WPE_USERNAME="user"
export WPE_PASSWORD="pass"
```

2. **Protect SSH keys:**
```bash
chmod 600 ~/.ssh/wpengine_rsa
chmod 644 ~/.ssh/wpengine_rsa.pub
```

3. **Use separate keys per machine:**
```bash
ssh-keygen -t rsa -f ~/.ssh/wpengine_laptop
ssh-keygen -t rsa -f ~/.ssh/wpengine_desktop
```

4. **Rotate credentials regularly:**
- Update SSH keys every 6 months
- Change API passwords quarterly

## Next Steps

- Review [Troubleshooting](TROUBLESHOOTING.md) for common issues
- Learn about [Multisite](MULTISITE.md) configurations
- Explore [Development](DEVELOPMENT.md) for extending Stax

## Getting Help

- WP Engine Support: [help.wpengine.com](https://help.wpengine.com)
- Stax Issues: [GitHub](https://github.com/Firecrown-Media/stax/issues)
- Team Support: Contact your DevOps lead
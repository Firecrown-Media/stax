# WPEngine Provider Documentation

## Overview

The WPEngine provider enables Stax to integrate with WPEngine WordPress hosting platform. This document covers WPEngine-specific configuration, capabilities, limitations, and best practices.

## Table of Contents

- [Capabilities](#capabilities)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Features](#features)
- [Limitations](#limitations)
- [Common Operations](#common-operations)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

## Capabilities

### Supported Features

The WPEngine provider supports the following capabilities:

| Capability | Supported | Method |
|-----------|-----------|--------|
| Authentication | Yes | API + SSH |
| Site Management | Yes | API |
| Database Export | Yes | SSH + WP-CLI |
| Database Import | No | Portal only |
| File Sync | Yes | Rsync over SSH |
| File Upload | No | Git only |
| Remote Execution | Yes | SSH gateway |
| WP-CLI | Yes | SSH gateway |
| Backups | Yes | API |
| Environments | Yes | Production/Staging |
| CDN | Yes | BunnyCDN |
| Deployments | Limited | Git push |

### Provider Capabilities Structure

```go
Authentication:  true   // API credentials + SSH key
SiteManagement:  true   // List sites, get details
DatabaseExport:  true   // Via SSH + mysqldump
DatabaseImport:  false  // Not supported (security)
FileSync:        true   // Rsync over SSH
Deployment:      false  // Git available but not implemented
Environments:    true   // Production, Staging, Dev
Backups:         true   // Point-in-time backups
RemoteExecution: true   // SSH gateway access
MediaManagement: true   // BunnyCDN integration
SSHAccess:       true   // SSH gateway
APIAccess:       true   // WPEngine API v1
```

## Prerequisites

### 1. WPEngine Account

- Active WPEngine hosting account
- Access to WPEngine User Portal
- Install/site already created

### 2. API Credentials

1. Log in to [WPEngine User Portal](https://my.wpengine.com/)
2. Navigate to **Account** → **API Access**
3. Generate API credentials
4. Save API username and password securely

### 3. SSH Key

1. Generate SSH key pair (if you don't have one):
   ```bash
   ssh-keygen -t ed25519 -C "your-email@example.com" -f ~/.ssh/wpengine
   ```

2. Add public key to WPEngine:
   - Navigate to **Account** → **SSH Keys**
   - Click **Add SSH Key**
   - Paste your public key (`~/.ssh/wpengine.pub`)

3. Test SSH access:
   ```bash
   ssh your-install@your-install.ssh.wpengine.net
   ```

## Configuration

### Project Configuration (.stax.yml)

```yaml
project:
  name: my-site
  type: wordpress

provider:
  name: wpengine

  wpengine:
    site: my-install-name               # WPEngine install name
    environment: production              # production, staging, development
    ssh_gateway: my-install-name.ssh.wpengine.net
```

### Credentials (~/.stax/credentials.yml)

**Important**: Set file permissions to 600:
```bash
chmod 600 ~/.stax/credentials.yml
```

```yaml
wpengine:
  api_user: your-api-username
  api_password: your-api-password
  ssh_key: |
    -----BEGIN OPENSSH PRIVATE KEY-----
    b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
    ...
    -----END OPENSSH PRIVATE KEY-----
```

### Environment Variables (Optional)

```bash
export WPENGINE_API_USER="your-api-username"
export WPENGINE_API_PASSWORD="your-api-password"
export WPENGINE_SSH_KEY_PATH="~/.ssh/wpengine"
```

## Features

### 1. Site Management

**List Sites:**
```bash
stax provider show wpengine
```

**Get Site Details:**
```bash
# Automatically uses configured install
stax status --provider=wpengine
```

### 2. Database Operations

**Export Database:**
```bash
# Pull to local DDEV
stax db pull --provider=wpengine

# Export to file
stax db export backup.sql --provider=wpengine

# Export with exclusions
stax db export backup.sql \
  --provider=wpengine \
  --skip-logs \
  --exclude-tables=wp_analytics
```

**Import Database:**

WPEngine does not support direct database import via API or SSH for security reasons.

**Workaround:**
1. Export from local DDEV:
   ```bash
   stax db export wpengine-import.sql
   ```

2. Import via WPEngine portal:
   - Navigate to your install in User Portal
   - Go to **Backup Points** → **Import**
   - Upload SQL file

3. Or use staging:
   ```bash
   # Push to staging first
   # Then copy staging to production in portal
   ```

### 3. File Sync

**Sync wp-content from WPEngine:**
```bash
# Sync all wp-content
stax files sync --provider=wpengine

# Sync only uploads
stax files sync --provider=wpengine --include="uploads/**"

# Exclude cache
stax files sync --provider=wpengine --exclude="cache/**"

# Dry run to preview
stax files sync --provider=wpengine --dry-run
```

**Upload Files:**

WPEngine has a read-only filesystem except for Git deployments.

**Workaround:**
- Use Git push deployments for code
- Upload media via WordPress admin or SFTP
- Use WPEngine portal for bulk uploads

### 4. Remote WP-CLI

Execute WP-CLI commands on WPEngine:

```bash
# Core version
stax wp core version --provider=wpengine

# Plugin list
stax wp plugin list --provider=wpengine

# Search/replace (careful!)
stax wp search-replace old.com new.com --provider=wpengine

# User management
stax wp user list --provider=wpengine

# Database operations
stax wp db query "SELECT * FROM wp_options WHERE option_name='siteurl'" \
  --provider=wpengine
```

### 5. Backups

**List Backups:**
```bash
# TODO: Implement backup listing
# stax backup list --provider=wpengine
```

**Create Backup:**
```bash
# TODO: Implement manual backup
# stax backup create "Pre-deployment backup" --provider=wpengine
```

**Via WPEngine Portal:**
- Navigate to **Backup Points**
- Create manual backup
- Download or restore

### 6. Environments

WPEngine supports multiple environments per install:

- **Production**: Live site
- **Staging**: Testing environment
- **Development**: Development environment (some plans)

**Switch Environment:**
```yaml
# In .stax.yml
provider:
  wpengine:
    environment: staging  # or production, development
```

**Copy Between Environments:**
Via WPEngine portal only:
1. Navigate to **Staging** tab
2. **Copy Production to Staging** or
3. **Push Staging to Production**

### 7. CDN (BunnyCDN)

WPEngine includes BunnyCDN integration:

- Automatic image optimization
- Global CDN distribution
- Automatic cache purging on content updates

**Configuration:**
Managed via WPEngine portal in install settings.

## Limitations

### Database Import Restriction

**Issue**: Cannot import databases directly via SSH or API

**Reason**: Security policy to prevent SQL injection and malicious code

**Workarounds**:
1. Use WPEngine portal import feature
2. Copy from staging to production (if already imported to staging)
3. Contact WPEngine support for large imports

### Read-Only Filesystem

**Issue**: Cannot upload files via SSH/SFTP

**Reason**: Security policy, Git-based deployments only

**Workarounds**:
1. Use Git deployments for code
2. Upload media via WordPress admin
3. Use WPEngine portal for bulk uploads
4. Use SFTP for limited file access (with WPEngine support approval)

### No Root Access

**Issue**: No root or sudo access to servers

**Reason**: Managed platform

**Impact**:
- Cannot install system packages
- Cannot modify server configuration
- Cannot restart services

**Workarounds**:
- Use WPEngine's optimized environment
- Request features via WPEngine support
- Use wp-config.php for PHP settings

### Limited Environment Variables

**Issue**: Limited access to environment variables

**Workarounds**:
- Define in wp-config.php
- Use WPEngine portal for environment-specific settings

## Common Operations

### Daily Development Workflow

```bash
# 1. Pull latest from WPEngine
stax db pull --provider=wpengine
stax files sync --provider=wpengine

# 2. Start local environment
stax start

# 3. Make changes locally
# ...

# 4. Deploy via Git
git add .
git commit -m "Your changes"
git push wpengine main
```

### Multisite Development

```yaml
# .stax.yml
wordpress:
  multisite: true
  multisite_type: subdirectory

provider:
  wpengine:
    site: firecrown-multisite
```

```bash
# Pull multisite database
stax db pull --provider=wpengine

# List sites
stax wp site list --provider=wpengine

# Work with specific site
stax wp --url=site1.example.com plugin list --provider=wpengine
```

### Staging to Production Deployment

```bash
# 1. Test on staging
# Configure .stax.yml with environment: staging

# 2. Pull staging to local
stax db pull --provider=wpengine

# 3. Test locally
stax start

# 4. Push to production (via WPEngine portal)
# Portal → Staging → "Push Staging to Production"
```

## Troubleshooting

### SSH Connection Failures

**Error**: "Connection refused" or "Permission denied"

**Solutions**:
1. Verify SSH key is added to WPEngine account
2. Test SSH manually:
   ```bash
   ssh install-name@install-name.ssh.wpengine.net
   ```
3. Check SSH key permissions:
   ```bash
   chmod 600 ~/.ssh/wpengine
   ```
4. Verify install name is correct

### API Authentication Failures

**Error**: "Authentication failed" or "Invalid credentials"

**Solutions**:
1. Verify API credentials in WPEngine portal
2. Check credentials file has correct format
3. Regenerate API credentials if needed
4. Ensure no extra whitespace in credentials

### Database Export Timeouts

**Error**: "Export timed out" or "Connection lost"

**Solutions**:
1. Exclude large tables:
   ```bash
   stax db export --exclude-tables=wp_large_table
   ```
2. Export without logs:
   ```bash
   stax db export --skip-logs
   ```
3. Split export into smaller chunks
4. Contact WPEngine support for very large databases

### File Sync Slow Performance

**Issue**: File sync takes too long

**Solutions**:
1. Use bandwidth limit:
   ```bash
   stax files sync --bandwidth-limit=5000
   ```
2. Sync specific directories:
   ```bash
   stax files sync --include="uploads/**"
   ```
3. Exclude unnecessary files:
   ```bash
   stax files sync --exclude="cache/**,*.log"
   ```
4. Use rsync compression (enabled by default)

## Best Practices

### 1. Use Staging for Testing

Always test database imports and major changes on staging first:
```yaml
provider:
  wpengine:
    environment: staging
```

### 2. Regular Backups

Create manual backups before major changes:
- Via WPEngine portal
- Export to local:
  ```bash
  stax db export "backup-$(date +%Y%m%d).sql" --provider=wpengine
  ```

### 3. Git Workflow

Use Git for deployments:
```bash
# .gitignore
wp-content/uploads/
.sql
*.log
.env
```

```bash
# Deploy
git push wpengine main
```

### 4. Environment-Specific Configuration

Use environment-specific wp-config.php settings:
```php
// wp-config.php
if ( getenv('WPE_ENVTYPE') === 'production' ) {
    define('WP_DEBUG', false);
} else {
    define('WP_DEBUG', true);
}
```

### 5. Monitor Performance

Use WPEngine's built-in tools:
- Site Analytics in portal
- PHP error logs
- Slow query logs
- CDN analytics

### 6. Security

- Rotate API credentials regularly
- Use strong SSH keys (ed25519)
- Keep credentials file permissions at 600
- Never commit credentials to Git

## WPEngine-Specific Features

### Smart Plugin Manager

WPEngine automatically manages certain plugins:
- Object caching (Redis)
- Security plugins
- Performance optimizations

**Note**: Some plugins are restricted or automatically enabled.

### Global Edge Security

Automatic DDoS protection and firewall.

**No configuration needed** - managed by WPEngine.

### Automatic WordPress Updates

WPEngine can automatically update WordPress core.

**Configure** in portal: Auto-updates settings

### PHP Version Management

Switch PHP versions via portal:
- Portal → Install → Configuration → PHP Version

**Supported versions**: 7.4, 8.0, 8.1, 8.2, 8.3

## Additional Resources

- [WPEngine API Documentation](https://wpengineapi.com/)
- [WPEngine SSH Gateway Guide](https://wpengine.com/support/ssh-gateway/)
- [WPEngine Git Push Guide](https://wpengine.com/support/git/)
- [WPEngine Support](https://wpengine.com/support/)

## See Also

- [MULTI_PROVIDER.md](./MULTI_PROVIDER.md) - Multi-provider guide
- [COMMANDS.md](./COMMANDS.md) - Command reference
- [CONFIG_SPEC.md](./CONFIG_SPEC.md) - Configuration specification

---

**Version**: 2.0.0
**Last Updated**: 2025-11-08
**Provider Implementation**: `/Users/geoff/_projects/fc/stax/pkg/providers/wpengine/`

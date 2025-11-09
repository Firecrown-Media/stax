# WPEngine Integration Guide

Complete guide to using Stax with WPEngine hosting.

---

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Database Operations](#database-operations)
- [File Synchronization](#file-synchronization)
- [WPEngine Environments](#wpengine-environments)
- [Remote Media](#remote-media)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## Overview

Stax provides seamless integration with WPEngine, allowing you to:
- Pull databases from production or staging
- Sync files from remote environments
- Access multiple WPEngine environments
- Proxy media files without local downloads
- Work with WPEngine multisite installations

### What is WPEngine?

WPEngine is a managed WordPress hosting platform. It handles:
- Server management and optimization
- Automatic backups
- Security and updates
- CDN and caching
- Staging environments
- Git-based deployments

Stax connects to WPEngine to bring your production/staging data to your local environment.

---

## Getting Started

### Prerequisites

Before using WPEngine features:

1. **WPEngine Account** with API access
2. **WPEngine API Credentials**
3. **SSH Key** added to WPEngine
4. **Install name** (found in WPEngine portal)

### Getting Your WPEngine API Credentials

**Step 1: Access WPEngine Portal**
1. Log in to [my.wpengine.com](https://my.wpengine.com)
2. Click your account name (top right)
3. Go to **Account** → **API Access**

**Step 2: Create API User** (if needed)
1. Click "Add User"
2. Note the username (usually email@company.com)
3. Set a password
4. Save credentials securely

**Step 3: Get Your Install Name**
1. In WPEngine portal, go to **Sites**
2. Click on your site
3. Note the "Install Name" (e.g., `mysite` or `mycompany`)

### Setting Up SSH Access

**Step 1: Generate SSH Key** (if you don't have one)
```bash
ssh-keygen -t ed25519 -C "your_email@example.com" -f ~/.ssh/wpengine
```

Press Enter for no passphrase (or add one for security).

**Step 2: Add Public Key to WPEngine**
```bash
# Copy your public key
cat ~/.ssh/wpengine.pub
```

In WPEngine portal:
1. Go to **Account** → **SSH Keys**
2. Click "Add SSH Key"
3. Paste your public key
4. Give it a name (e.g., "My Development Machine")
5. Click "Add"

**Step 3: Test SSH Connection**
```bash
ssh -i ~/.ssh/wpengine git@git.wpengine.com info
```

You should see a list of your WPEngine installs.

### Configuring Stax

```bash
stax setup
```

Enter your:
- **WPEngine API Username**: your_email@company.com
- **WPEngine API Password**: your_api_password
- **SSH Key Path**: ~/.ssh/wpengine

Stax stores these securely in macOS Keychain.

**Verify setup**:
```bash
stax doctor
```

Should show:
```
✓ WPEngine credentials valid
✓ SSH key configured for WPEngine
```

---

## Database Operations

### Pulling Database from WPEngine

**Basic pull** (from configured environment):
```bash
stax db pull
```

**From specific environment**:
```bash
# Production
stax db pull --environment=production

# Staging
stax db pull --environment=staging
```

**What happens**:
1. Creates automatic snapshot (safety backup)
2. Connects to WPEngine via SSH
3. Detects table prefix (usually `wp_`)
4. Exports database on WPEngine server
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
  Network: mysite.wpengine.com → my-project.local
  site1: site1.com → site1.my-project.local
  site2: site2.com → site2.my-project.local
✓ Flushing WordPress cache

Database pulled successfully!
Time: 2m 34s
```

### Optimizing Database Pulls

**Skip unnecessary data** (faster imports):
```bash
stax db pull \
  --skip-logs \
  --skip-transients \
  --skip-spam
```

This excludes:
- Log tables (action scheduler, error logs)
- Transient data (temporary cache)
- Spam comments and trash

Can reduce import time by 50%+.

**Exclude specific tables**:
```bash
stax db pull --exclude-tables=wp_actionscheduler_logs,wp_wc_admin_notes
```

**No automatic snapshot** (not recommended):
```bash
stax db pull --snapshot=false
```

### Sanitizing Production Data

**For testing with anonymized data**:
```bash
stax db pull --sanitize
```

This:
- Anonymizes user emails
- Resets user passwords
- Removes personal data
- Keeps site structure intact

Good for:
- Testing with sensitive production data
- Sharing databases with contractors
- Compliance with data protection laws

### Database Pull Frequency

**Recommended schedule**:
- **Daily**: If you need fresh content/users
- **Weekly**: For most development work
- **As needed**: When testing with specific production data
- **Never**: For purely local development

**Why not pull constantly?**
- Takes 2-5 minutes per pull
- Loses local database changes
- Creates large snapshots
- Uses bandwidth

---

## File Synchronization

### Syncing Files from WPEngine

Stax can sync files (like uploads) from WPEngine.

**Sync uploads directory**:
```bash
stax provider sync uploads
```

**Sync specific subdirectory**:
```bash
stax provider sync uploads/2024
```

**Sync with dry-run** (see what would sync):
```bash
stax provider sync uploads --dry-run
```

**What happens**:
1. Connects to WPEngine via SSH
2. Uses rsync to transfer files
3. Only syncs changed/new files
4. Preserves file permissions

**From staging**:
```bash
stax provider sync uploads --environment=staging
```

### When to Sync Files

**You should sync when**:
- Testing features that rely on specific uploads
- Debugging media-related issues
- Need exact production file structure

**You shouldn't sync when**:
- Just need to see images (use remote media proxy)
- Uploads are very large (use proxy instead)
- Files rarely change

**Alternative: Remote Media Proxy**
- Stax can proxy media from WPEngine/CDN
- No local download needed
- See "Remote Media" section below

---

## WPEngine Environments

### Available Environments

WPEngine provides multiple environments:

1. **Production**
   - Live site
   - Real users and data
   - Default for `stax db pull`

2. **Staging**
   - Testing environment
   - Copy of production
   - Safe for testing

3. **Development** (some plans)
   - Development environment
   - Isolated from production

### Listing Environments

```bash
stax provider info wpengine
```

Shows:
- Available environments
- Install names
- PHP/MySQL versions
- Domains

### Switching Environments

**Temporary switch** (one command):
```bash
stax db pull --environment=staging
```

**Permanent switch** (update config):
```bash
stax config set wpengine.environment staging
```

Then all pulls use staging by default:
```bash
stax db pull  # Now pulls from staging
```

**Switch back**:
```bash
stax config set wpengine.environment production
```

### Environment-Specific Domains

Production and staging often have different domains:

**Production**:
- mysite.wpengine.com
- customdomain.com

**Staging**:
- mysite-staging.wpengine.com

**Stax handles this automatically**.

In `.stax.yml`:
```yaml
wpengine:
  install: mysite
  environment: production

  domains:
    production:
      primary: mysite.wpengine.com
      sites:
        - site1.com
        - site2.com

    staging:
      primary: mysite-staging.wpengine.com
      sites:
        - staging-site1.com
        - staging-site2.com
```

When you pull from staging, Stax uses staging domains for search-replace.

---

## Remote Media

### How Remote Media Proxy Works

Instead of downloading gigabytes of media files, Stax can proxy them from production.

**How it works**:
1. WordPress requests `/wp-content/uploads/2024/01/image.jpg`
2. File doesn't exist locally
3. nginx proxies request to WPEngine/CDN
4. Image loads in browser
5. Optionally cached locally

**Advantages**:
- No need to download entire uploads directory
- Saves disk space (GBs)
- Always shows current production media
- Faster initial setup

**Disadvantages**:
- Requires internet connection
- Slightly slower than local files (first load)
- Can't test local uploads/modifications

### Enabling Remote Media Proxy

**In `.stax.yml`**:
```yaml
media:
  proxy:
    enabled: true
    remote_url: https://cdn.mysite.com
    # Or: https://mysite.wpengine.com
    fallback_url: https://mysite.wpengine.com
    cache_locally: true
    cache_duration: 7d
```

**Restart to apply**:
```bash
stax restart
```

### Configuration Options

**`enabled`**: Turn proxy on/off
```yaml
enabled: true  # Proxy enabled
enabled: false # Use local files only
```

**`remote_url`**: Primary source for media
```yaml
# Use CDN (fastest)
remote_url: https://cdn.mysite.com

# Or WPEngine directly
remote_url: https://mysite.wpengine.com
```

**`fallback_url`**: If primary fails
```yaml
fallback_url: https://mysite.wpengine.com
```

**`cache_locally`**: Cache downloaded files
```yaml
cache_locally: true  # Downloaded files stay local
cache_locally: false # Always fetch from remote
```

**`cache_duration`**: How long to cache
```yaml
cache_duration: 1d   # 1 day
cache_duration: 7d   # 1 week
cache_duration: 30d  # 1 month
```

### Per-Site Media Proxy (Multisite)

Different sites can proxy from different URLs:

```yaml
network:
  sites:
    - name: site1
      domain: site1.my-project.local
      wpengine_domain: site1.com
      media:
        proxy:
          enabled: true
          remote_url: https://cdn.site1.com

    - name: site2
      domain: site2.my-project.local
      wpengine_domain: site2.com
      media:
        proxy:
          enabled: true
          remote_url: https://cdn.site2.com
```

### Disabling Media Proxy

**Temporarily** (use local only):
```yaml
media:
  proxy:
    enabled: false
```

**Or sync files locally**:
```bash
stax provider sync uploads
```

Then disable proxy.

---

## Best Practices

### Database Management

**Do**:
- Create snapshots before risky operations
- Pull from staging for most development
- Use `--skip-logs` for faster imports
- Pull once a day or less

**Don't**:
- Pull too frequently (loses local changes)
- Pull from production without good reason
- Skip snapshots (no safety net)

### Environment Usage

**Use production when**:
- Debugging production-specific issues
- Testing with real content
- Need exact production data

**Use staging when**:
- Regular development work
- Testing features before production
- Learning/experimenting

### File Synchronization

**Use media proxy when**:
- Uploads are very large
- You don't modify uploads locally
- Internet connection is reliable

**Sync files locally when**:
- Working offline
- Testing upload functionality
- Need fast media loads
- Files are small

### Security

**Do**:
- Keep API credentials secure
- Use SSH keys (not passwords)
- Rotate credentials periodically
- Use `--sanitize` for sensitive data

**Don't**:
- Commit credentials to Git
- Share credentials in plain text
- Use production API user for testing

### Performance

**Optimize pulls**:
```bash
stax db pull \
  --environment=staging \
  --skip-logs \
  --skip-transients \
  --skip-spam
```

**Enable caching**:
```yaml
media:
  proxy:
    cache_locally: true
    cache_duration: 7d
```

---

## Troubleshooting

### Can't Connect to WPEngine

**Symptom**:
```
Error: WPEngine authentication failed
```

**Solutions**:

1. **Verify credentials**:
   ```bash
   stax setup
   # Re-enter credentials
   ```

2. **Test SSH connection**:
   ```bash
   ssh -i ~/.ssh/wpengine git@git.wpengine.com info
   ```

3. **Check SSH key in WPEngine portal**:
   - Account → SSH Keys
   - Verify your public key is listed

4. **Generate new SSH key**:
   ```bash
   ssh-keygen -t ed25519 -f ~/.ssh/wpengine
   cat ~/.ssh/wpengine.pub  # Add to WPEngine
   stax setup  # Update Stax
   ```

### Database Pull Fails

**Symptom**:
```
Error: Failed to pull database
```

**Solutions**:

1. **Check install name**:
   ```bash
   stax config get wpengine.install
   # Should match WPEngine portal
   ```

2. **Test SSH access**:
   ```bash
   ssh -i ~/.ssh/wpengine install-name@install-name.ssh.wpengine.net
   ```

3. **Try different environment**:
   ```bash
   stax db pull --environment=staging
   ```

4. **Reduce database size**:
   ```bash
   stax db pull --skip-logs --skip-transients
   ```

### Media Proxy Not Working

**Symptom**: Images show broken.

**Solutions**:

1. **Check configuration**:
   ```bash
   stax config get media.proxy.enabled
   # Should be: true

   stax config get media.proxy.remote_url
   # Should be valid URL
   ```

2. **Test URL manually**:
   ```bash
   curl -I https://cdn.mysite.com/wp-content/uploads/2024/01/test.jpg
   # Should return 200 OK
   ```

3. **Check nginx config**:
   ```bash
   stax ssh
   cat /etc/nginx/sites-enabled/wordpress.conf | grep proxy_pass
   exit
   ```

4. **Disable and use local**:
   ```yaml
   media:
     proxy:
       enabled: false
   ```

   Then sync files:
   ```bash
   stax provider sync uploads
   ```

### Wrong Environment

**Symptom**: Pulled wrong environment's data.

**Solution**:
```bash
# Check current setting
stax config get wpengine.environment

# Restore from snapshot
stax db list
stax db restore before-pull

# Pull from correct environment
stax db pull --environment=production
```

---

## Next Steps

- **User Guide**: [USER_GUIDE.md](./USER_GUIDE.md) - General usage
- **Multisite**: [MULTISITE.md](./MULTISITE.md) - Multisite with WPEngine
- **Examples**: [EXAMPLES.md](./EXAMPLES.md) - Real-world workflows
- **Troubleshooting**: [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - More solutions

---

**Questions?** Check the [FAQ](./FAQ.md) or contact your team!

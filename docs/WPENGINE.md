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

### Understanding WPEngine API Access

Before diving into credential setup, it's helpful to understand what you're configuring and why.

**What are WPEngine API credentials?**

WPEngine API credentials are special login credentials that allow external applications (like Stax) to programmatically access your WPEngine account and perform operations like pulling databases, accessing files, and managing environments. These credentials are separate from your regular WPEngine portal login.

**Why do I need them?**

Stax uses these credentials to authenticate with WPEngine's API and SSH gateway, enabling features like database pulls, file synchronization, and environment management. Without valid credentials, Stax cannot connect to your WPEngine hosting.

**Account requirements and permissions**

Not all WPEngine users can enable or use API access:

- **Owner role**: Required to enable API access for the account (one-time setup)
- **Full/Partial users**: Can use API access once enabled, with access to their assigned installs
- **Billing-only users**: Cannot access API features

Learn more about [WPEngine user roles and permissions](https://wpengine.com/support/users/).

**Security best practices**

- Store credentials securely (Stax uses macOS Keychain)
- Never commit credentials to version control
- Use unique passwords for API access
- Rotate credentials periodically
- Grant API access only to users who need it
- Consider creating separate API users for different team members or purposes

Official WPEngine API documentation:
- [WPEngine API Documentation](https://wpengineapi.com/)
- [Enabling API Access Guide](https://wpengine.com/support/enabling-wp-engine-api/)
- [Developer's Guide to WPEngine API](https://wpengine.com/builders/mastering-the-wp-engine-api-a-comprehensive-guide-for-developers/)

### Getting Your WPEngine API Credentials

Follow these four steps to set up your WPEngine API credentials for use with Stax.

#### Step 1: Verify Account Permissions

Before attempting to enable or use API access, confirm you have the appropriate permissions.

**Check your role**:
1. Log in to the [WPEngine User Portal](https://my.wpengine.com)
2. Click your account name in the top right corner
3. Go to **Account** → **Users**
4. Find your username and note your role

**Required roles**:
- To **enable API access**: Must have Owner role
- To **use API access**: Owner, Full, or Partial user roles

If you don't have Owner privileges and API access is not enabled, you'll need to:
- Contact your account owner to enable API access
- Or request role elevation from your organization's WPEngine administrator

Learn more: [WPEngine Account Users and Roles](https://wpengine.com/support/users/)

#### Step 2: Enable API Access (if not already enabled)

If your account doesn't have API access enabled yet, an Owner must enable it first.

**Check if API access is enabled**:
1. Log in to [my.wpengine.com](https://my.wpengine.com)
2. Click your account name (top right)
3. Look for **API Access** in the Account menu
4. If you see "API Access", it's already enabled
5. If you don't see it, it needs to be enabled

**Enable API access** (Owner only):
1. In the WPEngine portal, go to **Account Settings**
2. Find the **API Access** section
3. Click **Enable API Access** or **Turn On API Access**
4. Confirm the action

Note: Only account Owners can enable API access. If you're not an Owner and don't see this option, contact your account owner or [WPEngine support](https://help.wpengine.com/requests).

Official guide: [Enabling WP Engine API](https://wpengine.com/support/enabling-wp-engine-api/)

#### Step 3: Generate API Credentials

Once API access is enabled, any eligible user can create their API credentials.

**Create API user credentials**:
1. Log in to [my.wpengine.com](https://my.wpengine.com)
2. Click your account name (top right)
3. Go to **Account** → **API Access**
4. You'll see API access settings for your user

**Your API credentials are**:
- **Username**: Usually your WPEngine login email (e.g., `your_email@company.com`)
- **Password**: Set or view in the API Access section
  - Click "Reset Password" if you need to create or change it
  - Enter a strong, unique password
  - Click "Save" or "Update Password"

**Important: Save your credentials immediately**

Write down or securely store:
- Your API username (email)
- Your API password

You'll need these for Step 4 when configuring Stax. WPEngine may not show the password again after you close the page.

**Get your install name**:

You'll also need to know your WPEngine install name:
1. In WPEngine portal, go to **Sites**
2. Click on your site
3. Note the **Install Name** (e.g., `mysite` or `mycompany`)
4. This is different from your domain name

The install name is used to connect to the correct WPEngine environment.

#### Step 4: Verify API Access

Test your credentials to ensure they work before configuring Stax.

**Test API credentials** (optional but recommended):

You can verify your API access using WPEngine's API directly:

```bash
# Test API authentication
curl -u "your_email@company.com:your_api_password" \
  https://api.wpengineapi.com/v1/installs
```

Expected response: JSON list of your WPEngine installs.

If you see an authentication error:
- Double-check your username and password
- Ensure API access is enabled for your account
- Verify you have access to at least one install
- Try resetting your API password

**Configure Stax with credentials**:

Once verified, configure Stax (see "Configuring Stax" section below):

```bash
stax setup
```

Enter your:
- **WPEngine API Username**: your_email@company.com
- **WPEngine API Password**: your_api_password
- **SSH Key Path**: ~/.ssh/wpengine (covered in next section)

**What if verification fails?**

If you can't authenticate with the API:
1. Verify API access is enabled (Owner must enable it)
2. Check that your role has API permissions (not billing-only)
3. Ensure your password is correct (try resetting it)
4. Confirm you have access to at least one install
5. Contact [WPEngine support](https://help.wpengine.com/requests) if issues persist

Common issues are covered in the "Common Credential Issues" section in Troubleshooting.

### Setting Up SSH Access

SSH (Secure Shell) access allows Stax to securely connect to your WPEngine environments for operations like database pulls and file synchronization. You'll need to generate an SSH key pair and add the public key to your WPEngine account.

Official WPEngine SSH documentation:
- [SSH Keys for Shell Access](https://wpengine.com/support/ssh-keys-for-shell-access/)
- [SSH Gateway Guide](https://wpengine.com/support/ssh-gateway/)
- [Manage SSH Keys in Portal](https://my.wpengine.com/profile/ssh_keys)

**Step 1: Generate SSH Key** (if you don't have one)

Create a dedicated SSH key for WPEngine access:

```bash
ssh-keygen -t ed25519 -C "your_email@example.com" -f ~/.ssh/wpengine
```

When prompted:
- **Passphrase**: Press Enter for no passphrase (convenient) or enter one (more secure)
- If you use a passphrase, you'll need to enter it each time or use ssh-agent

This creates two files:
- `~/.ssh/wpengine` - Private key (keep secret, never share)
- `~/.ssh/wpengine.pub` - Public key (add to WPEngine)

**Step 2: Add Public Key to WPEngine**

Copy your public key to your clipboard:

```bash
# macOS
cat ~/.ssh/wpengine.pub | pbcopy

# Or just display it to copy manually
cat ~/.ssh/wpengine.pub
```

Add the key to WPEngine:
1. Log in to [my.wpengine.com](https://my.wpengine.com)
2. Go to **Your Profile** → **SSH Keys** (or click [here](https://my.wpengine.com/profile/ssh_keys))
3. Click **Add SSH Key**
4. Paste your public key (entire contents of `wpengine.pub`)
5. Give it a descriptive name (e.g., "MacBook Pro - Development")
6. Click **Add** or **Save**

The key should appear in your SSH keys list immediately.

**Step 3: Test SSH Connection**

Verify your SSH key works with WPEngine:

```bash
ssh -i ~/.ssh/wpengine git@git.wpengine.com info
```

Expected output: List of your WPEngine installs with their names and IDs.

Example:
```
Available repositories:
  mysite
  mysite-staging
  another-project
```

If you see this list, your SSH key is configured correctly!

#### Troubleshooting SSH Connection Issues

**"Permission denied (publickey)"**

This means WPEngine didn't accept your SSH key.

Solutions:
1. Verify the public key is added to WPEngine portal
2. Check you're using the correct private key file:
   ```bash
   ssh -i ~/.ssh/wpengine git@git.wpengine.com info
   ```
3. Ensure the key format is correct (should start with `ssh-ed25519` or `ssh-rsa`)
4. Try removing and re-adding the key in WPEngine portal
5. Generate a new key pair if the key is corrupted

**"Connection timed out" or "Could not resolve hostname"**

This indicates a network connectivity issue.

Solutions:
1. Check your internet connection
2. Verify you can reach WPEngine: `ping git.wpengine.com`
3. Check if a firewall or VPN is blocking SSH (port 22)
4. Try from a different network

**Key works in portal but not with Stax**

Solutions:
1. Verify the key path in Stax configuration:
   ```bash
   stax config get ssh.key_path
   ```
2. Re-run setup with correct path:
   ```bash
   stax setup
   # Enter: ~/.ssh/wpengine
   ```
3. Check file permissions (should be 600 for private key):
   ```bash
   chmod 600 ~/.ssh/wpengine
   chmod 644 ~/.ssh/wpengine.pub
   ```

**Multiple SSH keys for different WPEngine accounts**

If you work with multiple WPEngine accounts, create separate keys:

```bash
# Account 1
ssh-keygen -t ed25519 -f ~/.ssh/wpengine_client1

# Account 2
ssh-keygen -t ed25519 -f ~/.ssh/wpengine_client2
```

Configure Stax per project:
```bash
cd project1
stax setup
# SSH Key Path: ~/.ssh/wpengine_client1

cd project2
stax setup
# SSH Key Path: ~/.ssh/wpengine_client2
```

For more SSH help, see [WPEngine SSH Gateway documentation](https://wpengine.com/support/ssh-gateway/).

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

### Common Credential Issues

This section covers common problems users encounter when setting up or using WPEngine credentials with Stax.

#### Issue 1: API Access Not Enabled

**Symptom**:
```
Error: WPEngine API authentication failed
Error: 401 Unauthorized
```

Or when checking credentials:
```bash
curl -u "email@example.com:password" https://api.wpengineapi.com/v1/installs
# Returns: 401 Unauthorized
```

**Cause**: API access is not enabled for your WPEngine account. Only account Owners can enable API access.

**Solution**:

1. **Check if API access is enabled**:
   - Log in to [my.wpengine.com](https://my.wpengine.com)
   - Click your account name (top right)
   - Look for "API Access" in the menu
   - If you don't see it, API access is not enabled

2. **Enable API access** (Owner only):
   - Go to **Account Settings**
   - Find **API Access** section
   - Click **Enable API Access**
   - If you're not an Owner, contact your account owner to enable it

3. **Verify your role**:
   - Go to **Account** → **Users**
   - Find your username and check your role
   - Only Owner, Full, and Partial users can use API access
   - Billing-only users cannot access the API

**Documentation**: [Enabling WP Engine API](https://wpengine.com/support/enabling-wp-engine-api/)

#### Issue 2: Wrong Credential Format

**Symptom**:
```
Error: Authentication failed
Error: Invalid username or password
```

**Cause**: Using the wrong format for username. Common mistakes include:
- Using WPEngine account name instead of email
- Using install name as username
- Adding extra characters or spaces

**Solution**:

1. **Use the correct username format**:
   - Username should be your WPEngine login email (e.g., `your_email@company.com`)
   - NOT your install name
   - NOT your account name
   - NOT your domain name

2. **Verify credentials in portal**:
   - Log in to [my.wpengine.com](https://my.wpengine.com)
   - Go to **Account** → **API Access**
   - Your API username is shown there (usually your email)
   - Reset your API password if needed

3. **Test credentials directly**:
   ```bash
   curl -u "your_email@company.com:your_password" \
     https://api.wpengineapi.com/v1/installs
   ```

   Should return JSON list of installs, not 401 error.

4. **Re-configure Stax**:
   ```bash
   stax setup
   # Enter email exactly as shown in portal
   # Copy/paste password to avoid typos
   ```

#### Issue 3: API User Lacks Access to Specific Install

**Symptom**:
```
Error: Install not found
Error: You do not have access to this install
```

Or install doesn't appear in API response:
```bash
curl -u "email@company.com:password" https://api.wpengineapi.com/v1/installs
# Install you need is not in the list
```

**Cause**: Your API user account doesn't have permissions to access the specific WPEngine install you're trying to pull from.

**Solution**:

1. **Check your install access**:
   - Log in to [my.wpengine.com](https://my.wpengine.com)
   - Go to **Sites**
   - Note which installs are listed
   - You can only access installs visible to you

2. **Verify install name is correct**:
   ```bash
   stax config get wpengine.install
   # Should match exactly (case-sensitive)
   ```

   In WPEngine portal:
   - Go to **Sites** → Click your site
   - Note the exact install name (e.g., `mysite`, not `My Site`)

3. **Request access to the install**:
   - Contact your WPEngine account owner
   - Ask them to grant you access to the specific install
   - Owner can manage user access in **Account** → **Users**

4. **Use a different install**:
   ```bash
   # List installs you can access
   ssh -i ~/.ssh/wpengine git@git.wpengine.com info

   # Update Stax to use accessible install
   stax config set wpengine.install accessible-install-name
   ```

**Documentation**: [WPEngine Account Users and Roles](https://wpengine.com/support/users/)

#### Issue 4: Credentials Expired or Revoked

**Symptom**:
```
Error: Authentication failed
Error: Invalid credentials
```

Previously working credentials suddenly stop working.

**Cause**: API credentials were changed, reset, or revoked in the WPEngine portal by you or another administrator.

**Solution**:

1. **Reset your API password**:
   - Log in to [my.wpengine.com](https://my.wpengine.com)
   - Go to **Account** → **API Access**
   - Click **Reset Password**
   - Set a new password
   - Save it securely

2. **Update Stax with new credentials**:
   ```bash
   stax setup
   # Enter your email (username)
   # Enter new password
   ```

3. **Verify new credentials work**:
   ```bash
   stax doctor
   # Should show: ✓ WPEngine credentials valid
   ```

4. **Check if user was deactivated**:
   - Have an Owner check **Account** → **Users**
   - Ensure your account is still active
   - If deactivated, Owner must reactivate it

#### Issue 5: Network or Firewall Blocking API Calls

**Symptom**:
```
Error: Connection timeout
Error: Could not connect to WPEngine API
Error: Network unreachable
```

**Cause**: Your network, firewall, VPN, or security software is blocking connections to WPEngine's API or SSH gateway.

**Solution**:

1. **Test basic connectivity**:
   ```bash
   # Test API endpoint
   curl -I https://api.wpengineapi.com
   # Should return 200 or 401, not connection error

   # Test SSH gateway
   ping git.wpengine.com
   # Should show responses, not timeouts
   ```

2. **Check firewall rules**:
   - Ensure outbound HTTPS (port 443) is allowed
   - Ensure outbound SSH (port 22) is allowed
   - Whitelist these domains:
     - `api.wpengineapi.com`
     - `*.wpengine.com`
     - `*.wpengine.net`

3. **Try different network**:
   ```bash
   # Temporarily disable VPN
   # Try from different WiFi network
   # Test from phone hotspot
   stax db pull
   ```

   If it works on different network, your primary network has restrictions.

4. **Check corporate proxy/VPN**:
   - Corporate networks often block SSH or restrict API access
   - Contact your IT department
   - Request access to WPEngine domains
   - Or use a VPN that allows these connections

5. **DNS issues**:
   ```bash
   # Test DNS resolution
   nslookup api.wpengineapi.com
   nslookup git.wpengine.com

   # Try different DNS server temporarily
   # Google DNS: 8.8.8.8
   # Cloudflare DNS: 1.1.1.1
   ```

#### Issue 6: SSH Key Not Recognized

**Symptom**:
```
Error: Permission denied (publickey)
Error: Could not authenticate with SSH key
```

**Cause**: SSH key is not properly configured in WPEngine or Stax, or the key format is incorrect.

**Solution**:

1. **Verify key is in WPEngine portal**:
   - Go to [my.wpengine.com/profile/ssh_keys](https://my.wpengine.com/profile/ssh_keys)
   - Confirm your public key is listed
   - Key should start with `ssh-ed25519` or `ssh-rsa`

2. **Test SSH key directly**:
   ```bash
   ssh -i ~/.ssh/wpengine git@git.wpengine.com info
   ```

   Should show list of installs, not "Permission denied".

3. **Check key permissions**:
   ```bash
   ls -l ~/.ssh/wpengine*
   # Private key should be: -rw------- (600)
   # Public key should be: -rw-r--r-- (644)

   # Fix if needed:
   chmod 600 ~/.ssh/wpengine
   chmod 644 ~/.ssh/wpengine.pub
   ```

4. **Verify Stax is using correct key**:
   ```bash
   stax config get ssh.key_path
   # Should show: /Users/yourname/.ssh/wpengine

   # Update if wrong:
   stax setup
   # Enter correct path
   ```

5. **Regenerate key if corrupted**:
   ```bash
   # Backup old key
   mv ~/.ssh/wpengine ~/.ssh/wpengine.old
   mv ~/.ssh/wpengine.pub ~/.ssh/wpengine.pub.old

   # Generate new key
   ssh-keygen -t ed25519 -f ~/.ssh/wpengine

   # Add to WPEngine
   cat ~/.ssh/wpengine.pub
   # Copy and add to portal

   # Update Stax
   stax setup
   ```

**Documentation**: [SSH Keys for Shell Access](https://wpengine.com/support/ssh-keys-for-shell-access/)

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

### Getting Help with WPEngine Credentials

If you've tried the troubleshooting steps above and still can't get your credentials working, here's how to get help.

#### When to Contact WPEngine Support vs Stax Team

**Contact WPEngine Support for**:
- Enabling API access (if you're not an Owner)
- Resetting forgotten API passwords
- Account permission issues
- SSH key not being accepted by WPEngine
- Install access problems
- Questions about user roles and permissions
- WPEngine account or billing issues

**Contact Stax Team or Your Team Lead for**:
- Stax configuration issues
- How to configure `.stax.yml`
- Stax command errors
- Local environment problems
- Project-specific setup questions

#### WPEngine Support Channels

**WPEngine Support Center**:
- Browse articles: [wpengine.com/support](https://wpengine.com/support/)
- Search for credential setup guides
- Find SSH and API documentation

**WPEngine Ticketing System**:
- Create a support ticket: [help.wpengine.com/requests](https://help.wpengine.com/requests)
- Response time: Usually within 24 hours
- For urgent issues, mention it in the ticket

**When creating a ticket, include**:
- Your WPEngine account name
- Your install name
- Your email address (username)
- Description of the issue
- What you've already tried
- Any error messages (exact text)
- Screenshots if helpful

**Topics WPEngine can help with**:
- "API access not showing in my account"
- "Need API access enabled for my account"
- "SSH key not working after adding to portal"
- "Cannot access specific install with my credentials"
- "User role change request"

#### Common Credential Reset Procedures

**Reset API Password**:
1. Log in to [my.wpengine.com](https://my.wpengine.com)
2. Go to **Account** → **API Access**
3. Click **Reset Password**
4. Enter new password
5. Click **Save** or **Update Password**
6. Update Stax: `stax setup`

**Reset SSH Key**:
1. Go to [my.wpengine.com/profile/ssh_keys](https://my.wpengine.com/profile/ssh_keys)
2. Delete old key (if needed)
3. Generate new key: `ssh-keygen -t ed25519 -f ~/.ssh/wpengine`
4. Add new public key to portal
5. Update Stax: `stax setup`

**Request Access to Install**:
1. Contact your WPEngine account Owner
2. Provide your email address
3. Provide the install name you need access to
4. Owner grants access via **Account** → **Users**
5. Verify access: `ssh -i ~/.ssh/wpengine git@git.wpengine.com info`

#### Emergency Access Issues

**Locked out of WPEngine portal**:
- Go to [my.wpengine.com](https://my.wpengine.com)
- Click "Forgot Password?"
- Follow password reset email

**No one on your team has Owner access**:
- Contact [WPEngine support](https://help.wpengine.com/requests)
- Verify account ownership (may need billing information)
- Request Owner role assignment

**Entire account suspended or access revoked**:
- Contact WPEngine support immediately
- Usually due to billing or security issues
- Resolve with WPEngine directly

#### Additional Resources

- [WPEngine User Portal](https://my.wpengine.com) - Manage credentials
- [WPEngine Users Guide](https://wpengine.com/support/users/) - User roles
- [Enabling API Access](https://wpengine.com/support/enabling-wp-engine-api/) - Setup guide
- [SSH Key Management](https://wpengine.com/support/ssh-keys-for-shell-access/) - SSH setup
- [WPEngine API Documentation](https://wpengineapi.com/) - API reference

---

## Next Steps

- **User Guide**: [USER_GUIDE.md](./USER_GUIDE.md) - General usage
- **Multisite**: [MULTISITE.md](./MULTISITE.md) - Multisite with WPEngine
- **Examples**: [EXAMPLES.md](./EXAMPLES.md) - Real-world workflows
- **Troubleshooting**: [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - More solutions

---

**Questions?** Check the [FAQ](./FAQ.md) or contact your team!

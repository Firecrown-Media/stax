# Stax Quick Start Guide

Get your first WordPress multisite project up and running in 5 minutes or less.

---

## Before You Start

Make sure you have:
- [ ] Stax installed ([Installation Guide](./INSTALLATION.md))
- [ ] Docker Desktop running
- [ ] WPEngine API credentials
- [ ] GitHub access token (if using private repos)

**Quick check**:
```bash
stax doctor
```

If all checks pass, you're ready to go!

---

## Your First Project in 5 Steps

### Step 1: Configure Your Credentials (One-Time Setup)

This only needs to be done once on your machine.

```bash
stax setup
```

You'll be asked for:

**WPEngine API Username**: (e.g., `yourname@firecrown.com`)
```
? WPEngine API Username: yourname@firecrown.com
```

**WPEngine API Password**: (your WPEngine API password)
```
? WPEngine API Password: ********
```

**GitHub Token**: (optional, press Enter to skip if not using private repos)
```
? GitHub Personal Access Token (optional): ghp_xxxxxxxxxxxxx
```

**SSH Key**: (path to your WPEngine SSH key)
```
? SSH Key for WPEngine: ~/.ssh/wpengine
```

**Expected output**:
```
✓ Validating WPEngine credentials
✓ Validating GitHub token
✓ Saving credentials to macOS Keychain

Credentials saved successfully!
```

### Step 2: Create a Project Directory

```bash
# Create a directory for your project
mkdir -p ~/Sites/my-multisite
cd ~/Sites/my-multisite
```

You can use any directory you prefer. Stax works wherever you are.

### Step 3: Initialize Your Project

```bash
stax init
```

Stax will now guide you through an interactive setup. Here's what you'll be asked:

**Project name**: (defaults to directory name)
```
? Project name: my-multisite
```

**Multisite mode**: (subdomain or subdirectory)
```
? Multisite mode:
  ❯ subdomain
    subdirectory
```
Choose `subdomain` if your production site uses subdomains (e.g., site1.example.com).
Choose `subdirectory` if it uses paths (e.g., example.com/site1).

**Network domain**: (your local domain)
```
? Network domain: my-multisite.local
```
This will be your main network URL. Stax automatically configures SSL for `https://`.

**WPEngine install name**: (your WPEngine install)
```
? WPEngine install name: myinstall
```
Find this in your WPEngine portal under "Sites".

**Environment**: (production or staging)
```
? Environment to pull from:
  ❯ production
    staging
```
Choose which WPEngine environment to sync from.

**GitHub repository**: (your codebase repository)
```
? GitHub repository URL: https://github.com/mycompany/my-project.git
```

**Add sites**: Now add your subsites:
```
? Add a site? Yes

? Site name: site1
? Local domain: site1.my-multisite.local
? WPEngine domain: site1.example.com

✓ Added site: site1

? Add another site? Yes

? Site name: site2
? Local domain: site2.my-multisite.local
? WPEngine domain: site2.example.com

✓ Added site: site2

? Add another site? No
```

**What happens next**:

Stax will now:
1. Clone your repository from GitHub
2. Detect PHP/MySQL versions from WPEngine
3. Generate DDEV configuration
4. Start Docker containers
5. Install dependencies (Composer + npm)
6. Run build scripts
7. Pull database from WPEngine
8. Import and configure the database
9. Run search-replace for all sites

**Expected output** (this takes 2-5 minutes):
```
🚀 Initializing Stax project: my-multisite

✓ Validating WPEngine credentials
✓ Cloning repository from GitHub
  Repository: https://github.com/mycompany/my-project.git
  Branch: main
✓ Detecting versions from WPEngine
  PHP: 8.1
  MySQL: 8.0
  WordPress: 6.4.2
✓ Generating DDEV configuration
✓ Starting DDEV containers
  Web: https://my-multisite.local
  Database: MySQL 8.0
  MailHog: http://my-multisite.local:8025
✓ Installing Composer dependencies
  Installed 45 packages
✓ Installing NPM dependencies
  Added 234 packages
✓ Running build script
  Compiled assets
✓ Pulling database from WPEngine
  Size: 245 MB
  Tables: 127
✓ Importing database
  Rows imported: 1,245,678
✓ Running search-replace
  Network: myinstall.wpengine.com → my-multisite.local (43 replacements)
  site1: site1.example.com → site1.my-multisite.local (1,234 replacements)
  site2: site2.example.com → site2.my-multisite.local (987 replacements)
✓ Flushing WordPress cache

✓ Project initialized successfully!

Your sites are ready:
  Network:  https://my-multisite.local
  site1:    https://site1.my-multisite.local
  site2:    https://site2.my-multisite.local

  MailHog:  http://my-multisite.local:8025

Next steps:
  stax status       Check environment status
  stax ssh          SSH into web container
  stax logs -f      View logs
  stax db snapshot  Create database backup

WordPress Admin:
  URL:      https://my-multisite.local/wp-admin
  Users:    Your production users (unchanged)
```

### Step 4: Access Your Site

Open your browser and go to:
- **Network**: https://my-multisite.local
- **Site 1**: https://site1.my-multisite.local
- **Site 2**: https://site2.my-multisite.local

You'll see your WordPress site running locally with production data!

**Log in to WordPress**:
```
URL: https://my-multisite.local/wp-admin
Username: (your production admin username)
Password: (your production admin password)
```

### Step 5: Start Developing

Your environment is now running! Here are some commands to try:

**Check status**:
```bash
stax status
```

**SSH into the container**:
```bash
stax ssh
# You're now inside the container
wp plugin list
composer --version
npm --version
exit
```

**View logs**:
```bash
stax logs -f
# Press Ctrl+C to stop
```

**Make some changes**:
```bash
# Edit a file in your project
code .  # or your preferred editor

# Watch your changes
stax dev  # Starts build watcher
```

**When you're done**:
```bash
stax stop
```

**Tomorrow, start again**:
```bash
stax start  # Takes ~10 seconds
```

---

## What Just Happened?

Let's understand what Stax set up for you:

### 1. Directory Structure

```
~/Sites/my-multisite/
├── .stax.yml              # Your project configuration
├── .ddev/                 # DDEV configuration (auto-generated)
│   └── config.yaml
├── wp-admin/              # WordPress core
├── wp-content/            # Your themes and plugins
│   ├── themes/
│   ├── plugins/
│   ├── mu-plugins/
│   └── uploads/          # Media (proxied from production)
├── wp-config.php          # WordPress configuration
├── composer.json          # PHP dependencies
├── package.json           # JavaScript dependencies
└── ... (your other files)
```

### 2. Running Containers

Stax started several Docker containers:

- **Web container** (`ddev-<project>-web`):
  - Runs PHP 8.1 and nginx
  - Contains WordPress, WP-CLI, Composer, npm
  - This is where your code runs

- **Database container** (`ddev-<project>-db`):
  - Runs MySQL 8.0
  - Contains your imported database
  - Accessible from web container

- **Router container** (`ddev-router`):
  - Handles incoming requests
  - Routes to correct project
  - Manages SSL certificates
  - Shared across all DDEV projects

- **MailHog container**:
  - Catches all outgoing emails
  - View at http://my-multisite.local:8025
  - Prevents accidentally sending emails

### 3. Configuration Files

**`.stax.yml`** (Project configuration):
```yaml
version: 1

project:
  name: my-multisite
  type: wordpress-multisite
  mode: subdomain

wpengine:
  install: myinstall
  environment: production

network:
  domain: my-multisite.local
  sites:
    - name: site1
      domain: site1.my-multisite.local
      wpengine_domain: site1.example.com
    - name: site2
      domain: site2.my-multisite.local
      wpengine_domain: site2.example.com

ddev:
  php_version: "8.1"
  mysql_version: "8.0"

repository:
  url: https://github.com/mycompany/my-project.git
  branch: main
```

This file is version-controlled with your project so your whole team uses identical settings.

**`.ddev/config.yaml`** (DDEV configuration - auto-generated by Stax):
```yaml
name: my-multisite
type: wordpress
php_version: "8.1"
webserver_type: nginx-fpm
router_http_port: "80"
router_https_port: "443"
additional_fqdns:
  - "*.my-multisite.local"
  - site1.my-multisite.local
  - site2.my-multisite.local
```

You generally don't need to edit this - Stax regenerates it from `.stax.yml`.

---

## Common Tasks

Now that your environment is running, here are some common things you'll want to do:

### Refresh the Database

Pull the latest database from WPEngine:

```bash
stax db pull
```

This will:
1. Create a snapshot of your current database (automatic backup)
2. Download the latest database from WPEngine
3. Import it
4. Run search-replace
5. Flush caches

**Pull from staging instead**:
```bash
stax db pull --environment=staging
```

### Create a Database Snapshot

Before doing something risky, create a snapshot:

```bash
stax db snapshot before-testing
```

**Restore it later**:
```bash
stax db restore before-testing
```

**List all snapshots**:
```bash
stax db list
```

### Run WordPress Commands

Use WP-CLI for any WordPress operations:

```bash
# List plugins
stax wp plugin list

# List all sites
stax wp site list

# Flush cache
stax wp cache flush

# Create a user
stax wp user create newuser newuser@example.com --role=administrator

# Run search-replace
stax wp search-replace old-url.com new-url.com --dry-run
```

### Rebuild Assets

After changing CSS or JavaScript:

```bash
stax build
```

Or run in watch mode:
```bash
stax dev
```

### View Emails

WordPress sends emails to MailHog instead of real addresses:

1. Open http://my-multisite.local:8025
2. See all sent emails
3. Click to view full email content

Great for testing password resets, notifications, etc.

### Stop and Start

**Stop the environment**:
```bash
stax stop
```
Containers stop but data is preserved.

**Start again**:
```bash
stax start
```
Containers restart in ~10 seconds.

**Restart** (stop + start):
```bash
stax restart
```

**Check status**:
```bash
stax status
```

---

## Troubleshooting

### "Can't access my site"

**Check if containers are running**:
```bash
stax status
```

**If stopped, start them**:
```bash
stax start
```

**If running but still can't access**:
```bash
# Restart
stax restart

# Check logs
stax logs -f
```

### "SSL certificate error"

Browsers may warn about the certificate on first access. This is normal for local development.

**In Chrome/Edge**: Click "Advanced" → "Proceed to site (unsafe)"
**In Firefox**: Click "Advanced" → "Accept the Risk and Continue"
**In Safari**: Click "Show Details" → "visit this website"

The certificate is valid - browsers just don't trust local certificates by default.

### "Subdomain not accessible"

**Make sure you used the full domain**:
```
✗ site1.my-multisite          # Missing .local
✓ site1.my-multisite.local    # Correct
```

**Check your configuration**:
```bash
stax config list
```

**Restart to regenerate config**:
```bash
stax restart
```

### "Database import failed"

**Check WPEngine credentials**:
```bash
stax setup
```

**Try pulling again**:
```bash
stax db pull
```

**Check logs for details**:
```bash
stax logs -f
```

### "Port already in use"

Another service is using port 80 or 443.

**Find what's using it**:
```bash
sudo lsof -i :80
sudo lsof -i :443
```

**Common culprits**:
- Apache: `sudo apachectl stop`
- Other DDEV projects: `ddev poweroff`
- Docker containers: `docker stop $(docker ps -q)`

### "Build failed"

**Check if dependencies installed**:
```bash
stax ssh
composer install
npm install
exit
```

**Run build manually**:
```bash
stax build
```

**Check for errors**:
```bash
stax logs -f
```

---

## Next Steps

Now that you have a working environment, learn more:

1. **User Guide** ([docs/USER_GUIDE.md](./USER_GUIDE.md))
   - Daily workflows
   - Advanced features
   - Best practices

2. **Multisite Guide** ([docs/MULTISITE.md](./MULTISITE.md))
   - Understanding multisite
   - Managing subsites
   - Troubleshooting multisite issues

3. **Command Reference** ([docs/COMMAND_REFERENCE.md](./COMMAND_REFERENCE.md))
   - All available commands
   - Flags and options
   - Examples

4. **Real-World Examples** ([docs/EXAMPLES.md](./EXAMPLES.md))
   - Daily development workflow
   - Testing workflows
   - Team collaboration
   - Emergency recovery

---

## Quick Reference

### Most Common Commands

```bash
# Environment
stax start                    # Start environment
stax stop                     # Stop environment
stax restart                  # Restart environment
stax status                   # Show status
stax ssh                      # SSH into container

# Database
stax db pull                  # Pull from WPEngine
stax db snapshot <name>       # Create snapshot
stax db restore <name>        # Restore snapshot
stax db list                  # List snapshots

# Development
stax build                    # Build assets
stax dev                      # Watch mode
stax wp <command>             # WordPress commands
stax logs -f                  # View logs

# Configuration
stax config list              # Show config
stax doctor                   # Diagnose issues
```

### Important URLs

- **Main site**: https://my-multisite.local
- **WordPress admin**: https://my-multisite.local/wp-admin
- **MailHog**: http://my-multisite.local:8025
- **phpMyAdmin**: `ddev launch -p` (in project directory)

### Important Files

- **Project config**: `.stax.yml`
- **DDEV config**: `.ddev/config.yaml`
- **WordPress config**: `wp-config.php`
- **Logs**: `~/.stax/logs/stax.log`

---

**You're all set!** Start building amazing WordPress sites with Stax.

**Need help?** See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) or contact your team.

# Stax Frequently Asked Questions

Common questions and answers about Stax.

---

## General Questions

### What is Stax?

Stax is a command-line tool that automates WordPress development for both single sites and multisite networks. It replaces LocalWP for teams that need automated setup, WPEngine integration, and consistent team environments.

Think of it as: LocalWP + automation + team collaboration + hosting integration.

### Who should use Stax?

**Stax is perfect for**:
- WordPress developers (single site or multisite)
- Teams using WPEngine hosting
- Anyone who wants automated, consistent local development
- Developers comfortable with command-line tools

**Stax might not be for you if**:
- You prefer GUI tools over command-line
- You're not using macOS (Windows/Linux not supported yet)

### How is Stax different from LocalWP?

| Feature | LocalWP | Stax |
|---------|---------|------|
| **Setup time** | 10-30 minutes | 2-5 minutes |
| **Automation** | Manual steps | Fully automated |
| **Configuration** | GUI clicks | Version-controlled YAML |
| **Team consistency** | Variable | Identical environments |
| **WPEngine sync** | Manual export/import | One command |
| **Multisite subdomains** | Manual hosts file | Automatic |
| **Search-replace** | Manual | Automatic |
| **CLI support** | Limited | Full CLI |

**Bottom line**: Stax is faster, more automated, and better for teams.

### Can I use Stax for single-site WordPress?

**Absolutely!** Stax works great for standard single-site WordPress installations. In fact, the default project type is now `wordpress` (single site) rather than `wordpress-multisite`.

All core features work perfectly for single sites:
- Automatic database sync from WPEngine
- Remote media proxying
- Database snapshots and restore
- Build automation
- Team-friendly configuration

**When to use single site vs. multisite:**
- **Use single site** for most WordPress projects (blogs, client sites, marketing sites, etc.)
- **Use multisite** only when you need multiple sites sharing code and database

If you're not sure whether you need multisite, you probably don't. Most WordPress projects are single sites.

### How is Stax different from wp-env or other tools?

**vs wp-env** (WordPress's official tool):
- Stax supports both single sites and multisite networks
- Stax has WPEngine integration
- Stax has more features (snapshots, builds, linting)
- wp-env is simpler for basic development

**vs DDEV directly**:
- Stax adds WordPress-specific features on top of DDEV
- Stax handles database sync, search-replace automatically
- Stax is more opinionated (easier to use, less flexible)
- DDEV is more general-purpose

**vs Docker Compose directly**:
- Stax requires zero Docker knowledge
- Stax configures everything automatically
- Docker Compose gives you more control but much more complexity

### Does Stax work on Windows or Linux?

Currently **macOS only** (macOS 12 Monterey or later).

Windows and Linux support may come in future versions. The main blocker is credential storage (macOS Keychain) - we'd need equivalent secure storage for other platforms.

---

## Installation and Setup

### What are the system requirements?

**Minimum**:
- macOS 12.0 (Monterey) or later
- Intel or Apple Silicon processor
- 8GB RAM
- 10GB free disk space
- Docker Desktop 4.25+
- DDEV 1.22+

**Recommended**:
- macOS 13.0 (Ventura) or later
- Apple Silicon (M1/M2/M3) for best performance
- 16GB RAM
- 20GB+ free disk space

### How much disk space does Stax use?

**Stax itself**: <100MB

**Per project**:
- Docker images: ~3-5GB (shared across projects)
- Project files: ~100MB-1GB (your code)
- Database: ~100MB-500MB (depends on size)
- Snapshots: ~100MB each (can accumulate)

**Total for 2-3 projects**: 10-15GB

**Tip**: Clean up old snapshots regularly:
```bash
stax db list
stax db delete-snapshot old-snapshot
```

### Do I need to know Docker?

**No!** That's the point of Stax.

Stax handles all the Docker complexity. You just run `stax init` and everything works.

**You will use**:
- `stax start` (not `docker start`)
- `stax stop` (not `docker stop`)
- `stax ssh` (not `docker exec`)

**Behind the scenes** Stax uses Docker, but you don't need to know how.

### Can I use Stax without WPEngine?

**Yes**, but you'll lose some features:

**Works without WPEngine**:
- Local WordPress development
- DDEV container management
- Build automation
- Database snapshots
- Team configuration sharing

**Requires WPEngine**:
- `stax db pull` (pull from WPEngine)
- `stax provider sync` (sync files from WPEngine)
- Remote media proxying

**Alternative**: You can still manually import databases or use Stax as a pure local development tool.

---

## Using Stax

### How do I find my WPEngine install name?

Use the `stax list` command to see all available installs:

```bash
stax list
```

**Output**:
```
INSTALL NAME           ENVIRONMENT   PRIMARY DOMAIN              PHP   STATUS
myinstall              production    mysite.wpengine.com         8.1   active
myinstall-staging      staging       myinstall-staging.wpe       8.1   active
client-site            production    clientsite.com              8.2   active
```

The "Install Name" column shows what you need for `stax init`.

**You can also**:
- Filter by name: `stax list --filter="client.*"`
- Filter by environment: `stax list --environment=production`
- Get JSON output: `stax list --output=json`

**Alternative methods**:
1. Check WPEngine portal under "Sites"
2. Look in existing `.stax.yml` files from other projects
3. Ask your team lead or WPEngine account admin

### Can I use stax list without a project?

**Yes!** That's the whole point of `stax list`.

The `stax list` command is a global command that works anywhere, without needing:
- A `.stax.yml` file
- A project directory
- SSH keys
- Any existing Stax project

**All you need**:
- Stax installed
- WPEngine API credentials (run `stax setup` first)

**Example workflow**:
```bash
# Day 1: New computer, no projects yet
brew install stax
stax setup  # Configure API credentials
stax list   # See all available installs

# Now you know which install to use
cd ~/Sites/new-project
stax init   # Enter install name from list
```

### What if stax list shows no installs?

**Possible causes**:

**1. No installs in your WPEngine account**:
- Verify you have access to WPEngine installs
- Check with your WPEngine account admin
- Ensure you're using the correct API credentials

**2. API credentials incorrect**:
```bash
stax setup  # Reconfigure credentials
stax list   # Try again
```

**3. API access not enabled**:
- Log in to WPEngine portal
- Go to Account Settings > API Access
- Enable API access
- Create API credentials

**4. Using wrong WPEngine account**:
- Verify which WPEngine account your API credentials are for
- You might have multiple WPEngine accounts
- Use credentials for the account that has your installs

**5. Filters hiding results**:
```bash
# Remove all filters
stax list

# Instead of:
stax list --filter="something" --environment=production
```

### How do I update Stax?

**Via Homebrew**:
```bash
brew update
brew upgrade stax
```

**From source**:
```bash
cd ~/path/to/stax
git pull
make build
make install
```

**Check version**:
```bash
stax --version
```

### Can I run multiple projects at once?

**Yes!** Each project runs in its own containers:

```bash
# Project 1
cd ~/Sites/project1
stax start

# Project 2
cd ~/Sites/project2
stax start

# Both running simultaneously
```

Access them at their respective URLs:
- https://project1.local
- https://project2.local

**All projects share**:
- Docker Desktop
- ddev-router container (routes traffic)
- Docker images (saves space)

### How do I switch between projects?

Just navigate to the project directory:

```bash
# Work on project 1
cd ~/Sites/project1
stax start
stax status

# Switch to project 2
cd ~/Sites/project2
stax start
stax status
```

**Stopping projects**:
```bash
# Stop specific project
cd ~/Sites/project1
stax stop

# Stop all projects
ddev poweroff
```

### How often should I pull from production?

**Recommended**:
- **Weekly** for most development
- **Daily** if you need fresh content
- **As needed** for debugging production issues

**Don't pull too often** because:
- Takes 2-5 minutes
- Loses local database changes
- Creates large snapshots
- Uses bandwidth

**Better approach**: Pull from staging for most work, production only when needed.

### Do I need to download all my media files?

**No!** Stax has remote media proxying.

**How it works**:
- Request image: `/wp-content/uploads/2024/01/image.jpg`
- File doesn't exist locally
- Stax proxies from WPEngine/CDN
- Image displays in browser
- Optionally cached locally

**Advantages**:
- No large downloads
- Saves disk space
- Always current media

**When to download**:
- Working offline
- Testing upload functionality
- Need very fast media loads

**Enable proxying**:
```yaml
# .stax.yml
media:
  proxy:
    enabled: true
    remote_url: https://cdn.mysite.com
```

---

## WordPress and Multisite

### Can I use Stax for single-site WordPress?

**Yes!** Stax works great for single sites too.

Just choose during `stax init`:
```bash
? Project type:
  - wordpress (single site)
  ❯ wordpress-multisite
```

All Stax features work with single sites.

### How does Stax handle WordPress multisite?

Stax makes multisite easy:

**Subdomain mode** (recommended):
- Automatic wildcard DNS
- Automatic SSL certificates
- Per-site search-replace
- No /etc/hosts editing needed

**Subdirectory mode**:
- Also supported
- Simpler URLs (site.com/site1)
- All features work

**What Stax automates**:
- Multisite configuration
- Network and site setup
- SSL for all subdomains
- Search-replace for all sites
- Proper nginx configuration

See [MULTISITE.md](./MULTISITE.md) for details.

### Can I add custom domains to local multisite?

**Yes**, with some configuration.

**Example**: Production has `customdomain.com` pointing to a subsite.

**Local setup**:
```yaml
# .stax.yml
network:
  sites:
    - name: mysite
      domain: mysite.my-project.local
      wpengine_domain: mysite.com
      custom_domains:
        - customdomain.local  # Local version
```

Add to DDEV config:
```yaml
# .ddev/config.yaml
additional_fqdns:
  - customdomain.local
```

Restart:
```bash
stax restart
```

Now `customdomain.local` works locally.

---

## Database and Data

### Where are database snapshots stored?

**Location**: `~/.stax/snapshots/<project-name>/`

**Example**:
```bash
ls ~/.stax/snapshots/my-project/
# before-pull_2024-11-08.sql.gz
# before-migration_2024-11-07.sql.gz
```

**Compressed**: Snapshots are gzipped to save space.

**Cleanup**:
```bash
stax db list
stax db delete-snapshot old-snapshot
```

### What does --sanitize do?

`--sanitize` anonymizes sensitive data:

**Sanitized**:
- User emails → `user1@example.com`, `user2@example.com`
- Passwords → reset to `password`
- User personal data → randomized
- Comment author emails → anonymized

**Not sanitized**:
- Post content
- Page content
- Site structure
- Settings

**Use for**:
- Testing with production data
- Sharing with contractors
- Compliance (GDPR, etc.)
- Training

**Example**:
```bash
stax db pull --sanitize
```

### Can I edit the database directly?

**Yes**, several ways:

**WP-CLI** (recommended):
```bash
stax wp db query "SELECT * FROM wp_options WHERE option_name = 'siteurl'"
```

**MySQL CLI**:
```bash
stax ssh
mysql
USE db;
SHOW TABLES;
SELECT * FROM wp_options;
exit;
exit;
```

**phpMyAdmin**:
```bash
ddev launch -p
```

Opens phpMyAdmin in browser.

### How do I backup my local database?

**Snapshots** (recommended):
```bash
stax db snapshot my-backup
```

Stored in: `~/.stax/snapshots/`

**Export SQL**:
```bash
stax db export ~/backups/my-project.sql
```

**Both**:
```bash
stax db snapshot my-backup
stax db export ~/backups/my-project-$(date +%Y%m%d).sql
```

---

## Configuration and Customization

### Can team members have different settings?

**Yes!** Two config levels:

**Project config** (`.stax.yml`):
- In your project directory
- Committed to Git
- Shared by whole team
- Settings everyone needs

**Global config** (`~/.stax/config.yml`):
- In `~/.stax/config.yml`
- Not in Git
- Personal to each developer
- Personal preferences

**Example**: Different developers prefer different WPEngine environments:

Person A's `~/.stax/config.yml`:
```yaml
defaults:
  wpengine:
    environment: staging
```

Person B's `~/.stax/config.yml`:
```yaml
defaults:
  wpengine:
    environment: production
```

Same `.stax.yml`, different personal preferences.

### Can I customize the build process?

**Yes!** Configure in `.stax.yml`:

```yaml
build:
  pre_install:
    - echo "Starting build"

  install:
    - composer install --optimize-autoloader
    - npm ci

  post_install:
    - npm run build
    - scripts/custom-build.sh

  watch:
    enabled: true
    command: npm run watch
    paths:
      - wp-content/themes/*/assets/**
```

**Custom script** (`scripts/custom-build.sh`):
```bash
#!/bin/bash
echo "Running custom build"
npm run sass
npm run webpack
npm run imagemin
```

### Can I change PHP or MySQL versions?

**Yes**:

```bash
# Change PHP version
stax config set ddev.php_version 8.2

# Change MySQL version
stax config set ddev.mysql_version 8.0

# Restart to apply
stax restart
```

**Available versions**:
- **PHP**: 7.4, 8.0, 8.1, 8.2, 8.3
- **MySQL**: 5.7, 8.0
- **MariaDB**: 10.3, 10.4, 10.5, 10.6, 10.11

**Match WPEngine**:
Check your WPEngine environment's PHP/MySQL versions and match them locally.

---

## Performance and Troubleshooting

### Stax is slow - how can I speed it up?

**Enable Mutagen** (faster file sync on Mac):
```bash
echo "mutagen_enabled: true" >> .ddev/config.yaml
stax restart
```

**Increase Docker resources**:
- Docker Desktop → Settings → Resources
- Memory: 8GB
- CPUs: 4

**Skip unnecessary database data**:
```bash
stax db pull --skip-logs --skip-transients
```

**Disable Xdebug** (if enabled):
```bash
stax ssh
ddev xdebug off
exit
```

**Use SSD** (not HDD):
- Projects on SSD are much faster
- Check: System Preferences → About This Mac → Storage

### Why are my containers using so much RAM?

**Normal usage**:
- Each project: ~1-2GB
- Docker Desktop: ~2-4GB
- Total: ~4-8GB for 2 projects

**Reduce RAM usage**:

1. **Stop unused projects**:
   ```bash
   ddev poweroff
   ```

2. **Reduce Docker memory**:
   - Docker Desktop → Settings → Resources
   - Memory: 4GB (instead of 8GB)

3. **Stop Docker when not developing**:
   - Quit Docker Desktop
   - Frees all RAM

### What do I do if something breaks?

**Step 1: Run diagnostics**:
```bash
stax doctor
```

Checks common issues and suggests fixes.

**Step 2: Check logs**:
```bash
stax logs -f
```

Look for error messages.

**Step 3: Try restarting**:
```bash
stax restart
```

Fixes most container issues.

**Step 4: Nuclear option** (delete and recreate):
```bash
stax stop
ddev delete -Oy
stax init
```

**Step 5: Get help**:
- Check [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)
- Search GitHub issues
- Contact team

---

## Advanced Topics

### Can I use Stax in CI/CD?

**Possible but not recommended**. Stax is designed for local development.

**For CI/CD**:
- Use Docker Compose directly
- Or GitHub Actions with official WordPress actions
- Or platform-specific CI (WPEngine's CI, etc.)

Stax adds overhead (DDEV, keychain) not needed in CI.

### Can I extend Stax with custom commands?

**Not directly** (no plugin system yet).

**Workarounds**:

**Shell aliases**:
```bash
# ~/.zshrc
alias stax-deploy='stax build && git push'
```

**Scripts in your project**:
```bash
# scripts/custom-command.sh
#!/bin/bash
stax build
stax lint
stax ssh "composer test"
```

**DDEV custom commands**:
Create `.ddev/commands/web/custom-command`.

### Can I contribute to Stax?

**Yes!** Stax is an internal tool but welcomes contributions.

**Ways to contribute**:
- Report bugs on GitHub
- Suggest features
- Submit pull requests
- Improve documentation
- Share workflows and examples

**Development**:
```bash
git clone https://github.com/firecrown-media/stax.git
cd stax
make build
make test
```

See `ARCHITECTURE.md` for development details.

---

## Getting Help

### Where can I find more help?

**Documentation**:
- [README.md](../README.md) - Overview and quick start
- [INSTALLATION.md](./INSTALLATION.md) - Detailed installation
- [USER_GUIDE.md](./USER_GUIDE.md) - Complete user guide
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - Common problems
- [EXAMPLES.md](./EXAMPLES.md) - Real-world scenarios

**Run diagnostics**:
```bash
stax doctor
```

**Check logs**:
```bash
stax logs -f
```

**Get help**:
- Search GitHub issues
- Create new issue
- Internal: Slack #dev-tools
- Email: dev@firecrown.com

### How do I report a bug?

**Create a GitHub issue** with:

1. **Stax version**:
   ```bash
   stax --version
   ```

2. **System info**:
   - macOS version
   - DDEV version (`ddev version`)
   - Docker Desktop version

3. **What you were trying to do**:
   - Command you ran
   - Expected result
   - Actual result

4. **Error messages**:
   - Full error output
   - Output of `stax doctor`
   - Relevant logs

5. **Steps to reproduce**:
   Detailed steps so we can reproduce the issue.

### How do I request a feature?

**Create a GitHub issue** with:

1. **Feature description**:
   - What you want to do
   - Why it would be useful
   - Who would benefit

2. **Example usage**:
   - How you'd use it
   - Example commands
   - Expected output

3. **Alternatives**:
   - Current workarounds
   - How other tools do it

---

## Questions Not Answered Here?

**Ask!** We'll add it to this FAQ.

- GitHub Issues
- Internal: Slack #dev-tools
- Email: dev@firecrown.com

---

**Happy developing with Stax!**

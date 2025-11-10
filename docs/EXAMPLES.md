# Stax Real-World Examples

Practical examples for common workflows and scenarios.

> **Note**: The examples below show multisite configurations for demonstration purposes, but all workflows work equally well for single-site WordPress projects. Simply use `type: wordpress` instead of `type: wordpress-multisite` in your `.stax.yml` configuration.

---

## Table of Contents

- [Discovering Installs](#discovering-installs)
- [Daily Development](#daily-development)
- [Database Workflows](#database-workflows)
- [Feature Development](#feature-development)
- [Team Collaboration](#team-collaboration)
- [Testing and QA](#testing-and-qa)
- [Emergency Recovery](#emergency-recovery)
- [Onboarding](#onboarding)

---

## Discovering Installs

### Scenario 1: New Team Member Setup

**Goal**: Onboard a new developer with zero knowledge of the project.

```bash
# Day 1: Install prerequisites
brew tap firecrown-media/stax
brew install stax

# Configure WPEngine credentials
stax setup
# Enter API username: developer@company.com
# Enter API password: ********

# Discover available installs
stax list
```

**Output**:
```
INSTALL NAME           ENVIRONMENT   PRIMARY DOMAIN              PHP   STATUS
client-prod            production    client.com                  8.2   active
client-staging         staging       client-staging.wpe          8.2   active
internal-site          production    internal.company.com        8.1   active
```

```bash
# Clone project repository
git clone https://github.com/company/client-project.git
cd client-project

# Initialize Stax project (use install from list)
stax init
# Enter install name: client-staging

# Wait 3-5 minutes for setup...
# Open site
open https://client-project.local
```

**Expected result**: New developer productive in 15 minutes.

### Scenario 2: Finding the Right Install

**Goal**: You have multiple clients and need to find the right install.

```bash
# List all installs
stax list

# Too many results, filter by client name
stax list --filter="acme.*"
```

**Output**:
```
INSTALL NAME           ENVIRONMENT   PRIMARY DOMAIN              PHP   STATUS
acme-prod              production    acme.com                    8.2   active
acme-staging           staging       acme-staging.wpe            8.2   active
acme-dev               development   acme-dev.wpe                8.1   active
```

```bash
# Found it! Use staging for development
cd ~/Sites/acme-project
stax init
# Enter install name: acme-staging
```

### Scenario 3: Auditing All Installs

**Goal**: Generate a report of all WPEngine installs and their configurations.

```bash
# Get all installs as JSON
stax list --output=json > all-installs.json

# Get only production installs
stax list --environment=production --output=json > production-installs.json

# Process with jq to get specific info
stax list --output=json | jq '.[] | {name: .name, php: .metadata.php_version, domain: .primary_domain}'

# Find installs on old PHP versions
stax list --output=json | jq '.[] | select(.metadata.php_version < "8.1") | .name'

# Count installs by environment
stax list --output=json | jq 'group_by(.environment) | map({environment: .[0].environment, count: length})'
```

**Expected output**:
```json
[
  {"environment": "production", "count": 15},
  {"environment": "staging", "count": 15},
  {"environment": "development", "count": 3}
]
```

**Use cases**:
- Security audit (find old PHP versions)
- Inventory management
- Documentation generation
- Billing/usage tracking

### Scenario 4: Forgot Install Name Mid-Project

**Goal**: You're setting up an existing project but don't know the install name.

```bash
# Clone the repo
git clone https://github.com/company/mystery-project.git
cd mystery-project

# Check if there's a .stax.yml (maybe it's already configured)
cat .stax.yml
# No file found...

# List installs to find it
stax list

# Too many, filter by likely names from the repo
stax list --filter="mystery"
stax list --filter="company.*"

# Still not sure? Get help from README or ask team
cat README.md | grep -i "wpengine"

# Found it in README: "Uses WPEngine install: company-mystery-prod"
stax init
# Enter: company-mystery-prod
```

---

## Daily Development

### Scenario 1: Starting Your Work Day

**Goal**: Get your environment running and start coding.

```bash
# Monday morning, coffee in hand
cd ~/Sites/firecrown-multisite

# Start environment
stax start
# Takes ~10 seconds

# Check status
stax status

# Open in browser
open https://firecrown.local

# Open in VS Code
code .

# Start watching for file changes
stax dev
```

**Expected result**: Environment running, files watching, ready to code in ~30 seconds.

### Scenario 2: Working on a Theme Update

**Goal**: Update theme CSS and JavaScript.

```bash
# Start watch mode (in one terminal)
stax dev

# In another terminal, make changes
code wp-content/themes/my-theme/assets/scss/main.scss

# Save file - automatically rebuilds
# Check browser - changes appear

# When done, build production version
stax build --production

# Commit changes
git add .
git commit -m "Update theme header styles"
git push
```

**Tips**:
- Use two terminal windows: one for watching, one for commands
- Watch mode rebuilds faster than full builds
- Production builds are optimized (minified, compressed)

### Scenario 3: Testing a Plugin Change

**Goal**: Test a plugin modification.

```bash
# Edit plugin file
code wp-content/plugins/my-plugin/my-plugin.php

# Clear WordPress caches
stax wp cache flush

# Test in browser
open https://firecrown.local

# Check for errors in logs
stax logs -f

# Works! Commit changes
git add wp-content/plugins/my-plugin/
git commit -m "Fix cart calculation bug"
```

**Tips**:
- Always clear cache after plugin changes
- Watch logs for PHP errors
- Test on multiple subsites for multisite

### Scenario 4: Ending Your Day

**Goal**: Clean up and shut down.

```bash
# Stop watch mode if running (Ctrl+C)

# Commit any work in progress
git add .
git commit -m "WIP: header redesign"

# Stop environment
stax stop

# Or leave it running for faster access tomorrow
# Docker Desktop uses ~2-4GB RAM when idle
```

**Tips**:
- Stop if you want to free up RAM
- Leave running for instant access tomorrow
- Commit WIP branches so you don't lose work

---

## Database Workflows

### Scenario 5: Refreshing Your Database

**Goal**: Get the latest content from production.

```bash
# Create snapshot first (safety)
stax db snapshot before-pull

# Pull latest from production
stax db pull

# Check a few posts/pages to verify
open https://firecrown.local/wp-admin

# If something's wrong, restore
stax db restore before-pull
```

**Expected result**: Latest production content in ~3 minutes.

**Tips**:
- Always snapshot before pulls
- Use `--skip-logs` for faster imports
- Pull from staging for most work

### Scenario 6: Testing a Database Migration

**Goal**: Safely test a database migration script.

```bash
# 1. Start with fresh production data
stax db pull

# 2. Create snapshot
stax db snapshot before-migration

# 3. Run migration
stax ssh
wp db query < /var/www/html/migrations/001-add-user-meta.sql
exit

# 4. Verify migration
stax wp db query "SHOW COLUMNS FROM wp_usermeta"
stax wp db query "SELECT COUNT(*) FROM wp_usermeta WHERE meta_key = 'new_field'"

# 5. Test affected features
# ... test in browser ...

# If good:
git add migrations/001-add-user-meta.sql
git commit -m "Add user metadata migration"

# If bad:
stax db restore before-migration
# Fix migration script and try again
```

**Expected result**: Migration tested safely with easy rollback.

### Scenario 7: Creating a Clean Testing Database

**Goal**: Anonymize data for sharing or testing.

```bash
# Pull and sanitize
stax db pull --sanitize

# This anonymizes:
# - User emails → user1@example.com, user2@example.com
# - Passwords → reset to 'password'
# - Personal data → randomized

# Verify
stax wp user list
# Should show sanitized emails

# Create snapshot for reuse
stax db snapshot clean-test-data
```

**Expected result**: Safe, anonymized database for testing.

**When to use**:
- Sharing database with contractors
- Testing with sensitive production data
- Compliance testing
- Training new developers

---

## Feature Development

### Scenario 8: Developing a New Feature

**Goal**: Build a new feature with a proper workflow.

```bash
# 1. Create feature branch
git checkout -b feature/user-dashboard

# 2. Pull latest staging data
stax db pull --environment=staging

# 3. Create snapshot (safety)
stax db snapshot before-feature-work

# 4. Start development
stax dev  # Watch mode

# 5. Make changes
code wp-content/themes/my-theme/
code wp-content/mu-plugins/my-plugin/

# 6. Test frequently
# ... test in browser ...
stax wp cache flush

# 7. Commit incrementally
git add wp-content/themes/my-theme/templates/dashboard.php
git commit -m "Add dashboard template"

git add wp-content/mu-plugins/my-plugin/dashboard.php
git commit -m "Add dashboard data fetching"

# 8. Final test
stax build --production
stax lint
open https://firecrown.local/dashboard

# 9. Push for review
git push origin feature/user-dashboard

# 10. Create PR
# (use GitHub UI or gh CLI)
```

**Expected result**: Feature developed with clean commits and thorough testing.

### Scenario 9: Testing a Feature Branch

**Goal**: Review and test a teammate's feature branch.

```bash
# 1. Fetch latest branches
git fetch origin

# 2. Checkout feature branch
git checkout feature/team-member-feature

# 3. Pull dependencies (if needed)
stax ssh
composer install
npm install
exit

# 4. Build
stax build

# 5. Refresh database (if needed)
stax db pull --environment=staging

# 6. Test the feature
open https://firecrown.local

# 7. Run tests
stax lint
stax ssh
composer test
npm test
exit

# 8. Leave feedback
# Comment on PR with findings

# 9. Return to your branch
git checkout your-branch
```

**Expected result**: Feature thoroughly tested with feedback provided.

---

## Team Collaboration

### Scenario 10: Onboarding a New Team Member

**Goal**: Get a new developer up and running.

**New Developer**:
```bash
# 1. Install prerequisites
# Follow INSTALLATION.md

# 2. Configure credentials (one-time)
stax setup
# Enter WPEngine and GitHub credentials

# 3. Clone project
git clone https://github.com/mycompany/my-project.git
cd my-project

# 4. Initialize project
stax init
# Reads .stax.yml, sets up identical environment

# 5. Wait 3-5 minutes
# Stax pulls database, builds assets, etc.

# 6. Start coding!
open https://my-project.local
code .
```

**Expected time**: 10-20 minutes total (including prerequisite installation).

**Tips**:
- `.stax.yml` in Git ensures everyone has identical config
- New developer doesn't need to know DDEV/Docker details
- Stax handles all the complexity

### Scenario 11: Sharing a Database State

**Goal**: Share a specific database state with your team.

```bash
# Developer A: Create and share snapshot
stax db snapshot feature-xyz-test-data
stax db export ~/Dropbox/feature-xyz-test-data.sql

# Share file via Dropbox/Drive/etc.

# Developer B: Import snapshot
stax db import ~/Downloads/feature-xyz-test-data.sql
```

**Or use WPEngine staging**:
```bash
# All team members pull from staging
stax db pull --environment=staging
# Everyone gets the same data
```

### Scenario 12: Coordinating a Database Schema Change

**Goal**: Update database schema across team.

**Lead Developer**:
```bash
# 1. Create migration
code migrations/002-add-product-table.sql

# 2. Test locally
stax db snapshot before-schema-change
stax ssh
wp db query < migrations/002-add-product-table.sql
exit

# 3. Verify
stax wp db query "SHOW TABLES LIKE 'wp_products'"

# 4. Commit migration
git add migrations/002-add-product-table.sql
git commit -m "Add products table migration"
git push

# 5. Notify team in Slack
# "New migration: 002-add-product-table.sql - run after pulling"
```

**Team Members**:
```bash
# 1. Pull latest code
git pull

# 2. Create snapshot
stax db snapshot before-schema-change

# 3. Run migration
stax ssh
wp db query < migrations/002-add-product-table.sql
exit

# 4. Verify
stax wp db query "SHOW TABLES LIKE 'wp_products'"
```

---

## Testing and QA

### Scenario 13: Testing with Production Data

**Goal**: Debug a production issue locally.

```bash
# 1. Create snapshot of current state
stax db snapshot before-prod-testing

# 2. Pull production database
stax db pull --environment=production

# 3. Reproduce issue
# ... test in browser ...

# 4. Make fix
code wp-content/plugins/my-plugin/checkout.php

# 5. Test fix
stax wp cache flush
# ... verify fix in browser ...

# 6. Commit fix
git add wp-content/plugins/my-plugin/checkout.php
git commit -m "Fix checkout tax calculation"

# 7. Deploy fix
git push

# 8. Restore staging data
stax db restore before-prod-testing
# Or pull staging
stax db pull --environment=staging
```

**Expected result**: Bug fixed with production data, no production access needed.

### Scenario 14: QA Testing Before Release

**Goal**: Thoroughly test before deploying to production.

```bash
# 1. Checkout release branch
git checkout release/v2.1.0

# 2. Pull staging database (closest to production)
stax db pull --environment=staging

# 3. Build production assets
stax build --production

# 4. Run linters
stax lint

# 5. Run tests
stax ssh
composer test
npm test
exit

# 6. Manual testing checklist
# [ ] Login works
# [ ] Checkout flow works
# [ ] Admin panels work
# [ ] All subsites accessible
# etc.

# 7. Performance check
# Load pages, check Network tab in browser

# 8. If all good, approve for deployment
# Leave comment on PR or notify team
```

### Scenario 15: Regression Testing After Update

**Goal**: Verify WordPress/plugin update doesn't break things.

```bash
# 1. Snapshot before update
stax db snapshot before-wordpress-update

# 2. Update WordPress
stax wp core update

# 3. Update plugins
stax wp plugin update --all

# 4. Test critical paths
# [ ] Can log in
# [ ] Can create post
# [ ] Can publish post
# [ ] Frontend displays correctly
# [ ] Forms work
# [ ] Checkout works (if e-commerce)

# 5. If something broke:
stax db restore before-wordpress-update
# Debug the issue

# 6. If all good:
git add composer.lock  # If plugins updated via Composer
git commit -m "Update WordPress core and plugins"
```

---

## Emergency Recovery

### Scenario 16: Recovering from a Bad Database Import

**Goal**: Restore after accidentally importing wrong database.

```bash
# Oh no! I just imported the wrong database

# List snapshots
stax db list

# Output:
# before-import     5 minutes ago    245 MB
# daily-backup      1 day ago        243 MB

# Restore recent snapshot
stax db restore before-import

# Verify
open https://my-project.local

# All good!
```

**Expected result**: Back to working state in ~1 minute.

**Lesson**: Stax auto-creates snapshots before pulls/imports.

### Scenario 17: Recovering from Bad Migration

**Goal**: Roll back after failed database migration.

```bash
# Migration failed halfway through

# Check what snapshots exist
stax db list

# Restore pre-migration snapshot
stax db restore before-migration

# Fix migration script
code migrations/003-update-user-roles.sql

# Try again
stax db snapshot before-migration-v2
stax ssh
wp db query < migrations/003-update-user-roles.sql
exit

# Success!
```

### Scenario 18: Starting Fresh

**Goal**: Nuclear option - completely recreate environment.

```bash
# Everything is broken, let's start fresh

# 1. Delete environment
stax stop
ddev delete -Oy

# 2. Reinitialize
stax init
# Uses .stax.yml config
# Pulls fresh database
# Rebuilds everything

# 3. Wait 3-5 minutes

# 4. Back to working state
open https://my-project.local
```

**Expected result**: Clean environment in ~5 minutes.

**When to use**:
- Containers are corrupted
- Configuration is messed up
- Simpler than debugging

---

## Tips for All Scenarios

### General Workflow

1. **Always snapshot before risky operations**
   ```bash
   stax db snapshot before-<operation>
   ```

2. **Commit frequently**
   - Small, atomic commits
   - Clear commit messages
   - Push to remote often

3. **Test before committing**
   ```bash
   stax lint
   stax build
   # Manual testing
   ```

4. **Use branches for features**
   ```bash
   git checkout -b feature/my-feature
   ```

5. **Keep database fresh**
   - Pull from staging weekly
   - Pull from production as needed
   - Don't pull too often

### Performance Tips

1. **Use staging by default**
   ```bash
   stax config set wpengine.environment staging
   ```

2. **Skip unnecessary data**
   ```bash
   stax db pull --skip-logs --skip-transients
   ```

3. **Enable Mutagen** (faster file sync)
   ```bash
   echo "mutagen_enabled: true" >> .ddev/config.yaml
   stax restart
   ```

4. **Stop when not using**
   ```bash
   stax stop  # Frees up RAM
   ```

### Safety Tips

1. **Snapshots are your friend**
   - Before pulls
   - Before migrations
   - Before risky operations

2. **Test with staging first**
   - Use production data sparingly
   - Staging is usually good enough

3. **Sanitize sensitive data**
   ```bash
   stax db pull --sanitize
   ```

4. **Version control your config**
   - `.stax.yml` in Git
   - Team shares configuration
   - Easy recovery if lost

---

## Next Steps

- **User Guide**: [USER_GUIDE.md](./USER_GUIDE.md) - Complete usage guide
- **Troubleshooting**: [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - Fix problems
- **FAQ**: [FAQ.md](./FAQ.md) - Common questions
- **Command Reference**: [COMMAND_REFERENCE.md](./COMMAND_REFERENCE.md) - All commands

---

**Have a workflow to share?** Document it and contribute to this guide!

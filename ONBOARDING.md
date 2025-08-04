# Stax Developer Onboarding Guide

**Complete workflow guide for new Firecrown team members using Stax with WP Engine integration.**

This guide walks through the complete development lifecycle from initial setup through production deployment, demonstrating Stax's integration with Git workflows and WP Engine hosting.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Step 1: Install Stax](#step-1-install-stax)
- [Step 2: Configure WP Engine Access](#step-2-configure-wp-engine-access)
- [Step 3: Clone Client Repository](#step-3-clone-client-repository)
- [Step 4: Initialize Local Development Environment](#step-4-initialize-local-development-environment)
- [Step 5: Sync Database and Files from WP Engine](#step-5-sync-database-and-files-from-wp-engine)
- [Step 6: Make Development Changes](#step-6-make-development-changes)
- [Step 7: Commit Changes to Git](#step-7-commit-changes-to-git)
- [Step 8: Deploy to WP Engine Staging via GitHub Actions](#step-8-deploy-to-wp-engine-staging-via-github-actions)
- [Step 9: Test in WP Engine Staging](#step-9-test-in-wp-engine-staging)
- [Step 10: Create Production Pull Request](#step-10-create-production-pull-request)
- [Step 11: Deploy to WP Engine Production](#step-11-deploy-to-wp-engine-production)
- [Step 12: Verify Production Deployment](#step-12-verify-production-deployment)
- [Complete Workflow Summary](#complete-workflow-summary)
- [Troubleshooting Common Issues](#troubleshooting-common-issues)
- [Next Steps](#next-steps)
- [Support](#support)

## Prerequisites

Before starting, ensure you have:

### System Requirements
- **macOS** (Intel or Apple Silicon) or **Linux**
- **Docker Desktop** installed and running
- **Git** configured with your credentials
- **Access to Firecrown's GitHub organization**
- **WP Engine account access** (contact your team lead)

### Required Accounts & Permissions
- **GitHub**: Access to Firecrown-Media organization repositories
- **WP Engine**: Account with API and SSH access to client installations
- **SSH Key**: Added to your WP Engine account (see setup below)

---

## Step 1: Install Stax

### Option A: Install via Homebrew (Recommended)

```bash
# Add Firecrown's Stax tap
brew tap firecrown-media/stax

# Install Stax
brew install stax

# Verify installation
stax --version
```

### Option B: Install from Source (Development)

```bash
# Clone and build from source
git clone https://github.com/Firecrown-Media/stax.git
cd stax
make install
stax --version
```

### Install Dependencies

```bash
# Install additional dependencies (DDEV, WP-CLI)
make install-deps

# Verify DDEV is available
ddev version

# Verify WP-CLI is available  
wp --info
```

---

## Step 2: Configure WP Engine Access

### Set Up SSH Key for WP Engine

```bash
# Generate SSH key if you don't have one
ssh-keygen -t rsa -b 4096 -C "your.email@firecrown.com"

# Copy public key to clipboard (macOS)
cat ~/.ssh/id_rsa.pub | pbcopy

# Copy public key to clipboard (Linux)
cat ~/.ssh/id_rsa.pub | xclip -selection clipboard
```

**Add SSH key to WP Engine:**
1. Log into [WP Engine User Portal](https://my.wpengine.com/)
2. Navigate to **SSH Gateway** → **SSH Keys**
3. Click **Add Public Key**
4. Paste your public key and save

**Reference:** [WP Engine SSH Key Setup](https://wpengine.com/support/ssh-gateway/#Add_SSH_Key)

### Configure WP Engine API Access

**Get API credentials:**
1. Contact your Firecrown team lead for WP Engine API access
2. Follow [WP Engine API Setup Guide](https://wpengine.com/support/enabling-wp-engine-api/)
3. Save credentials securely (use a password manager)

**Set environment variables:**

```bash
# Add to your shell profile (.zshrc, .bashrc, etc.)
export WPE_USERNAME=your-wpe-username
export WPE_PASSWORD=your-secure-api-password

# Reload your shell or run:
source ~/.zshrc  # or ~/.bashrc
```

### Create Global Stax Configuration

Create `~/.stax.yaml` with Firecrown team standards:

```yaml
# Standard Firecrown development environment
php_version: "8.2"
webserver: "nginx-fpm"
database: "mysql:8.0"

# WordPress defaults
wordpress_defaults:
  admin_user: "fcadmin"
  admin_email: "dev@firecrown.com"

# WP Engine integration
hosting:
  wpengine:
    username: "your-wpe-username"
    sync_defaults:
      skip_media: true  # Use production CDN for faster sync
      exclude_dirs:
        - "wp-content/cache/"
        - "wp-content/uploads/backup-*"
    ssh_key_path: "~/.ssh/id_rsa"

# Standard Firecrown plugin stack
default_plugins:
  - "advanced-custom-fields-pro"
  - "yoast-seo"
```

---

## Step 3: Clone Client Repository

**Scenario**: Working on an existing client project hosted on WP Engine.

```bash
# Clone the client repository
git clone git@github.com:Firecrown-Media/client-website.git
cd client-website

# Check if project has Stax configuration
ls -la stax.yaml  # Project-specific Stax config

# If no stax.yaml exists, you'll create one in the next step
```

**Example project structure:**
```
client-website/
├── .git/
├── .github/
│   └── workflows/
│       ├── deploy-staging.yml
│       └── deploy-production.yml
├── wp-content/
│   ├── themes/
│   └── plugins/
├── stax.yaml           # Stax project configuration
└── README.md
```

---

## Step 4: Initialize Local Development Environment

### Create Local Environment

```bash
# Initialize Stax environment for this project
stax init client-website-local --php-version=8.2

# This creates a DDEV environment with the name 'client-website-local'
```

### Configure Project-Specific Settings

If `stax.yaml` doesn't exist, create it:

```yaml
name: "client-website-local"
type: "wordpress"

# Environment settings
php_version: "8.2"
database: "mysql:8.0"

# WordPress configuration
wordpress:
  url: "https://client-website-local.ddev.site"
  title: "Client Website - Development"

# WP Engine integration
hosting:
  wpengine:
    install_name: "clientwebsite"  # WP Engine installation name
    environment: "production"      # Default sync source

# Client-specific plugins
plugins:
  - "woocommerce"
  - "gravityforms"
  - "client-custom-plugin"
```

---

## Step 5: Sync Database and Files from WP Engine

### Full Sync (Database + Files)

```bash
# Sync everything from WP Engine production
stax wpe sync clientwebsite

# This will:
# 1. Download the database from WP Engine
# 2. Import it into your local environment
# 3. Download files (wp-content/uploads, etc.)
# 4. Update URLs for local development
```

### Database-Only Sync (Faster)

```bash
# For faster sync, get database only and use production CDN for media
stax wpe sync clientwebsite --skip-files

# This is recommended for regular development work
```

### Verify Local Site

```bash
# Check site status
stax status

# Expected output:
# client-website-local    ✅ Running  https://client-website-local.ddev.site

# Open in browser
open https://client-website-local.ddev.site

# Or get the URL
echo "https://client-website-local.ddev.site"
```

**Login to WordPress admin:**
- URL: `https://client-website-local.ddev.site/wp-admin`
- Username: `fcadmin` (from your global config)
- Password: Check with team lead or reset via WP-CLI

---

## Step 6: Make Development Changes

### Example: Update Theme Template

```bash
# Make changes to the active theme
# Edit wp-content/themes/client-theme/index.php
echo "<!-- Development change: $(date) -->" >> wp-content/themes/client-theme/index.php

# Or add a new feature
mkdir -p wp-content/themes/client-theme/template-parts
cat > wp-content/themes/client-theme/template-parts/hero-section.php << 'EOF'
<?php
/**
 * Hero Section Template Part
 * Added by Firecrown development team
 */
?>
<section class="hero-section">
    <div class="container">
        <h1>Welcome to our updated site!</h1>
        <p>This is a new feature developed by Firecrown.</p>
    </div>
</section>
EOF
```

### Test Changes Locally

```bash
# Clear any caches
stax wp cache flush

# Verify changes in browser
open https://client-website-local.ddev.site

# Run any local tests
stax wp plugin list --status=active
```

---

## Step 7: Commit Changes to Git

### Create Feature Branch

```bash
# Create feature branch from main
git checkout main
git pull origin main
git checkout -b feature/hero-section-update

# Add changes
git add wp-content/themes/client-theme/

# Commit with descriptive message
git commit -m "Add hero section template part

- Create new hero-section.php template part
- Add responsive styling for hero section
- Integrate with existing theme structure
- Tested locally in Stax development environment"

# Push feature branch
git push origin feature/hero-section-update
```

### Create Pull Request

```bash
# Create PR using GitHub CLI
gh pr create \
  --title "Add Hero Section Template Part" \
  --body "## Changes
- Added new hero section template part
- Responsive design implementation
- Tested in local Stax environment

## Testing
- ✅ Local development site working
- ✅ Template renders correctly
- ✅ No PHP errors or warnings

## Deployment Plan
1. Deploy to staging for client review
2. After approval, deploy to production

Ready for staging deployment and client review." \
  --base main

# This will output a PR URL for review
```

---

## Step 8: Deploy to WP Engine Staging via GitHub Actions

### GitHub Actions Workflow

Your repository should have `.github/workflows/deploy-staging.yml`:

```yaml
name: Deploy to WP Engine Staging

on:
  pull_request:
    branches: [ main ]
    types: [ opened, synchronize ]

jobs:
  deploy-staging:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
      
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        
    - name: Install dependencies
      run: npm install
      
    - name: Build assets
      run: npm run build
      
    - name: Deploy to WP Engine Staging
      uses: wpengine/github-action-wpe-site-deploy@v3
      with:
        WPE_SSHG_KEY_PRIVATE: ${{ secrets.WPE_SSHG_KEY_PRIVATE }}
        WPE_ENV: staging
        SRC_PATH: "wp-content"
        REMOTE_PATH: "wp-content" 
        
    - name: Comment on PR
      uses: actions/github-script@v6
      with:
        script: |
          github.rest.issues.createComment({
            issue_number: context.issue.number,
            owner: context.repo.owner,
            repo: context.repo.repo,
            body: '🚀 Deployed to staging: https://clientwebsite-staging.wpengine.com'
          })
```

### Automatic Staging Deployment

When you create the PR, GitHub Actions automatically:

1. **Builds assets** (CSS, JS compilation)
2. **Deploys to WP Engine staging** environment
3. **Comments on PR** with staging URL

**Monitor deployment:**
```bash
# Check GitHub Actions status
gh run list --branch feature/hero-section-update

# View logs if needed
gh run view [run-id] --log
```

---

## Step 9: Test in WP Engine Staging

### Access Staging Environment

```bash
# Staging URL (from GitHub Actions comment)
open https://clientwebsite-staging.wpengine.com

# Or check WP Engine portal for staging URL
```

### Staging Testing Checklist

**Frontend Testing:**
- ✅ Hero section displays correctly
- ✅ Responsive design works on mobile/tablet
- ✅ No layout issues or broken elements
- ✅ Page load speed acceptable

**Backend Testing:**
- ✅ WordPress admin accessible
- ✅ No PHP errors in logs
- ✅ All plugins functioning
- ✅ Theme customizations preserved

**Client Review:**
- Share staging URL with client
- Gather feedback and approval
- Document any requested changes

### Make Additional Changes (if needed)

```bash
# Switch back to local development
cd client-website

# Make updates based on feedback
# Edit files as needed...

# Commit and push updates
git add .
git commit -m "Update hero section based on client feedback

- Adjust hero section copy per client request
- Modify button styling for better contrast
- Ensure accessibility compliance"

git push origin feature/hero-section-update

# This triggers another staging deployment automatically
```

---

## Step 10: Create Production Pull Request

### PR Review Process

**Team Review:**
```bash
# Request review from team lead
gh pr edit feature/hero-section-update --add-reviewer @team-lead-username

# Add relevant labels
gh pr edit feature/hero-section-update --add-label "ready-for-production"
```

**Final Testing:**
- ✅ Code review completed
- ✅ Client approval received
- ✅ Staging tests passed
- ✅ Performance impact assessed
- ✅ Security review (if applicable)

### Approve and Merge PR

**Team Lead approval:**
```bash
# Team lead reviews and approves
gh pr review feature/hero-section-update --approve --body "Changes look good. Client approved staging. Ready for production deployment."

# Merge to main (this can be done via GitHub web interface)
gh pr merge feature/hero-section-update --merge --delete-branch
```

---

## Step 11: Deploy to WP Engine Production

### Production Deployment Workflow

Your repository should have `.github/workflows/deploy-production.yml`:

```yaml
name: Deploy to WP Engine Production

on:
  push:
    branches: [ main ]

jobs:
  deploy-production:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
      
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        
    - name: Install dependencies
      run: npm install
      
    - name: Build production assets
      run: npm run build:production
      
    - name: Deploy to WP Engine Production
      uses: wpengine/github-action-wpe-site-deploy@v3
      with:
        WPE_SSHG_KEY_PRIVATE: ${{ secrets.WPE_SSHG_KEY_PRIVATE }}
        WPE_ENV: production
        SRC_PATH: "wp-content"
        REMOTE_PATH: "wp-content"
        
    - name: Notify team
      uses: 8398a7/action-slack@v3
      with:
        status: ${{ job.status }}
        text: "🚀 Production deployment complete: https://clientwebsite.com"
      env:
        SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

### Automatic Production Deployment

When the PR is merged to `main`:

1. **GitHub Actions triggers** production deployment workflow
2. **Assets are built** for production (minified, optimized)
3. **Code deploys** to WP Engine production environment
4. **Team notification** sent via Slack/email

**Monitor production deployment:**
```bash
# Check deployment status
gh run list --branch main --limit 1

# View detailed logs
gh run view --log
```

---

## Step 12: Verify Production Deployment

### Production Testing

**Immediate checks:**
```bash
# Check site accessibility
curl -I https://clientwebsite.com

# Verify specific feature
curl -s https://clientwebsite.com | grep -i "hero-section" && echo "✅ Hero section deployed"
```

**Comprehensive testing:**
- ✅ Site loads correctly
- ✅ New hero section visible
- ✅ No broken functionality
- ✅ WordPress admin accessible
- ✅ Performance metrics acceptable

### Post-Deployment Checklist

**Monitor for issues:**
- Check error logs in WP Engine portal
- Monitor site performance
- Verify analytics tracking still works
- Test contact forms and key functionality

**Client notification:**
```bash
# Send deployment notification
echo "✅ Hero section feature deployed to production
🌐 Live site: https://clientwebsite.com
📊 Performance: All metrics normal
🔍 Testing: All functionality verified

Changes are now live for your users!" | \
mail -s "Feature Deployment Complete" client@example.com
```

---

## Complete Workflow Summary

**Local Development:**
1. ✅ Install Stax and configure WP Engine access
2. ✅ Clone client repository
3. ✅ Initialize local environment with `stax init`
4. ✅ Sync data with `stax wpe sync`
5. ✅ Make development changes
6. ✅ Test locally

**Git Workflow:**
7. ✅ Create feature branch
8. ✅ Commit changes
9. ✅ Push branch and create PR

**Staging Deployment:**
10. ✅ GitHub Actions deploys to staging automatically
11. ✅ Test in staging environment
12. ✅ Get client approval

**Production Deployment:**
13. ✅ Team review and PR approval
14. ✅ Merge to main branch
15. ✅ GitHub Actions deploys to production automatically
16. ✅ Verify production deployment

---

## Troubleshooting Common Issues

### Stax Connection Issues

```bash
# Test WP Engine connectivity
stax wpe list

# Verify SSH access
ssh your-username@clientwebsite.ssh.wpengine.net

# Check environment status
stax status
```

### Local Development Issues

```bash
# Restart DDEV environment
ddev restart

# Clear WordPress caches
stax wp cache flush

# Reset file permissions
stax wp eval "echo 'Permissions reset'"
```

### Deployment Issues

```bash
# Check GitHub Actions logs
gh run view --log

# Verify WP Engine deployment status in portal
# Check error logs in WP Engine dashboard
```

---

## Next Steps

**For ongoing development:**
- Use `stax wpe sync --skip-files` for regular database updates
- Create feature branches for each new task
- Follow the same staging → production workflow
- Keep local environment updated with `brew upgrade stax`

**Advanced workflows:**
- Learn hot swap for testing different PHP versions: `stax swap preset modern`
- Use batch operations for multi-client projects
- Explore Stax automation features for repetitive tasks

**Team collaboration:**
- Share Stax configurations via project `stax.yaml` files
- Document client-specific requirements
- Contribute improvements back to Stax development

---

## Support

**For help:**
- **Command help**: `stax --help` or `stax [command] --help`
- **Team support**: Contact your Firecrown team lead
- **Technical issues**: [GitHub Issues](https://github.com/Firecrown-Media/stax/issues)
- **WP Engine support**: WP Engine customer portal

**Resources:**
- [Stax Documentation](https://github.com/Firecrown-Media/stax/blob/main/README.md)
- [WP Engine Developer Portal](https://wpengine.com/developers/)
- [DDEV Documentation](https://ddev.readthedocs.io/)

---

_This onboarding guide demonstrates the complete Firecrown development workflow using Stax. The same pattern applies to all client projects with WP Engine hosting._
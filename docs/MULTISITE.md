# WordPress Multisite with Stax

Guide for setting up and managing WordPress Multisite installations with Stax.

## Table of Contents
- [Understanding Multisite](#understanding-multisite)
- [Setting Up Multisite](#setting-up-multisite)
- [Managing Subsites](#managing-subsites)
- [Syncing Multisite from WP Engine](#syncing-multisite-from-wp-engine)
- [Domain Configuration](#domain-configuration)
- [Common Multisite Tasks](#common-multisite-tasks)
- [Troubleshooting](#troubleshooting)

## Understanding Multisite

### What is WordPress Multisite?

WordPress Multisite allows you to run multiple WordPress sites from a single installation. All sites share:
- WordPress core files
- Themes and plugins
- User accounts (with different roles per site)
- Single database (with separate tables per site)

### When to Use Multisite

**Good use cases:**
- Multiple brand sites for one organization
- Network of related blogs
- Franchise locations
- Multi-language sites
- University department sites

**Not ideal for:**
- Completely unrelated sites
- Sites needing different hosting
- Sites with vastly different performance needs
- Client sites for an agency (unless all for same client)

### Multisite Types

| Type | URL Structure | Example | Use Case |
|------|---------------|---------|----------|
| **Subdomain** | site1.network.com | blog.company.com | Different brands/sections |
| **Subdirectory** | network.com/site1 | company.com/blog | Related content sections |
| **Domain Mapping** | custom-domain.com | brand1.com, brand2.com | Completely separate domains |

## Setting Up Multisite

### New Multisite Installation

While Stax doesn't have built-in multisite commands yet, you can set it up using DDEV and WP-CLI:

```bash
# 1. Create standard WordPress installation
stax init my-network
stax setup my-network --install-wp
stax start my-network

# 2. Enable multisite in wp-config.php
cd ~/projects/my-network
ddev wp config set WP_ALLOW_MULTISITE true --raw

# 3. Install multisite network
# For subdomain setup:
ddev wp core multisite-install \
  --title="My Network" \
  --admin_user=admin \
  --admin_password=secure-pass123 \
  --admin_email=admin@example.com \
  --subdomains

# For subdirectory setup:
ddev wp core multisite-install \
  --title="My Network" \
  --admin_user=admin \
  --admin_password=secure-pass123 \
  --admin_email=admin@example.com
```

### Converting Existing Site to Multisite

```bash
# 1. Backup everything first!
stax wp db export backup-before-multisite.sql

# 2. Enable multisite capability
ddev wp config set WP_ALLOW_MULTISITE true --raw

# 3. Visit wp-admin and follow network setup
open https://my-site.ddev.site/wp-admin/network.php

# 4. Add multisite constants to wp-config.php
ddev wp config set MULTISITE true --raw
ddev wp config set SUBDOMAIN_INSTALL true --raw
ddev wp config set DOMAIN_CURRENT_SITE "my-network.ddev.site"
ddev wp config set PATH_CURRENT_SITE "/"
ddev wp config set SITE_ID_CURRENT_SITE 1 --raw
ddev wp config set BLOG_ID_CURRENT_SITE 1 --raw

# 5. Update .htaccess (for subdirectory installs)
# WordPress will provide the rules - copy them
```

### Configure DDEV for Multisite

For subdomain multisite, configure DDEV to handle wildcards:

```bash
# Configure wildcard domains
ddev config --additional-fqdns="*.my-network.ddev.site"

# For multiple specific subdomains
ddev config --additional-fqdns="site1.my-network.ddev.site,site2.my-network.ddev.site"

# Restart to apply
ddev restart
```

## Managing Subsites

### Creating New Sites

```bash
# Create a new subsite
ddev wp site create \
  --slug=marketing \
  --title="Marketing Site" \
  --email=marketing@company.com

# List all sites
ddev wp site list

# Output:
# +---------+-------------------------+---------------------+
# | blog_id | url                     | registered          |
# +---------+-------------------------+---------------------+
# | 1       | https://network.ddev.site/      | 2024-01-01 |
# | 2       | https://marketing.network.ddev.site/ | 2024-01-15 |
# +---------+-------------------------+---------------------+
```

### Working with Specific Sites

```bash
# Run commands on specific site
ddev wp --url=marketing.my-network.ddev.site option get blogname

# Install plugin on specific site
ddev wp plugin install akismet --activate \
  --url=marketing.my-network.ddev.site

# Install plugin network-wide
ddev wp plugin install akismet --activate-network

# Switch between sites in commands
ddev wp --url=site1.network.ddev.site post list
ddev wp --url=site2.network.ddev.site post list
```

### Managing Users Across Sites

```bash
# Create user with access to specific site
ddev wp user create john john@example.com \
  --role=editor \
  --url=marketing.network.ddev.site

# Grant existing user access to another site
ddev wp user set-role existing-user editor \
  --url=newsite.network.ddev.site

# Make user super admin (access to all sites)
ddev wp super-admin add username
```

## Syncing Multisite from WP Engine

### Full Network Sync

```bash
# 1. Sync entire multisite database
stax wpe sync my-network-install --skip-files

# 2. Update all domain references
# Main site
ddev wp search-replace \
  "network.wpengine.com" \
  "network.ddev.site" \
  --network \
  --skip-columns=guid

# Each subsite
ddev wp search-replace \
  "site1.network.com" \
  "site1.network.ddev.site" \
  --network \
  --skip-columns=guid
```

### Selective Site Sync

For large networks, you might want to sync only specific sites:

```bash
# 1. Export specific site tables from WP Engine
ssh install@install.ssh.wpengine.net
cd sites/install

# Export main site tables (blog_id 1)
wp db export main-site.sql --tables=$(wp db tables --all-tables --format=csv | grep -E "^wp_[^0-9]" | tr '\n' ',')

# Export specific subsite (e.g., blog_id 3)
wp db export site3.sql --tables=$(wp db tables --all-tables --format=csv | grep -E "(^wp_3_|^wp_site|^wp_blogs|^wp_users|^wp_usermeta)" | tr '\n' ',')

exit

# 2. Download and import
scp install@install.ssh.wpengine.net:sites/install/site3.sql ./
ddev import-db --src=site3.sql
```

### Media Files for Multisite

Multisite stores uploads in separate directories:

```bash
# Main site uploads: wp-content/uploads/
# Subsite uploads: wp-content/uploads/sites/[blog_id]/

# Sync specific site's media
rsync -avz \
  install@install.ssh.wpengine.net:/sites/install/wp-content/uploads/sites/2/ \
  ./wp-content/uploads/sites/2/

# Or use media redirect for all sites
stax wpe sync install --skip-files --create-upload-redirect
```

## Domain Configuration

### Local Development Domains

```bash
# 1. Configure DDEV for all your domains
ddev config --additional-fqdns="\
marketing.network.ddev.site,\
sales.network.ddev.site,\
support.network.ddev.site"

# 2. Update hosts file (if needed)
sudo nano /etc/hosts

# Add:
127.0.0.1 marketing.network.ddev.site
127.0.0.1 sales.network.ddev.site
127.0.0.1 support.network.ddev.site
```

### Domain Mapping

For completely custom domains per site:

```bash
# 1. Install domain mapping plugin
ddev wp plugin install wordpress-mu-domain-mapping --activate-network

# 2. Map domains to sites
ddev wp site list  # Note blog_id for each site

# 3. Configure in Network Admin
open https://network.ddev.site/wp-admin/network/settings.php

# 4. Add domain mappings via database
ddev wp db query "INSERT INTO wp_domain_mapping (blog_id, domain, active) VALUES (2, 'marketing.local', 1);"
```

### Production Domain Structure

When syncing from production with custom domains:

```bash
# Create a mapping script
cat > update-domains.sh << 'EOF'
#!/bin/bash

# Main network domain
ddev wp search-replace "network.com" "network.ddev.site" --network --skip-columns=guid

# Subsite domains
ddev wp search-replace "marketing.company.com" "marketing.network.ddev.site" --network --skip-columns=guid
ddev wp search-replace "sales.company.com" "sales.network.ddev.site" --network --skip-columns=guid

# Update site URLs in database
ddev wp db query "UPDATE wp_blogs SET domain = 'marketing.network.ddev.site' WHERE blog_id = 2;"
ddev wp db query "UPDATE wp_blogs SET domain = 'sales.network.ddev.site' WHERE blog_id = 3;"
EOF

chmod +x update-domains.sh
./update-domains.sh
```

## Common Multisite Tasks

### Network-Wide Operations

```bash
# Activate theme for all sites
for site in $(ddev wp site list --field=url); do
  ddev wp theme activate twentytwentythree --url=$site
done

# Deactivate plugin on all sites
ddev wp plugin deactivate akismet --network

# Run search-replace across network
ddev wp search-replace "old-text" "new-text" --network

# Export entire network database
ddev wp db export full-network-backup.sql
```

### Site Cloning

Clone an existing site within the network:

```bash
# 1. Create new empty site
ddev wp site create --slug=clone --title="Cloned Site"

# 2. Get new site's blog_id
NEW_ID=$(ddev wp site list --field=blog_id --url=clone.network.ddev.site)

# 3. Copy tables from source site (e.g., blog_id 2)
SOURCE_ID=2
ddev wp db query "SHOW TABLES LIKE 'wp_${SOURCE_ID}_%'" | while read table; do
  NEW_TABLE=$(echo $table | sed "s/wp_${SOURCE_ID}_/wp_${NEW_ID}_/")
  ddev wp db query "CREATE TABLE $NEW_TABLE LIKE $table"
  ddev wp db query "INSERT INTO $NEW_TABLE SELECT * FROM $table"
done

# 4. Update URLs in cloned site
ddev wp search-replace "source.network.ddev.site" "clone.network.ddev.site" \
  --url=clone.network.ddev.site
```

### Maintenance Tasks

```bash
# Optimize all site databases
for site in $(ddev wp site list --field=url); do
  echo "Optimizing $site"
  ddev wp db optimize --url=$site
done

# Clear all caches
ddev wp cache flush --network

# Check all sites are accessible
for site in $(ddev wp site list --field=url); do
  echo -n "Checking $site: "
  curl -s -o /dev/null -w "%{http_code}" $site
  echo ""
done
```

## Troubleshooting

### Sites Not Loading

**Problem:** Subsite returns 404 or redirects to main site.

**Solutions:**

1. **Check DDEV configuration:**
```bash
ddev describe | grep URLs
# Should show wildcard or all subdomains
```

2. **Verify database entries:**
```bash
ddev wp site list
ddev wp db query "SELECT * FROM wp_blogs"
```

3. **Check multisite constants:**
```bash
ddev wp config get MULTISITE
ddev wp config get SUBDOMAIN_INSTALL
```

### Login Issues

**Problem:** Can't log into subsites.

**Solutions:**

```bash
# 1. Check cookies domain
ddev wp config set COOKIE_DOMAIN ""

# 2. Clear browser cookies
# or use incognito mode

# 3. Verify user exists on subsite
ddev wp user list --url=subsite.network.ddev.site
```

### Plugin/Theme Visibility

**Problem:** Plugins/themes not showing on subsites.

**Solutions:**

```bash
# Network activate plugins
ddev wp plugin activate plugin-name --network

# Enable themes for network
ddev wp theme enable theme-name --network

# Or activate per site
ddev wp theme activate theme-name --url=subsite.network.ddev.site
```

### Database Table Issues

**Problem:** "Table doesn't exist" errors.

**Solutions:**

```bash
# Check which tables exist
ddev wp db query "SHOW TABLES"

# Create missing subsite tables
ddev wp site create --slug=temp --title="Temporary"
# This creates the table structure

# Then copy structure for your site
ddev wp db query "CREATE TABLE wp_3_posts LIKE wp_2_posts"
```

### URL Rewrite Problems

**Problem:** Links point to wrong domain.

**Solutions:**

```bash
# 1. Update site URL in database
ddev wp db query "UPDATE wp_blogs SET domain = 'correct.ddev.site' WHERE blog_id = 2"

# 2. Clear all caches
ddev wp cache flush

# 3. Run search-replace
ddev wp search-replace "wrong.domain" "correct.domain" --network
```

## Best Practices

### Development Workflow

1. **Use subdirectory for development** - Easier domain management
2. **Sync specific sites** - Don't sync entire network if not needed
3. **Use media redirects** - Avoid downloading all subsites' media
4. **Document site IDs** - Keep a mapping of blog_id to purpose

### Performance

1. **Limit active plugins** - Network-activated plugins load on all sites
2. **Use object caching** - Redis/Memcached helps multisite significantly
3. **Monitor database size** - Each site adds tables

### Security

1. **Limit super admin accounts** - Use site-specific admins
2. **Regular updates** - One vulnerable plugin affects all sites
3. **Separate staging/production** - Test on staging network first

## Helper Scripts

### Multisite Setup Script

Save as `setup-multisite.sh`:

```bash
#!/bin/bash

PROJECT=$1
TYPE=${2:-subdomain}  # subdomain or subdirectory

if [ -z "$PROJECT" ]; then
  echo "Usage: ./setup-multisite.sh project-name [subdomain|subdirectory]"
  exit 1
fi

# Create and setup
stax init $PROJECT
stax setup $PROJECT --install-wp
stax start $PROJECT

cd ~/projects/$PROJECT

# Enable multisite
ddev wp config set WP_ALLOW_MULTISITE true --raw

# Install network
if [ "$TYPE" = "subdomain" ]; then
  ddev wp core multisite-install \
    --title="$PROJECT Network" \
    --admin_user=admin \
    --admin_password=admin123 \
    --admin_email=admin@example.com \
    --subdomains

  # Configure DDEV for subdomains
  ddev config --additional-fqdns="*.$PROJECT.ddev.site"
else
  ddev wp core multisite-install \
    --title="$PROJECT Network" \
    --admin_user=admin \
    --admin_password=admin123 \
    --admin_email=admin@example.com
fi

ddev restart

echo "Multisite ready at https://$PROJECT.ddev.site"
echo "Network admin: https://$PROJECT.ddev.site/wp-admin/network/"
```

## Next Steps

- Review [User Guide](USER_GUIDE.md) for general workflows
- Check [WP Engine Guide](WPENGINE.md) for syncing multisite
- See [Troubleshooting](TROUBLESHOOTING.md) for common issues
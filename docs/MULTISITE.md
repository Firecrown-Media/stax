# Working with WordPress Multisite

A comprehensive guide to WordPress multisite development with Stax.

---

## Table of Contents

- [Understanding Multisite](#understanding-multisite)
- [Subdomain vs Subdirectory](#subdomain-vs-subdirectory)
- [How Stax Handles Multisite](#how-stax-handles-multisite)
- [Setting Up Multisite](#setting-up-multisite)
- [Working with Subsites](#working-with-subsites)
- [URL Structure](#url-structure)
- [Database Considerations](#database-considerations)
- [Theme and Plugin Management](#theme-and-plugin-management)
- [Troubleshooting Multisite](#troubleshooting-multisite)

---

## Understanding Multisite

### What is WordPress Multisite?

WordPress Multisite is a feature that allows you to run multiple WordPress sites from a single WordPress installation. Think of it as one WordPress installation managing many websites.

**Key concepts**:
- **Network**: The main WordPress installation that manages all sites
- **Sites** (or subsites): Individual websites within the network
- **Network Admin**: Special admin area for managing the entire network
- **Site Admin**: Regular WordPress admin for individual sites

**Common use cases**:
- Running multiple brand websites (e.g., Flying Magazine, Plane & Pilot, AVweb)
- Multi-tenant SaaS platforms
- University/school websites (one per department)
- Multi-language sites
- Client site management

### Multisite vs Multiple WordPress Installations

**Why use multisite?**

**Advantages**:
- Share code across all sites (one WordPress core, one set of plugins)
- Manage all sites from one place
- Easier updates (update once, applies to all sites)
- Shared users and authentication
- More efficient hosting (one database, one codebase)

**Disadvantages**:
- More complex than single sites
- Plugin compatibility issues (not all plugins work with multisite)
- Performance considerations (all sites share resources)
- Harder to separate sites later if needed

**When to use multisite**:
- You manage multiple related sites
- Sites share similar functionality
- You want centralized management
- Users need access to multiple sites

**When NOT to use multisite**:
- Sites are completely unrelated
- Different clients who should be isolated
- Very different performance requirements
- Maximum isolation needed

---

## Subdomain vs Subdirectory

WordPress multisite supports two modes. Understanding the difference is crucial.

### Subdomain Mode

**URL structure**:
```
Main site:   https://example.com
Site 1:      https://site1.example.com
Site 2:      https://site2.example.com
Site 3:      https://site3.example.com
```

**Characteristics**:
- Each site has its own subdomain
- Most "separate" feeling for users
- Requires wildcard DNS
- More flexibility with custom domains later
- **This is what most enterprises use**

**When to use**:
- Sites should feel separate and independent
- You might add custom domains later (e.g., site1.com → site1.example.com)
- Professional/enterprise use
- Better for SEO (each site can have its own identity)

### Subdirectory Mode

**URL structure**:
```
Main site:   https://example.com
Site 1:      https://example.com/site1
Site 2:      https://example.com/site2
Site 3:      https://example.com/site3
```

**Characteristics**:
- Sites are in subdirectories
- All under one domain
- Simpler DNS setup (no wildcard needed)
- All sites share domain authority (SEO)
- Harder to separate later

**When to use**:
- Sites are closely related
- Internal/intranet use
- All sites should share brand/domain
- Simpler DNS requirements

### Stax Default

Stax defaults to **subdomain mode** because:
- It's more common in professional/enterprise settings
- It's what WPEngine recommends
- It provides better separation between sites
- It's easier to add custom domains later

You can still use subdirectory mode - just specify it during `stax init`.

---

## How Stax Handles Multisite

Stax makes WordPress multisite development dramatically easier than doing it manually.

### Automatic Subdomain Configuration

**In subdomain mode**, Stax automatically:

1. **Configures wildcard DNS**:
   - You don't need to edit `/etc/hosts`
   - DDEV's router handles DNS automatically
   - All subdomains work instantly

2. **Generates SSL certificates**:
   - Wildcard SSL for `*.yoursite.local`
   - All subsites have HTTPS automatically
   - No certificate warnings

3. **Sets up nginx configuration**:
   - Proper multisite rewrites
   - Subdomain routing
   - Remote media proxying per site

4. **Handles search-replace**:
   - Maps production domains to local domains
   - Runs for network + each site
   - Updates all URL references

### Configuration Example

Here's what Stax does behind the scenes:

**Your `.stax.yml`**:
```yaml
project:
  name: firecrown-multisite
  type: wordpress-multisite
  mode: subdomain

network:
  domain: firecrown.local
  sites:
    - name: flyingmag
      domain: flyingmag.firecrown.local
      wpengine_domain: flyingmag.com

    - name: planeandpilot
      domain: planeandpilot.firecrown.local
      wpengine_domain: planeandpilotmag.com
```

**Generated DDEV config** (`.ddev/config.yaml`):
```yaml
name: firecrown-multisite
type: wordpress
additional_fqdns:
  - "*.firecrown.local"
  - flyingmag.firecrown.local
  - planeandpilot.firecrown.local
# ... more config
```

**Automatic search-replace**:
```bash
# Network
wp search-replace 'firecrown.wpengine.com' 'firecrown.local' --network

# Each site
wp search-replace 'flyingmag.com' 'flyingmag.firecrown.local' \
  --url=flyingmag.com

wp search-replace 'planeandpilotmag.com' 'planeandpilot.firecrown.local' \
  --url=planeandpilotmag.com
```

You don't have to do any of this manually - Stax handles it all.

---

## Setting Up Multisite

### New Multisite Project

**Initialize a subdomain multisite**:

```bash
cd ~/Sites/my-multisite
stax init
```

When prompted:

```
? Multisite mode:
  ❯ subdomain
    subdirectory

? Network domain: my-multisite.local

? Add a site? Yes
? Site name: site1
? Local domain: site1.my-multisite.local
? WPEngine domain: site1.example.com

? Add another site? Yes
? Site name: site2
? Local domain: site2.my-multisite.local
? WPEngine domain: site2.example.com

? Add another site? No
```

Stax will:
1. Set up multisite
2. Configure all sites
3. Import database
4. Run search-replace for each site
5. Configure SSL for all domains

**Result**:
- Main network: https://my-multisite.local
- Site 1: https://site1.my-multisite.local
- Site 2: https://site2.my-multisite.local

### Existing Multisite Project

If you're joining a team with an existing Stax multisite:

```bash
git clone https://github.com/mycompany/multisite-project.git
cd multisite-project

# .stax.yml already exists with multisite config
stax init
```

Stax reads the existing config and sets up everything automatically.

### Converting Single Site to Multisite

**Not recommended** - it's complex and error-prone.

If you must:
1. Back up everything
2. Follow WordPress documentation for multisite conversion
3. Update `.stax.yml` with multisite config
4. Run `stax restart` to regenerate DDEV config

Better: Start a new multisite and migrate content.

---

## Working with Subsites

### Listing All Sites

```bash
stax wp site list
```

**Output**:
```
+----+--------------------------------+--------+
| ID | url                            | state  |
+----+--------------------------------+--------+
| 1  | my-multisite.local             | public |
| 2  | site1.my-multisite.local       | public |
| 3  | site2.my-multisite.local       | public |
+----+--------------------------------+--------+
```

**More details**:
```bash
stax wp site list --fields=blog_id,url,registered,last_updated
```

### Accessing Subsites

**Web browser**:
- Network admin: https://my-multisite.local/wp-admin/network/
- Site 1 admin: https://site1.my-multisite.local/wp-admin/
- Site 2 admin: https://site2.my-multisite.local/wp-admin/

**WP-CLI** (specify site with `--url`):
```bash
# Run command on specific site
stax wp plugin list --url=site1.my-multisite.local

# Flush cache for specific site
stax wp cache flush --url=site1.my-multisite.local

# Get option from specific site
stax wp option get siteurl --url=site1.my-multisite.local
```

**Run command on all sites**:
```bash
# Network-wide
stax wp cache flush --network

# Loop through all sites
for site in $(stax wp site list --field=url); do
  echo "Flushing cache for $site"
  stax wp cache flush --url=$site
done
```

### Creating a New Site

**Via WP-CLI**:
```bash
# Basic creation
stax wp site create --slug=site3

# With full options
stax wp site create \
  --slug=site3 \
  --title="Site 3" \
  --email=admin@site3.com
```

**Via network admin**:
1. Go to https://my-multisite.local/wp-admin/network/
2. Sites → Add New
3. Fill in details
4. Click "Add Site"

**Update Stax config**:

After creating a site, add it to `.stax.yml`:

```yaml
network:
  sites:
    - name: site3
      domain: site3.my-multisite.local
      wpengine_domain: site3.example.com
```

Then restart:
```bash
stax restart
```

This ensures:
- Search-replace includes the new site
- SSL certificate covers the new domain
- Configuration is saved for team

### Deleting a Site

**Via WP-CLI**:
```bash
# Delete site (removes all content)
stax wp site delete 3

# Archive instead (keep content, hide site)
stax wp site archive 3
```

**Via network admin**:
1. Network Admin → Sites
2. Hover over site → Delete
3. Confirm deletion

**Important**: Deleting a site removes all its content permanently!

**Empty a site** (delete content but keep site):
```bash
stax wp site empty 3 --yes
```

### Site-Specific Operations

**Activate plugin for one site only**:
```bash
stax wp plugin activate wordpress-seo --url=site1.my-multisite.local
```

**Activate for all sites**:
```bash
stax wp plugin activate wordpress-seo --network
```

**Change theme for one site**:
```bash
stax wp theme activate twentytwentyfour --url=site1.my-multisite.local
```

**Export site-specific database**:
```bash
stax wp db export site1.sql \
  --tables=$(stax wp db tables --url=site1.my-multisite.local --format=csv)
```

---

## URL Structure

### Understanding Multisite URLs

**Database URLs** (stored in database):
- Network: `firecrown.wpengine.com`
- Site 1: `flyingmag.com`
- Site 2: `planeandpilotmag.com`

**Local URLs** (after search-replace):
- Network: `firecrown.local`
- Site 1: `flyingmag.firecrown.local`
- Site 2: `planeandpilot.firecrown.local`

**How Stax maps them**:

In `.stax.yml`:
```yaml
network:
  domain: firecrown.local

network sites:
    - name: flyingmag
      domain: flyingmag.firecrown.local
      wpengine_domain: flyingmag.com
```

When you run `stax db pull`, Stax:
1. Imports the database
2. Runs: `wp search-replace flyingmag.com flyingmag.firecrown.local --url=flyingmag.com`
3. Updates all references

### Custom Domain Mapping

If production uses custom domains:

**Production**:
- Site lives at: `site1.mainsite.com`
- Custom domain: `customdomain.com` → points to site1

**Local** (simple approach):
- Just use: `site1.mainsite.local`
- Custom domain won't work locally (and that's OK)

**Local** (advanced - if you need custom domain locally):

Add to `.stax.yml`:
```yaml
network:
  sites:
    - name: site1
      domain: site1.mainsite.local
      wpengine_domain: site1.mainsite.com
      custom_domains:
        - customdomain.local
```

Update DDEV config to include custom domain:
```yaml
additional_fqdns:
  - customdomain.local
```

Restart:
```bash
stax restart
```

### wp-config.php Multisite Settings

**Generated by WordPress during multisite setup**:

```php
/* Multisite */
define( 'WP_ALLOW_MULTISITE', true );
define( 'MULTISITE', true );
define( 'SUBDOMAIN_INSTALL', true );  // true for subdomain, false for subdirectory
define( 'DOMAIN_CURRENT_SITE', 'firecrown.local' );
define( 'PATH_CURRENT_SITE', '/' );
define( 'SITE_ID_CURRENT_SITE', 1 );
define( 'BLOG_ID_CURRENT_SITE', 1 );
```

**Stax updates** `DOMAIN_CURRENT_SITE` during search-replace to match your local domain.

---

## Database Considerations

### Multisite Database Structure

**Shared tables** (all sites use these):
```
wp_users
wp_usermeta
wp_sitemeta
```

**Per-site tables** (site 1 = blog ID 1):
```
wp_posts           # Main site (blog_id = 1)
wp_postmeta
wp_options
wp_commentmeta
wp_comments
wp_term_relationships
wp_term_taxonomy
wp_termmeta
wp_terms
```

**Site 2 tables** (blog ID 2):
```
wp_2_posts
wp_2_postmeta
wp_2_options
... etc
```

**Site 3 tables** (blog ID 3):
```
wp_3_posts
wp_3_postmeta
wp_3_options
... etc
```

### Search-Replace in Multisite

**Challenge**: Each site has different URLs.

**Stax's solution**:

1. **Network-wide replace** (for network options):
   ```bash
   wp search-replace \
     'firecrown.wpengine.com' \
     'firecrown.local' \
     --network
   ```

2. **Site-specific replace** (for each site's tables):
   ```bash
   wp search-replace \
     'flyingmag.com' \
     'flyingmag.firecrown.local' \
     --url=flyingmag.com
   ```

3. **Repeat for each site**

**Important**: Always use `--url` to specify which site's tables to update.

### Site-Specific Database Operations

**Export one site's data**:
```bash
# Get list of tables for site
TABLES=$(stax wp db tables --url=site1.my-multisite.local --format=csv)

# Export those tables
stax wp db export site1-export.sql --tables=$TABLES
```

**Import data to specific site**:
```bash
# Import SQL file
stax wp db import site1-data.sql

# Run search-replace for that site
stax wp search-replace \
  'production-domain.com' \
  'site1.my-multisite.local' \
  --url=site1.my-multisite.local
```

---

## Theme and Plugin Management

### Network-Wide vs Site-Specific

**Network Activated** (available to all sites):
- Set in Network Admin → Plugins → Network Activate
- All sites can use the plugin
- Individual sites can't deactivate
- Use for must-have plugins (security, performance, etc.)

**Site Activated** (per-site basis):
- Regular admin → Plugins → Activate
- Only that site uses the plugin
- Site admin can activate/deactivate
- Use for site-specific functionality

### Managing Plugins

**Network activate**:
```bash
stax wp plugin activate wordpress-seo --network
```

**Install and network activate**:
```bash
stax wp plugin install wordpress-seo --activate-network
```

**Activate for specific site only**:
```bash
stax wp plugin activate contact-form-7 --url=site1.my-multisite.local
```

**List network-activated plugins**:
```bash
stax wp plugin list --status=active-network
```

**Deactivate network-wide**:
```bash
stax wp plugin deactivate wordpress-seo --network
```

### Managing Themes

**Install theme** (available to all sites):
```bash
stax wp theme install twentytwentyfour
```

**Activate for specific site**:
```bash
stax wp theme activate twentytwentyfour --url=site1.my-multisite.local
```

**Enable theme network-wide** (so sites can choose it):
```bash
stax wp theme enable twentytwentyfour --network
```

**List enabled themes**:
```bash
stax wp theme list --status=enabled
```

### Custom Plugins and Themes

**Must-use plugins** (mu-plugins):
- Located in `wp-content/mu-plugins/`
- Always active on all sites
- Can't be deactivated
- Perfect for network-wide customizations

**Per-site themes**:
- Each site can use a different theme
- Manage via site admin
- Or via WP-CLI with `--url`

---

## Troubleshooting Multisite

### Subdomain Not Accessible

**Symptom**: https://site1.my-multisite.local doesn't load.

**Checks**:

1. **Verify DDEV config**:
   ```bash
   cat .ddev/config.yaml | grep additional_fqdns
   ```
   Should include your subdomain.

2. **Verify site exists**:
   ```bash
   stax wp site list
   ```
   Should show the site.

3. **Check database URL**:
   ```bash
   stax wp option get siteurl --url=site1.my-multisite.local
   ```
   Should be `https://site1.my-multisite.local`.

4. **Restart DDEV**:
   ```bash
   stax restart
   ```

5. **Check nginx logs**:
   ```bash
   stax logs -f
   ```

**Fix**:
```bash
# Make sure .stax.yml has the site
cat .stax.yml

# Restart to regenerate DDEV config
stax restart

# Run search-replace
stax wp search-replace \
  'production-domain.com' \
  'site1.my-multisite.local' \
  --url=site1.my-multisite.local
```

### SSL Certificate Errors

**Symptom**: Browser shows "Not secure" or certificate error.

**Cause**: DDEV generates wildcard certificates, but browser might not trust them.

**Solution**:
1. **Trust mkcert CA**:
   ```bash
   mkcert -install
   ```

2. **Restart DDEV**:
   ```bash
   stax restart
   ```

3. **Clear browser cache**

4. **In browser**: Accept the certificate
   - Chrome: "Advanced" → "Proceed"
   - Firefox: "Advanced" → "Accept Risk"

### URLs Not Updating After Database Pull

**Symptom**: Site shows production URLs instead of local.

**Cause**: Search-replace didn't run or failed.

**Fix**:
```bash
# Manual search-replace for network
stax wp search-replace \
  'production-network.com' \
  'my-multisite.local' \
  --network

# For each site
stax wp search-replace \
  'site1-production.com' \
  'site1.my-multisite.local' \
  --url=site1-production.com

# Flush cache
stax wp cache flush --network
```

### Main Site Works, Subsites Don't

**Symptom**: https://my-multisite.local works, but https://site1.my-multisite.local doesn't.

**Checks**:

1. **Multisite mode**:
   ```bash
   stax wp eval 'echo SUBDOMAIN_INSTALL ? "subdomain" : "subdirectory";'
   ```
   Should match your `.stax.yml` mode.

2. **Site URL**:
   ```bash
   stax wp site list --fields=blog_id,url
   ```
   URLs should be correct.

3. **DNS resolution**:
   ```bash
   ping site1.my-multisite.local
   ```
   Should resolve (usually to 127.0.0.1).

**Fix**:

Update `.ddev/config.yaml` (or regenerate via Stax):
```yaml
additional_fqdns:
  - "*.my-multisite.local"
  - site1.my-multisite.local
  - site2.my-multisite.local
```

Restart:
```bash
ddev restart
```

### Site Shows Wrong Content

**Symptom**: Site 1 shows Site 2's content (or vice versa).

**Cause**: URL mapping is wrong in database.

**Fix**:

1. **Check site URL**:
   ```bash
   stax wp option get siteurl --url=site1.my-multisite.local
   ```

2. **Update if wrong**:
   ```bash
   stax wp option update siteurl 'https://site1.my-multisite.local' \
     --url=site1.my-multisite.local

   stax wp option update home 'https://site1.my-multisite.local' \
     --url=site1.my-multisite.local
   ```

3. **Run search-replace**:
   ```bash
   stax wp search-replace \
     'wrong-url.com' \
     'site1.my-multisite.local' \
     --url=site1.my-multisite.local
   ```

4. **Flush cache**:
   ```bash
   stax wp cache flush --network
   ```

### Can't Access Network Admin

**Symptom**: https://my-multisite.local/wp-admin/network/ shows "You do not have permission."

**Cause**: Your user isn't a Super Admin.

**Fix**:
```bash
# List current Super Admins
stax wp super-admin list

# Grant Super Admin to your user
stax wp super-admin add yourusername

# Verify
stax wp super-admin list
```

### New Site Not Appearing

**Symptom**: Created a site but it doesn't show up.

**Checks**:

1. **Verify site exists**:
   ```bash
   stax wp site list
   ```

2. **Check site status**:
   ```bash
   stax wp site list --fields=blog_id,url,public,archived,mature,spam,deleted
   ```

3. **Activate if archived**:
   ```bash
   stax wp site activate <blog_id>
   ```

4. **Update Stax config**:
   Add site to `.stax.yml` and restart.

---

## Next Steps

You now understand WordPress multisite with Stax! Continue learning:

- **[User Guide](./USER_GUIDE.md)** - General Stax usage
- **[WPEngine Guide](./WPENGINE.md)** - WPEngine-specific features
- **[Examples](./EXAMPLES.md)** - Real-world multisite workflows
- **[Troubleshooting](./TROUBLESHOOTING.md)** - More problem-solving

---

**Questions about multisite?** Check the [FAQ](./FAQ.md) or ask your team!

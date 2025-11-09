# Multi-Provider Guide

## Overview

Stax 2.0 introduces support for multiple WordPress hosting providers through a pluggable provider architecture. This guide explains how to work with multiple providers, switch between them, and migrate sites across platforms.

## Supported Providers

### Production-Ready Providers

- **WPEngine** - Full support for WPEngine WordPress hosting (API + SSH)
- **Local** - Local-only development (no remote hosting)

### Preview/Development Providers

- **AWS** - Amazon Web Services (EC2, Lightsail, RDS) - Coming soon
- **WordPress VIP** - WordPress.com VIP hosting - Coming soon

## Provider Comparison Matrix

| Feature | WPEngine | AWS | WordPress VIP | Local |
|---------|----------|-----|---------------|-------|
| Authentication | API + SSH | AWS SDK + SSH | API + VIP-CLI | None |
| Site Management | Yes | Yes | Yes | Limited |
| Database Export | Yes (SSH) | Yes (SSH/RDS) | Yes (VIP-CLI) | Yes (DDEV) |
| Database Import | No (Portal) | Yes | No (Support) | Yes (DDEV) |
| File Sync | Yes (Rsync) | Yes (Rsync/S3) | Yes (Git) | No |
| Deployments | Git | CodeDeploy | Git | N/A |
| Environments | Staging | Multi-instance | Dev/Preprod/Prod | N/A |
| Backups | Automatic | EBS/RDS | Automatic | Manual |
| SSH Access | Yes (Gateway) | Yes (Direct) | No | N/A |
| WP-CLI | Yes | Yes | Yes (VIP-CLI) | Yes (DDEV) |
| CDN | BunnyCDN | CloudFront | Photon | N/A |
| Scaling | Managed | Auto-scaling | Automatic | N/A |

## Getting Started

### 1. List Available Providers

```bash
stax provider list
```

Output:
```
PROVIDER         DESCRIPTION                                DEFAULT   CAPABILITIES
--------         -----------                                -------   ------------
wpengine         WPEngine WordPress Hosting Platform        *         5 core, 5 optional
aws              Amazon Web Services (EC2, Lightsail, RDS)            5 core, 8 optional
wordpress-vip    WordPress VIP (WordPress.com VIP Hosting)            4 core, 7 optional
local            Local Development Only (No Remote Hosting)           3 core, 0 optional
```

### 2. View Provider Details

```bash
stax provider show wpengine
```

### 3. Set Default Provider

```bash
# For current project
stax provider set wpengine

# Or set in .stax.yml
provider:
  name: wpengine
```

## Provider-Specific Setup

### WPEngine Setup

1. Obtain API credentials from WPEngine portal
2. Generate SSH key and add to WPEngine account
3. Configure credentials:

```yaml
# .stax.yml
provider:
  name: wpengine
  wpengine:
    site: my-site
    environment: production
```

```yaml
# ~/.stax/credentials.yml
wpengine:
  api_user: your-api-user
  api_password: your-api-password
  ssh_key: |
    -----BEGIN OPENSSH PRIVATE KEY-----
    ...
    -----END OPENSSH PRIVATE KEY-----
```

**See**: [PROVIDER_WPENGINE.md](./PROVIDER_WPENGINE.md) for complete WPEngine documentation.

### AWS Setup (Coming Soon)

1. Configure AWS credentials
2. Launch EC2 or Lightsail instance with WordPress
3. Configure Stax:

```yaml
provider:
  name: aws
  aws:
    region: us-east-1
    instance_id: i-1234567890abcdef
    ssh_user: ubuntu
    ssh_key_path: ~/.ssh/aws-key.pem
```

### WordPress VIP Setup (Coming Soon)

1. Obtain VIP Dashboard API token
2. Configure Stax:

```yaml
provider:
  name: wordpress-vip
  wordpress_vip:
    app_id: 12345
    org: my-org
    environment: production
```

### Local-Only Setup

For projects without remote hosting:

```yaml
provider:
  name: local
```

## Working with Multiple Providers

### Switching Providers

You can switch providers for a project at any time:

```bash
# Switch to AWS
stax provider set aws

# Or use --provider flag for one-off commands
stax db pull --provider=wpengine
stax db pull --provider=aws
```

### Multi-Provider Workflows

#### Development → Staging → Production

```yaml
# .stax.yml
provider:
  name: wpengine
  wpengine:
    production:
      site: mysite
      environment: production
    staging:
      site: mysite
      environment: staging
```

```bash
# Pull from staging
stax db pull --provider=wpengine --environment=staging

# Pull from production
stax db pull --provider=wpengine --environment=production
```

#### Testing Multiple Hosting Platforms

```bash
# Pull from WPEngine
stax db pull --provider=wpengine

# Compare with AWS
stax db pull --provider=aws

# Compare providers
stax provider compare wpengine aws
```

## Migrating Between Providers

### Migration Strategies

**1. Manual Migration** (Current Support)

Export from source provider → Import to DDEV → Deploy to target provider

```bash
# Step 1: Pull from source (WPEngine)
stax db pull --provider=wpengine
stax files sync --provider=wpengine

# Step 2: Test locally
stax start

# Step 3: Push to target (AWS)
stax db push --provider=aws
stax files push --provider=aws
```

**2. Direct Migration** (Future Enhancement)

```bash
# Future: Direct provider-to-provider migration
stax migrate --from=wpengine --to=aws
```

### Migration Checklist

Before migrating between providers:

- [ ] Backup current site
- [ ] Test locally first
- [ ] Update DNS records
- [ ] Verify SSL certificates
- [ ] Update CDN configuration
- [ ] Test media URLs
- [ ] Verify cron jobs
- [ ] Check plugin compatibility
- [ ] Update deployment workflows
- [ ] Monitor performance

### WPEngine → AWS Migration Example

```bash
# 1. Pull everything from WPEngine
stax db pull --provider=wpengine
stax files sync --provider=wpengine

# 2. Test locally
stax start
stax wp search-replace wpenginesite.wpengine.com mysite.local

# 3. Configure AWS provider
# Edit .stax.yml to add AWS configuration

# 4. Deploy to AWS
stax db push --provider=aws
stax files push --provider=aws

# 5. Update DNS
# Point domain to AWS instance

# 6. Verify
stax wp --provider=aws core version
```

## Provider-Specific Limitations

### WPEngine

- **No database import**: Must use WPEngine portal
- **No direct file upload**: Use Git deployments
- **Read-only filesystem**: Except via Git
- **No root access**: Managed platform

**Workarounds**:
- For DB import: Export from DDEV, import via WPEngine portal
- For file upload: Use Git push deployments
- For configuration changes: Use wp-config.php or portal

### AWS

- **Manual scaling**: Must configure auto-scaling groups
- **Security management**: Manage security groups, firewall rules
- **Updates**: Manually manage WordPress/plugin updates
- **Backup configuration**: Set up RDS/EBS snapshot schedules

### WordPress VIP

- **Controlled imports**: Database imports via support tickets
- **Code review**: All code changes reviewed by VIP
- **No SSH**: Access via VIP-CLI only
- **Required plugins**: Must use VIP Go mu-plugins
- **Enterprise only**: Premium pricing tier

### Local

- **No remote sync**: Can't pull from remote
- **Manual backups**: No automated backup system
- **Limited features**: DDEV operations only

## Provider Selection Guidelines

### Choose WPEngine if:
- You want managed WordPress hosting
- You need automatic backups and staging
- You prefer not to manage infrastructure
- You want built-in CDN and caching
- You need WordPress-specific support

### Choose AWS if:
- You need full infrastructure control
- You want cost optimization flexibility
- You need custom server configurations
- You're comfortable managing servers
- You need AWS service integrations

### Choose WordPress VIP if:
- You need enterprise-grade hosting
- You want automatic scaling
- You require high-traffic support
- You need 24/7 expert support
- Budget allows premium tier

### Choose Local if:
- You're building a greenfield project
- You don't need remote hosting yet
- You're learning WordPress development
- You have custom deployment workflows

## Common Patterns

### Pattern 1: Local Development + WPEngine Production

```yaml
provider:
  name: wpengine  # Default to production
  wpengine:
    site: mysite
    environment: production
```

```bash
# Daily workflow
stax start                    # Start local DDEV
stax db pull                  # Pull from WPEngine
stax files sync               # Sync media files

# Deploy changes
git push wpengine main        # Deploy via Git
```

### Pattern 2: Multi-Environment Testing

```yaml
provider:
  name: wpengine
  wpengine:
    site: mysite

  aws:
    instance_id: i-1234567890abcdef
```

```bash
# Test on both platforms
stax db pull --provider=wpengine
stax db push --provider=aws

# Compare performance
stax wp --provider=wpengine option get siteurl
stax wp --provider=aws option get siteurl
```

### Pattern 3: Gradual Migration

```yaml
provider:
  name: wpengine  # Current production

  wpengine:
    site: mysite

  aws:
    instance_id: i-1234567890abcdef  # New platform
```

```bash
# Week 1: Set up AWS instance
stax db push --provider=aws
stax files push --provider=aws

# Week 2: Test AWS performance
# Monitor both platforms

# Week 3: Switch DNS to AWS
# Update .stax.yml default provider

# Week 4: Verify, then decommission WPEngine
```

## Troubleshooting

### Provider Connection Failures

```bash
# Test provider connection
stax provider test

# Verify credentials
stax provider show wpengine

# Check configuration
stax config validate
```

### Missing Capabilities

```bash
# Check if provider supports feature
stax provider show <provider-name>

# Example: Check if provider supports database import
stax provider show wpengine | grep "Database Import"
```

### Migration Issues

**Database import fails**:
- Check provider supports imports (WPEngine doesn't)
- Use provider's portal/dashboard instead
- For WPEngine: Export from DDEV, import via portal

**File sync errors**:
- Verify SSH access for provider
- Check file permissions
- Use `--dry-run` flag first

**Performance differences**:
- Providers have different caching strategies
- CDN configurations may differ
- Database performance varies by platform

## Future Enhancements

- [ ] Direct provider-to-provider migrations
- [ ] Multi-provider monitoring dashboard
- [ ] Cost comparison tools
- [ ] Performance benchmarking
- [ ] Automated provider recommendations
- [ ] Provider-specific optimizations
- [ ] Backup sync between providers

## Additional Resources

- [PROVIDER_INTERFACE.md](./PROVIDER_INTERFACE.md) - Provider interface specification
- [PROVIDER_DEVELOPMENT.md](./PROVIDER_DEVELOPMENT.md) - Creating custom providers
- [PROVIDER_WPENGINE.md](./PROVIDER_WPENGINE.md) - WPEngine-specific documentation
- [CONFIG_SPEC.md](./CONFIG_SPEC.md) - Configuration reference
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture

---

**Version**: 2.0.0
**Last Updated**: 2025-11-08

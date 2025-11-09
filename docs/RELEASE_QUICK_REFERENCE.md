# Release Quick Reference

## Quick Commands

### Create Release (Recommended)

```bash
# Via GitHub UI
1. Go to: https://github.com/firecrown-media/stax/actions
2. Select "Version Bump" workflow
3. Click "Run workflow"
4. Choose: patch/minor/major
5. Click "Run workflow"
```

### Create Release (Manual)

```bash
# Determine version
git describe --tags --abbrev=0  # Current version

# Create tag
git tag -a v1.2.3 -m "Release v1.2.3"

# Push tag (triggers release)
git push origin v1.2.3
```

## Local Testing

```bash
# Validate GoReleaser config
make release-check

# Build snapshot locally
make release-snapshot

# Test dry-run
make release-dry-run

# Check current version
make version
```

## Verify Release

```bash
# Check GitHub release created
open https://github.com/firecrown-media/stax/releases/latest

# Check Homebrew formula updated
open https://github.com/firecrown-media/homebrew-tap/blob/main/Formula/stax.rb

# Test installation
brew upgrade stax
stax --version
```

## Troubleshooting

### Release Failed

```bash
# Delete bad tag
git tag -d v1.2.3
git push origin :refs/tags/v1.2.3

# Delete GitHub release via web UI

# Fix issue, create new release
```

### Homebrew Not Updated

```bash
# Check secret exists
# Settings → Secrets → Actions → HOMEBREW_TAP_TOKEN

# Check GoReleaser logs
# Actions → Release → Latest run
```

## Version Numbers

```
vMAJOR.MINOR.PATCH

Examples:
- Bug fix:          v1.2.3 → v1.2.4 (patch)
- New feature:      v1.2.3 → v1.3.0 (minor)
- Breaking change:  v1.2.3 → v2.0.0 (major)
```

## Resources

- [Full Release Process](RELEASE_PROCESS.md)
- [CI/CD Pipeline](CICD_PIPELINE.md)
- [Deployment Summary](DEPLOYMENT_SUMMARY.md)

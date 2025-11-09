# Stax Deployment System Summary

## Overview

This document provides a complete overview of the Stax deployment and release system, which uses GoReleaser and GitHub Actions to automate the entire release process from version tagging to Homebrew distribution.

## Architecture

### Components

1. **GoReleaser** - Build automation and release management
2. **GitHub Actions** - CI/CD pipeline orchestration
3. **Homebrew Tap** - Package distribution for macOS/Linux
4. **GitHub Releases** - Binary distribution and changelog

### Workflow Diagram

```
Developer Action
      │
      ├─── Push Code ────────────► Test Workflow (CI)
      │                            │
      │                            ├─► Unit Tests
      │                            ├─► Integration Tests
      │                            ├─► Security Tests
      │                            ├─► Code Quality
      │                            └─► Build Verification
      │
      └─── Create Tag ───────────► Release Workflow
           (via Version Bump         │
           workflow or manual)       ├─► Run Tests
                                    ├─► Build Binaries (4 platforms)
                                    ├─► Create Archives
                                    ├─► Generate Checksums
                                    ├─► Create GitHub Release
                                    └─► Update Homebrew Formula
                                         │
                                         └─► Users install via:
                                             - brew install stax
                                             - Direct download
```

## Files Created

### Configuration Files

1. **/.goreleaser.yml**
   - GoReleaser configuration
   - Defines build targets, archives, and Homebrew tap
   - Configures changelog generation

2. **/.github/workflows/release.yml**
   - Automated release workflow
   - Triggered by version tags
   - Runs tests and GoReleaser

3. **/.github/workflows/version-bump.yml**
   - Manual version bump workflow
   - Calculates and creates version tags
   - Triggers release workflow

4. **/Makefile** (updated)
   - Added release-related targets
   - Local testing commands
   - Version management

### Documentation Files

1. **/docs/RELEASE_PROCESS.md**
   - Complete release process guide
   - Troubleshooting for releases
   - Rollback procedures

2. **/docs/HOMEBREW_INSTALLATION.md**
   - End-user installation guide
   - Homebrew-specific instructions
   - Troubleshooting for users

3. **/docs/HOMEBREW_TAP_SETUP.md**
   - Repository setup instructions
   - Token configuration
   - Tap maintenance guide

4. **/docs/CICD_PIPELINE.md**
   - Pipeline architecture
   - Workflow details
   - Monitoring and metrics

5. **/docs/INSTALLATION.md** (updated)
   - Homebrew as primary method
   - Direct download option added
   - Updated uninstall instructions

## Release Process

### Automated Release (Recommended)

1. **Trigger Version Bump**:
   - Go to GitHub Actions
   - Select "Version Bump" workflow
   - Choose version type (patch/minor/major)
   - Click "Run workflow"

2. **Automatic Steps**:
   - Calculates new version
   - Creates git tag
   - Triggers release workflow
   - Runs tests
   - Builds binaries
   - Creates GitHub release
   - Updates Homebrew formula

3. **User Installation**:
   ```bash
   brew update
   brew upgrade stax
   ```

### Manual Release

1. **Create Tag**:
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```

2. **Automatic Steps**: Same as automated release

## Build Targets

GoReleaser builds for:

- **macOS Intel** (darwin/amd64)
- **macOS Apple Silicon** (darwin/arm64)
- **Linux Intel** (linux/amd64)
- **Linux ARM** (linux/arm64)

All binaries are:
- Stripped and optimized (`-s -w`)
- Include version information via ldflags
- Packaged in tar.gz archives
- Verified with SHA256 checksums

## Homebrew Distribution

### Tap Repository

**Repository**: `firecrown-media/homebrew-tap`
**Formula**: `Formula/stax.rb`

### Automatic Updates

On each release, GoReleaser:
1. Updates formula with new version
2. Updates download URLs
3. Calculates and updates SHA256 checksums
4. Commits changes
5. Pushes to tap repository

### User Installation

```bash
# First time
brew tap firecrown-media/tap
brew install stax

# Updates
brew upgrade stax
```

## Required Secrets

### GITHUB_TOKEN
- **Type**: Automatic
- **Used for**: Creating releases, uploading assets
- **Setup**: Provided automatically by GitHub Actions

### HOMEBREW_TAP_TOKEN
- **Type**: Personal Access Token
- **Scope**: `repo`
- **Used for**: Pushing formula updates to tap
- **Setup**: Manual (see HOMEBREW_TAP_SETUP.md)

## Testing Locally

Before releasing, test the process locally:

```bash
# Validate GoReleaser config
make release-check

# Build snapshot (no publish)
make release-snapshot

# Full dry-run
make release-dry-run

# Test binaries
./dist/stax_darwin_amd64_v1/stax --version
```

## Monitoring

### Success Indicators

- ✅ GitHub release created
- ✅ All binaries uploaded
- ✅ Checksums generated
- ✅ Homebrew formula updated
- ✅ Users can install via brew

### Check Points

1. **GitHub Actions**: https://github.com/firecrown-media/stax/actions
2. **Releases**: https://github.com/firecrown-media/stax/releases
3. **Homebrew Tap**: https://github.com/firecrown-media/homebrew-tap
4. **Formula**: https://github.com/firecrown-media/homebrew-tap/blob/main/Formula/stax.rb

## Troubleshooting

### Common Issues

1. **Release fails**
   - Check test results
   - Validate .goreleaser.yml
   - Review GitHub Actions logs

2. **Homebrew formula not updated**
   - Verify HOMEBREW_TAP_TOKEN secret
   - Check token permissions
   - Review GoReleaser logs

3. **Binary doesn't work**
   - Test locally with `make release-snapshot`
   - Check ldflags configuration
   - Verify Go version compatibility

### Recovery

**Delete bad release**:
```bash
# Delete tag
git tag -d v1.2.3
git push origin :refs/tags/v1.2.3

# Delete GitHub release (via web UI)

# Fix issue and re-release
```

## Versioning

Follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Breaking changes
- **MINOR** (1.X.0): New features (backward compatible)
- **PATCH** (1.2.X): Bug fixes (backward compatible)

## Changelog

Automatically generated from git commits:

- Includes commits since last release
- Excludes: docs, tests, chores, CI changes
- Sorted chronologically
- Formatted as markdown

## Installation Methods

### 1. Homebrew (Recommended)

**Pros**:
- Easy updates
- Automatic dependency handling
- Verified downloads

**Installation**:
```bash
brew tap firecrown-media/tap
brew install stax
```

### 2. Direct Download

**Pros**:
- No Homebrew required
- Works on any system

**Installation**:
1. Download from GitHub Releases
2. Extract archive
3. Move binary to PATH

### 3. Build from Source

**Pros**:
- Latest development version
- Full control

**Installation**:
```bash
git clone https://github.com/firecrown-media/stax.git
cd stax
make build
make install
```

## Security

### Build Security

- Checksums verify download integrity
- Binaries stripped of debug symbols
- Dependencies pinned in go.mod
- GitHub Actions uses pinned versions

### Token Security

- Secrets stored encrypted in GitHub
- Minimal required permissions
- Regular token rotation recommended
- Never logged or exposed

### Distribution Security

- HTTPS for all downloads
- SHA256 checksums for verification
- Homebrew formula audit
- Signed commits (future enhancement)

## Maintenance

### Regular Tasks

**Weekly**:
- Monitor GitHub Actions runs
- Check for failed workflows
- Review dependency updates

**Monthly**:
- Review download statistics
- Check Homebrew analytics
- Update documentation

**Quarterly**:
- Rotate HOMEBREW_TAP_TOKEN
- Audit security practices
- Review and optimize workflows

### Upgrading Components

**GoReleaser**:
```yaml
# In .github/workflows/release.yml
uses: goreleaser/goreleaser-action@v5
with:
  version: latest  # or pin to specific version
```

**GitHub Actions**:
- Review action versions quarterly
- Test upgrades in feature branches
- Pin to specific versions for stability

## Performance

### Build Times

- Test Workflow: ~5-10 minutes
- Release Workflow: ~5-10 minutes
- Total (tag to release): ~10-15 minutes

### Optimizations

- Parallel test execution
- Go module caching
- Matrix builds for multiple platforms
- Incremental builds where possible

## Future Enhancements

### Planned

1. **Signed Releases**: GPG signatures for binaries
2. **Docker Images**: Publish to Docker Hub/GHCR
3. **Chocolatey**: Windows package manager support
4. **Snap/AppImage**: Additional Linux distribution
5. **Auto-update**: Built-in update checker
6. **Metrics**: Download and usage analytics

### Under Consideration

1. **Multi-region CDN**: Faster downloads globally
2. **Canary Releases**: Staged rollouts
3. **Beta Channel**: Early access releases
4. **Release Automation**: Auto-release on main branch
5. **Security Scanning**: Automated vulnerability checks

## Resources

### Documentation

- [Release Process](RELEASE_PROCESS.md)
- [Homebrew Installation](HOMEBREW_INSTALLATION.md)
- [Homebrew Tap Setup](HOMEBREW_TAP_SETUP.md)
- [CI/CD Pipeline](CICD_PIPELINE.md)
- [Installation Guide](INSTALLATION.md)

### External Resources

- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Semantic Versioning](https://semver.org/)

## Support

### For Developers

- Review GitHub Actions logs
- Test locally with `make release-snapshot`
- Check GoReleaser documentation
- Contact development team

### For Users

- [Troubleshooting Guide](TROUBLESHOOTING.md)
- [GitHub Issues](https://github.com/firecrown-media/stax/issues)
- [Installation Guide](INSTALLATION.md)

## Summary

The Stax deployment system provides:

✅ **Automated Testing** - Every commit is tested
✅ **Easy Releases** - One-click version bumps
✅ **Multi-Platform Builds** - macOS and Linux support
✅ **Homebrew Distribution** - Easy installation and updates
✅ **Secure Downloads** - Checksums verify integrity
✅ **Clear Documentation** - Guides for all scenarios
✅ **Monitoring** - Visibility into all releases
✅ **Recovery Procedures** - Rollback capabilities

This ensures reliable, repeatable, and secure releases for all users.

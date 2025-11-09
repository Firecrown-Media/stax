# Automated Release Process with Release Please

## Overview

Stax uses [Release Please](https://github.com/googleapis/release-please) from Google to automate version bumping, changelog generation, and GitHub releases. This eliminates manual tag creation and ensures consistent semantic versioning.

## How It Works

### 1. Commit Format

Release Please reads your commit messages to determine version bumps:

```bash
# Patch release (1.0.0 → 1.0.1)
fix: resolve database connection timeout
fix(auth): handle expired credentials gracefully

# Minor release (1.0.0 → 1.1.0)
feat: add support for WordPress VIP provider
feat(cli): implement interactive mode for init command

# Major release (1.0.0 → 2.0.0)
feat!: redesign provider interface
feat(api)!: remove deprecated --legacy flag

BREAKING CHANGE: Provider interface now requires Validate() method
```

### 2. Commit Types

| Type | Version Bump | Description | In Changelog? |
|------|--------------|-------------|---------------|
| `feat` | Minor | New feature | ✅ Yes |
| `fix` | Patch | Bug fix | ✅ Yes |
| `perf` | Patch | Performance improvement | ✅ Yes |
| `docs` | Patch | Documentation changes | ✅ Yes |
| `refactor` | Patch | Code refactoring | ✅ Yes |
| `test` | None | Test changes | ❌ No |
| `build` | None | Build system changes | ❌ No |
| `ci` | None | CI/CD changes | ❌ No |
| `chore` | None | Other changes | ❌ No |

### 3. Breaking Changes

Add `!` after type or include `BREAKING CHANGE:` in commit body:

```bash
# Method 1: ! in type
git commit -m "feat!: redesign authentication flow"

# Method 2: BREAKING CHANGE in body
git commit -m "feat: redesign authentication flow

BREAKING CHANGE: AuthProvider interface changed. All providers must implement new Validate() method."
```

## Release Workflow

### Step 1: Create Feature Branch

```bash
# Create feature branch from main
git checkout main
git pull origin main
git checkout -b feature/add-aws-provider

# Work on your feature
# Commit with conventional format
git commit -m "feat: add AWS provider implementation"
git commit -m "docs: update provider documentation"
git commit -m "test: add tests for AWS provider"

# Push feature branch to GitHub
git push origin feature/add-aws-provider
```

### Step 2: Create Pull Request

```bash
# Create PR from feature branch to main
gh pr create \
  --base main \
  --head feature/add-aws-provider \
  --title "Add AWS Provider" \
  --body "Implements AWS provider support

## Changes
- Complete AWS provider implementation
- Updated documentation
- Added comprehensive tests"

# Wait for code review and approval
```

### Step 3: Merge PR to Main

```bash
# After PR approval, merge to main
# Use squash or merge (squash recommended for cleaner history)
gh pr merge --squash

# OR merge via GitHub UI
# The PR merge to main triggers Release Please
```

### Step 4: Release Please Creates Release PR

**Automatically triggers when feature PR is merged to main.**

After your feature branch is merged, Release Please will:
1. Analyze all commits in the merge (from your feature branch)
2. Determine version bump (major/minor/patch)
3. Generate CHANGELOG.md with all changes
4. Update version in `cmd/root.go`
5. Create a "Release PR" with all changes

**Example Release PR:**
```
Title: chore(main): release 1.1.0

Changes:
- Update CHANGELOG.md
- Update version in cmd/root.go to 1.1.0

Includes commits from:
- feat: add AWS provider implementation
- docs: update provider documentation
- test: add tests for AWS provider
```

**Important:** Release Please runs on every merge to main, but only creates a Release PR if there are version-bumping commits (feat, fix, etc.)

### Step 5: Review and Merge Release PR

```bash
# Review the release PR
gh pr view <PR-NUMBER>

# Check the version bump is correct
# Check the changelog looks good

# Merge the release PR
gh pr merge <PR-NUMBER> --merge
```

### Step 6: Release Please Creates GitHub Release

**Automatically triggers when Release PR is merged.**

Release Please will:
1. Create git tag (e.g., `v1.1.0`)
2. Create GitHub release
3. Trigger release workflow

### Step 7: Release Workflow Builds and Publishes

**Automatically triggered by Release Please.**

Our release workflow will:
1. Run full test suite
2. Run security tests
3. Build binaries for 4 platforms
4. Create archives
5. Upload to GitHub release
6. Update Homebrew formula

## Timeline

| Event | Duration | Automatic? |
|-------|----------|------------|
| Work on feature branch | Variable | Manual |
| Create PR from feature branch | Immediate | Manual |
| Code review | Variable | Manual |
| Merge PR to main | Immediate | Manual |
| Release Please creates PR | ~30 seconds | ✅ Automatic |
| Review Release PR | Variable | Manual |
| Merge Release PR | Immediate | Manual |
| Release Please creates tag | ~30 seconds | ✅ Automatic |
| GitHub Actions builds | 3-5 minutes | ✅ Automatic |
| Homebrew formula updated | ~30 seconds | ✅ Automatic |
| **Total Development** | **Variable** | **Mostly manual** |
| **Total Release** | **5-10 minutes** | **Fully automatic** |

## First Release (v1.0.0)

For the initial release, you need to bootstrap Release Please:

### Option 1: Automatic Bootstrap (Recommended)

```bash
# Create feature branch for initial release
git checkout -b feature/initial-release

# Commit with conventional format
git commit -m "feat: initial production release"

# Push and create PR
git push origin feature/initial-release
gh pr create --base main --title "Initial Production Release" --body "v1.0.0 release"

# Merge PR to main
gh pr merge --squash

# Release Please will create a Release PR for v1.0.0
# Review and merge it to trigger the release
```

### Option 2: Manual Bootstrap

```bash
# Create initial tag manually
git tag -a v1.0.0 -m "Release v1.0.0 - Initial production release"
git push origin v1.0.0

# This triggers the manual release workflow
# Future releases will use Release Please
```

## Version Strategy

### Pre-1.0 (0.x.y)
- **0.x.0** - Minor features, might have breaking changes
- **0.x.y** - Patches and small fixes

### Post-1.0 (x.y.z)
- **x.0.0** - Breaking changes (avoid unless necessary)
- **1.x.0** - New features, backward compatible
- **1.y.z** - Bug fixes and patches

## Common Workflows

### Release a Patch Fix

```bash
# Create fix branch from main
git checkout main
git pull origin main
git checkout -b fix/database-timeout

# Fix the bug
git commit -m "fix: resolve database connection timeout"

# Push and create PR
git push origin fix/database-timeout
gh pr create \
  --base main \
  --head fix/database-timeout \
  --title "Fix database timeout" \
  --body "Resolves timeout issues"

# After code review, merge to main
gh pr merge --squash

# Wait for Release Please to create Release PR
# Merge Release PR → automatic release
```

### Release a New Feature

```bash
# Create feature branch from main
git checkout main
git pull origin main
git checkout -b feature/interactive-mode

# Add the feature
git commit -m "feat: add interactive mode for init command"
git commit -m "docs: document interactive mode"
git commit -m "test: add tests for interactive mode"

# Push and create PR
git push origin feature/interactive-mode
gh pr create \
  --base main \
  --head feature/interactive-mode \
  --title "Add interactive mode" \
  --body "Implements interactive initialization"

# After code review, merge to main
gh pr merge --squash

# Wait for Release Please to create Release PR
# Merge Release PR → automatic release
```

### Release Multiple Changes

**Option 1: Multiple PRs Merged Before Release**

```bash
# Merge multiple feature PRs to main
# Each PR triggers Release Please to update its internal state
# But Release Please only creates ONE Release PR with all changes

# PR 1: AWS provider
gh pr merge <PR-1> --squash

# PR 2: Retry logic
gh pr merge <PR-2> --squash

# PR 3: Bug fix
gh pr merge <PR-3> --squash

# After last PR, Release Please creates/updates Release PR
# The Release PR will include all changes:
# - feat: add AWS provider
# - feat: add retry logic for API calls
# - fix: handle expired credentials
# Version bump: 1.0.0 → 1.1.0 (minor, due to features)
```

**Option 2: Multiple Features in Single Branch**

```bash
# Create feature branch
git checkout -b feature/api-improvements

# Work on multiple related features
git commit -m "feat: add AWS provider"
git commit -m "feat: add retry logic for API calls"
git commit -m "fix: handle expired credentials"
git commit -m "docs: update provider documentation"

# Push and create single PR
git push origin feature/api-improvements
gh pr create --base main --title "API Improvements" --body "Multiple improvements"

# Merge to main (all commits included)
gh pr merge --squash

# Release Please creates Release PR with all changes
# Version bump: 1.0.0 → 1.1.0
```

### Emergency Hotfix

```bash
# Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b hotfix/security-vulnerability

# Fix the critical issue
git commit -m "fix: critical security vulnerability"

# Push and create PR
git push origin hotfix/security-vulnerability
gh pr create \
  --base main \
  --head hotfix/security-vulnerability \
  --title "HOTFIX: Critical Security Vulnerability" \
  --body "Addresses CVE-XXXX"

# After quick review, merge to main
gh pr merge --squash

# For emergency releases, you can manually tag instead of waiting
# Skip Release Please and trigger manual release workflow
git pull origin main
git tag -a v1.0.1 -m "Hotfix: Security vulnerability"
git push origin v1.0.1

# This triggers manual release workflow immediately
```

## Monitoring Releases

### Check Release Please Status

```bash
# View open PRs (look for Release Please PR)
gh pr list

# View specific Release PR
gh pr view <PR-NUMBER>

# Check workflow runs
gh run list --workflow=release-please.yml
```

### Check Release Status

```bash
# View latest release
gh release view

# List all releases
gh release list

# Watch release workflow
gh run watch
```

## Configuration Files

### .release-please-manifest.json

```json
{
  ".": "0.0.0"
}
```

Tracks current version. Updated automatically by Release Please.

### release-please-config.json

```json
{
  "release-type": "go",
  "packages": {
    ".": {
      "package-name": "stax",
      "changelog-path": "CHANGELOG.md",
      "extra-files": [
        {
          "type": "go",
          "path": "cmd/root.go",
          "glob": false
        }
      ]
    }
  }
}
```

Configures Release Please behavior:
- Updates `cmd/root.go` with new version
- Generates `CHANGELOG.md`
- Uses Go semantic versioning rules

## Troubleshooting

### Release Please PR Not Created

**Possible causes:**
1. No conventional commits since last release
2. Only `chore:`, `test:`, or `ci:` commits (these don't trigger releases)
3. Feature PR was merged but didn't have conventional commits

**Solution:**
```bash
# Check recent commits on main
git checkout main
git pull origin main
git log --oneline -10

# Look for commits with feat:, fix:, etc.
# If none, create a new feature branch with conventional commit
git checkout -b chore/trigger-release
git commit --allow-empty -m "chore: trigger release"
git push origin chore/trigger-release

# Create and merge PR
gh pr create --base main --title "Trigger Release" --body "Empty commit to trigger release"
gh pr merge --squash

# Release Please should now create PR
```

### Wrong Version Bump

**Possible cause:** Commit type doesn't match intended change

**Solution:**
```bash
# Before merging Release PR, check version
gh pr view <PR-NUMBER>

# If wrong, close the PR
gh pr close <PR-NUMBER>

# Fix commit messages in your branch
git rebase -i main
# Change commit message types
git push --force-with-lease

# Release Please will create new PR with correct version
```

### Multiple Release PRs

**Possible cause:** Multiple commits to main before first PR was merged

**Solution:**
```bash
# Close all Release PRs except the newest
gh pr list | grep "chore(main): release"

# Keep the highest version, close others
gh pr close <OLD-PR-NUMBER>
```

### Manual Release Needed

**When to use:**
- Emergency hotfixes
- Release Please workflow broken
- Need immediate release

**How to:**
```bash
git tag -a v1.0.1 -m "Manual release: description"
git push origin v1.0.1

# This triggers manual-release.yml workflow
# (which is kept for backward compatibility)
```

## Best Practices

### Commit Messages

✅ **Good:**
```bash
feat: add support for WordPress VIP provider
fix(auth): handle expired token gracefully
perf(db): optimize database connection pooling
docs: update installation guide
refactor: simplify provider initialization
```

❌ **Bad:**
```bash
updated stuff
WIP
asdf
fixed bug
added feature
```

### PR Descriptions

Include context in PR body:
```markdown
## Changes
- Add AWS provider implementation
- Update provider documentation
- Add integration tests

## Breaking Changes
None

## Migration Guide
N/A
```

### Version Planning

- Plan major versions carefully (breaking changes)
- Group related features into minor releases
- Release patches frequently for bug fixes
- Use pre-release versions (1.0.0-beta.1) for testing

## Migration from Manual Releases

If you've been creating tags manually:

1. **Let Release Please take over:**
   ```bash
   # Just start using conventional commits
   git commit -m "feat: new feature"
   git push origin main

   # Release Please will create PR
   ```

2. **Update documentation:**
   - Point contributors to this guide
   - Update CONTRIBUTING.md with commit conventions
   - Add commit message linter (optional)

3. **Archive old process:**
   - Keep manual-release.yml for emergencies
   - Document that it's deprecated
   - Train team on new process

## Additional Resources

- [Release Please Documentation](https://github.com/googleapis/release-please)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)

## Quick Reference

```bash
# Start a feature
git commit -m "feat: description"

# Fix a bug
git commit -m "fix: description"

# Breaking change
git commit -m "feat!: description"

# Check Release Please status
gh pr list | grep "chore(main): release"

# Merge Release PR
gh pr merge <PR-NUMBER> --merge

# View release
gh release view

# Emergency manual release
git tag -a v1.0.1 -m "Hotfix"
git push origin v1.0.1
```

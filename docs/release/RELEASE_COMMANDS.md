# Stax v1.0.0 - Release Quick Reference

## ✨ NEW: Automated Release

```bash
# Commit with conventional format and push to main
git commit -m "feat: your feature description"
git push origin main

# Wait for Release Please to create PR, then merge it
gh pr list | grep "chore(main): release"
gh pr merge <PR-NUMBER> --merge

# Release is automatically created!
```

## One-Command Manual Release (Deprecated)

```bash
# For emergencies or v1.0.0 only
git tag -a v1.0.0 -m "Release v1.0.0" && git push origin v1.0.0
```

## One-Command Verification

```bash
# Install and verify in one step
brew tap Firecrown-Media/stax && brew install stax && stax --version
```

## One-Command Rollback

```bash
# Delete tag and release
git push origin :refs/tags/v1.0.0 && gh release delete v1.0.0 --repo firecrown-media/stax --yes
```

---

## Essential Commands

### Setup GitHub Secret (One-Time)

```bash
# Create Personal Access Token at: https://github.com/settings/tokens/new
# Scopes needed: repo
# Then set the secret:
gh secret set HOMEBREW_TAP_TOKEN --repo firecrown-media/stax
```

### Create Release Tag

```bash
cd /Users/geoff/_projects/fc/stax
git checkout main
git pull origin main
git tag -a v1.0.0 -m "Release v1.0.0 - Initial production release"
git push origin v1.0.0
```

### Monitor Release

```bash
# Watch GitHub Actions in real-time
gh run watch --repo firecrown-media/stax

# Or open in browser
open https://github.com/firecrown-media/stax/actions
```

### Verify Release

```bash
# Check GitHub release
gh release view v1.0.0 --repo firecrown-media/stax

# Check Homebrew formula
cd /Users/geoff/_projects/fc/homebrew-stax && git pull && cat Formula/stax.rb | head -15

# Test installation
brew tap Firecrown-Media/stax
brew install stax
stax --version
```

### Test Local Build

```bash
# Run full test suite
make test

# Run security tests
make test-security

# Test GoReleaser locally (snapshot build)
goreleaser release --snapshot --clean --skip=publish

# Check artifacts
ls -lh dist/*.tar.gz
```

---

## Emergency Commands

### Delete Release

```bash
# Delete git tag (local and remote)
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Delete GitHub release
gh release delete v1.0.0 --repo firecrown-media/stax --yes

# Revert Homebrew formula
cd /Users/geoff/_projects/fc/homebrew-stax
git revert HEAD
git push origin main
```

### Create Patch Release

```bash
# After fixing bug in main branch
git tag -a v1.0.1 -m "Release v1.0.1 - Bug fixes"
git push origin v1.0.1
```

---

## Pre-Release Checklist

```bash
# Verify all tests pass
make test && make test-security

# Verify clean working directory
git status

# Verify on main branch
git branch --show-current

# Verify HOMEBREW_TAP_TOKEN exists
gh secret list --repo firecrown-media/stax | grep HOMEBREW_TAP_TOKEN
```

---

## Post-Release Tasks

```bash
# Announce on GitHub Discussions
gh pr create --repo firecrown-media/stax --title "v1.0.0 Released" --body "Initial production release"

# Update project documentation
# Add release notes to changelog
# Notify users/team
```

---

## Common Workflows

### Test Release Locally Before Pushing

```bash
# 1. Create tag locally (don't push yet)
git tag -a v1.0.0-test -m "Test release"

# 2. Test GoReleaser
GITHUB_TOKEN=$(gh auth token) goreleaser release --snapshot --clean --skip=publish

# 3. Verify artifacts
ls -lh dist/

# 4. Delete test tag
git tag -d v1.0.0-test
```

### Manual Formula Update (If Automated Fails)

```bash
# 1. Build release locally
goreleaser release --snapshot --clean --skip=publish

# 2. Copy formula
cp dist/homebrew/Formula/stax.rb /Users/geoff/_projects/fc/homebrew-stax/Formula/

# 3. Commit and push
cd /Users/geoff/_projects/fc/homebrew-stax
git add Formula/stax.rb
git commit -m "Update formula to v1.0.0"
git push origin main
```

---

## Troubleshooting One-Liners

```bash
# Check if secret is set
gh secret list --repo firecrown-media/stax

# View latest workflow run
gh run list --repo firecrown-media/stax --limit 1

# View workflow logs
gh run view --repo firecrown-media/stax --log

# List all tags
git tag -l

# List all releases
gh release list --repo firecrown-media/stax

# Test Homebrew tap
brew tap Firecrown-Media/stax --debug

# Reinstall from tap
brew uninstall stax && brew install stax

# Check formula info
brew info Firecrown-Media/stax/stax
```

---

## Time Estimates

| Task | Time |
|------|------|
| Create GitHub PAT | 2 min |
| Set secret | 1 min |
| Create & push tag | 30 sec |
| GitHub Actions runs | 3-5 min |
| Verify installation | 2 min |
| **Total** | **~10 min** |

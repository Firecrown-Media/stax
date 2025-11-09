# Feature Branch Workflow with Automated Releases

## Overview

All changes to Stax are developed on feature branches and merged to `main` via Pull Requests. This document describes the complete workflow from feature development to automated release.

## Branch Strategy

```
main (protected)
  ↑
  └── feature/your-feature
  └── fix/bug-fix
  └── hotfix/emergency-fix
```

**Rules:**
- `main` is protected - no direct commits
- All changes come via Pull Requests from feature branches
- Feature branches are short-lived (hours to days, not weeks)
- Branch names use prefixes: `feature/`, `fix/`, `hotfix/`, `docs/`, `refactor/`

## Complete Workflow

### 1. Start New Feature

```bash
# Always start from latest main
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/add-aws-provider

# Verify you're on the right branch
git branch --show-current
# Output: feature/add-aws-provider
```

### 2. Develop with Conventional Commits

```bash
# Make your changes
# Commit with conventional format

# Feature commits (minor version bump)
git commit -m "feat: add AWS provider implementation"
git commit -m "feat(auth): implement AWS credential handling"

# Bug fixes (patch version bump)
git commit -m "fix: handle AWS rate limiting correctly"

# Documentation (patch version bump)
git commit -m "docs: add AWS provider documentation"

# Tests (no version bump, but good practice)
git commit -m "test: add AWS provider integration tests"

# Refactoring (patch version bump)
git commit -m "refactor: simplify provider initialization"

# Breaking changes (major version bump)
git commit -m "feat!: redesign provider interface"
# or
git commit -m "feat: redesign provider interface

BREAKING CHANGE: Provider interface now requires Validate() method"
```

### 3. Push Feature Branch

```bash
# Push to GitHub
git push origin feature/add-aws-provider

# If branch doesn't exist remotely yet
git push -u origin feature/add-aws-provider
```

### 4. Create Pull Request

```bash
# Using GitHub CLI
gh pr create \
  --base main \
  --head feature/add-aws-provider \
  --title "Add AWS Provider Support" \
  --body "## Changes
- Complete AWS provider implementation
- Credential handling via AWS SDK
- Integration tests
- Documentation

## Breaking Changes
None

## Testing
- Unit tests pass
- Integration tests with localstack
- Manual testing on AWS account"

# Or use GitHub UI
# Navigate to your branch on GitHub and click "Create Pull Request"
```

### 5. Code Review Process

```bash
# Reviewers provide feedback via GitHub UI
# Make requested changes on your branch

# Continue committing with conventional format
git commit -m "fix: address code review feedback"
git push origin feature/add-aws-provider

# PR automatically updates with new commits
```

### 6. Merge to Main

**Before merging, ensure:**
- [ ] All tests pass (CI/CD runs automatically)
- [ ] Code review approved
- [ ] Conventional commits used
- [ ] Documentation updated

**Merge options:**

**Option A: Squash Merge (Recommended)**
```bash
gh pr merge --squash

# Benefits:
# - Clean git history (one commit per feature)
# - Easy to revert if needed
# - Clear changelog entries
```

**Option B: Regular Merge**
```bash
gh pr merge --merge

# Benefits:
# - Preserves all individual commits
# - Shows detailed development history
# - Release Please sees all commits
```

**Option C: Rebase and Merge**
```bash
gh pr merge --rebase

# Benefits:
# - Linear history
# - Preserves all commits
# - No merge commits
```

### 7. Release Please Takes Over

**Automatically after merge to main:**

1. **Release Please Analyzes Commits**
   - Scans all commits since last release
   - Determines version bump based on commit types
   - Generates changelog entries

2. **Release Please Creates/Updates Release PR**
   - If no Release PR exists, creates new one
   - If Release PR exists, updates it with new changes
   - PR title: `chore(main): release X.Y.Z`

3. **You Review Release PR**
   ```bash
   # Check for Release PR
   gh pr list | grep "chore(main): release"

   # View the Release PR
   gh pr view <PR-NUMBER>

   # Verify:
   # - Version bump is correct
   # - Changelog looks good
   # - All changes included
   ```

4. **Merge Release PR**
   ```bash
   # Merge the Release PR (use merge, not squash)
   gh pr merge <PR-NUMBER> --merge

   # DO NOT squash Release PRs - use regular merge
   ```

5. **Automatic Release**
   - Release Please creates git tag
   - GitHub release created
   - GoReleaser builds binaries
   - Homebrew formula updated
   - All automatic!

## Multiple Features Workflow

### Scenario: Working on Multiple Features

**Option 1: Separate Branches (Recommended)**

```bash
# Feature 1: AWS Provider
git checkout -b feature/aws-provider
git commit -m "feat: add AWS provider"
git push origin feature/aws-provider
gh pr create --base main --title "AWS Provider"

# Feature 2: Retry Logic (independent)
git checkout main
git checkout -b feature/retry-logic
git commit -m "feat: add retry logic for API calls"
git push origin feature/retry-logic
gh pr create --base main --title "Retry Logic"

# Merge both PRs when ready
# Release Please will create ONE Release PR with both features
```

**Option 2: Dependent Features**

```bash
# Feature 1: Base implementation
git checkout -b feature/provider-base
git commit -m "feat: add provider base class"
git push origin feature/provider-base
gh pr create --base main --title "Provider Base"

# Wait for PR to be merged to main

# Feature 2: Build on feature 1
git checkout main
git pull origin main
git checkout -b feature/aws-provider
git commit -m "feat: implement AWS provider using base"
git push origin feature/aws-provider
gh pr create --base main --title "AWS Provider"
```

## Hotfix Workflow

For critical bugs in production:

```bash
# 1. Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b hotfix/critical-bug

# 2. Fix the bug
git commit -m "fix: resolve critical security vulnerability"

# 3. Push and create PR
git push origin hotfix/critical-bug
gh pr create \
  --base main \
  --head hotfix/critical-bug \
  --title "HOTFIX: Critical Security Fix" \
  --body "Addresses critical vulnerability"

# 4. Fast-track review and merge
gh pr merge --squash

# 5. Option A: Wait for Release Please (normal process)
# Release Please creates Release PR
# Merge Release PR → automatic release

# 5. Option B: Emergency manual release (skip Release Please)
git checkout main
git pull origin main
git tag -a v1.0.1 -m "Hotfix: Security vulnerability"
git push origin v1.0.1
# Triggers manual release workflow immediately
```

## Common Scenarios

### Updating Feature Branch with Main

```bash
# Your feature branch is behind main
git checkout feature/your-feature
git fetch origin

# Option 1: Merge main into feature (preserves history)
git merge origin/main

# Option 2: Rebase on main (cleaner history)
git rebase origin/main

# Push updated branch
git push origin feature/your-feature --force-with-lease
```

### Fixing Mistakes in Feature Branch

```bash
# Forgot to use conventional commit format
git checkout feature/your-feature

# Amend last commit message
git commit --amend -m "feat: add new feature"
git push --force-with-lease origin feature/your-feature

# Rewrite multiple commits (use interactive rebase)
git rebase -i origin/main
# Change commit messages in editor
git push --force-with-lease origin feature/your-feature
```

### Combining Multiple PRs into One Release

```bash
# Merge PR 1
gh pr merge <PR-1> --squash
# Release Please updates its internal state

# Merge PR 2
gh pr merge <PR-2> --squash
# Release Please updates Release PR

# Merge PR 3
gh pr merge <PR-3> --squash
# Release Please updates Release PR again

# The Release PR now includes all 3 changes
# Merge Release PR → single release with all features
```

## Best Practices

### Branch Naming

```bash
# Good
feature/add-aws-provider
fix/database-connection-timeout
hotfix/security-vulnerability
docs/update-installation-guide
refactor/simplify-provider-init

# Avoid
new-feature
bugfix
updates
my-changes
```

### Commit Messages

```bash
# Good - clear, conventional, descriptive
feat: add AWS provider with S3 and EC2 support
fix(auth): handle expired credentials gracefully
docs: add AWS provider configuration guide
test: add integration tests for AWS provider

# Avoid - unclear, non-conventional
added stuff
fix bug
update
WIP
changes
```

### PR Descriptions

**Good PR:**
```markdown
## Summary
Implements AWS provider support with S3 and EC2 integration.

## Changes
- AWS provider implementation
- Credential handling via AWS SDK
- Integration tests with localstack
- Documentation

## Breaking Changes
None

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Documentation reviewed

## Related Issues
Closes #123
```

**Avoid:**
```markdown
updates

WIP

see commits
```

### When to Merge

✅ **Merge when:**
- All tests pass
- Code review approved
- Documentation updated
- Conventional commits used
- No conflicts with main

❌ **Don't merge when:**
- Tests failing
- No code review
- Breaking changes not documented
- Non-conventional commit messages
- Merge conflicts exist

## Troubleshooting

### PR Won't Merge

```bash
# Update with latest main
git checkout feature/your-feature
git fetch origin
git merge origin/main

# Resolve conflicts if any
git add .
git commit -m "chore: resolve merge conflicts"
git push origin feature/your-feature
```

### Release Please Not Triggering

```bash
# Check recent commits on main
git checkout main
git pull origin main
git log --oneline -5

# Verify conventional commits exist
# If not, your commits may not have triggered Release Please

# Check GitHub Actions
gh run list --workflow=release-please.yml
```

### Wrong Version Bump

```bash
# If Release Please calculated wrong version:
# 1. Close the Release PR (don't merge it)
gh pr close <RELEASE-PR-NUMBER>

# 2. Fix commit messages in your feature branch
git checkout feature/your-feature
git rebase -i origin/main
# Fix commit message types in editor

# 3. Force push
git push --force-with-lease origin feature/your-feature

# 4. Merge corrected PR
gh pr merge --squash

# 5. Release Please will create new Release PR with correct version
```

## Quick Reference

```bash
# Start new feature
git checkout main && git pull && git checkout -b feature/name

# Commit with convention
git commit -m "feat: description"

# Push and create PR
git push -u origin feature/name
gh pr create --base main --title "Title"

# Merge PR
gh pr merge --squash

# Check for Release PR
gh pr list | grep "chore(main): release"

# Merge Release PR
gh pr merge <PR-NUMBER> --merge

# Verify release
gh release view
```

## Summary

**The golden rule:** All changes flow through feature branches → PRs → main → Release Please → automated release.

Never commit directly to `main`. Always use feature branches and conventional commits. Release Please handles the rest automatically!

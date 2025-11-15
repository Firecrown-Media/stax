# Contributing to Stax CLI

Welcome to the Stax CLI project! This guide will help you understand our development workflow, from creating a feature branch to seeing your changes automatically deployed via Homebrew.

## Table of Contents

1. [Development Workflow Overview](#development-workflow-overview)
2. [Feature Branch Promotion Strategy](#feature-branch-promotion-strategy)
3. [Conventional Commits](#conventional-commits)
4. [Release-Please Workflow](#release-please-workflow)
5. [GoReleaser Integration](#goreleaser-integration)
6. [Step-by-Step: Implementing a Feature](#step-by-step-implementing-a-feature)
7. [Homebrew Update Process](#homebrew-update-process)
8. [Testing Before Release](#testing-before-release)
9. [Troubleshooting](#troubleshooting)
10. [Best Practices](#best-practices)

## Development Workflow Overview

The Stax CLI uses an automated release workflow that connects feature development to production deployment:

```
Feature Branch → PR → Main → Release-Please → Release PR → Merge → Auto-Release → Homebrew Update
```

### Key Components

- **Git Branching**: Feature branches from `main`
- **Commit Format**: Conventional Commits for semantic versioning
- **PR Process**: Squash merges to maintain clean history
- **Release Automation**: Release-Please for version management
- **Build Automation**: GoReleaser for multi-platform builds
- **Distribution**: Automatic Homebrew formula updates

### Repository Information

- **GitHub Repository**: [Firecrown-Media/stax](https://github.com/Firecrown-Media/stax)
- **Main Branch**: `main`
- **Homebrew Tap**: [Firecrown-Media/homebrew-stax](https://github.com/Firecrown-Media/homebrew-stax)
- **Release Workflow**: `.github/workflows/release-please.yml`
- **Build Config**: `.goreleaser.yml`

## Feature Branch Promotion Strategy

### Complete Feature Lifecycle

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. CREATE FEATURE BRANCH                                        │
│    git checkout -b feature/my-feature                           │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 2. IMPLEMENT WITH CONVENTIONAL COMMITS                          │
│    feat(init): add new functionality                            │
│    fix(database): resolve connection issue                      │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 3. CREATE PULL REQUEST                                          │
│    gh pr create --title "feat: Add Feature" --body "..."        │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 4. CODE REVIEW & APPROVAL                                       │
│    - Automated tests pass                                       │
│    - Code review completed                                      │
│    - Changes requested/approved                                 │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 5. SQUASH MERGE TO MAIN                                         │
│    - One commit with conventional format                        │
│    - Clean git history maintained                               │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 6. RELEASE-PLEASE DETECTS CHANGES                               │
│    - Analyzes commit messages                                   │
│    - Calculates version bump                                    │
│    - Updates CHANGELOG.md                                       │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 7. CREATES/UPDATES RELEASE PR                                   │
│    - PR title: "chore(main): release X.Y.Z"                     │
│    - Contains version bumps + changelog                         │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 8. MERGE RELEASE PR                                             │
│    - Triggers release workflow                                  │
│    - Creates GitHub release                                     │
│    - Generates git tag                                          │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 9. GORELEASER BUILDS & PUBLISHES                                │
│    - Compiles for Darwin (amd64, arm64)                         │
│    - Compiles for Linux (amd64, arm64)                          │
│    - Uploads release artifacts                                  │
│    - Updates Homebrew formula                                   │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│ 10. USERS GET UPDATE                                            │
│     brew upgrade firecrown-media/stax/stax                      │
└─────────────────────────────────────────────────────────────────┘
```

### Branch Naming Conventions

Follow these naming patterns for feature branches:

- **Features**: `feature/descriptive-name`
  - Example: `feature/media-proxy-implementation`
  - Example: `feature/enhanced-doctor-diagnostics`

- **Bug Fixes**: `fix/descriptive-name`
  - Example: `fix/credential-storage-issue`
  - Example: `fix/multisite-detection`

- **Documentation**: `docs/descriptive-name`
  - Example: `docs/wpengine-credentials`
  - Example: `docs/media-proxy-guide`

- **Refactoring**: `refactor/descriptive-name`
  - Example: `refactor/database-package`

- **Performance**: `perf/descriptive-name`
  - Example: `perf/file-sync-optimization`

### Real-World Examples from Stax

Here are actual PRs from the Stax project:

```bash
# PR #53 - Feature implementation
feature/phase6-complete-init-integration
→ "Phase 6: Complete Init Integration"
→ Merged to main
→ Included in release 2.3.0

# PR #52 - Documentation
feature/media-proxy-docs
→ "docs: Add Comprehensive Media Proxy Documentation"
→ Merged to main
→ Included in release 2.3.0

# PR #51 - Documentation enhancement
docs/wpengine-credentials-enhancement
→ "docs: Enhance WPEngine Credentials Documentation"
→ Merged to main
→ Included in release 2.3.0

# PR #47 - Major feature
feature/complete-init-implementation
→ "Phase 3: Complete stax init Implementation"
→ Merged to main
→ Included in release 2.2.0
```

## Conventional Commits

Stax uses [Conventional Commits](https://www.conventionalcommits.org/) for semantic versioning and automated changelog generation.

### Commit Message Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Commit Types

| Type | Description | Version Bump | Example |
|------|-------------|--------------|---------|
| `feat` | New feature | **Minor** (0.1.0) | `feat(init): add interactive mode` |
| `fix` | Bug fix | **Patch** (0.0.1) | `fix(database): resolve connection timeout` |
| `docs` | Documentation only | None | `docs(readme): update installation steps` |
| `refactor` | Code refactoring | None | `refactor(ui): simplify spinner logic` |
| `test` | Adding tests | None | `test(init): add integration tests` |
| `chore` | Build/tooling | None | `chore: update dependencies` |
| `perf` | Performance improvements | **Patch** (0.0.1) | `perf(sync): optimize file transfer` |
| `style` | Code formatting | None | `style: format with gofmt` |
| `ci` | CI/CD changes | None | `ci: update release workflow` |

### Breaking Changes

Breaking changes trigger a **major** version bump (1.0.0):

```
feat(api)!: change configuration format

BREAKING CHANGE: Configuration format changed from YAML to TOML.
Users must migrate their .stax.yml files to .stax.toml format.
```

### Scope Guidelines

Use scopes to identify the affected component:

- `init` - Initialization commands
- `database` - Database operations
- `media` - Media proxy functionality
- `doctor` - Diagnostics
- `ui` - User interface components
- `config` - Configuration management
- `wpengine` - WPEngine provider
- `ddev` - DDEV integration

### Real Commit Examples

From the Stax repository:

```bash
# Feature with scope (minor bump)
feat(init): integrate file pull into init workflow

Replaces TODO placeholder with actual file pull integration.
The pullFiles function now calls the existing runFilesPull functionality.

Part of Phase 6: Complete Init Integration

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

# Bug fix (patch bump)
fix(ci): update release-please workflow to v4 configuration

Updates the release-please action from v3 to v4 to support
the latest configuration format and features.

# Documentation (no version bump)
docs(media): add comprehensive media proxy documentation

Adds complete documentation for media proxy configuration,
including setup instructions, troubleshooting, and examples.

# Multiple features in one release
feat(doctor): enhance diagnostics and add global WPEngine discovery

- Adds global WPEngine instance discovery
- Improves diagnostic output formatting
- Adds connectivity checks
```

### Writing Good Commit Messages

**DO:**
- Use present tense ("add feature" not "added feature")
- Be concise in the subject line (50-72 characters)
- Provide context in the body
- Reference issue numbers when applicable
- Explain "why" not just "what"

**DON'T:**
- Use vague descriptions ("update stuff", "fix bug")
- Include multiple unrelated changes
- Forget the scope when applicable
- Mix feature and fix changes

**Good Examples:**

```bash
feat(database): add connection pooling support

Implements connection pooling to improve database performance
for large imports. Pools are configured based on available
system resources.

Fixes #42
```

```bash
fix(init): prevent duplicate DDEV config generation

The init command was creating duplicate config entries when
run multiple times in the same directory. Now checks for
existing config before generating new entries.

Resolves #38
```

**Bad Examples:**

```bash
# Too vague
fix: bug fixes

# Missing type
add new feature to init

# Multiple unrelated changes
feat: add media proxy, fix database bug, update docs
```

## Release-Please Workflow

Release-Please automates the entire release process by analyzing commit messages and managing versions.

### How It Works

1. **Monitors Main Branch**
   - Workflow triggers on every push to `main`
   - Analyzes new commits since last release

2. **Analyzes Commits**
   - Parses conventional commit messages
   - Determines version bump type (major/minor/patch)
   - Generates changelog entries

3. **Creates/Updates Release PR**
   - Opens PR titled "chore(main): release X.Y.Z"
   - Updates `CHANGELOG.md` with new entries
   - Bumps version in relevant files

4. **Waits for Merge**
   - Release PR stays open for review
   - Accumulates changes from subsequent merges
   - Updates automatically with new commits

5. **Triggers Release on Merge**
   - Creates GitHub release
   - Generates git tag
   - Triggers GoReleaser workflow

### Version Bump Logic

Release-Please determines version bumps based on commit types:

```
Current version: 2.2.0

Commits merged:
- docs(readme): update installation → No bump
- fix(database): resolve timeout → Patch: 2.2.1
- feat(media): add proxy support → Minor: 2.3.0
- feat!: breaking change → Major: 3.0.0
```

**Priority Order:**
1. Breaking changes → Major bump (1.0.0)
2. Features → Minor bump (0.1.0)
3. Fixes → Patch bump (0.0.1)
4. Other → No bump

### Release PR Example

When you merge features to main, Release-Please creates a PR like this:

**Title:** `chore(main): release 2.3.0`

**Contents:**
```markdown
## [2.3.0](https://github.com/Firecrown-Media/stax/compare/v2.2.0...v2.3.0) (2025-11-15)

### Features

* **init:** integrate file pull into init workflow ([abc1234](link))
* **media:** implement media proxy configuration ([def5678](link))

### Bug Fixes

* **database:** resolve connection timeout issue ([ghi9012](link))

### Documentation

* **readme:** update installation instructions ([jkl3456](link))
```

### Workflow File

The workflow is defined in `.github/workflows/release-please.yml`:

```yaml
name: Release Please

on:
  push:
    branches:
      - main

permissions:
  contents: write
  pull-requests: write

jobs:
  release-please:
    runs-on: ubuntu-latest
    steps:
      - name: Release Please
        uses: googleapis/release-please-action@v4
        id: release
        with:
          release-type: go

      # When release is created, run tests and GoReleaser
      - name: Checkout
        if: ${{ steps.release.outputs.release_created }}
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        if: ${{ steps.release.outputs.release_created }}
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run tests
        if: ${{ steps.release.outputs.release_created }}
        run: make test

      - name: Run GoReleaser
        if: ${{ steps.release.outputs.release_created }}
        uses: goreleaser/goreleaser-action@v6
        with:
          version: '~> v2'
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

### Accumulating Changes

Release-Please intelligently accumulates changes:

```bash
# Day 1: Merge feature A
git merge feature/a  # feat(init): add feature A
→ Release PR created: "chore(main): release 2.3.0"

# Day 2: Merge feature B
git merge feature/b  # feat(media): add feature B
→ Release PR updated: "chore(main): release 2.3.0"
→ Both features now in changelog

# Day 3: Merge bug fix
git merge fix/c      # fix(database): fix issue C
→ Release PR updated: "chore(main): release 2.3.0"
→ All changes in changelog

# Day 4: Merge release PR
→ Release 2.3.0 published with all changes
```

## GoReleaser Integration

GoReleaser handles building, packaging, and distributing Stax across multiple platforms.

### What GoReleaser Does

1. **Compiles Binaries**
   - Darwin (macOS): amd64, arm64
   - Linux: amd64, arm64
   - Embeds version information

2. **Creates Archives**
   - `.tar.gz` for Unix systems
   - Includes README, LICENSE, docs

3. **Generates Checksums**
   - SHA256 checksums for verification
   - `checksums.txt` file

4. **Uploads to GitHub**
   - Attaches artifacts to release
   - Makes binaries downloadable

5. **Updates Homebrew Formula**
   - Commits to `homebrew-stax` repository
   - Updates version and SHA256
   - Users get automatic updates

### Configuration

The `.goreleaser.yml` file configures the build:

```yaml
version: 2

# Run before building
before:
  hooks:
    - go mod tidy
    - make test

# Build configuration
builds:
  - id: stax-darwin
    env:
      - CGO_ENABLED=0
    goos:
      - darwin
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X github.com/firecrown-media/stax/cmd.Version={{.Version}}
      - -X github.com/firecrown-media/stax/cmd.GitCommit={{.Commit}}
      - -X github.com/firecrown-media/stax/cmd.BuildDate={{.Date}}
    binary: stax

  - id: stax-linux
    # Similar configuration for Linux

# Archive configuration
archives:
  - id: stax
    name_template: >-
      {{ .ProjectName }}_
      {{- .Version }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else }}{{ .Arch }}{{ end }}
    files:
      - README.md
      - LICENSE
      - docs/**/*

# Homebrew formula configuration
brews:
  - name: stax
    homepage: https://github.com/firecrown-media/stax
    description: Powerful CLI tool for WordPress development workflows
    license: MIT
    repository:
      owner: Firecrown-Media
      name: homebrew-stax
      branch: main
      token: "{{ .Env.HOMEBREW_TAP_TOKEN }}"
    install: |
      bin.install "stax"
    test: |
      system "#{bin}/stax", "--version"
```

### Build Artifacts

After a successful release, GoReleaser produces:

```
dist/
├── stax_2.3.0_Darwin_x86_64.tar.gz
├── stax_2.3.0_Darwin_arm64.tar.gz
├── stax_2.3.0_Linux_x86_64.tar.gz
├── stax_2.3.0_Linux_arm64.tar.gz
├── checksums.txt
└── ...
```

### Version Information

The build embeds version information using ldflags:

```go
// In cmd/version.go
var (
    Version   string // Set by GoReleaser: "2.3.0"
    GitCommit string // Set by GoReleaser: "abc1234"
    BuildDate string // Set by GoReleaser: "2025-11-15"
)
```

Users can verify their installation:

```bash
stax --version
# Output: stax version 2.3.0 (abc1234) built on 2025-11-15
```

## Step-by-Step: Implementing a Feature

Let's walk through a complete example of implementing a feature and getting it released.

### Example: Adding a New Command

Suppose you want to add a `stax config:show` command to display current configuration.

#### Step 1: Create Feature Branch

```bash
# Make sure you're on main and up to date
cd /Users/geoff/_projects/fc/stax
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/config-show-command

# Verify you're on the new branch
git branch
# * feature/config-show-command
#   main
```

#### Step 2: Implement the Feature

```bash
# Create the new command file
# Edit cmd/config_show.go
# Implement the command logic
# Add tests

# Example files you might modify:
# - cmd/config_show.go (new file)
# - cmd/config_show_test.go (new file)
# - README.md (update with new command)
```

#### Step 3: Test Locally

```bash
# Run tests
make test

# Build and test manually
go build -o stax-dev
./stax-dev config:show

# Run linting
make lint

# Format code
gofmt -w .
```

#### Step 4: Commit with Conventional Format

```bash
# Stage your changes
git add cmd/config_show.go cmd/config_show_test.go README.md

# Commit with conventional format
git commit -m "feat(config): add config:show command

Implements a new command to display the current Stax configuration
in a user-friendly format. Supports JSON output via --json flag.

Features:
- Shows all configuration values
- Masks sensitive credentials
- Supports machine-readable JSON output
- Includes validation of config file

Fixes #65

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

#### Step 5: Push and Create PR

```bash
# Push feature branch
git push origin feature/config-show-command

# Create PR using GitHub CLI
gh pr create \
  --title "feat(config): Add config:show command" \
  --body "## Summary

Implements a new \`stax config:show\` command to display the current configuration.

## Changes

- Add \`config:show\` command implementation
- Add unit tests for config display
- Update README with new command documentation
- Add JSON output support via \`--json\` flag

## Testing

- ✅ Unit tests pass
- ✅ Manual testing completed
- ✅ Validates configuration file
- ✅ Properly masks sensitive data

## Related Issues

Fixes #65

## Screenshots

\`\`\`
$ stax config:show
Stax Configuration
==================
Project: my-wordpress-site
Environment: production
WPEngine Install: mysite
...
\`\`\`

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

#### Step 6: Code Review Process

```bash
# Monitor PR status
gh pr view

# Respond to review comments
# Make additional commits if needed
git add .
git commit -m "fix(config): address review feedback"
git push origin feature/config-show-command

# PR automatically updates
```

#### Step 7: Merge PR

Once approved, merge via GitHub:

```bash
# Via GitHub CLI (maintainers)
gh pr merge --squash --auto

# Or via GitHub web interface
# Click "Squash and merge"
```

The squash merge creates a single commit on `main`:

```
feat(config): Add config:show command (#66)

* Implements config:show command
* Adds JSON output support
* Includes tests and documentation

Fixes #65
```

#### Step 8: Wait for Release-Please

Within minutes, Release-Please will:

1. Detect the new `feat:` commit
2. Calculate version bump (2.2.0 → 2.3.0)
3. Create/update release PR

```bash
# Check for release PR
gh pr list

# Example output:
# #67  chore(main): release 2.3.0  release-please--branches--main
```

#### Step 9: Review Release PR

```bash
# View the release PR
gh pr view 67

# The PR will show:
# - Version bump: 2.2.0 → 2.3.0
# - Updated CHANGELOG.md
# - Your feature in the changelog

# Example CHANGELOG entry:
## [2.3.0](https://github.com/Firecrown-Media/stax/compare/v2.2.0...v2.3.0) (2025-11-15)

### Features

* **config:** add config:show command (#66) ([abc1234](link))
```

#### Step 10: Merge Release PR

```bash
# Merge the release PR (maintainers only)
gh pr merge 67 --squash

# This triggers:
# 1. Release-Please creates GitHub release
# 2. GoReleaser builds binaries
# 3. Homebrew formula updates automatically
```

#### Step 11: Verify Release

```bash
# Check the release was created
gh release view v2.3.0

# Verify artifacts
gh release view v2.3.0 --json assets

# Check Homebrew tap update
# Visit: https://github.com/Firecrown-Media/homebrew-stax
```

#### Step 12: Users Get the Update

Users can now upgrade:

```bash
# Update Homebrew
brew update

# Upgrade Stax
brew upgrade firecrown-media/stax/stax

# Verify version
stax --version
# stax version 2.3.0 (abc1234) built on 2025-11-15

# Use the new command
stax config:show
```

### Timeline Summary

```
00:00 - Create feature branch
01:00 - Implement feature + tests
02:00 - Create PR
02:30 - Code review begins
04:00 - Revisions made
08:00 - PR approved and merged
08:05 - Release-Please creates release PR #67
       (includes your feature)
24:00 - Maintainer merges release PR
24:05 - GitHub release v2.3.0 created
24:10 - GoReleaser builds complete
24:15 - Homebrew formula updated
24:20 - Users can upgrade!
```

## Homebrew Update Process

The Homebrew update is fully automated via GoReleaser and the `HOMEBREW_TAP_TOKEN`.

### How It Works

1. **GoReleaser Builds Binaries**
   - Compiles for Darwin and Linux
   - Generates SHA256 checksums

2. **Updates Formula**
   - Uses `HOMEBREW_TAP_TOKEN` for authentication
   - Commits to `Firecrown-Media/homebrew-stax`
   - Updates version and checksums

3. **Formula Changes**

Before (version 2.2.0):
```ruby
class Stax < Formula
  desc "Powerful CLI tool for WordPress development workflows"
  homepage "https://github.com/firecrown-media/stax"
  version "2.2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/Firecrown-Media/stax/releases/download/v2.2.0/stax_2.2.0_Darwin_x86_64.tar.gz"
      sha256 "abc123..."
    end
    if Hardware::CPU.arm?
      url "https://github.com/Firecrown-Media/stax/releases/download/v2.2.0/stax_2.2.0_Darwin_arm64.tar.gz"
      sha256 "def456..."
    end
  end

  def install
    bin.install "stax"
  end
end
```

After (version 2.3.0):
```ruby
class Stax < Formula
  desc "Powerful CLI tool for WordPress development workflows"
  homepage "https://github.com/firecrown-media/stax"
  version "2.3.0"  # ← Updated

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/Firecrown-Media/stax/releases/download/v2.3.0/stax_2.3.0_Darwin_x86_64.tar.gz"  # ← Updated
      sha256 "xyz789..."  # ← Updated
    end
    if Hardware::CPU.arm?
      url "https://github.com/Firecrown-Media/stax/releases/download/v2.3.0/stax_2.3.0_Darwin_arm64.tar.gz"  # ← Updated
      sha256 "uvw012..."  # ← Updated
    end
  end

  def install
    bin.install "stax"
  end
end
```

### HOMEBREW_TAP_TOKEN

The `HOMEBREW_TAP_TOKEN` is a GitHub Personal Access Token (PAT) with permissions to write to the `homebrew-stax` repository.

**Setup (for maintainers):**

1. Generate PAT in GitHub Settings
2. Grant permissions: `repo` scope
3. Add to repository secrets: `HOMEBREW_TAP_TOKEN`
4. GoReleaser uses it automatically

**Security:**
- Token is stored as GitHub secret
- Never committed to repository
- Only accessible during GitHub Actions runs
- Can be rotated if compromised

### User Upgrade Process

Users experience seamless updates:

```bash
# Check current version
stax --version
# stax version 2.2.0

# Update Homebrew package index
brew update

# Upgrade Stax
brew upgrade firecrown-media/stax/stax

# Verify new version
stax --version
# stax version 2.3.0

# Use new features
stax config:show
```

### Verification Steps

After a release, verify the Homebrew update:

```bash
# Check the homebrew-stax repository
gh repo view Firecrown-Media/homebrew-stax

# View recent commits
gh repo view Firecrown-Media/homebrew-stax \
  --web --branch main

# Check formula file
curl -L https://raw.githubusercontent.com/Firecrown-Media/homebrew-stax/main/Formula/stax.rb

# Test installation in clean environment
brew uninstall stax
brew install firecrown-media/stax/stax
stax --version
```

## Testing Before Release

Thorough testing prevents issues in production releases.

### Local Testing Checklist

Before creating a PR:

```bash
# 1. Run all tests
make test

# 2. Run security tests
make test-security

# 3. Build locally
go build -o stax-dev

# 4. Test the binary
./stax-dev --version
./stax-dev [your-new-command]

# 5. Run integration tests (if available)
make test-integration

# 6. Check for linting issues
make lint

# 7. Format code
gofmt -w .

# 8. Verify no unintended changes
git status
git diff
```

### Manual Testing

Test common user workflows:

```bash
# Test init workflow
cd /tmp/test-project
stax-dev init

# Test database operations
stax-dev db:pull

# Test configuration
stax-dev config:show

# Test doctor diagnostics
stax-dev doctor

# Test help output
stax-dev --help
stax-dev [command] --help
```

### Integration Testing

For significant changes:

```bash
# Create test WordPress site
mkdir -p /tmp/stax-integration-test
cd /tmp/stax-integration-test

# Initialize with real WPEngine credentials
export WPENGINE_USER_ID="your-user-id"
export WPENGINE_PASSWORD="your-password"

# Run complete workflow
stax-dev init --name testsite --install myinstall
stax-dev db:pull
stax-dev files:pull
stax-dev start

# Verify site is accessible
open http://testsite.ddev.site

# Clean up
stax-dev stop
cd -
rm -rf /tmp/stax-integration-test
```

### PR Testing

Automated tests run on every PR:

```yaml
# .github/workflows/test.yml
name: Test

on:
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run tests
        run: make test

      - name: Run security tests
        run: make test-security

      - name: Check formatting
        run: |
          gofmt -d .
          test -z "$(gofmt -d .)"
```

### Beta Testing

For major features, consider beta testing:

```bash
# Create pre-release
gh release create v2.3.0-beta.1 \
  --title "v2.3.0 Beta 1" \
  --notes "Beta release for testing new features" \
  --prerelease

# Build and upload beta binaries
goreleaser release --snapshot --clean

# Ask beta testers to test
# Collect feedback
# Fix issues
# Create next beta or final release
```

## Troubleshooting

Common issues and solutions:

### Release-Please Not Creating PR

**Symptom:** You merged a feature but no release PR was created.

**Causes:**
1. Commit message doesn't follow conventional format
2. Only `docs:`, `test:`, or `chore:` commits (no version bump)
3. Workflow failed

**Solutions:**

```bash
# Check workflow runs
gh run list --workflow=release-please.yml

# View latest run
gh run view [run-id]

# Check commit messages on main
git log --oneline -10

# If needed, create a trigger commit
git commit --allow-empty -m "chore: trigger release-please"
git push origin main
```

### Version Not Bumping Correctly

**Symptom:** Release version is wrong (e.g., patch instead of minor).

**Cause:** Commit types don't match expectations.

**Solution:**

```bash
# Review commits included in release
gh pr view [release-pr-number]

# Check commit history
git log v2.2.0..HEAD --oneline

# Verify commit messages use correct types
# - feat: for new features (minor bump)
# - fix: for bug fixes (patch bump)
# - BREAKING CHANGE: for breaking changes (major bump)
```

### Homebrew Formula Not Updating

**Symptom:** Release created but Homebrew formula wasn't updated.

**Causes:**
1. `HOMEBREW_TAP_TOKEN` missing or invalid
2. GoReleaser configuration error
3. Network/permissions issue

**Solutions:**

```bash
# Check GoReleaser logs in GitHub Actions
gh run view [run-id] --log

# Look for Homebrew-related errors
# Example: "token: invalid authentication credentials"

# Verify token is set (maintainers only)
gh secret list

# If needed, update token
gh secret set HOMEBREW_TAP_TOKEN

# Manually verify Homebrew tap
gh repo view Firecrown-Media/homebrew-stax
```

### Build Failures

**Symptom:** GoReleaser build fails.

**Common Causes:**
1. Test failures
2. Missing dependencies
3. Build errors
4. Incorrect `.goreleaser.yml` syntax

**Solutions:**

```bash
# Run tests locally
make test

# Build locally
go build

# Test GoReleaser config
goreleaser check

# Run GoReleaser in snapshot mode (no release)
goreleaser release --snapshot --clean

# Check build output
ls -la dist/
```

### Permission Denied Errors

**Symptom:** Can't merge release PR or push to main.

**Solution:**

```bash
# Verify you're a maintainer
gh repo view Firecrown-Media/stax

# Check branch protection rules
gh api repos/Firecrown-Media/stax/branches/main/protection

# Contact repository owner for access
```

### Duplicate Release PRs

**Symptom:** Multiple release PRs exist.

**Cause:** Previous release PR wasn't merged or closed properly.

**Solution:**

```bash
# List all release PRs
gh pr list --label "autorelease: pending"

# Close old release PRs
gh pr close [old-pr-number]

# Merge current release PR
gh pr merge [current-pr-number] --squash
```

## Best Practices

### When to Create Feature Branches

**Always:**
- New features
- Bug fixes
- Documentation updates
- Refactoring

**Never:**
- Direct commits to `main`
- Hotfixes without PR (except emergencies)

### Writing Good Commit Messages

**Template:**

```
<type>(<scope>): <short summary>
<blank line>
<detailed description>
<blank line>
<footer with issue references>
```

**Examples:**

```bash
# Good
feat(init): add interactive mode with guided prompts

Implements an interactive mode for the init command that guides
users through project setup with contextual help and validation.

Features:
- Step-by-step prompts for all configuration options
- Input validation with helpful error messages
- Ability to skip optional steps
- Preview of configuration before writing

Fixes #45

# Bad
update init

# Better
feat(init): add interactive mode
```

### PR Description Guidelines

**Good PR Description:**

```markdown
## Summary
Brief overview of the changes (2-3 sentences).

## Changes
- Bullet point list of specific changes
- Be clear and concise
- Group related changes

## Testing
- ✅ How you tested the changes
- ✅ What scenarios were covered
- ✅ Any manual testing performed

## Related Issues
Fixes #123
Related to #456

## Screenshots/Examples
Include command output or UI changes

## Breaking Changes
List any breaking changes and migration steps

## Checklist
- [x] Tests added/updated
- [x] Documentation updated
- [x] Follows conventional commits
- [x] No sensitive data exposed
```

### Code Review Checklist

**For Authors:**
- [ ] All tests pass locally
- [ ] Code is formatted (`gofmt`)
- [ ] Documentation updated
- [ ] Conventional commit format used
- [ ] No debug code or console.logs
- [ ] Error handling implemented
- [ ] No sensitive data in commits

**For Reviewers:**
- [ ] Code follows project conventions
- [ ] Tests are adequate
- [ ] Documentation is clear
- [ ] No security concerns
- [ ] Performance is acceptable
- [ ] Error handling is robust
- [ ] Commit messages are conventional

### Commit Frequency

**During Development:**
```bash
# Commit frequently with WIP prefix
git commit -m "WIP: add config show logic"
git commit -m "WIP: add tests"
git commit -m "WIP: update docs"
```

**Before Creating PR:**
```bash
# Squash or rebase WIP commits into logical units
git rebase -i HEAD~5

# Result: 1-3 well-structured commits
# - feat(config): add config:show command
# - test(config): add unit tests for config:show
# - docs(config): document config:show command
```

**Squash Merge:** The final squash merge combines all commits into one, so individual commit quality on the branch is less critical than the PR title.

### Security Best Practices

**Never Commit:**
- API keys or tokens
- Passwords or credentials
- Private keys
- `.env` files with secrets
- User data

**Always:**
- Use environment variables for secrets
- Add sensitive patterns to `.gitignore`
- Review diffs before committing
- Use GitHub secrets for CI/CD

### Documentation Standards

**Update Documentation When:**
- Adding new commands
- Changing command behavior
- Adding configuration options
- Fixing significant bugs
- Changing workflows

**Documentation Files:**
- `README.md` - User-facing overview
- `docs/USER_GUIDE.md` - Detailed usage
- `docs/DEVELOPMENT.md` - Developer guide
- `CHANGELOG.md` - Automated by Release-Please
- Code comments - Complex logic

### Release Timing

**When to Merge Release PRs:**
- During business hours
- When team is available for hotfixes
- After adequate testing
- Not on Fridays (avoid weekend issues)
- Coordinate with users if breaking changes

**Emergency Releases:**

```bash
# For critical security fixes
# 1. Create fix branch
git checkout -b fix/critical-security-issue

# 2. Implement fix
# ... make changes ...

# 3. Create PR with priority label
gh pr create --label "priority: critical" --title "fix: critical security issue"

# 4. Fast-track review
# 5. Merge immediately
# 6. Merge release PR ASAP

# 7. Notify users
gh release view v2.3.1
# Announce in communication channels
```

## Getting Help

### Resources

- **GitHub Discussions**: [Firecrown-Media/stax/discussions](https://github.com/Firecrown-Media/stax/discussions)
- **Issues**: [Firecrown-Media/stax/issues](https://github.com/Firecrown-Media/stax/issues)
- **Documentation**: `/docs` directory
- **Code of Conduct**: `CODE_OF_CONDUCT.md`

### Questions?

**For development questions:**
1. Check existing documentation
2. Search closed issues and PRs
3. Ask in GitHub Discussions
4. Create an issue with the `question` label

**For bug reports:**
1. Check existing issues
2. Create new issue with `bug` label
3. Include reproduction steps
4. Include version information

**For feature requests:**
1. Check existing issues
2. Create new issue with `enhancement` label
3. Describe the use case
4. Provide examples

## Summary

The Stax CLI uses a fully automated workflow from feature development to user installation:

```
1. Create feature branch: feature/my-feature
2. Commit with conventional format: feat(scope): description
3. Create PR and get it reviewed
4. Squash merge to main
5. Release-Please creates/updates release PR
6. Merge release PR
7. GitHub release created automatically
8. GoReleaser builds and uploads binaries
9. Homebrew formula updated automatically
10. Users upgrade: brew upgrade firecrown-media/stax/stax
```

Key takeaways:
- Use conventional commits for automatic versioning
- Feature branches for all changes
- Squash merge to maintain clean history
- Let automation handle releases
- Document everything
- Test thoroughly before merging

Happy contributing!

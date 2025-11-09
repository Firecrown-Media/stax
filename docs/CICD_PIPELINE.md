# CI/CD Pipeline Documentation

## Overview

Stax uses a comprehensive CI/CD pipeline built on GitHub Actions to ensure code quality, automate testing, and streamline releases. This document describes the complete pipeline architecture and workflows.

## Pipeline Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Developer Actions                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │
              ┌───────────────┼───────────────┐
              │               │               │
              ▼               ▼               ▼
        ┌──────────┐    ┌──────────┐    ┌──────────┐
        │  Push    │    │   Pull   │    │   Tag    │
        │  Code    │    │  Request │    │  Push    │
        └──────────┘    └──────────┘    └──────────┘
              │               │               │
              └───────────────┼───────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        GitHub Actions                            │
│                                                                   │
│  ┌─────────────────┐  ┌──────────────────┐  ┌────────────────┐ │
│  │  Test Workflow  │  │  Release         │  │  Version Bump  │ │
│  │                 │  │  Workflow        │  │  Workflow      │ │
│  │  - Unit Tests   │  │                  │  │                │ │
│  │  - Integration  │  │  - Build         │  │  - Calculate   │ │
│  │  - Security     │  │  - Test          │  │    Version     │ │
│  │  - Coverage     │  │  - Package       │  │  - Create Tag  │ │
│  │  - Linting      │  │  - Release       │  │                │ │
│  │  - Build Check  │  │  - Deploy        │  │                │ │
│  └─────────────────┘  └──────────────────┘  └────────────────┘ │
│                              │                                   │
└──────────────────────────────┼───────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │   GoReleaser     │
                    │                  │
                    │  - Build Bins    │
                    │  - Create SHA256 │
                    │  - Create Assets │
                    │  - Changelog     │
                    └──────────────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
              ▼               ▼               ▼
      ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
      │   GitHub    │  │  Homebrew   │  │  Archives   │
      │   Release   │  │     Tap     │  │  & SHA256   │
      └─────────────┘  └─────────────┘  └─────────────┘
              │               │               │
              └───────────────┼───────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │   End Users      │
                    │                  │
                    │  brew install    │
                    │  Direct download │
                    └──────────────────┘
```

## Workflows

### 1. Test Workflow

**Trigger**: Push to any branch, Pull Requests to main/develop

**File**: `.github/workflows/test.yml`

**Jobs**:

1. **Unit Tests**
   - Matrix: Go 1.22, 1.23
   - Runs: `make test-unit`
   - Runs: `make test-security`
   - Duration: ~2-3 minutes

2. **Integration Tests**
   - Depends on: Unit tests passing
   - Runs: `make test-integration`
   - Duration: ~3-5 minutes

3. **Code Coverage**
   - Depends on: Unit tests passing
   - Runs: `make test-coverage`
   - Uploads: Coverage report to Codecov
   - Artifacts: Coverage HTML report
   - Duration: ~2-3 minutes

4. **Code Quality**
   - Parallel with unit tests
   - Checks: `gofmt`, `go vet`, `golangci-lint`
   - Duration: ~1-2 minutes

5. **Build Verification**
   - Depends on: Unit tests, code quality
   - Runs: `make build`
   - Verifies: Binary executes (`./stax --version`)
   - Artifacts: Binary (7 day retention)
   - Duration: ~1 minute

6. **Test Summary**
   - Depends on: All jobs
   - Aggregates: Results from all jobs
   - Fails if: Any job failed
   - Duration: ~10 seconds

**Total Duration**: ~5-10 minutes (parallel execution)

**Failure Handling**:
- Stops deployment if any test fails
- Posts status to pull requests
- Sends notifications to team

### 2. Release Workflow

**Trigger**: Tag push matching `v*` pattern (e.g., v1.2.3)

**File**: `.github/workflows/release.yml`

**Jobs**:

1. **GoReleaser Job**
   - Checkout: Fetches full git history
   - Go Setup: Installs Go 1.22
   - Tests: Runs unit and security tests
   - Build: Compiles for all platforms
   - Package: Creates tar.gz archives
   - Checksum: Generates SHA256 sums
   - Release: Creates GitHub release
   - Homebrew: Updates formula in tap
   - Duration: ~5-10 minutes

**Platforms Built**:
- macOS Intel (darwin/amd64)
- macOS Apple Silicon (darwin/arm64)
- Linux Intel (linux/amd64)
- Linux ARM (linux/arm64)

**Assets Created**:
- Binaries for each platform
- tar.gz archives with docs
- checksums.txt
- Automated changelog

**Environment Variables**:
- `GITHUB_TOKEN`: Automatic (GitHub provides)
- `HOMEBREW_TAP_TOKEN`: Manual secret (required)

**Success Criteria**:
- All tests pass
- All binaries build successfully
- GitHub release created
- Homebrew formula updated
- Assets uploaded

**Failure Recovery**:
- Delete bad tag: `git push origin :refs/tags/vX.Y.Z`
- Fix issues
- Create new release

### 3. Version Bump Workflow

**Trigger**: Manual workflow dispatch via GitHub UI

**File**: `.github/workflows/version-bump.yml`

**Input**: Version type (patch/minor/major)

**Steps**:

1. **Get Current Version**
   - Fetches latest tag
   - Defaults to v0.0.0 if no tags

2. **Calculate New Version**
   - Parses version components
   - Increments based on type:
     - `patch`: 1.2.3 → 1.2.4
     - `minor`: 1.2.3 → 1.3.0
     - `major`: 1.2.3 → 2.0.0

3. **Create and Push Tag**
   - Creates annotated tag
   - Pushes to origin
   - Triggers release workflow automatically

**Duration**: ~30 seconds

**Usage**:
```
1. Go to Actions tab
2. Select "Version Bump"
3. Click "Run workflow"
4. Choose version type
5. Click "Run workflow"
```

## Pipeline Flow

### Development Flow

```
1. Developer pushes code
   ↓
2. Test workflow runs
   ↓
3. All tests must pass
   ↓
4. Code reviewed (manual)
   ↓
5. Merge to main
```

### Release Flow

```
1. Developer triggers version bump OR pushes tag
   ↓
2. Version bump creates tag (if using workflow)
   ↓
3. Tag push triggers release workflow
   ↓
4. Release workflow runs tests
   ↓
5. GoReleaser builds binaries
   ↓
6. GitHub release created
   ↓
7. Homebrew formula updated
   ↓
8. Users can install via brew
```

### Hotfix Flow

```
1. Create hotfix branch
   ↓
2. Fix critical issue
   ↓
3. Test workflow validates
   ↓
4. Merge to main
   ↓
5. Tag patch version
   ↓
6. Emergency release
   ↓
7. Notify users
```

## Secrets and Configuration

### Required Secrets

1. **GITHUB_TOKEN**
   - Type: Automatic
   - Scope: Repository access
   - Used for: Creating releases, uploading assets
   - Setup: Automatic by GitHub

2. **HOMEBREW_TAP_TOKEN**
   - Type: Manual Personal Access Token
   - Scope: `repo` (full repository control)
   - Used for: Pushing formula updates to homebrew-tap
   - Setup: Settings → Secrets → Actions → New secret
   - Expiration: None or very long
   - Create: https://github.com/settings/tokens

### Environment Variables

Set in workflows:
- `VERSION`: From git tag
- `GIT_COMMIT`: Short commit hash
- `BUILD_DATE`: ISO 8601 timestamp

Injected at build time via ldflags.

## Monitoring and Alerts

### Success Indicators

**Test Workflow**:
- All jobs green
- Coverage maintains or improves
- No linting warnings
- Build artifacts created

**Release Workflow**:
- GitHub release published
- All assets uploaded
- Homebrew formula updated
- Version tag exists

### Failure Indicators

**Test Workflow**:
- Red X on commit/PR
- Failed job in Actions tab
- Coverage drop
- Linting errors

**Release Workflow**:
- No GitHub release created
- Missing assets
- Formula not updated
- Build errors in logs

### Monitoring Locations

1. **GitHub Actions Tab**
   - https://github.com/firecrown-media/stax/actions
   - Shows all workflow runs
   - Provides detailed logs

2. **Releases Page**
   - https://github.com/firecrown-media/stax/releases
   - Lists all published releases
   - Shows download counts

3. **Codecov Dashboard**
   - Shows coverage trends
   - Highlights coverage drops

4. **Homebrew Tap**
   - https://github.com/firecrown-media/homebrew-tap
   - Verify formula updates
   - Check commit history

## Performance Optimization

### Test Workflow Optimizations

1. **Matrix Strategy**: Multiple Go versions in parallel
2. **Job Dependencies**: Unit tests gate integration tests
3. **Caching**: Go modules cached between runs
4. **Conditional Steps**: Skip redundant operations
5. **Parallel Jobs**: Independent jobs run simultaneously

### Release Workflow Optimizations

1. **Fetch Depth**: Only fetch necessary history
2. **Build Parallelization**: Multiple platforms at once
3. **Asset Compression**: Optimize archive sizes
4. **Changelog Generation**: Automated from commits

### Cache Strategy

**Cached Items**:
- Go modules (`go.sum`)
- Build dependencies
- Test results (for retries)

**Cache Keys**:
- Based on `go.sum` hash
- Platform-specific
- Updated on dependency changes

## Troubleshooting

### Common Issues

#### Tests Pass Locally But Fail in CI

**Causes**:
- Environment differences
- Race conditions
- Network dependencies

**Solutions**:
```bash
# Run tests with race detection
make test-unit

# Run in isolation
go test -v -count=1 ./...

# Check for timing issues
go test -race -count=10 ./...
```

#### Release Workflow Fails

**Check**:
1. Test results
2. GoReleaser configuration
3. Token permissions
4. Network connectivity

**Fix**:
```bash
# Test release locally
goreleaser release --snapshot --clean

# Validate config
goreleaser check
```

#### Homebrew Formula Not Updated

**Verify**:
1. HOMEBREW_TAP_TOKEN exists
2. Token has repo scope
3. Token hasn't expired
4. Repository name is correct

**Test Manually**:
```bash
# Clone tap
git clone https://github.com/firecrown-media/homebrew-tap.git

# Try pushing
git commit --allow-empty -m "test"
git push  # Should succeed
```

#### Version Bump Doesn't Trigger Release

**Check**:
1. Tag was pushed successfully
2. Tag matches pattern `v*`
3. Release workflow enabled

**Fix**:
```bash
# Verify tag exists
git tag -l

# Re-push tag
git push origin v1.2.3
```

## Best Practices

### For Developers

1. **Run tests locally** before pushing:
   ```bash
   make test-all
   make verify
   ```

2. **Keep commits atomic**: One logical change per commit

3. **Write descriptive commit messages**: Used in changelogs

4. **Wait for CI**: Don't merge failing PRs

### For Releases

1. **Use version bump workflow**: Ensures consistency

2. **Test releases locally**:
   ```bash
   goreleaser build --snapshot
   ```

3. **Document breaking changes**: Update docs before release

4. **Verify Homebrew update**: Test formula after release

### For Maintenance

1. **Monitor workflow runs**: Check for patterns in failures

2. **Update dependencies**: Keep actions up to date

3. **Review logs**: Understand performance bottlenecks

4. **Rotate tokens**: Renew HOMEBREW_TAP_TOKEN annually

## Metrics and KPIs

### Workflow Metrics

- **Test Workflow Success Rate**: Target >95%
- **Average Test Duration**: Track trends
- **Code Coverage**: Maintain >70%
- **Build Success Rate**: Target 100%

### Release Metrics

- **Release Frequency**: Track releases per month
- **Time to Release**: From tag to published
- **Download Counts**: Monitor adoption
- **Formula Install Success**: Track via Homebrew

## Security Considerations

### Token Security

- Never log secrets
- Rotate tokens regularly
- Use minimal required scopes
- Monitor token usage

### Build Security

- Pin action versions
- Verify checksums
- Sign releases (future)
- Audit dependencies

### Access Control

- Require reviews for merges
- Protect main branch
- Limit who can create tags
- Monitor workflow changes

## Future Enhancements

### Planned Improvements

1. **Signed Releases**: GPG signatures for binaries
2. **Docker Images**: Publish to registries
3. **Multi-arch Support**: More platforms
4. **Performance Tests**: Automated benchmarking
5. **Security Scanning**: SAST/DAST integration
6. **Release Notes**: Enhanced automation
7. **Rollback Automation**: One-click rollback
8. **Canary Releases**: Staged rollouts

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Homebrew Formula Docs](https://docs.brew.sh/Formula-Cookbook)
- [Semantic Versioning](https://semver.org/)
- [Codecov Documentation](https://docs.codecov.com/)

## Summary

The Stax CI/CD pipeline provides:

- **Automated Testing**: Every commit is tested
- **Quality Gates**: Code must meet standards
- **Easy Releases**: One-click version bumps
- **Automated Distribution**: Homebrew updates automatically
- **Visibility**: Clear status for all changes
- **Reliability**: Consistent, reproducible builds

This ensures high code quality and enables rapid, confident releases.

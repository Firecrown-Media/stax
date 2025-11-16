# Mirror Sync Quick Reference

**Quick reference card for the Stax public mirror sync workflow**

## Overview

- **Private Repo**: `Firecrown-Media/stax` (development)
- **Public Repo**: `Firecrown-Media/stax-public` (distribution)
- **Sync**: Automatic on release + manual trigger
- **Purpose**: Secure open-source distribution

## How It Works

```
Private Repo Release → Sync Workflow → Public Repo → GoReleaser → Homebrew
```

## Automatic Sync

**When**: Automatically when you publish a release
**What**: Syncs release tag to public repository
**No Action Required**: Just publish release as normal

## Manual Sync

### Via GitHub UI

1. Go to: https://github.com/Firecrown-Media/stax
2. Click: Actions → Sync Public Mirror
3. Click: Run workflow
4. Enter: Tag name (or leave empty for latest)
5. Click: Run workflow button

### What Gets Synced

**Included**:
- All source code (*.go, cmd/, pkg/, internal/)
- Build configs (.goreleaser.yml, Makefile)
- Documentation (docs/)
- README (public version)

**Excluded**:
- .claude/ directory
- *.claude.md files
- dist/ directory
- .DS_Store files
- Build artifacts

## Quick Commands

### Check Last Sync
```bash
# In private repo
gh run list --workflow=sync-public-mirror.yml --limit 5
```

### Verify Public Repo
```bash
# Clone and check
git clone https://github.com/Firecrown-Media/stax-public.git
cd stax-public
ls -la .claude/  # Should not exist
git tag -l       # Should have release tags
```

### Manual Sync via CLI
```bash
# In private repo
gh workflow run sync-public-mirror.yml -f tag=v2.4.0
```

## Common Tasks

### Re-sync Failed Release
1. Go to Actions → Sync Public Mirror
2. Run workflow with specific tag
3. Check logs for errors

### Verify Sync Succeeded
1. Go to: https://github.com/Firecrown-Media/stax-public
2. Check: Tags page has your release
3. Check: Main branch updated recently
4. Check: No .claude/ directory

### Check Workflow Status
1. Go to: https://github.com/Firecrown-Media/stax/actions
2. Filter: Sync Public Mirror workflow
3. Check: Recent runs succeeded

## Troubleshooting

### Sync Failed

**Check**:
1. Workflow logs in Actions tab
2. Deploy key is configured
3. Tag exists in private repo
4. Public repo is accessible

**Fix**:
- Re-run workflow manually
- Check `PUBLIC_MIRROR_DEPLOY_KEY` secret
- Verify deploy key has write access

### Wrong Files in Public Repo

**Fix**:
1. Update cleaning rules in workflow
2. Re-run sync for tag
3. Force push if needed

### GoReleaser Failed

**Check**:
1. GoReleaser targets stax-public
2. Release created in public repo
3. `HOMEBREW_TAP_TOKEN` secret set

## Security Checklist

- [ ] `PUBLIC_MIRROR_DEPLOY_KEY` secret configured
- [ ] Deploy key added to stax-public
- [ ] Deploy key has write access only
- [ ] No .claude/ in public repo
- [ ] No *.claude.md in public repo
- [ ] No secrets in public repo

## Key Files

| File | Purpose |
|------|---------|
| `.github/workflows/sync-public-mirror.yml` | Sync workflow |
| `.goreleaser.yml` | Release config (targets stax-public) |
| `docs/PUBLIC_MIRROR_README.md` | Public README template |
| `docs/MIRROR_SYNC.md` | Full documentation |
| `docs/MIRROR_SYNC_TESTING.md` | Testing procedures |

## Important Notes

1. **Force Push**: Workflow uses force push - this is intentional (mirror, not collaborative)
2. **README Swap**: Public repo gets different README from `docs/PUBLIC_MIRROR_README.md`
3. **No History**: Public repo doesn't have full git history, just release states
4. **Tag Immutability**: Once synced, tags shouldn't change (re-sync if needed)

## Monitoring

### Daily
- Check recent workflow runs succeeded

### Weekly
- Review sync status in Actions
- Verify public repo matches releases

### Monthly
- Audit file cleaning effectiveness
- Review deploy key security
- Update documentation if needed

## Getting Help

1. **Documentation**: `docs/MIRROR_SYNC.md`
2. **Testing Guide**: `docs/MIRROR_SYNC_TESTING.md`
3. **Workflow File**: `.github/workflows/sync-public-mirror.yml`
4. **Team**: Contact DevOps
5. **Issues**: Create issue in private repo

## Quick URLs

- **Private Repo**: https://github.com/Firecrown-Media/stax
- **Public Repo**: https://github.com/Firecrown-Media/stax-public
- **Actions**: https://github.com/Firecrown-Media/stax/actions
- **Workflow**: https://github.com/Firecrown-Media/stax/actions/workflows/sync-public-mirror.yml

## Secrets Required

| Secret | Location | Purpose |
|--------|----------|---------|
| `PUBLIC_MIRROR_DEPLOY_KEY` | Private repo secrets | SSH auth to public repo |
| `HOMEBREW_TAP_TOKEN` | Private repo secrets | Update Homebrew formula |

## Workflow Diagram

```
┌─────────────┐
│   Release   │
│  Published  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Trigger   │
│   Workflow  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Checkout   │
│  at Tag     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│    Clean    │
│   Files     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Update    │
│   README    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Force Push  │
│  to Public  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Verify    │
│    Sync     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Success   │
└─────────────┘
```

## File Cleaning Patterns

```bash
# Removed automatically
.claude/
*.claude.md
CLAUDE.md
dist/
stax
.DS_Store
.cache/
tmp/

# Kept in public repo
*.go
cmd/
pkg/
internal/
docs/
.goreleaser.yml
Makefile
go.mod
go.sum
```

## Version Info

- **Implemented**: 2025-11-15
- **Workflow Version**: 1.0
- **Documentation**: docs/MIRROR_SYNC.md

---

**Need more details?** See `docs/MIRROR_SYNC.md` for complete documentation.

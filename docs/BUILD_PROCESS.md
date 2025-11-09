# Build Process Documentation

## Overview

The Stax build system integrates with existing WordPress project build workflows, providing an easier interface while preserving battle-tested build scripts. This documentation covers the build process, available commands, and how to customize the build for your needs.

## Build Architecture

### Standard Build Workflow

WordPress projects typically use a modular build system with these components:

1. **scripts/build.sh** - Main build orchestrator
2. **scripts/build/10-mu-plugins.sh** - Builds MU plugins (composer install)
3. **scripts/build/20-theme.sh** - Builds themes (npm install + build)

Stax wraps these scripts with additional functionality:
- Build status detection
- Dependency caching
- Progress feedback
- Error handling
- Validation

## Build Commands

### Main Build Command

```bash
# Run full build
stax build

# Force rebuild even if not needed
stax build --force

# Clean and rebuild
stax build --clean

# Verbose output
stax build --verbose
```

The main build command:
1. Checks if a build is needed (skip with `--force`)
2. Optionally cleans build artifacts (`--clean`)
3. Executes `scripts/build.sh`
4. Validates build output
5. Reports build time and status

### Build Subcommands

#### Composer Only

```bash
# Install composer dependencies only
stax build:composer
```

Runs composer install for MU plugins without building themes.

#### NPM Only

```bash
# Install npm dependencies and build themes
stax build:npm
```

Builds themes without installing PHP dependencies.

#### Specific Theme

```bash
# Build parent theme
stax build:theme firecrown-parent

# Build child theme
stax build:theme firecrown-child
```

Builds a single theme. Runs:
- `npm install`
- `npm run build`
- `composer install` (if composer.json exists)

#### Clean Build Artifacts

```bash
# Remove all build artifacts
stax build:clean
```

Removes:
- `vendor/` directories
- `node_modules/` directories
- `build/` output directories

#### Build Status

```bash
# Check if build is needed
stax build:status
```

Shows:
- Composer dependency status
- NPM dependency status
- Last build time
- Whether build is needed and why

#### List Build Scripts

```bash
# Show available build scripts
stax build:scripts
```

Lists all scripts in `scripts/build/` with their execution order.

## Code Quality (Linting)

### Lint Commands

```bash
# Run PHPCS on all PHP files
stax lint

# Auto-fix issues with PHPCBF
stax lint:fix

# Lint only staged files (pre-commit)
stax lint:staged

# Check specific files
stax lint wp-content/mu-plugins/firecrown/src/
```

### PHPCS Configuration

Stax looks for PHPCS configuration in this order:
1. `.phpcs.xml.dist`
2. `phpcs.xml`
3. `.phpcs.xml`

Example `.phpcs.xml.dist`:
```xml
<?xml version="1.0"?>
<ruleset>
    <description>Firecrown coding standards</description>
    <file>.</file>

    <!-- Exclude patterns -->
    <exclude-pattern>*/vendor/*</exclude-pattern>
    <exclude-pattern>*/node_modules/*</exclude-pattern>
    <exclude-pattern>*/build/*</exclude-pattern>

    <!-- Coding standard -->
    <rule ref="TCC" />
</ruleset>
```

### Pre-commit Hooks

Firecrown uses Husky for pre-commit hooks that automatically run linting:

```bash
# Verify Husky configuration
ls .husky/pre-commit

# Test pre-commit hook
stax lint:staged
```

The pre-commit hook runs `composer run-script lint` before allowing commits.

To bypass hooks (use with caution):
```bash
git commit --no-verify
```

## Development Mode

### Development Commands

```bash
# Start dev mode (npm start with HMR)
stax dev

# Start dev mode for specific theme
stax dev --theme=firecrown-child

# Stop dev mode
stax dev:stop

# Watch for changes and auto-rebuild
stax dev:watch
```

### Development Workflow

1. **Start dev mode**: `stax dev`
   - Runs `npm start` in theme directory
   - Enables Hot Module Reloading (HMR)
   - Watches for file changes
   - Press Ctrl+C to stop

2. **File watching**: `stax dev:watch`
   - Watches theme `src/` directories
   - Watches MU plugin `src/` directory
   - Auto-rebuilds on changes
   - Includes debouncing (500ms)

### What's Watched

The file watcher monitors:
- `wp-content/mu-plugins/firecrown/src/**`
- `wp-content/themes/firecrown-parent/src/**`
- `wp-content/themes/firecrown-child/src/**`

And ignores:
- `node_modules/`
- `vendor/`
- `build/`
- `.git/`
- `.DS_Store`

## Build Status Detection

Stax automatically determines if a build is needed by checking:

1. **Build artifacts exist**
   - `wp-content/mu-plugins/firecrown/vendor/`
   - `wp-content/themes/firecrown-parent/build/`
   - `wp-content/themes/firecrown-parent/node_modules/`

2. **Dependencies installed**
   - Composer lock file exists
   - NPM lock file exists
   - Vendor directories exist

3. **Dependencies up to date**
   - `composer.json` not modified after `composer.lock`
   - `package.json` not modified after `package-lock.json`

4. **Source files modified**
   - Source files not newer than build output

### Example Status Check

```bash
$ stax build:status

==> Build Status

✓ Build script found: scripts/build.sh
  Custom build scripts (2):
  - 10-mu-plugins.sh
  - 20-theme.sh

Composer Dependencies
✓ Installed
  Lock file: 2024-11-08 10:30:15
  Vendor dir: 2024-11-08 10:30:25

NPM Dependencies
✓ Installed
  Lock file: 2024-11-08 10:31:00
  node_modules: 2024-11-08 10:31:45

Last Build
  Time: 2024-11-08 10:32:00
  Age: 2 hours ago

Build Status
✓ Build is up to date
```

## Configuration

Build configuration in `.stax.yml`:

```yaml
build:
  scripts:
    main: scripts/build.sh
    pre_build:
      - echo "Starting build..."
    post_build:
      - echo "Build complete!"

  composer:
    install_args: "--no-dev --prefer-dist --ignore-platform-reqs"
    timeout: 300
    ignore_platform_reqs: true
    optimize: true
    no_dev: true

  npm:
    install_args: "--legacy-peer-deps"
    build_command: "build"
    dev_command: "start"
    timeout: 600
    legacy_peer_deps: true

  phpcs:
    config: .phpcs.xml.dist
    standard: TCC
    extensions: php
    ignore: "vendor/*,node_modules/*,build/*"
    show_sniffs: true

  hooks:
    pre_commit: true
    pre_push: false
    commit_msg: false

  watch:
    enabled: true
    paths:
      - wp-content/mu-plugins/firecrown/src
      - wp-content/themes/firecrown-parent/src
      - wp-content/themes/firecrown-child/src
```

## Build Scripts

### Creating Build Scripts

If your project doesn't have build scripts, generate them:

```bash
# This will be available in a future version
stax build:generate
```

This creates:
- `scripts/build.sh` - Main orchestrator
- `scripts/build/10-mu-plugins.sh` - MU plugin build
- `scripts/build/20-theme.sh` - Theme build

### Custom Build Scripts

Add custom build scripts to `scripts/build/`:

```bash
scripts/build/30-custom.sh
```

Scripts are executed in numerical order (10, 20, 30, etc.).

Example custom script:
```bash
#!/bin/bash
set -e

echo "Running custom build step..."

# Your custom build logic here
# For example, compile additional assets, run tests, etc.

echo "Custom build complete"
```

### Build Script Best Practices

1. **Use `set -e`** - Exit on error
2. **Echo progress** - Show what's happening
3. **Check prerequisites** - Verify files/directories exist
4. **Use relative paths** - Scripts run from project root
5. **Handle errors gracefully** - Exit cleanly on failure

## Integration with DDEV

Builds work both locally and inside DDEV containers:

```bash
# Build locally
stax build

# Build inside DDEV
ddev exec stax build

# Or via DDEV custom command
ddev build
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Build and Test

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v2

      - name: Setup PHP
        uses: shivammathur/setup-php@v2
        with:
          php-version: 8.1

      - name: Setup Node
        uses: actions/setup-node@v2
        with:
          node-version: 20

      - name: Install Stax
        run: |
          curl -L https://github.com/firecrown-media/stax/releases/latest/download/stax -o /usr/local/bin/stax
          chmod +x /usr/local/bin/stax

      - name: Build project
        run: stax build

      - name: Run linting
        run: stax lint

      - name: Validate build
        run: stax build:status
```

## Troubleshooting

### Build Fails

```bash
# Check build status
stax build:status

# Clean and rebuild
stax build --clean

# Verbose output for debugging
stax build --verbose
```

### Missing Dependencies

```bash
# Check composer
composer --version

# Check npm
npm --version

# Check node
node --version
```

### PHPCS Not Found

```bash
# Install via composer
composer require --dev squizlabs/php_codesniffer

# Or install globally
composer global require squizlabs/php_codesniffer
```

### Permission Errors

```bash
# Make build scripts executable
chmod +x scripts/build.sh
chmod +x scripts/build/*.sh
```

### Clean Start

```bash
# Remove everything and start fresh
stax build:clean
rm -rf vendor node_modules
stax build
```

## Performance Tips

1. **Use build status** - Skip unnecessary rebuilds
   ```bash
   stax build:status && stax build
   ```

2. **Parallel builds** (future feature)
   ```bash
   stax build --parallel
   ```

3. **Cache dependencies** - Don't clean unless needed
   - `node_modules/` and `vendor/` can be cached in CI/CD

4. **Use dev mode** - Faster feedback during development
   ```bash
   stax dev
   ```

## Error Messages and Solutions

### "Build script not found"
**Solution**: Generate build scripts:
```bash
# Copy from templates
cp -r templates/scripts/* scripts/
```

### "Composer dependencies not installed"
**Solution**: Run composer install:
```bash
stax build:composer
```

### "NPM dependencies not installed"
**Solution**: Run npm install:
```bash
stax build:npm
```

### "PHPCS validation failed"
**Solution**: Fix code quality issues:
```bash
stax lint:fix
```

### "Build artifacts missing"
**Solution**: Run full build:
```bash
stax build --force
```

## Additional Resources

- [PHPCS Documentation](https://github.com/squizlabs/PHP_CodeSniffer)
- [NPM Scripts](https://docs.npmjs.com/cli/v8/using-npm/scripts)
- [Composer Documentation](https://getcomposer.org/doc/)
- [WordPress Coding Standards](https://developer.wordpress.org/coding-standards/)

## Summary

The Stax build system provides:
- ✓ Automated build process
- ✓ Intelligent build detection
- ✓ Code quality enforcement
- ✓ Development mode with HMR
- ✓ Pre-commit hooks
- ✓ CI/CD integration
- ✓ Comprehensive error handling

For questions or issues, refer to the main Stax documentation or open an issue on GitHub.

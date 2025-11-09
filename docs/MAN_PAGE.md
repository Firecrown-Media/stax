# Stax Man Page

## Overview

Stax includes a comprehensive Unix man page that can be viewed with the standard `man` command.

## Viewing the Man Page

### If Installed via Homebrew

The man page is automatically installed:
```bash
man stax
```

### If Installed from Source

After running `make install`:
```bash
man stax
```

### Preview Without Installing

```bash
# Generate and preview
make man-preview

# Or manually
make man
man dist/man/stax.1
```

## Generating the Man Page

### Automatic Generation

```bash
# Using the CLI
./stax man

# Using make
make man
```

### Manual Installation

```bash
# Generate
make man

# Install
sudo cp dist/man/stax.1 /usr/local/share/man/man1/
sudo mandb  # Update man database (Linux)
# or
sudo makewhatis  # Update man database (macOS)
```

## Man Page Sections

The stax man page includes:

1. **NAME** - Brief description
2. **SYNOPSIS** - Command syntax
3. **DESCRIPTION** - Detailed description
4. **OPTIONS** - Global options and flags
5. **COMMANDS** - All available commands organized by category
6. **EXAMPLES** - Common usage examples
7. **FILES** - Configuration and data files
8. **ENVIRONMENT VARIABLES** - Available environment variables
9. **EXIT STATUS** - Exit codes and meanings
10. **DEPENDENCIES** - Required software
11. **SEE ALSO** - Related commands and documentation
12. **BUGS** - How to report issues
13. **AUTHOR** - Project information
14. **COPYRIGHT** - License information
15. **VERSION** - Current version

## Updating the Man Page

The man page is automatically updated when:
- Creating a new release (via GoReleaser)
- Running `make install`
- Running `make man`

To manually update the template:
1. Edit `docs/stax.1.template`
2. Run `make man` to regenerate
3. Preview with `make man-preview`

## Man Page Format

The man page is written in groff/troff format, the standard Unix man page format. This ensures:
- Compatibility with all Unix-like systems
- Proper rendering in terminal pagers
- Searchability via `man -k stax`
- Integration with system documentation

## Makefile Targets

### `make man`
Generate the man page to `dist/man/stax.1`

### `make man-preview`
Generate and preview the man page using the `man` command

### `make man-install`
Generate and install the man page to `/usr/local/share/man/man1/`

### `make man-uninstall`
Remove the installed man page

## CLI Command

The `stax man` command generates the man page:

```bash
# Generate to current directory
stax man

# Generate to specific location
stax man -o /usr/local/share/man/man1/

# View the generated man page
man ./stax.1
```

## Integration with Build Process

The man page is automatically generated during:

1. **Local builds**: `make install` generates and installs the man page
2. **Release builds**: GoReleaser generates the man page via `scripts/generate-man.sh`
3. **Homebrew installation**: The formula installs the man page to the appropriate location

## Troubleshooting

### Man page not found

```bash
# Rebuild man database
sudo mandb  # Linux
sudo makewhatis  # macOS

# Check installation
ls -l /usr/local/share/man/man1/stax.1
```

### Man page outdated

```bash
# Regenerate and reinstall
make man-install
```

### Preview shows old version

```bash
# Clear man cache
rm -rf ~/.cache/man
man stax
```

## Man Page Sources

The man page is generated from two sources:

1. **Template**: `docs/stax.1.template` - Static man page template
2. **Cobra**: `cmd/man.go` - Dynamic generation from Cobra command metadata

The build process uses both sources to create a comprehensive man page.

## Searching Man Pages

### Find man pages related to stax
```bash
man -k stax
```

### Search within the man page
```bash
# Open man page
man stax

# Then press '/' to search
# Type your search term
# Press 'n' for next match
# Press 'N' for previous match
```

## Man Page Navigation

When viewing the man page:
- `Space` - Page down
- `b` - Page up
- `/pattern` - Search forward
- `?pattern` - Search backward
- `n` - Next match
- `N` - Previous match
- `q` - Quit
- `h` - Help

## Man Page Sections Explained

### NAME
The name and brief description of the tool.

### SYNOPSIS
Shows the command syntax and available options.

### DESCRIPTION
Detailed explanation of what stax does and its key features.

### OPTIONS
Global flags that can be used with any command.

### COMMANDS
All available commands organized by functional category:
- Project Management
- Database Operations
- Build & Development
- Configuration
- Provider Management

### EXAMPLES
Real-world usage examples showing common workflows.

### FILES
Configuration files and data directories used by stax.

### ENVIRONMENT VARIABLES
Environment variables that affect stax behavior.

### EXIT STATUS
Numeric exit codes and their meanings for scripting.

### DEPENDENCIES
External tools required for stax to function.

### SEE ALSO
Related commands and online documentation.

## Best Practices

### For Users
- Install the man page via `make install` or Homebrew
- Use `man stax` for quick reference
- Search within the man page using `/`
- Combine with `--help` for interactive help

### For Developers
- Update `docs/stax.1.template` when adding new features
- Test man page generation with `make man-preview`
- Ensure man page is included in releases
- Keep examples up-to-date

## Version Information

The man page displays the current version in:
- The header (via `{{.Version}}` template variable)
- The VERSION section
- The SOURCE field (e.g., "Stax v1.0.0")

This ensures users can verify they're viewing the correct documentation version.

## Distribution

The man page is distributed:
1. **Source**: In the `docs/` directory as a template
2. **Binary**: Generated during build to `dist/man/stax.1`
3. **Archive**: Included in release tarballs
4. **Homebrew**: Installed automatically via the formula

## Related Documentation

- [Installation Guide](INSTALLATION.md) - How to install stax
- [User Guide](USER_GUIDE.md) - Comprehensive usage guide
- [Command Reference](COMMAND_REFERENCE.md) - Detailed command documentation
- [Quick Start](QUICK_START.md) - Get started quickly

## Support

If you find issues with the man page:
1. Check that it's up-to-date: `stax --version`
2. Regenerate: `make man-install`
3. Report bugs: https://github.com/firecrown-media/stax/issues

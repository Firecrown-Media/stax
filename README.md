# Stax - WordPress Development CLI

A powerful command-line tool that streamlines WordPress development with DDEV and seamless hosting provider integration.

## What is Stax?

Stax automates the complex setup of local WordPress development environments. It manages Docker containers through DDEV, syncs with hosting providers like WP Engine, and enables instant PHP/MySQL version switching without data loss. Built for teams managing multiple client sites.

## Quick Start (5 minutes)

```bash
# 1. Install dependencies
brew install docker ddev/ddev/ddev

# 2. Install Stax
brew tap firecrown-media/stax
brew install stax

# 3. Create your first WordPress site
stax init my-first-site
stax setup my-first-site --install-wp
stax start my-first-site

# 4. Visit your site
open https://my-first-site.ddev.site
```

That's it! You now have a fully functional WordPress development environment.

## Documentation

- **[Setup Guide](docs/SETUP.md)** - Complete installation and configuration instructions
- **[User Guide](docs/USER_GUIDE.md)** - Daily usage, workflows, and best practices
- **[WP Engine Guide](docs/WPENGINE.md)** - Hosting integration and syncing
- **[Troubleshooting](docs/TROUBLESHOOTING.md)** - Common issues and solutions
- **[Development](docs/DEVELOPMENT.md)** - Contributing to Stax
- **[Multisite](docs/MULTISITE.md)** - WordPress Multisite setup and management

## Quick Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `stax init [name]` | Create new project | `stax init client-site` |
| `stax start [name]` | Start development environment | `stax start client-site` |
| `stax stop [name]` | Stop development environment | `stax stop client-site` |
| `stax status` | Show all running environments | `stax status` |
| `stax wpe sync [install]` | Sync from WP Engine | `stax wpe sync clientinstall` |
| `stax swap php [version]` | Change PHP version | `stax swap php 8.3` |
| `stax swap preset [name]` | Apply environment preset | `stax swap preset modern` |
| `stax wp [command]` | Run WP-CLI commands | `stax wp plugin list` |

## Key Features

✅ **Rapid Setup** - WordPress sites ready in under 2 minutes
✅ **Hot Swap** - Change PHP/MySQL versions without losing data
✅ **WP Engine Integration** - One-command production sync
✅ **Multi-Project** - Manage dozens of sites simultaneously
✅ **Team Friendly** - Shared configurations for consistency

## System Requirements

- macOS or Linux
- Docker Desktop (or compatible Docker provider)
- DDEV
- Homebrew (for easy installation)

For detailed requirements and installation instructions, see the [Setup Guide](docs/SETUP.md).

## Getting Help

- **Documentation**: Start with our [User Guide](docs/USER_GUIDE.md)
- **Issues**: [Report bugs or request features](https://github.com/Firecrown-Media/stax/issues)
- **Support**: Contact dev@firecrown.com

## Contributing

We welcome contributions! See our [Development Guide](docs/DEVELOPMENT.md) for information on:
- Setting up the development environment
- Code structure and architecture
- Adding new features
- Submitting pull requests

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Project Status

Actively maintained by Firecrown Media. Built specifically for agency WordPress development workflows but designed to be extensible for other use cases.

---

**Quick Links**: [Setup](docs/SETUP.md) | [User Guide](docs/USER_GUIDE.md) | [Troubleshooting](docs/TROUBLESHOOTING.md) | [WP Engine](docs/WPENGINE.md)
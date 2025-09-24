# Stax Setup Guide

This guide provides complete instructions for installing and configuring Stax on your system.

## Table of Contents
- [System Requirements](#system-requirements)
- [Installation Methods](#installation-methods)
- [Dependency Installation](#dependency-installation)
- [Configuration](#configuration)
- [Verification](#verification)
- [Next Steps](#next-steps)

## System Requirements

### Operating System
- **macOS**: 12.0 (Monterey) or later
- **Linux**: Ubuntu 20.04+, Debian 10+, Fedora 34+, or compatible distributions
- **Windows**: Not currently supported (use WSL2 for Linux compatibility)

### Hardware Requirements
- **RAM**: Minimum 8GB, 16GB recommended
- **Storage**: 20GB free space for Docker images and project files
- **CPU**: Intel/AMD x64 or Apple Silicon (M1/M2/M3)

### Required Software

| Software | Purpose | Minimum Version |
|----------|---------|-----------------|
| Docker | Container runtime | 20.10+ |
| DDEV | Development environment manager | 1.21+ |
| Git | Version control | 2.0+ |
| Homebrew (Mac) | Package manager | 3.0+ |

### Optional Software
- **WP-CLI**: WordPress command-line interface
- **Go**: Only needed if building from source (1.19+)

## Installation Methods

### Option 1: Homebrew Installation (Recommended for Mac)

The easiest way to install Stax on macOS:

```bash
# Add the Stax tap
brew tap firecrown-media/stax

# Install Stax
brew install stax

# Verify installation
stax --version
```

**Note**: The Homebrew formula does NOT automatically install dependencies. Continue to [Dependency Installation](#dependency-installation).

### Option 2: Install from Source

For customization or development:

```bash
# 1. Install Go (if not already installed)
brew install go  # On Mac
# or
sudo apt install golang-go  # On Ubuntu/Debian

# 2. Clone the repository
git clone https://github.com/Firecrown-Media/stax.git
cd stax

# 3. Build and install
make install

# 4. Verify installation
stax --version
```

If you get "command not found", add Go's bin directory to your PATH:

```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc  # Mac
# or
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc  # Linux
source ~/.zshrc  # or source ~/.bashrc
```

### Option 3: Direct Binary Download

For manual installation without building:

```bash
# Download the latest release (replace VERSION and PLATFORM)
curl -LO https://github.com/Firecrown-Media/stax/releases/download/vVERSION/stax-PLATFORM

# Make it executable
chmod +x stax-PLATFORM

# Move to a directory in your PATH
sudo mv stax-PLATFORM /usr/local/bin/stax

# Verify
stax --version
```

Platform options:
- `darwin-amd64` (Intel Mac)
- `darwin-arm64` (Apple Silicon Mac)
- `linux-amd64` (Linux x64)

## Dependency Installation

### Step 1: Install Docker

Docker provides the container runtime that DDEV uses.

#### macOS

**Option A: Docker Desktop (Easiest)**
```bash
brew install --cask docker
# Open Docker Desktop from Applications
# Complete the setup wizard
```

**Option B: OrbStack (Faster, Commercial)**
```bash
brew install --cask orbstack
# Open OrbStack and complete setup
```

**Option C: Colima (Free, Open Source)**
```bash
brew install colima docker
colima start --cpu 4 --memory 8
```

#### Linux

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker

# Add your user to docker group
sudo usermod -aG docker $USER
# Log out and back in for group changes to take effect
```

#### Verify Docker Installation

```bash
docker --version
docker run hello-world
```

### Step 2: Install DDEV

DDEV manages Docker containers specifically for web development.

#### All Platforms

```bash
# Mac with Homebrew
brew install ddev/ddev/ddev

# Linux
curl -fsSL https://ddev.readthedocs.io/en/stable/users/install/ddev-installation-script/ | bash

# Verify installation
ddev version
```

### Step 3: Install WP-CLI (Optional but Recommended)

WP-CLI provides command-line control of WordPress.

```bash
# Mac with Homebrew
brew install wp-cli

# Linux/Manual installation
curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar
chmod +x wp-cli.phar
sudo mv wp-cli.phar /usr/local/bin/wp

# Verify
wp --version
```

## Configuration

### Global Configuration

Stax uses a global configuration file at `~/.stax.yaml` for defaults:

```bash
# Create the configuration file
cat > ~/.stax.yaml << 'EOF'
# Organization name (optional)
organization: mycompany

# Default environment settings
default_php_version: "8.2"
default_mysql_version: "8.0"
default_webserver: "nginx-fpm"

# WordPress defaults
wordpress:
  admin_user: "admin"
  admin_email: "admin@example.com"

# Verbose output for debugging
verbose: false
EOF
```

### Project Configuration

Each project can have its own `stax.yaml` in the project directory:

```yaml
# Project-specific stax.yaml
name: "client-project"
php_version: "8.1"
database: "mysql:8.0"

wordpress:
  url: "https://client-project.ddev.site"
  title: "Client Project Development"
```

### Environment Variables

For WP Engine integration, set these environment variables:

```bash
# Add to ~/.zshrc (Mac) or ~/.bashrc (Linux)
export WPE_USERNAME="your-wpengine-username"
export WPE_PASSWORD="your-wpengine-password"

# Optional: Set default environment
export WPE_DEFAULT_ENV="staging"

# Reload your shell configuration
source ~/.zshrc  # or source ~/.bashrc
```

## Verification

### System Check Script

Run this script to verify all dependencies are installed correctly:

```bash
#!/bin/bash
echo "=== Stax System Check ==="
echo ""

# Check Docker
if command -v docker >/dev/null 2>&1; then
    echo "✅ Docker installed: $(docker --version)"
    if docker ps >/dev/null 2>&1; then
        echo "✅ Docker is running"
    else
        echo "❌ Docker is not running - please start Docker"
    fi
else
    echo "❌ Docker not found - please install Docker"
fi

# Check DDEV
if command -v ddev >/dev/null 2>&1; then
    echo "✅ DDEV installed: $(ddev version | head -1)"
else
    echo "❌ DDEV not found - please install DDEV"
fi

# Check Stax
if command -v stax >/dev/null 2>&1; then
    echo "✅ Stax installed: $(stax --version)"
else
    echo "❌ Stax not found - please check installation"
fi

# Check WP-CLI (optional)
if command -v wp >/dev/null 2>&1; then
    echo "✅ WP-CLI installed: $(wp --version)"
else
    echo "⚠️  WP-CLI not found (optional but recommended)"
fi

# Check disk space
available=$(df -h . | awk 'NR==2 {print $4}')
echo ""
echo "💾 Available disk space: $available"

# Check memory
if [[ "$OSTYPE" == "darwin"* ]]; then
    memory=$(sysctl -n hw.memsize | awk '{print $1/1024/1024/1024 " GB"}')
else
    memory=$(free -h | awk 'NR==2 {print $2}')
fi
echo "💾 Total memory: $memory"

echo ""
echo "=== Check Complete ==="
```

Save this as `check-stax.sh`, make it executable (`chmod +x check-stax.sh`), and run it.

### Test Installation

Create a test project to verify everything works:

```bash
# Create a test site
stax init test-site
stax setup test-site --install-wp
stax start test-site

# Verify it's running
stax status

# Clean up
stax stop test-site
stax delete test-site --yes
```

## Troubleshooting Installation

### Common Issues

#### "Command not found: stax"
- Check if Stax is in your PATH: `echo $PATH`
- Find where Stax was installed: `find / -name stax 2>/dev/null`
- Add the directory to PATH in your shell configuration

#### "Cannot connect to Docker daemon"
- Ensure Docker Desktop is running (check for whale icon in menu bar)
- On Linux, ensure your user is in the docker group: `groups $USER`
- Try: `docker ps` to test Docker connection

#### "DDEV requires Docker"
- Docker must be running before DDEV can work
- Check Docker status: `docker info`
- Restart Docker if needed

#### Port conflicts
- Check if ports 80/443 are in use: `lsof -i :80` and `lsof -i :443`
- Stop conflicting services or configure DDEV to use different ports

## Next Steps

Congratulations! Stax is now installed and configured. Continue with:

1. **[User Guide](USER_GUIDE.md)** - Learn daily workflows and commands
2. **[WP Engine Setup](WPENGINE.md)** - Configure hosting integration
3. **[Your First Project](#)** - Create your first WordPress site

## Getting Help

- Run `stax --help` for built-in documentation
- Check [Troubleshooting](TROUBLESHOOTING.md) for common issues
- Report bugs at [GitHub Issues](https://github.com/Firecrown-Media/stax/issues)
- Contact support at dev@firecrown.com
# Installing Stax via Homebrew

## Prerequisites

- **macOS** 12.0 (Monterey) or later (Intel or Apple Silicon)
- **Homebrew** installed
- (Recommended) **Docker Desktop**
- (Recommended) **DDEV**

## Installation

### Step 1: Add the Firecrown Tap

First, add the Firecrown Homebrew tap to your system:

```bash
brew tap firecrown-media/tap
```

This tells Homebrew where to find the Stax formula.

### Step 2: Install Stax

Install Stax with a single command:

```bash
brew install stax
```

Homebrew will:
- Download the appropriate binary for your system (Intel or Apple Silicon)
- Verify the download with SHA256 checksum
- Install stax to your PATH
- Set up proper permissions

### Step 3: Verify Installation

Check that Stax is installed correctly:

```bash
stax --version
```

You should see output like:
```
stax version 1.0.0
```

Check that Stax is in your PATH:

```bash
which stax
```

Should show:
- Apple Silicon: `/opt/homebrew/bin/stax`
- Intel: `/usr/local/bin/stax`

## Install Recommended Dependencies

### Install DDEV

DDEV is required for running local WordPress environments:

```bash
brew install ddev/ddev/ddev
```

Configure DDEV's SSL certificates:

```bash
mkcert -install
```

### Install Docker Desktop

Stax requires Docker Desktop to run containers:

**Option 1: Direct Download**
1. Download from [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
2. Choose the version for your Mac (Intel or Apple Silicon)
3. Install and open Docker Desktop
4. Wait for Docker to start (green icon in menu bar)

**Option 2: Homebrew Cask**
```bash
brew install --cask docker
```

After installation, open Docker Desktop from Applications and let it complete its setup.

### Verify All Dependencies

Check that everything is installed:

```bash
# Check Stax
stax --version

# Check DDEV
ddev version

# Check Docker
docker --version
docker ps  # Should show empty table
```

## Configure Stax

After installation, configure your credentials:

```bash
stax setup
```

You'll be prompted for:
- WPEngine API username and password
- GitHub personal access token (optional)
- SSH key path for WPEngine

See [INSTALLATION.md](INSTALLATION.md#post-installation-setup) for detailed setup instructions.

## Updating Stax

### Update to Latest Version

Keep Stax up to date with Homebrew:

```bash
brew update
brew upgrade stax
```

### Check for Updates

See if a newer version is available:

```bash
brew outdated stax
```

If a newer version exists, it will be listed.

### Update All Homebrew Packages

Update all your Homebrew packages at once:

```bash
brew update
brew upgrade
```

## Uninstalling

### Remove Stax

To uninstall Stax:

```bash
brew uninstall stax
```

### Remove the Tap

To remove the Firecrown tap:

```bash
brew untap firecrown-media/tap
```

### Clean Up Stax Data (Optional)

Remove Stax configuration and data:

```bash
# Remove configuration directory
rm -rf ~/.stax

# Remove credentials from macOS Keychain
# (Must be done before uninstalling stax)
stax setup --remove
```

## Troubleshooting

### Installation Issues

#### "Error: Invalid formula"

**Problem**: Homebrew can't find the Stax formula.

**Solution**:
```bash
# Update Homebrew
brew update

# Re-add the tap
brew untap firecrown-media/tap
brew tap firecrown-media/tap

# Try installing again
brew install stax
```

#### "Error: Checksum mismatch"

**Problem**: Downloaded file doesn't match expected checksum.

**Solution**:
```bash
# Clean Homebrew cache
brew cleanup stax

# Try again
brew reinstall stax
```

#### "Permission denied"

**Problem**: Homebrew doesn't have permission to install.

**Solution**:
```bash
# Fix Homebrew permissions (Apple Silicon)
sudo chown -R $(whoami) /opt/homebrew/*

# Fix Homebrew permissions (Intel)
sudo chown -R $(whoami) /usr/local/*
```

#### "Command not found: brew"

**Problem**: Homebrew isn't in your PATH.

**Solution**:
```bash
# For Apple Silicon Macs:
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"

# For Intel Macs:
echo 'eval "$(/usr/local/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/usr/local/bin/brew shellenv)"

# Restart your terminal
```

### Runtime Issues

#### "Command not found: stax"

**Problem**: Stax is installed but not in PATH.

**Solution**:
```bash
# Check if Homebrew is in PATH
echo $PATH | grep -q homebrew && echo "Homebrew is in PATH" || echo "Homebrew not in PATH"

# Add Homebrew to PATH (see above)

# Verify stax location
ls -l /opt/homebrew/bin/stax  # Apple Silicon
ls -l /usr/local/bin/stax     # Intel

# If file exists, restart terminal
```

#### Wrong Version Installed

**Problem**: `stax --version` shows old version.

**Solution**:
```bash
# Uninstall old version
brew uninstall stax

# Clean cache
brew cleanup

# Update tap
brew update

# Reinstall
brew install stax

# Verify
stax --version
```

#### Multiple Versions Installed

**Problem**: Both Homebrew and manual installation exist.

**Solution**:
```bash
# Find all stax installations
which -a stax

# Remove manual installation
sudo rm /usr/local/bin/stax  # If it exists

# Keep only Homebrew version
brew reinstall stax
```

### Dependency Issues

#### "DDEV not found"

**Solution**:
```bash
brew install ddev/ddev/ddev
```

#### "Docker not running"

**Solution**:
1. Open Docker Desktop from Applications
2. Wait for Docker to start (green icon in menu bar)
3. Test: `docker ps`

#### "mkcert not installed"

**Solution**:
```bash
brew install mkcert
brew install nss  # For Firefox support
mkcert -install
```

## Advanced Usage

### Install Specific Version

Install a specific version of Stax:

```bash
# See available versions
brew info stax

# Install specific version (if available)
brew install stax@1.2.3
```

### Install from Local Tap

For testing or development:

```bash
# Clone the tap repository
git clone https://github.com/firecrown-media/homebrew-tap.git
cd homebrew-tap

# Install from local formula
brew install --build-from-source ./Formula/stax.rb
```

### View Formula Information

See details about the Stax formula:

```bash
brew info stax
```

Output shows:
- Version
- Homepage
- Description
- Dependencies
- Installation location
- Caveats

### Check Formula Health

Verify the formula is valid:

```bash
# Audit the formula
brew audit stax

# Test the formula
brew test stax
```

## Comparison with Other Installation Methods

| Method | Pros | Cons | Best For |
|--------|------|------|----------|
| **Homebrew** | - Easy updates<br>- Automatic dependency handling<br>- Trusted source | - Requires Homebrew<br>- macOS/Linux only | **Most users** |
| **Direct Download** | - No Homebrew needed<br>- Works anywhere | - Manual updates<br>- No dependency management | Quick testing |
| **Build from Source** | - Latest development version<br>- Full control | - Requires Go<br>- Manual updates<br>- More complex | Developers |

## Getting Help

If you encounter issues:

1. **Check stax diagnostics**:
   ```bash
   stax doctor
   ```

2. **Check Homebrew**:
   ```bash
   brew doctor
   ```

3. **View installation logs**:
   ```bash
   brew install --verbose stax
   ```

4. **Get support**:
   - [Stax Troubleshooting Guide](TROUBLESHOOTING.md)
   - [GitHub Issues](https://github.com/firecrown-media/stax/issues)
   - Contact Firecrown development team

## Next Steps

After installing Stax via Homebrew:

1. **Configure credentials**: `stax setup`
2. **Verify installation**: `stax doctor`
3. **Read the Quick Start**: [QUICK_START.md](QUICK_START.md)
4. **Explore commands**: `stax --help`

---

**Installation via Homebrew is the recommended method for all macOS users.**

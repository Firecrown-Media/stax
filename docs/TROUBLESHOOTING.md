# Troubleshooting Guide

Solutions for common issues when using Stax.

## Table of Contents
- [Installation Issues](#installation-issues)
- [Runtime Issues](#runtime-issues)
- [WP Engine Issues](#wp-engine-issues)
- [Performance Issues](#performance-issues)
- [Docker/DDEV Issues](#dockerddev-issues)
- [Database Issues](#database-issues)
- [Debugging Commands](#debugging-commands)
- [Getting Help](#getting-help)

## Installation Issues

### "Command not found: stax"

**Problem:** After installation, `stax` command is not recognized.

**Solutions:**

1. **Check if Stax is installed:**
```bash
# Find Stax binary
find / -name stax 2>/dev/null

# Common locations
ls -la /usr/local/bin/stax
ls -la ~/go/bin/stax
ls -la /opt/homebrew/bin/stax
```

2. **Add to PATH:**
```bash
# For Go installation
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc

# For Homebrew (M1 Mac)
echo 'export PATH=/opt/homebrew/bin:$PATH' >> ~/.zshrc
source ~/.zshrc
```

3. **Reinstall:**
```bash
# Via Homebrew
brew uninstall stax
brew install stax

# Via source
cd /path/to/stax
make clean
make install
```

### "Docker is not running"

**Problem:** Stax commands fail with Docker errors.

**Solutions:**

1. **Start Docker Desktop:**
```bash
# Mac
open /Applications/Docker.app

# Wait for whale icon in menu bar
```

2. **Check Docker status:**
```bash
docker version
docker ps
```

3. **Alternative Docker providers:**
```bash
# Using Colima
colima start --cpu 4 --memory 8

# Using OrbStack
open /Applications/OrbStack.app
```

### "DDEV not found"

**Problem:** DDEV is not installed or not in PATH.

**Solutions:**

```bash
# Install DDEV
brew install ddev/ddev/ddev

# Verify installation
ddev version

# Check PATH
which ddev
```

## Runtime Issues

### Port Conflicts

**Problem:** "Port 80 is already in use" or "Port 443 is already in use"

**Solutions:**

1. **Find conflicting process:**
```bash
# Mac/Linux
sudo lsof -i :80
sudo lsof -i :443

# Common culprits: Apache, nginx, MAMP, XAMPP
```

2. **Stop conflicting services:**
```bash
# Mac - Stop Apache
sudo apachectl stop

# Stop MAMP
# Open MAMP and click Stop Servers

# Stop other DDEV projects
ddev poweroff
```

3. **Use alternative ports:**
```bash
# Configure DDEV for different ports
ddev config --router-http-port=8080 --router-https-port=8443
```

### Container Startup Failures

**Problem:** "Failed to start containers"

**Solutions:**

1. **Check Docker resources:**
```bash
# Check if Docker is running
docker ps

# Check available resources
docker system df
docker stats --no-stream
```

2. **Clean up Docker:**
```bash
# Remove unused containers and images
docker system prune -a

# Remove all DDEV projects
ddev poweroff
ddev delete --all --omit-snapshot
```

3. **Restart Docker:**
```bash
# Mac
killall Docker && open /Applications/Docker.app

# Linux
sudo systemctl restart docker
```

### Project Won't Start

**Problem:** `stax start` fails or hangs.

**Solutions:**

1. **Check project status:**
```bash
stax status
ddev describe
```

2. **Reset project:**
```bash
stax stop project-name
stax delete project-name --yes
stax init project-name
```

3. **Check logs:**
```bash
ddev logs
ddev logs -s db
ddev logs -s web
```

## WP Engine Issues

### SSH Connection Failed

**Problem:** "Permission denied (publickey)" or connection timeout.

**Solutions:**

1. **Check SSH key:**
```bash
# List loaded keys
ssh-add -l

# Add key if missing
ssh-add ~/.ssh/wpengine_rsa

# Check key permissions
ls -la ~/.ssh/wpengine_rsa
# Should be 600
chmod 600 ~/.ssh/wpengine_rsa
```

2. **Test direct SSH:**
```bash
# Verbose connection test
ssh -vvv install@install.ssh.wpengine.net

# Check with specific key
ssh -i ~/.ssh/wpengine_rsa install@install.ssh.wpengine.net
```

3. **Verify WP Engine setup:**
- Log into my.wpengine.com
- Check SSH keys are added
- Verify install name is correct
- Check account has SSH access enabled

### Sync Failures

**Problem:** `stax wpe sync` fails during database or file transfer.

**Solutions:**

1. **Check credentials:**
```bash
# Verify environment variables
echo $WPE_USERNAME
echo $WPE_PASSWORD

# Test API access
stax wpe list
```

2. **Manual sync process:**
```bash
# SSH into WP Engine
ssh install@install.ssh.wpengine.net

# Export database manually
cd sites/install
wp db export backup.sql
exit

# Download database
scp install@install.ssh.wpengine.net:sites/install/backup.sql ./

# Import locally
ddev import-db --src=backup.sql
```

3. **Check WP Engine status:**
```bash
# Check if WP Engine is having issues
curl https://wpengine.com/support/status/
```

### Warp Terminal Issues

**Problem:** SSH hangs when using Warp terminal.

**Solution:** Use alternative terminal:
```bash
# Mac
open -a Terminal
# or
open -a iTerm

# Then run Stax commands
stax wpe sync install
```

## Performance Issues

### Slow Sync Operations

**Problem:** WP Engine sync takes too long.

**Solutions:**

1. **Skip unnecessary files:**
```bash
# Database only (fastest)
stax wpe sync install --skip-files

# Skip media files
stax wpe sync install --skip-media
```

2. **Use staging environment:**
```bash
# Staging databases are often smaller
stax wpe sync install --environment=staging
```

3. **Check network speed:**
```bash
# Test connection speed
ssh install@install.ssh.wpengine.net "dd if=/dev/zero bs=1M count=10" | dd of=/dev/null
```

### High Memory Usage

**Problem:** Docker using too much memory.

**Solutions:**

1. **Check resource usage:**
```bash
docker stats --no-stream
docker system df
```

2. **Limit DDEV resources:**
```bash
# Stop unused projects
stax poweroff

# Configure memory limits
ddev config --web-working-dir=/var/www/html --php-memory-limit=256M
```

3. **Clean up Docker:**
```bash
# Remove unused data
docker system prune -a --volumes
```

### Slow Site Performance

**Problem:** Local site runs slowly.

**Solutions:**

1. **Check container resources:**
```bash
ddev describe
docker stats
```

2. **Optimize database:**
```bash
stax wp db optimize
stax wp transient delete --all
```

3. **Disable unnecessary plugins:**
```bash
stax wp plugin deactivate --all
stax wp plugin activate essential-plugin-only
```

## Docker/DDEV Issues

### "Cannot connect to Docker daemon"

**Solutions:**

1. **Check Docker service:**
```bash
# Mac
docker version
# If failed, start Docker Desktop

# Linux
sudo systemctl status docker
sudo systemctl start docker
```

2. **Check user permissions (Linux):**
```bash
# Add user to docker group
sudo usermod -aG docker $USER
# Log out and back in
```

### DDEV Router Issues

**Problem:** "Router is not running" or "ddev-router container is not running"

**Solutions:**

```bash
# Restart router
ddev poweroff
ddev start

# Reset router ports
ddev config global --router-http-port=80 --router-https-port=443
```

### Container Disk Space

**Problem:** "No space left on device"

**Solutions:**

```bash
# Check disk usage
df -h
docker system df

# Clean everything
docker system prune -a --volumes
ddev delete --all --omit-snapshot

# Check Docker Desktop settings
# Increase disk image size in preferences
```

## Database Issues

### Import Failures

**Problem:** Database import fails or hangs.

**Solutions:**

1. **Check file size and format:**
```bash
# Check file
ls -lh database.sql
file database.sql
head -n 20 database.sql
```

2. **Import with progress:**
```bash
# Install pv for progress
brew install pv  # Mac
sudo apt install pv  # Linux

# Import with progress bar
pv database.sql | ddev mysql
```

3. **Split large files:**
```bash
# Split into smaller chunks
split -l 50000 database.sql chunk_

# Import chunks
for file in chunk_*; do
    ddev mysql < $file
done
```

### URL Rewrite Issues

**Problem:** Site redirects to production URL.

**Solutions:**

```bash
# Search and replace URLs
stax wp search-replace "production.com" "local.ddev.site" --all-tables

# Check site URL settings
stax wp option get siteurl
stax wp option get home

# Update if needed
stax wp option update siteurl "https://local.ddev.site"
stax wp option update home "https://local.ddev.site"
```

### Character Encoding Issues

**Problem:** Special characters appear broken.

**Solutions:**

```bash
# Check database charset
stax wp db query "SHOW VARIABLES LIKE 'character_set_%';"

# Convert database
stax wp db export backup.sql
iconv -f ISO-8859-1 -t UTF-8 backup.sql > converted.sql
stax wp db reset --yes
stax wp db import converted.sql
```

## Debugging Commands

### General Debugging

```bash
# Verbose output
stax --verbose [command]

# Check versions
stax --version
ddev version
docker version

# System information
uname -a
df -h
free -h  # Linux
vm_stat  # Mac
```

### DDEV Debugging

```bash
# Project information
ddev describe
ddev list

# Logs
ddev logs
ddev logs -f  # Follow
ddev logs -s web  # Web container only
ddev logs -s db  # Database only

# SSH into containers
ddev ssh  # Web container
ddev ssh -s db  # Database container

# Execute commands in container
ddev exec pwd
ddev exec php -v
```

### WordPress Debugging

```bash
# Enable debug mode
stax wp config set WP_DEBUG true --raw
stax wp config set WP_DEBUG_LOG true --raw
stax wp config set WP_DEBUG_DISPLAY false --raw

# Check debug log
stax wp eval 'echo WP_CONTENT_DIR . "/debug.log";'
tail -f wp-content/debug.log

# Database check
stax wp db check
stax wp db repair
```

### Network Debugging

```bash
# Check ports
netstat -an | grep -E ":(80|443|3306)"
lsof -i :80
lsof -i :443

# DNS resolution
nslookup local.ddev.site
ping local.ddev.site

# Check hosts file
cat /etc/hosts | grep ddev
```

## Getting Help

### Self-Help Resources

1. **Built-in help:**
```bash
stax --help
stax [command] --help
ddev help
```

2. **Check logs:**
```bash
# Stax logs (if verbose enabled)
stax --verbose [command] 2>&1 | tee stax.log

# DDEV logs
ddev logs > ddev.log

# Docker logs
docker logs [container-id]
```

3. **Version information for bug reports:**
```bash
# Collect system info
stax --version > debug-info.txt
ddev version >> debug-info.txt
docker version >> debug-info.txt
uname -a >> debug-info.txt
```

### Reporting Issues

When reporting issues, include:

1. **Environment details:**
   - Operating system and version
   - Stax version (`stax --version`)
   - DDEV version (`ddev version`)
   - Docker version (`docker version`)

2. **Steps to reproduce:**
   - Exact commands run
   - Expected behavior
   - Actual behavior

3. **Error messages:**
   - Full error output
   - Relevant log files
   - Screenshots if applicable

4. **What you've tried:**
   - Troubleshooting steps taken
   - Workarounds attempted

### Support Channels

- **GitHub Issues:** [github.com/Firecrown-Media/stax/issues](https://github.com/Firecrown-Media/stax/issues)
- **Email:** dev@firecrown.com
- **Documentation:** Check all guides in `/docs` directory

### Emergency Recovery

If everything is broken:

```bash
# Nuclear option - reset everything
ddev poweroff
docker system prune -a --volumes
brew uninstall stax ddev
brew install ddev/ddev/ddev stax

# Start fresh
stax init test-project
stax setup test-project --install-wp
stax start test-project
```

## Prevention Tips

1. **Regular maintenance:**
```bash
# Weekly cleanup
docker system prune
ddev cleanup
```

2. **Monitor resources:**
```bash
# Check before starting new projects
docker system df
df -h
```

3. **Keep software updated:**
```bash
brew upgrade stax ddev docker
```

4. **Backup before major changes:**
```bash
stax wp db export backup-$(date +%Y%m%d).sql
```
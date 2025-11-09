# Security Documentation - Stax CLI

## Table of Contents

- [Overview](#overview)
- [Security Features](#security-features)
- [For Users](#for-users)
- [For Developers](#for-developers)
- [Reporting Vulnerabilities](#reporting-vulnerabilities)
- [Security Best Practices](#security-best-practices)

---

## Overview

Stax handles sensitive information including:
- WPEngine API credentials
- SSH private keys
- Database credentials and dumps
- WordPress configuration data
- System-level access (/etc/hosts modification)

This document outlines security features, best practices, and guidelines for secure usage and development.

---

## Security Features

### Credential Management

**macOS Keychain Integration**
- All credentials stored in macOS Keychain
- OS-level encryption and access control
- No credentials in configuration files
- No credentials in environment variables
- Secure credential retrieval with user authentication

**Supported Credentials:**
- WPEngine API credentials (username/password)
- GitHub personal access tokens
- SSH private keys

**Storage Location:**
```
Service Name: com.firecrown.stax.wpengine
Account: [install-name]
```

### Network Security

**HTTPS Enforcement**
- All API communications use HTTPS
- WPEngine API: `https://api.wpengineapi.com/v1`
- No plaintext credential transmission

**SSH Security**
- Public key authentication only
- Encrypted data transfer for database dumps and files
- Connection timeout: 30 seconds
- ⚠️ **Note:** Host key verification currently disabled (see [Known Limitations](#known-limitations))

### File System Security

**Temporary Files**
- Created with restrictive permissions (0600)
- SSH keys stored temporarily only during operations
- Automatic cleanup after use
- Secure directory: `/tmp` with stax-specific naming

**System Modifications**
- /etc/hosts modifications require sudo
- Automatic backup before modification
- Restore capability from backups
- Marker-based section management

### Data Protection

**Database Dumps**
- Never stored in git repositories (.gitignore included)
- Local storage only during sync operations
- Option to exclude sensitive tables
- Compression support for reduced exposure time

**Configuration Files**
- No credentials in `.stax.yml`
- Project configurations safe to commit
- Global config in `~/.stax/config.yml` (no credentials)

---

## For Users

### Initial Setup Security

**1. Configure Credentials Securely**

```bash
# Set up WPEngine credentials (stored in keychain)
stax setup

# Verify credentials are in keychain, not files
security find-generic-password -s "com.firecrown.stax.wpengine"
```

**2. Protect Your Configuration**

```bash
# Ensure .stax.yml is in project root (safe to commit)
# Never commit files containing credentials

# Example .gitignore (already included):
.env
*.sql
*.sql.gz
database-*.sql
```

**3. Secure Your SSH Keys**

```bash
# Generate SSH key if needed
ssh-keygen -t ed25519 -C "stax-wpengine"

# Set restrictive permissions
chmod 600 ~/.ssh/id_ed25519

# Add to WPEngine account via portal
# Store in keychain using stax setup
```

### Daily Usage Security

**Database Operations**

```bash
# Pull database (will be stored locally temporarily)
stax db pull

# IMPORTANT: Never commit database dumps
# They contain sensitive user data

# Delete database dump after import
rm database-*.sql*
```

**File Synchronization**

```bash
# Sync files from WPEngine
stax sync files

# Review synced files before committing
# Exclude sensitive files:
# - wp-config.php (if present)
# - .htaccess with sensitive rules
# - Any files with credentials
```

**Credential Management**

```bash
# View stored credentials (requires keychain auth)
stax config credentials list

# Update credentials
stax setup --update-credentials

# Remove credentials when done with project
stax config credentials delete
```

### Security Checklist for Users

- [ ] Credentials stored in keychain only
- [ ] Never commit `.env` files
- [ ] Database dumps not in git repository
- [ ] SSH keys have 0600 permissions
- [ ] `.gitignore` includes sensitive files
- [ ] Remove database dumps after use
- [ ] Use strong passwords for WPEngine account
- [ ] Enable 2FA on WPEngine account
- [ ] Regularly rotate API credentials
- [ ] Review file permissions before commit

### Team Security Considerations

**Sharing Projects**

✅ **Safe to share:**
- `.stax.yml` configuration file
- Project structure and code
- Build scripts and configurations

❌ **Never share:**
- Keychain credentials (each team member sets up their own)
- Database dumps
- SSH private keys
- API credentials

**Team Setup Process**

1. Each team member runs `stax setup` independently
2. Each member uses their own WPEngine credentials
3. Share only project configuration (`.stax.yml`)
4. Document environment-specific settings separately

### API Key Rotation

**When to Rotate:**
- Every 90 days (recommended)
- After team member departure
- If credentials potentially compromised
- After security incident

**How to Rotate:**

```bash
# 1. Generate new API credentials in WPEngine portal
# 2. Update keychain
stax setup --update-credentials

# 3. Test new credentials
stax status

# 4. Delete old credentials from WPEngine portal
```

### Database Dump Security

**Best Practices:**

```bash
# 1. Pull database
stax db pull

# 2. Import immediately
stax db import database-latest.sql

# 3. DELETE source file
rm database-latest.sql*

# 4. Never upload database dumps to:
#    - Git repositories
#    - Cloud storage
#    - Slack/email
#    - Public locations
```

**If Database Dump is Accidentally Committed:**

```bash
# 1. Remove from git history immediately
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch database-*.sql*" \
  --prune-empty --tag-name-filter cat -- --all

# 2. Force push (coordinate with team)
git push origin --force --all

# 3. Rotate all credentials in dump:
#    - WordPress admin passwords
#    - API keys in database
#    - Database credentials
```

### Known Limitations

⚠️ **Current Security Limitations:**

1. **SSH Host Key Verification Disabled**
   - Vulnerable to MITM attacks
   - Use trusted networks only
   - Fix planned for v1.1

2. **No Built-in Encryption for Database Dumps**
   - Store dumps only temporarily
   - Use encrypted disk for extra protection
   - Delete immediately after use

3. **Sudo Required for /etc/hosts**
   - Cannot fully automate on first run
   - Use with caution on shared systems

---

## For Developers

### Secure Coding Guidelines

**1. Credential Handling**

```go
// ✅ CORRECT: Use keychain
creds, err := credentials.GetWPEngineCredentials(install)
if err != nil {
    return err
}

// ❌ WRONG: Never hardcode
const apiKey = "wpengine_api_key_abc123" // NEVER DO THIS

// ❌ WRONG: Never log credentials
log.Printf("API Key: %s", apiKey) // NEVER DO THIS

// ✅ CORRECT: Sanitize errors
return fmt.Errorf("authentication failed") // Don't include credentials
```

**2. Input Validation**

```go
// ✅ CORRECT: Validate all user input
func processUserInput(input string) error {
    if !isValidInput(input) {
        return fmt.Errorf("invalid input format")
    }
    // ... process
}

// ❌ WRONG: Trust user input
func processUserInput(input string) error {
    cmd := fmt.Sprintf("wp %s", input) // VULNERABLE TO INJECTION
    exec.Command("bash", "-c", cmd)
}

// ✅ CORRECT: Use argument arrays
func processUserInput(input string) error {
    // Validate input first
    if !isValidInput(input) {
        return fmt.Errorf("invalid input")
    }

    // Use argument array, not shell string
    cmd := exec.Command("wp", sanitize(input))
    return cmd.Run()
}
```

**3. Command Execution**

```go
// ✅ CORRECT: Use exec.Command with separate arguments
cmd := exec.Command("rsync", "-avz", source, destination)

// ❌ WRONG: Shell expansion with user input
cmd := exec.Command("bash", "-c", fmt.Sprintf("rsync -avz %s %s", userInput1, userInput2))

// ✅ CORRECT: Sanitize before command construction
func sanitizePath(path string) (string, error) {
    // Remove dangerous characters
    if strings.ContainsAny(path, ";|&$`<>()") {
        return "", fmt.Errorf("invalid characters in path")
    }

    // Prevent path traversal
    if strings.Contains(path, "../") {
        return "", fmt.Errorf("path traversal detected")
    }

    return path, nil
}
```

**4. Error Handling**

```go
// ✅ CORRECT: Sanitize error messages
func sanitizeError(err error) error {
    msg := err.Error()

    // Remove credential patterns
    patterns := map[string]string{
        `password=\S+`:       "password=***",
        `token=\S+`:          "token=***",
        `Authorization: \S+`: "Authorization: ***",
    }

    for pattern, replacement := range patterns {
        re := regexp.MustCompile(pattern)
        msg = re.ReplaceAllString(msg, replacement)
    }

    return errors.New(msg)
}

// ❌ WRONG: Expose credentials in errors
return fmt.Errorf("API request failed with key %s", apiKey)
```

**5. File Operations**

```go
// ✅ CORRECT: Create temp files with secure permissions
func createSecureTempFile(data []byte) (string, error) {
    tmpFile, err := os.OpenFile(
        filepath.Join(os.TempDir(), "stax-"+randomString()),
        os.O_RDWR|os.O_CREATE|os.O_EXCL,
        0600, // Secure permissions at creation
    )
    if err != nil {
        return "", err
    }
    defer tmpFile.Close()

    if _, err := tmpFile.Write(data); err != nil {
        os.Remove(tmpFile.Name())
        return "", err
    }

    return tmpFile.Name(), nil
}

// ❌ WRONG: Insecure temp file creation
func createTempFile(data []byte) (string, error) {
    tmpFile, err := os.CreateTemp("", "stax-*") // Created with 0666
    // ... race condition before chmod
}
```

**6. Path Validation**

```go
// ✅ CORRECT: Validate paths
func validatePath(path string, allowedDir string) error {
    // Convert to absolute path
    absPath, err := filepath.Abs(path)
    if err != nil {
        return err
    }

    // Ensure within allowed directory
    if !strings.HasPrefix(absPath, allowedDir) {
        return fmt.Errorf("path outside allowed directory")
    }

    // Check for symlinks
    if info, err := os.Lstat(absPath); err == nil {
        if info.Mode()&os.ModeSymlink != 0 {
            return fmt.Errorf("symlinks not allowed")
        }
    }

    return nil
}
```

### Security Testing Requirements

**Required Tests:**

```go
// Test command injection prevention
func TestCommandInjectionPrevention(t *testing.T) {
    malicious := []string{
        "; rm -rf /",
        "| cat /etc/passwd",
        "$(whoami)",
        "`id`",
    }

    for _, input := range malicious {
        err := processInput(input)
        if err == nil {
            t.Errorf("Failed to prevent injection: %s", input)
        }
    }
}

// Test path traversal prevention
func TestPathTraversalPrevention(t *testing.T) {
    malicious := []string{
        "../../../etc/passwd",
        "..\\..\\..\\windows\\system32",
        "/etc/passwd",
    }

    for _, input := range malicious {
        err := validatePath(input)
        if err == nil {
            t.Errorf("Failed to prevent traversal: %s", input)
        }
    }
}

// Test credential sanitization
func TestCredentialSanitization(t *testing.T) {
    errWithCreds := fmt.Errorf("Failed with password=secret123")
    sanitized := sanitizeError(errWithCreds)

    if strings.Contains(sanitized.Error(), "secret123") {
        t.Error("Failed to sanitize credential from error")
    }
}
```

### Code Review Checklist

**Security Review Checklist:**

- [ ] No hardcoded credentials or secrets
- [ ] All user input validated
- [ ] No shell expansion with user input
- [ ] exec.Command uses argument arrays
- [ ] Errors don't leak sensitive data
- [ ] File operations use secure permissions
- [ ] No path traversal vulnerabilities
- [ ] SQL queries use safe construction
- [ ] Credentials retrieved from keychain only
- [ ] Temporary files cleaned up properly
- [ ] Network calls use HTTPS
- [ ] Certificate validation enabled
- [ ] Security tests added for new features
- [ ] Dependencies scanned for vulnerabilities

### Dependency Management

**Adding Dependencies:**

```bash
# 1. Research dependency security
# 2. Check for known vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# 3. Add dependency
go get package@version

# 4. Verify in go.sum
git diff go.sum

# 5. Document why dependency is needed
```

**Updating Dependencies:**

```bash
# Check for vulnerabilities
govulncheck ./...

# Update specific package
go get package@latest

# Update all dependencies
go get -u ./...

# Test thoroughly
go test ./...
```

### Pre-commit Hooks

**Recommended pre-commit hook:**

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Run security checks
echo "Running security checks..."

# Check for secrets
if git diff --cached | grep -i "password\|secret\|api_key" > /dev/null; then
    echo "ERROR: Potential secret in commit"
    exit 1
fi

# Run vulnerability check
govulncheck ./...
if [ $? -ne 0 ]; then
    echo "ERROR: Vulnerabilities detected"
    exit 1
fi

# Run static analysis
gosec ./...
if [ $? -ne 0 ]; then
    echo "ERROR: Security issues detected"
    exit 1
fi

echo "Security checks passed"
```

---

## Reporting Vulnerabilities

### How to Report

If you discover a security vulnerability in Stax, please report it responsibly:

**Do NOT:**
- Create public GitHub issues for security vulnerabilities
- Disclose vulnerability details publicly
- Exploit the vulnerability

**Do:**
1. Email security report to: security@firecrown-media.com
2. Include:
   - Description of vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)
3. Allow 48 hours for initial response

### Security Report Template

```
Subject: [SECURITY] Vulnerability in Stax CLI

Component: [e.g., SSH client, API client]
Severity: [Critical/High/Medium/Low]

Description:
[Detailed description of the vulnerability]

Steps to Reproduce:
1. [Step 1]
2. [Step 2]
3. [etc]

Impact:
[What an attacker could do]

Suggested Fix:
[Optional - your recommendation]

Additional Information:
[Any other relevant details]
```

### Response Timeline

- **Initial Response:** Within 48 hours
- **Vulnerability Assessment:** Within 1 week
- **Fix Development:** Based on severity
  - Critical: 1-2 weeks
  - High: 2-4 weeks
  - Medium: 4-8 weeks
  - Low: Next release cycle
- **Public Disclosure:** After fix is released

### Disclosure Policy

**Coordinated Disclosure:**
1. Vulnerability reported privately
2. Fix developed and tested
3. Security advisory published
4. Fix released to users
5. Public disclosure after users have time to update (typically 30 days)

**Security Advisories:**
Published at: https://github.com/firecrown-media/stax/security/advisories

---

## Security Best Practices

### Development Environment

**Secure Development Setup:**

```bash
# 1. Use separate credentials for dev/staging/production
stax setup --environment development

# 2. Never use production credentials locally
# 3. Keep development dependencies updated
go get -u ./...

# 4. Enable security scanning in IDE
# Install: gopls, staticcheck, gosec
```

### Production Deployment

**Security Hardening:**

1. **Credential Isolation**
   - Use separate credentials per environment
   - Rotate credentials regularly
   - Limit credential scope (least privilege)

2. **Access Control**
   - Restrict who can run stax commands
   - Use separate user accounts for automation
   - Audit credential access

3. **Monitoring**
   - Log security-relevant events
   - Monitor for unusual activity
   - Set up alerts for credential access

### Incident Response

**If Credentials Compromised:**

1. **Immediate Actions:**
   ```bash
   # 1. Delete compromised credentials from keychain
   stax config credentials delete

   # 2. Generate new credentials in WPEngine portal
   # 3. Update keychain with new credentials
   stax setup --update-credentials

   # 4. Revoke old credentials in WPEngine
   ```

2. **Investigation:**
   - Review recent activity logs
   - Check for unauthorized database access
   - Review file modifications
   - Audit user account changes

3. **Recovery:**
   - Reset all potentially affected credentials
   - Review and update access controls
   - Document incident and lessons learned

### Security Training

**Recommended Training:**

- OWASP Top 10 awareness
- Secure coding practices for Go
- Git security and secret management
- SSH key management
- Incident response procedures

---

## Additional Resources

**Security Tools:**

- `govulncheck` - Go vulnerability scanner
- `gosec` - Go security checker
- `git-secrets` - Prevent committing secrets
- `truffleHog` - Scan for secrets in git history

**Documentation:**

- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)
- [Go Security Best Practices](https://github.com/golang/go/wiki/Security)
- [WPEngine Security](https://wpengine.com/support/security/)

**Contact:**

- Security Issues: security@firecrown-media.com
- General Support: support@firecrown-media.com
- Documentation: https://github.com/firecrown-media/stax

---

**Last Updated:** November 8, 2025
**Next Review:** February 8, 2026

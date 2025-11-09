# Security Quick Reference Guide

**Quick access to security guidelines for Stax developers**

---

## For Code Reviewers

### Security Checklist for PRs

Before approving any PR, verify:

```
Credentials & Secrets
[ ] No hardcoded credentials
[ ] No credentials in logs/errors
[ ] Keychain used for credential storage
[ ] No sensitive data in test files

Input Validation
[ ] All user input validated
[ ] No shell metacharacters in commands
[ ] Paths validated for traversal
[ ] Config validation present

Command Execution
[ ] exec.Command uses argument arrays (not shell strings)
[ ] No user input in fmt.Sprintf for commands
[ ] SSH commands sanitized
[ ] No shell expansion with user data

Error Handling
[ ] Errors don't leak credentials
[ ] Stack traces sanitized
[ ] Error messages user-friendly

File Operations
[ ] Secure permissions (0600 for keys, 0644 for configs)
[ ] Temp files cleaned up
[ ] No race conditions

Tests
[ ] Security tests added for new code
[ ] Edge cases covered
[ ] Malicious input tested
```

---

## Common Vulnerabilities & How to Avoid Them

### ❌ Command Injection

**WRONG:**
```go
cmd := fmt.Sprintf("wp %s", userInput)
exec.Command("bash", "-c", cmd)
```

**RIGHT:**
```go
// Validate input first
if !isValidInput(userInput) {
    return fmt.Errorf("invalid input")
}

// Use argument array
cmd := exec.Command("wp", userInput)
```

---

### ❌ Path Traversal

**WRONG:**
```go
func downloadFile(path string) {
    // No validation - dangerous!
    cmd := fmt.Sprintf("cat %s", path)
}
```

**RIGHT:**
```go
func downloadFile(path string) error {
    // Validate path
    if strings.Contains(path, "../") {
        return fmt.Errorf("invalid path")
    }

    // Verify within allowed directory
    absPath, _ := filepath.Abs(path)
    if !strings.HasPrefix(absPath, allowedDir) {
        return fmt.Errorf("path outside allowed directory")
    }

    // Now safe to use
    cmd := fmt.Sprintf("cat %s", path)
}
```

---

### ❌ Credential Leakage

**WRONG:**
```go
log.Printf("Connecting with API key: %s", apiKey)

return fmt.Errorf("auth failed with key %s", apiKey)
```

**RIGHT:**
```go
log.Printf("Connecting to API")

return fmt.Errorf("authentication failed")
```

---

### ❌ Insecure Temp Files

**WRONG:**
```go
// Race condition - created with 0666, then chmod
tmpFile, _ := os.CreateTemp("", "key-*")
tmpFile.WriteString(privateKey)
os.Chmod(tmpFile.Name(), 0600)  // TOO LATE!
```

**RIGHT:**
```go
// Secure from creation
tmpFile, err := os.OpenFile(
    filepath.Join(os.TempDir(), "key-"+randomString()),
    os.O_RDWR|os.O_CREATE|os.O_EXCL,
    0600,  // Secure permissions at creation
)
defer os.Remove(tmpFile.Name())
tmpFile.WriteString(privateKey)
```

---

### ❌ SQL Injection

**WRONG:**
```go
query := fmt.Sprintf("SELECT * FROM %s_options", userPrefix)
```

**RIGHT:**
```go
// Validate prefix first
if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(userPrefix) {
    return fmt.Errorf("invalid table prefix")
}

query := fmt.Sprintf("SELECT * FROM %s_options", userPrefix)
```

---

## Code Snippets Library

### Input Validation

```go
// Validate project name
func ValidateProjectName(name string) error {
    if name == "" {
        return fmt.Errorf("project name required")
    }

    // Alphanumeric, hyphens, underscores only
    if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
        return fmt.Errorf("invalid characters in project name")
    }

    if len(name) > 64 {
        return fmt.Errorf("project name too long (max 64 chars)")
    }

    return nil
}
```

```go
// Validate hostname
func ValidateHostname(hostname string) error {
    if hostname == "" || len(hostname) > 253 {
        return fmt.Errorf("invalid hostname length")
    }

    // RFC 1123
    pattern := `^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*` +
                `[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`

    if !regexp.MustCompile(pattern).MatchString(hostname) {
        return fmt.Errorf("invalid hostname format")
    }

    return nil
}
```

```go
// Validate file path
func ValidatePath(path string, allowedDir string) error {
    // Prevent traversal
    if strings.Contains(path, "../") || strings.Contains(path, "..\\") {
        return fmt.Errorf("path traversal detected")
    }

    // Convert to absolute path
    absPath, err := filepath.Abs(path)
    if err != nil {
        return err
    }

    // Ensure within allowed directory
    if !strings.HasPrefix(absPath, allowedDir) {
        return fmt.Errorf("path outside allowed directory")
    }

    // Check if symlink
    info, err := os.Lstat(absPath)
    if err == nil && info.Mode()&os.ModeSymlink != 0 {
        return fmt.Errorf("symlinks not allowed")
    }

    return nil
}
```

### Sanitization

```go
// Sanitize error messages
func SanitizeError(err error) error {
    if err == nil {
        return nil
    }

    msg := err.Error()

    // Remove credential patterns
    patterns := map[string]string{
        `password=\S+`:       "password=***",
        `token=\S+`:          "token=***",
        `api_key=\S+`:        "api_key=***",
        `Authorization: \S+`: "Authorization: ***",
    }

    for pattern, replacement := range patterns {
        re := regexp.MustCompile(pattern)
        msg = re.ReplaceAllString(msg, replacement)
    }

    return errors.New(msg)
}
```

```go
// Sanitize for shell (whitelist approach)
func SanitizeForShell(input string) (string, error) {
    // Allow only safe characters
    if !regexp.MustCompile(`^[a-zA-Z0-9/_.-]+$`).MatchString(input) {
        return "", fmt.Errorf("invalid characters in input")
    }

    return input, nil
}
```

### Secure File Operations

```go
// Create secure temp file
func CreateSecureTempFile(prefix string, data []byte) (string, error) {
    // Generate random filename
    random := make([]byte, 16)
    rand.Read(random)
    filename := fmt.Sprintf("%s-%x", prefix, random)

    tmpPath := filepath.Join(os.TempDir(), filename)

    // Create with secure permissions
    file, err := os.OpenFile(
        tmpPath,
        os.O_RDWR|os.O_CREATE|os.O_EXCL,
        0600,
    )
    if err != nil {
        return "", err
    }
    defer file.Close()

    if _, err := file.Write(data); err != nil {
        os.Remove(tmpPath)
        return "", err
    }

    return tmpPath, nil
}
```

```go
// Secure cleanup
func SecureCleanup(path string) error {
    // Ensure file exists and get size
    info, err := os.Stat(path)
    if err != nil {
        return err
    }

    // Open for writing
    file, err := os.OpenFile(path, os.O_WRONLY, 0)
    if err != nil {
        return err
    }
    defer file.Close()

    // Overwrite with zeros
    zeros := make([]byte, info.Size())
    if _, err := file.Write(zeros); err != nil {
        return err
    }

    file.Sync()

    // Delete
    return os.Remove(path)
}
```

### Credential Handling

```go
// Store credentials
func StoreCredentials(install string, creds *Credentials) error {
    // ALWAYS use keychain
    return credentials.SetWPEngineCredentials(install, creds)

    // NEVER do this:
    // config.APIKey = creds.APIKey  ❌
    // os.Setenv("API_KEY", key)     ❌
}
```

```go
// Retrieve credentials
func GetCredentials(install string) (*Credentials, error) {
    creds, err := credentials.GetWPEngineCredentials(install)
    if err != nil {
        // Don't leak credential details in error
        return nil, fmt.Errorf("failed to retrieve credentials")
    }
    return creds, nil
}
```

---

## Testing Patterns

### Security Test Template

```go
func TestCommandInjectionPrevention(t *testing.T) {
    maliciousInputs := []string{
        "; rm -rf /",
        "| cat /etc/passwd",
        "$(whoami)",
        "`id`",
        "&& echo pwned",
        "' OR '1'='1",
    }

    for _, input := range maliciousInputs {
        t.Run(input, func(t *testing.T) {
            err := processUserInput(input)
            if err == nil {
                t.Errorf("Failed to prevent injection: %s", input)
            }
        })
    }
}
```

```go
func TestPathTraversalPrevention(t *testing.T) {
    testCases := []struct {
        path        string
        shouldFail  bool
        description string
    }{
        {"../../../etc/passwd", true, "Unix path traversal"},
        {"..\\..\\..\\windows\\system32", true, "Windows path traversal"},
        {"/etc/passwd", true, "Absolute path"},
        {"valid/path/file.txt", false, "Valid relative path"},
        {"./file.txt", false, "Current directory"},
    }

    for _, tc := range testCases {
        t.Run(tc.description, func(t *testing.T) {
            err := validatePath(tc.path, "/allowed/dir")
            failed := (err != nil)

            if failed != tc.shouldFail {
                t.Errorf("Path %s: expected fail=%v, got fail=%v",
                    tc.path, tc.shouldFail, failed)
            }
        })
    }
}
```

```go
func TestCredentialSanitization(t *testing.T) {
    testCases := []struct {
        input    error
        shouldContain string
        shouldNotContain string
    }{
        {
            fmt.Errorf("Failed with password=secret123"),
            "password=***",
            "secret123",
        },
        {
            fmt.Errorf("API error with token=abc123"),
            "token=***",
            "abc123",
        },
    }

    for _, tc := range testCases {
        sanitized := SanitizeError(tc.input)
        msg := sanitized.Error()

        if !strings.Contains(msg, tc.shouldContain) {
            t.Errorf("Expected %q in error message", tc.shouldContain)
        }

        if strings.Contains(msg, tc.shouldNotContain) {
            t.Errorf("Credential leaked in error message")
        }
    }
}
```

---

## Git Hooks

### Pre-commit Hook

Save as `.git/hooks/pre-commit`:

```bash
#!/bin/bash

echo "Running security checks..."

# Check for potential secrets
if git diff --cached | grep -Ei "password\s*=|api_key\s*=|secret\s*=|token\s*=" > /dev/null; then
    echo "❌ ERROR: Potential secret in staged files"
    echo "Review your changes and use keychain instead"
    exit 1
fi

# Check for SQL files (database dumps)
if git diff --cached --name-only | grep -E "\.sql(\.gz)?$" > /dev/null; then
    echo "❌ ERROR: SQL dump file in commit"
    echo "Database dumps should never be committed"
    exit 1
fi

# Check for .env files
if git diff --cached --name-only | grep "\.env" > /dev/null; then
    echo "❌ ERROR: .env file in commit"
    exit 1
fi

# Run go fmt
echo "Running go fmt..."
go fmt ./...

# Run security scan (if available)
if command -v gosec &> /dev/null; then
    echo "Running gosec..."
    gosec -quiet ./... || exit 1
fi

echo "✅ Security checks passed"
```

Make executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

## Emergency Response

### If Credentials Compromised

```bash
# 1. IMMEDIATELY delete from keychain
security delete-generic-password -s "com.firecrown.stax.wpengine"

# 2. Generate new credentials in WPEngine portal

# 3. Store new credentials
stax setup --update-credentials

# 4. Verify old credentials deleted from WPEngine

# 5. Document incident
```

### If Secrets Committed to Git

```bash
# 1. DON'T push if local only!

# 2. If already pushed, remove from history
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch path/to/secret" \
  --prune-empty --tag-name-filter cat -- --all

# 3. Force push (COORDINATE WITH TEAM FIRST)
git push origin --force --all

# 4. Rotate ALL credentials in the file

# 5. Notify security team
```

---

## Quick Links

**Full Documentation:**
- [Complete Security Audit](SECURITY_AUDIT.md)
- [Security Best Practices](SECURITY.md)
- [Pre-Release Checklist](SECURITY_CHECKLIST.md)
- [Scan Results](SECURITY_SCAN_RESULTS.md)

**External Resources:**
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Guidelines](https://github.com/golang/go/wiki/Security)
- [CWE Database](https://cwe.mitre.org/)

**Tools:**
- govulncheck: `go install golang.org/x/vuln/cmd/govulncheck@latest`
- gosec: `go install github.com/securego/gosec/v2/cmd/gosec@latest`
- staticcheck: `go install honnef.co/go/tools/cmd/staticcheck@latest`

**Contacts:**
- Security Email: security@firecrown-media.com
- Security Slack: #security
- Incident Response: [On-call rotation]

---

## Quick Command Reference

```bash
# Run all security checks
make security-check

# Or manually:
go build ./...                    # Verify build
govulncheck ./...                 # Check vulnerabilities
gosec ./...                       # Security scan
staticcheck ./...                 # Static analysis
go test -tags=security ./...      # Security tests

# Check for secrets in git
git log -p | grep -i "password\|secret\|api_key"

# Update dependencies
go get -u ./...
go mod tidy

# Verify dependencies
go mod verify
```

---

**Keep This Handy:** Bookmark this page for quick security reference during development.

**Last Updated:** November 8, 2025
**Version:** 1.0

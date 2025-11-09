# Security Checklist - Stax CLI

## Pre-Release Security Checklist

This checklist must be completed before each release. Check off each item and provide evidence/notes where applicable.

**Release Version:** _____________
**Date:** _____________
**Reviewed By:** _____________

---

## Code Security

### Credentials & Secrets

- [ ] No hardcoded credentials in source code
  - [ ] Searched codebase for "password", "api_key", "secret", "token"
  - [ ] Verified all credentials use keychain
  - [ ] No credentials in test files

- [ ] No credentials in configuration files
  - [ ] `.stax.yml` contains no sensitive data
  - [ ] Example configs are sanitized
  - [ ] Template files contain no secrets

- [ ] Credentials never logged
  - [ ] Reviewed all log statements
  - [ ] Error messages sanitized
  - [ ] Debug output doesn't expose secrets

- [ ] No credentials in error messages
  - [ ] Tested error paths
  - [ ] Verified error sanitization
  - [ ] Stack traces reviewed

**Evidence/Notes:**
```
Grep results for credential patterns:
$ grep -r "password\|api_key\|secret" --include="*.go" .

Tools used:
- git-secrets
- gitleaks
```

---

## Network Security

### HTTPS & TLS

- [ ] All API calls use HTTPS
  - [ ] WPEngine API uses HTTPS
  - [ ] No HTTP fallback
  - [ ] URLs hardcoded as HTTPS

- [ ] Certificate validation enabled
  - [ ] No InsecureSkipVerify in production code
  - [ ] TLS configuration reviewed
  - [ ] Minimum TLS version: 1.2

- [ ] HTTP client configured securely
  - [ ] Timeout set appropriately
  - [ ] Redirect policy defined
  - [ ] User-Agent set

**Evidence/Notes:**
```
Verified API calls:
- WPEngine API: https://api.wpengineapi.com/v1

TLS Configuration:
- Minimum version: TLS 1.2
- Certificate validation: Enabled
```

### SSH Security

- [ ] SSH host key verification status documented
  - [ ] Known limitation documented if disabled
  - [ ] Mitigation measures in place
  - [ ] Roadmap for fix defined

- [ ] SSH keys stored securely
  - [ ] Keychain integration verified
  - [ ] Temporary files have 0600 permissions
  - [ ] Keys cleaned up after use

- [ ] SSH configuration hardened
  - [ ] Timeout configured
  - [ ] Key size requirements documented
  - [ ] Supported algorithms documented

**Evidence/Notes:**
```
Current status:
- Host key verification: Disabled (documented limitation)
- Key storage: Keychain
- Temp file permissions: 0600
```

---

## Input Validation

### User Input

- [ ] All user inputs validated
  - [ ] Project names validated
  - [ ] Hostnames validated
  - [ ] URLs validated
  - [ ] File paths validated

- [ ] No command injection vulnerabilities
  - [ ] exec.Command uses argument arrays
  - [ ] No shell expansion with user input
  - [ ] Shell metacharacters rejected or escaped

- [ ] No path traversal vulnerabilities
  - [ ] Paths validated against traversal
  - [ ] No "../" in user-provided paths
  - [ ] Symlinks handled securely

- [ ] Configuration validated
  - [ ] Schema validation present
  - [ ] Required fields checked
  - [ ] Format validation applied
  - [ ] Range checks implemented

**Evidence/Notes:**
```
Validation functions tested:
- ValidateProjectName()
- ValidateHostname()
- ValidatePath()

Test results:
$ go test -v ./pkg/config/...
$ go test -v ./pkg/system/...
```

### Command Execution

- [ ] No shell injection in commands
  - [ ] Reviewed all exec.Command calls
  - [ ] Verified argument separation
  - [ ] Tested with malicious inputs

- [ ] Command arguments sanitized
  - [ ] Special characters handled
  - [ ] Whitelist validation where appropriate
  - [ ] Length limits enforced

- [ ] No SQL injection vulnerabilities
  - [ ] Reviewed database queries
  - [ ] WP-CLI commands constructed safely
  - [ ] Table/column names validated

**Evidence/Notes:**
```
exec.Command usage reviewed:
$ grep -r "exec.Command" --include="*.go" .

Findings: [document any issues found]
```

---

## File System Security

### File Operations

- [ ] File permissions set correctly
  - [ ] Config files: 0644
  - [ ] SSH keys: 0600
  - [ ] Temporary files: 0600
  - [ ] Directories: 0755

- [ ] Temporary files handled securely
  - [ ] Created with secure permissions
  - [ ] Cleaned up properly
  - [ ] Unique names generated
  - [ ] Race conditions prevented

- [ ] Sensitive files deleted securely
  - [ ] SSH keys wiped before deletion
  - [ ] Database dumps removed after use
  - [ ] No sensitive data in temp directories

- [ ] Path operations validated
  - [ ] Absolute paths used where appropriate
  - [ ] Symlinks detected and handled
  - [ ] Directory traversal prevented

**Evidence/Notes:**
```
File permission audit:
$ find . -name "*.go" -exec grep -l "os.Create\|os.OpenFile" {} \;

Verified permissions in:
- pkg/wpengine/files.go
- pkg/credentials/keychain.go
```

### System Modifications

- [ ] /etc/hosts modifications secure
  - [ ] Backup created before modification
  - [ ] Restore capability tested
  - [ ] Input validation applied
  - [ ] Sudo requirements documented

- [ ] File backup process tested
  - [ ] Backups created successfully
  - [ ] Restore process works
  - [ ] Backup cleanup implemented

**Evidence/Notes:**
```
Tested hosts file modification:
- Backup created: ✓
- Modification applied: ✓
- Restore successful: ✓
```

---

## Dependency Security

### Dependency Scanning

- [ ] Dependencies scanned for vulnerabilities
  - [ ] govulncheck run successfully
  - [ ] No critical vulnerabilities
  - [ ] Known issues documented

- [ ] Dependencies up to date
  - [ ] go.mod reviewed
  - [ ] Outdated packages identified
  - [ ] Update plan documented

- [ ] No unmaintained packages
  - [ ] Last commit dates checked
  - [ ] Active maintenance verified
  - [ ] Alternatives identified for EOL packages

- [ ] License compliance verified
  - [ ] All licenses compatible
  - [ ] Attribution requirements met
  - [ ] Copyleft concerns addressed

**Evidence/Notes:**
```
Vulnerability scan results:
$ govulncheck ./...

Dependency updates:
$ go list -u -m all

Unmaintained packages:
[List any concerns]
```

### Supply Chain

- [ ] go.sum file committed
  - [ ] Integrity checksums present
  - [ ] No missing dependencies
  - [ ] Verified against go.mod

- [ ] No suspicious dependencies
  - [ ] Package names verified
  - [ ] Sources reviewed
  - [ ] Maintainers researched

**Evidence/Notes:**
```
go.sum verification:
$ go mod verify
```

---

## Secure Defaults

### Configuration Defaults

- [ ] HTTPS used by default
  - [ ] No HTTP URLs in defaults
  - [ ] No insecure fallbacks

- [ ] SSH host key verification default documented
  - [ ] Current setting documented
  - [ ] Security implications explained
  - [ ] Mitigation documented

- [ ] Certificate validation enabled by default
  - [ ] No skip-verify flags
  - [ ] Proper TLS configuration

- [ ] Restrictive file permissions by default
  - [ ] 0600 for sensitive files
  - [ ] 0644 for configs
  - [ ] 0755 for directories

- [ ] No verbose logging by default
  - [ ] Sensitive data not logged
  - [ ] Debug mode disabled
  - [ ] Error messages sanitized

- [ ] Credentials never in default configs
  - [ ] Examples sanitized
  - [ ] Templates cleaned
  - [ ] Documentation reviewed

**Evidence/Notes:**
```
Default configuration review:
- HTTPS: Enforced
- Logging: Info level (not debug)
- Permissions: Restrictive
```

---

## Testing & Quality

### Security Testing

- [ ] Security test suite exists
  - [ ] Command injection tests
  - [ ] Path traversal tests
  - [ ] Input validation tests
  - [ ] Credential leak tests

- [ ] Fuzzing tests implemented
  - [ ] Input validation fuzzing
  - [ ] File path fuzzing
  - [ ] Configuration fuzzing

- [ ] Static analysis passing
  - [ ] gosec run successfully
  - [ ] staticcheck passing
  - [ ] No critical findings

**Evidence/Notes:**
```
Test results:
$ go test -v ./...
$ gosec ./...
$ staticcheck ./...

Coverage:
$ go test -cover ./...
```

### Code Quality

- [ ] Code review completed
  - [ ] Security checklist used
  - [ ] All changes reviewed
  - [ ] No security anti-patterns

- [ ] Documentation updated
  - [ ] Security features documented
  - [ ] Known limitations noted
  - [ ] Best practices included

**Evidence/Notes:**
```
Code review: [PR numbers or commit hashes]
Documentation: [files updated]
```

---

## Documentation

### Security Documentation

- [ ] SECURITY.md is current
  - [ ] Features documented
  - [ ] Limitations noted
  - [ ] Best practices included

- [ ] Vulnerability reporting process documented
  - [ ] Contact information current
  - [ ] Process clear
  - [ ] SLA defined

- [ ] Security checklist updated
  - [ ] This document current
  - [ ] All items relevant
  - [ ] Examples updated

- [ ] Known security limitations documented
  - [ ] SSH host key verification
  - [ ] Any temporary workarounds
  - [ ] Mitigation strategies

**Evidence/Notes:**
```
Documentation review date: __________
Reviewed by: __________
```

### User Documentation

- [ ] Security best practices documented
  - [ ] Credential management
  - [ ] Git repository security
  - [ ] Team security
  - [ ] Database dump handling

- [ ] Warning for sensitive operations
  - [ ] Database export warnings
  - [ ] /etc/hosts modification
  - [ ] Credential storage

**Evidence/Notes:**
```
User-facing security docs:
- docs/SECURITY.md
- README.md security section
```

---

## Git Repository Security

### Repository Configuration

- [ ] .gitignore includes sensitive files
  - [ ] .env files
  - [ ] *.sql and *.sql.gz
  - [ ] stax binary
  - [ ] .stax.yaml (local config)

- [ ] No credentials in git history
  - [ ] git-secrets scan passed
  - [ ] Manual review completed
  - [ ] History clean

- [ ] Example configs sanitized
  - [ ] No real credentials
  - [ ] Placeholder values used
  - [ ] Comments explain placeholders

**Evidence/Notes:**
```
git-secrets scan:
$ git secrets --scan

History audit:
$ git log --all --full-history -- '*.env' '*.sql'
```

### CI/CD Security

- [ ] Security scanning in CI/CD
  - [ ] govulncheck runs on PRs
  - [ ] gosec runs on PRs
  - [ ] Failed checks block merge

- [ ] Secrets not in CI/CD configs
  - [ ] GitHub Actions reviewed
  - [ ] No hardcoded credentials
  - [ ] Secrets use proper storage

**Evidence/Notes:**
```
CI/CD configuration files:
- .github/workflows/*.yml

Security checks configured: [list]
```

---

## Release Process

### Binary Security

- [ ] Binaries built securely
  - [ ] Clean build environment
  - [ ] Reproducible builds
  - [ ] No debug symbols in release

- [ ] Code signing planned/implemented
  - [ ] macOS: Developer ID
  - [ ] Windows: Authenticode (if applicable)
  - [ ] Linux: GPG signatures

- [ ] Release notes include security info
  - [ ] Security fixes listed
  - [ ] CVE numbers if applicable
  - [ ] Upgrade priority indicated

**Evidence/Notes:**
```
Build process:
$ make release

Signing status:
- macOS: [status]
- Windows: [status]
- Linux: [status]
```

### Distribution Security

- [ ] Download URLs use HTTPS
  - [ ] GitHub releases
  - [ ] Documentation links
  - [ ] Update checks

- [ ] Checksums provided
  - [ ] SHA256 sums
  - [ ] Signature verification
  - [ ] Instructions for verification

**Evidence/Notes:**
```
Distribution channels:
- GitHub releases: https://...
- Checksums: [link]
```

---

## Incident Response

### Preparation

- [ ] Incident response plan exists
  - [ ] Contact procedures defined
  - [ ] Escalation path clear
  - [ ] Communication templates ready

- [ ] Security contacts current
  - [ ] Email monitored
  - [ ] Response time committed
  - [ ] Backup contacts available

**Evidence/Notes:**
```
Security contact: security@firecrown-media.com
Response SLA: 48 hours for initial response
```

---

## Sign-off

### Pre-Release Approval

**Security Review Completed By:**

Name: ________________________
Date: ________________________
Signature: ________________________

**Critical Issues:**
- [ ] No critical issues OR
- [ ] All critical issues resolved OR
- [ ] Critical issues documented with accepted risk

**High Issues:**
- [ ] No high issues OR
- [ ] All high issues resolved OR
- [ ] High issues documented with mitigation plan

**Release Approval:**
- [ ] Security review complete
- [ ] All required items checked
- [ ] Documentation updated
- [ ] Team notified of any limitations

**Approved for Release:** Yes / No

**Conditions/Notes:**
```
[Any conditions for release or special notes]
```

---

## Post-Release

### Monitoring

- [ ] Security monitoring in place
  - [ ] Vulnerability scanning scheduled
  - [ ] Dependency updates tracked
  - [ ] Security advisories monitored

- [ ] Incident response ready
  - [ ] Contacts verified
  - [ ] Process tested
  - [ ] Communication plan ready

### Next Review

**Scheduled Security Review Date:** _______________

**Items to address in next release:**
```
1.
2.
3.
```

---

**Checklist Version:** 1.0
**Last Updated:** November 8, 2025
**Next Review:** With each release

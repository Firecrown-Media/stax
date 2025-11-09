# Provider Development Guide

## Overview

This guide explains how to create custom providers for Stax. Whether you're adding support for a new hosting platform or creating a specialized provider for your organization's infrastructure, this document will walk you through the process.

## Prerequisites

- Go 1.19 or later
- Understanding of the hosting platform's API/access methods
- Familiarity with Go interfaces and error handling

## Quick Start

### 1. Create Provider Package

```bash
mkdir -p pkg/providers/myprovider
```

### 2. Implement Provider Interface

```go
// pkg/providers/myprovider/provider.go
package myprovider

import (
    "github.com/firecrown/stax/pkg/provider"
)

type MyProvider struct {
    // Provider-specific fields
}

func init() {
    provider.RegisterProvider("myprovider", &MyProvider{})
}

// Implement all required interface methods...
```

### 3. Test Your Provider

```go
// pkg/providers/myprovider/provider_test.go
func TestMyProvider_ImplementsInterface(t *testing.T) {
    var _ provider.Provider = (*MyProvider)(nil)
}
```

## Step-by-Step Implementation

### Step 1: Define Provider Structure

```go
package myprovider

import (
    "github.com/firecrown/stax/pkg/provider"
)

type MyProvider struct {
    // Authentication fields
    apiKey      string
    apiSecret   string

    // Client/connection fields
    httpClient  *http.Client
    sshClient   *ssh.Client

    // Configuration fields
    region      string
    instanceID  string
}
```

### Step 2: Register Provider

```go
func init() {
    provider.RegisterProvider("myprovider", &MyProvider{})
}

func New() *MyProvider {
    return &MyProvider{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}
```

### Step 3: Implement Metadata Methods

```go
func (p *MyProvider) Name() string {
    return "myprovider"
}

func (p *MyProvider) Description() string {
    return "My Custom WordPress Hosting Provider"
}

func (p *MyProvider) Capabilities() provider.ProviderCapabilities {
    return provider.ProviderCapabilities{
        // Core capabilities
        Authentication: true,
        SiteManagement: true,
        DatabaseExport: true,
        DatabaseImport: true,
        FileSync:       true,

        // Optional capabilities
        Deployment:      false,
        Environments:    false,
        Backups:         true,
        RemoteExecution: true,
        MediaManagement: false,
        SSHAccess:       true,
        APIAccess:       true,
    }
}
```

### Step 4: Implement Authentication

```go
func (p *MyProvider) ValidateCredentials(credentials map[string]string) error {
    required := []string{"api_key", "api_secret", "region"}

    for _, key := range required {
        if credentials[key] == "" {
            return fmt.Errorf("missing required credential: %s", key)
        }
    }

    return nil
}

func (p *MyProvider) Authenticate(credentials map[string]string) error {
    if err := p.ValidateCredentials(credentials); err != nil {
        return err
    }

    p.apiKey = credentials["api_key"]
    p.apiSecret = credentials["api_secret"]
    p.region = credentials["region"]

    // Initialize API client
    p.httpClient.Transport = &authenticatedTransport{
        apiKey:    p.apiKey,
        apiSecret: p.apiSecret,
    }

    return p.TestConnection()
}

func (p *MyProvider) TestConnection() error {
    // Test API connectivity
    resp, err := p.httpClient.Get("https://api.myprovider.com/v1/ping")
    if err != nil {
        return fmt.Errorf("connection test failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("connection test failed: status %d", resp.StatusCode)
    }

    return nil
}
```

### Step 5: Implement Site Management

```go
func (p *MyProvider) ListSites() ([]provider.Site, error) {
    resp, err := p.httpClient.Get(fmt.Sprintf(
        "https://api.myprovider.com/v1/sites?region=%s", p.region))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var apiSites []APISite
    if err := json.NewDecoder(resp.Body).Decode(&apiSites); err != nil {
        return nil, err
    }

    sites := make([]provider.Site, len(apiSites))
    for i, apiSite := range apiSites {
        sites[i] = provider.Site{
            ID:            apiSite.ID,
            Name:          apiSite.Name,
            PrimaryDomain: apiSite.Domain,
            Environment:   apiSite.Environment,
            Status:        apiSite.Status,
            Provider:      "myprovider",
            Metadata: map[string]string{
                "region": apiSite.Region,
            },
        }
    }

    return sites, nil
}

func (p *MyProvider) GetSite(identifier string) (*provider.Site, error) {
    // Try to find by ID, name, or domain
    sites, err := p.ListSites()
    if err != nil {
        return nil, err
    }

    for _, site := range sites {
        if site.ID == identifier ||
           site.Name == identifier ||
           site.PrimaryDomain == identifier {
            return &site, nil
        }
    }

    return nil, fmt.Errorf("site not found: %s", identifier)
}

func (p *MyProvider) GetSiteMetadata(site *provider.Site) (*provider.SiteMetadata, error) {
    resp, err := p.httpClient.Get(fmt.Sprintf(
        "https://api.myprovider.com/v1/sites/%s/details", site.ID))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var details APISiteDetails
    if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
        return nil, err
    }

    return &provider.SiteMetadata{
        Site:             site,
        PHPVersion:       details.PHPVersion,
        MySQLVersion:     details.MySQLVersion,
        WordPressVersion: details.WordPressVersion,
        DiskUsage: provider.DiskUsage{
            Used:  details.DiskUsed,
            Total: details.DiskTotal,
        },
        Domains:  details.Domains,
        Features: details.Features,
    }, nil
}
```

### Step 6: Implement Database Operations

```go
func (p *MyProvider) ExportDatabase(site *provider.Site, options provider.DatabaseExportOptions) (io.ReadCloser, error) {
    // Option 1: API-based export
    // Trigger export via API, return download stream

    // Option 2: SSH-based export
    // Use SSH to run mysqldump or wp db export

    // Example SSH approach:
    if p.sshClient == nil {
        return nil, fmt.Errorf("SSH not configured")
    }

    session, err := p.sshClient.NewSession()
    if err != nil {
        return nil, err
    }

    stdout, err := session.StdoutPipe()
    if err != nil {
        session.Close()
        return nil, err
    }

    // Build export command
    cmd := buildExportCommand(options)

    if err := session.Start(cmd); err != nil {
        session.Close()
        return nil, err
    }

    return &sessionReadCloser{
        ReadCloser: stdout,
        session:    session,
    }, nil
}

func (p *MyProvider) ImportDatabase(site *provider.Site, data io.Reader, options provider.DatabaseImportOptions) error {
    // Stream SQL to database via SSH or API

    if p.sshClient == nil {
        return fmt.Errorf("SSH not configured")
    }

    session, err := p.sshClient.NewSession()
    if err != nil {
        return err
    }
    defer session.Close()

    session.Stdin = data

    cmd := buildImportCommand(options)

    return session.Run(cmd)
}

func (p *MyProvider) GetDatabaseCredentials(site *provider.Site) (*provider.DatabaseCredentials, error) {
    // Retrieve from API or parse wp-config.php via SSH

    resp, err := p.httpClient.Get(fmt.Sprintf(
        "https://api.myprovider.com/v1/sites/%s/database", site.ID))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var dbCreds APIDBCredentials
    if err := json.NewDecoder(resp.Body).Decode(&dbCreds); err != nil {
        return nil, err
    }

    return &provider.DatabaseCredentials{
        Host:     dbCreds.Host,
        Port:     dbCreds.Port,
        Database: dbCreds.Database,
        Username: dbCreds.Username,
        Password: dbCreds.Password,
        SSL:      dbCreds.SSL,
    }, nil
}
```

### Step 7: Implement File Operations

```go
func (p *MyProvider) SyncFiles(site *provider.Site, destination string, options provider.SyncOptions) error {
    // Use rsync, SFTP, or API-based sync

    // Example rsync approach:
    source := fmt.Sprintf("%s@%s:/var/www/html/wp-content/",
        p.sshUser, site.Metadata["server_ip"])

    args := buildRsyncArgs(source, destination, options)

    cmd := exec.Command("rsync", args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    return cmd.Run()
}

func (p *MyProvider) DownloadFile(site *provider.Site, remotePath string) (io.ReadCloser, error) {
    // Use SFTP, SCP, or API download

    session, err := p.sshClient.NewSession()
    if err != nil {
        return nil, err
    }

    stdout, err := session.StdoutPipe()
    if err != nil {
        session.Close()
        return nil, err
    }

    cmd := fmt.Sprintf("cat %s", remotePath)
    if err := session.Start(cmd); err != nil {
        session.Close()
        return nil, err
    }

    return &sessionReadCloser{
        ReadCloser: stdout,
        session:    session,
    }, nil
}

func (p *MyProvider) UploadFile(site *provider.Site, localPath, remotePath string) error {
    // Use SFTP, SCP, or API upload

    file, err := os.Open(localPath)
    if err != nil {
        return err
    }
    defer file.Close()

    session, err := p.sshClient.NewSession()
    if err != nil {
        return err
    }
    defer session.Close()

    session.Stdin = file

    cmd := fmt.Sprintf("cat > %s", remotePath)
    return session.Run(cmd)
}
```

### Step 8: Implement Environment Information

```go
func (p *MyProvider) GetPHPVersion(site *provider.Site) (string, error) {
    // Query via SSH or API

    output, err := p.executeCommand(site, "php -v")
    if err != nil {
        return "", err
    }

    // Parse version from output
    return parseVersion(output, "PHP"), nil
}

func (p *MyProvider) GetMySQLVersion(site *provider.Site) (string, error) {
    output, err := p.executeCommand(site, "mysql --version")
    if err != nil {
        return "", err
    }

    return parseVersion(output, "mysql"), nil
}

func (p *MyProvider) GetWordPressVersion(site *provider.Site) (string, error) {
    output, err := p.executeCommand(site, "wp core version")
    if err != nil {
        return "", err
    }

    return strings.TrimSpace(output), nil
}
```

## Implementing Optional Capabilities

### BackupManager Interface

```go
var _ provider.BackupManager = (*MyProvider)(nil)

func (p *MyProvider) ListBackups(site *provider.Site) ([]provider.Backup, error) {
    // Implementation
}

func (p *MyProvider) CreateBackup(site *provider.Site, description string) (*provider.Backup, error) {
    // Implementation
}

// ... other BackupManager methods
```

### RemoteExecutor Interface

```go
var _ provider.RemoteExecutor = (*MyProvider)(nil)

func (p *MyProvider) ExecuteCommand(site *provider.Site, command string) (string, error) {
    // Implementation
}

func (p *MyProvider) ExecuteWPCLI(site *provider.Site, args []string) (string, error) {
    cmd := "wp " + strings.Join(args, " ")
    return p.ExecuteCommand(site, cmd)
}

func (p *MyProvider) StreamCommand(site *provider.Site, command string, stdout, stderr io.Writer) error {
    // Implementation
}
```

## Testing Your Provider

### Unit Tests

```go
func TestMyProvider_ValidateCredentials(t *testing.T) {
    p := &MyProvider{}

    tests := []struct {
        name        string
        credentials map[string]string
        wantErr     bool
    }{
        {
            name: "valid credentials",
            credentials: map[string]string{
                "api_key":    "test-key",
                "api_secret": "test-secret",
                "region":     "us-west-2",
            },
            wantErr: false,
        },
        {
            name: "missing api_key",
            credentials: map[string]string{
                "api_secret": "test-secret",
                "region":     "us-west-2",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := p.ValidateCredentials(tt.credentials)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateCredentials() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests

```go
func TestMyProvider_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Set up test provider
    p := setupTestProvider(t)

    // Test connection
    err := p.TestConnection()
    if err != nil {
        t.Fatalf("TestConnection failed: %v", err)
    }

    // Test site listing
    sites, err := p.ListSites()
    if err != nil {
        t.Fatalf("ListSites failed: %v", err)
    }

    if len(sites) == 0 {
        t.Skip("No sites available for testing")
    }

    // Test site details
    site := &sites[0]
    metadata, err := p.GetSiteMetadata(site)
    if err != nil {
        t.Fatalf("GetSiteMetadata failed: %v", err)
    }

    t.Logf("Site: %s, PHP: %s, WordPress: %s",
        site.Name, metadata.PHPVersion, metadata.WordPressVersion)
}
```

## Configuration Documentation

Document your provider's configuration schema in `docs/PROVIDER_<NAME>.md`:

```markdown
# MyProvider Configuration

## Required Credentials

- `api_key`: API key from MyProvider dashboard
- `api_secret`: API secret from MyProvider dashboard
- `region`: Provider region (us-west-2, eu-west-1, etc.)

## Configuration Example

yaml
provider:
  name: myprovider
  myprovider:
    api_key: ${MYPROVIDER_API_KEY}
    api_secret: ${MYPROVIDER_API_SECRET}
    region: us-west-2
    instance_id: inst-12345


## Getting Credentials

1. Log in to MyProvider dashboard
2. Navigate to Settings → API
3. Generate new API key pair
4. Save credentials in `~/.stax/credentials.yml`
```

## Best Practices

1. **Error Handling**: Return descriptive errors with context
2. **Resource Cleanup**: Use `defer` to ensure cleanup
3. **Streaming**: Use `io.Reader/Writer` for large data
4. **Credentials**: Never log credentials
5. **Rate Limiting**: Implement rate limiting for API calls
6. **Retries**: Use exponential backoff for transient failures
7. **Logging**: Use structured logging
8. **Testing**: Provide both unit and integration tests
9. **Documentation**: Document all provider-specific behavior
10. **Versioning**: Handle API version changes gracefully

## Common Pitfalls

1. **Forgetting to register provider** in `init()`
2. **Not implementing all interface methods**
3. **Leaking SSH/HTTP connections**
4. **Not handling API pagination**
5. **Hardcoding timeouts**
6. **Not validating inputs**
7. **Ignoring context cancellation**
8. **Poor error messages**

## Example: Complete Minimal Provider

See `pkg/providers/local/provider.go` for a minimal provider implementation.

## Publishing Your Provider

1. **Fork stax repository**
2. **Add provider to `pkg/providers/<name>/`**
3. **Add tests and documentation**
4. **Submit pull request**
5. **Add to provider registry** (maintainer will handle)

## Getting Help

- GitHub Issues: Report bugs or ask questions
- Discussions: Architecture and design questions
- Examples: See existing providers in `pkg/providers/`

---

**Last Updated**: 2025-11-08

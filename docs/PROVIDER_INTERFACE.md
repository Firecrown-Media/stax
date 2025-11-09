# Provider Interface Specification

## Overview

The Provider Interface is the contract that all WordPress hosting providers must implement to integrate with Stax. This document provides the complete specification, including required methods, optional capabilities, common types, and implementation guidelines.

## Core Provider Interface

All providers must implement the `Provider` interface defined in `pkg/provider/interface.go`:

```go
package provider

import (
    "io"
)

// Provider is the core interface that all hosting providers must implement
type Provider interface {
    // ===== Metadata =====

    // Name returns the provider's unique identifier (e.g., "wpengine", "aws", "wordpress-vip")
    Name() string

    // Description returns a human-readable description of the provider
    Description() string

    // Capabilities returns the provider's capability set
    Capabilities() ProviderCapabilities

    // ===== Authentication & Setup =====

    // Authenticate authenticates with the provider using provided credentials
    // credentials map contains provider-specific authentication parameters
    // Returns error if authentication fails
    Authenticate(credentials map[string]string) error

    // TestConnection tests the connection to the provider
    // Returns error if connection test fails
    TestConnection() error

    // ValidateCredentials validates credentials without establishing a connection
    // Returns error if credentials are invalid or incomplete
    ValidateCredentials(credentials map[string]string) error

    // ===== Site Management =====

    // ListSites lists all sites/installations available on this provider
    ListSites() ([]Site, error)

    // GetSite retrieves information about a specific site
    // identifier can be site ID, name, or domain (provider-dependent)
    GetSite(identifier string) (*Site, error)

    // GetSiteMetadata retrieves detailed metadata about a site
    GetSiteMetadata(site *Site) (*SiteMetadata, error)

    // ===== Database Operations =====

    // ExportDatabase exports the database from the remote site
    // Returns an io.ReadCloser that streams the database dump
    // Caller is responsible for closing the ReadCloser
    ExportDatabase(site *Site, options DatabaseExportOptions) (io.ReadCloser, error)

    // ImportDatabase imports a database to the remote site
    // data is an io.Reader containing the SQL dump
    ImportDatabase(site *Site, data io.Reader, options DatabaseImportOptions) error

    // GetDatabaseCredentials retrieves database connection credentials
    // Useful for direct database access or debugging
    GetDatabaseCredentials(site *Site) (*DatabaseCredentials, error)

    // ===== File Operations =====

    // SyncFiles synchronizes files from remote site to local destination
    // Typically syncs wp-content or specific directories
    SyncFiles(site *Site, destination string, options SyncOptions) error

    // DownloadFile downloads a single file from the remote site
    DownloadFile(site *Site, remotePath string) (io.ReadCloser, error)

    // UploadFile uploads a single file to the remote site
    UploadFile(site *Site, localPath, remotePath string) error

    // ===== Environment Information =====

    // GetPHPVersion returns the PHP version for the site
    GetPHPVersion(site *Site) (string, error)

    // GetMySQLVersion returns the MySQL/MariaDB version for the site
    GetMySQLVersion(site *Site) (string, error)

    // GetWordPressVersion returns the WordPress version for the site
    GetWordPressVersion(site *Site) (string, error)
}
```

## Common Types

### Site

Represents a WordPress site/installation on a provider:

```go
// Site represents a WordPress site on a hosting provider
type Site struct {
    ID            string            `json:"id"`
    Name          string            `json:"name"`
    PrimaryDomain string            `json:"primary_domain"`
    Environment   string            `json:"environment"`   // e.g., "production", "staging"
    Status        string            `json:"status"`        // e.g., "active", "suspended"
    Provider      string            `json:"provider"`      // Provider name
    Metadata      map[string]string `json:"metadata"`      // Provider-specific metadata
}
```

### SiteMetadata

Detailed information about a site:

```go
// SiteMetadata contains detailed information about a site
type SiteMetadata struct {
    Site             *Site             `json:"site"`
    PHPVersion       string            `json:"php_version"`
    MySQLVersion     string            `json:"mysql_version"`
    WordPressVersion string            `json:"wordpress_version"`
    DiskUsage        DiskUsage         `json:"disk_usage"`
    Domains          []string          `json:"domains"`
    Features         []string          `json:"features"`      // e.g., "ssl", "cdn", "backups"
    CreatedAt        string            `json:"created_at"`
    UpdatedAt        string            `json:"updated_at"`
}

// DiskUsage represents disk space usage
type DiskUsage struct {
    Used  int64 `json:"used"`   // Bytes used
    Total int64 `json:"total"`  // Total bytes available
}
```

### Environment

Represents a site environment (production, staging, development):

```go
// Environment represents a site environment
type Environment struct {
    Name          string `json:"name"`           // e.g., "production", "staging"
    URL           string `json:"url"`
    Status        string `json:"status"`         // e.g., "active", "inactive"
    IsDefault     bool   `json:"is_default"`
    LastDeployAt  string `json:"last_deploy_at"`
}
```

### Database Types

```go
// DatabaseExportOptions configures database export behavior
type DatabaseExportOptions struct {
    ExcludeTables  []string `json:"exclude_tables"`   // Tables to exclude
    SkipLogs       bool     `json:"skip_logs"`        // Skip log tables
    SkipTransients bool     `json:"skip_transients"`  // Skip transient data
    SkipSpam       bool     `json:"skip_spam"`        // Skip spam comments
    Compress       bool     `json:"compress"`         // Compress output (gzip)
    IncludePrefix  bool     `json:"include_prefix"`   // Include table prefix detection
}

// DatabaseImportOptions configures database import behavior
type DatabaseImportOptions struct {
    DropExisting  bool     `json:"drop_existing"`    // Drop existing tables
    SearchReplace []string `json:"search_replace"`   // Search/replace patterns
    SkipErrors    bool     `json:"skip_errors"`      // Continue on SQL errors
}

// DatabaseCredentials contains database connection information
type DatabaseCredentials struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
    SSL      bool   `json:"ssl"`
}
```

### File Sync Types

```go
// SyncOptions configures file synchronization
type SyncOptions struct {
    Source         string   `json:"source"`           // Source path (optional, provider determines default)
    Destination    string   `json:"destination"`      // Local destination path
    Include        []string `json:"include"`          // Patterns to include
    Exclude        []string `json:"exclude"`          // Patterns to exclude
    Delete         bool     `json:"delete"`           // Delete files not on remote
    DryRun         bool     `json:"dry_run"`          // Perform dry run
    BandwidthLimit int      `json:"bandwidth_limit"`  // KB/s limit
    Progress       bool     `json:"progress"`         // Show progress
}
```

## Provider Capabilities

Capabilities define what features a provider supports:

```go
// ProviderCapabilities describes what features a provider supports
type ProviderCapabilities struct {
    // Core capabilities (should always be true if implemented correctly)
    Authentication  bool `json:"authentication"`
    SiteManagement  bool `json:"site_management"`
    DatabaseExport  bool `json:"database_export"`
    DatabaseImport  bool `json:"database_import"`
    FileSync        bool `json:"file_sync"`

    // Optional capabilities
    Deployment      bool `json:"deployment"`        // Git-based deployments
    Environments    bool `json:"environments"`      // Multi-environment support
    Backups         bool `json:"backups"`           // Automated backups
    RemoteExecution bool `json:"remote_execution"`  // SSH/WP-CLI
    MediaManagement bool `json:"media_management"`  // CDN/media proxy
    SSHAccess       bool `json:"ssh_access"`        // Direct SSH access
    APIAccess       bool `json:"api_access"`        // REST API access

    // Advanced capabilities
    Scaling         bool `json:"scaling"`           // Auto-scaling support
    Monitoring      bool `json:"monitoring"`        // Performance monitoring
    Logging         bool `json:"logging"`           // Centralized logging
}
```

## Optional Capability Interfaces

Providers can implement additional interfaces to provide optional features:

### Deployer Interface

For providers that support Git-based deployments:

```go
// Deployer interface for providers that support deployments
type Deployer interface {
    Provider

    // Deploy deploys code to the site
    Deploy(site *Site, options DeployOptions) (*Deployment, error)

    // GetDeploymentStatus checks the status of a deployment
    GetDeploymentStatus(site *Site, deploymentID string) (*DeploymentStatus, error)

    // ListDeployments lists recent deployments
    ListDeployments(site *Site) ([]Deployment, error)
}

// DeployOptions configures deployment behavior
type DeployOptions struct {
    Branch      string            `json:"branch"`       // Git branch to deploy
    Commit      string            `json:"commit"`       // Specific commit (optional)
    Message     string            `json:"message"`      // Deployment message
    Environment string            `json:"environment"`  // Target environment
    Metadata    map[string]string `json:"metadata"`     // Provider-specific options
}

// Deployment represents a deployment
type Deployment struct {
    ID          string `json:"id"`
    Status      string `json:"status"`       // "pending", "in_progress", "completed", "failed"
    Branch      string `json:"branch"`
    Commit      string `json:"commit"`
    Message     string `json:"message"`
    DeployedAt  string `json:"deployed_at"`
    DeployedBy  string `json:"deployed_by"`
}

// DeploymentStatus represents deployment status details
type DeploymentStatus struct {
    Deployment  *Deployment `json:"deployment"`
    Progress    int         `json:"progress"`      // 0-100
    Phase       string      `json:"phase"`         // Current phase
    Logs        []string    `json:"logs"`          // Recent log lines
    Error       string      `json:"error"`         // Error message if failed
}
```

### EnvironmentManager Interface

For providers that support multiple environments (staging, production):

```go
// EnvironmentManager interface for providers with multi-environment support
type EnvironmentManager interface {
    Provider

    // ListEnvironments lists available environments for a site
    ListEnvironments(site *Site) ([]Environment, error)

    // GetEnvironment retrieves information about a specific environment
    GetEnvironment(site *Site, environmentName string) (*Environment, error)

    // SwitchEnvironment switches to a different environment
    SwitchEnvironment(site *Site, environmentName string) error

    // CreateEnvironment creates a new environment (if supported)
    CreateEnvironment(site *Site, environmentName string, options EnvironmentOptions) error

    // DeleteEnvironment deletes an environment (if supported)
    DeleteEnvironment(site *Site, environmentName string) error
}

// EnvironmentOptions configures environment creation
type EnvironmentOptions struct {
    CloneFrom   string            `json:"clone_from"`   // Clone from existing environment
    PHPVersion  string            `json:"php_version"`
    Domain      string            `json:"domain"`
    Metadata    map[string]string `json:"metadata"`
}
```

### BackupManager Interface

For providers with backup capabilities:

```go
// BackupManager interface for providers that support backups
type BackupManager interface {
    Provider

    // ListBackups lists available backups for a site
    ListBackups(site *Site) ([]Backup, error)

    // CreateBackup creates a manual backup
    CreateBackup(site *Site, description string) (*Backup, error)

    // RestoreBackup restores a site from a backup
    RestoreBackup(site *Site, backupID string, options RestoreOptions) error

    // DeleteBackup deletes a backup
    DeleteBackup(site *Site, backupID string) error

    // DownloadBackup downloads a backup archive
    DownloadBackup(site *Site, backupID string) (io.ReadCloser, error)
}

// Backup represents a backup
type Backup struct {
    ID          string `json:"id"`
    Type        string `json:"type"`         // "manual", "automatic", "scheduled"
    Description string `json:"description"`
    Size        int64  `json:"size"`
    CreatedAt   string `json:"created_at"`
    Status      string `json:"status"`       // "pending", "completed", "failed"
    ExpiresAt   string `json:"expires_at"`
}

// RestoreOptions configures backup restoration
type RestoreOptions struct {
    DatabaseOnly bool   `json:"database_only"`  // Restore only database
    FilesOnly    bool   `json:"files_only"`     // Restore only files
    Environment  string `json:"environment"`    // Target environment
}
```

### RemoteExecutor Interface

For providers that support remote command execution:

```go
// RemoteExecutor interface for providers that support remote execution
type RemoteExecutor interface {
    Provider

    // ExecuteCommand executes a shell command on the remote server
    ExecuteCommand(site *Site, command string) (string, error)

    // ExecuteWPCLI executes a WP-CLI command
    ExecuteWPCLI(site *Site, args []string) (string, error)

    // StreamCommand executes a command and streams output
    StreamCommand(site *Site, command string, stdout, stderr io.Writer) error
}
```

### MediaManager Interface

For providers with CDN or media proxy capabilities:

```go
// MediaManager interface for providers with media/CDN support
type MediaManager interface {
    Provider

    // GetMediaURL returns the CDN or media proxy URL for a site
    GetMediaURL(site *Site) (string, error)

    // SupportsRemoteMedia indicates if provider supports remote media serving
    SupportsRemoteMedia() bool

    // ConfigureMedia configures media settings
    ConfigureMedia(site *Site, options MediaOptions) error

    // PurgeMediaCache purges the media cache (if CDN)
    PurgeMediaCache(site *Site, paths []string) error
}

// MediaOptions configures media/CDN settings
type MediaOptions struct {
    CDNEnabled  bool     `json:"cdn_enabled"`
    CDNDomain   string   `json:"cdn_domain"`
    CacheTTL    int      `json:"cache_ttl"`      // Seconds
    Excludes    []string `json:"excludes"`       // Paths to exclude from CDN
}
```

## Implementation Guidelines

### 1. Provider Registration

Register your provider in the `init()` function:

```go
package myprovider

import "github.com/firecrown/stax/pkg/provider"

func init() {
    provider.RegisterProvider("myprovider", NewMyProvider())
}

func NewMyProvider() *MyProvider {
    return &MyProvider{}
}
```

### 2. Credential Handling

Credentials are passed as a `map[string]string`. Define expected keys:

```go
func (p *MyProvider) ValidateCredentials(creds map[string]string) error {
    required := []string{"api_key", "api_secret"}
    for _, key := range required {
        if creds[key] == "" {
            return fmt.Errorf("missing required credential: %s", key)
        }
    }
    return nil
}

func (p *MyProvider) Authenticate(creds map[string]string) error {
    if err := p.ValidateCredentials(creds); err != nil {
        return err
    }

    p.apiKey = creds["api_key"]
    p.apiSecret = creds["api_secret"]

    return p.TestConnection()
}
```

### 3. Error Handling

Return descriptive errors with context:

```go
func (p *MyProvider) ExportDatabase(site *Site, options DatabaseExportOptions) (io.ReadCloser, error) {
    if site == nil {
        return nil, fmt.Errorf("site cannot be nil")
    }

    reader, err := p.performExport(site, options)
    if err != nil {
        return nil, fmt.Errorf("failed to export database for site %s: %w", site.Name, err)
    }

    return reader, nil
}
```

### 4. Capability Detection

Implement the `Capabilities()` method accurately:

```go
func (p *MyProvider) Capabilities() provider.ProviderCapabilities {
    return provider.ProviderCapabilities{
        Authentication:  true,
        SiteManagement:  true,
        DatabaseExport:  true,
        DatabaseImport:  true,
        FileSync:        true,
        Deployment:      false,  // Not supported
        Environments:    false,  // Not supported
        Backups:         true,
        RemoteExecution: true,
        MediaManagement: false,  // Not supported
        SSHAccess:       true,
        APIAccess:       true,
    }
}
```

### 5. Streaming Large Data

Use `io.Reader` and `io.Writer` for large operations:

```go
func (p *MyProvider) ExportDatabase(site *Site, options DatabaseExportOptions) (io.ReadCloser, error) {
    // Open SSH session
    session, err := p.sshClient.NewSession()
    if err != nil {
        return nil, err
    }

    stdout, err := session.StdoutPipe()
    if err != nil {
        session.Close()
        return nil, err
    }

    // Start export command (streams to stdout)
    if err := session.Start("wp db export -"); err != nil {
        session.Close()
        return nil, err
    }

    // Return reader that closes session when done
    return &exportReadCloser{
        ReadCloser: stdout,
        session:    session,
    }, nil
}
```

### 6. Site Identifier Flexibility

Support multiple identifier types:

```go
func (p *MyProvider) GetSite(identifier string) (*Site, error) {
    // Try as site ID first
    if site := p.findByID(identifier); site != nil {
        return site, nil
    }

    // Try as site name
    if site := p.findByName(identifier); site != nil {
        return site, nil
    }

    // Try as primary domain
    if site := p.findByDomain(identifier); site != nil {
        return site, nil
    }

    return nil, fmt.Errorf("site not found: %s", identifier)
}
```

### 7. Provider Metadata

Use the `Metadata` field for provider-specific data:

```go
site := &provider.Site{
    ID:            "12345",
    Name:          "my-site",
    PrimaryDomain: "example.com",
    Environment:   "production",
    Provider:      "myprovider",
    Metadata: map[string]string{
        "region":      "us-west-2",
        "instance_id": "i-1234567890abcdef",
        "account_id":  "123456789012",
    },
}
```

## Testing Requirements

### Unit Tests

Test each method independently:

```go
func TestMyProvider_Authenticate(t *testing.T) {
    p := NewMyProvider()

    creds := map[string]string{
        "api_key":    "test-key",
        "api_secret": "test-secret",
    }

    err := p.Authenticate(creds)
    assert.NoError(t, err)
}
```

### Integration Tests

Test with mock server or real provider (if API available):

```go
func TestMyProvider_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    p := setupTestProvider(t)

    sites, err := p.ListSites()
    assert.NoError(t, err)
    assert.NotEmpty(t, sites)
}
```

### Interface Compliance

Verify interface implementation:

```go
func TestMyProvider_ImplementsProvider(t *testing.T) {
    var _ provider.Provider = (*MyProvider)(nil)
}

func TestMyProvider_ImplementsDeployer(t *testing.T) {
    var _ provider.Deployer = (*MyProvider)(nil)
}
```

## Example Implementation

See `pkg/providers/wpengine/provider.go` for a complete reference implementation.

## Best Practices

1. **Fail Fast**: Validate inputs early and return errors immediately
2. **Resource Cleanup**: Use `defer` to ensure resources (SSH sessions, HTTP connections) are cleaned up
3. **Context Support**: Consider adding `context.Context` to long-running operations (future enhancement)
4. **Rate Limiting**: Implement rate limiting for API calls
5. **Retry Logic**: Implement exponential backoff for transient failures
6. **Logging**: Use structured logging for debugging
7. **Documentation**: Document provider-specific quirks and limitations
8. **Versioning**: Handle API version changes gracefully

## Provider-Specific Configuration

Document required configuration in `docs/PROVIDER_<NAME>.md`:

```yaml
provider:
  name: myprovider

  myprovider:
    api_key: ${MYPROVIDER_API_KEY}
    api_secret: ${MYPROVIDER_API_SECRET}
    region: us-west-2
    custom_option: value
```

## Migration Support

If your provider supports migration from other providers, implement:

```go
type Migrator interface {
    Provider

    // ImportFromProvider imports a site from another provider
    ImportFromProvider(sourceProvider Provider, sourceSite *Site, options MigrationOptions) error
}
```

---

**Version**: 1.0.0
**Last Updated**: 2025-11-08
**Maintainer**: Firecrown Development Team

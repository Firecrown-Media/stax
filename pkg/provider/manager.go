package provider

import (
	"fmt"
	"io"
)

// Manager provides high-level provider management operations
type Manager struct {
	currentProvider Provider
	providerName    string
}

// NewManager creates a new provider manager with the given provider
func NewManager(provider Provider) *Manager {
	return &Manager{
		currentProvider: provider,
		providerName:    provider.Name(),
	}
}

// NewManagerFromConfig creates a new manager from configuration
func NewManagerFromConfig(config ProviderConfig) (*Manager, error) {
	provider, err := NewProvider(config)
	if err != nil {
		return nil, err
	}

	return NewManager(provider), nil
}

// GetProvider returns the current provider
func (m *Manager) GetProvider() Provider {
	return m.currentProvider
}

// GetProviderName returns the current provider name
func (m *Manager) GetProviderName() string {
	return m.providerName
}

// SwitchProvider switches to a different provider
func (m *Manager) SwitchProvider(config ProviderConfig) error {
	provider, err := NewProvider(config)
	if err != nil {
		return fmt.Errorf("failed to switch provider: %w", err)
	}

	m.currentProvider = provider
	m.providerName = provider.Name()

	return nil
}

// TestCurrentProvider tests the connection to the current provider
func (m *Manager) TestCurrentProvider() error {
	if m.currentProvider == nil {
		return fmt.Errorf("no provider configured")
	}

	return m.currentProvider.TestConnection()
}

// GetCurrentProviderInfo returns information about the current provider
func (m *Manager) GetCurrentProviderInfo() (*ProviderInfo, error) {
	if m.currentProvider == nil {
		return nil, fmt.Errorf("no provider configured")
	}

	return &ProviderInfo{
		Name:         m.currentProvider.Name(),
		Description:  m.currentProvider.Description(),
		Capabilities: m.currentProvider.Capabilities(),
		IsDefault:    m.providerName == GetDefaultProvider(),
	}, nil
}

// ===== High-Level Operations =====

// ListSites lists all sites from the current provider
func (m *Manager) ListSites() ([]Site, error) {
	if m.currentProvider == nil {
		return nil, fmt.Errorf("no provider configured")
	}

	return m.currentProvider.ListSites()
}

// GetSite retrieves a site by identifier
func (m *Manager) GetSite(identifier string) (*Site, error) {
	if m.currentProvider == nil {
		return nil, fmt.Errorf("no provider configured")
	}

	return m.currentProvider.GetSite(identifier)
}

// ExportDatabase exports a database from the current provider
func (m *Manager) ExportDatabase(site *Site, options DatabaseExportOptions) (io.ReadCloser, error) {
	if m.currentProvider == nil {
		return nil, fmt.Errorf("no provider configured")
	}

	return m.currentProvider.ExportDatabase(site, options)
}

// ImportDatabase imports a database to the current provider
func (m *Manager) ImportDatabase(site *Site, data io.Reader, options DatabaseImportOptions) error {
	if m.currentProvider == nil {
		return fmt.Errorf("no provider configured")
	}

	return m.currentProvider.ImportDatabase(site, data, options)
}

// SyncFiles synchronizes files from the current provider
func (m *Manager) SyncFiles(site *Site, destination string, options SyncOptions) error {
	if m.currentProvider == nil {
		return fmt.Errorf("no provider configured")
	}

	return m.currentProvider.SyncFiles(site, destination, options)
}

// ===== Optional Capability Helpers =====

// Deploy deploys code (if provider supports it)
func (m *Manager) Deploy(site *Site, options DeployOptions) (*Deployment, error) {
	deployer, ok := m.currentProvider.(Deployer)
	if !ok {
		return nil, fmt.Errorf("provider %s does not support deployment", m.providerName)
	}

	return deployer.Deploy(site, options)
}

// ListEnvironments lists environments (if provider supports it)
func (m *Manager) ListEnvironments(site *Site) ([]Environment, error) {
	envManager, ok := m.currentProvider.(EnvironmentManager)
	if !ok {
		return nil, fmt.Errorf("provider %s does not support environment management", m.providerName)
	}

	return envManager.ListEnvironments(site)
}

// CreateBackup creates a backup (if provider supports it)
func (m *Manager) CreateBackup(site *Site, description string) (*Backup, error) {
	backupManager, ok := m.currentProvider.(BackupManager)
	if !ok {
		return nil, fmt.Errorf("provider %s does not support backups", m.providerName)
	}

	return backupManager.CreateBackup(site, description)
}

// ExecuteWPCLI executes a WP-CLI command (if provider supports it)
func (m *Manager) ExecuteWPCLI(site *Site, args []string) (string, error) {
	executor, ok := m.currentProvider.(RemoteExecutor)
	if !ok {
		return "", fmt.Errorf("provider %s does not support remote execution", m.providerName)
	}

	return executor.ExecuteWPCLI(site, args)
}

// GetMediaURL gets the media URL (if provider supports it)
func (m *Manager) GetMediaURL(site *Site) (string, error) {
	mediaManager, ok := m.currentProvider.(MediaManager)
	if !ok {
		return "", fmt.Errorf("provider %s does not support media management", m.providerName)
	}

	return mediaManager.GetMediaURL(site)
}

// ===== Multi-Provider Operations =====

// MigrateOptions contains options for migrating between providers
type MigrateOptions struct {
	SourceSite      *Site
	TargetProvider  Provider
	TargetSiteName  string
	IncludeDatabase bool
	IncludeFiles    bool
	DryRun          bool
}

// MigrateSite migrates a site from the current provider to another provider
func (m *Manager) MigrateSite(options MigrateOptions) error {
	if m.currentProvider == nil {
		return fmt.Errorf("no source provider configured")
	}

	if options.TargetProvider == nil {
		return fmt.Errorf("target provider is required")
	}

	if options.SourceSite == nil {
		return fmt.Errorf("source site is required")
	}

	// Check if target provider supports migration interface
	migrator, ok := options.TargetProvider.(Migrator)
	if ok {
		// Use provider's built-in migration
		migOpts := MigrationOptions{
			IncludeDatabase: options.IncludeDatabase,
			IncludeFiles:    options.IncludeFiles,
			DryRun:          options.DryRun,
		}
		return migrator.ImportFromProvider(m.currentProvider, options.SourceSite, migOpts)
	}

	// Manual migration
	return m.manualMigration(options)
}

// manualMigration performs a manual migration between providers
func (m *Manager) manualMigration(options MigrateOptions) error {
	if options.DryRun {
		fmt.Printf("Dry run: Would migrate site %s from %s to %s\n",
			options.SourceSite.Name,
			m.currentProvider.Name(),
			options.TargetProvider.Name())
		return nil
	}

	// TODO: Implement manual migration logic
	// This would involve:
	// 1. Export database from source
	// 2. Import database to target
	// 3. Sync files from source to local
	// 4. Upload files to target

	return fmt.Errorf("manual migration not yet implemented")
}

// CompareProviders compares capabilities between two providers
func CompareProviders(provider1Name, provider2Name string) (*ProviderComparison, error) {
	p1, err := GetProvider(provider1Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s: %w", provider1Name, err)
	}

	p2, err := GetProvider(provider2Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s: %w", provider2Name, err)
	}

	return &ProviderComparison{
		Provider1:      provider1Name,
		Provider2:      provider2Name,
		Capabilities1:  p1.Capabilities(),
		Capabilities2:  p2.Capabilities(),
		SharedFeatures: findSharedCapabilities(p1.Capabilities(), p2.Capabilities()),
	}, nil
}

// ProviderComparison contains comparison results between two providers
type ProviderComparison struct {
	Provider1      string               `json:"provider1"`
	Provider2      string               `json:"provider2"`
	Capabilities1  ProviderCapabilities `json:"capabilities1"`
	Capabilities2  ProviderCapabilities `json:"capabilities2"`
	SharedFeatures []string             `json:"shared_features"`
}

// findSharedCapabilities finds capabilities that both providers support
func findSharedCapabilities(caps1, caps2 ProviderCapabilities) []string {
	shared := []string{}

	if caps1.Authentication && caps2.Authentication {
		shared = append(shared, "authentication")
	}
	if caps1.SiteManagement && caps2.SiteManagement {
		shared = append(shared, "site_management")
	}
	if caps1.DatabaseExport && caps2.DatabaseExport {
		shared = append(shared, "database_export")
	}
	if caps1.DatabaseImport && caps2.DatabaseImport {
		shared = append(shared, "database_import")
	}
	if caps1.FileSync && caps2.FileSync {
		shared = append(shared, "file_sync")
	}
	if caps1.Deployment && caps2.Deployment {
		shared = append(shared, "deployment")
	}
	if caps1.Environments && caps2.Environments {
		shared = append(shared, "environments")
	}
	if caps1.Backups && caps2.Backups {
		shared = append(shared, "backups")
	}
	if caps1.RemoteExecution && caps2.RemoteExecution {
		shared = append(shared, "remote_execution")
	}
	if caps1.MediaManagement && caps2.MediaManagement {
		shared = append(shared, "media_management")
	}
	if caps1.SSHAccess && caps2.SSHAccess {
		shared = append(shared, "ssh_access")
	}
	if caps1.APIAccess && caps2.APIAccess {
		shared = append(shared, "api_access")
	}

	return shared
}

// GetProviderRecommendation recommends a provider based on requirements
func GetProviderRecommendation(requirements []string) (string, error) {
	providers := GetAllProviders()

	var bestMatch string
	var bestScore int

	for name, provider := range providers {
		caps := provider.Capabilities()
		score := 0

		for _, req := range requirements {
			if hasCapability(caps, req) {
				score++
			}
		}

		if score > bestScore {
			bestScore = score
			bestMatch = name
		}
	}

	if bestMatch == "" {
		return "", fmt.Errorf("no provider matches the requirements")
	}

	return bestMatch, nil
}

// hasCapability checks if capabilities struct has a specific capability
func hasCapability(caps ProviderCapabilities, capability string) bool {
	switch capability {
	case "authentication":
		return caps.Authentication
	case "site_management":
		return caps.SiteManagement
	case "database_export":
		return caps.DatabaseExport
	case "database_import":
		return caps.DatabaseImport
	case "file_sync":
		return caps.FileSync
	case "deployment":
		return caps.Deployment
	case "environments":
		return caps.Environments
	case "backups":
		return caps.Backups
	case "remote_execution":
		return caps.RemoteExecution
	case "media_management":
		return caps.MediaManagement
	case "ssh_access":
		return caps.SSHAccess
	case "api_access":
		return caps.APIAccess
	case "scaling":
		return caps.Scaling
	case "monitoring":
		return caps.Monitoring
	case "logging":
		return caps.Logging
	default:
		return false
	}
}

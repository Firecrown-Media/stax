package provider

import (
	"fmt"
	"os"
)

// ProviderConfig contains configuration for creating a provider instance
type ProviderConfig struct {
	Name        string            `json:"name"`        // Provider name
	Credentials map[string]string `json:"credentials"` // Provider-specific credentials
	Options     map[string]string `json:"options"`     // Provider-specific options
}

// NewProvider creates a new provider instance with the given configuration
func NewProvider(config ProviderConfig) (Provider, error) {
	if config.Name == "" {
		return nil, fmt.Errorf("provider name is required")
	}

	// Get provider from registry
	provider, err := GetProvider(config.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	// Authenticate if credentials provided
	if len(config.Credentials) > 0 {
		if err := provider.Authenticate(config.Credentials); err != nil {
			return nil, fmt.Errorf("authentication failed for provider %s: %w", config.Name, err)
		}
	}

	return provider, nil
}

// NewProviderFromName creates a provider instance by name without authentication
func NewProviderFromName(name string) (Provider, error) {
	if name == "" {
		return nil, fmt.Errorf("provider name is required")
	}

	return GetProvider(name)
}

// DetectProviderFromConfig attempts to detect the provider from configuration
// This checks for provider-specific configuration sections
func DetectProviderFromConfig(config map[string]interface{}) (string, error) {
	// Check if provider is explicitly specified
	if providerName, ok := config["name"].(string); ok && providerName != "" {
		return providerName, nil
	}

	// Try to detect from provider-specific config sections
	knownProviders := []string{"wpengine", "aws", "wordpress-vip", "local"}

	for _, name := range knownProviders {
		if _, ok := config[name]; ok {
			return name, nil
		}
	}

	// Default to configured default provider
	return GetDefaultProvider(), nil
}

// ResolveProvider resolves the provider name using priority order:
// 1. Explicit name parameter
// 2. Environment variable STAX_PROVIDER
// 3. Project configuration
// 4. Global configuration default
// 5. Built-in default (wpengine)
func ResolveProvider(explicitName string, projectConfig map[string]interface{}, globalDefaultProvider string) (string, error) {
	// Priority 1: Explicit name
	if explicitName != "" {
		if !ProviderExists(explicitName) {
			return "", fmt.Errorf("provider %s not found", explicitName)
		}
		return explicitName, nil
	}

	// Priority 2: Environment variable
	if envProvider := os.Getenv("STAX_PROVIDER"); envProvider != "" {
		if !ProviderExists(envProvider) {
			return "", fmt.Errorf("provider %s (from STAX_PROVIDER) not found", envProvider)
		}
		return envProvider, nil
	}

	// Priority 3: Project configuration
	if projectConfig != nil {
		if provider, err := DetectProviderFromConfig(projectConfig); err == nil && provider != "" {
			if !ProviderExists(provider) {
				return "", fmt.Errorf("provider %s (from project config) not found", provider)
			}
			return provider, nil
		}
	}

	// Priority 4: Global default
	if globalDefaultProvider != "" {
		if !ProviderExists(globalDefaultProvider) {
			return "", fmt.Errorf("provider %s (global default) not found", globalDefaultProvider)
		}
		return globalDefaultProvider, nil
	}

	// Priority 5: Built-in default
	defaultProvider := GetDefaultProvider()
	if !ProviderExists(defaultProvider) {
		return "", fmt.Errorf("default provider %s not found", defaultProvider)
	}

	return defaultProvider, nil
}

// CreateProviderFromResolution resolves and creates a provider in one step
func CreateProviderFromResolution(
	explicitName string,
	projectConfig map[string]interface{},
	globalDefaultProvider string,
	credentials map[string]string,
) (Provider, error) {
	// Resolve provider name
	providerName, err := ResolveProvider(explicitName, projectConfig, globalDefaultProvider)
	if err != nil {
		return nil, err
	}

	// Create provider
	config := ProviderConfig{
		Name:        providerName,
		Credentials: credentials,
	}

	return NewProvider(config)
}

// ValidateProviderConfig validates a provider configuration
func ValidateProviderConfig(config ProviderConfig) error {
	if config.Name == "" {
		return fmt.Errorf("provider name is required")
	}

	// Check provider exists
	if !ProviderExists(config.Name) {
		return fmt.Errorf("provider %s not found", config.Name)
	}

	// Get provider instance
	provider, err := GetProvider(config.Name)
	if err != nil {
		return err
	}

	// Validate credentials if provided
	if len(config.Credentials) > 0 {
		if err := provider.ValidateCredentials(config.Credentials); err != nil {
			return fmt.Errorf("invalid credentials for provider %s: %w", config.Name, err)
		}
	}

	return nil
}

// GetProviderCapabilities returns the capabilities for a provider by name
func GetProviderCapabilities(name string) (*ProviderCapabilities, error) {
	provider, err := GetProvider(name)
	if err != nil {
		return nil, err
	}

	caps := provider.Capabilities()
	return &caps, nil
}

// SupportsCapability checks if a provider supports a specific capability
func SupportsCapability(providerName, capability string) (bool, error) {
	caps, err := GetProviderCapabilities(providerName)
	if err != nil {
		return false, err
	}

	switch capability {
	case "authentication":
		return caps.Authentication, nil
	case "site_management":
		return caps.SiteManagement, nil
	case "database_export":
		return caps.DatabaseExport, nil
	case "database_import":
		return caps.DatabaseImport, nil
	case "file_sync":
		return caps.FileSync, nil
	case "deployment":
		return caps.Deployment, nil
	case "environments":
		return caps.Environments, nil
	case "backups":
		return caps.Backups, nil
	case "remote_execution":
		return caps.RemoteExecution, nil
	case "media_management":
		return caps.MediaManagement, nil
	case "ssh_access":
		return caps.SSHAccess, nil
	case "api_access":
		return caps.APIAccess, nil
	case "scaling":
		return caps.Scaling, nil
	case "monitoring":
		return caps.Monitoring, nil
	case "logging":
		return caps.Logging, nil
	default:
		return false, fmt.Errorf("unknown capability: %s", capability)
	}
}

// GetRequiredCredentials returns the required credential keys for a provider
// This is a helper function that attempts to validate with empty credentials
// to get meaningful error messages
func GetRequiredCredentials(providerName string) ([]string, error) {
	provider, err := GetProvider(providerName)
	if err != nil {
		return nil, err
	}

	// Try to validate with empty credentials
	err = provider.ValidateCredentials(map[string]string{})
	if err != nil {
		// Parse error message to extract required fields
		// This is provider-dependent, so we return a generic error
		return nil, fmt.Errorf("provider %s requires credentials: %w", providerName, err)
	}

	return []string{}, nil
}

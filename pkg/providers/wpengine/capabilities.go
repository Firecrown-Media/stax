package wpengine

import "github.com/firecrown-media/stax/pkg/provider"

// GetWPEngineCapabilities returns the capabilities of the WPEngine provider
func GetWPEngineCapabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{
		// Core capabilities
		Authentication: true,
		SiteManagement: true,
		DatabaseExport: true,
		DatabaseImport: false, // WPEngine doesn't allow direct DB imports via API
		FileSync:       true,

		// Optional capabilities
		Deployment:      false, // Git deployments available but not yet implemented
		Environments:    true,  // Production/Staging environments
		Backups:         true,  // Point-in-time backups
		RemoteExecution: true,  // SSH access with WP-CLI
		MediaManagement: true,  // BunnyCDN integration
		SSHAccess:       true,  // SSH gateway access
		APIAccess:       true,  // WPEngine API v1

		// Advanced capabilities
		Scaling:    false, // Managed by WPEngine, not user-controlled
		Monitoring: false, // Available in portal but no API access
		Logging:    false, // Available in portal but no API access
	}
}

// WPEngineFeatures lists specific features of WPEngine hosting
var WPEngineFeatures = []string{
	"SSL Certificates (Let's Encrypt)",
	"CDN (BunnyCDN)",
	"Automatic Backups (Daily)",
	"Point-in-time Restore",
	"Staging Environment",
	"Git Push Deployments",
	"SSH Gateway Access",
	"WP-CLI Pre-installed",
	"PHP Version Management",
	"Object Caching (Redis)",
	"Global Edge Security",
	"Smart Plugin Manager",
	"Automated WordPress Updates",
}

// WPEngineLimitations lists known limitations of the WPEngine provider
var WPEngineLimitations = []string{
	"No direct database import via API/SSH",
	"Read-only filesystem (except Git deployments)",
	"No direct file upload via SSH",
	"No backup download via API",
	"Backup restoration requires WPEngine portal",
	"Limited environment variables access",
	"Cannot modify server configuration",
	"No root/sudo access",
}

// WPEngineAPIEndpoints lists available API endpoints
var WPEngineAPIEndpoints = map[string]string{
	"list_installs":     "/installs",
	"get_install":       "/installs/{id}",
	"list_backups":      "/installs/{id}/backups",
	"create_backup":     "/installs/{id}/backups",
	"list_domains":      "/installs/{id}/domains",
	"get_install_stats": "/installs/{id}/stats",
}

// WPEngineSSHGateway is the default SSH gateway
const WPEngineSSHGateway = "ssh.wpengine.net"

// WPEngineRequiredCredentials lists required credentials for WPEngine
var WPEngineRequiredCredentials = []string{
	"api_user",     // WPEngine API username
	"api_password", // WPEngine API password
	"install",      // Installation name
}

// WPEngineOptionalCredentials lists optional credentials for WPEngine
var WPEngineOptionalCredentials = []string{
	"ssh_key",     // SSH private key for SSH operations
	"ssh_gateway", // SSH gateway hostname (defaults to ssh.wpengine.net)
}

// GetWPEngineDefaultExclusions returns default file exclusions for rsync
func GetWPEngineDefaultExclusions() []string {
	return []string{
		"*.log",
		"cache/",
		"object-cache.php", // WPEngine's managed object cache
		".DS_Store",
		"Thumbs.db",
		"*.tmp",
		"*.swp",
		"node_modules/",
		".git/",
		".gitignore",
		"wp-config.php", // Don't sync config
		".htaccess",     // WPEngine manages this
	}
}

// GetWPEngineEnvironments returns available environments
func GetWPEngineEnvironments() []string {
	return []string{
		"production",
		"staging",
		"development", // Some plans include dev environment
	}
}

// WPEnginePHPVersions lists supported PHP versions
var WPEnginePHPVersions = []string{
	"7.4",
	"8.0",
	"8.1",
	"8.2",
	"8.3",
}

// WPEngineSupportsFeature checks if a specific feature is supported
func WPEngineSupportsFeature(feature string) bool {
	supported := map[string]bool{
		"ssh_access":      true,
		"wp_cli":          true,
		"git_deployments": true,
		"staging":         true,
		"cdn":             true,
		"ssl":             true,
		"backups":         true,
		"redis":           true,
		"database_import": false,
		"file_upload":     false,
		"root_access":     false,
		"custom_php_ini":  false,
		"backup_download": false,
	}

	return supported[feature]
}

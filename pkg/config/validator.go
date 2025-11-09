package config

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult represents the result of configuration validation
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationError
}

// Validate validates the configuration
func Validate(cfg *Config) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}

	// Validate required fields
	validateRequired(cfg, result)

	// Validate field formats
	validateFormats(cfg, result)

	// Validate cross-field constraints
	validateConstraints(cfg, result)

	// Validate version compatibility
	validateVersions(cfg, result)

	// Set overall validity
	result.Valid = len(result.Errors) == 0

	return result
}

// validateRequired checks that required fields are present
func validateRequired(cfg *Config, result *ValidationResult) {
	if cfg.Project.Name == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "project.name",
			Message: "is required",
		})
	}

	if cfg.Project.Type == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "project.type",
			Message: "is required",
		})
	}

	if cfg.Project.Mode == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "project.mode",
			Message: "is required",
		})
	}

	if cfg.WPEngine.Install == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "wpengine.install",
			Message: "is required",
		})
	}

	if cfg.WPEngine.Environment == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "wpengine.environment",
			Message: "is required",
		})
	}

	if cfg.Network.Domain == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "network.domain",
			Message: "is required",
		})
	}
}

// validateFormats checks that field formats are valid
func validateFormats(cfg *Config, result *ValidationResult) {
	// Validate project name (alphanumeric, hyphens, underscores)
	if cfg.Project.Name != "" {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, cfg.Project.Name)
		if !matched {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "project.name",
				Message: "must contain only alphanumeric characters, hyphens, and underscores",
			})
		}
	}

	// Validate project type
	validTypes := []string{"wordpress", "wordpress-multisite"}
	if !contains(validTypes, cfg.Project.Type) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "project.type",
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validTypes, ", ")),
		})
	}

	// Validate project mode
	validModes := []string{"subdomain", "subdirectory", "single"}
	if !contains(validModes, cfg.Project.Mode) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "project.mode",
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validModes, ", ")),
		})
	}

	// Validate WPEngine environment
	validEnvs := []string{"production", "staging", "development"}
	if !contains(validEnvs, cfg.WPEngine.Environment) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "wpengine.environment",
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validEnvs, ", ")),
		})
	}

	// Validate domain format
	if cfg.Network.Domain != "" {
		if !isValidDomain(cfg.Network.Domain) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "network.domain",
				Message: "is not a valid domain format",
			})
		}
	}

	// Validate site domains
	for i, site := range cfg.Network.Sites {
		if !isValidDomain(site.Domain) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("network.sites[%d].domain", i),
				Message: "is not a valid domain format",
			})
		}
	}

	// Validate PHP version
	validPHPVersions := []string{"7.4", "8.0", "8.1", "8.2", "8.3"}
	if !contains(validPHPVersions, cfg.DDEV.PHPVersion) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "ddev.php_version",
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validPHPVersions, ", ")),
		})
	}
}

// validateConstraints checks cross-field constraints
func validateConstraints(cfg *Config, result *ValidationResult) {
	// If mode is subdomain, ensure site domains are subdomains of network domain
	if cfg.Project.Mode == "subdomain" {
		for i, site := range cfg.Network.Sites {
			if !strings.HasSuffix(site.Domain, "."+cfg.Network.Domain) {
				result.Errors = append(result.Errors, ValidationError{
					Field:   fmt.Sprintf("network.sites[%d].domain", i),
					Message: fmt.Sprintf("must be a subdomain of %s", cfg.Network.Domain),
				})
			}
		}
	}

	// Check for duplicate site domains
	domainMap := make(map[string]bool)
	for i, site := range cfg.Network.Sites {
		if domainMap[site.Domain] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("network.sites[%d].domain", i),
				Message: "duplicate domain",
			})
		}
		domainMap[site.Domain] = true
	}
}

// validateVersions checks version compatibility
func validateVersions(cfg *Config, result *ValidationResult) {
	// Check if PHP version is older than recommended
	if cfg.DDEV.PHPVersion == "7.4" {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:   "ddev.php_version",
			Message: "PHP 7.4 is end-of-life, consider upgrading to 8.1 or later",
		})
	}

	// Check if Xdebug is enabled (performance warning)
	if cfg.DDEV.XdebugEnabled {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:   "ddev.xdebug_enabled",
			Message: "Xdebug is enabled, which may impact performance",
		})
	}
}

// Helper functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isValidDomain(domain string) bool {
	// Simple domain validation regex
	matched, _ := regexp.MatchString(`^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`, domain)
	return matched
}

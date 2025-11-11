package errors

// NewConfigNotFoundError creates an error for missing configuration file
func NewConfigNotFoundError(path string, err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeConfigNotFound,
		Message: "Configuration file not found",
		Details: "Stax requires a .stax.yml configuration file in the project directory. This file defines your WordPress project settings, WPEngine connection details, and development environment configuration.",
		Solutions: []Solution{
			{
				Description: "Create a new .stax.yml configuration file",
				Command:     "stax config template > .stax.yml",
				Steps: []string{
					"Generate a template config file",
					"Edit .stax.yml and update the values for your project",
					"Run 'stax config validate' to check your configuration",
				},
			},
			{
				Description: "Initialize a new project",
				Command:     "stax init",
				Steps: []string{
					"Run stax init to interactively create configuration",
					"Follow the prompts to set up your project",
				},
			},
		},
		DocsURL: "https://github.com/Firecrown-Media/stax#configuration",
		Err:     err,
	}
}

// NewCredentialsNotFoundError creates an error for missing WPEngine credentials
func NewCredentialsNotFoundError(err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeCredentialsNotFound,
		Message: "WPEngine credentials not configured",
		Details: "Stax requires WPEngine API credentials to sync databases and files. These credentials are securely stored in your macOS Keychain.",
		Solutions: []Solution{
			{
				Description: "Configure WPEngine credentials using the setup command",
				Command:     "stax setup wpengine",
				Steps: []string{
					"Run the setup command",
					"Enter your WPEngine username and password when prompted",
					"Credentials will be stored securely in macOS Keychain",
				},
			},
			{
				Description: "Check credential status",
				Command:     "stax setup --check",
				Steps: []string{
					"Run the check command to see credential diagnostics",
					"Follow the suggestions to fix any issues",
				},
			},
		},
		DocsURL: "https://github.com/Firecrown-Media/stax#wpengine-setup",
		Err:     err,
	}
}

// NewSSHKeyNotFoundError creates an error for missing SSH key
func NewSSHKeyNotFoundError(keyPath string, err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeSSHKeyNotFound,
		Message: "SSH key not found or not configured",
		Details: "Stax requires an SSH key to connect to WPEngine servers. The key should be located at ~/.ssh/id_rsa or another location specified in your configuration.",
		Solutions: []Solution{
			{
				Description: "Generate a new SSH key",
				Command:     "ssh-keygen -t rsa -b 4096 -C \"your_email@example.com\"",
				Steps: []string{
					"Run the ssh-keygen command",
					"Press Enter to accept the default location",
					"Add the public key to your WPEngine account",
				},
			},
			{
				Description: "Add existing SSH key to WPEngine",
				Steps: []string{
					"Copy your public key: cat ~/.ssh/id_rsa.pub",
					"Log in to WPEngine User Portal",
					"Navigate to SSH Keys settings",
					"Add your public key",
				},
			},
			{
				Description: "Specify a different SSH key path in .stax.yml",
				Steps: []string{
					"Edit your .stax.yml file",
					"Add or update the ssh_key_path setting",
					"Point it to your SSH key location",
				},
			},
		},
		DocsURL: "https://github.com/Firecrown-Media/stax#ssh-setup",
		Err:     err,
	}
}

// NewDDEVNotInstalledError creates an error for missing DDEV installation
func NewDDEVNotInstalledError(err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeDDEVNotInstalled,
		Message: "DDEV is not installed",
		Details: "Stax uses DDEV for local WordPress development. DDEV manages Docker containers for PHP, MySQL, and other services required for WordPress.",
		Solutions: []Solution{
			{
				Description: "Install DDEV using Homebrew (recommended for macOS)",
				Command:     "brew install ddev/ddev/ddev",
				Steps: []string{
					"Install DDEV via Homebrew",
					"Run 'mkcert -install' to set up SSL certificates",
					"Verify installation with 'ddev version'",
				},
			},
			{
				Description: "Install DDEV manually",
				Steps: []string{
					"Download DDEV from https://ddev.readthedocs.io/",
					"Follow the installation instructions for your OS",
					"Ensure Docker Desktop is installed and running",
					"Run 'ddev version' to verify installation",
				},
			},
		},
		DocsURL: "https://ddev.readthedocs.io/en/stable/#installation",
		Err:     err,
	}
}

// NewDDEVNotConfiguredError creates an error for missing DDEV configuration
func NewDDEVNotConfiguredError(projectDir string, err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeDDEVNotConfigured,
		Message: "DDEV is not configured for this project",
		Details: "This project directory does not have a DDEV configuration. A .ddev/config.yaml file is required to define the development environment settings.",
		Solutions: []Solution{
			{
				Description: "Initialize DDEV configuration with stax init",
				Command:     "stax init",
				Steps: []string{
					"Run stax init to set up your project",
					"DDEV will be configured automatically",
					"Database and files can be synced from WPEngine",
				},
			},
			{
				Description: "Initialize DDEV configuration manually",
				Command:     "ddev config --project-type=wordpress --php-version=8.1",
				Steps: []string{
					"Navigate to your project directory",
					"Run the ddev config command",
					"Adjust PHP and MySQL versions as needed",
					"Run 'ddev start' to start the environment",
				},
			},
		},
		DocsURL: "https://github.com/Firecrown-Media/stax#manual-setup",
		Err:     err,
	}
}

// NewCommandNotImplementedError creates an error for unimplemented commands
func NewCommandNotImplementedError(command string, workaround string, steps []string) *EnhancedError {
	solutions := []Solution{}

	if workaround != "" {
		solutions = append(solutions, Solution{
			Description: "Use this workaround instead",
			Command:     workaround,
		})
	}

	if len(steps) > 0 {
		solutions = append(solutions, Solution{
			Description: "Manual setup instructions",
			Steps:       steps,
		})
	}

	solutions = append(solutions, Solution{
		Description: "Track implementation progress",
		Steps: []string{
			"This command is planned for a future release",
			"Check the GitHub repository for updates",
			"Consider contributing if you'd like to help implement it",
		},
	})

	return &EnhancedError{
		Code:      ErrCodeCommandNotImplemented,
		Message:   "Command not yet implemented",
		Details:   "The '" + command + "' command is currently under development. While Stax is evolving, you can use the workarounds below to achieve similar functionality.",
		Solutions: solutions,
		DocsURL:   "https://github.com/Firecrown-Media/stax#current-limitations",
		Err:       nil,
	}
}

// NewInvalidConfigError creates an error for invalid configuration
func NewInvalidConfigError(message string, err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeInvalidConfig,
		Message: "Invalid configuration",
		Details: message,
		Solutions: []Solution{
			{
				Description: "Validate your configuration file",
				Command:     "stax config validate",
				Steps: []string{
					"Run the validation command to see specific issues",
					"Fix each validation error",
					"Re-run validation until all errors are resolved",
				},
			},
			{
				Description: "Generate a fresh configuration template",
				Command:     "stax config template",
				Steps: []string{
					"View the template to see correct format",
					"Compare with your current configuration",
					"Update your config to match the template structure",
				},
			},
		},
		DocsURL: "https://github.com/Firecrown-Media/stax#configuration",
		Err:     err,
	}
}

// NewWPEngineAPIError creates an error for WPEngine API failures
func NewWPEngineAPIError(message string, err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeWPEngineAPI,
		Message: "WPEngine API error",
		Details: message,
		Solutions: []Solution{
			{
				Description: "Check your WPEngine credentials",
				Command:     "stax setup --check",
				Steps: []string{
					"Verify credentials are configured correctly",
					"Test API connectivity",
					"Re-run setup if needed",
				},
			},
			{
				Description: "Check WPEngine service status",
				Steps: []string{
					"Visit https://my.wpengine.com/ to check if services are running",
					"Check if your IP is whitelisted (if applicable)",
					"Try again in a few minutes if WPEngine is experiencing issues",
				},
			},
		},
		DocsURL: "https://github.com/Firecrown-Media/stax#troubleshooting",
		Err:     err,
	}
}

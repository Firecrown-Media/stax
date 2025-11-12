package cmd

import (
	"fmt"

	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "🚧 Configuration management",
	Long:  `Manage Stax configuration files and values.`,
}

var (
	configGlobal bool
	configJSON   bool
)

// configGetCmd represents the config:get command
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Example: `  # Get project config
  stax config get wpengine.environment

  # Get global config
  stax config get wpengine.api_user --global`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGet,
}

// configSetCmd represents the config:set command
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Example: `  # Set project config
  stax config set wpengine.environment staging

  # Set global config
  stax config set ddev.php_version 8.2 --global`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

// configListCmd represents the config:list command
var configListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all configuration values",
	Example: `  # List project config
  stax config list

  # List global config
  stax config list --global

  # List as JSON
  stax config list --json`,
	RunE: runConfigList,
}

// configValidateCmd represents the config:validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Example: `  # Validate config
  stax config validate`,
	RunE: runConfigValidate,
}

// configTemplateCmd represents the config:template command
var configTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Generate a configuration template",
	Example: `  # Generate template to stdout
  stax config template

  # Generate template to file
  stax config template > .stax.yml`,
	RunE: runConfigTemplate,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configTemplateCmd)

	// Flags for get
	configGetCmd.Flags().BoolVar(&configGlobal, "global", false, "get from global config")

	// Flags for set
	configSetCmd.Flags().BoolVar(&configGlobal, "global", false, "set in global config")

	// Flags for list
	configListCmd.Flags().BoolVar(&configGlobal, "global", false, "list global config")
	configListCmd.Flags().BoolVar(&configJSON, "json", false, "output as JSON")
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	// TODO: Load config from file
	// TODO: Get value by key path
	// TODO: Display value

	ui.Info(fmt.Sprintf("Getting config value: %s", key))
	ui.Info("Config get is not yet implemented")

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// TODO: Load config from file
	// TODO: Set value by key path
	// TODO: Save config file
	// TODO: Display confirmation

	ui.Info(fmt.Sprintf("Setting config: %s = %s", key, value))
	ui.Info("Config set is not yet implemented")

	return nil
}

func runConfigList(cmd *cobra.Command, args []string) error {
	// TODO: Load config from file
	// TODO: Format as YAML or JSON
	// TODO: Display config

	ui.PrintHeader("Configuration")
	ui.Info("Config list is not yet implemented")

	return nil
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Validating Configuration")

	// Load config from current directory
	cfg := getConfig()
	if cfg == nil {
		return fmt.Errorf("no configuration found - run 'stax config template' to create one")
	}

	ui.Info("Validating .stax.yml configuration...")

	// Basic validation checks
	errors := []string{}
	warnings := []string{}

	// Check required fields
	if cfg.Project.Name == "" {
		errors = append(errors, "project.name is required")
	}
	if cfg.Project.Type == "" {
		errors = append(errors, "project.type is required (wordpress or wordpress-multisite)")
	}
	if cfg.WPEngine.Install == "" {
		errors = append(errors, "wpengine.install is required")
	}

	// Check optional but recommended fields
	if cfg.Network.Domain == "" {
		warnings = append(warnings, "network.domain not set - recommended for multisite")
	}

	// Display results
	if len(errors) > 0 {
		ui.Error("Validation failed:")
		for _, err := range errors {
			ui.Info(fmt.Sprintf("  - %s", err))
		}
		return fmt.Errorf("configuration has %d error(s)", len(errors))
	}

	if len(warnings) > 0 {
		ui.Warning("Validation warnings:")
		for _, warn := range warnings {
			ui.Info(fmt.Sprintf("  - %s", warn))
		}
	}

	ui.Success("Configuration is valid!")
	return nil
}

func runConfigTemplate(cmd *cobra.Command, args []string) error {
	// Generate a template .stax.yml configuration
	template := `# Stax Configuration
# Version: 1

version: 1

# Project metadata
project:
  name: my-wordpress-project
  type: wordpress-multisite  # wordpress or wordpress-multisite
  mode: subdomain            # subdomain, subdirectory, or single
  description: My WordPress Multisite Project

# WPEngine integration
wpengine:
  install: myinstall           # WPEngine install name
  environment: production      # production, staging, or development
  account_name: myaccount      # WPEngine account name (optional)
  ssh_gateway: ssh.wpengine.net

  # Backup preferences
  backup:
    auto_snapshot: true        # Create snapshot before DB pull
    skip_logs: true           # Skip log tables
    skip_transients: true     # Skip transient data
    skip_spam: true           # Skip spam comments
    exclude_tables: []        # Additional tables to exclude

  # Domain mapping (optional)
  domains:
    production:
      primary: example.com
      sites:
        - site1.example.com
        - site2.example.com
    staging:
      primary: staging.example.com

# Network configuration (for multisite)
network:
  domain: example.local      # Local development domain
  title: My Network         # Network title
  admin_email: admin@example.local

  # Individual sites (optional - can be auto-detected)
  sites:
    - name: Main Site
      slug: main
      title: Main Site
      domain: example.local
      path: /
    - name: Site One
      slug: site1
      title: Site One
      domain: site1.example.local
      path: /

# DDEV configuration
ddev:
  name: ""                  # Auto-generated from project name
  type: wordpress
  php_version: "8.1"
  mysql_version: "8.0"
  webserver_type: nginx-fpm
  router_http_port: "80"
  router_https_port: "443"
  xdebug_enabled: false
  additional_hostnames: []
  additional_fqdns: []

# GitHub repository (optional)
repository:
  url: https://github.com/username/repo.git
  branch: main
  private: true

# Build process (optional)
build:
  enabled: true
  composer:
    enabled: true
    install_args: "--no-dev --optimize-autoloader"
  npm:
    enabled: true
    install_command: npm ci
    build_command: npm run build
    dev_command: npm run dev
    watch_command: npm run watch
  quality:
    phpcs: true
    phpcbf: true

# WordPress configuration (optional)
wordpress:
  version: latest
  locale: en_US
  timezone: America/New_York
  debug: true
  debug_log: true

# Remote media configuration
media:
  enabled: true
  strategy: bunnycdn         # bunnycdn, wpengine, or disabled
  cache_ttl: 2592000         # 30 days
  cache_max_size: 10737418240  # 10GB

# Logging and debugging (optional)
logging:
  level: info               # debug, info, warning, error
  file: .stax/logs/stax.log
  rotate: true

# Snapshots (optional)
snapshots:
  directory: .stax/snapshots
  retention_days: 30
  auto_cleanup: true

# Performance tuning (optional)
performance:
  db_import_chunk_size: 1048576  # 1MB
  parallel_downloads: 4
  connection_timeout: 30
`

	fmt.Println(template)
	return nil
}

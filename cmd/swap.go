package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Firecrown-Media/stax/pkg/config"
	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var swapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Hot swap environment components (PHP, WordPress, MySQL, etc.)",
	Long: `Hot swap allows you to quickly change versions of core components in your 
development environment without losing data or requiring full environment recreation.

Examples:
  stax swap php 8.3                    # Switch to PHP 8.3
  stax swap mysql 8.4                  # Switch to MySQL 8.4
  stax swap wordpress 6.4              # Switch to WordPress 6.4
  stax swap preset php8.3-wp6.4        # Use predefined preset
  stax swap --rollback                 # Rollback to previous configuration`,
	RunE: swapRun,
}

var swapListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available versions and presets",
	Long:  `List all available versions for each component and predefined presets.`,
	RunE:  swapListRun,
}

var swapStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current environment configuration",
	Long:  `Display the current versions of all components and swap history.`,
	RunE:  swapStatusRun,
}

var (
	rollback bool
	force    bool
)

func init() {
	rootCmd.AddCommand(swapCmd)
	swapCmd.AddCommand(swapListCmd)
	swapCmd.AddCommand(swapStatusCmd)

	swapCmd.Flags().BoolVar(&rollback, "rollback", false, "rollback to previous configuration")
	swapCmd.Flags().BoolVar(&force, "force", false, "force swap without confirmation prompts")
}

func swapRun(cmd *cobra.Command, args []string) error {
	if rollback {
		return rollbackEnvironment()
	}

	if len(args) < 2 {
		return fmt.Errorf("usage: stax swap <component> <version> or stax swap preset <preset-name>")
	}

	component := args[0]
	version := args[1]

	// Handle preset swaps
	if component == "preset" {
		return swapPreset(version)
	}

	// Handle individual component swaps
	return swapComponent(component, version)
}

func swapComponent(component, version string) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Load current config
	cfg, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Save current config as backup before changing
	if err := saveSwapBackup(projectPath, cfg); err != nil {
		return fmt.Errorf("failed to save backup: %w", err)
	}

	// Update configuration based on component
	switch strings.ToLower(component) {
	case "php":
		if !isValidPHPVersion(version) {
			return fmt.Errorf("invalid PHP version: %s", version)
		}
		cfg.PHPVersion = version
	case "mysql", "database", "db":
		if !isValidDatabaseVersion(version) {
			return fmt.Errorf("invalid database version: %s", version)
		}
		cfg.Database = "mysql:" + version
	case "webserver", "web":
		if !isValidWebServer(version) {
			return fmt.Errorf("invalid webserver: %s", version)
		}
		cfg.WebServer = version
	default:
		return fmt.Errorf("unsupported component: %s. Supported: php, mysql, webserver", component)
	}

	// Apply the configuration change
	if err := applyEnvironmentChange(projectPath, cfg); err != nil {
		return fmt.Errorf("failed to apply environment change: %w", err)
	}

	fmt.Printf("✅ Successfully swapped %s to %s\n", component, version)
	if viper.GetBool("verbose") {
		fmt.Printf("Environment will restart to apply changes\n")
	}

	return nil
}

func swapPreset(presetName string) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Load current config
	cfg, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Save current config as backup
	if err := saveSwapBackup(projectPath, cfg); err != nil {
		return fmt.Errorf("failed to save backup: %w", err)
	}

	// Apply preset configuration
	preset, err := getPresetConfig(presetName)
	if err != nil {
		return fmt.Errorf("failed to get preset config: %w", err)
	}

	// Update config with preset values
	if preset.PHPVersion != "" {
		cfg.PHPVersion = preset.PHPVersion
	}
	if preset.Database != "" {
		cfg.Database = preset.Database
	}
	if preset.WebServer != "" {
		cfg.WebServer = preset.WebServer
	}

	// Apply the configuration change
	if err := applyEnvironmentChange(projectPath, cfg); err != nil {
		return fmt.Errorf("failed to apply preset: %w", err)
	}

	fmt.Printf("✅ Successfully applied preset: %s\n", presetName)
	return nil
}

func rollbackEnvironment() error {
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Load backup configuration
	backupCfg, err := loadSwapBackup(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load backup config: %w", err)
	}

	// Apply the backup configuration
	if err := applyEnvironmentChange(projectPath, backupCfg); err != nil {
		return fmt.Errorf("failed to rollback environment: %w", err)
	}

	fmt.Println("✅ Successfully rolled back to previous configuration")
	return nil
}

func applyEnvironmentChange(projectPath string, cfg *config.ProjectConfig) error {
	// Stop the current environment
	if ddev.IsProject(projectPath) {
		if viper.GetBool("verbose") {
			fmt.Println("Stopping current environment...")
		}
		if err := ddev.Stop(projectPath); err != nil {
			return fmt.Errorf("failed to stop environment: %w", err)
		}
	}

	// Update DDEV configuration
	if err := ddev.UpdateConfig(projectPath, ddev.Config{
		ProjectName:  cfg.Name,
		ProjectType:  cfg.Type,
		PHPVersion:   cfg.PHPVersion,
		WebServer:    cfg.WebServer,
		DatabaseType: cfg.Database,
	}); err != nil {
		return fmt.Errorf("failed to update DDEV config: %w", err)
	}

	// Save the updated project configuration
	if err := cfg.Save(projectPath); err != nil {
		return fmt.Errorf("failed to save project config: %w", err)
	}

	// Start the environment with new configuration
	if viper.GetBool("verbose") {
		fmt.Println("Starting environment with new configuration...")
	}
	if err := ddev.Start(projectPath); err != nil {
		return fmt.Errorf("failed to start environment: %w", err)
	}

	return nil
}

func swapListRun(cmd *cobra.Command, args []string) error {
	fmt.Println("Available versions and presets:")
	fmt.Println()

	fmt.Println("📱 PHP Versions:")
	phpVersions := []string{"7.4", "8.0", "8.1", "8.2", "8.3", "8.4"}
	for _, v := range phpVersions {
		fmt.Printf("  • php %s\n", v)
	}
	fmt.Println()

	fmt.Println("🗄️  MySQL Versions:")
	mysqlVersions := []string{"5.7", "8.0", "8.4"}
	for _, v := range mysqlVersions {
		fmt.Printf("  • mysql %s\n", v)
	}
	fmt.Println()

	fmt.Println("🌐 Web Servers:")
	webServers := []string{"nginx-fpm", "apache-fpm", "nginx-fpm-arm64"}
	for _, v := range webServers {
		fmt.Printf("  • webserver %s\n", v)
	}
	fmt.Println()

	fmt.Println("🎯 Presets:")
	presets := getAvailablePresets()
	for name, desc := range presets {
		fmt.Printf("  • preset %-20s %s\n", name, desc)
	}

	return nil
}

func swapStatusRun(cmd *cobra.Command, args []string) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	cfg, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	fmt.Printf("Current Environment Configuration:\n")
	fmt.Printf("  Project:    %s\n", cfg.Name)
	fmt.Printf("  PHP:        %s\n", cfg.PHPVersion)
	fmt.Printf("  Database:   %s\n", cfg.Database)
	fmt.Printf("  WebServer:  %s\n", cfg.WebServer)
	fmt.Println()

	// Check if backup exists
	if hasSwapBackup(projectPath) {
		fmt.Println("💾 Backup configuration available (use --rollback to restore)")
	} else {
		fmt.Println("ℹ️  No backup configuration found")
	}

	return nil
}

// Validation functions
func isValidPHPVersion(version string) bool {
	validVersions := []string{"7.4", "8.0", "8.1", "8.2", "8.3", "8.4"}
	for _, v := range validVersions {
		if v == version {
			return true
		}
	}
	return false
}

func isValidDatabaseVersion(version string) bool {
	validVersions := []string{"5.7", "8.0", "8.4"}
	for _, v := range validVersions {
		if v == version {
			return true
		}
	}
	return false
}

func isValidWebServer(server string) bool {
	validServers := []string{"nginx-fpm", "apache-fpm", "nginx-fpm-arm64"}
	for _, v := range validServers {
		if v == server {
			return true
		}
	}
	return false
}
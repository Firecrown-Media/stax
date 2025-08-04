package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Firecrown-Media/stax/pkg/config"
	"github.com/spf13/viper"
)

// PresetConfig defines a configuration preset
type PresetConfig struct {
	PHPVersion  string `yaml:"php_version"`
	Database    string `yaml:"database"`
	WebServer   string `yaml:"webserver"`
	Description string `yaml:"description"`
}

// getAvailablePresets returns a map of preset names to descriptions
func getAvailablePresets() map[string]string {
	return map[string]string{
		"legacy":          "PHP 7.4, MySQL 5.7 for older WordPress sites",
		"stable":          "PHP 8.2, MySQL 8.0 - a recommended stable configuration",
		"modern":          "PHP 8.3, MySQL 8.4 for latest WordPress features",
		"bleeding-edge":   "PHP 8.4, MySQL 8.4 for testing new features",
		"performance":     "PHP 8.3, MySQL 8.4, nginx-fpm optimized for speed",
		"compatibility":   "PHP 8.1, MySQL 8.0 for maximum plugin compatibility",
		"development":     "PHP 8.3, MySQL 8.0 with enhanced debugging",
	}
}

// getPresetConfig returns the configuration for a named preset
func getPresetConfig(presetName string) (*PresetConfig, error) {
	presets := map[string]PresetConfig{
		"legacy": {
			PHPVersion:  "7.4",
			Database:    "mysql:5.7",
			WebServer:   "nginx-fpm",
			Description: "PHP 7.4, MySQL 5.7 for older WordPress sites",
		},
		"stable": {
			PHPVersion:  "8.2",
			Database:    "mysql:8.0",
			WebServer:   "nginx-fpm",
			Description: "PHP 8.2, MySQL 8.0 - a recommended stable configuration",
		},
		"modern": {
			PHPVersion:  "8.3",
			Database:    "mysql:8.4",
			WebServer:   "nginx-fpm",
			Description: "PHP 8.3, MySQL 8.4 for latest WordPress features",
		},
		"bleeding-edge": {
			PHPVersion:  "8.4",
			Database:    "mysql:8.4",
			WebServer:   "nginx-fpm",
			Description: "PHP 8.4, MySQL 8.4 for testing new features",
		},
		"performance": {
			PHPVersion:  "8.3",
			Database:    "mysql:8.4",
			WebServer:   "nginx-fpm",
			Description: "PHP 8.3, MySQL 8.4, nginx-fpm optimized for speed",
		},
		"compatibility": {
			PHPVersion:  "8.1",
			Database:    "mysql:8.0",
			WebServer:   "nginx-fpm",
			Description: "PHP 8.1, MySQL 8.0 for maximum plugin compatibility",
		},
		"development": {
			PHPVersion:  "8.3",
			Database:    "mysql:8.0",
			WebServer:   "nginx-fpm",
			Description: "PHP 8.3, MySQL 8.0 with enhanced debugging",
		},
	}

	preset, exists := presets[presetName]
	if !exists {
		return nil, fmt.Errorf("preset '%s' not found", presetName)
	}

	return &preset, nil
}

// saveSwapBackup saves the current configuration as a backup
func saveSwapBackup(projectPath string, cfg *config.ProjectConfig) error {
	backupPath := filepath.Join(projectPath, ".stax.backup.yaml")
	
	v := viper.New()
	v.SetConfigFile(backupPath)
	v.SetConfigType("yaml")
	
	// Set all values
	v.Set("name", cfg.Name)
	v.Set("type", cfg.Type)
	v.Set("php_version", cfg.PHPVersion)
	v.Set("webserver", cfg.WebServer)
	v.Set("database", cfg.Database)
	v.Set("wordpress", cfg.WordPress)
	v.Set("wpengine", cfg.WPEngine)
	v.Set("plugins", cfg.Plugins)
	v.Set("environment", cfg.Environment)
	
	return v.WriteConfig()
}

// loadSwapBackup loads the backup configuration
func loadSwapBackup(projectPath string) (*config.ProjectConfig, error) {
	backupPath := filepath.Join(projectPath, ".stax.backup.yaml")
	
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no backup configuration found")
	}
	
	v := viper.New()
	v.SetConfigFile(backupPath)
	v.SetConfigType("yaml")
	
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read backup config: %w", err)
	}
	
	var cfg config.ProjectConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup config: %w", err)
	}
	
	return &cfg, nil
}

// hasSwapBackup checks if a backup configuration exists
func hasSwapBackup(projectPath string) bool {
	backupPath := filepath.Join(projectPath, ".stax.backup.yaml")
	_, err := os.Stat(backupPath)
	return err == nil
}
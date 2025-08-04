package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ProjectConfig struct {
	Name         string            `mapstructure:"name" yaml:"name" json:"name"`
	Type         string            `mapstructure:"type" yaml:"type" json:"type"`
	PHPVersion   string            `mapstructure:"php_version" yaml:"php_version" json:"php_version"`
	WebServer    string            `mapstructure:"webserver" yaml:"webserver" json:"webserver"`
	Database     string            `mapstructure:"database" yaml:"database" json:"database"`
	WordPress    WordPressConfig   `mapstructure:"wordpress" yaml:"wordpress" json:"wordpress,omitempty"`
	WPEngine     WPEngineConfig    `mapstructure:"wpengine" yaml:"wpengine" json:"wpengine,omitempty"`
	Plugins      []string          `mapstructure:"plugins" yaml:"plugins" json:"plugins,omitempty"`
	Environment  map[string]string `mapstructure:"environment" yaml:"environment" json:"environment,omitempty"`
}

type WordPressConfig struct {
	URL       string `mapstructure:"url" yaml:"url" json:"url"`
	Title     string `mapstructure:"title" yaml:"title" json:"title"`
	AdminUser string `mapstructure:"admin_user" yaml:"admin_user" json:"admin_user"`
	AdminEmail string `mapstructure:"admin_email" yaml:"admin_email" json:"admin_email"`
	Theme     string `mapstructure:"theme" yaml:"theme" json:"theme,omitempty"`
}

type WPEngineConfig struct {
	InstallName    string            `mapstructure:"install_name" yaml:"install_name" json:"install_name"`
	Environment    string            `mapstructure:"environment" yaml:"environment" json:"environment"`
	Username       string            `mapstructure:"username" yaml:"username" json:"username,omitempty"`
	RemoteMediaURL string            `mapstructure:"remote_media_url" yaml:"remote_media_url" json:"remote_media_url,omitempty"`
	SyncOptions    WPEngineSyncConfig `mapstructure:"sync" yaml:"sync" json:"sync,omitempty"`
}

type WPEngineSyncConfig struct {
	SkipMedia     bool     `mapstructure:"skip_media" yaml:"skip_media" json:"skip_media"`
	SkipPlugins   bool     `mapstructure:"skip_plugins" yaml:"skip_plugins" json:"skip_plugins"`
	SkipThemes    bool     `mapstructure:"skip_themes" yaml:"skip_themes" json:"skip_themes"`
	ExcludeDirs   []string `mapstructure:"exclude_dirs" yaml:"exclude_dirs" json:"exclude_dirs,omitempty"`
	PreservePaths []string `mapstructure:"preserve_paths" yaml:"preserve_paths" json:"preserve_paths,omitempty"`
}

const ConfigFileName = "stax"

func Load(projectPath string) (*ProjectConfig, error) {
	v := viper.New()
	v.SetConfigName(ConfigFileName)
	v.SetConfigType("yaml")
	v.AddConfigPath(projectPath)
	
	// Set defaults
	setDefaults(v)
	
	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, return defaults
			var config ProjectConfig
			if err := v.Unmarshal(&config); err != nil {
				return nil, fmt.Errorf("failed to unmarshal default config: %w", err)
			}
			return &config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ProjectConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *ProjectConfig) Save(projectPath string) error {
	v := viper.New()
	v.SetConfigName(ConfigFileName)
	v.SetConfigType("yaml")
	v.AddConfigPath(projectPath)
	
	// Set all values
	v.Set("name", c.Name)
	v.Set("type", c.Type)
	v.Set("php_version", c.PHPVersion)
	v.Set("webserver", c.WebServer)
	v.Set("database", c.Database)
	v.Set("wordpress", c.WordPress)
	v.Set("wpengine", c.WPEngine)
	v.Set("plugins", c.Plugins)
	v.Set("environment", c.Environment)
	
	configPath := filepath.Join(projectPath, ConfigFileName+".yaml")
	return v.WriteConfigAs(configPath)
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("type", "wordpress")
	v.SetDefault("php_version", "8.2")
	v.SetDefault("webserver", "nginx-fpm")
	v.SetDefault("database", "mysql:8.0")
	v.SetDefault("wordpress.url", "https://localhost")
	v.SetDefault("wordpress.title", "My WordPress Site")
	v.SetDefault("wordpress.admin_user", "admin")
	v.SetDefault("wordpress.admin_email", "admin@localhost.local")
}

func Exists(projectPath string) bool {
	configPath := filepath.Join(projectPath, ConfigFileName+".yaml")
	_, err := os.Stat(configPath)
	return err == nil
}
package wpengine

import "time"

type Config struct {
	Username    string `mapstructure:"username" yaml:"username" json:"username"`
	APIKey      string `mapstructure:"api_key" yaml:"api_key" json:"api_key,omitempty"`
	InstallName string `mapstructure:"install_name" yaml:"install_name" json:"install_name"`
	Environment string `mapstructure:"environment" yaml:"environment" json:"environment"` // production, staging, development
	SSHKey      string `mapstructure:"ssh_key" yaml:"ssh_key" json:"ssh_key,omitempty"`
}

type InstallInfo struct {
	Name         string    `json:"name"`
	Environment  string    `json:"environment"`
	PHPVersion   string    `json:"php_version"`
	Status       string    `json:"status"`
	Domain       string    `json:"domain"`
	CreatedAt    time.Time `json:"created_at"`
	LastBackup   time.Time `json:"last_backup"`
	DatabaseName string    `json:"database_name"`
}

type BackupInfo struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // database, files, full
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	DownloadURL string    `json:"download_url,omitempty"`
}

type SyncOptions struct {
	SkipMedia     bool     `json:"skip_media"`
	SkipPlugins   bool     `json:"skip_plugins"`
	SkipThemes    bool     `json:"skip_themes"`
	ExcludeDirs   []string `json:"exclude_dirs"`
	IncludeDirs   []string `json:"include_dirs"`
	PreservePaths []string `json:"preserve_paths"`
	DeleteLocal   bool     `json:"delete_local"`   // Allow deletion of local files not on remote
	SuppressDebug bool     `json:"suppress_debug"` // Suppress WordPress debug notices and warnings
}

type DatabaseSyncResult struct {
	Success       bool   `json:"success"`
	BackupID      string `json:"backup_id"`
	DatabaseFile  string `json:"database_file"`
	ImportedRows  int    `json:"imported_rows"`
	RewrittenURLs int    `json:"rewritten_urls"`
	Error         string `json:"error,omitempty"`
}

type FilesSyncResult struct {
	Success       bool     `json:"success"`
	SyncedFiles   int      `json:"synced_files"`
	SkippedFiles  int      `json:"skipped_files"`
	TotalSize     int64    `json:"total_size"`
	ExcludedPaths []string `json:"excluded_paths"`
	Error         string   `json:"error,omitempty"`
}
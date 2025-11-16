package ddev

import "time"

// DDEVConfig represents the DDEV configuration structure
type DDEVConfig struct {
	Name                  string         `yaml:"name"`
	Type                  string         `yaml:"type"`
	DocRoot               string         `yaml:"docroot"`
	PHPVersion            string         `yaml:"php_version"`
	WebImageExtraPackages []string       `yaml:"webimage_extra_packages,omitempty"`
	Database              DatabaseConfig `yaml:"database"`
	AdditionalHostnames   []string       `yaml:"additional_hostnames,omitempty"`
	AdditionalFQDNs       []string       `yaml:"additional_fqdns,omitempty"`
	RouterHTTPPort        string         `yaml:"router_http_port"`
	RouterHTTPSPort       string         `yaml:"router_https_port"`
	XdebugEnabled         bool           `yaml:"xdebug_enabled"`
	UseDNSWhenPossible    bool           `yaml:"use_dns_when_possible"`
	ComposerVersion       string         `yaml:"composer_version"`
	WebEnvironment        []string       `yaml:"web_environment,omitempty"`
	NodeJSVersion         string         `yaml:"nodejs_version,omitempty"`
	Hooks                 *DDEVHooks     `yaml:"hooks,omitempty"`
	MutagenEnabled        bool           `yaml:"mutagen_enabled"`
	PerformanceMode       string         `yaml:"performance_mode,omitempty"`
}

// DatabaseConfig represents DDEV database configuration
type DatabaseConfig struct {
	Type    string `yaml:"type"`
	Version string `yaml:"version"`
}

// DDEVHooks represents DDEV lifecycle hooks
type DDEVHooks struct {
	PostStart  []HookCommand `yaml:"post-start,omitempty"`
	PreStop    []HookCommand `yaml:"pre-stop,omitempty"`
	PostImport []HookCommand `yaml:"post-import,omitempty"`
}

// HookCommand represents a single hook command
type HookCommand struct {
	Exec     string `yaml:"exec,omitempty"`
	ExecHost string `yaml:"exec-host,omitempty"`
}

// DDEVStatus represents the current status of a DDEV project
type DDEVStatus struct {
	State       string // running, stopped, paused, etc.
	ProjectName string
	Type        string
	Location    string
	URLs        []string // All URLs (http, https, *.ddev.site)
	Services    []ServiceStatus
	PHPVersion  string
	DBVersion   string
	RouterHTTP  string
	RouterHTTPS string
}

// ServiceStatus represents the status of a single DDEV service
type ServiceStatus struct {
	Name   string
	State  string // running, stopped, etc.
	Ports  []string
	Image  string
	Health string // healthy, unhealthy, starting, etc.
}

// ConfigOptions contains options for generating DDEV configuration
type ConfigOptions struct {
	ProjectName           string
	DocRoot               string
	Type                  string // php, wordpress, etc.
	PHPVersion            string
	DatabaseType          string // mysql, mariadb, postgres
	DatabaseVersion       string
	AdditionalHostnames   []string // For multisite subdomains
	AdditionalFQDNs       []string // Full domain names
	RouterHTTPPort        string
	RouterHTTPSPort       string
	MutagenEnabled        bool   // For Mac performance
	PerformanceMode       string // mutagen, nfs, none
	WebImageExtraPackages []string
	WebEnvironment        []string // Environment variables for web container
	NodeJSVersion         string
	ComposerVersion       string
	PostStartHooks        []string
	PostImportHooks       []string
	XdebugEnabled         bool
	UseDNSWhenPossible    bool // Use DNS for .ddev.site domains
}

// MediaProxyOptions contains options for Nginx media proxy configuration
type MediaProxyOptions struct {
	Enabled      bool
	CDNName      string
	CDNURL       string
	WPEngineURL  string
	WPEngineHost string
	CacheTTL     string // e.g., "30d", "24h"
	CacheMaxSize string // e.g., "10g", "1g"
	CacheEnabled bool
	ProxyHeaders map[string]string
}

// ProjectInfo represents detailed information about a DDEV project
type ProjectInfo struct {
	Name            string
	Type            string
	Location        string
	AppRoot         string // Alias for Location
	URLs            []string
	PrimaryURL      string // First URL in the URLs slice
	PHPVersion      string
	DatabaseType    string
	DatabaseVersion string
	RouterHTTPPort  string
	RouterHTTPSPort string
	Hostnames       []string
	Status          string
	Running         bool // Whether containers are running
	Healthy         bool // Whether containers are healthy
	Services        []ServiceStatus
	Router          string
	RouterStatus    string
	Webserver       string
	XdebugEnabled   bool
	MailhogURL      string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ExecOptions contains options for executing commands in DDEV container
type ExecOptions struct {
	Service     string   // web, db, etc. (default: web)
	Dir         string   // Working directory
	Environment []string // Environment variables
	User        string   // User to run as
	NoTTY       bool     // Disable TTY allocation
}

// LogOptions contains options for viewing logs
type LogOptions struct {
	Service    string // web, db, etc. (empty for all)
	Follow     bool   // Tail logs
	Tail       int    // Number of lines to show
	Timestamps bool   // Show timestamps
	Since      string // Show logs since timestamp (e.g., "1h", "30m")
}

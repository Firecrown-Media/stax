package wordpress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// MultisiteType represents the type of multisite installation
type MultisiteType string

const (
	MultisiteTypeSubdomain      MultisiteType = "subdomain"
	MultisiteTypeSubdirectory   MultisiteType = "subdirectory"
	MultisiteTypeNone           MultisiteType = "none"
)

// Subsite represents a single site in a multisite network
type Subsite struct {
	ID         int    `json:"blog_id"`
	SiteID     int    `json:"site_id"`
	Domain     string `json:"domain"`
	Path       string `json:"path"`
	Registered string `json:"registered"`
	LastUpdated string `json:"last_updated"`
	Public     bool   `json:"public"`
	Archived   bool   `json:"archived"`
	Mature     bool   `json:"mature"`
	Spam       bool   `json:"spam"`
	Deleted    bool   `json:"deleted"`
	LangID     int    `json:"lang_id"`
	SiteURL    string `json:"siteurl"`
	Name       string `json:"blogname"`
}

// MultisiteConfig represents multisite configuration
type MultisiteConfig struct {
	Enabled    bool
	Type       MultisiteType
	Domain     string
	Path       string
	Subsites   []Subsite
	NetworkURL string
}

// NetworkConfig represents the network-level configuration
type NetworkConfig struct {
	ID          int
	Domain      string
	Path        string
	SiteCount   int
	AdminEmail  string
	NetworkName string
}

// IsMultisite checks if WordPress is configured as multisite
func IsMultisite(projectPath string) (bool, error) {
	wpConfigPath := filepath.Join(projectPath, "wp-config.php")

	// Check if wp-config.php exists
	data, err := os.ReadFile(wpConfigPath)
	if err != nil {
		return false, fmt.Errorf("failed to read wp-config.php: %w", err)
	}

	content := string(data)

	// Look for multisite constants
	hasMultisite := strings.Contains(content, "WP_ALLOW_MULTISITE") ||
		strings.Contains(content, "MULTISITE") ||
		strings.Contains(content, "SUBDOMAIN_INSTALL")

	return hasMultisite, nil
}

// GetMultisiteType determines if the multisite uses subdomains or subdirectories
func GetMultisiteType(projectPath string) (MultisiteType, error) {
	wpConfigPath := filepath.Join(projectPath, "wp-config.php")

	data, err := os.ReadFile(wpConfigPath)
	if err != nil {
		return MultisiteTypeNone, fmt.Errorf("failed to read wp-config.php: %w", err)
	}

	content := string(data)

	// Check for SUBDOMAIN_INSTALL constant
	if strings.Contains(content, "define( 'SUBDOMAIN_INSTALL', true )") ||
		strings.Contains(content, "define('SUBDOMAIN_INSTALL', true)") ||
		strings.Contains(content, "define(\"SUBDOMAIN_INSTALL\", true)") {
		return MultisiteTypeSubdomain, nil
	}

	if strings.Contains(content, "define( 'SUBDOMAIN_INSTALL', false )") ||
		strings.Contains(content, "define('SUBDOMAIN_INSTALL', false)") ||
		strings.Contains(content, "define(\"SUBDOMAIN_INSTALL\", false)") {
		return MultisiteTypeSubdirectory, nil
	}

	// If MULTISITE is enabled but SUBDOMAIN_INSTALL is not defined, default to subdirectory
	if strings.Contains(content, "MULTISITE") {
		return MultisiteTypeSubdirectory, nil
	}

	return MultisiteTypeNone, nil
}

// GetSubsites retrieves all subsites using WP-CLI
func GetSubsites(cli *CLI) ([]Subsite, error) {
	// Execute wp site list command
	args := []string{"site", "list", "--format=json"}
	output, err := cli.ExecuteWithOutput(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list subsites: %w", err)
	}

	// Parse JSON output
	var subsites []Subsite
	if err := json.Unmarshal([]byte(output), &subsites); err != nil {
		return nil, fmt.Errorf("failed to parse subsites: %w", err)
	}

	// Enrich with additional information
	for i := range subsites {
		// Get site URL
		urlArgs := []string{"option", "get", "siteurl", fmt.Sprintf("--url=%s", subsites[i].Domain)}
		url, err := cli.ExecuteWithOutput(urlArgs...)
		if err == nil {
			subsites[i].SiteURL = strings.TrimSpace(url)
		}

		// Get site name
		nameArgs := []string{"option", "get", "blogname", fmt.Sprintf("--url=%s", subsites[i].Domain)}
		name, err := cli.ExecuteWithOutput(nameArgs...)
		if err == nil {
			subsites[i].Name = strings.TrimSpace(name)
		}
	}

	return subsites, nil
}

// DetectSubsites attempts to detect subsites from the database
func DetectSubsites(cli *CLI) ([]Subsite, error) {
	// Query wp_blogs table
	query := "SELECT blog_id, site_id, domain, path, registered, last_updated, public, archived, mature, spam, deleted, lang_id FROM wp_blogs ORDER BY blog_id"

	args := []string{"db", "query", query, "--skip-column-names"}
	output, err := cli.ExecuteWithOutput(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query subsites: %w", err)
	}

	// Parse output
	var subsites []Subsite
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 12 {
			continue
		}

		blogID, _ := strconv.Atoi(fields[0])
		siteID, _ := strconv.Atoi(fields[1])
		public, _ := strconv.Atoi(fields[6])
		archived, _ := strconv.Atoi(fields[7])
		mature, _ := strconv.Atoi(fields[8])
		spam, _ := strconv.Atoi(fields[9])
		deleted, _ := strconv.Atoi(fields[10])
		langID, _ := strconv.Atoi(fields[11])

		subsite := Subsite{
			ID:          blogID,
			SiteID:      siteID,
			Domain:      fields[2],
			Path:        fields[3],
			Registered:  fields[4] + " " + fields[5],
			Public:      public == 1,
			Archived:    archived == 1,
			Mature:      mature == 1,
			Spam:        spam == 1,
			Deleted:     deleted == 1,
			LangID:      langID,
		}

		subsites = append(subsites, subsite)
	}

	return subsites, nil
}

// GenerateHostsEntries creates hosts file entries for multisite domains
func GenerateHostsEntries(subsites []Subsite, localIP string) []string {
	if localIP == "" {
		localIP = "127.0.0.1"
	}

	entries := []string{}
	seen := make(map[string]bool)

	for _, subsite := range subsites {
		if subsite.Deleted || subsite.Spam || subsite.Archived {
			continue
		}

		// Skip if we've already added this domain
		if seen[subsite.Domain] {
			continue
		}
		seen[subsite.Domain] = true

		entry := fmt.Sprintf("%s %s", localIP, subsite.Domain)
		entries = append(entries, entry)
	}

	return entries
}

// GetNetworkSiteURL returns the main network site URL
func GetNetworkSiteURL(cli *CLI) (string, error) {
	// Get site URL for site 1 (main site)
	args := []string{"option", "get", "siteurl", "--url=1"}
	output, err := cli.ExecuteWithOutput(args...)
	if err != nil {
		// Try without --url flag
		args = []string{"option", "get", "siteurl"}
		output, err = cli.ExecuteWithOutput(args...)
		if err != nil {
			return "", fmt.Errorf("failed to get network site URL: %w", err)
		}
	}

	return strings.TrimSpace(output), nil
}

// GetSubsiteURL returns the URL for a specific subsite
func GetSubsiteURL(cli *CLI, siteID int) (string, error) {
	args := []string{"option", "get", "siteurl", fmt.Sprintf("--url=%d", siteID)}
	output, err := cli.ExecuteWithOutput(args...)
	if err != nil {
		return "", fmt.Errorf("failed to get subsite URL: %w", err)
	}

	return strings.TrimSpace(output), nil
}

// UpdateSubsiteURL changes a subsite's URL
func UpdateSubsiteURL(cli *CLI, siteID int, newURL string) error {
	args := []string{
		"option", "update", "siteurl", newURL,
		fmt.Sprintf("--url=%d", siteID),
	}

	if err := cli.Execute(args...); err != nil {
		return fmt.Errorf("failed to update subsite URL: %w", err)
	}

	// Also update home URL
	args = []string{
		"option", "update", "home", newURL,
		fmt.Sprintf("--url=%d", siteID),
	}

	if err := cli.Execute(args...); err != nil {
		return fmt.Errorf("failed to update home URL: %w", err)
	}

	return nil
}

// GetMultisiteConfig retrieves complete multisite configuration
func GetMultisiteConfig(projectPath string, cli *CLI) (*MultisiteConfig, error) {
	// Check if multisite
	isMulti, err := IsMultisite(projectPath)
	if err != nil {
		return nil, err
	}

	if !isMulti {
		return &MultisiteConfig{
			Enabled: false,
			Type:    MultisiteTypeNone,
		}, nil
	}

	// Get multisite type
	msType, err := GetMultisiteType(projectPath)
	if err != nil {
		return nil, err
	}

	// Get subsites
	subsites, err := GetSubsites(cli)
	if err != nil {
		// Fallback to database detection
		subsites, err = DetectSubsites(cli)
		if err != nil {
			return nil, err
		}
	}

	// Get network URL
	networkURL, err := GetNetworkSiteURL(cli)
	if err != nil {
		networkURL = ""
	}

	// Extract domain and path from network URL
	domain := ""
	path := "/"
	if networkURL != "" {
		parts := strings.SplitN(strings.TrimPrefix(strings.TrimPrefix(networkURL, "http://"), "https://"), "/", 2)
		domain = parts[0]
		if len(parts) > 1 {
			path = "/" + parts[1]
		}
	}

	return &MultisiteConfig{
		Enabled:    true,
		Type:       msType,
		Domain:     domain,
		Path:       path,
		Subsites:   subsites,
		NetworkURL: networkURL,
	}, nil
}

// GetNetworkConfig retrieves network-level configuration
func GetNetworkConfig(cli *CLI) (*NetworkConfig, error) {
	// Get network admin email
	emailArgs := []string{"option", "get", "admin_email", "--network"}
	email, err := cli.ExecuteWithOutput(emailArgs...)
	if err != nil {
		email = ""
	}

	// Get network name
	nameArgs := []string{"option", "get", "site_name", "--network"}
	name, err := cli.ExecuteWithOutput(nameArgs...)
	if err != nil {
		name = ""
	}

	// Count sites
	countArgs := []string{"site", "list", "--format=count"}
	countStr, err := cli.ExecuteWithOutput(countArgs...)
	count := 0
	if err == nil {
		count, _ = strconv.Atoi(strings.TrimSpace(countStr))
	}

	return &NetworkConfig{
		ID:          1,
		SiteCount:   count,
		AdminEmail:  strings.TrimSpace(email),
		NetworkName: strings.TrimSpace(name),
	}, nil
}

// GenerateLocalDomains creates local domain names for subsites
func GenerateLocalDomains(subsites []Subsite, baseDomain string) map[int]string {
	domains := make(map[int]string)

	for _, subsite := range subsites {
		if subsite.Deleted || subsite.Spam || subsite.Archived {
			continue
		}

		// For subdomain multisite, convert production domain to local
		if subsite.ID == 1 {
			// Main site
			domains[subsite.ID] = baseDomain
		} else {
			// Extract subdomain from production domain
			parts := strings.Split(subsite.Domain, ".")
			if len(parts) > 0 {
				subdomain := parts[0]
				domains[subsite.ID] = fmt.Sprintf("%s.%s", subdomain, baseDomain)
			}
		}
	}

	return domains
}

// GetSubsitesByStatus filters subsites by their status
func GetSubsitesByStatus(subsites []Subsite, activeOnly bool) []Subsite {
	filtered := []Subsite{}

	for _, subsite := range subsites {
		if activeOnly && (subsite.Deleted || subsite.Spam || subsite.Archived) {
			continue
		}
		filtered = append(filtered, subsite)
	}

	return filtered
}

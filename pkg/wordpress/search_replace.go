package wordpress

import (
	"fmt"
)

// SearchReplaceOptions represents options for search-replace operations
type SearchReplaceOptions struct {
	Network     bool
	SkipColumns []string
	SkipTables  []string
	DryRun      bool
	URL         string
}

// SearchReplacePair represents a search-replace operation
type SearchReplacePair struct {
	Old string
	New string
	URL string
}

// MultisiteSearchReplaceConfig represents configuration for multisite search-replace
type MultisiteSearchReplaceConfig struct {
	Network NetworkReplace
	Sites   []SiteReplace
}

// NetworkReplace represents network-wide search-replace
type NetworkReplace struct {
	Old string
	New string
}

// SiteReplace represents site-specific search-replace
type SiteReplace struct {
	Old    string
	New    string
	URL    string
	BlogID int
}

// SearchReplace performs a single search-replace operation
func (c *CLI) SearchReplaceWithOptions(old, new string, options SearchReplaceOptions) error {
	args := []string{"search-replace", old, new}

	// Skip GUID column by default
	skipColumns := options.SkipColumns
	if len(skipColumns) == 0 {
		skipColumns = []string{"guid"}
	}
	args = append(args, "--skip-columns="+joinStrings(skipColumns, ","))

	// Network-wide
	if options.Network {
		args = append(args, "--network")
	}

	// Specific site
	if options.URL != "" {
		args = append(args, "--url="+options.URL)
	}

	// Dry run
	if options.DryRun {
		args = append(args, "--dry-run")
	}

	return c.Execute(args...)
}

// MultisiteSearchReplace performs search-replace for all sites in a multisite network
func (c *CLI) MultisiteSearchReplace(config MultisiteSearchReplaceConfig) error {
	// Network-wide replacement
	if config.Network.Old != "" && config.Network.New != "" {
		opts := SearchReplaceOptions{
			Network:     true,
			SkipColumns: []string{"guid"},
		}
		if err := c.SearchReplaceWithOptions(config.Network.Old, config.Network.New, opts); err != nil {
			return fmt.Errorf("network search-replace failed: %w", err)
		}
	}

	// Site-specific replacements
	for _, site := range config.Sites {
		opts := SearchReplaceOptions{
			Network:     false,
			SkipColumns: []string{"guid"},
			URL:         site.URL,
		}
		if err := c.SearchReplaceWithOptions(site.Old, site.New, opts); err != nil {
			return fmt.Errorf("search-replace failed for %s: %w", site.URL, err)
		}
	}

	return nil
}

// PreviewSearchReplace shows what would be replaced without making changes
func (c *CLI) PreviewSearchReplace(old, new string) (string, error) {
	args := []string{"search-replace", old, new, "--dry-run"}
	return c.ExecuteWithOutput(args...)
}

// GetSubsites gets all subsites from the database
func (c *CLI) GetSubsites() ([]Site, error) {
	return c.GetSites()
}

// BuildFirecrownSearchReplaceConfig builds the standard Firecrown multisite search-replace configuration
func BuildFirecrownSearchReplaceConfig(localDomain string) MultisiteSearchReplaceConfig {
	return MultisiteSearchReplaceConfig{
		Network: NetworkReplace{
			Old: "fsmultisite.wpenginepowered.com",
			New: localDomain,
		},
		Sites: []SiteReplace{
			{
				Old: "flyingmag.com",
				New: "flyingmag." + localDomain,
				URL: "flyingmag.com",
			},
			{
				Old: "planeandpilotmag.com",
				New: "planeandpilot." + localDomain,
				URL: "planeandpilotmag.com",
			},
			{
				Old: "finescale.com",
				New: "finescale." + localDomain,
				URL: "finescale.com",
			},
			{
				Old: "avweb.com",
				New: "avweb." + localDomain,
				URL: "avweb.com",
			},
		},
	}
}

// Helper function to join strings
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}
	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += separator + strings[i]
	}
	return result
}

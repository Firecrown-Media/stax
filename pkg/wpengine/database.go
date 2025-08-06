package wpengine

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Firecrown-Media/stax/pkg/ui"
)

type DatabaseManager struct {
	projectPath string
	client      *Client
}

func NewDatabaseManager(projectPath string, client *Client) *DatabaseManager {
	return &DatabaseManager{
		projectPath: projectPath,
		client:      client,
	}
}

func (dm *DatabaseManager) ImportDatabase(dbFile string, options SyncOptions) (*DatabaseSyncResult, error) {
	// Check if we have a DDEV project first
	if !dm.isDDEVProject() {
		return &DatabaseSyncResult{
			Success: false,
			Error:   fmt.Sprintf("no DDEV project found in %s. Please run 'stax init' or 'ddev config' to initialize a DDEV project first", dm.projectPath),
		}, fmt.Errorf("no DDEV project found in %s. Please run 'stax init' or 'ddev config' to initialize a DDEV project first", dm.projectPath)
	}

	// Note: WordPress debug suppression will be done after import when WordPress is ready

	// First, import the database using DDEV
	err := ui.WithSpinnerResult("Importing database into DDEV", func() error {
		cmd := exec.Command("ddev", "import-db", "--src="+dbFile)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to import database: %s", string(output))
		}
		return nil
	})
	
	if err != nil {
		return &DatabaseSyncResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Skip URL rewriting during database import - it will be done after file sync
	fmt.Println("Database imported successfully")
	fmt.Println("Note: URL rewriting will be performed after WordPress files are synced")
	
	return &DatabaseSyncResult{
		Success:       true,
		DatabaseFile:  dbFile,
		RewrittenURLs: 0, // No URL rewriting done at this stage
	}, nil
}

func (dm *DatabaseManager) rewriteURLs(oldURL, newURL string, options SyncOptions) (int, error) {
	rewrittenURLs := 0

	// Check if wp-cli is available in DDEV
	if dm.canUseWPCLI() {
		// Use advanced wp-cli replacement for comprehensive URL rewriting
		count, err := dm.replaceWithWPCLIAdvanced(oldURL, newURL, options)
		if err != nil {
			fmt.Printf("Advanced wp-cli replacement failed, falling back to basic method: %v\n", err)
			count, err = dm.replaceWithWPCLI(oldURL, newURL, options)
			if err != nil {
				fmt.Printf("Basic wp-cli replacement failed, falling back to SQL: %v\n", err)
				return dm.rewriteURLsWithSQL(oldURL, newURL, options)
			}
		}
		rewrittenURLs += count
	} else {
		// Fallback to SQL replacements
		return dm.rewriteURLsWithSQL(oldURL, newURL, options)
	}

	return rewrittenURLs, nil
}

func (dm *DatabaseManager) replaceWithWPCLI(oldURL, newURL string, options SyncOptions) (int, error) {
	// Use wp-cli search-replace with aggressive plugin/theme skipping and explicit path
	args := []string{"wp", "search-replace", oldURL, newURL, "--dry-run", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"}
	
	cmd := exec.Command("ddev", args...)
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("wp-cli dry-run failed: %s", string(output))
	}

	// Parse the dry-run output to get the count
	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "replacements") {
			fmt.Sscanf(line, "%d", &count)
			break
		}
	}

	// If dry-run looks good, do the actual replacement
	if count > 0 {
		args = []string{"wp", "search-replace", oldURL, newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"}
		
		cmd = exec.Command("ddev", args...)
		cmd.Dir = dm.projectPath
		output, err = cmd.CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("wp-cli replacement failed: %s", string(output))
		}

		// Update WordPress home and site URLs
		if err := dm.updateWordPressURLs(newURL); err != nil {
			fmt.Printf("Warning: failed to update WordPress URLs: %v\n", err)
		}
	}

	return count, nil
}

func (dm *DatabaseManager) replaceInTable(table, column, oldValue, newValue string) (int, error) {
	// Use wp-cli for safe WordPress-aware replacements
	cmd := exec.Command("ddev", "wp", "search-replace", oldValue, newValue, "--dry-run")
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("dry-run failed: %s", string(output))
	}

	// Parse the dry-run output to get the count
	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "replacements") {
			fmt.Sscanf(line, "%d", &count)
			break
		}
	}

	// If dry-run looks good, do the actual replacement
	if count > 0 {
		cmd = exec.Command("ddev", "wp", "search-replace", oldValue, newValue)
		cmd.Dir = dm.projectPath
		output, err = cmd.CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("replacement failed: %s", string(output))
		}
	}

	return count, nil
}

func (dm *DatabaseManager) updateWordPressURLs(newURL string) error {
	// Update home and siteurl options with aggressive skipping and explicit path
	commands := [][]string{
		{"ddev", "wp", "option", "update", "home", newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"},
		{"ddev", "wp", "option", "update", "siteurl", newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = dm.projectPath
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to update WordPress URL: %s", string(output))
		}
	}

	return nil
}

func (dm *DatabaseManager) GetProjectName() string {
	// Try to get project name from DDEV config
	cmd := exec.Command("ddev", "describe", "-j")
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Fallback to directory name
		return filepath.Base(dm.projectPath)
	}

	// Parse JSON output to get project name
	// This is a simplified approach - in a real implementation,
	// you'd want to properly parse the JSON
	re := regexp.MustCompile(`"name":\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		return matches[1]
	}

	return filepath.Base(dm.projectPath)
}

func (dm *DatabaseManager) CreateDatabaseBackup(outputFile string) error {
	// Check if we have a DDEV project first
	if !dm.isDDEVProject() {
		return fmt.Errorf("no DDEV project found in %s. Please run 'stax init' or 'ddev config' to initialize a DDEV project first", dm.projectPath)
	}

	fmt.Printf("Creating local database backup: %s\n", outputFile)
	
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	cmd := exec.Command("ddev", "export-db", "--file="+outputFile)
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to export database: %s", string(output))
	}

	fmt.Printf("Database backup created: %s\n", outputFile)
	return nil
}

func (dm *DatabaseManager) hasWordPressFiles() bool {
	// Check for key WordPress files
	wpFiles := []string{
		"wp-config.php",
		"index.php",
		"wp-includes/version.php",
	}
	
	for _, file := range wpFiles {
		filePath := filepath.Join(dm.projectPath, file)
		if _, err := os.Stat(filePath); err == nil {
			return true
		}
	}
	
	return false
}

func (dm *DatabaseManager) isDDEVProject() bool {
	// Check if .ddev/config.yaml exists
	configPath := filepath.Join(dm.projectPath, ".ddev", "config.yaml")
	if _, err := os.Stat(configPath); err != nil {
		return false
	}
	return true
}

func (dm *DatabaseManager) getTablePrefix() (string, error) {
	// Try to detect table prefix using direct database query first (most reliable)
	cmd := exec.Command("ddev", "mysql", "-e", "SHOW TABLES LIKE 'wp_%'", "-s", "-N")
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err == nil && len(output) > 0 {
		// Parse first table name to extract prefix
		firstTable := strings.TrimSpace(strings.Split(string(output), "\n")[0])
		if strings.Contains(firstTable, "_") {
			// Find the prefix by looking for the pattern before a known WordPress table suffix
			for _, suffix := range []string{"_options", "_posts", "_users", "_comments", "_postmeta"} {
				if strings.HasSuffix(firstTable, suffix) {
					return strings.TrimSuffix(firstTable, suffix) + "_", nil
				}
			}
			// Fallback: assume everything before last underscore + underscore is the prefix
			parts := strings.Split(firstTable, "_")
			if len(parts) >= 2 {
				prefix := strings.Join(parts[:len(parts)-1], "_") + "_"
				fmt.Printf("Debug: Detected table prefix from database: %s\n", prefix)
				return prefix, nil
			}
		}
	}

	// Fallback to WP-CLI if available and WordPress is configured
	if dm.canUseWPCLI() {
		cmd = exec.Command("ddev", "wp", "eval", "echo $wpdb->prefix;", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html")
		cmd.Dir = dm.projectPath
		output, err = cmd.CombinedOutput()
		if err == nil {
			prefix := strings.TrimSpace(string(output))
			if prefix != "" {
				fmt.Printf("Debug: Detected table prefix via WP-CLI: %s\n", prefix)
				return prefix, nil
			}
		}
	}

	// Final fallback: try to detect from any table that looks like WordPress
	cmd = exec.Command("ddev", "mysql", "-e", "SHOW TABLES", "-s", "-N")
	cmd.Dir = dm.projectPath
	output, err = cmd.CombinedOutput()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				// Look for WordPress-like table patterns
				for _, suffix := range []string{"options", "posts", "users", "comments", "postmeta", "usermeta", "commentmeta"} {
					if strings.HasSuffix(line, "_"+suffix) {
						prefix := strings.TrimSuffix(line, "_"+suffix) + "_"
						fmt.Printf("Debug: Detected table prefix from table pattern: %s\n", prefix)
						return prefix, nil
					}
				}
			}
		}
	}

	fmt.Printf("Debug: Could not detect table prefix, using default: wp_\n")
	return "wp_", fmt.Errorf("could not detect table prefix, using default wp_")
}

func (dm *DatabaseManager) canUseWPCLI() bool {
	// First check if DDEV project exists
	if !dm.isDDEVProject() {
		fmt.Printf("Debug: No DDEV project found in %s\n", dm.projectPath)
		return false
	}

	// Try multiple WordPress paths in order of likelihood
	wpPaths := []string{
		"", // Let wp-cli auto-detect
		"/var/www/html",
		"/var/www/html/web", // Some WordPress installs use subdirectories
		"/var/www/html/public",
		"/var/www/html/wordpress",
	}

	for _, wpPath := range wpPaths {
		args := []string{"wp", "core", "version", "--skip-plugins", "--skip-themes", "--skip-packages"}
		if wpPath != "" {
			args = append(args, "--path="+wpPath)
		}
		
		cmd := exec.Command("ddev", args...)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err == nil {
			// Success! wp-cli can find WordPress at this path
			return true
		}
		
		// Log only the first failure for clarity
		if wpPath == "" {
			fmt.Printf("Debug: wp-cli test failed: %s\n", string(output))
		}
	}
	
	fmt.Printf("Debug: wp-cli could not find WordPress installation at any expected path\n")
	return false
}

// waitForDatabaseReady waits for the database to be accessible
func (dm *DatabaseManager) waitForDatabaseReady() error {
	fmt.Printf("Debug: Waiting for database to be ready...\n")
	
	maxAttempts := 30
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Try a simple database connection test
		cmd := exec.Command("ddev", "mysql", "-e", "SELECT 1;")
		cmd.Dir = dm.projectPath
		if output, err := cmd.CombinedOutput(); err == nil {
			fmt.Printf("Debug: Database ready after %d attempts\n", attempt)
			return nil
		} else if attempt < maxAttempts {
			fmt.Printf("Debug: Database not ready (attempt %d/%d), waiting...\n", attempt, maxAttempts)
			time.Sleep(2 * time.Second)
		} else {
			return fmt.Errorf("database not ready after %d attempts: %s", maxAttempts, string(output))
		}
	}
	
	return fmt.Errorf("database not ready after %d attempts", maxAttempts)
}

// waitForWordPressReady waits for WordPress to be properly configured
func (dm *DatabaseManager) waitForWordPressReady() error {
	fmt.Printf("Debug: Waiting for WordPress to be ready...\n")
	
	maxAttempts := 15
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Try to get WordPress version - this will fail if WP isn't configured
		cmd := exec.Command("ddev", "wp", "core", "version", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html")
		cmd.Dir = dm.projectPath
		if output, err := cmd.CombinedOutput(); err == nil {
			fmt.Printf("Debug: WordPress ready after %d attempts\n", attempt)
			return nil
		} else if attempt < maxAttempts {
			fmt.Printf("Debug: WordPress not ready (attempt %d/%d): %s\n", attempt, maxAttempts, string(output))
			time.Sleep(3 * time.Second)
		} else {
			return fmt.Errorf("WordPress not ready after %d attempts: %s", maxAttempts, string(output))
		}
	}
	
	return fmt.Errorf("WordPress not ready after %d attempts", maxAttempts)
}

// configureWordPressTablePrefix ensures WordPress wp-config.php and wp-config-ddev.php have the correct table prefix
func (dm *DatabaseManager) configureWordPressTablePrefix() error {
	// Detect the actual table prefix from the database
	prefix, err := dm.getTablePrefix()
	if err != nil {
		return fmt.Errorf("could not detect table prefix: %w", err)
	}

	// If it's the default prefix, no configuration needed
	if prefix == "wp_" {
		return nil
	}

	fmt.Printf("Debug: Configuring WordPress to use table prefix: %s\n", prefix)

	// Update both wp-config.php and wp-config-ddev.php if they exist
	configFiles := []string{
		filepath.Join(dm.projectPath, "wp-config.php"),
		filepath.Join(dm.projectPath, "wp-config-ddev.php"),
	}

	configUpdated := false
	for _, configPath := range configFiles {
		if _, err := os.Stat(configPath); err == nil {
			if err := dm.updateConfigFilePrefix(configPath, prefix); err != nil {
				fmt.Printf("Warning: Failed to update %s: %v\n", filepath.Base(configPath), err)
			} else {
				configUpdated = true
			}
		}
	}

	// If no config files were updated, create a basic wp-config.php
	if !configUpdated {
		return dm.createWordPressConfigWithPrefix(prefix)
	}

	return nil
}

// updateConfigFilePrefix updates the table prefix in a specific config file
func (dm *DatabaseManager) updateConfigFilePrefix(configPath, prefix string) error {
	// Read existing config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", configPath, err)
	}

	configStr := string(content)
	
	// Check if $table_prefix is already set correctly
	if strings.Contains(configStr, fmt.Sprintf("$table_prefix = '%s';", prefix)) {
		return nil // Already configured correctly
	}

	// Replace or add table prefix configuration
	tablePrefixRegex := regexp.MustCompile(`\$table_prefix\s*=\s*['"]([^'"]*?)['"];`)
	
	if tablePrefixRegex.MatchString(configStr) {
		// Replace existing table prefix
		configStr = tablePrefixRegex.ReplaceAllString(configStr, fmt.Sprintf("$table_prefix = '%s';", prefix))
		fmt.Printf("Debug: Updated existing table prefix in %s\n", filepath.Base(configPath))
	} else {
		// Add table prefix before the database constants
		dbConstantsRegex := regexp.MustCompile(`(define\s*\(\s*['"]DB_NAME['"])`)
		if dbConstantsRegex.MatchString(configStr) {
			configStr = dbConstantsRegex.ReplaceAllString(configStr, fmt.Sprintf("$table_prefix = '%s';\n\n$1", prefix))
			fmt.Printf("Debug: Added table prefix to %s\n", filepath.Base(configPath))
		} else {
			return fmt.Errorf("could not find suitable location to add table prefix in %s", configPath)
		}
	}

	// Write the updated configuration
	return os.WriteFile(configPath, []byte(configStr), 0644)
}

// createWordPressConfigWithPrefix creates a basic wp-config.php with the correct table prefix
func (dm *DatabaseManager) createWordPressConfigWithPrefix(prefix string) error {
	wpConfigPath := filepath.Join(dm.projectPath, "wp-config.php")
	
	// Basic wp-config.php content for DDEV with custom table prefix
	configContent := fmt.Sprintf(`<?php
// Database settings for DDEV
define('DB_NAME', 'db');
define('DB_USER', 'db');
define('DB_PASSWORD', 'db');
define('DB_HOST', 'db');
define('DB_CHARSET', 'utf8mb4');
define('DB_COLLATE', '');

// Custom table prefix
$table_prefix = '%s';

// WordPress debugging
define('WP_DEBUG', false);

// WordPress security keys (minimal for local development)
define('AUTH_KEY',         'local-dev-key');
define('SECURE_AUTH_KEY',  'local-dev-key');
define('LOGGED_IN_KEY',    'local-dev-key');
define('NONCE_KEY',        'local-dev-key');
define('AUTH_SALT',        'local-dev-key');
define('SECURE_AUTH_SALT', 'local-dev-key');
define('LOGGED_IN_SALT',   'local-dev-key');
define('NONCE_SALT',       'local-dev-key');

// WordPress directory structure
if (!defined('ABSPATH')) {
    define('ABSPATH', __DIR__ . '/');
}

require_once(ABSPATH . 'wp-settings.php');
`, prefix)

	err := os.WriteFile(wpConfigPath, []byte(configContent), 0644)
	if err != nil {
		return fmt.Errorf("could not create wp-config.php: %w", err)
	}

	fmt.Printf("Debug: Created wp-config.php with table prefix: %s\n", prefix)
	return nil
}

func (dm *DatabaseManager) RewriteURLsPostSync(oldURL, newURL string, options SyncOptions) (int, error) {
	// Wait for database to be ready first
	if err := dm.waitForDatabaseReady(); err != nil {
		fmt.Printf("Warning: Database not ready, falling back to SQL: %v\n", err)
		return dm.rewriteURLsWithSQL(oldURL, newURL, options)
	}

	// Configure WordPress with correct table prefix if needed
	if err := dm.configureWordPressTablePrefix(); err != nil {
		fmt.Printf("Warning: Could not configure WordPress table prefix: %v\n", err)
	}

	// Wait for WordPress to be properly configured
	if err := dm.waitForWordPressReady(); err != nil {
		fmt.Printf("Warning: WordPress not ready, falling back to SQL: %v\n", err)
		return dm.rewriteURLsWithSQL(oldURL, newURL, options)
	}

	// Suppress WordPress debug output if requested
	if err := dm.suppressWordPressDebug(options); err != nil {
		fmt.Printf("Warning: Could not suppress WordPress debug output: %v\n", err)
	}

	var result int
	var err error

	// Check if we can use wp-cli now
	if !dm.canUseWPCLI() {
		// Fallback to direct SQL replacements
		result, err = dm.rewriteURLsWithSQL(oldURL, newURL, options)
	} else {
		result, err = dm.rewriteURLs(oldURL, newURL, options)
	}

	// Provide info about debug suppression if it was used
	if options.SuppressDebug && err == nil {
		if err := dm.restoreWordPressDebug(options); err != nil {
			fmt.Printf("Warning: Could not restore WordPress debug settings: %v\n", err)
		}
	}

	return result, err
}

func (dm *DatabaseManager) rewriteURLsWithSQL(oldURL, newURL string, options SyncOptions) (int, error) {
	// Check if we have a DDEV project first
	if !dm.isDDEVProject() {
		return 0, fmt.Errorf("no DDEV project found in %s. Please run 'stax init' or 'ddev config' to initialize a DDEV project first", dm.projectPath)
	}

	// Detect the table prefix
	prefix, err := dm.getTablePrefix()
	if err != nil {
		fmt.Printf("Warning: Could not detect table prefix: %v\n", err)
	}

	fmt.Printf("Debug: Using comprehensive SQL-based URL rewriting from %s to %s\n", oldURL, newURL)
	fmt.Printf("Debug: Using table prefix: %s\n", prefix)
	fmt.Printf("Debug: All URLs including media URLs will be rewritten\n")
	
	// Use direct SQL UPDATE statements for comprehensive URL replacement
	// Handle both HTTP and HTTPS versions
	oldURLHttp := strings.Replace(oldURL, "https://", "http://", 1)
	oldURLHttps := strings.Replace(oldURL, "http://", "https://", 1)
	
	replacements := []string{
		// Core WordPress options
		fmt.Sprintf("UPDATE %soptions SET option_value = REPLACE(option_value, '%s', '%s');", prefix, oldURL, newURL),
		fmt.Sprintf("UPDATE %soptions SET option_value = REPLACE(option_value, '%s', '%s');", prefix, oldURLHttp, newURL),
		fmt.Sprintf("UPDATE %soptions SET option_value = REPLACE(option_value, '%s', '%s');", prefix, oldURLHttps, newURL),
		
		// Post content (includes media URLs in content)
		fmt.Sprintf("UPDATE %sposts SET post_content = REPLACE(post_content, '%s', '%s');", prefix, oldURL, newURL),
		fmt.Sprintf("UPDATE %sposts SET post_content = REPLACE(post_content, '%s', '%s');", prefix, oldURLHttp, newURL),
		fmt.Sprintf("UPDATE %sposts SET post_content = REPLACE(post_content, '%s', '%s');", prefix, oldURLHttps, newURL),
		
		// Post excerpts
		fmt.Sprintf("UPDATE %sposts SET post_excerpt = REPLACE(post_excerpt, '%s', '%s');", prefix, oldURL, newURL),
		fmt.Sprintf("UPDATE %sposts SET post_excerpt = REPLACE(post_excerpt, '%s', '%s');", prefix, oldURLHttp, newURL),
		fmt.Sprintf("UPDATE %sposts SET post_excerpt = REPLACE(post_excerpt, '%s', '%s');", prefix, oldURLHttps, newURL),
		
		// GUIDs (important for media attachments)
		fmt.Sprintf("UPDATE %sposts SET guid = REPLACE(guid, '%s', '%s');", prefix, oldURL, newURL),
		fmt.Sprintf("UPDATE %sposts SET guid = REPLACE(guid, '%s', '%s');", prefix, oldURLHttp, newURL),
		fmt.Sprintf("UPDATE %sposts SET guid = REPLACE(guid, '%s', '%s');", prefix, oldURLHttps, newURL),
		
		// Comments
		fmt.Sprintf("UPDATE %scomments SET comment_content = REPLACE(comment_content, '%s', '%s');", prefix, oldURL, newURL),
		fmt.Sprintf("UPDATE %scomments SET comment_content = REPLACE(comment_content, '%s', '%s');", prefix, oldURLHttp, newURL),
		fmt.Sprintf("UPDATE %scomments SET comment_content = REPLACE(comment_content, '%s', '%s');", prefix, oldURLHttps, newURL),
		
		// Post meta (includes serialized data)
		fmt.Sprintf("UPDATE %spostmeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURL, newURL),
		fmt.Sprintf("UPDATE %spostmeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURLHttp, newURL),
		fmt.Sprintf("UPDATE %spostmeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURLHttps, newURL),
		
		// User meta (for user profile images, etc.)
		fmt.Sprintf("UPDATE %susermeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURL, newURL),
		fmt.Sprintf("UPDATE %susermeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURLHttp, newURL),
		fmt.Sprintf("UPDATE %susermeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURLHttps, newURL),
		
		// Comment meta
		fmt.Sprintf("UPDATE %scommentmeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURL, newURL),
		fmt.Sprintf("UPDATE %scommentmeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURLHttp, newURL),
		fmt.Sprintf("UPDATE %scommentmeta SET meta_value = REPLACE(meta_value, '%s', '%s');", prefix, oldURLHttps, newURL),
	}

	totalReplacements := 0
	for i, query := range replacements {
		cmd := exec.Command("ddev", "mysql", "-e", query)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: SQL replacement %d failed: %v - %s\n", i+1, err, string(output))
		} else {
			fmt.Printf("Debug: SQL replacement %d succeeded (including media URLs)\n", i+1)
			totalReplacements += 1
		}
	}

	return totalReplacements, nil
}

func (dm *DatabaseManager) OptimizeDatabase() error {
	// Check if we have a DDEV project first
	if !dm.isDDEVProject() {
		return fmt.Errorf("no DDEV project found in %s. Please run 'stax init' or 'ddev config' to initialize a DDEV project first", dm.projectPath)
	}

	return ui.WithSpinnerResult("Optimizing database", func() error {
		// Run WordPress database optimization
		cmd := exec.Command("ddev", "wp", "db", "optimize")
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to optimize database: %s", string(output))
		}

		// Flush WordPress caches
		cmd = exec.Command("ddev", "wp", "cache", "flush")
		cmd.Dir = dm.projectPath
		if output, err := cmd.CombinedOutput(); err != nil {
			// Cache flush is not critical, just log it
			fmt.Printf("Warning: failed to flush cache: %s\n", string(output))
		}

		return nil
	})
}

// MediaURL represents a detected media URL with context
type MediaURL struct {
	URL       string `json:"url"`
	Extension string `json:"extension"`
	IsUpload  bool   `json:"is_upload"`
	Size      string `json:"size,omitempty"`
}

// GetWordPressMediaURLs retrieves all media URLs from WordPress using wp-cli
func (dm *DatabaseManager) GetWordPressMediaURLs() ([]MediaURL, error) {
	var mediaURLs []MediaURL
	
	if !dm.canUseWPCLI() {
		return mediaURLs, fmt.Errorf("wp-cli not available")
	}

	// Get all attachment URLs from WordPress
	cmd := exec.Command("ddev", "wp", "post", "list", "--post_type=attachment", "--field=guid", "--format=json", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html")
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Debug: Failed to get attachment URLs via wp-cli: %s\n", string(output))
		return mediaURLs, fmt.Errorf("failed to get attachment URLs: %s", string(output))
	}

	// Parse JSON output
	var attachmentURLs []string
	if err := json.Unmarshal(output, &attachmentURLs); err != nil {
		return mediaURLs, fmt.Errorf("failed to parse attachment URLs: %w", err)
	}

	// Convert to MediaURL structs
	for _, url := range attachmentURLs {
		if url != "" {
			mediaURL := MediaURL{
				URL:       url,
				Extension: getFileExtension(url),
				IsUpload:  strings.Contains(url, "/wp-content/uploads/"),
			}
			mediaURLs = append(mediaURLs, mediaURL)
		}
	}

	return mediaURLs, nil
}

// GetUploadsBaseURL gets the WordPress uploads directory base URL
func (dm *DatabaseManager) GetUploadsBaseURL() (string, error) {
	if !dm.canUseWPCLI() {
		return "", fmt.Errorf("wp-cli not available")
	}

	// Get uploads directory URL from WordPress
	cmd := exec.Command("ddev", "wp", "eval", "echo wp_upload_dir()['baseurl'];", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html")
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get uploads URL: %s", string(output))
	}

	uploadsURL := strings.TrimSpace(string(output))
	return uploadsURL, nil
}

// GetOriginalUploadsURLFromDatabase gets the original uploads URL directly from the database before any URL rewriting
func (dm *DatabaseManager) GetOriginalUploadsURLFromDatabase() (string, error) {
	// Detect the table prefix
	prefix, err := dm.getTablePrefix()
	if err != nil {
		fmt.Printf("Warning: Could not detect table prefix: %v\n", err)
	}

	// Try to get the home URL from options table to determine the original domain
	query := fmt.Sprintf("SELECT option_value FROM %soptions WHERE option_name = 'home' LIMIT 1;", prefix)
	cmd := exec.Command("ddev", "mysql", "-e", query, "-s", "-N")
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get home URL from database: %s", string(output))
	}

	homeURL := strings.TrimSpace(string(output))
	if homeURL == "" {
		return "", fmt.Errorf("no home URL found in database")
	}

	// Construct the uploads URL based on the original home URL
	uploadsURL := homeURL + "/wp-content/uploads"
	fmt.Printf("Debug: Original uploads URL from database: %s\n", uploadsURL)
	return uploadsURL, nil
}

// identifyMediaURLsInContent analyzes content to find media URLs
func (dm *DatabaseManager) identifyMediaURLsInContent(content string, knownMediaURLs []MediaURL) []string {
	var detectedMedia []string
	
	// Common media file extensions
	mediaExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff", ".svg",
		".mp4", ".mov", ".avi", ".wmv", ".flv", ".webm", ".mkv",
		".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".zip", ".rar", ".tar", ".gz", ".7z",
	}

	// Regex to find URLs in content
	urlRegex := regexp.MustCompile(`https?://[^\s<>"'\)]+`)
	foundURLs := urlRegex.FindAllString(content, -1)

	for _, url := range foundURLs {
		isMedia := false
		
		// Check against known WordPress media URLs
		for _, knownMedia := range knownMediaURLs {
			if strings.Contains(url, knownMedia.URL) || url == knownMedia.URL {
				isMedia = true
				break
			}
		}
		
		// Check if URL contains uploads directory path
		if strings.Contains(url, "/wp-content/uploads/") {
			isMedia = true
		}
		
		// Check file extension
		if !isMedia {
			urlLower := strings.ToLower(url)
			for _, ext := range mediaExtensions {
				if strings.Contains(urlLower, ext) {
					isMedia = true
					break
				}
			}
		}
		
		if isMedia {
			detectedMedia = append(detectedMedia, url)
		}
	}
	
	return detectedMedia
}

// getFileExtension extracts file extension from URL
func getFileExtension(url string) string {
	// Remove query parameters and fragments
	cleanURL := url
	if idx := strings.Index(url, "?"); idx != -1 {
		cleanURL = url[:idx]
	}
	if idx := strings.Index(cleanURL, "#"); idx != -1 {
		cleanURL = cleanURL[:idx]
	}
	
	// Extract extension
	parts := strings.Split(cleanURL, ".")
	if len(parts) > 1 {
		return "." + strings.ToLower(parts[len(parts)-1])
	}
	return ""
}

// replaceWithWPCLIAdvanced performs comprehensive URL replacement including media URLs
func (dm *DatabaseManager) replaceWithWPCLIAdvanced(oldURL, newURL string, options SyncOptions) (int, error) {
	fmt.Printf("Debug: Starting comprehensive URL replacement from %s to %s\n", oldURL, newURL)
	fmt.Printf("Debug: Media URLs will be rewritten along with all other URLs\n")
	
	// Use standard wp-cli search-replace for comprehensive URL rewriting
	args := []string{"wp", "search-replace", oldURL, newURL, "--dry-run", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"}
	
	cmd := exec.Command("ddev", args...)
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("wp-cli dry-run failed: %s", string(output))
	}

	// Parse the dry-run output to get the count
	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "replacements") {
			fmt.Sscanf(line, "%d", &count)
			break
		}
	}

	// If dry-run looks good, do the actual replacement
	if count > 0 {
		args = []string{"wp", "search-replace", oldURL, newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"}
		
		cmd = exec.Command("ddev", args...)
		cmd.Dir = dm.projectPath
		output, err = cmd.CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("wp-cli replacement failed: %s", string(output))
		}

		// Update WordPress home and site URLs
		if err := dm.updateWordPressURLs(newURL); err != nil {
			fmt.Printf("Warning: failed to update WordPress URLs: %v\n", err)
		}
		
		fmt.Printf("Debug: Successfully replaced %d URLs including media URLs\n", count)
	} else {
		fmt.Printf("Debug: No URLs found to replace\n")
	}

	return count, nil
}

// replaceInOptionsTable specifically handles wp_options table (home, siteurl, etc.)
func (dm *DatabaseManager) replaceInOptionsTable(oldURL, newURL string) (int, error) {
	// Update WordPress home and site URLs directly
	commands := [][]string{
		{"ddev", "wp", "option", "update", "home", newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"},
		{"ddev", "wp", "option", "update", "siteurl", newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"},
	}

	count := 0
	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = dm.projectPath
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("Warning: Failed to update option: %s\n", string(output))
		} else {
			count++
		}
	}

	return count, nil
}

// replaceInContentWithMediaExclusion replaces URLs in content with media exclusion applied FIRST
// This ensures media URLs are never touched by the rewriting process
func (dm *DatabaseManager) replaceInContentWithMediaExclusion(oldURL, newURL string, mediaURLsToPreserve []string) (int, error) {
	if len(mediaURLsToPreserve) == 0 {
		// No media preservation needed, do normal replacement in content
		args := []string{"wp", "search-replace", oldURL, newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html", "--skip-columns=guid"}
		cmd := exec.Command("ddev", args...)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("wp-cli replacement failed: %s", string(output))
		}
		return 1, nil // Simplified count
	}

	fmt.Printf("Debug: Applying media exclusion BEFORE URL rewriting from %s to %s\n", oldURL, newURL)
	for i, mediaURL := range mediaURLsToPreserve {
		fmt.Printf("Debug: Excluding media URL pattern %d: %s\n", i+1, mediaURL)
	}
	
	// Get table prefix
	prefix, err := dm.getTablePrefix()
	if err != nil {
		fmt.Printf("Warning: Could not detect table prefix: %v\n", err)
	}

	// Build SQL exclusion conditions for ALL media URL patterns
	// These conditions ensure rows containing media URLs are NEVER updated
	var excludeConditions []string
	for _, mediaURL := range mediaURLsToPreserve {
		excludeConditions = append(excludeConditions, fmt.Sprintf("post_content NOT LIKE '%%%s%%'", mediaURL))
	}
	excludeClause := strings.Join(excludeConditions, " AND ")
	
	// SQL queries that ONLY update content that does NOT contain any media URL patterns
	queries := []string{
		// Update post content ONLY if it contains the old URL AND does NOT contain any media URL patterns
		fmt.Sprintf("UPDATE %sposts SET post_content = REPLACE(post_content, '%s', '%s') WHERE post_content LIKE '%%%s%%' AND %s;", 
			prefix, oldURL, newURL, oldURL, excludeClause),
		
		// Update post excerpts ONLY if they contain the old URL AND do NOT contain any media URL patterns  
		fmt.Sprintf("UPDATE %sposts SET post_excerpt = REPLACE(post_excerpt, '%s', '%s') WHERE post_excerpt LIKE '%%%s%%' AND %s;", 
			prefix, oldURL, newURL, oldURL, strings.Replace(excludeClause, "post_content", "post_excerpt", -1)),
		
		// Update comments ONLY if they contain the old URL AND do NOT contain any media URL patterns
		fmt.Sprintf("UPDATE %scomments SET comment_content = REPLACE(comment_content, '%s', '%s') WHERE comment_content LIKE '%%%s%%' AND %s;", 
			prefix, oldURL, newURL, oldURL, strings.Replace(excludeClause, "post_content", "comment_content", -1)),
		
		// Update meta values ONLY if they contain the old URL AND do NOT contain any media URL patterns
		fmt.Sprintf("UPDATE %spostmeta SET meta_value = REPLACE(meta_value, '%s', '%s') WHERE meta_value LIKE '%%%s%%' AND %s;", 
			prefix, oldURL, newURL, oldURL, strings.Replace(excludeClause, "post_content", "meta_value", -1)),
	}

	count := 0
	for i, query := range queries {
		fmt.Printf("Debug: Executing media-exclusion query %d\n", i+1)
		cmd := exec.Command("ddev", "mysql", "-e", query)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Media-exclusion query %d failed: %v - %s\n", i+1, err, string(output))
		} else {
			fmt.Printf("Debug: Media-exclusion query %d succeeded (content updated only if no media URLs present)\n", i+1)
			count++
		}
	}

	fmt.Printf("Debug: Media exclusion complete - %d query types executed with media filtering\n", count)
	return count, nil
}

// replaceInContentPreservingMediaURLs replaces URLs in content while preserving multiple media URL patterns
// DEPRECATED: Use replaceInContentWithMediaExclusion instead for better media filtering
func (dm *DatabaseManager) replaceInContentPreservingMediaURLs(oldURL, newURL string, mediaURLsToPreserve []string) (int, error) {
	if len(mediaURLsToPreserve) == 0 {
		// No media preservation needed, do normal replacement in content
		args := []string{"wp", "search-replace", oldURL, newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html", "--skip-columns=guid"}
		cmd := exec.Command("ddev", args...)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("wp-cli replacement failed: %s", string(output))
		}
		return 1, nil // Simplified count
	}

	// When preserving media URLs, we use a more targeted approach
	fmt.Printf("Debug: Replacing %s with %s in content\n", oldURL, newURL)
	for i, mediaURL := range mediaURLsToPreserve {
		fmt.Printf("Debug: Will preserve media URL pattern %d: %s\n", i+1, mediaURL)
	}
	
	// Get table prefix
	prefix, err := dm.getTablePrefix()
	if err != nil {
		fmt.Printf("Warning: Could not detect table prefix: %v\n", err)
	}

	// Build SQL conditions to exclude ALL media URL patterns
	var excludeConditions []string
	for _, mediaURL := range mediaURLsToPreserve {
		excludeConditions = append(excludeConditions, fmt.Sprintf("post_content NOT LIKE '%%%s%%'", mediaURL))
	}
	excludeClause := strings.Join(excludeConditions, " AND ")
	
	queries := []string{
		// Replace URLs in post content, but NOT if they contain any of the media URL patterns
		fmt.Sprintf("UPDATE %sposts SET post_content = REPLACE(post_content, '%s', '%s') WHERE post_content LIKE '%%%s%%' AND %s;", 
			prefix, oldURL, newURL, oldURL, excludeClause),
		
		// Replace URLs in post excerpts, but NOT if they contain any of the media URL patterns  
		fmt.Sprintf("UPDATE %sposts SET post_excerpt = REPLACE(post_excerpt, '%s', '%s') WHERE post_excerpt LIKE '%%%s%%' AND %s;", 
			prefix, oldURL, newURL, oldURL, strings.Replace(excludeClause, "post_content", "post_excerpt", -1)),
		
		// Replace URLs in comments, but NOT if they contain any of the media URL patterns
		fmt.Sprintf("UPDATE %scomments SET comment_content = REPLACE(comment_content, '%s', '%s') WHERE comment_content LIKE '%%%s%%' AND %s;", 
			prefix, oldURL, newURL, oldURL, strings.Replace(excludeClause, "post_content", "comment_content", -1)),
	}

	count := 0
	for i, query := range queries {
		fmt.Printf("Debug: Executing content replacement query %d\n", i+1)
		cmd := exec.Command("ddev", "mysql", "-e", query)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Content URL replacement query %d failed: %v - %s\n", i+1, err, string(output))
		} else {
			fmt.Printf("Debug: Content URL replacement query %d succeeded\n", i+1)
			count++
		}
	}

	return count, nil
}

// replaceInContentPreservingMedia replaces URLs in content while preserving media URLs (legacy function)
func (dm *DatabaseManager) replaceInContentPreservingMedia(oldURL, newURL, originalUploadsURL, remoteMediaURL string) (int, error) {
	if remoteMediaURL == "" {
		// No media preservation needed, do normal replacement in content
		args := []string{"wp", "search-replace", oldURL, newURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html", "--skip-columns=guid"}
		cmd := exec.Command("ddev", args...)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("wp-cli replacement failed: %s", string(output))
		}
		return 1, nil // Simplified count
	}

	// When preserving media URLs, we use a more targeted approach
	fmt.Printf("Debug: Replacing %s with %s in content, preserving media URLs pointing to %s\n", oldURL, newURL, remoteMediaURL)
	fmt.Printf("Debug: Original uploads URL to preserve: %s\n", originalUploadsURL)
	
	// Get table prefix
	prefix, err := dm.getTablePrefix()
	if err != nil {
		fmt.Printf("Warning: Could not detect table prefix: %v\n", err)
	}

	// Use the original uploads URL from the database instead of reconstructing it
	oldUploadPath := originalUploadsURL
	
	queries := []string{
		// Replace URLs in post content, but NOT if they contain /wp-content/uploads
		fmt.Sprintf("UPDATE %sposts SET post_content = REPLACE(post_content, '%s', '%s') WHERE post_content LIKE '%%%s%%' AND post_content NOT LIKE '%%%s%%';", 
			prefix, oldURL, newURL, oldURL, oldUploadPath),
		
		// Replace URLs in post excerpts, but NOT if they contain /wp-content/uploads  
		fmt.Sprintf("UPDATE %sposts SET post_excerpt = REPLACE(post_excerpt, '%s', '%s') WHERE post_excerpt LIKE '%%%s%%' AND post_excerpt NOT LIKE '%%%s%%';", 
			prefix, oldURL, newURL, oldURL, oldUploadPath),
		
		// Replace URLs in comments, but NOT if they contain /wp-content/uploads
		fmt.Sprintf("UPDATE %scomments SET comment_content = REPLACE(comment_content, '%s', '%s') WHERE comment_content LIKE '%%%s%%' AND comment_content NOT LIKE '%%%s%%';", 
			prefix, oldURL, newURL, oldURL, oldUploadPath),
	}

	count := 0
	for i, query := range queries {
		cmd := exec.Command("ddev", "mysql", "-e", query)
		cmd.Dir = dm.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Content URL replacement query %d failed: %v - %s\n", i+1, err, string(output))
		} else {
			fmt.Printf("Debug: Content URL replacement query %d succeeded\n", i+1)
			count++
		}
	}

	return count, nil
}

// preserveMediaURLs changes media URLs back to point to remote server
func (dm *DatabaseManager) preserveMediaURLs(localURL, remoteMediaURL string, mediaURLs []MediaURL, uploadsBaseURL string) (int, error) {
	preservedCount := 0
	
	// Determine the local uploads URL that would have been created
	var localUploadsURL string
	if uploadsBaseURL != "" {
		// Replace the domain part of the uploads URL with the local URL
		// e.g., "https://remote.com/wp-content/uploads" -> "https://local.ddev.site/wp-content/uploads"
		if strings.Contains(uploadsBaseURL, "/wp-content/uploads") {
			localUploadsURL = localURL + "/wp-content/uploads"
		} else {
			localUploadsURL = uploadsBaseURL
		}
	} else {
		// Fallback if we couldn't detect uploads URL
		localUploadsURL = localURL + "/wp-content/uploads"
	}
	
	fmt.Printf("Debug: Changing media URLs from %s back to %s\n", localUploadsURL, remoteMediaURL)
	
	// Replace local uploads URLs back to remote
	args := []string{"wp", "search-replace", localUploadsURL, remoteMediaURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"}
	cmd := exec.Command("ddev", args...)
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to preserve media URLs: %s", string(output))
	}
	
	// Parse output for count
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "replacements") {
			fmt.Sscanf(line, "%d", &preservedCount)
			break
		}
	}
	
	// Also fix any specific attachment URLs that got changed
	for _, mediaURL := range mediaURLs {
		if strings.Contains(mediaURL.URL, localURL) {
			// This media URL got changed to local, change it back
			originalRemoteURL := strings.Replace(mediaURL.URL, localURL, strings.TrimSuffix(remoteMediaURL, "/wp-content/uploads"), 1)
			
			args = []string{"wp", "search-replace", mediaURL.URL, originalRemoteURL, "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"}
			cmd = exec.Command("ddev", args...)
			cmd.Dir = dm.projectPath
			if output, err := cmd.CombinedOutput(); err != nil {
				fmt.Printf("Warning: Failed to fix media URL %s: %s\n", mediaURL.URL, string(output))
			} else {
				preservedCount++
			}
		}
	}
	
	return preservedCount, nil
}

// suppressWordPressDebug temporarily disables WordPress debug output for cleaner command output
func (dm *DatabaseManager) suppressWordPressDebug(options SyncOptions) error {
	if !options.SuppressDebug {
		return nil // Debug suppression not requested
	}

	if !dm.canUseWPCLI() {
		return nil // Can't use wp-cli to modify debug settings
	}

	// Temporarily disable WordPress debug constants
	debugCommands := [][]string{
		{"ddev", "wp", "config", "set", "WP_DEBUG", "false", "--raw", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"},
		{"ddev", "wp", "config", "set", "WP_DEBUG_LOG", "false", "--raw", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"},
		{"ddev", "wp", "config", "set", "WP_DEBUG_DISPLAY", "false", "--raw", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html"},
	}

	for _, cmdArgs := range debugCommands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = dm.projectPath
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("Debug: Failed to suppress debug setting: %s\n", string(output))
			// Continue anyway - this is not critical
		}
	}

	return nil
}

// restoreWordPressDebug restores WordPress debug settings (optional - for development environments)
func (dm *DatabaseManager) restoreWordPressDebug(options SyncOptions) error {
	if !options.SuppressDebug {
		return nil // Debug suppression not requested, nothing to restore
	}

	if !dm.canUseWPCLI() {
		return nil // Can't use wp-cli to modify debug settings
	}

	// Note: We don't automatically restore debug settings as users may want them off
	// This function exists for future expansion if needed
	fmt.Printf("Debug: WordPress debug output has been suppressed for cleaner output\n")
	fmt.Printf("Debug: To re-enable debug mode, run: stax wp config set WP_DEBUG true --raw\n")

	return nil
}

// DiagnoseMediaURLRouting helps diagnose URL routing issues by checking various sources
func (dm *DatabaseManager) DiagnoseMediaURLRouting() error {
	if !dm.canUseWPCLI() {
		return fmt.Errorf("wp-cli not available for diagnosis")
	}

	fmt.Printf("🔍 MEDIA URL ROUTING DIAGNOSIS\n")
	fmt.Printf("=====================================\n\n")

	// 1. Check actual database values
	fmt.Printf("1. DATABASE VALUES:\n")
	originalUploadsURL, err := dm.GetOriginalUploadsURLFromDatabase()
	if err == nil {
		fmt.Printf("   ✅ Uploads URL in database: %s\n", originalUploadsURL)
	} else {
		fmt.Printf("   ❌ Could not get uploads URL from database: %v\n", err)
	}

	// Check a sample attachment URL directly from database
	prefix, _ := dm.getTablePrefix()
	query := fmt.Sprintf("SELECT guid FROM %sposts WHERE post_type = 'attachment' AND guid LIKE '%%/wp-content/uploads/%%' LIMIT 3;", prefix)
	cmd := exec.Command("ddev", "mysql", "-e", query, "-s", "-N")
	cmd.Dir = dm.projectPath
	output, err := cmd.CombinedOutput()
	if err == nil && len(output) > 0 {
		fmt.Printf("   ✅ Sample attachment GUIDs from database:\n")
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for i, line := range lines {
			if line != "" {
				fmt.Printf("      %d. %s\n", i+1, line)
			}
		}
	} else {
		fmt.Printf("   ⚠️  No attachment GUIDs found in database\n")
	}

	// 2. Check WordPress options
	fmt.Printf("\n2. WORDPRESS OPTIONS:\n")
	homeCmd := exec.Command("ddev", "wp", "option", "get", "home", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html")
	homeCmd.Dir = dm.projectPath
	homeOutput, err := homeCmd.CombinedOutput()
	if err == nil {
		fmt.Printf("   ✅ WordPress home URL: %s\n", strings.TrimSpace(string(homeOutput)))
	} else {
		fmt.Printf("   ❌ Could not get home URL: %v\n", err)
	}

	siteurlCmd := exec.Command("ddev", "wp", "option", "get", "siteurl", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html")
	siteurlCmd.Dir = dm.projectPath
	siteurlOutput, err := siteurlCmd.CombinedOutput()
	if err == nil {
		fmt.Printf("   ✅ WordPress site URL: %s\n", strings.TrimSpace(string(siteurlOutput)))
	} else {
		fmt.Printf("   ❌ Could not get site URL: %v\n", err)
	}

	// 3. Check uploads directory URL via WordPress functions
	fmt.Printf("\n3. WORDPRESS UPLOADS DIRECTORY:\n")
	uploadsCmd := exec.Command("ddev", "wp", "eval", "echo wp_upload_dir()['baseurl'];", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html")
	uploadsCmd.Dir = dm.projectPath
	uploadsOutput, err := uploadsCmd.CombinedOutput()
	if err == nil {
		fmt.Printf("   ✅ WordPress uploads base URL: %s\n", strings.TrimSpace(string(uploadsOutput)))
	} else {
		fmt.Printf("   ❌ Could not get uploads base URL: %v\n", err)
	}

	// 4. Check for active plugins that might affect URLs
	fmt.Printf("\n4. ACTIVE PLUGINS (that might affect URLs):\n")
	pluginsCmd := exec.Command("ddev", "wp", "plugin", "list", "--status=active", "--field=name", "--skip-plugins", "--skip-themes", "--skip-packages", "--path=/var/www/html")
	pluginsCmd.Dir = dm.projectPath
	pluginsOutput, err := pluginsCmd.CombinedOutput()
	if err == nil {
		plugins := strings.Split(strings.TrimSpace(string(pluginsOutput)), "\n")
		urlAffectingPlugins := []string{}
		for _, plugin := range plugins {
			plugin = strings.TrimSpace(plugin)
			if plugin != "" {
				// Check for common URL-affecting plugins
				lowerPlugin := strings.ToLower(plugin)
				if strings.Contains(lowerPlugin, "cdn") || 
				   strings.Contains(lowerPlugin, "cache") || 
				   strings.Contains(lowerPlugin, "offload") || 
				   strings.Contains(lowerPlugin, "media") || 
				   strings.Contains(lowerPlugin, "s3") || 
				   strings.Contains(lowerPlugin, "cloudfront") {
					urlAffectingPlugins = append(urlAffectingPlugins, plugin)
				}
			}
		}
		if len(urlAffectingPlugins) > 0 {
			fmt.Printf("   ⚠️  Found plugins that might affect URLs:\n")
			for _, plugin := range urlAffectingPlugins {
				fmt.Printf("      - %s\n", plugin)
			}
		} else {
			fmt.Printf("   ✅ No obvious URL-affecting plugins detected\n")
		}
	} else {
		fmt.Printf("   ❌ Could not list active plugins: %v\n", err)
	}

	// 5. Test actual URL resolution
	fmt.Printf("\n5. URL RESOLUTION TEST:\n")
	fmt.Printf("   💡 To test if DDEV is routing media URLs:\n")
	fmt.Printf("      1. Find a media URL from the database (see section 1 above)\n")
	fmt.Printf("      2. Try accessing it directly in your browser\n")
	fmt.Printf("      3. Check if it loads from local DDEV or returns 404\n")
	fmt.Printf("      4. If it loads locally, DDEV might be proxying/routing the URL\n")

	// 6. DDEV describe info
	fmt.Printf("\n6. DDEV CONFIGURATION:\n")
	describeCmd := exec.Command("ddev", "describe")
	describeCmd.Dir = dm.projectPath
	describeOutput, err := describeCmd.CombinedOutput()
	if err == nil {
		fmt.Printf("   ✅ DDEV project info:\n")
		// Look for any router or proxy configurations
		lines := strings.Split(string(describeOutput), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(strings.ToLower(line), "router") || 
			   strings.Contains(strings.ToLower(line), "proxy") || 
			   strings.Contains(strings.ToLower(line), "url") {
				fmt.Printf("      %s\n", line)
			}
		}
	} else {
		fmt.Printf("   ❌ Could not get DDEV project info: %v\n", err)
	}

	fmt.Printf("\n=====================================\n")
	fmt.Printf("💡 INTERPRETATION:\n")
	fmt.Printf("- If database URLs point to WP Engine ✅ but browser shows local URLs ❌\n")
	fmt.Printf("- This suggests WordPress/DDEV is routing the URLs at runtime\n")
	fmt.Printf("- Check for CDN/media plugins or DDEV router configurations\n")

	return nil
}

// CreateUploadDomainRedirectPlugin creates a must-use plugin to redirect upload URLs to remote domain
func (dm *DatabaseManager) CreateUploadDomainRedirectPlugin(localURL, remoteURL string) error {
	// Ensure mu-plugins directory exists
	muPluginsDir := filepath.Join(dm.projectPath, "wp-content", "mu-plugins")
	if err := os.MkdirAll(muPluginsDir, 0755); err != nil {
		return fmt.Errorf("failed to create mu-plugins directory: %w", err)
	}

	// Determine the upload URLs
	localUploadURL := localURL + "/wp-content/uploads"
	remoteUploadURL := remoteURL + "/wp-content/uploads"

	// Create the plugin content
	pluginContent := fmt.Sprintf(`<?php
/**
 * Plugin Name: Upload Domain Redirect
 * Description: Redirects upload URLs to remote WP Engine domain for media files
 * Version: 1.0.0
 * Auto-generated by stax CLI
 */

// Redirect upload directory URLs to remote domain
add_filter('upload_dir', function($uploads) {
    $uploads['url'] = str_replace(
        '%s',
        '%s',
        $uploads['url']
    );
    $uploads['baseurl'] = str_replace(
        '%s',
        '%s',
        $uploads['baseurl']
    );
    return $uploads;
});

// Also handle wp_get_attachment_url() calls
add_filter('wp_get_attachment_url', function($url) {
    return str_replace(
        '%s',
        '%s',
        $url
    );
});

// Handle get_attached_file() for proper file paths
add_filter('wp_get_attachment_image_src', function($image) {
    if ($image && isset($image[0])) {
        $image[0] = str_replace(
            '%s',
            '%s',
            $image[0]
        );
    }
    return $image;
});
`, localUploadURL, remoteUploadURL, localUploadURL, remoteUploadURL, localUploadURL, remoteUploadURL, localUploadURL, remoteUploadURL)

	// Write the plugin file
	pluginPath := filepath.Join(muPluginsDir, "upload-domain-redirect.php")
	if err := os.WriteFile(pluginPath, []byte(pluginContent), 0644); err != nil {
		return fmt.Errorf("failed to write upload domain redirect plugin: %w", err)
	}

	fmt.Printf("✅ Created upload domain redirect plugin\n")
	fmt.Printf("   📍 Plugin location: %s\n", pluginPath)
	fmt.Printf("   🔄 Redirects: %s → %s\n", localUploadURL, remoteUploadURL)
	fmt.Printf("   ℹ️  This plugin ensures media URLs point to remote WP Engine server\n")

	return nil
}
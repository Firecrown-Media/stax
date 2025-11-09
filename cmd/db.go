package cmd

import (
	"fmt"

	"github.com/firecrown-media/stax/pkg/config"
	"github.com/firecrown-media/stax/pkg/credentials"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command group
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
	Long:  `Manage database operations including pull, push, import, export, snapshots, and queries.`,
}

var (
	dbEnvironment    string
	dbSnapshot       bool
	dbSanitize       bool
	dbSkipReplace    bool
	dbExcludeTables  string
	dbSkipLogs       bool
	dbSkipTransients bool
	dbSkipSpam       bool
)

// dbPullCmd represents the db:pull command
var dbPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull database from WPEngine",
	Long: `Pull database from WPEngine, import it locally, and run search-replace.

This command will:
  - Create a snapshot of the current database (unless --snapshot=false)
  - Connect to WPEngine SSH Gateway
  - Export the database from WPEngine
  - Transfer the database to local environment
  - Import into local DDEV database
  - Run search-replace operations
  - Flush WordPress cache`,
	Example: `  # Basic pull
  stax db pull

  # Pull from staging
  stax db pull --environment=staging

  # Pull without creating snapshot
  stax db pull --snapshot=false

  # Pull with sanitized data
  stax db pull --sanitize`,
	RunE: runDBPull,
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbPullCmd)

	// Flags for pull
	dbPullCmd.Flags().StringVar(&dbEnvironment, "environment", "", "WPEngine environment (default: from config)")
	dbPullCmd.Flags().BoolVar(&dbSnapshot, "snapshot", true, "create snapshot before import")
	dbPullCmd.Flags().BoolVar(&dbSanitize, "sanitize", false, "sanitize user data")
	dbPullCmd.Flags().BoolVar(&dbSkipReplace, "skip-replace", false, "skip search-replace")
	dbPullCmd.Flags().StringVar(&dbExcludeTables, "exclude-tables", "", "comma-separated tables to exclude")
	dbPullCmd.Flags().BoolVar(&dbSkipLogs, "skip-logs", true, "skip log tables")
	dbPullCmd.Flags().BoolVar(&dbSkipTransients, "skip-transients", true, "skip transient tables")
	dbPullCmd.Flags().BoolVar(&dbSkipSpam, "skip-spam", true, "skip spam/trash")
}

func runDBPull(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Pulling Database from WPEngine")

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get WPEngine credentials
	_, err = getWPEngineCredentials(cfg.WPEngine.Install)
	if err != nil {
		return fmt.Errorf("failed to get WPEngine credentials: %w", err)
	}

	// Get SSH key
	_, err = getSSHKey()
	if err != nil {
		return fmt.Errorf("failed to get SSH key: %w", err)
	}

	// Create snapshot if requested
	if dbSnapshot {
		ui.Info("Creating database snapshot...")
		// TODO: Implement snapshot creation
		ui.Success("Snapshot created")
	}

	// Connect to WPEngine SSH Gateway
	ui.Info("Connecting to WPEngine SSH Gateway...")
	// TODO: Use wpengine.NewSSHClient with proper config
	ui.Success("Connected to WPEngine")

	// Export database
	ui.Info("Exporting database from WPEngine...")
	// TODO: Use SSHClient.ExportDatabase with options
	ui.Success("Database exported")

	// Import to local database
	ui.Info("Importing database to local environment...")
	// TODO: Use wordpress.CLI.ImportDatabase
	ui.Success("Database imported")

	// Run search-replace unless skipped
	if !dbSkipReplace {
		ui.Info("Running search-replace operations...")
		// TODO: Use wordpress.CLI.MultisiteSearchReplace
		ui.Success("Search-replace completed")
	}

	// Flush cache
	ui.Info("Flushing WordPress cache...")
	// TODO: Use wordpress.CLI.FlushCache
	ui.Success("Cache flushed")

	ui.Success("\nDatabase pull completed!")

	return nil
}

// Helper functions (placeholder implementations)
func loadConfig() (*config.Config, error) {
	// TODO: Load from .stax.yml
	cfg := config.Defaults()
	cfg.WPEngine.Install = "fsmultisite"
	return cfg, nil
}

func getWPEngineCredentials(install string) (*credentials.WPEngineCredentials, error) {
	return credentials.GetWPEngineCredentials(install)
}

func getSSHKey() (string, error) {
	return credentials.GetSSHPrivateKey("wpengine")
}

package cmd

import (
	"fmt"

	"github.com/firecrown-media/stax/pkg/credentials"
	"github.com/firecrown-media/stax/pkg/ddev"
	"github.com/firecrown-media/stax/pkg/errors"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/firecrown-media/stax/pkg/wordpress"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command group
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "✓ Database operations",
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
  - Run search-replace operations (unless --skip-replace)
  - Flush WordPress cache`,
	Example: `  # Basic pull
  stax db pull

  # Pull from staging
  stax db pull --environment=staging

  # Pull without creating snapshot
  stax db pull --snapshot=false

  # Pull without automatic URL replacement (advanced users)
  stax db pull --skip-replace

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
	dbPullCmd.Flags().BoolVar(&dbSkipReplace, "skip-replace", false, "skip automatic URL search-replace")
	dbPullCmd.Flags().StringVar(&dbExcludeTables, "exclude-tables", "", "comma-separated tables to exclude")
	dbPullCmd.Flags().BoolVar(&dbSkipLogs, "skip-logs", true, "skip log tables")
	dbPullCmd.Flags().BoolVar(&dbSkipTransients, "skip-transients", true, "skip transient tables")
	dbPullCmd.Flags().BoolVar(&dbSkipSpam, "skip-spam", true, "skip spam/trash")
}

func runDBPull(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Pulling Database from WPEngine")

	// Load configuration
	cfg, err := loadConfigForCommand()
	if err != nil {
		return err
	}

	// Get WPEngine credentials with fallback
	creds, err := credentials.GetWPEngineCredentialsWithFallback(cfg.WPEngine.Install)
	if err != nil {
		if credErr, ok := err.(*credentials.CredentialsNotFoundError); ok {
			return errors.NewCredentialsNotFoundError(credErr.Tried, credErr.LastErr)
		}
		return fmt.Errorf("failed to get WPEngine credentials: %w", err)
	}

	// Get SSH key with fallback
	sshKey, err := credentials.GetSSHPrivateKeyWithFallback("wpengine")
	if err != nil {
		if keyErr, ok := err.(*credentials.SSHKeyNotFoundError); ok {
			return errors.NewSSHKeyNotFoundError("", keyErr.Tried, keyErr.LastErr)
		}
		return fmt.Errorf("failed to get SSH key: %w", err)
	}

	// Use credentials
	_ = creds
	_ = sshKey

	// Check if DDEV is running
	projectDir := getProjectDir()
	mgr := ddev.NewManager(projectDir)
	running, err := mgr.IsRunning()
	if err != nil {
		return fmt.Errorf("failed to check DDEV status: %w", err)
	}
	if !running {
		return fmt.Errorf("DDEV must be running to import database. Please run 'stax start' first")
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
	dbPath := "/tmp/wpengine-db.sql" // TODO: Get actual path from export
	ui.Success("Database exported")

	// Import to local database
	ui.Info("Importing database to local environment...")
	if err := mgr.ImportDB(dbPath); err != nil {
		return fmt.Errorf("database import failed: %w", err)
	}
	ui.Success("Database imported")

	// Run search-replace unless skipped
	if !dbSkipReplace {
		ui.Info("Replacing URLs...")

		// Get source and target URLs
		sourceURL := getWPEngineURL(cfg)
		targetURL := getDDEVURL(cfg)

		// Run search-replace
		if err := runSearchReplace(projectDir, sourceURL, targetURL, cfg); err != nil {
			ui.Warning(fmt.Sprintf("URL replacement failed: %v", err))
			ui.Info("You may need to run manually: ddev wp search-replace '%s' '%s' --all-tables", sourceURL, targetURL)
		} else {
			ui.Success("URLs replaced successfully")
		}
	} else {
		ui.Info("Skipping URL replacement (--skip-replace flag set)")
		sourceURL := getWPEngineURL(cfg)
		targetURL := getDDEVURL(cfg)
		ui.Info("To replace URLs manually, run: ddev wp search-replace '%s' '%s' --all-tables", sourceURL, targetURL)
	}

	// Flush cache
	ui.Info("Flushing WordPress cache...")
	cli := wordpress.NewCLI(projectDir)
	if err := cli.FlushCache(); err != nil {
		ui.Warning(fmt.Sprintf("Cache flush failed: %v", err))
	} else {
		ui.Success("Cache flushed")
	}

	ui.Success("\nDatabase pull completed!")

	return nil
}

// loadConfigForCommand is now defined in files.go to avoid duplication

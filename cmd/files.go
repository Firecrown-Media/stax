package cmd

import (
	"fmt"
	"strings"

	"github.com/firecrown-media/stax/pkg/config"
	"github.com/firecrown-media/stax/pkg/credentials"
	"github.com/firecrown-media/stax/pkg/errors"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/firecrown-media/stax/pkg/wpengine"
	"github.com/spf13/cobra"
)

// filesCmd represents the files command group
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "✓ File synchronization operations",
	Long:  `Manage file synchronization between WPEngine and local environment including themes, plugins, and uploads.`,
}

var (
	filesEnvironment         string
	filesThemesOnly          bool
	filesPluginsOnly         bool
	filesMuPluginsOnly       bool
	filesExcludeUploads      bool
	filesDryRun              bool
	filesDelete              bool
	filesBandwidthLimit      int
	filesInclude             string
	filesExclude             string
	filesVerify              bool
	filesPreservePermissions bool
)

// filesPullCmd represents the files:pull command
var filesPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull files from WPEngine",
	Long: `Pull files from WPEngine to your local environment.

This command will:
  - Connect to WPEngine via SSH
  - Sync wp-content directory (or specific subdirectories)
  - Transfer files using rsync over SSH
  - Verify file integrity after transfer

By default, this syncs the entire wp-content directory. Use flags to
limit the sync to specific directories like themes or plugins.`,
	Example: `  # Basic pull (all wp-content)
  stax files pull

  # Pull only themes
  stax files pull --themes-only

  # Pull only plugins
  stax files pull --plugins-only

  # Pull only mu-plugins
  stax files pull --mu-plugins-only

  # Pull without uploads directory
  stax files pull --exclude-uploads

  # Pull from staging environment
  stax files pull --environment=staging

  # Dry run to see what would be synced
  stax files pull --dry-run

  # Delete local files not on remote
  stax files pull --delete

  # Limit bandwidth to 1000 KB/s
  stax files pull --bandwidth-limit=1000

  # Preserve file permissions
  stax files pull --preserve-permissions

  # Custom includes and excludes
  stax files pull --include="*.php,*.js" --exclude="*.log,cache/"`,
	RunE: runFilesPull,
}

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(filesPullCmd)

	// Flags for pull
	filesPullCmd.Flags().StringVar(&filesEnvironment, "environment", "", "WPEngine environment (default: from config)")
	filesPullCmd.Flags().BoolVar(&filesThemesOnly, "themes-only", false, "sync only themes directory")
	filesPullCmd.Flags().BoolVar(&filesPluginsOnly, "plugins-only", false, "sync only plugins directory")
	filesPullCmd.Flags().BoolVar(&filesMuPluginsOnly, "mu-plugins-only", false, "sync only mu-plugins directory")
	filesPullCmd.Flags().BoolVar(&filesExcludeUploads, "exclude-uploads", false, "exclude uploads directory")
	filesPullCmd.Flags().BoolVar(&filesDryRun, "dry-run", false, "show what would be transferred without syncing")
	filesPullCmd.Flags().BoolVar(&filesDelete, "delete", false, "delete local files not present on remote")
	filesPullCmd.Flags().IntVar(&filesBandwidthLimit, "bandwidth-limit", 0, "bandwidth limit in KB/s (0 = unlimited)")
	filesPullCmd.Flags().StringVar(&filesInclude, "include", "", "comma-separated patterns to include")
	filesPullCmd.Flags().StringVar(&filesExclude, "exclude", "", "comma-separated patterns to exclude")
	filesPullCmd.Flags().BoolVar(&filesVerify, "verify", false, "verify file checksums after sync (slower for large sites)")
	filesPullCmd.Flags().BoolVar(&filesPreservePermissions, "preserve-permissions", false, "preserve file permissions during sync")
}

func runFilesPull(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Pulling Files from WPEngine")

	// Load configuration
	cfg, err := loadConfigForCommand()
	if err != nil {
		return err
	}

	// Determine environment
	environment := filesEnvironment
	if environment == "" {
		environment = cfg.WPEngine.Environment
	}

	ui.Info(fmt.Sprintf("Environment: %s", environment))
	ui.Info(fmt.Sprintf("Install: %s", cfg.WPEngine.Install))

	// Get credentials with fallback
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

	// Create SSH client
	ui.Info("Connecting to WPEngine SSH Gateway...")
	sshConfig := wpengine.SSHConfig{
		Host:       cfg.WPEngine.SSHGateway,
		Port:       22,
		User:       creds.SSHUser,
		PrivateKey: sshKey,
		Install:    cfg.WPEngine.Install,
	}

	sshClient, err := wpengine.NewSSHClient(sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to WPEngine: %w", err)
	}
	defer sshClient.Close()

	ui.Success("Connected to WPEngine")

	// Build sync options
	syncOptions := buildSyncOptions(cfg)

	// Determine what to sync
	var remotePath, localPath string
	if filesThemesOnly {
		ui.Info("Syncing themes only...")
		remotePath = fmt.Sprintf("/sites/%s/wp-content/themes/", cfg.WPEngine.Install)
		localPath = getProjectDir() + "/wp-content/themes/"
	} else if filesPluginsOnly {
		ui.Info("Syncing plugins only...")
		remotePath = fmt.Sprintf("/sites/%s/wp-content/plugins/", cfg.WPEngine.Install)
		localPath = getProjectDir() + "/wp-content/plugins/"
	} else if filesMuPluginsOnly {
		ui.Info("Syncing mu-plugins only...")
		remotePath = fmt.Sprintf("/sites/%s/wp-content/mu-plugins/", cfg.WPEngine.Install)
		localPath = getProjectDir() + "/wp-content/mu-plugins/"
	} else {
		ui.Info("Syncing wp-content directory...")
		remotePath = fmt.Sprintf("/sites/%s/wp-content/", cfg.WPEngine.Install)
		localPath = getProjectDir() + "/wp-content/"
	}

	// Execute sync
	if filesDryRun {
		ui.Info("DRY RUN - No files will be transferred")
	}

	ui.Info("Starting file synchronization...")
	if err := sshClient.SyncDirectory(remotePath, localPath, syncOptions); err != nil {
		return fmt.Errorf("file sync failed: %w", err)
	}

	if !filesDryRun {
		ui.Success("Files synchronized successfully")

		// Verify integrity if not a dry run
		ui.Info("Verifying file integrity...")
		if err := sshClient.VerifyFileIntegrity(remotePath, localPath); err != nil {
			ui.Warning(fmt.Sprintf("File integrity check failed: %v", err))
			ui.Info("Files were transferred but counts may differ (this is usually OK)")
		} else {
			ui.Success("File integrity verified")
		}

		// Perform checksum verification if requested
		if filesVerify {
			ui.Section("Checksum Verification")
			ui.Info("Generating checksums (this may take a while for large sites)...")

			spinner := ui.NewSpinner("Verifying checksums...")
			spinner.Start()

			result, err := sshClient.VerifyFileChecksums(remotePath, localPath)
			spinner.Stop()

			if err != nil {
				ui.Warning(fmt.Sprintf("Checksum verification failed: %v", err))
			} else {
				// Print verification results
				printChecksumResults(result)
			}
		}
	} else {
		ui.Info("Dry run completed")
	}

	ui.Success("\nFile pull completed!")

	return nil
}

// buildSyncOptions builds rsync sync options from command flags
func buildSyncOptions(cfg *config.Config) wpengine.SyncOptions {
	options := wpengine.SyncOptions{
		DryRun:              filesDryRun,
		Delete:              filesDelete,
		BandwidthLimit:      filesBandwidthLimit,
		Progress:            true,
		PreservePermissions: filesPreservePermissions,
		Include:             []string{},
		Exclude:             wpengine.GetExcludePatterns(),
		ProjectDir:          getProjectDir(), // Enable .staxignore support
	}

	// Add custom includes
	if filesInclude != "" {
		options.Include = strings.Split(filesInclude, ",")
	}

	// Add custom excludes
	if filesExclude != "" {
		customExcludes := strings.Split(filesExclude, ",")
		options.Exclude = append(options.Exclude, customExcludes...)
	}

	// Exclude uploads if requested
	if filesExcludeUploads {
		options.Exclude = append(options.Exclude, "uploads/")
	}

	// Apply bandwidth limit from config if not specified via flag
	if filesBandwidthLimit == 0 && cfg.Performance.RsyncBandwidthLimit > 0 {
		options.BandwidthLimit = cfg.Performance.RsyncBandwidthLimit
	}

	return options
}

// loadConfigForCommand loads configuration for a command
func loadConfigForCommand() (*config.Config, error) {
	cfg, err := config.Load(cfgFile, projectDir)
	if err != nil {
		return nil, errors.NewConfigNotFoundError(cfgFile, err)
	}
	return cfg, nil
}

// printChecksumResults prints the results of checksum verification
func printChecksumResults(result *wpengine.ChecksumResult) {
	// Print summary
	ui.Info(fmt.Sprintf("Total files checked: %d", result.TotalFiles))
	ui.Success(fmt.Sprintf("Matched files: %d", result.MatchedFiles))

	// Print mismatches if any
	if result.MismatchedFiles > 0 {
		ui.Warning(fmt.Sprintf("Mismatched checksums: %d", result.MismatchedFiles))
		ui.Info("Files with different checksums:")
		for i, mismatch := range result.Mismatches {
			if i >= 10 {
				ui.Info(fmt.Sprintf("  ... and %d more", len(result.Mismatches)-10))
				break
			}
			ui.Info(fmt.Sprintf("  - %s", mismatch.RelativePath))
			ui.Verbose(fmt.Sprintf("    Remote: %s", mismatch.RemoteChecksum))
			ui.Verbose(fmt.Sprintf("    Local:  %s", mismatch.LocalChecksum))
		}
	}

	// Print missing local files if any
	if result.MissingLocal > 0 {
		ui.Warning(fmt.Sprintf("Missing locally: %d", result.MissingLocal))
		ui.Info("Files that exist remotely but not locally:")
		for i, file := range result.MissingLocally {
			if i >= 10 {
				ui.Info(fmt.Sprintf("  ... and %d more", len(result.MissingLocally)-10))
				break
			}
			ui.Info(fmt.Sprintf("  - %s", file))
		}
	}

	// Print missing remote files if any
	if result.MissingRemote > 0 {
		ui.Warning(fmt.Sprintf("Missing remotely: %d", result.MissingRemote))
		ui.Info("Files that exist locally but not remotely:")
		for i, file := range result.MissingRemotely {
			if i >= 10 {
				ui.Info(fmt.Sprintf("  ... and %d more", len(result.MissingRemotely)-10))
				break
			}
			ui.Info(fmt.Sprintf("  - %s", file))
		}
	}

	// Final verdict
	if result.MismatchedFiles == 0 && result.MissingLocal == 0 && result.MissingRemote == 0 {
		ui.Success("All files verified successfully - checksums match!")
	} else {
		ui.Warning("Some files have checksum differences - review the details above")
	}
}

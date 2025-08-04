package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/Firecrown-Media/stax/pkg/wpengine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var wpengineCmd = &cobra.Command{
	Use:   "wpe",
	Short: "WP Engine integration commands",
	Long: `Commands for integrating with WP Engine hosted sites.
Allows you to sync databases and files from WP Engine to local DDEV environment.`,
}

var wpeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List WP Engine installs",
	Long:  `List all WP Engine installs you have access to and optionally select one interactively.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create minimal config for API access
		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			username = viper.GetString("wpengine.username")
			if username == "" {
				username = os.Getenv("WPE_USERNAME")
			}
		}
		
		if username == "" {
			return fmt.Errorf("WP Engine username is required. Set via --username flag, config file, or WPE_USERNAME environment variable")
		}

		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			apiKey = viper.GetString("wpengine.api_key")
			if apiKey == "" {
				apiKey = os.Getenv("WPE_API_KEY")
			}
		}
		
		config := wpengine.Config{
			Username: username,
			APIKey:   apiKey,
		}
		
		client := wpengine.NewClient(config)
		
		fmt.Printf("📋 Listing WP Engine installs for %s...\n\n", username)
		
		installs, err := client.ListInstalls()
		if err != nil {
			return fmt.Errorf("failed to list installs: %w", err)
		}
		
		if len(installs) == 0 {
			fmt.Printf("No installs found for user %s\n", username)
			return nil
		}
		
		fmt.Printf("Found %d install(s):\n\n", len(installs))
		for i, install := range installs {
			fmt.Printf("%d. %s\n", i+1, install.Name)
			fmt.Printf("   Environment: %s\n", install.Environment)
			fmt.Printf("   Domain: %s\n", install.Domain)
			fmt.Printf("   Status: %s\n", install.Status)
			fmt.Printf("   PHP Version: %s\n", install.PHPVersion)
			fmt.Printf("\n")
		}
		
		return nil
	},
}

var wpeConnectCmd = &cobra.Command{
	Use:   "connect [install-name]",
	Short: "Connect to a WP Engine install",
	Long: `Connect to a WP Engine install and test the connection.
This will validate your credentials and install access.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := getWPEngineConfig(cmd, args)
		if err != nil {
			return err
		}

		client := wpengine.NewClient(*config)

		fmt.Printf("Connecting to WP Engine install: %s\n", config.InstallName)

		if err := client.TestConnection(); err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}

		info, err := client.GetInstallInfo()
		if err != nil {
			return fmt.Errorf("failed to get install info: %w", err)
		}

		fmt.Printf("✅ Successfully connected to %s\n", config.InstallName)
		fmt.Printf("   Environment: %s\n", info.Environment)
		fmt.Printf("   Domain: %s\n", info.Domain)
		fmt.Printf("   PHP Version: %s\n", info.PHPVersion)
		fmt.Printf("   Status: %s\n", info.Status)

		return nil
	},
}

var wpeSyncCmd = &cobra.Command{
	Use:   "sync [install-name]",
	Short: "Sync database and files from WP Engine",
	Long: `Sync database and files from WP Engine to local DDEV environment.
This will download the latest database backup and sync WordPress files
with comprehensive URL rewriting including media URLs for local development.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		config, err := getWPEngineConfig(cmd, args)
		if err != nil {
			return err
		}

		client := wpengine.NewClient(*config)

		// Get sync options
		skipMedia, _ := cmd.Flags().GetBool("skip-media")
		skipFiles, _ := cmd.Flags().GetBool("skip-files")
		skipDatabase, _ := cmd.Flags().GetBool("skip-database")
		deleteLocal, _ := cmd.Flags().GetBool("delete-local")
		suppressDebug, _ := cmd.Flags().GetBool("suppress-debug")
		createUploadRedirect, _ := cmd.Flags().GetBool("create-upload-redirect")

		syncOptions := wpengine.SyncOptions{
			SkipMedia:     skipMedia,
			DeleteLocal:   deleteLocal,
			SuppressDebug: suppressDebug,
			ExcludeDirs: []string{
				"wp-content/uploads/",
				"wp-content/cache/",
				"wp-content/backup/",
			},
		}

		fmt.Printf("🚀 Starting sync from WP Engine install: %s\n", config.InstallName)

		// Ensure DDEV project is started before syncing
		fmt.Printf("🔧 Ensuring DDEV environment is running...\n")
		if err := ensureDDEVStarted(projectPath); err != nil {
			return fmt.Errorf("failed to start DDEV environment: %w", err)
		}

		// Sync database  
		databaseSynced := false
		if !skipDatabase {
			dbFile := filepath.Join(projectPath, "tmp", fmt.Sprintf("wpe-db-%s.sql", time.Now().Format("20060102-150405")))

			_, err := client.DownloadDatabase(dbFile)
			if err != nil {
				return fmt.Errorf("failed to download database: %w", err)
			}

			// Import database (URL rewriting will happen later after file sync)
			dbManager := wpengine.NewDatabaseManager(projectPath, client)
			_, err = dbManager.ImportDatabase(dbFile, syncOptions)
			if err != nil {
				return fmt.Errorf("failed to import database: %w", err)
			}

			fmt.Printf("✅ Database imported successfully\n")
			databaseSynced = true

			// Clean up temporary database file
			os.Remove(dbFile)
		}

		// Sync files
		if !skipFiles {
			filesResult, err := client.SyncFiles(projectPath, syncOptions)
			if err != nil {
				return fmt.Errorf("failed to sync files: %w", err)
			}

			fmt.Printf("📄 Synced files: %d\n", filesResult.SyncedFiles)
		}

		// Perform URL rewriting after file sync if database was synced
		if databaseSynced {
			fmt.Printf("🔄 Performing URL rewriting now that WordPress files are present...\n")
			
			// Get the install info to determine URLs
			info, err := client.GetInstallInfo()
			if err == nil {
				dbManager := wpengine.NewDatabaseManager(projectPath, client)
				remoteURL := fmt.Sprintf("https://%s", info.Domain)
				
				// Get the local project name from DDEV, not the WP Engine install name
				localProjectName := dbManager.GetProjectName()
				localURL := fmt.Sprintf("https://%s.ddev.site", localProjectName)
				
				var rewrittenURLs int
				rewrittenURLs, err = dbManager.RewriteURLsPostSync(remoteURL, localURL, syncOptions)
				if err != nil {
					fmt.Printf("   ⚠️  Warning: URL rewriting failed: %v\n", err)
				} else if rewrittenURLs > 0 {
					fmt.Printf("   ✅ URL rewriting completed (including media URLs)\n")
					fmt.Printf("   📊 Rewritten URLs: %d\n", rewrittenURLs)
				}

				// Create upload domain redirect plugin if requested
				if createUploadRedirect {
					fmt.Printf("🔌 Creating upload domain redirect plugin...\n")
					if err := dbManager.CreateUploadDomainRedirectPlugin(localURL, remoteURL); err != nil {
						fmt.Printf("   ⚠️  Warning: Failed to create upload redirect plugin: %v\n", err)
					}
				}
			}
		}

		// Show media configuration
		if !skipMedia {
			fmt.Printf("🖼️  Media configuration:\n")
			if createUploadRedirect {
				fmt.Printf("   📍 Database URLs: Rewritten to local development URLs\n")
				fmt.Printf("   🔌 WordPress URLs: Redirected to WP Engine via must-use plugin\n")
				fmt.Printf("   ✅ Result: Media loads from WP Engine, database works locally\n")
			} else {
				fmt.Printf("   📍 Media URLs have been rewritten to local development URLs\n")
				fmt.Printf("   ⚠️  Media files will need to exist locally or use --create-upload-redirect flag\n")
			}
			fmt.Printf("   ℹ️  Media files remain on WP Engine server\n")
		}

		fmt.Printf("✅ Sync completed successfully!\n")
		return nil
	},
}

var wpeDbCmd = &cobra.Command{
	Use:   "db [install-name]",
	Short: "Database operations with WP Engine",
	Long:  `Download and import database from WP Engine.`,
}

var wpeDbDownloadCmd = &cobra.Command{
	Use:   "download [install-name]",
	Short: "Download database from WP Engine",
	Long:  `Download the latest database backup from WP Engine.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		config, err := getWPEngineConfig(cmd, args)
		if err != nil {
			return err
		}

		client := wpengine.NewClient(*config)

		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = filepath.Join(projectPath, "tmp", fmt.Sprintf("wpe-db-%s.sql", time.Now().Format("20060102-150405")))
		}

		fmt.Printf("📥 Downloading database from %s...\n", config.InstallName)

		var result *wpengine.DatabaseSyncResult
		result, err = client.DownloadDatabase(outputFile)
		if err != nil {
			return fmt.Errorf("failed to download database: %w", err)
		}

		fmt.Printf("✅ Database downloaded: %s\n", result.DatabaseFile)
		fmt.Printf("   📄 Backup ID: %s\n", result.BackupID)

		return nil
	},
}

var wpeDbImportCmd = &cobra.Command{
	Use:   "import [database-file]",
	Short: "Import WP Engine database to local environment",
	Long:  `Import a WP Engine database file to the local DDEV environment with URL rewriting.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		dbFile := args[0]
		if !filepath.IsAbs(dbFile) {
			dbFile = filepath.Join(projectPath, dbFile)
		}

		// Check if database file exists
		if _, err := os.Stat(dbFile); os.IsNotExist(err) {
			return fmt.Errorf("database file not found: %s", dbFile)
		}

		config, err := getWPEngineConfig(cmd, []string{})
		if err != nil {
			return err
		}

		client := wpengine.NewClient(*config)

		suppressDebug, _ := cmd.Flags().GetBool("suppress-debug")
		createUploadRedirect, _ := cmd.Flags().GetBool("create-upload-redirect")

		syncOptions := wpengine.SyncOptions{
			SkipMedia:     true,
			SuppressDebug: suppressDebug,
		}

		fmt.Printf("📥 Importing database: %s\n", dbFile)

		// Ensure DDEV project is started before importing
		fmt.Printf("🔧 Ensuring DDEV environment is running...\n")
		if err := ensureDDEVStarted(projectPath); err != nil {
			return fmt.Errorf("failed to start DDEV environment: %w", err)
		}

		dbManager := wpengine.NewDatabaseManager(projectPath, client)
		var result *wpengine.DatabaseSyncResult
		result, err = dbManager.ImportDatabase(dbFile, syncOptions)
		if err != nil {
			return fmt.Errorf("failed to import database: %w", err)
		}

		fmt.Printf("✅ Database imported successfully\n")
		fmt.Printf("   📊 Rewritten URLs: %d (including media URLs)\n", result.RewrittenURLs)

		// Create upload domain redirect plugin if requested
		if createUploadRedirect {
			fmt.Printf("🔌 Creating upload domain redirect plugin...\n")
			
			// Get install info to determine remote URL
			info, err := client.GetInstallInfo()
			if err == nil {
				remoteURL := fmt.Sprintf("https://%s", info.Domain)
				localProjectName := dbManager.GetProjectName()
				localURL := fmt.Sprintf("https://%s.ddev.site", localProjectName)
				
				if err := dbManager.CreateUploadDomainRedirectPlugin(localURL, remoteURL); err != nil {
					fmt.Printf("   ⚠️  Warning: Failed to create upload redirect plugin: %v\n", err)
				}
			} else {
				fmt.Printf("   ⚠️  Warning: Could not get install info for plugin creation: %v\n", err)
			}
		}

		return nil
	},
}

var wpeDbAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze media URLs in WordPress database",
	Long:  `Analyze the WordPress database to identify media URLs and show what would be preserved during sync.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		// Ensure DDEV project is started
		if err := ensureDDEVStarted(projectPath); err != nil {
			return fmt.Errorf("failed to start DDEV environment: %w", err)
		}

		// Create a minimal client for database operations
		dbManager := wpengine.NewDatabaseManager(projectPath, nil)
		
		// Get uploads base URL directly from database (more reliable than wp-cli)
		originalUploadsURL, err := dbManager.GetOriginalUploadsURLFromDatabase()
		if err != nil {
			fmt.Printf("⚠️  Could not get original uploads URL from database: %v\n", err)
			fmt.Printf("💡 Media URL analysis not available due to database connection issues\n")
		} else {
			fmt.Printf("📂 Original uploads base URL: %s\n", originalUploadsURL)
			fmt.Printf("\n🖼️  Media URL Preservation Strategy:\n")
			fmt.Printf("   • Pattern-based detection (more reliable than wp-cli attachment queries)\n")
			fmt.Printf("   • Preserves URLs containing: %s\n", originalUploadsURL)
			fmt.Printf("   • Also preserves URLs containing: %s\n", "https://astronomystage.wpengine.com/wp-content/uploads")
			fmt.Printf("   • Non-media URLs will be rewritten to local development URLs\n")
		}

		fmt.Printf("\n💡 During WP Engine sync with --remote-media-url, media URLs will be preserved\n")
		fmt.Printf("   to continue pointing to your WP Engine server instead of being rewritten\n")
		fmt.Printf("   to local URLs.\n")

		return nil
	},
}

var wpeDbRewriteCmd = &cobra.Command{
	Use:   "rewrite [install-name]",
	Short: "Rewrite URLs in database for local development",
	Long:  `Manually rewrite URLs in the database from WP Engine URLs to local DDEV URLs.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		config, err := getWPEngineConfig(cmd, args)
		if err != nil {
			return err
		}

		client := wpengine.NewClient(*config)
		
		suppressDebug, _ := cmd.Flags().GetBool("suppress-debug")
		createUploadRedirect, _ := cmd.Flags().GetBool("create-upload-redirect")

		syncOptions := wpengine.SyncOptions{
			SkipMedia:     true,
			SuppressDebug: suppressDebug,
		}

		// Get URLs
		info, err := client.GetInstallInfo()
		if err != nil {
			return fmt.Errorf("failed to get install info: %w", err)
		}

		remoteURL := fmt.Sprintf("https://%s", info.Domain)
		localURL := fmt.Sprintf("https://%s.ddev.site", info.Name)

		fmt.Printf("🔄 Rewriting URLs from %s to %s\n", remoteURL, localURL)
		
		// Ensure DDEV project is started before rewriting
		fmt.Printf("🔧 Ensuring DDEV environment is running...\n")
		if err := ensureDDEVStarted(projectPath); err != nil {
			return fmt.Errorf("failed to start DDEV environment: %w", err)
		}

		dbManager := wpengine.NewDatabaseManager(projectPath, client)
		var rewrittenURLs int
		rewrittenURLs, err = dbManager.RewriteURLsPostSync(remoteURL, localURL, syncOptions)
		if err != nil {
			return fmt.Errorf("failed to rewrite URLs: %w", err)
		}

		fmt.Printf("✅ URL rewriting completed\n")
		fmt.Printf("📊 Rewritten URLs: %d (including media URLs)\n", rewrittenURLs)

		// Create upload domain redirect plugin if requested
		if createUploadRedirect {
			fmt.Printf("🔌 Creating upload domain redirect plugin...\n")
			
			dbManager := wpengine.NewDatabaseManager(projectPath, client)
			localProjectName := dbManager.GetProjectName()
			localURL := fmt.Sprintf("https://%s.ddev.site", localProjectName)
			
			if err := dbManager.CreateUploadDomainRedirectPlugin(localURL, remoteURL); err != nil {
				fmt.Printf("   ⚠️  Warning: Failed to create upload redirect plugin: %v\n", err)
			}
		}

		return nil
	},
}

var wpeDbDiagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Diagnose media URL routing and rewriting behavior",
	Long:  `Analyze the WordPress database and DDEV configuration to understand how media URLs are being handled and whether they're being routed or rewritten.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		// Ensure DDEV project is started
		if err := ensureDDEVStarted(projectPath); err != nil {
			return fmt.Errorf("failed to start DDEV environment: %w", err)
		}

		// Create a minimal database manager for diagnostics
		dbManager := wpengine.NewDatabaseManager(projectPath, nil)
		
		return dbManager.DiagnoseMediaURLRouting()
	},
}

func getWPEngineConfig(cmd *cobra.Command, args []string) (*wpengine.Config, error) {
	var installName string
	if len(args) > 0 {
		installName = args[0]
	} else {
		var err error
		installName, err = cmd.Flags().GetString("install")
		if err != nil || installName == "" {
			return nil, fmt.Errorf("install name is required")
		}
	}

	username, _ := cmd.Flags().GetString("username")
	environment, _ := cmd.Flags().GetString("environment")

	// Try to get from environment variables if not provided
	if username == "" {
		username = viper.GetString("wpengine.username")
		if username == "" {
			username = os.Getenv("WPE_USERNAME")
		}
	}

	if username == "" {
		return nil, fmt.Errorf("WP Engine username is required. Set via --username flag, config file, or WPE_USERNAME environment variable")
	}

	if environment == "" {
		environment = "production"
	}

	return &wpengine.Config{
		Username:    username,
		InstallName: installName,
		Environment: environment,
	}, nil
}

func ensureDDEVStarted(projectPath string) error {
	// Check if DDEV project exists
	configPath := filepath.Join(projectPath, ".ddev", "config.yaml")
	if _, err := os.Stat(configPath); err != nil {
		// Check if WordPress files exist to provide helpful guidance
		wpConfigPath := filepath.Join(projectPath, "wp-config.php")
		if _, wpErr := os.Stat(wpConfigPath); wpErr == nil {
			return fmt.Errorf(`DDEV project not found, but WordPress files detected in %s

It looks like you already have WordPress files but no DDEV configuration.
To fix this, run:

  cd %s
  stax init your-project-name --php-version=8.2
  
Then try the WP Engine sync again.`, projectPath, projectPath)
		}
		return fmt.Errorf("no DDEV project found in %s. Please run 'stax init <project-name>' first", projectPath)
	}

	// Try to start DDEV (this is idempotent - won't fail if already running)
	return ddev.Start(projectPath)
}

func init() {
	rootCmd.AddCommand(wpengineCmd)

	wpengineCmd.AddCommand(wpeListCmd)
	wpengineCmd.AddCommand(wpeConnectCmd)
	wpengineCmd.AddCommand(wpeSyncCmd)
	wpengineCmd.AddCommand(wpeDbCmd)

	wpeDbCmd.AddCommand(wpeDbDownloadCmd)
	wpeDbCmd.AddCommand(wpeDbImportCmd)
	wpeDbCmd.AddCommand(wpeDbAnalyzeCmd)
	wpeDbCmd.AddCommand(wpeDbRewriteCmd)
	wpeDbCmd.AddCommand(wpeDbDiagnoseCmd)

	// Global flags for WP Engine commands
	wpengineCmd.PersistentFlags().StringP("path", "p", "", "path to project (default: current directory)")
	wpengineCmd.PersistentFlags().String("install", "", "WP Engine install name")
	wpengineCmd.PersistentFlags().String("username", "", "WP Engine username for API access (or set WPE_USERNAME env var)")
	wpengineCmd.PersistentFlags().String("api-key", "", "WP Engine API key (or set WPE_API_KEY env var)")
	wpengineCmd.PersistentFlags().String("environment", "production", "WP Engine environment (production, staging, development)")

	// Sync command flags
	wpeSyncCmd.Flags().Bool("skip-media", true, "skip syncing media files (default: true)")
	wpeSyncCmd.Flags().Bool("skip-files", false, "skip syncing WordPress files")
	wpeSyncCmd.Flags().Bool("skip-database", false, "skip syncing database")
	wpeSyncCmd.Flags().Bool("delete-local", false, "delete local files not present on remote (WARNING: dangerous for development files)")
	wpeSyncCmd.Flags().Bool("suppress-debug", false, "suppress WordPress debug notices and warnings for cleaner output")
	wpeSyncCmd.Flags().Bool("create-upload-redirect", false, "create must-use plugin to redirect upload URLs to remote WP Engine domain")

	// Database download flags
	wpeDbDownloadCmd.Flags().String("output", "", "output file path")

	// Database import flags
	wpeDbImportCmd.Flags().Bool("suppress-debug", false, "suppress WordPress debug notices and warnings for cleaner output")
	wpeDbImportCmd.Flags().Bool("create-upload-redirect", false, "create must-use plugin to redirect upload URLs to remote WP Engine domain")
	
	// Database rewrite flags  
	wpeDbRewriteCmd.Flags().Bool("suppress-debug", false, "suppress WordPress debug notices and warnings for cleaner output")
	wpeDbRewriteCmd.Flags().Bool("create-upload-redirect", false, "create must-use plugin to redirect upload URLs to remote WP Engine domain")
}

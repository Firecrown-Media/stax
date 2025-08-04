package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Firecrown-Media/stax/pkg/config"
	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/Firecrown-Media/stax/pkg/wordpress"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [project-name]",
	Short: "Quick setup of a complete WordPress development environment",
	Long: `Quick setup that initializes DDEV, downloads WordPress, creates configuration,
and optionally installs WordPress with sensible defaults. This is the fastest way
to get a WordPress development environment up and running.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := "wordpress-site"
		if len(args) > 0 {
			projectName = args[0]
		}

		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		projectPath, err := filepath.Abs(projectPath)
		if err != nil {
			return fmt.Errorf("failed to resolve project path: %w", err)
		}

		// Load or create config
		cfg, err := config.Load(projectPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cfg.Name = projectName

		// Override with command line flags
		if phpVersion, _ := cmd.Flags().GetString("php-version"); phpVersion != "" {
			cfg.PHPVersion = phpVersion
		}
		if webServer, _ := cmd.Flags().GetString("webserver"); webServer != "" {
			cfg.WebServer = webServer
		}
		if database, _ := cmd.Flags().GetString("database"); database != "" {
			cfg.Database = database
		}

		fmt.Printf("🚀 Setting up WordPress project '%s' in %s\n", projectName, projectPath)

		// Step 1: Initialize DDEV
		fmt.Printf("📦 Initializing DDEV...\n")
		ddevConfig := ddev.Config{
			ProjectName:  cfg.Name,
			ProjectType:  cfg.Type,
			PHPVersion:   cfg.PHPVersion,
			WebServer:    cfg.WebServer,
			DatabaseType: cfg.Database,
		}

		if !ddev.IsProject(projectPath) {
			if err := ddev.Init(projectPath, ddevConfig); err != nil {
				return fmt.Errorf("failed to initialize DDEV: %w", err)
			}
		}

		// Step 2: Start DDEV
		fmt.Printf("🔄 Starting DDEV...\n")
		if err := ddev.Start(projectPath); err != nil {
			return fmt.Errorf("failed to start DDEV: %w", err)
		}

		// Step 3: Download WordPress if not present
		if !wordpress.HasWordPress(projectPath) {
			fmt.Printf("⬇️  Downloading WordPress core...\n")
			if err := wordpress.DownloadCore(projectPath); err != nil {
				return fmt.Errorf("failed to download WordPress: %w", err)
			}

			fmt.Printf("⚙️  Creating WordPress configuration...\n")
			if err := wordpress.CreateConfig(projectPath); err != nil {
				return fmt.Errorf("failed to create WordPress config: %w", err)
			}
		}

		// Step 4: Install WordPress if requested
		installWP, _ := cmd.Flags().GetBool("install-wp")
		if installWP {
			fmt.Printf("🔧 Installing WordPress...\n")
			wpConfig := wordpress.InstallConfig{
				URL:      cfg.WordPress.URL,
				Title:    cfg.WordPress.Title,
				Username: cfg.WordPress.AdminUser,
				Password: "admin", // Default password
				Email:    cfg.WordPress.AdminEmail,
			}

			if err := wordpress.Install(projectPath, wpConfig); err != nil {
				return fmt.Errorf("failed to install WordPress: %w", err)
			}
		}

		// Step 5: Save configuration
		if err := cfg.Save(projectPath); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ WordPress development environment setup complete!\n")
		fmt.Printf("\n📍 Project: %s\n", cfg.Name)
		fmt.Printf("🌐 URL: %s\n", cfg.WordPress.URL)

		if installWP {
			fmt.Printf("👤 Admin: %s / admin\n", cfg.WordPress.AdminUser)
			fmt.Printf("📧 Email: %s\n", cfg.WordPress.AdminEmail)
		}

		fmt.Printf("\n🛠️  Next steps:\n")
		if !installWP {
			fmt.Printf("   • Run 'stax wp install' to install WordPress\n")
		}
		fmt.Printf("   • Run 'stax status' to see project status\n")
		fmt.Printf("   • Run 'stax stop' to stop the environment\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().StringP("path", "p", "", "path to setup project (default: current directory)")
	setupCmd.Flags().String("php-version", "8.2", "PHP version to use")
	setupCmd.Flags().String("webserver", "nginx-fpm", "web server type")
	setupCmd.Flags().String("database", "mysql:8.0", "database type and version")
	setupCmd.Flags().Bool("install-wp", false, "install WordPress after setup")
}

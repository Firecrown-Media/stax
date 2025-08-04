package cmd

import (
	"fmt"
	"os"

	"github.com/Firecrown-Media/stax/pkg/wordpress"
	"github.com/spf13/cobra"
)

var wordpressCmd = &cobra.Command{
	Use:   "wp",
	Short: "WordPress-specific commands",
	Long:  `Commands for managing WordPress installations within DDEV.`,
}

var wpDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download WordPress core",
	Long:  `Download the latest WordPress core files using wp-cli.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		fmt.Printf("Downloading WordPress core...\n")

		if err := wordpress.DownloadCore(projectPath); err != nil {
			return fmt.Errorf("failed to download WordPress core: %w", err)
		}

		fmt.Printf("✅ WordPress core downloaded successfully!\n")
		return nil
	},
}

var wpInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install WordPress",
	Long:  `Install WordPress with initial configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		url, _ := cmd.Flags().GetString("url")
		title, _ := cmd.Flags().GetString("title")
		username, _ := cmd.Flags().GetString("admin-user")
		password, _ := cmd.Flags().GetString("admin-password")
		email, _ := cmd.Flags().GetString("admin-email")

		config := wordpress.InstallConfig{
			URL:      url,
			Title:    title,
			Username: username,
			Password: password,
			Email:    email,
		}

		fmt.Printf("Installing WordPress...\n")

		if err := wordpress.Install(projectPath, config); err != nil {
			return fmt.Errorf("failed to install WordPress: %w", err)
		}

		fmt.Printf("✅ WordPress installed successfully!\n")
		fmt.Printf("Admin URL: %s/wp-admin\n", url)
		fmt.Printf("Username: %s\n", username)
		return nil
	},
}

var wpConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Create WordPress configuration",
	Long:  `Create wp-config.php with DDEV database settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		fmt.Printf("Creating WordPress configuration...\n")

		if err := wordpress.CreateConfig(projectPath); err != nil {
			return fmt.Errorf("failed to create WordPress config: %w", err)
		}

		fmt.Printf("✅ WordPress configuration created successfully!\n")
		return nil
	},
}

var wpPluginCmd = &cobra.Command{
	Use:   "plugin [install|list] [plugin-name]",
	Short: "Manage WordPress plugins",
	Long:  `Install and manage WordPress plugins.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		action := args[0]

		switch action {
		case "install":
			if len(args) < 2 {
				return fmt.Errorf("plugin name required for install command")
			}

			plugin := args[1]
			activate, _ := cmd.Flags().GetBool("activate")

			fmt.Printf("Installing plugin: %s\n", plugin)

			if err := wordpress.InstallPlugin(projectPath, plugin, activate); err != nil {
				return fmt.Errorf("failed to install plugin: %w", err)
			}

			fmt.Printf("✅ Plugin installed successfully!\n")
		default:
			return fmt.Errorf("unsupported action: %s", action)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(wordpressCmd)

	wordpressCmd.AddCommand(wpDownloadCmd)
	wordpressCmd.AddCommand(wpInstallCmd)
	wordpressCmd.AddCommand(wpConfigCmd)
	wordpressCmd.AddCommand(wpPluginCmd)

	// Global flags
	wordpressCmd.PersistentFlags().StringP("path", "p", "", "path to project (default: current directory)")

	// Install command flags
	wpInstallCmd.Flags().String("url", "https://localhost", "WordPress site URL")
	wpInstallCmd.Flags().String("title", "My WordPress Site", "WordPress site title")
	wpInstallCmd.Flags().String("admin-user", "admin", "WordPress admin username")
	wpInstallCmd.Flags().String("admin-password", "admin", "WordPress admin password")
	wpInstallCmd.Flags().String("admin-email", "admin@localhost.local", "WordPress admin email")

	// Plugin command flags
	wpPluginCmd.Flags().Bool("activate", false, "activate plugin after installation")
}

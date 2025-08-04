package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new WordPress development environment",
	Long: `Initialize a new WordPress development environment using DDEV.
This will create a DDEV configuration optimized for WordPress development.`,
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

		phpVersion, _ := cmd.Flags().GetString("php-version")
		webServer, _ := cmd.Flags().GetString("webserver")
		database, _ := cmd.Flags().GetString("database")

		config := ddev.Config{
			ProjectName:  projectName,
			ProjectType:  "wordpress",
			PHPVersion:   phpVersion,
			WebServer:    webServer,
			DatabaseType: database,
		}

		fmt.Printf("Initializing WordPress project '%s' in %s\n", projectName, projectPath)

		if err := ddev.Init(projectPath, config); err != nil {
			return fmt.Errorf("failed to initialize DDEV project: %w", err)
		}

		fmt.Printf("✅ DDEV project initialized successfully!\n")
		fmt.Printf("🚀 Starting DDEV environment...\n")

		// Automatically start DDEV after initialization
		if err := ddev.Start(projectPath); err != nil {
			fmt.Printf("⚠️  DDEV project initialized but failed to start automatically: %v\n", err)
			fmt.Printf("You can start it manually with 'stax start'\n")
			return nil // Don't fail completely, project was initialized
		}

		fmt.Printf("✅ Development environment ready!\n")
		fmt.Printf("🌐 Your project will be available at: https://%s.ddev.site\n", projectName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("path", "p", "", "path to initialize project (default: current directory)")
	initCmd.Flags().String("php-version", "8.2", "PHP version to use")
	initCmd.Flags().String("webserver", "nginx-fpm", "web server type (nginx-fpm, apache-fpm)")
	initCmd.Flags().String("database", "mysql:8.0", "database type and version")
}

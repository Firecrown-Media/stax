package cmd

import (
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	initName            string
	initType            string
	initMode            string
	initPHPVersion      string
	initMySQLVersion    string
	initRepo            string
	initBranch          string
	initWPEngineInstall string
	initWPEngineEnv     string
	initInteractive     bool
	initSkipDB          bool
	initSkipBuild       bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "[construction] Initialize a new Stax project",
	Long: `Initialize a new Stax project in the current directory.

This command will:
  - Create a .stax.yml configuration file
  - Validate WPEngine credentials
  - Clone the GitHub repository (if specified)
  - Generate DDEV configuration
  - Start DDEV containers
  - Run the build process
  - Pull and import the database from WPEngine
  - Run search-replace operations

By default, this command runs in interactive mode, prompting for all
required information. You can skip prompts by providing all flags.`,
	Example: `  # Interactive mode (default)
  stax init

  # Non-interactive with all flags
  stax init \
    --name=firecrown-multisite \
    --mode=subdomain \
    --php-version=8.1 \
    --mysql-version=8.0 \
    --repo=https://github.com/Firecrown-Media/firecrown-multisite.git \
    --wpengine-install=fsmultisite \
    --no-interactive

  # Initialize without database import
  stax init --skip-db

  # Initialize without build process
  stax init --skip-build`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags
	initCmd.Flags().StringVar(&initName, "name", "", "project name (default: current directory name)")
	initCmd.Flags().StringVar(&initType, "type", "wordpress", "project type (wordpress or wordpress-multisite)")
	initCmd.Flags().StringVar(&initMode, "mode", "subdomain", "multisite mode (subdomain/subdirectory)")
	initCmd.Flags().StringVar(&initPHPVersion, "php-version", "8.1", "PHP version")
	initCmd.Flags().StringVar(&initMySQLVersion, "mysql-version", "8.0", "MySQL version")
	initCmd.Flags().StringVar(&initRepo, "repo", "", "GitHub repository URL")
	initCmd.Flags().StringVar(&initBranch, "branch", "main", "repository branch")
	initCmd.Flags().StringVar(&initWPEngineInstall, "wpengine-install", "", "WPEngine install name")
	initCmd.Flags().StringVar(&initWPEngineEnv, "wpengine-env", "production", "WPEngine environment")
	initCmd.Flags().BoolVar(&initInteractive, "interactive", true, "enable interactive prompts")
	initCmd.Flags().BoolVar(&initSkipDB, "skip-db", false, "skip database import")
	initCmd.Flags().BoolVar(&initSkipBuild, "skip-build", false, "skip build process")
}

func runInit(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Initializing Stax Project")

	ui.Warning("[construction] This command is not yet fully implemented")
	ui.Info("")
	ui.Info("To set up your project manually:")
	ui.Info("1. Clone your repository to the current directory")
	ui.Info("2. Run 'ddev config --project-type=wordpress'")
	ui.Info("3. Run 'ddev start'")
	ui.Info("4. Run 'stax db pull' to sync your database")
	ui.Info("")
	ui.Info("The interactive init workflow is coming soon!")

	return nil
}

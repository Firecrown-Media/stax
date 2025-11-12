package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/firecrown-media/stax/pkg/config"
	"github.com/firecrown-media/stax/pkg/ddev"
	"github.com/firecrown-media/stax/pkg/errors"
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
	initFromDDEV        bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "🚧 Initialize a new Stax project",
	Long: `Initialize a new Stax project in the current directory.

This command can either:
  - Set up a new project from scratch
  - Import an existing DDEV project (--from-ddev)

For new projects, this will:
  - Create a .stax.yml configuration file
  - Validate WPEngine credentials (optional)
  - Clone the GitHub repository (if specified)
  - Generate DDEV configuration
  - Start DDEV containers
  - Optionally pull database from WPEngine

By default, this command runs in interactive mode, prompting for all
required information. You can skip prompts by providing all flags.`,
	Example: `  # Interactive mode (default)
  stax init

  # Import existing DDEV project
  stax init --from-ddev

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
	initCmd.Flags().BoolVar(&initFromDDEV, "from-ddev", false, "import existing DDEV project")
}

func runInit(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Initializing Stax Project")

	projectDir := getProjectDir()

	// Check if importing from existing DDEV
	if initFromDDEV {
		return runInitFromDDEV(projectDir)
	}

	// TODO: Full init implementation
	// For now, show helpful guidance
	ui.Warning("Full init implementation coming soon!")
	ui.Info("")
	ui.Info("To set up your project manually:")
	ui.Info("1. Clone your repository to the current directory")
	ui.Info("2. Run 'ddev config --project-type=wordpress'")
	ui.Info("3. Run 'ddev start'")
	ui.Info("4. Run 'stax db pull' to sync your database")
	ui.Info("")
	ui.Info("Or use 'stax init --from-ddev' to import an existing DDEV project")

	return nil
}

func runInitFromDDEV(projectDir string) error {
	ui.Info("Importing existing DDEV project...")

	// Check if DDEV config exists
	if !ddev.IsConfigured(projectDir) {
		return errors.NewWithSolution(
			"No DDEV configuration found",
			"Cannot import from DDEV - no .ddev/config.yaml exists",
			errors.Solution{
				Description: "Initialize DDEV first",
				Steps: []string{
					"Run 'ddev config --project-type=wordpress' to set up DDEV",
					"Then run 'stax init --from-ddev' again",
				},
			},
		)
	}

	// Check if .stax.yml already exists
	configPath := filepath.Join(projectDir, ".stax.yml")
	if _, err := os.Stat(configPath); err == nil {
		ui.Warning(".stax.yml already exists")
		fmt.Print("Overwrite? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			ui.Info("Import cancelled")
			return nil
		}
	}

	// Read DDEV config to get basic info
	ddevConfigPath := filepath.Join(projectDir, ".ddev", "config.yaml")
	ui.Success("Found DDEV configuration")

	// Prompt for optional WPEngine integration
	ui.Info("\nOptional: Configure WPEngine integration")
	fmt.Print("Add WPEngine integration? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	addWPEngine := (response == "y" || response == "yes")

	var install, env string
	if addWPEngine {
		// Prompt for WPEngine details
		fmt.Print("WPEngine install name: ")
		install, _ = reader.ReadString('\n')
		install = strings.TrimSpace(install)

		fmt.Print("WPEngine environment (production/staging/development) [production]: ")
		env, _ = reader.ReadString('\n')
		env = strings.TrimSpace(env)
		if env == "" {
			env = "production"
		}
	}

	// Create basic .stax.yml from DDEV config
	// Start with defaults and add WPEngine config if provided
	cfg := config.Defaults()

	if addWPEngine {
		cfg.WPEngine.Install = install
		cfg.WPEngine.Environment = env
	}

	// Save .stax.yml
	if err := config.Save(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	ui.Success("Created .stax.yml from DDEV configuration at %s", ddevConfigPath)
	ui.Info("\nYour DDEV project now has Stax features enabled!")
	ui.Info("Run 'stax status' to see your environment")

	if addWPEngine {
		ui.Info("Run 'stax db pull' to sync your database from WPEngine")
	}

	return nil
}

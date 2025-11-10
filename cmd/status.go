package cmd

import (
	"github.com/firecrown-media/stax/pkg/ddev"
	"github.com/firecrown-media/stax/pkg/errors"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	statusJSON bool
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show environment status",
	Long: `Show detailed status information about the DDEV environment,
including container health, URLs, configuration, database info, and
WPEngine sync status.`,
	Aliases: []string{"s"},
	Example: `  # Show status
  stax status

  # Show status as JSON
  stax status --json`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "output as JSON")
}

func runStatus(cmd *cobra.Command, args []string) error {
	projectDir := getProjectDir()

	// Check if we have .stax.yml config
	hasStaxConfig := false
	if cfg != nil {
		hasStaxConfig = true
	}

	// Check if we have DDEV config
	hasDDEVConfig := ddev.IsConfigured(projectDir)

	if !hasStaxConfig && !hasDDEVConfig {
		return errors.NewWithSolution(
			"No project configuration found",
			"Neither .stax.yml nor .ddev/config.yaml exists",
			errors.Solution{
				Description: "Initialize your project",
				Steps: []string{
					"Run 'stax init' to set up a new Stax project",
					"Or run 'ddev config' if you just want basic DDEV",
				},
			},
		)
	}

	if !hasStaxConfig {
		ui.Warning("Using DDEV configuration only (no .stax.yml found)")
		ui.Info("Run 'stax init' to enable Stax features like WPEngine sync")
	}

	ui.PrintHeader("Environment Status")

	// TODO: Get DDEV status
	// TODO: Get container health
	// TODO: Get URLs from config
	// TODO: Get database info
	// TODO: Get WPEngine sync info
	// TODO: Format as JSON if requested

	ui.Info("Environment status is not yet implemented")
	ui.Info("This is a placeholder for DDEV integration")

	ui.Section("Environment: Not Running")

	return nil
}

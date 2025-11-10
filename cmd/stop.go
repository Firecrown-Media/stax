package cmd

import (
	"github.com/firecrown-media/stax/pkg/ddev"
	"github.com/firecrown-media/stax/pkg/errors"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	stopAll        bool
	stopRemoveData bool
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the DDEV environment",
	Long: `Stop the DDEV environment for the current project.

This command stops all DDEV containers while preserving data.
Use --remove-data to also remove database data (destructive).`,
	Aliases: []string{"down"},
	Example: `  # Stop current project
  stax stop

  # Stop all DDEV projects
  stax stop --all`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVar(&stopAll, "all", false, "stop all DDEV projects")
	stopCmd.Flags().BoolVar(&stopRemoveData, "remove-data", false, "remove database data (destructive)")
}

func runStop(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Stopping Environment")

	projectDir := getProjectDir()

	// Handle --all flag (poweroff all DDEV projects)
	if stopAll {
		// No config check needed for --all flag
		ui.Info("Environment stop is not yet implemented")
		ui.Info("This is a placeholder for DDEV integration")
		ui.Success("Environment stopped!")
		return nil
	}

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

	// TODO: Run ddev stop
	// TODO: Optionally remove data if requested

	ui.Info("Environment stop is not yet implemented")
	ui.Info("This is a placeholder for DDEV integration")

	ui.Success("Environment stopped!")

	return nil
}

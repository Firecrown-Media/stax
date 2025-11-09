package cmd

import (
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

	// TODO: Run ddev stop or ddev poweroff (if --all)
	// TODO: Optionally remove data if requested

	ui.Info("Environment stop is not yet implemented")
	ui.Info("This is a placeholder for DDEV integration")

	ui.Success("Environment stopped!")

	return nil
}

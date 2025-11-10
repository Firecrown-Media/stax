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
	Short: "[warning] Stop the DDEV environment",
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

	ui.Warning("[warning] Using DDEV directly (stax wrapper coming soon)")
	ui.Info("")
	ui.Info("For now, use DDEV commands directly:")

	if stopAll {
		ui.Info("  ddev poweroff")
	} else {
		ui.Info("  ddev stop")
	}

	if stopRemoveData {
		ui.Warning("  ddev delete --omit-snapshot")
		ui.Warning("  WARNING: This will permanently delete your database!")
	}

	ui.Info("")
	ui.Info("Future stax stop will integrate these steps automatically.")

	return nil
}

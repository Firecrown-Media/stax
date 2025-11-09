package cmd

import (
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	restartBuild  bool
	restartXdebug bool
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the DDEV environment",
	Long:  `Restart the DDEV environment (stop followed by start).`,
	Example: `  # Basic restart
  stax restart

  # Restart with build
  stax restart --build`,
	RunE: runRestart,
}

func init() {
	rootCmd.AddCommand(restartCmd)

	restartCmd.Flags().BoolVar(&restartBuild, "build", false, "run build process after restart")
	restartCmd.Flags().BoolVar(&restartXdebug, "xdebug", false, "enable Xdebug")
}

func runRestart(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Restarting Environment")

	// TODO: Run ddev restart
	// TODO: Enable Xdebug if requested
	// TODO: Run build process if requested

	ui.Info("Environment restart is not yet implemented")

	ui.Success("Environment restarted!")

	return nil
}

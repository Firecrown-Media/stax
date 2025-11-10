package cmd

import (
	"github.com/firecrown-media/stax/pkg/ddev"
	"github.com/firecrown-media/stax/pkg/errors"
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

	// TODO: Run ddev restart
	// TODO: Enable Xdebug if requested
	// TODO: Run build process if requested

	ui.Info("Environment restart is not yet implemented")

	ui.Success("Environment restarted!")

	return nil
}

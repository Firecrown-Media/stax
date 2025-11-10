package cmd

import (
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	startBuild  bool
	startXdebug bool
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "[warning] Start the DDEV environment",
	Long: `Start the DDEV environment for the current project.

This command starts all DDEV containers (web, database, router) and
optionally runs the build process and enables Xdebug.`,
	Aliases: []string{"up"},
	Example: `  # Basic start
  stax start

  # Start with Xdebug enabled
  stax start --xdebug

  # Start and rebuild
  stax start --build`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVar(&startBuild, "build", false, "run build process after start")
	startCmd.Flags().BoolVar(&startXdebug, "xdebug", false, "enable Xdebug")
}

func runStart(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Starting Environment")

	ui.Warning("[warning] Using DDEV directly (stax wrapper coming soon)")
	ui.Info("")
	ui.Info("For now, use DDEV commands directly:")
	ui.Info("  ddev start")

	if startXdebug {
		ui.Info("  ddev xdebug on")
	}

	if startBuild {
		ui.Info("  npm run build")
	}

	ui.Info("")
	ui.Info("Future stax start will integrate these steps automatically.")

	return nil
}

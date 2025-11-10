package cmd

import (
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	statusJSON bool
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "[warning] Show environment status",
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
	ui.PrintHeader("Environment Status")

	ui.Warning("[warning] Using DDEV directly (stax wrapper coming soon)")
	ui.Info("")
	ui.Info("For now, use DDEV commands directly:")
	ui.Info("  ddev describe")

	if statusJSON {
		ui.Info("  ddev describe --json")
	}

	ui.Info("")
	ui.Info("Future stax status will provide additional WordPress-specific info.")

	return nil
}

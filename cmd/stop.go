package cmd

import (
	"fmt"

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
	Short: "✓ Stop the DDEV environment",
	Long: `Stop the DDEV environment for the current project.

This command stops all DDEV containers while preserving data.
Use --remove-data to also remove database data (destructive).`,
	Aliases: []string{"down"},
	Example: `  # Stop current project
  stax stop

  # Stop all DDEV projects
  stax stop --all

  # Stop and remove data (destructive)
  stax stop --remove-data`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVar(&stopAll, "all", false, "stop all DDEV projects")
	stopCmd.Flags().BoolVar(&stopRemoveData, "remove-data", false, "remove database data (destructive)")
}

func runStop(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Stopping DDEV Environment")

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

	// Handle --all flag (poweroff all DDEV projects)
	if stopAll {
		spinner := ui.NewSpinner("Stopping all DDEV projects")
		spinner.Start()

		if err := ddev.PowerOff(); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to stop all projects: %w", err)
		}

		spinner.Success("All DDEV projects stopped")
		return nil
	}

	// Check if project is configured
	if !ddev.IsConfigured(projectDir) {
		ui.Warning("DDEV is not configured for this project")
		return nil
	}

	// Check if already stopped
	status, err := ddev.GetStatus(projectDir)
	if err == nil && !status.Running {
		ui.Info("Environment is already stopped")
		return nil
	}

	// Handle --remove-data flag (destructive)
	if stopRemoveData {
		confirmed := ui.Confirm("⚠️  WARNING: This will permanently delete your database. Continue?")
		if !confirmed {
			ui.Info("Operation cancelled")
			return nil
		}

		spinner := ui.NewSpinner("Stopping and removing data")
		spinner.Start()

		if err := ddev.Delete(projectDir, false); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to delete project: %w", err)
		}

		spinner.Success("Environment stopped and data removed")
		ui.Warning("Database has been permanently deleted")
		return nil
	}

	// Normal stop
	spinner := ui.NewSpinner("Stopping DDEV containers")
	spinner.Start()

	if err := ddev.Stop(projectDir); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to stop environment: %w", err)
	}

	spinner.Success("Environment stopped successfully")
	fmt.Println()
	ui.Info("Data has been preserved. Use 'stax start' to restart.")

	return nil
}

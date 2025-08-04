package cmd

import (
	"fmt"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
)

var poweroffCmd = &cobra.Command{
	Use:   "poweroff",
	Short: "Stop all DDEV projects and clean up",
	Long: `Stop all running DDEV projects and clean up Docker resources.
This is useful for completely shutting down the development environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🔌 Stopping all DDEV projects and cleaning up...\n")

		if err := ddev.Poweroff(); err != nil {
			return fmt.Errorf("failed to power off DDEV: %w", err)
		}

		fmt.Printf("✅ All DDEV projects stopped and cleaned up\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(poweroffCmd)
}
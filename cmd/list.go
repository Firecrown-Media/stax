package cmd

import (
	"fmt"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DDEV projects",
	Long:  `List all DDEV projects and their current status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ddev.List(); err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
package cmd

import (
	"fmt"
	"os"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of DDEV projects",
	Long:  `Show the status of all DDEV projects or the current project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		if err := ddev.Status(projectPath); err != nil {
			return fmt.Errorf("failed to get DDEV status: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().StringP("path", "p", "", "path to project (default: current directory)")
}

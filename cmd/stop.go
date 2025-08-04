package cmd

import (
	"fmt"
	"os"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the WordPress development environment",
	Long: `Stop the DDEV WordPress development environment.
This will stop all running containers and services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		fmt.Printf("Stopping WordPress development environment...\n")

		if err := ddev.Stop(projectPath); err != nil {
			return fmt.Errorf("failed to stop DDEV project: %w", err)
		}

		fmt.Printf("✅ Development environment stopped successfully!\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().StringP("path", "p", "", "path to project (default: current directory)")
}

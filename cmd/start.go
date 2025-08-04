package cmd

import (
	"fmt"
	"os"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the WordPress development environment",
	Long: `Start the DDEV WordPress development environment.
This will start all necessary containers and services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		fmt.Printf("Starting WordPress development environment...\n")

		if err := ddev.Start(projectPath); err != nil {
			return fmt.Errorf("failed to start DDEV project: %w", err)
		}

		fmt.Printf("✅ Development environment started successfully!\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP("path", "p", "", "path to project (default: current directory)")
}

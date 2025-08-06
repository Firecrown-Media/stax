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

		// Show general DDEV project list
		fmt.Printf("📋 DDEV Projects Overview:\n")
		fmt.Printf("========================\n")
		if err := ddev.Status(projectPath); err != nil {
			return fmt.Errorf("failed to get DDEV status: %w", err)
		}

		// If we're in a DDEV project directory, also show detailed info
		if ddev.IsProject(projectPath) {
			fmt.Printf("\n🔍 Current Project Details:\n")
			fmt.Printf("==========================\n")
			if err := ddev.Describe(projectPath); err != nil {
				fmt.Printf("⚠️  Failed to get project details: %v\n", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().StringP("path", "p", "", "path to project (default: current directory)")
}

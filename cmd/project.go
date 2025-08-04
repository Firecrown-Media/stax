package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  `Commands for managing projects.`,
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all DDEV projects and their current status. This is an alias for 'stax list'.`,
	RunE:  listCmd.RunE, // Use the same function as listCmd
}

var projectDeployCmd = &cobra.Command{
	Use:   "deploy [project-name]",
	Short: "Deploy a project to a remote environment",
	Long:  `Deploy a project to a remote environment like staging or production. (Not yet implemented)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Sorry, the 'project deploy' command is not yet implemented.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectDeployCmd)

	projectDeployCmd.Flags().String("environment", "staging", "The environment to deploy to (e.g., staging, production)")
}

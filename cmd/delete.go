package cmd

import (
	"fmt"
	"os"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [project-name]",
	Short: "Delete a DDEV project",
	Long: `Delete a DDEV project, optionally removing all data including database.
Use --omit-snapshot to remove all data including database and files.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		omitSnapshot, _ := cmd.Flags().GetBool("omit-snapshot")
		yes, _ := cmd.Flags().GetBool("yes")

		var projectName string
		if len(args) > 0 {
			projectName = args[0]
		}

		if err := ddev.Delete(projectPath, projectName, omitSnapshot, yes); err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}

		if omitSnapshot {
			fmt.Printf("✅ Project deleted completely (including database)\n")
		} else {
			fmt.Printf("✅ Project deleted (database snapshot preserved)\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("path", "p", "", "path to project (default: current directory)")
	deleteCmd.Flags().BoolP("omit-snapshot", "O", false, "omit database snapshot (remove all data)")
	deleteCmd.Flags().BoolP("yes", "y", false, "skip confirmation prompts")
}
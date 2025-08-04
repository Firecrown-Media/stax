package cmd

import (
	"fmt"
	"os"

	"github.com/Firecrown-Media/stax/pkg/ddev"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database operations",
	Long:  `Perform database operations like export, import, and management.`,
}

var exportDbCmd = &cobra.Command{
	Use:   "export-db",
	Short: "Export database to file",
	Long:  `Export the current project's database to a SQL file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		filename, _ := cmd.Flags().GetString("file")
		if filename == "" {
			return fmt.Errorf("--file parameter is required")
		}

		fmt.Printf("📤 Exporting database to %s...\n", filename)

		if err := ddev.ExportDB(projectPath, filename); err != nil {
			return fmt.Errorf("failed to export database: %w", err)
		}

		fmt.Printf("✅ Database exported successfully\n")
		return nil
	},
}

var importDbCmd = &cobra.Command{
	Use:   "import-db",
	Short: "Import database from file",
	Long:  `Import a SQL file into the current project's database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("path")
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}

		src, _ := cmd.Flags().GetString("src")
		if src == "" {
			return fmt.Errorf("--src parameter is required")
		}

		fmt.Printf("📥 Importing database from %s...\n", src)

		if err := ddev.ImportDB(projectPath, src); err != nil {
			return fmt.Errorf("failed to import database: %w", err)
		}

		fmt.Printf("✅ Database imported successfully\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportDbCmd)
	rootCmd.AddCommand(importDbCmd)

	exportDbCmd.Flags().StringP("path", "p", "", "path to project (default: current directory)")
	exportDbCmd.Flags().String("file", "", "output file path (required)")

	importDbCmd.Flags().StringP("path", "p", "", "path to project (default: current directory)")
	importDbCmd.Flags().String("src", "", "source SQL file path (required)")
}
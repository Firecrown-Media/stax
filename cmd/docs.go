package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for stax",
	Long:  `Generate documentation for stax in various formats.`,
}

var manCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man pages",
	Long:  `Generate man pages for all stax commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			return fmt.Errorf("--dir flag is required")
		}

		// Create the directory if it doesn't exist
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Define the header for the man page
		header := &doc.GenManHeader{
			Title:   "STAX",
			Section: "1",
			Source:  "stax",
			Manual:  "stax Manual",
		}

		fmt.Printf("Generating man pages in %s...\n", dir)

		// Generate the man pages
		if err := doc.GenManTree(rootCmd, header, dir); err != nil {
			return fmt.Errorf("failed to generate man pages: %w", err)
		}

		fmt.Println("✅ Man pages generated successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.AddCommand(manCmd)

	manCmd.Flags().StringP("dir", "d", "./man", "directory to write man pages to")
}

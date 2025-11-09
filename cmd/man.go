package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var manCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man page for stax",
	Long: `Generate a man page for stax in groff format.

The man page can be installed to your system's man directory
or viewed directly with: man ./stax.1`,
	Example: `  # Generate man page to current directory
  stax man

  # Generate to specific location
  stax man -o /usr/local/share/man/man1/

  # View the generated man page
  man ./stax.1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir, _ := cmd.Flags().GetString("output")

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Generate man page
		header := &doc.GenManHeader{
			Title:   "STAX",
			Section: "1",
			Source:  "Stax " + Version,
			Manual:  "Stax Manual",
		}

		if err := doc.GenManTree(rootCmd, header, outputDir); err != nil {
			return fmt.Errorf("failed to generate man page: %w", err)
		}

		manFile := filepath.Join(outputDir, "stax.1")
		fmt.Printf("Man page generated: %s\n", manFile)
		fmt.Printf("\nView with: man %s\n", manFile)
		fmt.Printf("Install with: sudo cp %s /usr/local/share/man/man1/\n", manFile)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(manCmd)
	manCmd.Flags().StringP("output", "o", ".", "output directory for man page")
}

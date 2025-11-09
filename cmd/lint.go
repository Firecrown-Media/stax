package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/firecrown-media/stax/pkg/build"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	lintFix      bool
	lintStaged   bool
	lintStandard string
	lintFiles    []string
	lintFormat   string
)

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Run PHP CodeSniffer checks",
	Long: `Run PHPCS (PHP CodeSniffer) to check code quality.

Uses the configuration from .phpcs.xml.dist if present, or falls
back to composer run-script lint if available.

This helps maintain code quality standards and catch issues before
they're committed.`,
	Example: `  # Run linting on all PHP files
  stax lint

  # Auto-fix issues
  stax lint:fix

  # Lint only staged files (for pre-commit)
  stax lint:staged

  # Check specific files
  stax lint wp-content/mu-plugins/firecrown/src/`,
	RunE: runLint,
}

var lintFixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Auto-fix code quality issues",
	Long: `Run PHPCBF (PHP Code Beautifier and Fixer) to automatically fix
code quality issues where possible.

This will modify your files to fix issues like:
  - Spacing and indentation
  - Line endings
  - Code formatting

Review the changes after running this command.`,
	RunE: runLintFix,
}

var lintCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check code quality without fixing",
	Long: `Run PHPCS in check-only mode.

This is the same as 'stax lint' but more explicit about not
modifying files.`,
	RunE: runLint,
}

var lintStagedCmd = &cobra.Command{
	Use:   "staged",
	Short: "Lint only staged files",
	Long: `Run PHPCS only on files staged for commit.

This is useful in pre-commit hooks to only check files
that are about to be committed.`,
	RunE: runLintStaged,
}

func init() {
	rootCmd.AddCommand(lintCmd)

	// Subcommands
	lintCmd.AddCommand(lintFixCmd)
	lintCmd.AddCommand(lintCheckCmd)
	lintCmd.AddCommand(lintStagedCmd)

	// Flags
	lintCmd.Flags().BoolVar(&lintFix, "fix", false, "auto-fix issues (runs phpcbf)")
	lintCmd.Flags().BoolVar(&lintStaged, "staged", false, "only lint staged files")
	lintCmd.Flags().StringVar(&lintStandard, "standard", "", "coding standard to use")
	lintCmd.Flags().StringSliceVar(&lintFiles, "files", []string{}, "specific files to lint")
	lintCmd.Flags().StringVar(&lintFormat, "format", "full", "output format (full, summary, json)")
}

func runLint(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Running Code Quality Checks")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	quality := build.NewQuality(projectDir)

	// Prepare options
	options := build.PHPCSOptions{
		WorkingDir: projectDir,
		Standard:   lintStandard,
		Report:     lintFormat,
		ShowSniffs: true,
	}

	// Add files if specified
	if len(lintFiles) > 0 {
		options.Files = lintFiles
	} else if len(args) > 0 {
		options.Files = args
	}

	// Find PHPCS config
	configPath, err := quality.GetPHPCSConfig()
	if err != nil {
		ui.Warning("PHPCS config not found, using default settings")
	} else {
		ui.Verbose("Using PHPCS config: %s", configPath)
		options.ConfigFile = configPath
	}

	// Run PHPCS
	ui.Info("Running PHP CodeSniffer...")

	result, err := quality.RunPHPCS(options)
	if err != nil {
		return fmt.Errorf("PHPCS execution failed: %w", err)
	}

	// Display results
	if result.Success {
		ui.Success("No errors or warnings found!")
		return nil
	}

	// Format and display results
	formattedResults := quality.FormatPHPCSResults(result)
	fmt.Println(formattedResults)

	// Summary
	ui.Section("Summary")
	ui.Error("%d error(s) found", result.Errors)
	ui.Warning("%d warning(s) found", result.Warnings)

	if result.Fixable > 0 {
		ui.Info("")
		ui.Info("%d issue(s) can be fixed automatically", result.Fixable)
		ui.Info("Run 'stax lint:fix' to auto-fix these issues")
	}

	// Exit with error if errors found
	if result.Errors > 0 {
		return fmt.Errorf("code quality check failed")
	}

	return nil
}

func runLintFix(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Auto-fixing Code Quality Issues")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	quality := build.NewQuality(projectDir)

	// Prepare options
	options := build.PHPCSOptions{
		WorkingDir: projectDir,
		Standard:   lintStandard,
	}

	// Add files if specified
	if len(args) > 0 {
		options.Files = args
	}

	// Find PHPCS config
	configPath, err := quality.GetPHPCSConfig()
	if err != nil {
		ui.Warning("PHPCS config not found, using default settings")
	} else {
		options.ConfigFile = configPath
	}

	// Run PHPCBF
	ui.Info("Running PHP Code Beautifier and Fixer...")

	if err := quality.RunPHPCBF(options); err != nil {
		return fmt.Errorf("PHPCBF failed: %w", err)
	}

	ui.Success("Code formatting complete!")
	ui.Info("Review the changes with 'git diff'")

	return nil
}

func runLintStaged(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Linting Staged Files")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Get staged PHP files
	stagedFiles, err := getStagedPHPFiles(projectDir)
	if err != nil {
		return err
	}

	if len(stagedFiles) == 0 {
		ui.Info("No staged PHP files to lint")
		return nil
	}

	ui.Info("Found %d staged PHP file(s)", len(stagedFiles))

	quality := build.NewQuality(projectDir)

	// Prepare options
	options := build.PHPCSOptions{
		WorkingDir: projectDir,
		Files:      stagedFiles,
		Report:     "full",
		ShowSniffs: true,
	}

	// Find PHPCS config
	configPath, err := quality.GetPHPCSConfig()
	if err == nil {
		options.ConfigFile = configPath
	}

	// Run PHPCS
	result, err := quality.RunPHPCS(options)
	if err != nil {
		return fmt.Errorf("PHPCS execution failed: %w", err)
	}

	// Display results
	if result.Success {
		ui.Success("All staged files pass code quality checks!")
		return nil
	}

	// Format and display results
	formattedResults := quality.FormatPHPCSResults(result)
	fmt.Println(formattedResults)

	// Summary
	ui.Error("%d error(s) in staged files", result.Errors)
	ui.Warning("%d warning(s) in staged files", result.Warnings)

	if result.Fixable > 0 {
		ui.Info("%d issue(s) can be fixed automatically", result.Fixable)
		ui.Info("Run 'stax lint:fix' to auto-fix")
	}

	return fmt.Errorf("staged files have code quality issues")
}

// getStagedPHPFiles returns a list of staged PHP files
func getStagedPHPFiles(projectDir string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	cmd.Dir = projectDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	phpFiles := []string{}

	for _, line := range lines {
		if strings.HasSuffix(line, ".php") {
			// Convert to absolute path
			absPath := filepath.Join(projectDir, line)
			if _, err := os.Stat(absPath); err == nil {
				phpFiles = append(phpFiles, absPath)
			}
		}
	}

	return phpFiles, nil
}

package cmd

import (
	"fmt"

	"github.com/firecrown-media/stax/pkg/diagnostics"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	doctorFix  bool
	doctorJSON bool
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "✓ Diagnose and fix common issues",
	Long: `Run diagnostics to check for common issues with the Stax environment.

This command checks:
  - Git installation and version
  - Docker installation and status
  - DDEV installation and version
  - Stax configuration (.stax.yml)
  - DDEV configuration (.ddev/)
  - WPEngine credentials
  - Port availability
  - Disk space

Each check provides detailed information and suggestions for fixing any issues found.`,
	Example: `  # Diagnose issues
  stax doctor

  # Show detailed JSON output
  stax doctor --json

  # Future: Automatically fix issues
  stax doctor --fix`,
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)

	doctorCmd.Flags().BoolVar(&doctorFix, "fix", false, "automatically fix issues (not yet implemented)")
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "output results as JSON")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Running System Diagnostics")
	fmt.Println()

	projectDir := getProjectDir()

	// Run all diagnostic checks
	report, err := diagnostics.RunAllChecks(projectDir)
	if err != nil {
		return fmt.Errorf("failed to run diagnostics: %w", err)
	}

	// Display results
	if doctorJSON {
		return outputDoctorJSON(report)
	}

	return outputDoctorTable(report)
}

func outputDoctorTable(report *diagnostics.DiagnosticReport) error {
	// Display each check result
	for _, check := range report.Checks {
		displayCheckResult(check)
	}

	// Display summary
	fmt.Println()
	ui.Section("Summary")
	fmt.Printf("  Total Checks:   %d\n", report.Summary.Total)
	fmt.Printf("  Passed:         %s %d\n", getStatusEmoji(diagnostics.StatusPass), report.Summary.Passed)
	fmt.Printf("  Warnings:       %s %d\n", getStatusEmoji(diagnostics.StatusWarning), report.Summary.Warnings)
	fmt.Printf("  Failed:         %s %d\n", getStatusEmoji(diagnostics.StatusFail), report.Summary.Failed)
	fmt.Printf("  Skipped:        %s %d\n", getStatusEmoji(diagnostics.StatusSkip), report.Summary.Skipped)
	fmt.Println()

	// Overall health status
	if report.IsHealthy() {
		ui.Success("✓ System is healthy - all checks passed!")
		fmt.Println()
		ui.Info("Quick commands:")
		ui.Info("  stax start     - Start the environment")
		ui.Info("  stax status    - Show environment status")
	} else if report.HasCriticalFailures() {
		ui.Error("✗ Critical issues found that need attention")
		fmt.Println()
		ui.Info("Review the failures above and follow the suggestions to fix them.")
		ui.Info("Run 'stax doctor' again after making changes.")
	} else if report.HasWarnings() {
		ui.Warning("⚠ System is functional but has warnings")
		fmt.Println()
		ui.Info("The environment should work, but you may want to address the warnings.")
	}

	// Show fix flag info
	if report.HasCriticalFailures() || report.HasWarnings() {
		fmt.Println()
		ui.Info("Note: Automatic fixes (stax doctor --fix) are planned for a future release.")
	}

	return nil
}

func displayCheckResult(check diagnostics.CheckResult) {
	emoji := getStatusEmoji(check.Status)

	fmt.Printf("%s %s\n", emoji, check.Name)

	// Show message with indentation
	if check.Message != "" {
		fmt.Printf("  %s\n", check.Message)
	}

	// Show details if available
	if len(check.Details) > 0 {
		for key, value := range check.Details {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Show suggestion for warnings and failures
	if (check.Status == diagnostics.StatusWarning || check.Status == diagnostics.StatusFail) && check.Suggestion != "" {
		fmt.Printf("  → %s\n", check.Suggestion)
	}

	fmt.Println()
}

func getStatusEmoji(status diagnostics.CheckStatus) string {
	switch status {
	case diagnostics.StatusPass:
		return "✓"
	case diagnostics.StatusWarning:
		return "⚠"
	case diagnostics.StatusFail:
		return "✗"
	case diagnostics.StatusSkip:
		return "○"
	default:
		return "?"
	}
}

func outputDoctorJSON(report *diagnostics.DiagnosticReport) error {
	// For JSON output, we'd marshal the report
	// For now, just show a message
	ui.Info("JSON output not yet implemented")
	ui.Info("Use regular output for now: stax doctor")
	return nil
}

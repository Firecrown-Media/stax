package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/firecrown-media/stax/pkg/diagnostics"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	doctorFix     bool
	doctorJSON    bool
	doctorVerbose bool
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose and fix common issues",
	Long: `Run comprehensive diagnostics to check for common issues with the Stax environment.

This command checks:
  System Requirements:
  - Git installation and version
  - Docker installation and status
  - DDEV installation and version
  - Go installation (optional, for development)

  Configuration:
  - Stax configuration (.stax.yml)
  - DDEV configuration (.ddev/)

  Credentials:
  - WPEngine API credentials
  - WPEngine SSH credentials
  - SSH key availability and permissions
  - GitHub token (optional)

  Network:
  - Port availability (80, 443, 3306, 8025, 8036)
  - WPEngine API connectivity
  - WPEngine SSH gateway connectivity
  - Internet connectivity

  Environment:
  - Disk space availability
  - DDEV project status
  - Database connectivity
  - WordPress installation

Each check provides detailed information and suggestions for fixing any issues found.`,
	Example: `  # Run diagnostics
  stax doctor

  # Show detailed output including skipped checks
  stax doctor --verbose

  # Show JSON output
  stax doctor --json

  # Future: Automatically fix issues
  stax doctor --fix`,
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)

	doctorCmd.Flags().BoolVar(&doctorFix, "fix", false, "automatically fix issues (not yet implemented)")
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "output results as JSON")
	doctorCmd.Flags().BoolVarP(&doctorVerbose, "verbose", "v", false, "show detailed output including skipped checks")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	if !doctorJSON {
		ui.PrintHeader("Running System Diagnostics")
		fmt.Println()
	}

	projectDir := getProjectDir()

	// Run all diagnostic checks
	report, err := diagnostics.RunAllChecks(projectDir, doctorVerbose)
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
		// Skip displaying skipped checks unless verbose mode
		if check.Status == diagnostics.StatusSkip && !report.Verbose {
			continue
		}
		displayCheckResult(check, report.Verbose)
	}

	// Display summary
	fmt.Println()
	ui.Section("Summary")
	fmt.Printf("  Total Checks:   %d\n", report.Summary.Total)
	fmt.Printf("  Passed:         %s %d\n", getStatusEmoji(diagnostics.StatusPass), report.Summary.Passed)
	fmt.Printf("  Warnings:       %s %d\n", getStatusEmoji(diagnostics.StatusWarning), report.Summary.Warnings)
	fmt.Printf("  Failed:         %s %d\n", getStatusEmoji(diagnostics.StatusFail), report.Summary.Failed)
	if report.Verbose {
		fmt.Printf("  Skipped:        %s %d\n", getStatusEmoji(diagnostics.StatusSkip), report.Summary.Skipped)
	}
	fmt.Println()

	// Calculate health score
	healthScore := calculateHealthScore(report)
	fmt.Printf("  Health Score:   %d%%\n", healthScore)
	fmt.Println()

	// Overall health status
	if report.IsHealthy() {
		ui.Success("System is healthy - all checks passed!")
		fmt.Println()
		ui.Info("Quick commands:")
		ui.Info("  stax start     - Start the environment")
		ui.Info("  stax status    - Show environment status")
	} else if report.HasCriticalFailures() {
		ui.Error("Critical issues found that need attention")
		fmt.Println()
		ui.Info("Review the failures above and follow the suggestions to fix them.")
		ui.Info("Run 'stax doctor' again after making changes.")
	} else if report.HasWarnings() {
		ui.Warning("System is functional but has warnings")
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

func displayCheckResult(check diagnostics.CheckResult, verbose bool) {
	emoji := getStatusEmoji(check.Status)

	fmt.Printf("%s %s\n", emoji, check.Name)

	// Show message with indentation
	if check.Message != "" {
		fmt.Printf("  %s\n", check.Message)
	}

	// Show details if available and in verbose mode
	if verbose && len(check.Details) > 0 {
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

func calculateHealthScore(report *diagnostics.DiagnosticReport) int {
	// Skip skipped checks from the calculation
	totalRelevant := report.Summary.Passed + report.Summary.Warnings + report.Summary.Failed
	if totalRelevant == 0 {
		return 100
	}

	// Passed checks get full points, warnings get half points, failures get zero
	score := float64(report.Summary.Passed) + (float64(report.Summary.Warnings) * 0.5)
	percentage := int((score / float64(totalRelevant)) * 100)

	return percentage
}

func outputDoctorJSON(report *diagnostics.DiagnosticReport) error {
	output, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report to JSON: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

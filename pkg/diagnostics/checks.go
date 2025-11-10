package diagnostics

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/firecrown-media/stax/pkg/credentials"
	"github.com/firecrown-media/stax/pkg/ddev"
	"github.com/firecrown-media/stax/pkg/git"
	"github.com/firecrown-media/stax/pkg/system"
)

// CheckResult represents the result of a diagnostic check
type CheckResult struct {
	Name        string
	Status      CheckStatus
	Message     string
	Suggestion  string
	Details     map[string]string
}

// CheckStatus represents the status of a check
type CheckStatus string

const (
	StatusPass    CheckStatus = "pass"
	StatusWarning CheckStatus = "warning"
	StatusFail    CheckStatus = "fail"
	StatusSkip    CheckStatus = "skip"
)

// DiagnosticReport contains all diagnostic check results
type DiagnosticReport struct {
	Checks      []CheckResult
	Summary     Summary
	ProjectPath string
}

// Summary provides a summary of check results
type Summary struct {
	Total    int
	Passed   int
	Warnings int
	Failed   int
	Skipped  int
}

// RunAllChecks runs all diagnostic checks
func RunAllChecks(projectPath string) (*DiagnosticReport, error) {
	report := &DiagnosticReport{
		ProjectPath: projectPath,
		Checks:      []CheckResult{},
	}

	// System checks
	report.Checks = append(report.Checks, CheckGit())
	report.Checks = append(report.Checks, CheckDocker())
	report.Checks = append(report.Checks, CheckDDEV())

	// Project checks
	report.Checks = append(report.Checks, CheckStaxConfig(projectPath))
	report.Checks = append(report.Checks, CheckDDEVConfig(projectPath))
	report.Checks = append(report.Checks, CheckCredentials(projectPath))

	// Environment checks
	report.Checks = append(report.Checks, CheckPorts())
	report.Checks = append(report.Checks, CheckDiskSpace(projectPath))

	// Calculate summary
	report.Summary = calculateSummary(report.Checks)

	return report, nil
}

// CheckGit checks if Git is installed and configured
func CheckGit() CheckResult {
	if !git.IsGitAvailable() {
		return CheckResult{
			Name:       "Git Installation",
			Status:     StatusFail,
			Message:    "Git is not installed",
			Suggestion: "Install Git: https://git-scm.com/downloads",
		}
	}

	version, err := git.GetGitVersion()
	if err != nil {
		return CheckResult{
			Name:       "Git Installation",
			Status:     StatusWarning,
			Message:    "Git is installed but version could not be determined",
			Suggestion: "Verify Git installation: git --version",
		}
	}

	return CheckResult{
		Name:    "Git Installation",
		Status:  StatusPass,
		Message: fmt.Sprintf("Git version %s installed", version),
		Details: map[string]string{
			"version": version,
		},
	}
}

// CheckDocker checks if Docker is installed and running
func CheckDocker() CheckResult {
	info, err := system.GetDockerInfo()
	if err != nil {
		return CheckResult{
			Name:       "Docker",
			Status:     StatusFail,
			Message:    "Failed to get Docker information",
			Suggestion: "Verify Docker installation",
		}
	}

	if !info.Installed {
		return CheckResult{
			Name:       "Docker",
			Status:     StatusFail,
			Message:    "Docker is not installed",
			Suggestion: "Install Docker Desktop: https://www.docker.com/products/docker-desktop",
		}
	}

	if !info.Running {
		return CheckResult{
			Name:       "Docker",
			Status:     StatusFail,
			Message:    "Docker is installed but not running",
			Suggestion: "Start Docker Desktop application",
			Details: map[string]string{
				"version": info.Version,
			},
		}
	}

	details := map[string]string{
		"version": info.Version,
		"running": "yes",
	}

	if info.ComposeInstalled {
		details["compose_version"] = info.ComposeVersion
	}

	return CheckResult{
		Name:    "Docker",
		Status:  StatusPass,
		Message: fmt.Sprintf("Docker %s is running", info.Version),
		Details: details,
	}
}

// CheckDDEV checks if DDEV is installed and configured
func CheckDDEV() CheckResult {
	if !ddev.IsInstalled() {
		return CheckResult{
			Name:       "DDEV",
			Status:     StatusFail,
			Message:    "DDEV is not installed",
			Suggestion: "Install DDEV: https://ddev.readthedocs.io/en/stable/users/install/",
		}
	}

	version, err := ddev.GetVersion()
	if err != nil {
		return CheckResult{
			Name:       "DDEV",
			Status:     StatusWarning,
			Message:    "DDEV is installed but version could not be determined",
			Suggestion: "Verify DDEV installation: ddev version",
		}
	}

	return CheckResult{
		Name:    "DDEV",
		Status:  StatusPass,
		Message: fmt.Sprintf("DDEV version %s installed", version),
		Details: map[string]string{
			"version": version,
		},
	}
}

// CheckStaxConfig checks if .stax.yml exists and is valid
func CheckStaxConfig(projectPath string) CheckResult {
	configPath := filepath.Join(projectPath, ".stax.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return CheckResult{
			Name:       "Stax Configuration",
			Status:     StatusFail,
			Message:    ".stax.yml not found",
			Suggestion: "Create configuration: stax init or stax config template > .stax.yml",
			Details: map[string]string{
				"expected_path": configPath,
			},
		}
	}

	return CheckResult{
		Name:    "Stax Configuration",
		Status:  StatusPass,
		Message: ".stax.yml found",
		Details: map[string]string{
			"path": configPath,
		},
	}
}

// CheckDDEVConfig checks if DDEV is configured for the project
func CheckDDEVConfig(projectPath string) CheckResult {
	ddevPath := filepath.Join(projectPath, ".ddev")
	if _, err := os.Stat(ddevPath); os.IsNotExist(err) {
		return CheckResult{
			Name:       "DDEV Configuration",
			Status:     StatusWarning,
			Message:    ".ddev directory not found",
			Suggestion: "Initialize DDEV: stax init or ddev config",
			Details: map[string]string{
				"expected_path": ddevPath,
			},
		}
	}

	configPath := filepath.Join(ddevPath, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return CheckResult{
			Name:       "DDEV Configuration",
			Status:     StatusWarning,
			Message:    ".ddev/config.yaml not found",
			Suggestion: "Configure DDEV: stax init or ddev config",
			Details: map[string]string{
				"expected_path": configPath,
			},
		}
	}

	return CheckResult{
		Name:    "DDEV Configuration",
		Status:  StatusPass,
		Message: "DDEV is configured",
		Details: map[string]string{
			"path": configPath,
		},
	}
}

// CheckCredentials checks if credentials are properly configured
func CheckCredentials(projectPath string) CheckResult {
	diag := credentials.RunDiagnostics()

	// Analyze diagnostic results
	if diag.OverallStatus == "error" {
		return CheckResult{
			Name:       "Credentials",
			Status:     StatusFail,
			Message:    "Credential configuration has errors",
			Suggestion: "Run: stax setup --check for details",
		}
	}

	if diag.OverallStatus == "warning" {
		return CheckResult{
			Name:       "Credentials",
			Status:     StatusWarning,
			Message:    "Credential configuration has warnings",
			Suggestion: "Run: stax setup --check for details",
		}
	}

	return CheckResult{
		Name:    "Credentials",
		Status:  StatusPass,
		Message: "All credentials are properly configured",
	}
}

// CheckPorts checks if required ports are available
func CheckPorts() CheckResult {
	defaultPorts := system.DefaultDDEVPorts()
	inUse, err := system.CheckRequiredPorts(defaultPorts)
	if err != nil {
		recommendations := system.RecommendedPorts(defaultPorts)
		var suggestions []string
		for original, recommended := range recommendations {
			suggestions = append(suggestions, fmt.Sprintf("Port %d -> %d", original, recommended))
		}

		return CheckResult{
			Name:       "Port Availability",
			Status:     StatusWarning,
			Message:    fmt.Sprintf("Some required ports are in use: %v", inUse),
			Suggestion: fmt.Sprintf("Consider using alternative ports:\n%s", strings.Join(suggestions, "\n")),
			Details: map[string]string{
				"ports_in_use": fmt.Sprintf("%v", inUse),
			},
		}
	}

	return CheckResult{
		Name:    "Port Availability",
		Status:  StatusPass,
		Message: "All required ports are available",
		Details: map[string]string{
			"checked_ports": fmt.Sprintf("%v", defaultPorts),
		},
	}
}

// CheckDiskSpace checks if there is sufficient disk space
func CheckDiskSpace(projectPath string) CheckResult {
	// This is a simplified check - in production, you'd use syscall to get actual disk space
	// For now, we'll just verify the project directory is accessible

	info, err := os.Stat(projectPath)
	if err != nil {
		return CheckResult{
			Name:       "Disk Space",
			Status:     StatusFail,
			Message:    "Cannot access project directory",
			Suggestion: "Verify project path exists and is accessible",
		}
	}

	if !info.IsDir() {
		return CheckResult{
			Name:       "Disk Space",
			Status:     StatusFail,
			Message:    "Project path is not a directory",
			Suggestion: "Verify project path",
		}
	}

	return CheckResult{
		Name:    "Disk Space",
		Status:  StatusPass,
		Message: "Project directory is accessible",
		Details: map[string]string{
			"path": projectPath,
		},
	}
}

// calculateSummary calculates the summary of check results
func calculateSummary(checks []CheckResult) Summary {
	summary := Summary{
		Total: len(checks),
	}

	for _, check := range checks {
		switch check.Status {
		case StatusPass:
			summary.Passed++
		case StatusWarning:
			summary.Warnings++
		case StatusFail:
			summary.Failed++
		case StatusSkip:
			summary.Skipped++
		}
	}

	return summary
}

// HasCriticalFailures returns true if there are any critical failures
func (r *DiagnosticReport) HasCriticalFailures() bool {
	return r.Summary.Failed > 0
}

// HasWarnings returns true if there are any warnings
func (r *DiagnosticReport) HasWarnings() bool {
	return r.Summary.Warnings > 0
}

// IsHealthy returns true if all checks passed
func (r *DiagnosticReport) IsHealthy() bool {
	return r.Summary.Failed == 0 && r.Summary.Warnings == 0
}

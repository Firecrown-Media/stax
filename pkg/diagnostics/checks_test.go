package diagnostics

import (
	"testing"
)

func TestCheckGit(t *testing.T) {
	result := CheckGit()

	if result.Name != "Git Installation" {
		t.Errorf("CheckGit() name = %v, want Git Installation", result.Name)
	}

	// We can't assert the status since it depends on system state,
	// but we can verify it's a valid status
	validStatuses := map[CheckStatus]bool{
		StatusPass:    true,
		StatusWarning: true,
		StatusFail:    true,
		StatusSkip:    true,
	}

	if !validStatuses[result.Status] {
		t.Errorf("CheckGit() status = %v, which is not a valid status", result.Status)
	}

	t.Logf("Git check result: %s - %s", result.Status, result.Message)
}

func TestCheckDocker(t *testing.T) {
	result := CheckDocker()

	if result.Name != "Docker" {
		t.Errorf("CheckDocker() name = %v, want Docker", result.Name)
	}

	t.Logf("Docker check result: %s - %s", result.Status, result.Message)
}

func TestCheckDDEV(t *testing.T) {
	result := CheckDDEV()

	if result.Name != "DDEV" {
		t.Errorf("CheckDDEV() name = %v, want DDEV", result.Name)
	}

	t.Logf("DDEV check result: %s - %s", result.Status, result.Message)
}

func TestCalculateSummary(t *testing.T) {
	checks := []CheckResult{
		{Status: StatusPass},
		{Status: StatusPass},
		{Status: StatusWarning},
		{Status: StatusFail},
		{Status: StatusSkip},
	}

	summary := calculateSummary(checks)

	if summary.Total != 5 {
		t.Errorf("summary.Total = %d, want 5", summary.Total)
	}

	if summary.Passed != 2 {
		t.Errorf("summary.Passed = %d, want 2", summary.Passed)
	}

	if summary.Warnings != 1 {
		t.Errorf("summary.Warnings = %d, want 1", summary.Warnings)
	}

	if summary.Failed != 1 {
		t.Errorf("summary.Failed = %d, want 1", summary.Failed)
	}

	if summary.Skipped != 1 {
		t.Errorf("summary.Skipped = %d, want 1", summary.Skipped)
	}
}

func TestDiagnosticReport_HasCriticalFailures(t *testing.T) {
	tests := []struct {
		name    string
		summary Summary
		want    bool
	}{
		{
			name:    "no failures",
			summary: Summary{Failed: 0},
			want:    false,
		},
		{
			name:    "has failures",
			summary: Summary{Failed: 1},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &DiagnosticReport{
				Summary: tt.summary,
			}

			got := report.HasCriticalFailures()
			if got != tt.want {
				t.Errorf("HasCriticalFailures() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiagnosticReport_IsHealthy(t *testing.T) {
	tests := []struct {
		name    string
		summary Summary
		want    bool
	}{
		{
			name:    "all passed",
			summary: Summary{Passed: 5, Failed: 0, Warnings: 0},
			want:    true,
		},
		{
			name:    "has warnings",
			summary: Summary{Passed: 4, Warnings: 1, Failed: 0},
			want:    false,
		},
		{
			name:    "has failures",
			summary: Summary{Passed: 4, Failed: 1, Warnings: 0},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &DiagnosticReport{
				Summary: tt.summary,
			}

			got := report.IsHealthy()
			if got != tt.want {
				t.Errorf("IsHealthy() = %v, want %v", got, tt.want)
			}
		})
	}
}

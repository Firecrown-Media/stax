package system

import (
	"testing"
)

func TestIsDockerAvailable(t *testing.T) {
	// This test will skip if Docker is not installed
	result := IsDockerAvailable()
	t.Logf("Docker available: %v", result)
}

func TestGetDockerInfo(t *testing.T) {
	info, err := GetDockerInfo()
	if err != nil {
		t.Fatalf("GetDockerInfo() error = %v", err)
	}

	if info.Installed && info.Version == "" {
		t.Error("Docker is installed but version is empty")
	}

	t.Logf("Docker info: %+v", info)
}

func TestIsVersionAtLeast(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		minVersion string
		want       bool
	}{
		{
			name:       "exact match",
			version:    "20.10.0",
			minVersion: "20.10.0",
			want:       true,
		},
		{
			name:       "newer version",
			version:    "24.0.5",
			minVersion: "20.10.0",
			want:       true,
		},
		{
			name:       "older version",
			version:    "19.03.0",
			minVersion: "20.10.0",
			want:       false,
		},
		{
			name:       "newer patch",
			version:    "20.10.1",
			minVersion: "20.10.0",
			want:       true,
		},
		{
			name:       "older patch",
			version:    "20.10.0",
			minVersion: "20.10.1",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isVersionAtLeast(tt.version, tt.minVersion)
			if got != tt.want {
				t.Errorf("isVersionAtLeast(%s, %s) = %v, want %v", tt.version, tt.minVersion, got, tt.want)
			}
		})
	}
}

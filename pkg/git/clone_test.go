package git

import (
	"testing"
)

func TestIsGitAvailable(t *testing.T) {
	// This should pass on any system with git installed
	if !IsGitAvailable() {
		t.Skip("git is not installed, skipping test")
	}
}

func TestGetGitVersion(t *testing.T) {
	if !IsGitAvailable() {
		t.Skip("git is not installed, skipping test")
	}

	version, err := GetGitVersion()
	if err != nil {
		t.Fatalf("GetGitVersion() error = %v", err)
	}

	if version == "" {
		t.Error("GetGitVersion() returned empty version")
	}
}

func TestGetRemoteURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "SSH URL with .git",
			url:     "git@github.com:user/repo.git",
			want:    "user/repo",
			wantErr: false,
		},
		{
			name:    "SSH URL without .git",
			url:     "git@github.com:user/repo",
			want:    "user/repo",
			wantErr: false,
		},
		{
			name:    "HTTPS URL with .git",
			url:     "https://github.com/user/repo.git",
			want:    "user/repo",
			wantErr: false,
		},
		{
			name:    "HTTPS URL without .git",
			url:     "https://github.com/user/repo",
			want:    "user/repo",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			url:     "not-a-git-url",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRemoteURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemoteURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRemoteURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid SSH URL",
			url:     "git@github.com:user/repo.git",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL",
			url:     "https://github.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "valid HTTP URL",
			url:     "http://github.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			url:     "not-a-git-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

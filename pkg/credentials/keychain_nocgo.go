//go:build (darwin && !cgo) || (!darwin && !linux)
// +build darwin,!cgo !darwin,!linux

package credentials

import (
	"fmt"
	"os"
)

// WPEngineCredentials represents WPEngine API credentials
type WPEngineCredentials struct {
	APIUser     string `json:"api_user"`
	APIPassword string `json:"api_password"`
	SSHUser     string `json:"ssh_user"`
	SSHGateway  string `json:"ssh_gateway"`
}

// GitHubCredentials holds GitHub API credentials
type GitHubCredentials struct {
	Token string `json:"token"`
}

// SSHCredentials holds SSH private key
type SSHCredentials struct {
	PrivateKey string `json:"private_key"`
}

// GetWPEngineCredentials retrieves WPEngine credentials
// Falls back to environment variables or config file
func GetWPEngineCredentials(install string) (*WPEngineCredentials, error) {
	// Try environment variables first
	if apiUser := os.Getenv("WPENGINE_API_USER"); apiUser != "" {
		return &WPEngineCredentials{
			APIUser:     apiUser,
			APIPassword: os.Getenv("WPENGINE_API_PASSWORD"),
			SSHUser:     os.Getenv("WPENGINE_SSH_USER"),
			SSHGateway:  getEnvOrDefault("WPENGINE_SSH_GATEWAY", "ssh.wpengine.net"),
		}, nil
	}

	// Try credentials file
	creds, err := LoadCredentialsFile()
	if err == nil {
		return &WPEngineCredentials{
			APIUser:     creds.WPEngine.APIUser,
			APIPassword: creds.WPEngine.APIPassword,
			SSHUser:     creds.WPEngine.SSHUser,
			SSHGateway:  creds.WPEngine.SSHGateway,
		}, nil
	}

	// Return helpful error
	return nil, &KeychainUnavailableError{Operation: "get WPEngine credentials"}
}

// SetWPEngineCredentials stores WPEngine credentials (not supported)
func SetWPEngineCredentials(install string, creds *WPEngineCredentials) error {
	return &KeychainUnavailableError{Operation: "store WPEngine credentials"}
}

// GetGitHubToken retrieves GitHub token
func GetGitHubToken(organization string) (string, error) {
	// Try environment variable first
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	// Try credentials file
	creds, err := LoadCredentialsFile()
	if err == nil && creds.GitHub.Token != "" {
		return creds.GitHub.Token, nil
	}

	return "", &KeychainUnavailableError{Operation: "get GitHub token"}
}

// SetGitHubToken stores GitHub token (not supported)
func SetGitHubToken(organization string, token string) error {
	return &KeychainUnavailableError{Operation: "store GitHub token"}
}

// GetSSHPrivateKey retrieves SSH private key
func GetSSHPrivateKey(account string) (string, error) {
	// Try environment variable first (base64 encoded or direct path)
	if key := os.Getenv("WPENGINE_SSH_KEY"); key != "" {
		return key, nil
	}

	// Try credentials file (stores path to key)
	creds, err := LoadCredentialsFile()
	if err == nil && creds.SSH.PrivateKeyPath != "" {
		keyData, err := os.ReadFile(creds.SSH.PrivateKeyPath)
		if err != nil {
			return "", fmt.Errorf("failed to read SSH key from %s: %w",
				creds.SSH.PrivateKeyPath, err)
		}
		return string(keyData), nil
	}

	return "", &KeychainUnavailableError{Operation: "get SSH private key"}
}

// SetSSHPrivateKey stores SSH private key (not supported)
func SetSSHPrivateKey(account string, privateKey string) error {
	return &KeychainUnavailableError{Operation: "store SSH private key"}
}

// DeleteWPEngineCredentials deletes WPEngine credentials (not supported)
func DeleteWPEngineCredentials(install string) error {
	return &KeychainUnavailableError{Operation: "delete WPEngine credentials"}
}

// DeleteGitHubToken deletes GitHub token (not supported)
func DeleteGitHubToken(organization string) error {
	return &KeychainUnavailableError{Operation: "delete GitHub token"}
}

// DeleteSSHPrivateKey deletes SSH private key (not supported)
func DeleteSSHPrivateKey(account string) error {
	return &KeychainUnavailableError{Operation: "delete SSH private key"}
}

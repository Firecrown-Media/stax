//go:build (darwin && !cgo) || (!darwin && !linux)
// +build darwin,!cgo !darwin,!linux

package credentials

import (
	"fmt"
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

// GetWPEngineCredentials retrieves WPEngine credentials (not supported without CGO on Darwin or on unsupported platforms)
func GetWPEngineCredentials(install string) (*WPEngineCredentials, error) {
	return nil, fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

// SetWPEngineCredentials stores WPEngine credentials (not supported without CGO on Darwin or on unsupported platforms)
func SetWPEngineCredentials(install string, creds *WPEngineCredentials) error {
	return fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

// GetGitHubToken retrieves GitHub token (not supported without CGO on Darwin or on unsupported platforms)
func GetGitHubToken(organization string) (string, error) {
	return "", fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

// SetGitHubToken stores GitHub token (not supported without CGO on Darwin or on unsupported platforms)
func SetGitHubToken(organization string, token string) error {
	return fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

// GetSSHPrivateKey retrieves SSH private key (not supported without CGO on Darwin or on unsupported platforms)
func GetSSHPrivateKey(account string) (string, error) {
	return "", fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

// SetSSHPrivateKey stores SSH private key (not supported without CGO on Darwin or on unsupported platforms)
func SetSSHPrivateKey(account string, privateKey string) error {
	return fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

// DeleteWPEngineCredentials deletes WPEngine credentials (not supported without CGO on Darwin or on unsupported platforms)
func DeleteWPEngineCredentials(install string) error {
	return fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

// DeleteGitHubToken deletes GitHub token (not supported without CGO on Darwin or on unsupported platforms)
func DeleteGitHubToken(organization string) error {
	return fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

// DeleteSSHPrivateKey deletes SSH private key (not supported without CGO on Darwin or on unsupported platforms)
func DeleteSSHPrivateKey(account string) error {
	return fmt.Errorf("keychain storage is not supported - please use environment variables or config files")
}

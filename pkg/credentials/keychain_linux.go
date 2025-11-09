//go:build linux
// +build linux

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

// KeychainManager provides a stub implementation for Linux
type KeychainManager struct{}

// NewKeychainManager creates a new keychain manager
func NewKeychainManager() *KeychainManager {
	return &KeychainManager{}
}

// Store stores credentials (not supported on Linux)
func (k *KeychainManager) Store(service, account string, data map[string]string) error {
	return fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// Retrieve retrieves credentials (not supported on Linux)
func (k *KeychainManager) Retrieve(service, account string) (map[string]string, error) {
	return nil, fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// Delete deletes credentials (not supported on Linux)
func (k *KeychainManager) Delete(service, account string) error {
	return fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// GetWPEngineCredentials retrieves WPEngine credentials (not supported on Linux)
func GetWPEngineCredentials(install string) (*WPEngineCredentials, error) {
	return nil, fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// SetWPEngineCredentials stores WPEngine credentials (not supported on Linux)
func SetWPEngineCredentials(install string, creds *WPEngineCredentials) error {
	return fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// GetGitHubToken retrieves GitHub token (not supported on Linux)
func GetGitHubToken(organization string) (string, error) {
	return "", fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// SetGitHubToken stores GitHub token (not supported on Linux)
func SetGitHubToken(organization string, token string) error {
	return fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// GetSSHPrivateKey retrieves SSH private key (not supported on Linux)
func GetSSHPrivateKey(account string) (string, error) {
	return "", fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// SetSSHPrivateKey stores SSH private key (not supported on Linux)
func SetSSHPrivateKey(account string, privateKey string) error {
	return fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// DeleteWPEngineCredentials deletes WPEngine credentials (not supported on Linux)
func DeleteWPEngineCredentials(install string) error {
	return fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// DeleteGitHubToken deletes GitHub token (not supported on Linux)
func DeleteGitHubToken(organization string) error {
	return fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

// DeleteSSHPrivateKey deletes SSH private key (not supported on Linux)
func DeleteSSHPrivateKey(account string) error {
	return fmt.Errorf("keychain storage is not supported on Linux - please use environment variables or config files")
}

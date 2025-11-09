// +build darwin

package credentials

import (
	"encoding/json"
	"fmt"

	"github.com/keybase/go-keychain"
)

const (
	// Keychain service names
	ServiceWPEngine = "com.firecrown.stax.wpengine"
	ServiceGitHub   = "com.firecrown.stax.github"
	ServiceSSH      = "com.firecrown.stax.ssh"
)

// WPEngineCredentials represents WPEngine API credentials
type WPEngineCredentials struct {
	APIUser     string `json:"api_user"`
	APIPassword string `json:"api_password"`
	SSHUser     string `json:"ssh_user"`
	SSHGateway  string `json:"ssh_gateway"`
}

// GitHubCredentials represents GitHub credentials
type GitHubCredentials struct {
	Token string `json:"token"`
}

// SSHCredentials represents SSH credentials
type SSHCredentials struct {
	PrivateKey string `json:"private_key"`
}

// GetWPEngineCredentials retrieves WPEngine credentials from Keychain
func GetWPEngineCredentials(install string) (*WPEngineCredentials, error) {
	password, err := getPassword(ServiceWPEngine, install)
	if err != nil {
		return nil, fmt.Errorf("failed to get WPEngine credentials: %w", err)
	}

	var creds WPEngineCredentials
	if err := json.Unmarshal([]byte(password), &creds); err != nil {
		return nil, fmt.Errorf("failed to parse WPEngine credentials: %w", err)
	}

	return &creds, nil
}

// SetWPEngineCredentials stores WPEngine credentials in Keychain
func SetWPEngineCredentials(install string, creds *WPEngineCredentials) error {
	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to marshal WPEngine credentials: %w", err)
	}

	if err := setPassword(ServiceWPEngine, install, string(data)); err != nil {
		return fmt.Errorf("failed to set WPEngine credentials: %w", err)
	}

	return nil
}

// GetGitHubToken retrieves GitHub token from Keychain
func GetGitHubToken(organization string) (string, error) {
	password, err := getPassword(ServiceGitHub, organization)
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub token: %w", err)
	}

	return password, nil
}

// SetGitHubToken stores GitHub token in Keychain
func SetGitHubToken(organization string, token string) error {
	if err := setPassword(ServiceGitHub, organization, token); err != nil {
		return fmt.Errorf("failed to set GitHub token: %w", err)
	}

	return nil
}

// GetSSHPrivateKey retrieves SSH private key from Keychain
func GetSSHPrivateKey(account string) (string, error) {
	password, err := getPassword(ServiceSSH, account)
	if err != nil {
		return "", fmt.Errorf("failed to get SSH private key: %w", err)
	}

	return password, nil
}

// SetSSHPrivateKey stores SSH private key in Keychain
func SetSSHPrivateKey(account string, privateKey string) error {
	if err := setPassword(ServiceSSH, account, privateKey); err != nil {
		return fmt.Errorf("failed to set SSH private key: %w", err)
	}

	return nil
}

// DeleteWPEngineCredentials removes WPEngine credentials from Keychain
func DeleteWPEngineCredentials(install string) error {
	return deletePassword(ServiceWPEngine, install)
}

// DeleteGitHubToken removes GitHub token from Keychain
func DeleteGitHubToken(organization string) error {
	return deletePassword(ServiceGitHub, organization)
}

// DeleteSSHPrivateKey removes SSH private key from Keychain
func DeleteSSHPrivateKey(account string) error {
	return deletePassword(ServiceSSH, account)
}

// Low-level Keychain operations

func getPassword(service, account string) (string, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetAccount(account)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", fmt.Errorf("no credentials found for %s/%s", service, account)
	}

	return string(results[0].Data), nil
}

func setPassword(service, account, password string) error {
	// First, try to delete any existing item
	deletePassword(service, account)

	// Add new item
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(account)
	item.SetLabel(fmt.Sprintf("Stax - %s", account))
	item.SetData([]byte(password))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	return keychain.AddItem(item)
}

func deletePassword(service, account string) error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(account)

	return keychain.DeleteItem(item)
}

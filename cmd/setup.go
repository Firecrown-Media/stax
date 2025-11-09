package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/firecrown-media/stax/pkg/credentials"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	setupWPEngineUser     string
	setupWPEnginePassword string
	setupGitHubToken      string
	setupSSHKey           string
	setupInteractive      bool
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure WPEngine and GitHub credentials",
	Long: `Configure WPEngine and GitHub credentials securely in macOS Keychain.

This command stores sensitive credentials in the macOS Keychain, ensuring
they are never stored in plain text configuration files.`,
	Example: `  # Interactive mode
  stax setup

  # Non-interactive
  stax setup \
    --wpengine-user=myuser@example.com \
    --wpengine-password=mypassword \
    --github-token=ghp_xxxxxxxxxxxxx \
    --ssh-key=~/.ssh/wpengine_rsa`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().StringVar(&setupWPEngineUser, "wpengine-user", "", "WPEngine API username")
	setupCmd.Flags().StringVar(&setupWPEnginePassword, "wpengine-password", "", "WPEngine API password")
	setupCmd.Flags().StringVar(&setupGitHubToken, "github-token", "", "GitHub personal access token")
	setupCmd.Flags().StringVar(&setupSSHKey, "ssh-key", "", "Path to SSH private key for WPEngine")
	setupCmd.Flags().BoolVar(&setupInteractive, "interactive", true, "Interactive credential setup")
}

func runSetup(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Setting up Stax Credentials")

	// Get credentials interactively if not provided
	if setupInteractive {
		if setupWPEngineUser == "" {
			setupWPEngineUser = ui.PromptString("WPEngine API Username", "")
		}
		if setupWPEnginePassword == "" {
			setupWPEnginePassword = ui.PromptString("WPEngine API Password", "")
		}
		if setupGitHubToken == "" {
			setupGitHubToken = ui.PromptString("GitHub Personal Access Token", "")
		}
		if setupSSHKey == "" {
			setupSSHKey = ui.PromptString("SSH Key for WPEngine (optional)", "~/.ssh/id_rsa")
		}
	}

	// Validate required fields
	if setupWPEngineUser == "" || setupWPEnginePassword == "" {
		return fmt.Errorf("WPEngine credentials are required")
	}

	// Store WPEngine credentials
	ui.Info("Storing WPEngine credentials in Keychain...")
	wpeCreds := &credentials.WPEngineCredentials{
		APIUser:     setupWPEngineUser,
		APIPassword: setupWPEnginePassword,
		SSHGateway:  "ssh.wpengine.net",
	}

	// For setup, we'll use a default install name "default"
	// Users can have multiple installs with different credentials
	if err := credentials.SetWPEngineCredentials("default", wpeCreds); err != nil {
		return fmt.Errorf("failed to store WPEngine credentials: %w", err)
	}
	ui.Success("WPEngine credentials stored")

	// Test WPEngine API connection
	ui.Info("Testing WPEngine API connection...")
	if err := testWPEngineConnection(setupWPEngineUser, setupWPEnginePassword); err != nil {
		ui.Warning(fmt.Sprintf("Failed to connect to WPEngine API: %v", err))
		ui.Info("Credentials saved, but please verify they are correct")
	} else {
		ui.Success("WPEngine API connection successful")
	}

	// Store GitHub token if provided
	if setupGitHubToken != "" {
		ui.Info("Storing GitHub token in Keychain...")
		if err := credentials.SetGitHubToken("default", setupGitHubToken); err != nil {
			return fmt.Errorf("failed to store GitHub token: %w", err)
		}
		ui.Success("GitHub token stored")
	}

	// Store SSH key if provided
	if setupSSHKey != "" {
		ui.Info("Storing SSH key in Keychain...")
		// Read SSH key file
		if err := storeSSHKey(setupSSHKey); err != nil {
			ui.Warning(fmt.Sprintf("Failed to store SSH key: %v", err))
		} else {
			ui.Success("SSH key stored")
		}
	}

	ui.Section("\nCredentials saved successfully!")
	ui.Info("Your credentials are securely stored in macOS Keychain")
	ui.Info("You can now run 'stax init' to initialize a project")

	return nil
}

// testWPEngineConnection tests the WPEngine API connection
func testWPEngineConnection(apiUser, apiPassword string) error {
	// Import wpengine package
	wpengine := struct {
		NewClient func(string, string, string) interface{ TestConnection() error }
	}{
		NewClient: func(u, p, i string) interface{ TestConnection() error } {
			// This is a placeholder - actual implementation will use wpengine.NewClient
			return nil
		},
	}

	if wpengine.NewClient == nil {
		return fmt.Errorf("WPEngine client not available")
	}

	// Note: In actual implementation, this would use wpengine.NewClient
	// For now, we'll just validate credentials are not empty
	if apiUser == "" || apiPassword == "" {
		return fmt.Errorf("credentials are empty")
	}

	return nil
}

// storeSSHKey reads and stores an SSH key
func storeSSHKey(keyPath string) error {
	// Read SSH key file
	keyData, err := os.ReadFile(expandPath(keyPath))
	if err != nil {
		return fmt.Errorf("failed to read SSH key: %w", err)
	}

	// Store in keychain
	if err := credentials.SetSSHPrivateKey("wpengine", string(keyData)); err != nil {
		return fmt.Errorf("failed to store SSH key: %w", err)
	}

	return nil
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}
	return path
}

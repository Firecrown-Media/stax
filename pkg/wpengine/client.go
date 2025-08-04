package wpengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Firecrown-Media/stax/pkg/ui"
)

type Client struct {
	config     Config
	httpClient *http.Client
	baseURL    string
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.wpengineapi.com/v1",
	}
}

func (c *Client) GetInstallInfo() (*InstallInfo, error) {
	// Test SSH connection to WP Engine using the install name as both username and hostname
	sshHost := fmt.Sprintf("%s@%s.ssh.wpengine.net", c.config.InstallName, c.config.InstallName)
	
	// Simple SSH test command
	cmd := exec.Command("ssh", 
		"-o", "ConnectTimeout=30",
		"-o", "ServerAliveInterval=10", 
		"-o", "ServerAliveCountMax=3",
		"-o", "StrictHostKeyChecking=accept-new",
		"-o", "BatchMode=yes", 
		sshHost, "echo 'connection_test'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("SSH connection failed. Make sure your SSH key is configured with WP Engine and you can connect to %s. Error: %s", sshHost, string(output))
	}

	// If SSH works, get basic info
	cmd = exec.Command("ssh", 
		"-o", "ConnectTimeout=30",
		"-o", "ServerAliveInterval=10", 
		"-o", "ServerAliveCountMax=3",
		"-o", "StrictHostKeyChecking=accept-new",
		sshHost, "pwd")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get install info via SSH: %s", string(output))
	}

	// Parse the path to determine install name and environment
	_ = strings.TrimSpace(string(output)) // installPath not used currently
	
	return &InstallInfo{
		Name:        c.config.InstallName,
		Environment: c.config.Environment,
		Domain:      fmt.Sprintf("%s.wpenginepowered.com", c.config.InstallName),
		Status:      "running",
		PHPVersion:  "8.0", // Default, can be detected via SSH if needed
	}, nil
}

func (c *Client) ListInstalls() ([]InstallInfo, error) {
	req, err := c.newRequest("GET", "/installs", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response struct {
		Results []InstallInfo `json:"results"`
	}
	
	if err := c.doRequest(req, &response); err != nil {
		return nil, fmt.Errorf("failed to list installs: %w", err)
	}

	return response.Results, nil
}

func (c *Client) ListBackups() ([]BackupInfo, error) {
	endpoint := fmt.Sprintf("/installs/%s/backups", c.config.InstallName)
	
	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response struct {
		Results []BackupInfo `json:"results"`
	}
	
	if err := c.doRequest(req, &response); err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	return response.Results, nil
}

func (c *Client) CreateBackup(backupType string) (*BackupInfo, error) {
	endpoint := fmt.Sprintf("/installs/%s/backups", c.config.InstallName)
	
	payload := map[string]string{
		"type":        backupType,
		"description": fmt.Sprintf("Stax sync backup - %s", time.Now().Format("2006-01-02 15:04:05")),
	}

	req, err := c.newRequest("POST", endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var backup BackupInfo
	if err := c.doRequest(req, &backup); err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	return &backup, nil
}

func (c *Client) GetBackup(backupID string) (*BackupInfo, error) {
	endpoint := fmt.Sprintf("/installs/%s/backups/%s", c.config.InstallName, backupID)
	
	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var backup BackupInfo
	if err := c.doRequest(req, &backup); err != nil {
		return nil, fmt.Errorf("failed to get backup: %w", err)
	}

	return &backup, nil
}

func (c *Client) DownloadDatabase(outputPath string) (*DatabaseSyncResult, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	var result *DatabaseSyncResult
	
	err := ui.WithSpinnerResult(fmt.Sprintf("Downloading database from %s", c.config.InstallName), func() error {
		// Use mysqldump via SSH to WP Engine using the install name as both username and hostname
		sshHost := fmt.Sprintf("%s@%s.ssh.wpengine.net", c.config.InstallName, c.config.InstallName)
		
		// WP Engine provides database credentials in the environment
		// Use their standard backup location or create a fresh dump with explicit path
		dumpCmd := fmt.Sprintf("wp db export - --single-transaction --path=/home/wpe-user/sites/%s", c.config.InstallName)
		
		fmt.Printf("Debug: SSH command: ssh %s %s\n", sshHost, dumpCmd)
		
		// Add SSH options to help with connection issues
		cmd := exec.Command("ssh", 
			"-o", "ConnectTimeout=30",
			"-o", "ServerAliveInterval=10", 
			"-o", "ServerAliveCountMax=3",
			"-o", "StrictHostKeyChecking=accept-new",
			"-v", // Verbose output for debugging
			sshHost, dumpCmd)
		
		// Create output file
		outFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outFile.Close()
		
		cmd.Stdout = outFile
		
		// Capture stderr for debugging
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		
		// Run the command
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to download database via SSH: %w - SSH Error: %s", err, stderr.String())
		}

		result = &DatabaseSyncResult{
			Success:      true,
			BackupID:     "ssh-dump",
			DatabaseFile: outputPath,
		}
		
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) SyncFiles(localPath string, options SyncOptions) (*FilesSyncResult, error) {
	var result *FilesSyncResult
	
	fmt.Printf("🔄 Starting file sync from %s...\n", c.config.InstallName)
	
	err := ui.WithSpinnerResult(fmt.Sprintf("Syncing files from %s (excluding media)", c.config.InstallName), func() error {
		// Use the correct WP Engine path structure using the install name as both username and hostname
		sshHost := fmt.Sprintf("%s@%s.ssh.wpengine.net", c.config.InstallName, c.config.InstallName)
		remotePath := fmt.Sprintf("%s:/home/wpe-user/sites/%s/", sshHost, c.config.InstallName)

		excludeArgs := []string{}
		
		// Default excludes for media and other large files
		if options.SkipMedia {
			excludeArgs = append(excludeArgs, "--exclude=wp-content/uploads/")
		}
		
		// Add custom excludes
		for _, exclude := range options.ExcludeDirs {
			excludeArgs = append(excludeArgs, fmt.Sprintf("--exclude=%s", exclude))
		}

		// Add common excludes
		commonExcludes := []string{
			"--exclude=.git/",
			"--exclude=node_modules/",
			"--exclude=.DS_Store",
			"--exclude=*.log",
			"--exclude=wp-config.php",      // Production config - don't overwrite local
			"--exclude=wp-config-ddev.php", // DDEV-specific config
			"--exclude=.ddev/",             // DDEV configuration directory
			"--exclude=.env",               // Environment files
			"--exclude=.env.local",         // Local environment overrides
			"--exclude=tmp/",               // Temporary files
			"--exclude=*.sql",              // Database dumps
		}
		excludeArgs = append(excludeArgs, commonExcludes...)

		args := []string{
			"-avz",
			"--progress",
			"--stats", // Add statistics to see what's happening
		}
		
		fmt.Printf("Debug: Syncing from %s to %s\n", remotePath, localPath)
		
		// Only add --delete flag if explicitly requested (dangerous for local development)
		if options.DeleteLocal {
			args = append(args, "--delete")
		}
		args = append(args, excludeArgs...)
		args = append(args, remotePath, localPath)

		cmd := exec.Command("rsync", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("rsync failed: %s", string(output))
		}

		// Parse rsync output for statistics
		lines := strings.Split(string(output), "\n")
		syncedFiles := 0
		for _, line := range lines {
			if strings.Contains(line, "files transferred") {
				// Parse the number of files transferred
				// This is a simplified parser
				fmt.Sscanf(line, "%d", &syncedFiles)
				break
			}
		}

		result = &FilesSyncResult{
			Success:     true,
			SyncedFiles: syncedFiles,
		}
		
		// Store output for debugging outside spinner
		result.Error = string(output) // Temporarily use Error field for debug output
		
		return nil
	})

	if err != nil {
		return &FilesSyncResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Show debug output after spinner completes
	if result != nil && result.Error != "" {
		fmt.Printf("Debug: rsync output:\n%s\n", result.Error)
		result.Error = "" // Clear debug output from error field
	}

	return result, nil
}

func (c *Client) newRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	url := c.baseURL + endpoint
	
	var buf io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	
	// Use API key authentication if available
	if c.config.APIKey != "" {
		req.SetBasicAuth(c.config.Username, c.config.APIKey)
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

func (c *Client) downloadFile(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create output directory if it doesn't exist
	if err := exec.Command("mkdir", "-p", filepath.Dir(outputPath)).Run(); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use curl for downloading large files with progress
	cmd := exec.Command("curl", "-L", "-o", outputPath, url)
	return cmd.Run()
}

func (c *Client) waitForBackup(backupID string) (*BackupInfo, error) {
	maxAttempts := 60 // Wait up to 30 minutes (30s intervals)
	
	for attempt := 0; attempt < maxAttempts; attempt++ {
		backup, err := c.GetBackup(backupID)
		if err != nil {
			return nil, err
		}

		switch backup.Status {
		case "completed":
			return backup, nil
		case "failed", "cancelled":
			return nil, fmt.Errorf("backup %s failed with status: %s", backupID, backup.Status)
		default:
			fmt.Printf("Backup in progress... (%s)\n", backup.Status)
			time.Sleep(30 * time.Second)
		}
	}

	return nil, fmt.Errorf("backup did not complete within timeout")
}

func (c *Client) TestConnection() error {
	_, err := c.GetInstallInfo()
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	return nil
}
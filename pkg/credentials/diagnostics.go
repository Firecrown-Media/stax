package credentials

import (
	"fmt"
	"os"
	"path/filepath"
)

// DiagnosticResult represents a single diagnostic check result
type DiagnosticResult struct {
	Name    string
	Status  string // "ok", "warning", "error"
	Message string
	Details []string
}

// CredentialDiagnostics contains all credential diagnostic results
type CredentialDiagnostics struct {
	KeychainAvailable  DiagnosticResult
	WPEngineAPI        DiagnosticResult
	WPEngineSSH        DiagnosticResult
	GitHubToken        DiagnosticResult
	CredentialsFile    DiagnosticResult
	EnvironmentVars    DiagnosticResult
	SSHKeyFile         DiagnosticResult
	OverallStatus      string
	RecommendedActions []string
}

// RunDiagnostics performs comprehensive credential diagnostics
func RunDiagnostics() *CredentialDiagnostics {
	diag := &CredentialDiagnostics{}

	// Check keychain availability
	diag.KeychainAvailable = checkKeychain()

	// Check WPEngine API credentials
	diag.WPEngineAPI = checkWPEngineAPI()

	// Check WPEngine SSH credentials
	diag.WPEngineSSH = checkWPEngineSSH()

	// Check GitHub token
	diag.GitHubToken = checkGitHubToken()

	// Check credentials file
	diag.CredentialsFile = checkCredentialsFile()

	// Check environment variables
	diag.EnvironmentVars = checkEnvironmentVariables()

	// Check SSH key file
	diag.SSHKeyFile = checkSSHKeyFile()

	// Determine overall status and recommendations
	diag.calculateOverallStatus()
	diag.generateRecommendations()

	return diag
}

// checkKeychain checks if keychain is available
func checkKeychain() DiagnosticResult {
	result := DiagnosticResult{
		Name: "Keychain Storage",
	}

	if IsKeychainAvailable() {
		result.Status = "ok"
		result.Message = "macOS Keychain is available"
		result.Details = []string{
			"Credentials can be stored securely in system keychain",
			"This is the recommended storage method",
		}
	} else {
		result.Status = "warning"
		result.Message = "Keychain storage not available"
		result.Details = []string{
			"Falling back to file-based or environment variable storage",
			"Consider rebuilding with CGO enabled for keychain support",
		}
	}

	return result
}

// checkWPEngineAPI checks WPEngine API credentials
func checkWPEngineAPI() DiagnosticResult {
	result := DiagnosticResult{
		Name: "WPEngine API Credentials",
	}

	creds, err := GetWPEngineCredentials("default")
	if err != nil || creds == nil {
		result.Status = "error"
		result.Message = "WPEngine API credentials not found"
		result.Details = []string{
			"Run 'stax setup wpengine' to configure credentials",
			"Credentials are required for database and file operations",
		}
	} else if creds.APIUser == "" || creds.APIPassword == "" {
		result.Status = "error"
		result.Message = "WPEngine API credentials incomplete"
		result.Details = []string{
			"API user or password is missing",
			"Run 'stax setup wpengine' to reconfigure",
		}
	} else {
		result.Status = "ok"
		result.Message = "WPEngine API credentials configured"
		result.Details = []string{
			fmt.Sprintf("API User: %s", creds.APIUser),
			"Credentials are ready for use",
		}
	}

	return result
}

// checkWPEngineSSH checks WPEngine SSH credentials
func checkWPEngineSSH() DiagnosticResult {
	result := DiagnosticResult{
		Name: "WPEngine SSH Credentials",
	}

	creds, err := GetWPEngineCredentials("default")
	if err != nil || creds == nil {
		result.Status = "error"
		result.Message = "SSH credentials not found"
		result.Details = []string{
			"SSH credentials are required for file operations",
		}
	} else if creds.SSHUser == "" {
		result.Status = "warning"
		result.Message = "SSH user not configured"
		result.Details = []string{
			"SSH user will default to API user",
			"Set explicitly in setup if different",
		}
	} else {
		result.Status = "ok"
		result.Message = "WPEngine SSH credentials configured"
		result.Details = []string{
			fmt.Sprintf("SSH User: %s", creds.SSHUser),
			fmt.Sprintf("SSH Gateway: %s", creds.SSHGateway),
		}
	}

	return result
}

// checkGitHubToken checks GitHub token
func checkGitHubToken() DiagnosticResult {
	result := DiagnosticResult{
		Name: "GitHub Token",
	}

	token, err := GetGitHubToken("default")
	if err != nil || token == "" {
		result.Status = "warning"
		result.Message = "GitHub token not configured"
		result.Details = []string{
			"GitHub token is optional but recommended",
			"Required for private repository access",
			"Prevents API rate limiting",
		}
	} else {
		result.Status = "ok"
		result.Message = "GitHub token configured"
		result.Details = []string{
			"Token is available for repository operations",
		}
	}

	return result
}

// checkCredentialsFile checks credentials file
func checkCredentialsFile() DiagnosticResult {
	result := DiagnosticResult{
		Name: "Credentials File",
	}

	path, err := GetCredentialsFilePath()
	if err != nil {
		result.Status = "error"
		result.Message = "Cannot determine credentials file path"
		result.Details = []string{err.Error()}
		return result
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		result.Status = "ok"
		result.Message = "Credentials file not in use"
		result.Details = []string{
			"Using keychain or environment variables",
			fmt.Sprintf("File would be at: %s", path),
		}
	} else if err != nil {
		result.Status = "error"
		result.Message = "Cannot access credentials file"
		result.Details = []string{err.Error()}
	} else {
		// File exists, check permissions
		mode := info.Mode()
		if mode.Perm() != 0600 {
			result.Status = "warning"
			result.Message = "Credentials file has insecure permissions"
			result.Details = []string{
				fmt.Sprintf("Current permissions: %v", mode.Perm()),
				fmt.Sprintf("Run: chmod 600 %s", path),
				"Recommended: 0600 (read/write for owner only)",
			}
		} else {
			result.Status = "ok"
			result.Message = "Credentials file configured securely"
			result.Details = []string{
				fmt.Sprintf("Location: %s", path),
				"Permissions are secure (0600)",
			}
		}
	}

	return result
}

// checkEnvironmentVariables checks environment variables
func checkEnvironmentVariables() DiagnosticResult {
	result := DiagnosticResult{
		Name: "Environment Variables",
	}

	envVars := []string{
		"WPENGINE_API_USER",
		"WPENGINE_API_PASSWORD",
		"WPENGINE_SSH_USER",
		"WPENGINE_SSH_GATEWAY",
		"GITHUB_TOKEN",
	}

	foundVars := []string{}
	for _, envVar := range envVars {
		if os.Getenv(envVar) != "" {
			foundVars = append(foundVars, envVar)
		}
	}

	if len(foundVars) == 0 {
		result.Status = "ok"
		result.Message = "Environment variables not in use"
		result.Details = []string{
			"Using keychain or credentials file instead",
		}
	} else {
		result.Status = "ok"
		result.Message = "Environment variables configured"
		result.Details = append([]string{
			fmt.Sprintf("Found %d credential variables", len(foundVars)),
		}, foundVars...)
	}

	return result
}

// checkSSHKeyFile checks SSH key file
func checkSSHKeyFile() DiagnosticResult {
	result := DiagnosticResult{
		Name: "SSH Key File",
	}

	home, err := os.UserHomeDir()
	if err != nil {
		result.Status = "error"
		result.Message = "Cannot determine home directory"
		result.Details = []string{err.Error()}
		return result
	}

	sshKeyPath := filepath.Join(home, ".ssh", "id_rsa")
	info, err := os.Stat(sshKeyPath)

	if os.IsNotExist(err) {
		result.Status = "warning"
		result.Message = "Default SSH key not found"
		result.Details = []string{
			fmt.Sprintf("Checked: %s", sshKeyPath),
			"SSH key is required for WPEngine file operations",
			"Generate with: ssh-keygen -t rsa -b 4096",
			"Add public key to WPEngine User Portal",
		}
	} else if err != nil {
		result.Status = "error"
		result.Message = "Cannot access SSH key file"
		result.Details = []string{err.Error()}
	} else {
		mode := info.Mode()
		if mode.Perm()&0077 != 0 {
			result.Status = "warning"
			result.Message = "SSH key has insecure permissions"
			result.Details = []string{
				fmt.Sprintf("Current permissions: %v", mode.Perm()),
				fmt.Sprintf("Run: chmod 600 %s", sshKeyPath),
				"SSH requires private keys to be secure (0600)",
			}
		} else {
			result.Status = "ok"
			result.Message = "SSH key found and secure"
			result.Details = []string{
				fmt.Sprintf("Location: %s", sshKeyPath),
				"Permissions are secure",
			}
		}
	}

	return result
}

// calculateOverallStatus determines the overall diagnostic status
func (d *CredentialDiagnostics) calculateOverallStatus() {
	hasError := false
	hasWarning := false

	results := []DiagnosticResult{
		d.KeychainAvailable,
		d.WPEngineAPI,
		d.WPEngineSSH,
		d.GitHubToken,
		d.CredentialsFile,
		d.EnvironmentVars,
		d.SSHKeyFile,
	}

	for _, result := range results {
		if result.Status == "error" {
			hasError = true
		} else if result.Status == "warning" {
			hasWarning = true
		}
	}

	if hasError {
		d.OverallStatus = "error"
	} else if hasWarning {
		d.OverallStatus = "warning"
	} else {
		d.OverallStatus = "ok"
	}
}

// generateRecommendations generates actionable recommendations
func (d *CredentialDiagnostics) generateRecommendations() {
	d.RecommendedActions = []string{}

	// Check for critical missing credentials
	if d.WPEngineAPI.Status == "error" {
		d.RecommendedActions = append(d.RecommendedActions,
			"Run 'stax setup wpengine' to configure WPEngine credentials")
	}

	// Check for SSH issues
	if d.SSHKeyFile.Status == "error" || d.SSHKeyFile.Status == "warning" {
		if d.SSHKeyFile.Message == "Default SSH key not found" {
			d.RecommendedActions = append(d.RecommendedActions,
				"Generate SSH key with: ssh-keygen -t rsa -b 4096",
				"Add public key (~/.ssh/id_rsa.pub) to WPEngine User Portal")
		} else if d.SSHKeyFile.Message == "SSH key has insecure permissions" {
			home, _ := os.UserHomeDir()
			sshKeyPath := filepath.Join(home, ".ssh", "id_rsa")
			d.RecommendedActions = append(d.RecommendedActions,
				fmt.Sprintf("Fix SSH key permissions: chmod 600 %s", sshKeyPath))
		}
	}

	// Check for credentials file permission issues
	if d.CredentialsFile.Status == "warning" {
		path, _ := GetCredentialsFilePath()
		d.RecommendedActions = append(d.RecommendedActions,
			fmt.Sprintf("Fix credentials file permissions: chmod 600 %s", path))
	}

	// Check for missing GitHub token (optional but recommended)
	if d.GitHubToken.Status == "warning" && len(d.RecommendedActions) == 0 {
		d.RecommendedActions = append(d.RecommendedActions,
			"Consider adding GitHub token with 'stax setup github' for private repos")
	}

	// If everything is OK, provide confirmation
	if len(d.RecommendedActions) == 0 {
		d.RecommendedActions = append(d.RecommendedActions,
			"All credentials are properly configured!")
	}
}

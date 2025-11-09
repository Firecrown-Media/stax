package wpengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default WPEngine API base URL
	DefaultBaseURL = "https://api.wpengineapi.com/v1"

	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

// Client handles WPEngine API operations
type Client struct {
	baseURL     string
	httpClient  *http.Client
	apiUser     string
	apiPassword string
	install     string
}

// NewClient creates a new WPEngine API client
func NewClient(apiUser, apiPassword, install string) *Client {
	return &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		apiUser:     apiUser,
		apiPassword: apiPassword,
		install:     install,
	}
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// makeRequest makes an HTTP request to the WPEngine API
func (c *Client) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var buf io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		buf = bytes.NewBuffer(data)
	}

	url := c.baseURL + path
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiUser, c.apiPassword)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.httpClient.Do(req)
}

// makeRequestWithRetry makes an HTTP request with retry logic
func (c *Client) makeRequestWithRetry(method, path string, body interface{}, maxAttempts int) (*http.Response, error) {
	var lastErr error
	delay := 1 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := c.makeRequest(method, path, body)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		lastErr = err
		if err == nil {
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
		}

		if attempt < maxAttempts {
			time.Sleep(delay)
			delay *= 2 // Exponential backoff
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// ListInstalls lists all WPEngine installations for the account
func (c *Client) ListInstalls() ([]Install, error) {
	resp, err := c.makeRequestWithRetry("GET", "/installs", nil, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to list installs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result ListInstallsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Results, nil
}

// GetInstall gets detailed information about a specific installation
func (c *Client) GetInstall(installID string) (*InstallDetails, error) {
	resp, err := c.makeRequestWithRetry("GET", fmt.Sprintf("/installs/%s", installID), nil, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to get install: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var details InstallDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &details, nil
}

// GetInstallByName gets installation details by name
func (c *Client) GetInstallByName(name string) (*InstallDetails, error) {
	installs, err := c.ListInstalls()
	if err != nil {
		return nil, err
	}

	for _, install := range installs {
		if install.Name == name {
			return c.GetInstall(install.ID)
		}
	}

	return nil, fmt.Errorf("install %s not found", name)
}

// ListBackups lists available backups for an installation
func (c *Client) ListBackups(installID string) ([]Backup, error) {
	resp, err := c.makeRequestWithRetry("GET", fmt.Sprintf("/installs/%s/backups", installID), nil, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result ListBackupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Results, nil
}

// CreateBackup creates a manual backup
func (c *Client) CreateBackup(installID, description string) (string, error) {
	request := CreateBackupRequest{
		Description: description,
	}

	resp, err := c.makeRequestWithRetry("POST", fmt.Sprintf("/installs/%s/backups", installID), request, 3)
	if err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", c.handleErrorResponse(resp)
	}

	var result CreateBackupResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.ID, nil
}

// GetInstallInfo retrieves information about the configured WPEngine install
func (c *Client) GetInstallInfo() (*InstallDetails, error) {
	if c.install == "" {
		return nil, fmt.Errorf("no install configured")
	}
	return c.GetInstallByName(c.install)
}

// handleErrorResponse handles API error responses
func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d (failed to read error response)", resp.StatusCode)
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("WPEngine API error (%d): %s", resp.StatusCode, errorResp.Message)
}

// handleRateLimit handles API rate limiting
func (c *Client) handleRateLimit(resp *http.Response) error {
	if resp.StatusCode != http.StatusTooManyRequests {
		return nil
	}

	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		return fmt.Errorf("rate limited but no Retry-After header")
	}

	var waitSeconds int
	if _, err := fmt.Sscanf(retryAfter, "%d", &waitSeconds); err != nil {
		return fmt.Errorf("invalid Retry-After header: %s", retryAfter)
	}

	time.Sleep(time.Duration(waitSeconds) * time.Second)
	return nil
}

// TestConnection tests the WPEngine API connection
func (c *Client) TestConnection() error {
	_, err := c.ListInstalls()
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	return nil
}

package mocks

import (
	"fmt"

	"github.com/firecrown-media/stax/pkg/wpengine"
)

// MockWPEngineClient is a mock implementation of WPEngine client
type MockWPEngineClient struct {
	GetInstallFunc       func(install string) (*wpengine.Install, error)
	ListBackupsFunc      func(install string) ([]wpengine.Backup, error)
	CreateBackupFunc     func(install string, description string) (*wpengine.Backup, error)
	DownloadBackupFunc   func(install, backupID, destination string) error
	GetSitesFunc         func(install string) ([]wpengine.Site, error)
	DownloadFilesFunc    func(install, remotePath, localPath string) error
	ExecuteCommandFunc   func(install, command string) (string, error)
	TestConnectionFunc   func() error
}

// NewMockWPEngineClient creates a new mock WPEngine client
func NewMockWPEngineClient() *MockWPEngineClient {
	return &MockWPEngineClient{
		GetInstallFunc: func(install string) (*wpengine.Install, error) {
			return &wpengine.Install{
				Name:        install,
				Environment: "production",
				PHPVersion:  "8.1",
				Status:      "running",
			}, nil
		},
		ListBackupsFunc: func(install string) ([]wpengine.Backup, error) {
			return []wpengine.Backup{
				{
					ID:          "backup-1",
					Description: "Test backup 1",
					Created:     "2024-01-01T12:00:00Z",
					Size:        1024000,
					Type:        "manual",
				},
				{
					ID:          "backup-2",
					Description: "Test backup 2",
					Created:     "2024-01-02T12:00:00Z",
					Size:        2048000,
					Type:        "automatic",
				},
			}, nil
		},
		CreateBackupFunc: func(install string, description string) (*wpengine.Backup, error) {
			return &wpengine.Backup{
				ID:          "backup-new",
				Description: description,
				Created:     "2024-01-03T12:00:00Z",
				Size:        1500000,
				Type:        "manual",
			}, nil
		},
		DownloadBackupFunc: func(install, backupID, destination string) error {
			return nil
		},
		GetSitesFunc: func(install string) ([]wpengine.Site, error) {
			return []wpengine.Site{
				{
					ID:     1,
					Domain: "example.wpengine.com",
					Path:   "/",
				},
				{
					ID:     2,
					Domain: "site1.wpengine.com",
					Path:   "/",
				},
				{
					ID:     3,
					Domain: "site2.wpengine.com",
					Path:   "/",
				},
			}, nil
		},
		DownloadFilesFunc: func(install, remotePath, localPath string) error {
			return nil
		},
		ExecuteCommandFunc: func(install, command string) (string, error) {
			return "Command executed successfully", nil
		},
		TestConnectionFunc: func() error {
			return nil
		},
	}
}

// GetInstall mocks getting install information
func (m *MockWPEngineClient) GetInstall(install string) (*wpengine.Install, error) {
	if m.GetInstallFunc != nil {
		return m.GetInstallFunc(install)
	}
	return nil, fmt.Errorf("GetInstallFunc not implemented")
}

// ListBackups mocks listing backups
func (m *MockWPEngineClient) ListBackups(install string) ([]wpengine.Backup, error) {
	if m.ListBackupsFunc != nil {
		return m.ListBackupsFunc(install)
	}
	return nil, fmt.Errorf("ListBackupsFunc not implemented")
}

// CreateBackup mocks creating a backup
func (m *MockWPEngineClient) CreateBackup(install string, description string) (*wpengine.Backup, error) {
	if m.CreateBackupFunc != nil {
		return m.CreateBackupFunc(install, description)
	}
	return nil, fmt.Errorf("CreateBackupFunc not implemented")
}

// DownloadBackup mocks downloading a backup
func (m *MockWPEngineClient) DownloadBackup(install, backupID, destination string) error {
	if m.DownloadBackupFunc != nil {
		return m.DownloadBackupFunc(install, backupID, destination)
	}
	return fmt.Errorf("DownloadBackupFunc not implemented")
}

// GetSites mocks getting sites
func (m *MockWPEngineClient) GetSites(install string) ([]wpengine.Site, error) {
	if m.GetSitesFunc != nil {
		return m.GetSitesFunc(install)
	}
	return nil, fmt.Errorf("GetSitesFunc not implemented")
}

// DownloadFiles mocks downloading files
func (m *MockWPEngineClient) DownloadFiles(install, remotePath, localPath string) error {
	if m.DownloadFilesFunc != nil {
		return m.DownloadFilesFunc(install, remotePath, localPath)
	}
	return fmt.Errorf("DownloadFilesFunc not implemented")
}

// ExecuteCommand mocks executing a command
func (m *MockWPEngineClient) ExecuteCommand(install, command string) (string, error) {
	if m.ExecuteCommandFunc != nil {
		return m.ExecuteCommandFunc(install, command)
	}
	return "", fmt.Errorf("ExecuteCommandFunc not implemented")
}

// TestConnection mocks testing the connection
func (m *MockWPEngineClient) TestConnection() error {
	if m.TestConnectionFunc != nil {
		return m.TestConnectionFunc()
	}
	return fmt.Errorf("TestConnectionFunc not implemented")
}

// WithError returns a mock client that returns errors
func (m *MockWPEngineClient) WithError(err error) *MockWPEngineClient {
	m.GetInstallFunc = func(install string) (*wpengine.Install, error) {
		return nil, err
	}
	m.ListBackupsFunc = func(install string) ([]wpengine.Backup, error) {
		return nil, err
	}
	m.TestConnectionFunc = func() error {
		return err
	}
	return m
}

// WithEmptyResults returns a mock client with empty results
func (m *MockWPEngineClient) WithEmptyResults() *MockWPEngineClient {
	m.ListBackupsFunc = func(install string) ([]wpengine.Backup, error) {
		return []wpengine.Backup{}, nil
	}
	m.GetSitesFunc = func(install string) ([]wpengine.Site, error) {
		return []wpengine.Site{}, nil
	}
	return m
}

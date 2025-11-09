package mocks

import (
	"fmt"
	"strings"

	"github.com/firecrown-media/stax/pkg/ddev"
)

// MockDDEVManager is a mock implementation of DDEV manager
type MockDDEVManager struct {
	IsRunningFunc func() (bool, error)
	StartFunc     func() error
	StopFunc      func() error
	RestartFunc   func() error
	DescribeFunc  func() (*ddev.ProjectInfo, error)
	ExecFunc      func(command string, args ...string) (string, error)
	ImportDBFunc  func(dbPath string) error
	ExportDBFunc  func(destination string) error
	GetLogsFunc   func(service string) (string, error)
	running       bool
	lastCommand   string
	lastArgs      []string
}

// NewMockDDEVManager creates a new mock DDEV manager
func NewMockDDEVManager() *MockDDEVManager {
	return &MockDDEVManager{
		running: false,
		IsRunningFunc: func() (bool, error) {
			return false, nil
		},
		StartFunc: func() error {
			return nil
		},
		StopFunc: func() error {
			return nil
		},
		RestartFunc: func() error {
			return nil
		},
		DescribeFunc: func() (*ddev.ProjectInfo, error) {
			return &ddev.ProjectInfo{
				Name:            "test-project",
				Status:          "running",
				PHPVersion:      "8.1",
				DatabaseType:    "mysql",
				DatabaseVersion: "8.0",
				URLs:            []string{"https://test-project.ddev.site"},
			}, nil
		},
		ExecFunc: func(command string, args ...string) (string, error) {
			return fmt.Sprintf("Executed: %s %s", command, strings.Join(args, " ")), nil
		},
		ImportDBFunc: func(dbPath string) error {
			return nil
		},
		ExportDBFunc: func(destination string) error {
			return nil
		},
		GetLogsFunc: func(service string) (string, error) {
			return fmt.Sprintf("Logs for %s", service), nil
		},
	}
}

// IsRunning mocks checking if DDEV is running
func (m *MockDDEVManager) IsRunning() (bool, error) {
	if m.IsRunningFunc != nil {
		return m.IsRunningFunc()
	}
	return m.running, nil
}

// Start mocks starting DDEV
func (m *MockDDEVManager) Start() error {
	if m.StartFunc != nil {
		err := m.StartFunc()
		if err == nil {
			m.running = true
		}
		return err
	}
	m.running = true
	return nil
}

// Stop mocks stopping DDEV
func (m *MockDDEVManager) Stop() error {
	if m.StopFunc != nil {
		err := m.StopFunc()
		if err == nil {
			m.running = false
		}
		return err
	}
	m.running = false
	return nil
}

// Restart mocks restarting DDEV
func (m *MockDDEVManager) Restart() error {
	if m.RestartFunc != nil {
		return m.RestartFunc()
	}
	m.running = true
	return nil
}

// Describe mocks getting project information
func (m *MockDDEVManager) Describe() (*ddev.ProjectInfo, error) {
	if m.DescribeFunc != nil {
		return m.DescribeFunc()
	}
	return nil, fmt.Errorf("DescribeFunc not implemented")
}

// Exec mocks executing a command in DDEV
func (m *MockDDEVManager) Exec(command string, args ...string) (string, error) {
	m.lastCommand = command
	m.lastArgs = args

	if m.ExecFunc != nil {
		return m.ExecFunc(command, args...)
	}
	return "", fmt.Errorf("ExecFunc not implemented")
}

// ImportDB mocks importing a database
func (m *MockDDEVManager) ImportDB(dbPath string) error {
	if m.ImportDBFunc != nil {
		return m.ImportDBFunc(dbPath)
	}
	return fmt.Errorf("ImportDBFunc not implemented")
}

// ExportDB mocks exporting a database
func (m *MockDDEVManager) ExportDB(destination string) error {
	if m.ExportDBFunc != nil {
		return m.ExportDBFunc(destination)
	}
	return fmt.Errorf("ExportDBFunc not implemented")
}

// GetLogs mocks getting logs
func (m *MockDDEVManager) GetLogs(service string) (string, error) {
	if m.GetLogsFunc != nil {
		return m.GetLogsFunc(service)
	}
	return "", fmt.Errorf("GetLogsFunc not implemented")
}

// GetLastCommand returns the last executed command
func (m *MockDDEVManager) GetLastCommand() (string, []string) {
	return m.lastCommand, m.lastArgs
}

// WithRunningState sets the running state
func (m *MockDDEVManager) WithRunningState(running bool) *MockDDEVManager {
	m.running = running
	m.IsRunningFunc = func() (bool, error) {
		return running, nil
	}
	return m
}

// WithStartError returns a mock that fails to start
func (m *MockDDEVManager) WithStartError(err error) *MockDDEVManager {
	m.StartFunc = func() error {
		return err
	}
	return m
}

// WithStopError returns a mock that fails to stop
func (m *MockDDEVManager) WithStopError(err error) *MockDDEVManager {
	m.StopFunc = func() error {
		return err
	}
	return m
}

// WithImportDBError returns a mock that fails to import database
func (m *MockDDEVManager) WithImportDBError(err error) *MockDDEVManager {
	m.ImportDBFunc = func(dbPath string) error {
		return err
	}
	return m
}

// MockDDEVConfig is a mock DDEV configuration
type MockDDEVConfig struct {
	ProjectName     string
	PHPVersion      string
	MySQLVersion    string
	WebserverType   string
	AdditionalFQDNs []string
}

// NewMockDDEVConfig creates a new mock DDEV config
func NewMockDDEVConfig() *MockDDEVConfig {
	return &MockDDEVConfig{
		ProjectName:     "test-project",
		PHPVersion:      "8.1",
		MySQLVersion:    "8.0",
		WebserverType:   "nginx-fpm",
		AdditionalFQDNs: []string{"test.local", "site1.test.local"},
	}
}

// ToYAML converts the mock config to YAML
func (m *MockDDEVConfig) ToYAML() string {
	return fmt.Sprintf(`name: %s
type: php
docroot: ""
php_version: "%s"
mysql_version: "%s"
webserver_type: %s
additional_fqdns:
%s`, m.ProjectName, m.PHPVersion, m.MySQLVersion, m.WebserverType, formatFQDNs(m.AdditionalFQDNs))
}

func formatFQDNs(fqdns []string) string {
	var result string
	for _, fqdn := range fqdns {
		result += fmt.Sprintf("  - %s\n", fqdn)
	}
	return result
}

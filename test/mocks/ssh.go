package mocks

import (
	"fmt"
	"io"
	"os"
)

// MockSSHClient is a mock implementation of an SSH client
type MockSSHClient struct {
	ConnectFunc      func(host, user, keyPath string) error
	DisconnectFunc   func() error
	ExecuteFunc      func(command string) (string, error)
	DownloadFileFunc func(remotePath, localPath string) error
	UploadFileFunc   func(localPath, remotePath string) error
	responses        map[string]string
	downloadHandler  func(remotePath, localPath string) error
}

// NewMockSSHClient creates a new mock SSH client
func NewMockSSHClient() *MockSSHClient {
	return &MockSSHClient{
		responses: make(map[string]string),
		ConnectFunc: func(host, user, keyPath string) error {
			return nil
		},
		DisconnectFunc: func() error {
			return nil
		},
		ExecuteFunc: func(command string) (string, error) {
			return "Command executed", nil
		},
		DownloadFileFunc: func(remotePath, localPath string) error {
			// Create an empty file
			return os.WriteFile(localPath, []byte("mock file content"), 0644)
		},
		UploadFileFunc: func(localPath, remotePath string) error {
			return nil
		},
	}
}

// Connect mocks SSH connection
func (m *MockSSHClient) Connect(host, user, keyPath string) error {
	if m.ConnectFunc != nil {
		return m.ConnectFunc(host, user, keyPath)
	}
	return nil
}

// Disconnect mocks SSH disconnection
func (m *MockSSHClient) Disconnect() error {
	if m.DisconnectFunc != nil {
		return m.DisconnectFunc()
	}
	return nil
}

// Execute mocks command execution
func (m *MockSSHClient) Execute(command string) (string, error) {
	// Check if we have a specific response for this command
	if response, ok := m.responses[command]; ok {
		return response, nil
	}

	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(command)
	}
	return "", fmt.Errorf("ExecuteFunc not implemented")
}

// DownloadFile mocks file download
func (m *MockSSHClient) DownloadFile(remotePath, localPath string) error {
	if m.downloadHandler != nil {
		return m.downloadHandler(remotePath, localPath)
	}

	if m.DownloadFileFunc != nil {
		return m.DownloadFileFunc(remotePath, localPath)
	}
	return fmt.Errorf("DownloadFileFunc not implemented")
}

// UploadFile mocks file upload
func (m *MockSSHClient) UploadFile(localPath, remotePath string) error {
	if m.UploadFileFunc != nil {
		return m.UploadFileFunc(localPath, remotePath)
	}
	return fmt.Errorf("UploadFileFunc not implemented")
}

// AddCommandResponse adds a specific response for a command
func (m *MockSSHClient) AddCommandResponse(command, response string) {
	m.responses[command] = response
}

// WithConnectionError returns a mock that fails to connect
func (m *MockSSHClient) WithConnectionError(err error) *MockSSHClient {
	m.ConnectFunc = func(host, user, keyPath string) error {
		return err
	}
	return m
}

// WithExecutionError returns a mock that fails to execute commands
func (m *MockSSHClient) WithExecutionError(err error) *MockSSHClient {
	m.ExecuteFunc = func(command string) (string, error) {
		return "", err
	}
	return m
}

// WithDownloadHandler sets a custom download handler
func (m *MockSSHClient) WithDownloadHandler(handler func(remotePath, localPath string) error) *MockSSHClient {
	m.downloadHandler = handler
	return m
}

// MockSCPClient is a mock implementation of an SCP client
type MockSCPClient struct {
	CopyFromRemoteFunc func(remotePath, localPath string) error
	CopyToRemoteFunc   func(localPath, remotePath string) error
}

// NewMockSCPClient creates a new mock SCP client
func NewMockSCPClient() *MockSCPClient {
	return &MockSCPClient{
		CopyFromRemoteFunc: func(remotePath, localPath string) error {
			return os.WriteFile(localPath, []byte("mock scp content"), 0644)
		},
		CopyToRemoteFunc: func(localPath, remotePath string) error {
			return nil
		},
	}
}

// CopyFromRemote mocks copying a file from remote
func (m *MockSCPClient) CopyFromRemote(remotePath, localPath string) error {
	if m.CopyFromRemoteFunc != nil {
		return m.CopyFromRemoteFunc(remotePath, localPath)
	}
	return fmt.Errorf("CopyFromRemoteFunc not implemented")
}

// CopyToRemote mocks copying a file to remote
func (m *MockSCPClient) CopyToRemote(localPath, remotePath string) error {
	if m.CopyToRemoteFunc != nil {
		return m.CopyToRemoteFunc(localPath, remotePath)
	}
	return fmt.Errorf("CopyToRemoteFunc not implemented")
}

// MockRsyncClient is a mock implementation of an rsync client
type MockRsyncClient struct {
	SyncFunc func(source, destination string, options []string) error
}

// NewMockRsyncClient creates a new mock rsync client
func NewMockRsyncClient() *MockRsyncClient {
	return &MockRsyncClient{
		SyncFunc: func(source, destination string, options []string) error {
			return nil
		},
	}
}

// Sync mocks rsync synchronization
func (m *MockRsyncClient) Sync(source, destination string, options []string) error {
	if m.SyncFunc != nil {
		return m.SyncFunc(source, destination, options)
	}
	return fmt.Errorf("SyncFunc not implemented")
}

// WithError returns a mock that returns an error
func (m *MockRsyncClient) WithError(err error) *MockRsyncClient {
	m.SyncFunc = func(source, destination string, options []string) error {
		return err
	}
	return m
}

// MockFileTransfer is a helper for testing file transfers
type MockFileTransfer struct {
	files map[string][]byte
}

// NewMockFileTransfer creates a new mock file transfer
func NewMockFileTransfer() *MockFileTransfer {
	return &MockFileTransfer{
		files: make(map[string][]byte),
	}
}

// AddFile adds a file to the mock transfer
func (m *MockFileTransfer) AddFile(path string, content []byte) {
	m.files[path] = content
}

// GetFile retrieves a file from the mock transfer
func (m *MockFileTransfer) GetFile(path string) ([]byte, error) {
	content, ok := m.files[path]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return content, nil
}

// WriteToFile writes a mock file to disk
func (m *MockFileTransfer) WriteToFile(remotePath, localPath string) error {
	content, ok := m.files[remotePath]
	if !ok {
		return fmt.Errorf("file not found: %s", remotePath)
	}
	return os.WriteFile(localPath, content, 0644)
}

// ReadFromFile reads a file from disk into the mock transfer
func (m *MockFileTransfer) ReadFromFile(localPath, remotePath string) error {
	content, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}
	m.files[remotePath] = content
	return nil
}

// MockProgressReader is a mock io.Reader that tracks progress
type MockProgressReader struct {
	reader       io.Reader
	total        int64
	current      int64
	onProgress   func(current, total int64)
}

// NewMockProgressReader creates a new mock progress reader
func NewMockProgressReader(reader io.Reader, total int64, onProgress func(current, total int64)) *MockProgressReader {
	return &MockProgressReader{
		reader:     reader,
		total:      total,
		onProgress: onProgress,
	}
}

// Read implements io.Reader
func (m *MockProgressReader) Read(p []byte) (int, error) {
	n, err := m.reader.Read(p)
	m.current += int64(n)
	if m.onProgress != nil {
		m.onProgress(m.current, m.total)
	}
	return n, err
}

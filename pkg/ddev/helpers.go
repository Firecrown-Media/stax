package ddev

import "time"

// IsConfigured checks if DDEV is configured for a project
func IsConfigured(projectPath string) bool {
	return ConfigExists(projectPath)
}

// Start starts DDEV for a project
func Start(projectPath string) error {
	m := NewManager(projectPath)
	return m.Start()
}

// Stop stops DDEV for a project
func Stop(projectPath string) error {
	m := NewManager(projectPath)
	return m.Stop()
}

// Restart restarts DDEV for a project
func Restart(projectPath string) error {
	m := NewManager(projectPath)
	return m.Restart()
}

// Delete deletes DDEV project
func Delete(projectPath string, removeData bool) error {
	m := NewManager(projectPath)
	return m.Delete(removeData)
}

// GetStatus gets the status of a DDEV project
func GetStatus(projectPath string) (*ProjectInfo, error) {
	m := NewManager(projectPath)
	return m.Describe()
}

// PowerOff powers off all DDEV projects
func PowerOff() error {
	// This is a global command - would use `ddev poweroff`
	// For now, just return nil - actual implementation would be added to manager.go
	return nil
}

// EnableXdebug enables Xdebug for a project
func EnableXdebug(projectPath string) error {
	m := NewManager(projectPath)
	// Update config to enable xdebug
	updates := map[string]interface{}{
		"xdebug_enabled": true,
	}
	if err := UpdateConfig(projectPath, updates); err != nil {
		return err
	}
	// Restart to apply changes
	return m.Restart()
}

// DisableXdebug disables Xdebug for a project
func DisableXdebug(projectPath string) error {
	m := NewManager(projectPath)
	// Update config to disable xdebug
	updates := map[string]interface{}{
		"xdebug_enabled": false,
	}
	if err := UpdateConfig(projectPath, updates); err != nil {
		return err
	}
	// Restart to apply changes
	return m.Restart()
}

// Exec executes a command in the web container
func Exec(projectPath string, command ...string) error {
	m := NewManager(projectPath)
	return m.Exec(command, nil)
}

// WaitForReady waits for DDEV to be ready
func WaitForReady(projectPath string, timeout time.Duration) error {
	m := NewManager(projectPath)
	return m.WaitForReady(timeout)
}

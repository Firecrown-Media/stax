package ddev

import "os/exec"

// IsConfigured checks if DDEV is configured for a project
func IsConfigured(projectPath string) bool {
	return ConfigExists(projectPath)
}

// Start starts the DDEV project
func Start(projectPath string) error {
	mgr := NewManager(projectPath)
	return mgr.Start()
}

// Stop stops the DDEV project
func Stop(projectPath string) error {
	mgr := NewManager(projectPath)
	return mgr.Stop()
}

// Restart restarts the DDEV project
func Restart(projectPath string) error {
	mgr := NewManager(projectPath)
	return mgr.Restart()
}

// GetStatus gets the DDEV project status
func GetStatus(projectPath string) (*ProjectInfo, error) {
	mgr := NewManager(projectPath)
	return mgr.Describe()
}

// EnableXdebug enables Xdebug
func EnableXdebug(projectPath string) error {
	mgr := NewManager(projectPath)
	return mgr.ExecCommand("xdebug", "on")
}

// DisableXdebug disables Xdebug
func DisableXdebug(projectPath string) error {
	mgr := NewManager(projectPath)
	return mgr.ExecCommand("xdebug", "off")
}

// Exec executes a command in the DDEV container
func Exec(projectPath string, args ...string) error {
	mgr := NewManager(projectPath)
	return mgr.Exec(args, nil)
}

// Delete deletes the DDEV project
func Delete(projectPath string, removeData bool) error {
	mgr := NewManager(projectPath)
	return mgr.Delete(removeData)
}

// PowerOff stops all DDEV projects
func PowerOff() error {
	cmd := exec.Command("ddev", "poweroff")
	return cmd.Run()
}

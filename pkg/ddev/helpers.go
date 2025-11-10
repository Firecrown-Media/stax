package ddev

// IsConfigured checks if DDEV is configured for a project
func IsConfigured(projectPath string) bool {
	return ConfigExists(projectPath)
}

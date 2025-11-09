package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// StatusChecker checks build status and determines if rebuild is needed
type StatusChecker struct {
	projectPath string
}

// NewStatusChecker creates a new status checker
func NewStatusChecker(projectPath string) *StatusChecker {
	return &StatusChecker{
		projectPath: projectPath,
	}
}

// GetStatus returns the current build status
func (s *StatusChecker) GetStatus() (*BuildStatus, error) {
	status := &BuildStatus{
		NeedsBuild:         false,
		Reasons:            []string{},
		CustomBuildScripts: []string{},
	}

	// Check for build script
	buildScriptPath := filepath.Join(s.projectPath, "scripts", "build.sh")
	if _, err := os.Stat(buildScriptPath); err == nil {
		status.BuildScriptExists = true
	} else {
		status.Reasons = append(status.Reasons, "No build script found")
	}

	// Detect custom build scripts
	buildDir := filepath.Join(s.projectPath, "scripts", "build")
	if entries, err := os.ReadDir(buildDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sh") {
				status.CustomBuildScripts = append(status.CustomBuildScripts, entry.Name())
			}
		}
	}

	// Check composer status
	composerStatus, err := s.getComposerStatus()
	if err == nil {
		status.ComposerStatus = *composerStatus
		if !composerStatus.Installed {
			status.NeedsBuild = true
			status.Reasons = append(status.Reasons, "Composer dependencies not installed")
		} else if composerStatus.NeedsUpdate {
			status.NeedsBuild = true
			status.Reasons = append(status.Reasons, "Composer dependencies need updating")
		}
	}

	// Check npm status
	npmStatus, err := s.getNPMStatus()
	if err == nil {
		status.NPMStatus = *npmStatus
		if !npmStatus.Installed {
			status.NeedsBuild = true
			status.Reasons = append(status.Reasons, "NPM dependencies not installed")
		} else if npmStatus.NeedsUpdate {
			status.NeedsBuild = true
			status.Reasons = append(status.Reasons, "NPM dependencies need updating")
		}
	}

	// Check build artifacts
	if !s.buildArtifactsExist() {
		status.NeedsBuild = true
		status.Reasons = append(status.Reasons, "Build artifacts missing")
	}

	// Get last build time
	lastBuildTime, err := s.getLastBuildTime()
	if err == nil {
		status.LastBuildTime = lastBuildTime
	}

	// Check if source files are newer than build
	if s.sourceFilesNewer(lastBuildTime) {
		status.NeedsBuild = true
		status.Reasons = append(status.Reasons, "Source files modified since last build")
	}

	return status, nil
}

// NeedsBuild checks if a build is needed
func (s *StatusChecker) NeedsBuild() (bool, []string) {
	status, err := s.GetStatus()
	if err != nil {
		return true, []string{"Error checking build status"}
	}

	return status.NeedsBuild, status.Reasons
}

// GetLastBuildTime returns when the last build occurred
func (s *StatusChecker) GetLastBuildTime() (time.Time, error) {
	return s.getLastBuildTime()
}

// getLastBuildTime finds the most recent build artifact modification time
func (s *StatusChecker) getLastBuildTime() (time.Time, error) {
	var latestTime time.Time

	// Check various build artifacts
	artifactPaths := []string{
		filepath.Join(s.projectPath, "wp-content", "mu-plugins", "firecrown", "vendor"),
		filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-parent", "build"),
		filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-child", "build"),
	}

	for _, path := range artifactPaths {
		if info, err := os.Stat(path); err == nil {
			modTime := info.ModTime()
			if modTime.After(latestTime) {
				latestTime = modTime
			}
		}
	}

	if latestTime.IsZero() {
		return latestTime, fmt.Errorf("no build artifacts found")
	}

	return latestTime, nil
}

// getComposerStatus checks composer dependencies status
func (s *StatusChecker) getComposerStatus() (*DependencyStatus, error) {
	// Check in MU plugins directory
	muPluginPath := filepath.Join(s.projectPath, "wp-content", "mu-plugins", "firecrown")
	composer := NewComposer(muPluginPath)
	return composer.GetStatus()
}

// getNPMStatus checks npm dependencies status
func (s *StatusChecker) getNPMStatus() (*DependencyStatus, error) {
	// Check in parent theme directory
	themePath := filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-parent")
	npm := NewNPM(themePath)
	return npm.GetStatus()
}

// buildArtifactsExist checks if build artifacts exist
func (s *StatusChecker) buildArtifactsExist() bool {
	requiredArtifacts := []string{
		filepath.Join(s.projectPath, "wp-content", "mu-plugins", "firecrown", "vendor"),
		filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-parent", "build"),
		filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-parent", "node_modules"),
	}

	for _, artifact := range requiredArtifacts {
		if _, err := os.Stat(artifact); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// sourceFilesNewer checks if source files are newer than build artifacts
func (s *StatusChecker) sourceFilesNewer(buildTime time.Time) bool {
	if buildTime.IsZero() {
		return true
	}

	// Check source files that should trigger rebuild
	sourceFiles := []string{
		filepath.Join(s.projectPath, "wp-content", "mu-plugins", "firecrown", "composer.json"),
		filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-parent", "package.json"),
		filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-parent", "src"),
		filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-child", "package.json"),
	}

	for _, source := range sourceFiles {
		if info, err := os.Stat(source); err == nil {
			if info.ModTime().After(buildTime) {
				return true
			}
		}
	}

	return false
}

// CompareTimestamps compares modification times
func (s *StatusChecker) CompareTimestamps(source, dest string) (bool, error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return false, err
	}

	destInfo, err := os.Stat(dest)
	if err != nil {
		return true, nil // Dest doesn't exist, needs build
	}

	return sourceInfo.ModTime().After(destInfo.ModTime()), nil
}

// HasChanges checks if there are uncommitted git changes
func (s *StatusChecker) HasChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = s.projectPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// If output is not empty, there are changes
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// HasChangesSince checks if there are git changes since a specific time
func (s *StatusChecker) HasChangesSince(since time.Time) (bool, error) {
	// Get commits since the timestamp
	sinceStr := since.Format("2006-01-02 15:04:05")
	cmd := exec.Command("git", "log", "--since", sinceStr, "--oneline")
	cmd.Dir = s.projectPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

// ValidateBuild verifies that build output exists and is valid
func (s *StatusChecker) ValidateBuild() error {
	errors := []string{}

	// Check vendor directory
	vendorPath := filepath.Join(s.projectPath, "wp-content", "mu-plugins", "firecrown", "vendor")
	if _, err := os.Stat(filepath.Join(vendorPath, "autoload.php")); os.IsNotExist(err) {
		errors = append(errors, "MU plugin vendor/autoload.php missing")
	}

	// Check parent theme build
	parentBuildPath := filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-parent", "build")
	if _, err := os.Stat(parentBuildPath); os.IsNotExist(err) {
		errors = append(errors, "Parent theme build directory missing")
	}

	// Check for compiled CSS/JS in parent theme
	parentCSSPath := filepath.Join(parentBuildPath, "scripts.css")
	if _, err := os.Stat(parentCSSPath); os.IsNotExist(err) {
		errors = append(errors, "Parent theme compiled CSS missing")
	}

	// Check child theme build
	childBuildPath := filepath.Join(s.projectPath, "wp-content", "themes", "firecrown-child", "build")
	if _, err := os.Stat(childBuildPath); os.IsNotExist(err) {
		errors = append(errors, "Child theme build directory missing")
	}

	if len(errors) > 0 {
		return fmt.Errorf("build validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// GetBuildDuration estimates how long a build will take
func (s *StatusChecker) GetBuildDuration() time.Duration {
	// This is an estimate based on typical build times
	// Could be made more sophisticated by tracking actual build times
	return 3 * time.Minute
}

// IsBuildRunning checks if a build is currently running
func (s *StatusChecker) IsBuildRunning() bool {
	// Check for lock file
	lockFile := filepath.Join(s.projectPath, ".stax", "build.lock")
	if _, err := os.Stat(lockFile); err == nil {
		// Lock file exists, but check if process is still running
		// For now, just return true
		return true
	}

	return false
}

// CreateBuildLock creates a lock file to prevent concurrent builds
func (s *StatusChecker) CreateBuildLock() error {
	lockDir := filepath.Join(s.projectPath, ".stax")
	if err := os.MkdirAll(lockDir, 0755); err != nil {
		return err
	}

	lockFile := filepath.Join(lockDir, "build.lock")
	pidContent := fmt.Sprintf("%d", os.Getpid())
	return os.WriteFile(lockFile, []byte(pidContent), 0644)
}

// RemoveBuildLock removes the build lock file
func (s *StatusChecker) RemoveBuildLock() error {
	lockFile := filepath.Join(s.projectPath, ".stax", "build.lock")
	return os.Remove(lockFile)
}

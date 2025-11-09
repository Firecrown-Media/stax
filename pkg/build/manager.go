package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Manager handles build operations
type Manager struct {
	projectPath string
	verbose     bool
}

// NewManager creates a new build manager
func NewManager(projectPath string) *Manager {
	return &Manager{
		projectPath: projectPath,
		verbose:     false,
	}
}

// SetVerbose enables/disables verbose output
func (m *Manager) SetVerbose(verbose bool) {
	m.verbose = verbose
}

// RunBuildScript executes the main build script (scripts/build.sh)
func (m *Manager) RunBuildScript() (*BuildResult, error) {
	buildScriptPath := filepath.Join(m.projectPath, "scripts", "build.sh")

	if _, err := os.Stat(buildScriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("build script not found: %s", buildScriptPath)
	}

	startTime := time.Now()

	cmd := exec.Command("bash", buildScriptPath)
	cmd.Dir = m.projectPath
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	result := &BuildResult{
		Success:  err == nil,
		Duration: duration,
		Output:   string(output),
		Error:    err,
		Steps: []BuildStep{
			{
				Name:     "scripts/build.sh",
				Command:  "bash scripts/build.sh",
				Success:  err == nil,
				Duration: duration,
				Output:   string(output),
				Error:    err,
			},
		},
	}

	return result, err
}

// RunComposerInstall runs composer install in a specific directory
func (m *Manager) RunComposerInstall(path string, options ComposerOptions) error {
	composer := NewComposer(path)
	return composer.Install(options)
}

// RunNPMInstall runs npm install in a specific directory
func (m *Manager) RunNPMInstall(path string, options NPMOptions) error {
	npm := NewNPM(path)
	return npm.Install(options)
}

// RunNPMBuild runs npm run build in a specific directory
func (m *Manager) RunNPMBuild(path string, options NPMOptions) error {
	npm := NewNPM(path)
	return npm.Build(options)
}

// RunNPMStart runs npm start in a specific directory
func (m *Manager) RunNPMStart(path string, background bool, options NPMOptions) error {
	npm := NewNPM(path)
	return npm.Start(background, options)
}

// RunPHPCS runs PHPCS linting
func (m *Manager) RunPHPCS(options PHPCSOptions) (*PHPCSResult, error) {
	quality := NewQuality(m.projectPath)
	return quality.RunPHPCS(options)
}

// RunPHPCBF runs PHPCS auto-fixer
func (m *Manager) RunPHPCBF(options PHPCSOptions) error {
	quality := NewQuality(m.projectPath)
	return quality.RunPHPCBF(options)
}

// DetectBuildScripts finds all build scripts in the project
func (m *Manager) DetectBuildScripts() ([]ScriptInfo, error) {
	scripts := []ScriptInfo{}

	// Check for main build script
	mainBuildScript := filepath.Join(m.projectPath, "scripts", "build.sh")
	if _, err := os.Stat(mainBuildScript); err == nil {
		scripts = append(scripts, ScriptInfo{
			Name:        "build.sh",
			Path:        mainBuildScript,
			Type:        "build",
			Description: "Main build script",
			Order:       0,
		})
	}

	// Check for individual build scripts
	buildDir := filepath.Join(m.projectPath, "scripts", "build")
	if _, err := os.Stat(buildDir); err == nil {
		entries, err := os.ReadDir(buildDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sh") {
					scriptPath := filepath.Join(buildDir, entry.Name())
					scriptType := "custom"
					description := "Custom build script"
					order := 999

					// Parse order from filename (e.g., 10-mu-plugins.sh)
					if len(entry.Name()) >= 3 && entry.Name()[0] >= '0' && entry.Name()[0] <= '9' {
						order = int(entry.Name()[0]-'0') * 10
						if entry.Name()[1] >= '0' && entry.Name()[1] <= '9' {
							order += int(entry.Name()[1] - '0')
						}
					}

					// Determine type and description from filename
					name := entry.Name()
					if strings.Contains(name, "mu-plugin") {
						scriptType = "composer"
						description = "Build MU plugins (composer install)"
					} else if strings.Contains(name, "theme") {
						scriptType = "npm"
						description = "Build themes (npm install & build)"
					} else if strings.Contains(name, "plugin") {
						scriptType = "composer"
						description = "Build plugins (composer install)"
					}

					scripts = append(scripts, ScriptInfo{
						Name:        entry.Name(),
						Path:        scriptPath,
						Type:        scriptType,
						Description: description,
						Order:       order,
					})
				}
			}
		}
	}

	// Sort by order
	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].Order < scripts[j].Order
	})

	return scripts, nil
}

// GetBuildStatus checks if a build is needed
func (m *Manager) GetBuildStatus() (*BuildStatus, error) {
	checker := NewStatusChecker(m.projectPath)
	return checker.GetStatus()
}

// Clean removes build artifacts
func (m *Manager) Clean() error {
	dirsToClean := []string{
		filepath.Join(m.projectPath, "wp-content", "mu-plugins", "firecrown", "vendor"),
		filepath.Join(m.projectPath, "wp-content", "themes", "firecrown-parent", "node_modules"),
		filepath.Join(m.projectPath, "wp-content", "themes", "firecrown-parent", "build"),
		filepath.Join(m.projectPath, "wp-content", "themes", "firecrown-child", "node_modules"),
		filepath.Join(m.projectPath, "wp-content", "themes", "firecrown-child", "build"),
	}

	for _, dir := range dirsToClean {
		if _, err := os.Stat(dir); err == nil {
			if m.verbose {
				fmt.Printf("Removing %s\n", dir)
			}
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("failed to remove %s: %w", dir, err)
			}
		}
	}

	return nil
}

// BuildMUPlugins builds MU plugins (composer install)
func (m *Manager) BuildMUPlugins(options ComposerOptions) error {
	muPluginPath := filepath.Join(m.projectPath, "wp-content", "mu-plugins", "firecrown")

	if _, err := os.Stat(filepath.Join(muPluginPath, "composer.json")); os.IsNotExist(err) {
		return fmt.Errorf("composer.json not found in %s", muPluginPath)
	}

	composer := NewComposer(muPluginPath)
	return composer.Install(options)
}

// BuildTheme builds a theme (npm install && npm run build)
func (m *Manager) BuildTheme(themeName string, options NPMOptions) error {
	themePath := filepath.Join(m.projectPath, "wp-content", "themes", themeName)

	if _, err := os.Stat(filepath.Join(themePath, "package.json")); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in %s", themePath)
	}

	npm := NewNPM(themePath)

	// Run npm install
	if err := npm.Install(options); err != nil {
		return fmt.Errorf("npm install failed: %w", err)
	}

	// Run npm build
	if err := npm.Build(options); err != nil {
		return fmt.Errorf("npm build failed: %w", err)
	}

	// Also run composer install if composer.json exists
	if _, err := os.Stat(filepath.Join(themePath, "composer.json")); err == nil {
		composer := NewComposer(themePath)
		composerOpts := ComposerOptions{
			WorkingDir:         themePath,
			NoDev:              true,
			IgnorePlatformReqs: true,
			PreferDist:         true,
			Timeout:            options.Timeout,
			Verbose:            options.Verbose,
		}
		if err := composer.Install(composerOpts); err != nil {
			return fmt.Errorf("composer install failed: %w", err)
		}
	}

	return nil
}

// BuildThemes builds all themes
func (m *Manager) BuildThemes(options NPMOptions) error {
	themes := []string{"firecrown-parent", "firecrown-child"}

	for _, theme := range themes {
		themePath := filepath.Join(m.projectPath, "wp-content", "themes", theme)
		if _, err := os.Stat(filepath.Join(themePath, "package.json")); err == nil {
			if err := m.BuildTheme(theme, options); err != nil {
				return fmt.Errorf("failed to build theme %s: %w", theme, err)
			}
		}
	}

	return nil
}

// WatchForChanges watches files for changes and triggers rebuilds
func (m *Manager) WatchForChanges(callback func()) error {
	watcher := NewWatcher(m.projectPath)
	return watcher.Watch(callback)
}

// GenerateBuildScript creates a build script if one doesn't exist
func (m *Manager) GenerateBuildScript() error {
	scriptsDir := filepath.Join(m.projectPath, "scripts")
	buildDir := filepath.Join(scriptsDir, "build")

	// Create directories
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create scripts directory: %w", err)
	}

	// Create main build script
	mainBuildScript := filepath.Join(scriptsDir, "build.sh")
	if _, err := os.Stat(mainBuildScript); os.IsNotExist(err) {
		content := `#!/bin/bash

set -e

echo "Running build scripts..."

script_folder="scripts/build"

# Check if the script subfolder exists
if [ -d "$script_folder" ]; then
    # Loop through each script file in the subfolder
    for script_file in "$script_folder"/*.sh; do
        # Check if the file is a regular file
        if [ -f "$script_file" ]; then
            # Run the script file
            bash "$script_file"
        fi
    done
else
    echo "$script_folder does not exist"
fi
`
		if err := os.WriteFile(mainBuildScript, []byte(content), 0755); err != nil {
			return fmt.Errorf("failed to create build script: %w", err)
		}
	}

	// Create MU plugins build script
	muPluginsScript := filepath.Join(buildDir, "10-mu-plugins.sh")
	if _, err := os.Stat(muPluginsScript); os.IsNotExist(err) {
		content := `#!/bin/bash

set -e

echo "Running mu-plugins scripts..."
cd wp-content/mu-plugins/firecrown
composer install --ignore-platform-reqs
cd -
`
		if err := os.WriteFile(muPluginsScript, []byte(content), 0755); err != nil {
			return fmt.Errorf("failed to create mu-plugins build script: %w", err)
		}
	}

	// Create theme build script
	themeScript := filepath.Join(buildDir, "20-theme.sh")
	if _, err := os.Stat(themeScript); os.IsNotExist(err) {
		content := `#!/bin/bash

set -e

echo "Running theme scripts..."

cd wp-content/themes/firecrown-parent
npm install
npm run build
composer install --ignore-platform-reqs
cd -

cd wp-content/themes/firecrown-child
npm install
npm run build
cd -
`
		if err := os.WriteFile(themeScript, []byte(content), 0755); err != nil {
			return fmt.Errorf("failed to create theme build script: %w", err)
		}
	}

	return nil
}

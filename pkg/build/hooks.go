package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Hooks handles git hooks integration (Husky)
type Hooks struct {
	projectPath string
}

// NewHooks creates a new hooks manager
func NewHooks(projectPath string) *Hooks {
	return &Hooks{
		projectPath: projectPath,
	}
}

// SetupHusky initializes Husky hooks
func (h *Hooks) SetupHusky() error {
	// Check if Husky is already configured
	huskyDir := filepath.Join(h.projectPath, ".husky")
	if _, err := os.Stat(huskyDir); err == nil {
		return fmt.Errorf("Husky is already configured")
	}

	// Create .husky directory
	if err := os.MkdirAll(huskyDir, 0755); err != nil {
		return fmt.Errorf("failed to create .husky directory: %w", err)
	}

	// Create pre-commit hook
	preCommitPath := filepath.Join(huskyDir, "pre-commit")
	preCommitContent := h.GetPreCommitHook()
	if err := os.WriteFile(preCommitPath, []byte(preCommitContent), 0755); err != nil {
		return fmt.Errorf("failed to create pre-commit hook: %w", err)
	}

	return nil
}

// VerifyHusky checks if Husky is configured
func (h *Hooks) VerifyHusky() (*HuskyConfig, error) {
	config := &HuskyConfig{
		Enabled: false,
	}

	// Check for .husky directory
	huskyDir := filepath.Join(h.projectPath, ".husky")
	if _, err := os.Stat(huskyDir); os.IsNotExist(err) {
		return config, fmt.Errorf("Husky not configured (.husky directory not found)")
	}

	config.Enabled = true
	config.ConfigFile = huskyDir

	// Check for pre-commit hook
	preCommitPath := filepath.Join(huskyDir, "pre-commit")
	if data, err := os.ReadFile(preCommitPath); err == nil {
		config.PreCommit = string(data)
	}

	// Check for pre-push hook
	prePushPath := filepath.Join(huskyDir, "pre-push")
	if data, err := os.ReadFile(prePushPath); err == nil {
		config.PrePush = string(data)
	}

	// Check for commit-msg hook
	commitMsgPath := filepath.Join(huskyDir, "commit-msg")
	if data, err := os.ReadFile(commitMsgPath); err == nil {
		config.CommitMsg = string(data)
	}

	return config, nil
}

// GetPreCommitHook returns the content of the pre-commit hook
func (h *Hooks) GetPreCommitHook() string {
	return `#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

# Run composer lint script
if [ -f "composer.json" ]; then
    composer run-script lint || exit 1
fi

# If stax is available, could also run stax lint
# if command -v stax &> /dev/null; then
#     stax lint:staged || exit 1
# fi
`
}

// TestPreCommitHook tests the pre-commit hook execution
func (h *Hooks) TestPreCommitHook() error {
	// Try running composer lint
	composer := NewComposer(h.projectPath)
	scripts, err := composer.ListScripts()
	if err != nil {
		return fmt.Errorf("failed to list composer scripts: %w", err)
	}

	if _, exists := scripts["lint"]; !exists {
		return fmt.Errorf("lint script not found in composer.json")
	}

	// Try executing the lint script
	if err := composer.Lint(); err != nil {
		return fmt.Errorf("lint script failed: %w", err)
	}

	return nil
}

// BypassHook returns instructions for bypassing hooks
func (h *Hooks) BypassHook() string {
	return `To bypass git hooks (use with caution):

  git commit --no-verify

Or set environment variable:

  HUSKY=0 git commit

Note: Bypassing hooks means code quality checks will be skipped.
Only use when necessary and ensure code is checked manually.`
}

// GetHookStatus returns the status of git hooks
func (h *Hooks) GetHookStatus() map[string]bool {
	status := map[string]bool{
		"pre-commit":         false,
		"pre-push":           false,
		"commit-msg":         false,
		"prepare-commit-msg": false,
	}

	huskyDir := filepath.Join(h.projectPath, ".husky")

	for hook := range status {
		hookPath := filepath.Join(huskyDir, hook)
		if _, err := os.Stat(hookPath); err == nil {
			status[hook] = true
		}
	}

	return status
}

// InstallHook installs a specific git hook
func (h *Hooks) InstallHook(hookName string, content string) error {
	huskyDir := filepath.Join(h.projectPath, ".husky")

	// Create .husky directory if it doesn't exist
	if err := os.MkdirAll(huskyDir, 0755); err != nil {
		return fmt.Errorf("failed to create .husky directory: %w", err)
	}

	hookPath := filepath.Join(huskyDir, hookName)
	return os.WriteFile(hookPath, []byte(content), 0755)
}

// RemoveHook removes a specific git hook
func (h *Hooks) RemoveHook(hookName string) error {
	huskyDir := filepath.Join(h.projectPath, ".husky")
	hookPath := filepath.Join(huskyDir, hookName)

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return fmt.Errorf("hook %s does not exist", hookName)
	}

	return os.Remove(hookPath)
}

// GetPrePushHook returns a pre-push hook template
func (h *Hooks) GetPrePushHook() string {
	return `#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

# Run tests before pushing
if [ -f "composer.json" ]; then
    composer run-script test || exit 1
fi
`
}

// GetCommitMsgHook returns a commit-msg hook template
func (h *Hooks) GetCommitMsgHook() string {
	return `#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

# Validate commit message format
# Example: require conventional commits format
# commit_msg=$(cat "$1")
# if ! echo "$commit_msg" | grep -qE "^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?: .+"; then
#     echo "Error: Commit message must follow conventional commits format"
#     echo "Example: feat(build): add new build command"
#     exit 1
# fi
`
}

// EnableHook enables a git hook
func (h *Hooks) EnableHook(hookName string) error {
	var content string

	switch hookName {
	case "pre-commit":
		content = h.GetPreCommitHook()
	case "pre-push":
		content = h.GetPrePushHook()
	case "commit-msg":
		content = h.GetCommitMsgHook()
	default:
		return fmt.Errorf("unknown hook: %s", hookName)
	}

	return h.InstallHook(hookName, content)
}

// DisableHook disables a git hook by removing it
func (h *Hooks) DisableHook(hookName string) error {
	return h.RemoveHook(hookName)
}

// ListHooks lists all installed hooks
func (h *Hooks) ListHooks() ([]string, error) {
	huskyDir := filepath.Join(h.projectPath, ".husky")

	if _, err := os.Stat(huskyDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(huskyDir)
	if err != nil {
		return nil, err
	}

	hooks := []string{}
	for _, entry := range entries {
		// Skip directories and special files
		if entry.IsDir() || strings.HasPrefix(entry.Name(), "_") || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		hooks = append(hooks, entry.Name())
	}

	return hooks, nil
}

// CheckHookExecutable verifies that a hook is executable
func (h *Hooks) CheckHookExecutable(hookName string) error {
	huskyDir := filepath.Join(h.projectPath, ".husky")
	hookPath := filepath.Join(huskyDir, hookName)

	info, err := os.Stat(hookPath)
	if err != nil {
		return fmt.Errorf("hook not found: %w", err)
	}

	// Check if executable bit is set
	mode := info.Mode()
	if mode&0111 == 0 {
		return fmt.Errorf("hook is not executable (run: chmod +x %s)", hookPath)
	}

	return nil
}

// MakeHookExecutable makes a hook executable
func (h *Hooks) MakeHookExecutable(hookName string) error {
	huskyDir := filepath.Join(h.projectPath, ".husky")
	hookPath := filepath.Join(huskyDir, hookName)

	return os.Chmod(hookPath, 0755)
}

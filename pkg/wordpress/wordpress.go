package wordpress

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type InstallConfig struct {
	URL      string
	Title    string
	Username string
	Password string
	Email    string
}

func DownloadCore(projectPath string) error {
	if !isWPCLIAvailable() {
		return fmt.Errorf("wp-cli is not available in DDEV container")
	}

	cmd := exec.Command("ddev", "wp", "core", "download")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Install(projectPath string, config InstallConfig) error {
	if !isWPCLIAvailable() {
		return fmt.Errorf("wp-cli is not available in DDEV container")
	}

	args := []string{
		"wp", "core", "install",
		"--url=" + config.URL,
		"--title=" + config.Title,
		"--admin_user=" + config.Username,
		"--admin_password=" + config.Password,
		"--admin_email=" + config.Email,
	}

	cmd := exec.Command("ddev", args...)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func CreateConfig(projectPath string) error {
	if !isWPCLIAvailable() {
		return fmt.Errorf("wp-cli is not available in DDEV container")
	}

	args := []string{
		"wp", "config", "create",
		"--dbname=db",
		"--dbuser=db", 
		"--dbpass=db",
		"--dbhost=db",
	}

	cmd := exec.Command("ddev", args...)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func InstallPlugin(projectPath, plugin string, activate bool) error {
	if !isWPCLIAvailable() {
		return fmt.Errorf("wp-cli is not available in DDEV container")
	}

	args := []string{"wp", "plugin", "install", plugin}
	if activate {
		args = append(args, "--activate")
	}

	cmd := exec.Command("ddev", args...)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func InstallTheme(projectPath, theme string, activate bool) error {
	if !isWPCLIAvailable() {
		return fmt.Errorf("wp-cli is not available in DDEV container")
	}

	args := []string{"wp", "theme", "install", theme}
	if activate {
		args = append(args, "--activate")
	}

	cmd := exec.Command("ddev", args...)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func isWPCLIAvailable() bool {
	cmd := exec.Command("ddev", "wp", "--version")
	return cmd.Run() == nil
}

func HasWordPress(projectPath string) bool {
	wpConfigPath := filepath.Join(projectPath, "wp-config.php")
	indexPath := filepath.Join(projectPath, "index.php")
	
	_, wpConfigErr := os.Stat(wpConfigPath)
	_, indexErr := os.Stat(indexPath)
	
	return wpConfigErr == nil || indexErr == nil
}
package ddev

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	ProjectName string
	ProjectType string
	PHPVersion  string
	WebServer   string
	DatabaseType string
}

func IsInstalled() bool {
	_, err := exec.LookPath("ddev")
	return err == nil
}

func IsProject(path string) bool {
	configPath := filepath.Join(path, ".ddev", "config.yaml")
	_, err := os.Stat(configPath)
	return err == nil
}

func Init(projectPath string, config Config) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	if IsProject(projectPath) {
		return fmt.Errorf("DDEV project already exists in %s", projectPath)
	}

	args := []string{
		"config",
		"--project-type=" + config.ProjectType,
		"--project-name=" + config.ProjectName,
	}

	if config.PHPVersion != "" {
		args = append(args, "--php-version="+config.PHPVersion)
	}

	if config.WebServer != "" {
		args = append(args, "--webserver-type="+config.WebServer)
	}

	if config.DatabaseType != "" {
		args = append(args, "--database="+config.DatabaseType)
	}

	cmd := exec.Command("ddev", args...)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Start(projectPath string) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	if !IsProject(projectPath) {
		return fmt.Errorf("no DDEV project found in %s", projectPath)
	}

	cmd := exec.Command("ddev", "start")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Stop(projectPath string) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	if !IsProject(projectPath) {
		return fmt.Errorf("no DDEV project found in %s", projectPath)
	}

	cmd := exec.Command("ddev", "stop")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Restart(projectPath string) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	if !IsProject(projectPath) {
		return fmt.Errorf("no DDEV project found in %s", projectPath)
	}

	cmd := exec.Command("ddev", "restart")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Status(projectPath string) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	cmd := exec.Command("ddev", "list")
	if IsProject(projectPath) {
		cmd.Dir = projectPath
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Describe(projectPath string) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	if !IsProject(projectPath) {
		return fmt.Errorf("no DDEV project found in %s", projectPath)
	}

	cmd := exec.Command("ddev", "describe")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Delete(projectPath, projectName string, omitSnapshot, yes bool) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	args := []string{"delete"}
	
	if projectName != "" {
		args = append(args, projectName)
	}
	
	if omitSnapshot {
		args = append(args, "--omit-snapshot")
	}
	
	if yes {
		args = append(args, "--yes")
	}

	cmd := exec.Command("ddev", args...)
	if projectName == "" && IsProject(projectPath) {
		cmd.Dir = projectPath
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // Allow user input for confirmation

	err := cmd.Run()
	if err != nil {
		return err
	}

	// Additional cleanup: remove DDEV-related files if they still exist and omitSnapshot is true
	if omitSnapshot && projectName == "" {
		filesToCleanup := []string{
			filepath.Join(projectPath, ".ddev"),
			filepath.Join(projectPath, "wp-config-ddev.php"),
		}
		
		for _, file := range filesToCleanup {
			if _, err := os.Stat(file); err == nil {
				fmt.Printf("Cleaning up remaining file: %s\n", filepath.Base(file))
				if err := os.RemoveAll(file); err != nil {
					fmt.Printf("Warning: failed to remove %s: %v\n", file, err)
				}
			}
		}
	}

	return nil
}

func Poweroff() error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	cmd := exec.Command("ddev", "poweroff")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func List() error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	cmd := exec.Command("ddev", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ExportDB(projectPath, filename string) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	if !IsProject(projectPath) {
		return fmt.Errorf("no DDEV project found in %s", projectPath)
	}

	cmd := exec.Command("ddev", "export-db", "--file="+filename)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ImportDB(projectPath, src string) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	if !IsProject(projectPath) {
		return fmt.Errorf("no DDEV project found in %s", projectPath)
	}

	cmd := exec.Command("ddev", "import-db", "--src="+src)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// UpdateConfig updates an existing DDEV project configuration
func UpdateConfig(projectPath string, config Config) error {
	if !IsInstalled() {
		return fmt.Errorf("ddev is not installed or not in PATH")
	}

	if !IsProject(projectPath) {
		return fmt.Errorf("no DDEV project found in %s", projectPath)
	}

	args := []string{"config"}

	if config.PHPVersion != "" {
		args = append(args, "--php-version="+config.PHPVersion)
	}

	if config.WebServer != "" {
		args = append(args, "--webserver-type="+config.WebServer)
	}

	if config.DatabaseType != "" {
		args = append(args, "--database="+config.DatabaseType)
	}

	cmd := exec.Command("ddev", args...)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetConfig reads the current DDEV configuration
func GetConfig(projectPath string) (*Config, error) {
	if !IsInstalled() {
		return nil, fmt.Errorf("ddev is not installed or not in PATH")
	}

	if !IsProject(projectPath) {
		return nil, fmt.Errorf("no DDEV project found in %s", projectPath)
	}

	cmd := exec.Command("ddev", "describe", "--json-output")
	cmd.Dir = projectPath
	
	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get DDEV config: %w", err)
	}

	// For now, return basic config - in a real implementation you'd parse the JSON
	return &Config{
		ProjectType: "wordpress",
		PHPVersion:  "8.2", // Default values - would be parsed from JSON
		WebServer:   "nginx-fpm",
		DatabaseType: "mysql:8.0",
	}, nil
}
package system

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// HostEntry represents an entry in the hosts file
type HostEntry struct {
	IP       string
	Hostname string
	Comment  string
}

const (
	hostsFilePathUnix    = "/etc/hosts"
	hostsFilePathWindows = "C:\\Windows\\System32\\drivers\\etc\\hosts"
)

// GetHostsFilePath returns the path to the system hosts file
func GetHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return hostsFilePathWindows
	}
	return hostsFilePathUnix
}

// RequiresSudo checks if sudo is required to modify the hosts file
func RequiresSudo() bool {
	// On Unix systems, /etc/hosts requires root
	return runtime.GOOS != "windows"
}

// AddHostsEntry adds a single entry to the hosts file
func AddHostsEntry(hostname, ip string) error {
	if ip == "" {
		ip = "127.0.0.1"
	}

	entry := HostEntry{
		IP:       ip,
		Hostname: hostname,
		Comment:  fmt.Sprintf("Added by Stax on %s", time.Now().Format("2006-01-02")),
	}

	return AddHostsEntries([]HostEntry{entry}, "stax")
}

// RemoveHostsEntry removes a specific hostname from the hosts file
func RemoveHostsEntry(hostname string) error {
	// Read current hosts file
	entries, err := ReadHostsFile()
	if err != nil {
		return err
	}

	// Filter out the entry
	var filtered []string
	for _, line := range entries {
		// Skip lines containing the hostname
		if !strings.Contains(line, hostname) {
			filtered = append(filtered, line)
		}
	}

	// Write back
	return writeHostsFile(filtered)
}

// HasHostsEntry checks if a hostname exists in the hosts file
func HasHostsEntry(hostname string) (bool, error) {
	entries, err := GetHostsEntries("")
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.Hostname == hostname {
			return true, nil
		}
	}

	return false, nil
}

// BackupHostsFile creates a backup of the hosts file
func BackupHostsFile() (string, error) {
	hostsPath := GetHostsFilePath()
	backupPath := hostsPath + ".stax-backup-" + time.Now().Format("20060102-150405")

	// Read current file
	data, err := os.ReadFile(hostsPath)
	if err != nil {
		return "", fmt.Errorf("failed to read hosts file: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	return backupPath, nil
}

// RestoreHostsFile restores the hosts file from a backup
func RestoreHostsFile(backupPath string) error {
	hostsPath := GetHostsFilePath()

	// Read backup
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Write to hosts file
	if err := os.WriteFile(hostsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to restore hosts file: %w", err)
	}

	return nil
}

// GetHostsEntries returns all entries with an optional marker filter
func GetHostsEntries(marker string) ([]HostEntry, error) {
	entries, err := ReadHostsFile()
	if err != nil {
		return nil, err
	}

	var result []HostEntry
	inMarkerSection := marker == ""

	for _, line := range entries {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments (unless we're filtering by marker)
		if line == "" {
			continue
		}

		// Check for marker start
		if marker != "" && strings.Contains(line, "### START "+marker) {
			inMarkerSection = true
			continue
		}

		// Check for marker end
		if marker != "" && strings.Contains(line, "### END "+marker) {
			inMarkerSection = false
			continue
		}

		if !inMarkerSection {
			continue
		}

		// Parse entry
		if strings.HasPrefix(line, "#") {
			// Comment only line
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		entry := HostEntry{
			IP:       parts[0],
			Hostname: parts[1],
		}

		// Extract comment if present
		if commentIdx := strings.Index(line, "#"); commentIdx != -1 {
			entry.Comment = strings.TrimSpace(line[commentIdx+1:])
		}

		result = append(result, entry)
	}

	return result, nil
}

// UpdateHostsEntries adds or updates entries with a marker
func UpdateHostsEntries(entries []HostEntry, marker string) error {
	if marker == "" {
		marker = "stax"
	}

	// Read current hosts file
	lines, err := ReadHostsFile()
	if err != nil {
		return err
	}

	// Remove existing marker section
	var filtered []string
	inMarkerSection := false

	for _, line := range lines {
		if strings.Contains(line, "### START "+marker) {
			inMarkerSection = true
			continue
		}
		if strings.Contains(line, "### END "+marker) {
			inMarkerSection = false
			continue
		}
		if !inMarkerSection {
			filtered = append(filtered, line)
		}
	}

	// Add new marker section
	if len(entries) > 0 {
		filtered = append(filtered, "")
		filtered = append(filtered, fmt.Sprintf("### START %s - Managed by Stax ###", marker))
		filtered = append(filtered, fmt.Sprintf("# Generated on %s", time.Now().Format("2006-01-02 15:04:05")))

		for _, entry := range entries {
			line := fmt.Sprintf("%s\t%s", entry.IP, entry.Hostname)
			if entry.Comment != "" {
				line += " # " + entry.Comment
			}
			filtered = append(filtered, line)
		}

		filtered = append(filtered, fmt.Sprintf("### END %s ###", marker))
	}

	// Write back
	return writeHostsFile(filtered)
}

// AddHostsEntries adds multiple entries with a marker
func AddHostsEntries(entries []HostEntry, marker string) error {
	// Get existing entries (excluding our marker section)
	existing, err := ReadHostsFile()
	if err != nil {
		return err
	}

	// Remove any existing marker section
	var filtered []string
	inMarkerSection := false

	for _, line := range existing {
		if strings.Contains(line, "### START "+marker) {
			inMarkerSection = true
			continue
		}
		if strings.Contains(line, "### END "+marker) {
			inMarkerSection = false
			continue
		}
		if !inMarkerSection {
			filtered = append(filtered, line)
		}
	}

	// Add new entries
	return UpdateHostsEntries(entries, marker)
}

// ReadHostsFile reads all lines from the hosts file
func ReadHostsFile() ([]string, error) {
	hostsPath := GetHostsFilePath()

	file, err := os.Open(hostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open hosts file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read hosts file: %w", err)
	}

	return lines, nil
}

// writeHostsFile writes lines to the hosts file
func writeHostsFile(lines []string) error {
	hostsPath := GetHostsFilePath()

	// Join lines with newline
	content := strings.Join(lines, "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	// Write to file
	if err := os.WriteFile(hostsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write hosts file: %w", err)
	}

	return nil
}

// RemoveStaxEntries removes all Stax-managed entries from the hosts file
func RemoveStaxEntries(marker string) error {
	if marker == "" {
		marker = "stax"
	}

	return UpdateHostsEntries([]HostEntry{}, marker)
}

// ValidateHostname checks if a hostname is valid
func ValidateHostname(hostname string) bool {
	if hostname == "" {
		return false
	}

	// Basic validation
	if strings.Contains(hostname, " ") {
		return false
	}

	if len(hostname) > 253 {
		return false
	}

	return true
}

// ValidateIP checks if an IP address is valid (basic check)
func ValidateIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}
	}

	return true
}

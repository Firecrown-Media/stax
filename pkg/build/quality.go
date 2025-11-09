package build

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Quality handles code quality operations (PHPCS, PHPCBF)
type Quality struct {
	projectPath string
}

// NewQuality creates a new Quality instance
func NewQuality(projectPath string) *Quality {
	return &Quality{
		projectPath: projectPath,
	}
}

// RunPHPCS runs PHP CodeSniffer
func (q *Quality) RunPHPCS(options PHPCSOptions) (*PHPCSResult, error) {
	// Find phpcs executable
	phpcsPath, err := q.findPHPCS()
	if err != nil {
		return nil, err
	}

	// Build command arguments
	args := q.buildPHPCSArgs(options)

	// Execute PHPCS
	cmd := exec.Command(phpcsPath, args...)
	cmd.Dir = q.projectPath
	output, err := cmd.CombinedOutput()

	// Parse result
	result := q.parsePHPCSOutput(string(output), options.Report)

	// PHPCS returns exit code 1 if errors found, 2 if warnings found
	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
		// Don't treat this as an error - it's expected behavior
		err = nil
	}

	result.Success = result.ExitCode == 0

	return result, err
}

// RunPHPCBF runs PHP Code Beautifier and Fixer
func (q *Quality) RunPHPCBF(options PHPCSOptions) error {
	// Find phpcbf executable
	phpcbfPath, err := q.findPHPCBF()
	if err != nil {
		return err
	}

	// Build command arguments (same as PHPCS)
	args := q.buildPHPCSArgs(options)

	// Execute PHPCBF
	cmd := exec.Command(phpcbfPath, args...)
	cmd.Dir = q.projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// PHPCBF returns exit code 1 if it fixed files
	err = cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 1 {
			// This means files were fixed - not an error
			return nil
		}
	}

	return err
}

// RunComposerLint runs composer lint script
func (q *Quality) RunComposerLint() error {
	composer := NewComposer(q.projectPath)
	return composer.Lint()
}

// RunComposerFix runs composer fix script
func (q *Quality) RunComposerFix() error {
	composer := NewComposer(q.projectPath)
	return composer.Fix()
}

// GetPHPCSConfig finds the PHPCS configuration file
func (q *Quality) GetPHPCSConfig() (string, error) {
	// Check for phpcs.xml.dist
	distConfig := filepath.Join(q.projectPath, ".phpcs.xml.dist")
	if _, err := os.Stat(distConfig); err == nil {
		return distConfig, nil
	}

	// Check for phpcs.xml
	xmlConfig := filepath.Join(q.projectPath, "phpcs.xml")
	if _, err := os.Stat(xmlConfig); err == nil {
		return xmlConfig, nil
	}

	// Check for .phpcs.xml
	dotConfig := filepath.Join(q.projectPath, ".phpcs.xml")
	if _, err := os.Stat(dotConfig); err == nil {
		return dotConfig, nil
	}

	return "", fmt.Errorf("PHPCS configuration file not found")
}

// ValidatePHPCSConfig validates the PHPCS configuration
func (q *Quality) ValidatePHPCSConfig() error {
	configPath, err := q.GetPHPCSConfig()
	if err != nil {
		return err
	}

	// Basic validation - just check if file is readable XML
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read PHPCS config: %w", err)
	}

	// Check for basic XML structure
	if !strings.Contains(string(data), "<ruleset") {
		return fmt.Errorf("invalid PHPCS config: missing <ruleset> tag")
	}

	return nil
}

// GetCodingStandards returns available coding standards
func (q *Quality) GetCodingStandards() ([]string, error) {
	phpcsPath, err := q.findPHPCS()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(phpcsPath, "-i")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse output like: "The installed coding standards are PSR1, PSR2, PSR12, etc."
	standardsStr := strings.TrimSpace(string(output))
	standardsStr = strings.TrimPrefix(standardsStr, "The installed coding standards are ")
	standardsStr = strings.TrimSuffix(standardsStr, ".")

	standards := strings.Split(standardsStr, ", ")

	// Trim whitespace from each standard
	for i, std := range standards {
		standards[i] = strings.TrimSpace(std)
	}

	return standards, nil
}

// FormatPHPCSResults formats PHPCS results for display
func (q *Quality) FormatPHPCSResults(result *PHPCSResult) string {
	if result.Success {
		return "No errors or warnings found"
	}

	var output strings.Builder

	output.WriteString(fmt.Sprintf("\nFound %d error(s) and %d warning(s)", result.Errors, result.Warnings))
	if result.Fixable > 0 {
		output.WriteString(fmt.Sprintf(" (%d fixable)", result.Fixable))
	}
	output.WriteString("\n\n")

	for _, file := range result.Files {
		if len(file.Messages) == 0 {
			continue
		}

		output.WriteString(fmt.Sprintf("FILE: %s\n", file.File))
		output.WriteString(fmt.Sprintf("FOUND %d ERROR(S) AND %d WARNING(S)\n", file.Errors, file.Warnings))
		output.WriteString(strings.Repeat("-", 80) + "\n")

		for _, msg := range file.Messages {
			fixable := ""
			if msg.Fixable {
				fixable = " [x]"
			}

			output.WriteString(fmt.Sprintf("%d:%d | %s | %s%s\n",
				msg.Line, msg.Column, msg.Type, msg.Message, fixable))

			if msg.Source != "" {
				output.WriteString(fmt.Sprintf("     | (%s)\n", msg.Source))
			}
		}

		output.WriteString("\n")
	}

	return output.String()
}

// findPHPCS locates the phpcs executable
func (q *Quality) findPHPCS() (string, error) {
	// Check vendor/bin/phpcs
	vendorPhpcs := filepath.Join(q.projectPath, "vendor", "bin", "phpcs")
	if _, err := os.Stat(vendorPhpcs); err == nil {
		return vendorPhpcs, nil
	}

	// Check global phpcs
	if path, err := exec.LookPath("phpcs"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("phpcs not found (install via composer or globally)")
}

// findPHPCBF locates the phpcbf executable
func (q *Quality) findPHPCBF() (string, error) {
	// Check vendor/bin/phpcbf
	vendorPhpcbf := filepath.Join(q.projectPath, "vendor", "bin", "phpcbf")
	if _, err := os.Stat(vendorPhpcbf); err == nil {
		return vendorPhpcbf, nil
	}

	// Check global phpcbf
	if path, err := exec.LookPath("phpcbf"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("phpcbf not found (install via composer or globally)")
}

// buildPHPCSArgs builds command-line arguments for PHPCS/PHPCBF
func (q *Quality) buildPHPCSArgs(options PHPCSOptions) []string {
	args := []string{}

	// Config file
	if options.ConfigFile != "" {
		args = append(args, fmt.Sprintf("--standard=%s", options.ConfigFile))
	} else if configPath, err := q.GetPHPCSConfig(); err == nil {
		args = append(args, fmt.Sprintf("--standard=%s", configPath))
	}

	// Standard (overrides config)
	if options.Standard != "" {
		args = append(args, fmt.Sprintf("--standard=%s", options.Standard))
	}

	// Extensions
	if len(options.Extensions) > 0 {
		args = append(args, fmt.Sprintf("--extensions=%s", strings.Join(options.Extensions, ",")))
	}

	// Ignore pattern
	if options.Ignore != "" {
		args = append(args, fmt.Sprintf("--ignore=%s", options.Ignore))
	}

	// Report format
	if options.Report != "" {
		args = append(args, fmt.Sprintf("--report=%s", options.Report))
	} else {
		args = append(args, "--report=json")
	}

	// Show sniff codes
	if options.ShowSniffs {
		args = append(args, "-s")
	}

	// Severity levels
	if options.Severity > 0 {
		args = append(args, fmt.Sprintf("--severity=%d", options.Severity))
	}

	if options.ErrorSeverity > 0 {
		args = append(args, fmt.Sprintf("--error-severity=%d", options.ErrorSeverity))
	}

	if options.WarningSeverity > 0 {
		args = append(args, fmt.Sprintf("--warning-severity=%d", options.WarningSeverity))
	}

	// Files/directories to check
	if len(options.Files) > 0 {
		args = append(args, options.Files...)
	} else {
		args = append(args, ".")
	}

	return args
}

// parsePHPCSOutput parses PHPCS output into a structured result
func (q *Quality) parsePHPCSOutput(output string, reportFormat string) *PHPCSResult {
	result := &PHPCSResult{
		Output: output,
		Files:  []PHPCSFileResult{},
	}

	// If JSON format, parse as JSON
	if reportFormat == "json" || strings.Contains(output, `"totals"`) {
		q.parseJSONOutput(output, result)
	} else {
		q.parseTextOutput(output, result)
	}

	return result
}

// parseJSONOutput parses JSON PHPCS output
func (q *Quality) parseJSONOutput(output string, result *PHPCSResult) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
		return
	}

	// Parse totals
	if totals, ok := jsonData["totals"].(map[string]interface{}); ok {
		if errors, ok := totals["errors"].(float64); ok {
			result.Errors = int(errors)
		}
		if warnings, ok := totals["warnings"].(float64); ok {
			result.Warnings = int(warnings)
		}
		if fixable, ok := totals["fixable"].(float64); ok {
			result.Fixable = int(fixable)
		}
	}

	// Parse files
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		for filePath, fileData := range files {
			fileDataMap := fileData.(map[string]interface{})
			fileResult := PHPCSFileResult{
				File:     filePath,
				Messages: []PHPCSMessage{},
			}

			if errors, ok := fileDataMap["errors"].(float64); ok {
				fileResult.Errors = int(errors)
			}
			if warnings, ok := fileDataMap["warnings"].(float64); ok {
				fileResult.Warnings = int(warnings)
			}

			// Parse messages
			if messages, ok := fileDataMap["messages"].([]interface{}); ok {
				for _, msg := range messages {
					msgMap := msg.(map[string]interface{})
					message := PHPCSMessage{}

					if line, ok := msgMap["line"].(float64); ok {
						message.Line = int(line)
					}
					if column, ok := msgMap["column"].(float64); ok {
						message.Column = int(column)
					}
					if msgType, ok := msgMap["type"].(string); ok {
						message.Type = msgType
					}
					if msgText, ok := msgMap["message"].(string); ok {
						message.Message = msgText
					}
					if source, ok := msgMap["source"].(string); ok {
						message.Source = source
					}
					if severity, ok := msgMap["severity"].(float64); ok {
						message.Severity = int(severity)
					}
					if fixable, ok := msgMap["fixable"].(bool); ok {
						message.Fixable = fixable
					}

					fileResult.Messages = append(fileResult.Messages, message)
				}
			}

			result.Files = append(result.Files, fileResult)
		}
	}
}

// parseTextOutput parses text PHPCS output
func (q *Quality) parseTextOutput(output string, result *PHPCSResult) {
	lines := strings.Split(output, "\n")

	// Look for summary line like "FOUND 5 ERRORS AND 3 WARNINGS AFFECTING 2 FILES"
	summaryRegex := regexp.MustCompile(`FOUND (\d+) ERROR(?:S)?(?: AND (\d+) WARNING(?:S)?)?`)
	fixableRegex := regexp.MustCompile(`(\d+) FIXABLE`)

	for _, line := range lines {
		if matches := summaryRegex.FindStringSubmatch(line); matches != nil {
			if errors, err := strconv.Atoi(matches[1]); err == nil {
				result.Errors = errors
			}
			if len(matches) > 2 && matches[2] != "" {
				if warnings, err := strconv.Atoi(matches[2]); err == nil {
					result.Warnings = warnings
				}
			}
		}

		if matches := fixableRegex.FindStringSubmatch(line); matches != nil {
			if fixable, err := strconv.Atoi(matches[1]); err == nil {
				result.Fixable = fixable
			}
		}
	}
}

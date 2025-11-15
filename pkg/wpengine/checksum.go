package wpengine

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/firecrown-media/stax/pkg/security"
)

// ChecksumResult represents the result of checksum verification
type ChecksumResult struct {
	TotalFiles      int
	MatchedFiles    int
	MismatchedFiles int
	MissingLocal    int
	MissingRemote   int
	Mismatches      []FileMismatch
	MissingLocally  []string
	MissingRemotely []string
}

// FileMismatch represents a file with different checksums
type FileMismatch struct {
	RelativePath   string
	RemoteChecksum string
	LocalChecksum  string
}

// GenerateRemoteChecksums generates MD5 checksums for files in a remote directory
// Returns a map of relative_path -> checksum
func (c *SSHClient) GenerateRemoteChecksums(remotePath string) (map[string]string, error) {
	// Sanitize path to prevent command injection
	safePath, err := security.SanitizeForShell(remotePath)
	if err != nil {
		return nil, fmt.Errorf("invalid remote path: %w", err)
	}

	// Use find + md5sum to generate checksums
	// The command finds all files and generates MD5 checksums
	// Format: checksum  filepath
	cmd := fmt.Sprintf("cd %s && find . -type f -exec md5sum {} \\;", safePath)

	output, err := c.ExecuteCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to generate remote checksums: %w", err)
	}

	checksums := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse md5sum output: "checksum  filepath"
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		checksum := parts[0]
		// Join remaining parts as the path (in case path has spaces)
		relativePath := strings.Join(parts[1:], " ")

		// Remove leading "./" from path
		relativePath = strings.TrimPrefix(relativePath, "./")

		checksums[relativePath] = checksum
	}

	return checksums, nil
}

// GenerateLocalChecksums generates MD5 checksums for files in a local directory
// Returns a map of relative_path -> checksum
func GenerateLocalChecksums(localPath string) (map[string]string, error) {
	checksums := make(map[string]string)

	// Validate local path exists
	if _, err := os.Stat(localPath); err != nil {
		return nil, fmt.Errorf("local path does not exist: %w", err)
	}

	// Walk the directory tree
	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate relative path
		relativePath, err := filepath.Rel(localPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Generate MD5 checksum
		checksum, err := calculateMD5(path)
		if err != nil {
			return fmt.Errorf("failed to calculate checksum for %s: %w", relativePath, err)
		}

		checksums[relativePath] = checksum
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk local directory: %w", err)
	}

	return checksums, nil
}

// calculateMD5 calculates the MD5 checksum of a file
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

// VerifyChecksums compares remote and local checksums
// Returns detailed results about matches, mismatches, and missing files
func VerifyChecksums(remoteChecksums, localChecksums map[string]string) *ChecksumResult {
	result := &ChecksumResult{
		TotalFiles:      len(remoteChecksums),
		Mismatches:      make([]FileMismatch, 0),
		MissingLocally:  make([]string, 0),
		MissingRemotely: make([]string, 0),
	}

	// Check all remote files
	for relativePath, remoteChecksum := range remoteChecksums {
		localChecksum, exists := localChecksums[relativePath]

		if !exists {
			// File exists remotely but not locally
			result.MissingLocal++
			result.MissingLocally = append(result.MissingLocally, relativePath)
		} else if localChecksum != remoteChecksum {
			// File exists but checksums don't match
			result.MismatchedFiles++
			result.Mismatches = append(result.Mismatches, FileMismatch{
				RelativePath:   relativePath,
				RemoteChecksum: remoteChecksum,
				LocalChecksum:  localChecksum,
			})
		} else {
			// File matches
			result.MatchedFiles++
		}
	}

	// Check for files that exist locally but not remotely
	for relativePath := range localChecksums {
		if _, exists := remoteChecksums[relativePath]; !exists {
			result.MissingRemote++
			result.MissingRemotely = append(result.MissingRemotely, relativePath)
		}
	}

	return result
}

// VerifyFileChecksums performs complete checksum verification
// This is a high-level function that combines remote and local checksum generation
// and comparison, returning detailed results
func (c *SSHClient) VerifyFileChecksums(remotePath, localPath string) (*ChecksumResult, error) {
	// Generate remote checksums
	remoteChecksums, err := c.GenerateRemoteChecksums(remotePath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate remote checksums: %w", err)
	}

	// Generate local checksums
	localChecksums, err := GenerateLocalChecksums(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate local checksums: %w", err)
	}

	// Compare and return results
	result := VerifyChecksums(remoteChecksums, localChecksums)
	return result, nil
}

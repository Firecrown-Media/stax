package snapshot

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/firecrown-media/stax/pkg/config"
	"github.com/firecrown-media/stax/pkg/ddev"
)

// Manager handles snapshot operations
type Manager struct {
	Config      *config.Config
	ProjectDir  string
	DDEVManager *ddev.Manager
}

// NewManager creates a new snapshot manager
func NewManager(cfg *config.Config, projectDir string) *Manager {
	return &Manager{
		Config:      cfg,
		ProjectDir:  projectDir,
		DDEVManager: ddev.NewManager(projectDir),
	}
}

// CreateSnapshot creates a database snapshot
// Returns the snapshot filename and any error
func (m *Manager) CreateSnapshot(projectName, snapshotType string) (string, error) {
	// Validate snapshot type
	var snapType SnapshotType
	switch snapshotType {
	case "auto":
		snapType = Auto
	case "manual":
		snapType = Manual
	default:
		return "", fmt.Errorf("invalid snapshot type: %s (must be 'auto' or 'manual')", snapshotType)
	}

	// Get snapshot directory and expand ~
	snapshotDir := expandPath(m.Config.Snapshots.Directory)
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// Generate snapshot filename: {project}-{timestamp}-{type}.sql.gz
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s-%s.sql.gz", projectName, timestamp, snapType)
	snapshotPath := filepath.Join(snapshotDir, filename)

	// Create temporary uncompressed file
	tmpFile, err := os.CreateTemp("", "stax-snapshot-*.sql")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Export database using DDEV
	if err := m.DDEVManager.ExportDB(tmpPath); err != nil {
		return "", fmt.Errorf("failed to export database: %w", err)
	}

	// Compress with gzip
	if err := compressFile(tmpPath, snapshotPath); err != nil {
		return "", fmt.Errorf("failed to compress snapshot: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(snapshotPath)
	if err != nil {
		return "", fmt.Errorf("failed to get snapshot file info: %w", err)
	}

	// Record metadata
	metadataPath := filepath.Join(snapshotDir, "metadata.json")
	store, err := LoadMetadata(metadataPath)
	if err != nil {
		return "", fmt.Errorf("failed to load metadata: %w", err)
	}

	metadata := SnapshotMetadata{
		File:      filename,
		Project:   projectName,
		Timestamp: time.Now(),
		Type:      snapType,
		Size:      fileInfo.Size(),
		CreatedBy: getCreatedByContext(snapshotType),
	}

	store.AddSnapshot(metadata)

	if err := SaveMetadata(metadataPath, store); err != nil {
		return "", fmt.Errorf("failed to save metadata: %w", err)
	}

	return filename, nil
}

// RestoreSnapshot restores a database from a snapshot
func (m *Manager) RestoreSnapshot(snapshotPath string) error {
	// Expand ~ in path
	snapshotPath = expandPath(snapshotPath)

	// Check if snapshot exists
	if _, err := os.Stat(snapshotPath); err != nil {
		return fmt.Errorf("snapshot not found: %w", err)
	}

	// Create temporary decompressed file
	tmpFile, err := os.CreateTemp("", "stax-restore-*.sql")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Decompress snapshot
	if err := decompressFile(snapshotPath, tmpPath); err != nil {
		return fmt.Errorf("failed to decompress snapshot: %w", err)
	}

	// Import to database using DDEV
	if err := m.DDEVManager.ImportDB(tmpPath); err != nil {
		return fmt.Errorf("failed to import database: %w", err)
	}

	return nil
}

// ListSnapshots lists all snapshots for a project
func (m *Manager) ListSnapshots(projectName string) ([]SnapshotMetadata, error) {
	// Get snapshot directory and expand ~
	snapshotDir := expandPath(m.Config.Snapshots.Directory)

	// Load metadata
	metadataPath := filepath.Join(snapshotDir, "metadata.json")
	store, err := LoadMetadata(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	// Filter by project if specified
	if projectName != "" {
		return store.GetSnapshotsByProject(projectName), nil
	}

	return store.Snapshots, nil
}

// CleanSnapshots deletes snapshots based on retention policy
func (m *Manager) CleanSnapshots(cfg *config.Config) error {
	// Get snapshot directory and expand ~
	snapshotDir := expandPath(cfg.Snapshots.Directory)

	// Load metadata
	metadataPath := filepath.Join(snapshotDir, "metadata.json")
	store, err := LoadMetadata(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	// Find old snapshots
	autoSnapshots := store.GetOldSnapshots(Auto, cfg.Snapshots.Retention.Auto)
	manualSnapshots := store.GetOldSnapshots(Manual, cfg.Snapshots.Retention.Manual)

	// Combine snapshots to delete
	toDelete := append(autoSnapshots, manualSnapshots...)

	// Delete each snapshot
	deletedCount := 0
	for _, snap := range toDelete {
		snapshotPath := filepath.Join(snapshotDir, snap.File)
		if err := os.Remove(snapshotPath); err != nil {
			// Log error but continue
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Warning: failed to delete snapshot %s: %v\n", snap.File, err)
			}
		} else {
			deletedCount++
		}

		// Remove from metadata
		store.RemoveSnapshot(snap.File)
	}

	// Save updated metadata
	if err := SaveMetadata(metadataPath, store); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// DeleteSnapshot deletes a specific snapshot
func (m *Manager) DeleteSnapshot(snapshotPath string) error {
	// Expand ~ in path
	snapshotPath = expandPath(snapshotPath)

	// Get snapshot directory
	snapshotDir := expandPath(m.Config.Snapshots.Directory)

	// Extract filename
	filename := filepath.Base(snapshotPath)

	// If only filename provided, construct full path
	if !filepath.IsAbs(snapshotPath) && !strings.Contains(snapshotPath, string(os.PathSeparator)) {
		snapshotPath = filepath.Join(snapshotDir, filename)
	}

	// Check if snapshot exists
	if _, err := os.Stat(snapshotPath); err != nil {
		return fmt.Errorf("snapshot not found: %w", err)
	}

	// Delete file
	if err := os.Remove(snapshotPath); err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	// Update metadata
	metadataPath := filepath.Join(snapshotDir, "metadata.json")
	store, err := LoadMetadata(metadataPath)
	if err != nil {
		// Log error but don't fail if metadata can't be loaded
		fmt.Fprintf(os.Stderr, "Warning: failed to update metadata: %v\n", err)
		return nil
	}

	store.RemoveSnapshot(filename)

	if err := SaveMetadata(metadataPath, store); err != nil {
		// Log error but don't fail
		fmt.Fprintf(os.Stderr, "Warning: failed to save metadata: %v\n", err)
	}

	return nil
}

// compressFile compresses a file using gzip
func compressFile(srcPath, dstPath string) error {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	// Copy data
	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}

	return nil
}

// decompressFile decompresses a gzipped file
func decompressFile(srcPath, dstPath string) error {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create gzip reader
	gzReader, err := gzip.NewReader(srcFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy data
	if _, err := io.Copy(dstFile, gzReader); err != nil {
		return fmt.Errorf("failed to decompress data: %w", err)
	}

	return nil
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			if len(path) == 1 {
				return home
			}
			return filepath.Join(home, path[1:])
		}
	}
	return path
}

// getCreatedByContext returns a descriptive string for what created the snapshot
func getCreatedByContext(snapshotType string) string {
	if snapshotType == "auto" {
		return "automatic snapshot"
	}
	return "manual snapshot"
}

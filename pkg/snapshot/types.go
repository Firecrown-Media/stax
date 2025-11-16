package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SnapshotType represents the type of snapshot
type SnapshotType string

const (
	// Auto represents an automatic snapshot
	Auto SnapshotType = "auto"
	// Manual represents a manual snapshot
	Manual SnapshotType = "manual"
)

// SnapshotMetadata represents metadata about a database snapshot
type SnapshotMetadata struct {
	File        string       `json:"file"`
	Project     string       `json:"project"`
	Timestamp   time.Time    `json:"timestamp"`
	Type        SnapshotType `json:"type"`
	Size        int64        `json:"size"`
	Description string       `json:"description,omitempty"`
	CreatedBy   string       `json:"created_by"`
}

// MetadataStore represents the entire metadata storage
type MetadataStore struct {
	Snapshots []SnapshotMetadata `json:"snapshots"`
}

// LoadMetadata loads snapshot metadata from disk
func LoadMetadata(metadataPath string) (*MetadataStore, error) {
	// Check if file exists
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		// Return empty store if file doesn't exist
		return &MetadataStore{
			Snapshots: []SnapshotMetadata{},
		}, nil
	}

	// Read file
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	// Parse JSON
	var store MetadataStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &store, nil
}

// SaveMetadata saves snapshot metadata to disk
func SaveMetadata(metadataPath string, store *MetadataStore) error {
	// Ensure directory exists
	dir := filepath.Dir(metadataPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write file
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// AddSnapshot adds a snapshot to the metadata store
func (s *MetadataStore) AddSnapshot(metadata SnapshotMetadata) {
	s.Snapshots = append(s.Snapshots, metadata)
}

// RemoveSnapshot removes a snapshot from the metadata store by filename
func (s *MetadataStore) RemoveSnapshot(filename string) {
	filtered := make([]SnapshotMetadata, 0, len(s.Snapshots))
	for _, snap := range s.Snapshots {
		if snap.File != filename {
			filtered = append(filtered, snap)
		}
	}
	s.Snapshots = filtered
}

// GetSnapshot retrieves snapshot metadata by filename
func (s *MetadataStore) GetSnapshot(filename string) (*SnapshotMetadata, bool) {
	for _, snap := range s.Snapshots {
		if snap.File == filename {
			return &snap, true
		}
	}
	return nil, false
}

// GetSnapshotsByProject retrieves all snapshots for a specific project
func (s *MetadataStore) GetSnapshotsByProject(projectName string) []SnapshotMetadata {
	filtered := make([]SnapshotMetadata, 0)
	for _, snap := range s.Snapshots {
		if snap.Project == projectName {
			filtered = append(filtered, snap)
		}
	}
	return filtered
}

// GetOldSnapshots returns snapshots older than the specified age in days
func (s *MetadataStore) GetOldSnapshots(snapshotType SnapshotType, maxAgeDays int) []SnapshotMetadata {
	cutoffDate := time.Now().AddDate(0, 0, -maxAgeDays)
	filtered := make([]SnapshotMetadata, 0)

	for _, snap := range s.Snapshots {
		if snap.Type == snapshotType && snap.Timestamp.Before(cutoffDate) {
			filtered = append(filtered, snap)
		}
	}

	return filtered
}

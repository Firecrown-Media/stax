package snapshot

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/firecrown-media/stax/pkg/config"
)

func TestCompressFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "compress small file",
			content:     "test content",
			expectError: false,
		},
		{
			name:        "compress large file",
			content:     strings.Repeat("test content\n", 1000),
			expectError: false,
		},
		{
			name:        "compress empty file",
			content:     "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create source file
			srcPath := filepath.Join(tmpDir, "source.sql")
			if err := os.WriteFile(srcPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create source file: %v", err)
			}

			// Compress
			dstPath := filepath.Join(tmpDir, "compressed.sql.gz")
			err := compressFile(srcPath, dstPath)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify compressed file exists
			if !tt.expectError {
				if _, err := os.Stat(dstPath); err != nil {
					t.Errorf("compressed file not created: %v", err)
				}

				// Verify it's a valid gzip file
				f, err := os.Open(dstPath)
				if err != nil {
					t.Fatalf("failed to open compressed file: %v", err)
				}
				defer f.Close()

				gr, err := gzip.NewReader(f)
				if err != nil {
					t.Errorf("not a valid gzip file: %v", err)
				} else {
					gr.Close()
				}
			}
		})
	}
}

func TestDecompressFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "decompress small file",
			content:     "test content",
			expectError: false,
		},
		{
			name:        "decompress large file",
			content:     strings.Repeat("test content\n", 1000),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create and compress source file
			srcPath := filepath.Join(tmpDir, "source.sql")
			if err := os.WriteFile(srcPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create source file: %v", err)
			}

			compressedPath := filepath.Join(tmpDir, "compressed.sql.gz")
			if err := compressFile(srcPath, compressedPath); err != nil {
				t.Fatalf("failed to compress file: %v", err)
			}

			// Decompress
			dstPath := filepath.Join(tmpDir, "decompressed.sql")
			err := decompressFile(compressedPath, dstPath)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify decompressed content matches original
			if !tt.expectError {
				content, err := os.ReadFile(dstPath)
				if err != nil {
					t.Fatalf("failed to read decompressed file: %v", err)
				}

				if string(content) != tt.content {
					t.Errorf("content mismatch:\nexpected: %s\ngot: %s", tt.content, string(content))
				}
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expand home directory",
			input:    "~/test",
			expected: "", // Will be set to actual home dir in test
		},
		{
			name:     "expand home directory only",
			input:    "~",
			expected: "", // Will be set to actual home dir in test
		},
		{
			name:     "absolute path unchanged",
			input:    "/usr/local/test",
			expected: "/usr/local/test",
		},
		{
			name:     "relative path unchanged",
			input:    "test/path",
			expected: "test/path",
		},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.expected
			if expected == "" {
				if tt.input == "~" {
					expected = home
				} else {
					expected = filepath.Join(home, strings.TrimPrefix(tt.input, "~/"))
				}
			}

			result := expandPath(tt.input)
			if result != expected {
				t.Errorf("expandPath(%s) = %s, expected %s", tt.input, result, expected)
			}
		})
	}
}

func TestGetCreatedByContext(t *testing.T) {
	tests := []struct {
		name         string
		snapshotType string
		expected     string
	}{
		{
			name:         "auto snapshot",
			snapshotType: "auto",
			expected:     "automatic snapshot",
		},
		{
			name:         "manual snapshot",
			snapshotType: "manual",
			expected:     "manual snapshot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCreatedByContext(tt.snapshotType)
			if result != tt.expected {
				t.Errorf("getCreatedByContext(%s) = %s, expected %s", tt.snapshotType, result, tt.expected)
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name: "test-project",
		},
		Snapshots: config.SnapshotsConfig{
			Directory: "~/.stax/snapshots",
		},
	}

	projectDir := "/test/project"

	mgr := NewManager(cfg, projectDir)

	if mgr == nil {
		t.Error("expected manager, got nil")
	}

	if mgr.Config != cfg {
		t.Error("config not set correctly")
	}

	if mgr.ProjectDir != projectDir {
		t.Error("project dir not set correctly")
	}

	if mgr.DDEVManager == nil {
		t.Error("DDEV manager not initialized")
	}
}

func TestDeleteSnapshot(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name: "test-project",
		},
		Snapshots: config.SnapshotsConfig{
			Directory: tmpDir,
		},
	}

	mgr := NewManager(cfg, tmpDir)

	// Create a test snapshot file
	snapshotFile := "test-20250115-120000-auto.sql.gz"
	snapshotPath := filepath.Join(tmpDir, snapshotFile)
	if err := os.WriteFile(snapshotPath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test snapshot: %v", err)
	}

	// Create metadata
	metadataPath := filepath.Join(tmpDir, "metadata.json")
	store := &MetadataStore{
		Snapshots: []SnapshotMetadata{
			{
				File:    snapshotFile,
				Project: "test-project",
			},
		},
	}
	if err := SaveMetadata(metadataPath, store); err != nil {
		t.Fatalf("failed to save metadata: %v", err)
	}

	// Delete snapshot
	err := mgr.DeleteSnapshot(snapshotFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify file was deleted
	if _, err := os.Stat(snapshotPath); !os.IsNotExist(err) {
		t.Error("snapshot file still exists")
	}

	// Verify metadata was updated
	store, err = LoadMetadata(metadataPath)
	if err != nil {
		t.Errorf("failed to load metadata: %v", err)
	}

	if len(store.Snapshots) != 0 {
		t.Errorf("expected 0 snapshots in metadata, got %d", len(store.Snapshots))
	}
}

func TestListSnapshots(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name: "test-project",
		},
		Snapshots: config.SnapshotsConfig{
			Directory: tmpDir,
		},
	}

	mgr := NewManager(cfg, tmpDir)

	// Create metadata with snapshots
	metadataPath := filepath.Join(tmpDir, "metadata.json")
	store := &MetadataStore{
		Snapshots: []SnapshotMetadata{
			{
				File:    "test-project-20250115-120000-auto.sql.gz",
				Project: "test-project",
				Type:    Auto,
			},
			{
				File:    "test-project-20250115-130000-manual.sql.gz",
				Project: "test-project",
				Type:    Manual,
			},
			{
				File:    "other-project-20250115-120000-auto.sql.gz",
				Project: "other-project",
				Type:    Auto,
			},
		},
	}
	if err := SaveMetadata(metadataPath, store); err != nil {
		t.Fatalf("failed to save metadata: %v", err)
	}

	tests := []struct {
		name          string
		projectName   string
		expectedCount int
	}{
		{
			name:          "list snapshots for specific project",
			projectName:   "test-project",
			expectedCount: 2,
		},
		{
			name:          "list all snapshots",
			projectName:   "",
			expectedCount: 3,
		},
		{
			name:          "list snapshots for nonexistent project",
			projectName:   "nonexistent",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshots, err := mgr.ListSnapshots(tt.projectName)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(snapshots) != tt.expectedCount {
				t.Errorf("expected %d snapshots, got %d", tt.expectedCount, len(snapshots))
			}
		})
	}
}

func TestRestoreSnapshot(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name: "test-project",
		},
		Snapshots: config.SnapshotsConfig{
			Directory: tmpDir,
		},
	}

	mgr := NewManager(cfg, tmpDir)

	// Create a compressed test snapshot
	testContent := "CREATE TABLE test (id INT);"
	tmpFile := filepath.Join(tmpDir, "test.sql")
	if err := os.WriteFile(tmpFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	snapshotPath := filepath.Join(tmpDir, "test-snapshot.sql.gz")
	if err := compressFile(tmpFile, snapshotPath); err != nil {
		t.Fatalf("failed to compress test file: %v", err)
	}

	// Test restore (will fail because DDEV is not running, but we can test decompression)
	// We just verify the method doesn't panic with a valid compressed file
	err := mgr.RestoreSnapshot(snapshotPath)
	// We expect an error because DDEV won't be running in tests
	if err == nil {
		t.Log("Note: RestoreSnapshot succeeded unexpectedly (DDEV may be running)")
	}
}

func TestCleanSnapshots(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name: "test-project",
		},
		Snapshots: config.SnapshotsConfig{
			Directory: tmpDir,
			Retention: config.RetentionConfig{
				Auto:   7,
				Manual: 30,
			},
		},
	}

	mgr := NewManager(cfg, tmpDir)

	// Create test snapshot files
	oldAutoFile := filepath.Join(tmpDir, "old-auto.sql.gz")
	recentAutoFile := filepath.Join(tmpDir, "recent-auto.sql.gz")
	oldManualFile := filepath.Join(tmpDir, "old-manual.sql.gz")

	for _, file := range []string{oldAutoFile, recentAutoFile, oldManualFile} {
		if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Create metadata with old and recent snapshots
	metadataPath := filepath.Join(tmpDir, "metadata.json")
	store := &MetadataStore{
		Snapshots: []SnapshotMetadata{
			{
				File:      "old-auto.sql.gz",
				Project:   "test-project",
				Timestamp: time.Now().AddDate(0, 0, -10), // 10 days old
				Type:      Auto,
			},
			{
				File:      "recent-auto.sql.gz",
				Project:   "test-project",
				Timestamp: time.Now().AddDate(0, 0, -3), // 3 days old
				Type:      Auto,
			},
			{
				File:      "old-manual.sql.gz",
				Project:   "test-project",
				Timestamp: time.Now().AddDate(0, 0, -40), // 40 days old
				Type:      Manual,
			},
		},
	}
	if err := SaveMetadata(metadataPath, store); err != nil {
		t.Fatalf("failed to save metadata: %v", err)
	}

	// Clean snapshots
	err := mgr.CleanSnapshots(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify old snapshots were deleted
	if _, err := os.Stat(oldAutoFile); !os.IsNotExist(err) {
		t.Error("old auto snapshot should be deleted")
	}

	if _, err := os.Stat(oldManualFile); !os.IsNotExist(err) {
		t.Error("old manual snapshot should be deleted")
	}

	// Verify recent snapshot still exists
	if _, err := os.Stat(recentAutoFile); err != nil {
		t.Error("recent auto snapshot should still exist")
	}

	// Verify metadata was updated
	store, err = LoadMetadata(metadataPath)
	if err != nil {
		t.Errorf("failed to load metadata: %v", err)
	}

	if len(store.Snapshots) != 1 {
		t.Errorf("expected 1 snapshot in metadata, got %d", len(store.Snapshots))
	}
}

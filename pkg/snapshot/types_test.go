package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadMetadata(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(string) error
		expectError bool
		expectCount int
	}{
		{
			name: "load existing metadata",
			setupFunc: func(path string) error {
				store := &MetadataStore{
					Snapshots: []SnapshotMetadata{
						{
							File:      "test-20250115-120000-auto.sql.gz",
							Project:   "test",
							Timestamp: time.Now(),
							Type:      Auto,
							Size:      1024,
							CreatedBy: "test",
						},
					},
				}
				return SaveMetadata(path, store)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name:        "load non-existent metadata",
			setupFunc:   func(path string) error { return nil },
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			metadataPath := filepath.Join(tmpDir, "metadata.json")

			if err := tt.setupFunc(metadataPath); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			store, err := LoadMetadata(metadataPath)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && len(store.Snapshots) != tt.expectCount {
				t.Errorf("expected %d snapshots, got %d", tt.expectCount, len(store.Snapshots))
			}
		})
	}
}

func TestSaveMetadata(t *testing.T) {
	tests := []struct {
		name        string
		store       *MetadataStore
		expectError bool
	}{
		{
			name: "save valid metadata",
			store: &MetadataStore{
				Snapshots: []SnapshotMetadata{
					{
						File:      "test-20250115-120000-auto.sql.gz",
						Project:   "test",
						Timestamp: time.Now(),
						Type:      Auto,
						Size:      1024,
						CreatedBy: "test",
					},
				},
			},
			expectError: false,
		},
		{
			name: "save empty metadata",
			store: &MetadataStore{
				Snapshots: []SnapshotMetadata{},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			metadataPath := filepath.Join(tmpDir, "metadata.json")

			err := SaveMetadata(metadataPath, tt.store)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify file was created
			if !tt.expectError {
				if _, err := os.Stat(metadataPath); err != nil {
					t.Errorf("metadata file not created: %v", err)
				}
			}
		})
	}
}

func TestAddSnapshot(t *testing.T) {
	store := &MetadataStore{
		Snapshots: []SnapshotMetadata{},
	}

	snapshot := SnapshotMetadata{
		File:      "test-20250115-120000-auto.sql.gz",
		Project:   "test",
		Timestamp: time.Now(),
		Type:      Auto,
		Size:      1024,
		CreatedBy: "test",
	}

	store.AddSnapshot(snapshot)

	if len(store.Snapshots) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(store.Snapshots))
	}

	if store.Snapshots[0].File != snapshot.File {
		t.Errorf("expected file %s, got %s", snapshot.File, store.Snapshots[0].File)
	}
}

func TestRemoveSnapshot(t *testing.T) {
	store := &MetadataStore{
		Snapshots: []SnapshotMetadata{
			{
				File:      "test1-20250115-120000-auto.sql.gz",
				Project:   "test1",
				Timestamp: time.Now(),
				Type:      Auto,
				Size:      1024,
				CreatedBy: "test",
			},
			{
				File:      "test2-20250115-120000-auto.sql.gz",
				Project:   "test2",
				Timestamp: time.Now(),
				Type:      Auto,
				Size:      2048,
				CreatedBy: "test",
			},
		},
	}

	store.RemoveSnapshot("test1-20250115-120000-auto.sql.gz")

	if len(store.Snapshots) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(store.Snapshots))
	}

	if store.Snapshots[0].File != "test2-20250115-120000-auto.sql.gz" {
		t.Errorf("wrong snapshot removed")
	}
}

func TestGetSnapshot(t *testing.T) {
	store := &MetadataStore{
		Snapshots: []SnapshotMetadata{
			{
				File:      "test-20250115-120000-auto.sql.gz",
				Project:   "test",
				Timestamp: time.Now(),
				Type:      Auto,
				Size:      1024,
				CreatedBy: "test",
			},
		},
	}

	tests := []struct {
		name        string
		filename    string
		expectFound bool
	}{
		{
			name:        "find existing snapshot",
			filename:    "test-20250115-120000-auto.sql.gz",
			expectFound: true,
		},
		{
			name:        "find non-existent snapshot",
			filename:    "nonexistent.sql.gz",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snap, found := store.GetSnapshot(tt.filename)

			if found != tt.expectFound {
				t.Errorf("expected found=%v, got %v", tt.expectFound, found)
			}

			if tt.expectFound && snap == nil {
				t.Error("expected snapshot, got nil")
			}

			if !tt.expectFound && snap != nil {
				t.Error("expected nil, got snapshot")
			}
		})
	}
}

func TestGetSnapshotsByProject(t *testing.T) {
	store := &MetadataStore{
		Snapshots: []SnapshotMetadata{
			{
				File:      "test1-20250115-120000-auto.sql.gz",
				Project:   "test1",
				Timestamp: time.Now(),
				Type:      Auto,
				Size:      1024,
				CreatedBy: "test",
			},
			{
				File:      "test2-20250115-120000-auto.sql.gz",
				Project:   "test2",
				Timestamp: time.Now(),
				Type:      Auto,
				Size:      2048,
				CreatedBy: "test",
			},
			{
				File:      "test1-20250115-130000-manual.sql.gz",
				Project:   "test1",
				Timestamp: time.Now(),
				Type:      Manual,
				Size:      1536,
				CreatedBy: "test",
			},
		},
	}

	tests := []struct {
		name          string
		projectName   string
		expectedCount int
	}{
		{
			name:          "get snapshots for test1",
			projectName:   "test1",
			expectedCount: 2,
		},
		{
			name:          "get snapshots for test2",
			projectName:   "test2",
			expectedCount: 1,
		},
		{
			name:          "get snapshots for nonexistent project",
			projectName:   "nonexistent",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshots := store.GetSnapshotsByProject(tt.projectName)

			if len(snapshots) != tt.expectedCount {
				t.Errorf("expected %d snapshots, got %d", tt.expectedCount, len(snapshots))
			}
		})
	}
}

func TestGetOldSnapshots(t *testing.T) {
	now := time.Now()
	oldTime := now.AddDate(0, 0, -10) // 10 days ago
	recentTime := now.AddDate(0, 0, -3) // 3 days ago

	store := &MetadataStore{
		Snapshots: []SnapshotMetadata{
			{
				File:      "old-auto.sql.gz",
				Project:   "test",
				Timestamp: oldTime,
				Type:      Auto,
				Size:      1024,
				CreatedBy: "test",
			},
			{
				File:      "recent-auto.sql.gz",
				Project:   "test",
				Timestamp: recentTime,
				Type:      Auto,
				Size:      1024,
				CreatedBy: "test",
			},
			{
				File:      "old-manual.sql.gz",
				Project:   "test",
				Timestamp: oldTime,
				Type:      Manual,
				Size:      1024,
				CreatedBy: "test",
			},
		},
	}

	tests := []struct {
		name          string
		snapshotType  SnapshotType
		maxAgeDays    int
		expectedCount int
	}{
		{
			name:          "get old auto snapshots (7 days)",
			snapshotType:  Auto,
			maxAgeDays:    7,
			expectedCount: 1,
		},
		{
			name:          "get old manual snapshots (7 days)",
			snapshotType:  Manual,
			maxAgeDays:    7,
			expectedCount: 1,
		},
		{
			name:          "get old auto snapshots (2 days)",
			snapshotType:  Auto,
			maxAgeDays:    2,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshots := store.GetOldSnapshots(tt.snapshotType, tt.maxAgeDays)

			if len(snapshots) != tt.expectedCount {
				t.Errorf("expected %d old snapshots, got %d", tt.expectedCount, len(snapshots))
			}
		})
	}
}

package wpengine

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestCalculateMD5(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	content := []byte("Hello, World!")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate checksum
	checksum, err := calculateMD5(tmpFile)
	if err != nil {
		t.Fatalf("Failed to calculate MD5: %v", err)
	}

	// Verify checksum
	hash := md5.New()
	hash.Write(content)
	expected := hex.EncodeToString(hash.Sum(nil))

	if checksum != expected {
		t.Errorf("Expected checksum %s, got %s", expected, checksum)
	}
}

func TestGenerateLocalChecksums(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	files := map[string]string{
		"file1.txt":               "content1",
		"file2.txt":               "content2",
		"subdir/file3.txt":        "content3",
		"subdir/nested/file4.txt": "content4",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Generate checksums
	checksums, err := GenerateLocalChecksums(tmpDir)
	if err != nil {
		t.Fatalf("Failed to generate checksums: %v", err)
	}

	// Verify we found all files
	if len(checksums) != len(files) {
		t.Errorf("Expected %d checksums, got %d", len(files), len(checksums))
	}

	// Verify each file has a checksum
	for path := range files {
		if _, exists := checksums[path]; !exists {
			t.Errorf("Missing checksum for file: %s", path)
		}
	}
}

func TestVerifyChecksums(t *testing.T) {
	tests := []struct {
		name              string
		remote            map[string]string
		local             map[string]string
		wantMatched       int
		wantMismatched    int
		wantMissingLocal  int
		wantMissingRemote int
	}{
		{
			name: "all files match",
			remote: map[string]string{
				"file1.txt": "abc123",
				"file2.txt": "def456",
			},
			local: map[string]string{
				"file1.txt": "abc123",
				"file2.txt": "def456",
			},
			wantMatched:       2,
			wantMismatched:    0,
			wantMissingLocal:  0,
			wantMissingRemote: 0,
		},
		{
			name: "mismatched checksums",
			remote: map[string]string{
				"file1.txt": "abc123",
				"file2.txt": "def456",
			},
			local: map[string]string{
				"file1.txt": "abc123",
				"file2.txt": "different",
			},
			wantMatched:       1,
			wantMismatched:    1,
			wantMissingLocal:  0,
			wantMissingRemote: 0,
		},
		{
			name: "missing local files",
			remote: map[string]string{
				"file1.txt": "abc123",
				"file2.txt": "def456",
			},
			local: map[string]string{
				"file1.txt": "abc123",
			},
			wantMatched:       1,
			wantMismatched:    0,
			wantMissingLocal:  1,
			wantMissingRemote: 0,
		},
		{
			name: "missing remote files",
			remote: map[string]string{
				"file1.txt": "abc123",
			},
			local: map[string]string{
				"file1.txt": "abc123",
				"file2.txt": "def456",
			},
			wantMatched:       1,
			wantMismatched:    0,
			wantMissingLocal:  0,
			wantMissingRemote: 1,
		},
		{
			name: "mixed scenario",
			remote: map[string]string{
				"file1.txt": "abc123",
				"file2.txt": "def456",
				"file3.txt": "ghi789",
			},
			local: map[string]string{
				"file1.txt": "abc123",
				"file2.txt": "different",
				"file4.txt": "jkl012",
			},
			wantMatched:       1,
			wantMismatched:    1,
			wantMissingLocal:  1,
			wantMissingRemote: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyChecksums(tt.remote, tt.local)

			if result.MatchedFiles != tt.wantMatched {
				t.Errorf("MatchedFiles = %d, want %d", result.MatchedFiles, tt.wantMatched)
			}
			if result.MismatchedFiles != tt.wantMismatched {
				t.Errorf("MismatchedFiles = %d, want %d", result.MismatchedFiles, tt.wantMismatched)
			}
			if result.MissingLocal != tt.wantMissingLocal {
				t.Errorf("MissingLocal = %d, want %d", result.MissingLocal, tt.wantMissingLocal)
			}
			if result.MissingRemote != tt.wantMissingRemote {
				t.Errorf("MissingRemote = %d, want %d", result.MissingRemote, tt.wantMissingRemote)
			}
		})
	}
}

func TestGenerateLocalChecksumsNonExistentPath(t *testing.T) {
	_, err := GenerateLocalChecksums("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
}

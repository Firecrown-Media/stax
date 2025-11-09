package wordpress

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// DBConfig represents database configuration
type DBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

// SnapshotInfo represents information about a database snapshot
type SnapshotInfo struct {
	Name        string
	Path        string
	Size        int64
	Created     time.Time
	Description string
}

// ImportDatabaseOptions represents options for database import
type ImportDatabaseOptions struct {
	SourceFile string
	Replace    bool
}

// ExportDatabaseOptions represents options for database export
type ExportDatabaseOptions struct {
	Destination    string
	ExcludeTables  []string
	SkipLogs       bool
	SkipTransients bool
	Compress       bool
}

// ImportDatabase imports a SQL file into the database
func (c *CLI) ImportDatabaseWithOptions(options ImportDatabaseOptions) error {
	return c.ImportDatabase(options.SourceFile)
}

// ExportDatabaseWithOptions exports the database with options
func (c *CLI) ExportDatabaseWithOptions(options ExportDatabaseOptions) error {
	args := []string{"db", "export"}

	// Add destination
	if options.Destination != "" {
		args = append(args, options.Destination)
	}

	// Exclude tables
	if len(options.ExcludeTables) > 0 {
		excludePattern := joinStrings(options.ExcludeTables, ",")
		args = append(args, "--exclude_tables="+excludePattern)
	}

	return c.Execute(args...)
}

// CreateSnapshot creates a database snapshot
func CreateSnapshot(cli *CLI, name string, snapshotDir string) (*SnapshotInfo, error) {
	// Create snapshot directory if it doesn't exist
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// Generate snapshot filename
	if name == "" {
		name = fmt.Sprintf("db_%s", time.Now().Format("2006-01-02_15-04-05"))
	}
	filename := filepath.Join(snapshotDir, name+".sql.gz")

	// Export database
	if err := cli.ExportDatabase(filename); err != nil {
		return nil, fmt.Errorf("failed to export database: %w", err)
	}

	// Get file info
	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot info: %w", err)
	}

	return &SnapshotInfo{
		Name:    name,
		Path:    filename,
		Size:    info.Size(),
		Created: time.Now(),
	}, nil
}

// RestoreSnapshot restores a database from a snapshot
func RestoreSnapshot(cli *CLI, snapshotPath string) error {
	// Check if snapshot exists
	if _, err := os.Stat(snapshotPath); err != nil {
		return fmt.Errorf("snapshot not found: %w", err)
	}

	// Import snapshot
	if err := cli.ImportDatabase(snapshotPath); err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	return nil
}

// ListSnapshots lists all database snapshots in a directory
func ListSnapshots(snapshotDir string) ([]SnapshotInfo, error) {
	// Check if directory exists
	if _, err := os.Stat(snapshotDir); os.IsNotExist(err) {
		return []SnapshotInfo{}, nil
	}

	// Read directory
	files, err := os.ReadDir(snapshotDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot directory: %w", err)
	}

	var snapshots []SnapshotInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only include .sql and .sql.gz files
		name := file.Name()
		if filepath.Ext(name) != ".sql" && filepath.Ext(name) != ".gz" {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		snapshots = append(snapshots, SnapshotInfo{
			Name:    name,
			Path:    filepath.Join(snapshotDir, name),
			Size:    info.Size(),
			Created: info.ModTime(),
		})
	}

	return snapshots, nil
}

// StreamDatabaseImport streams a database import from a reader
func StreamDatabaseImport(cli *CLI, reader io.Reader) error {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "stax-import-*.sql")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Copy reader to temp file
	if _, err := io.Copy(tmpFile, reader); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Import temp file
	return cli.ImportDatabase(tmpFile.Name())
}

// GetDatabaseSize gets the approximate size of the database
func (c *CLI) GetDatabaseSize() (int64, error) {
	query := `SELECT SUM(data_length + index_length) as size
	          FROM information_schema.TABLES
	          WHERE table_schema = DATABASE()`

	output, err := c.Query(query)
	if err != nil {
		return 0, err
	}

	var size int64
	if _, err := fmt.Sscanf(output, "%d", &size); err != nil {
		return 0, fmt.Errorf("failed to parse database size: %w", err)
	}

	return size, nil
}

// GetTableCount gets the number of tables in the database
func (c *CLI) GetTableCount() (int, error) {
	query := `SELECT COUNT(*) FROM information_schema.TABLES WHERE table_schema = DATABASE()`

	output, err := c.Query(query)
	if err != nil {
		return 0, err
	}

	var count int
	if _, err := fmt.Sscanf(output, "%d", &count); err != nil {
		return 0, fmt.Errorf("failed to parse table count: %w", err)
	}

	return count, nil
}

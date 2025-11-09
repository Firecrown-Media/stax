package build

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches files for changes and triggers rebuilds
type Watcher struct {
	projectPath string
	watcher     *fsnotify.Watcher
}

// NewWatcher creates a new file watcher
func NewWatcher(projectPath string) *Watcher {
	return &Watcher{
		projectPath: projectPath,
	}
}

// Watch starts watching for file changes
func (w *Watcher) Watch(callback func()) error {
	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer w.watcher.Close()

	// Add paths to watch
	paths := w.getWatchPaths()
	for _, path := range paths {
		if err := w.addRecursive(path); err != nil {
			log.Printf("Warning: failed to watch %s: %v\n", path, err)
		}
	}

	// Debounce settings
	debounceTime := 500 * time.Millisecond
	timer := time.NewTimer(0)
	<-timer.C // Drain the initial timer

	fmt.Println("Watching for file changes... (Press Ctrl+C to stop)")

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			// Ignore certain events
			if w.shouldIgnore(event.Name) {
				continue
			}

			// Only trigger on write/create/remove
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Remove == fsnotify.Remove {

				// Reset debounce timer
				timer.Reset(debounceTime)
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher error: %v\n", err)

		case <-timer.C:
			// Debounce period elapsed, trigger callback
			fmt.Println("\nFile changes detected, rebuilding...")
			callback()
		}
	}
}

// WatchTheme watches a specific theme for changes
func (w *Watcher) WatchTheme(themeName string, callback func()) error {
	themePath := filepath.Join(w.projectPath, "wp-content", "themes", themeName)

	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer w.watcher.Close()

	// Watch theme src directory
	srcPath := filepath.Join(themePath, "src")
	if err := w.addRecursive(srcPath); err != nil {
		return err
	}

	debounceTime := 500 * time.Millisecond
	timer := time.NewTimer(0)
	<-timer.C

	fmt.Printf("Watching theme '%s' for changes...\n", themeName)

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			if w.shouldIgnore(event.Name) {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				timer.Reset(debounceTime)
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher error: %v\n", err)

		case <-timer.C:
			fmt.Printf("\nChanges detected in %s, rebuilding...\n", themeName)
			callback()
		}
	}
}

// WatchPlugin watches a specific plugin for changes
func (w *Watcher) WatchPlugin(pluginPath string, callback func()) error {
	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer w.watcher.Close()

	// Watch plugin src directory
	srcPath := filepath.Join(pluginPath, "src")
	if err := w.addRecursive(srcPath); err != nil {
		return err
	}

	debounceTime := 500 * time.Millisecond
	timer := time.NewTimer(0)
	<-timer.C

	fmt.Printf("Watching plugin for changes...\n")

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			if w.shouldIgnore(event.Name) {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				timer.Reset(debounceTime)
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher error: %v\n", err)

		case <-timer.C:
			fmt.Println("\nPlugin changes detected, rebuilding...")
			callback()
		}
	}
}

// getWatchPaths returns the paths that should be watched
func (w *Watcher) getWatchPaths() []string {
	return []string{
		filepath.Join(w.projectPath, "wp-content", "mu-plugins", "firecrown", "src"),
		filepath.Join(w.projectPath, "wp-content", "themes", "firecrown-parent", "src"),
		filepath.Join(w.projectPath, "wp-content", "themes", "firecrown-child", "src"),
	}
}

// getIgnorePatterns returns patterns to ignore
func (w *Watcher) getIgnorePatterns() []string {
	return []string{
		"node_modules",
		"vendor",
		"build",
		".git",
		".stax",
		".DS_Store",
		"*.swp",
		"*.swo",
		"*~",
		".npm-start.log",
		".npm-start.pid",
	}
}

// shouldIgnore checks if a path should be ignored
func (w *Watcher) shouldIgnore(path string) bool {
	ignorePatterns := w.getIgnorePatterns()

	for _, pattern := range ignorePatterns {
		// Simple pattern matching
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}

		// Check if any part of the path contains the pattern
		if containsPath(path, pattern) {
			return true
		}
	}

	return false
}

// addRecursive adds a directory and all its subdirectories to the watcher
func (w *Watcher) addRecursive(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Path doesn't exist, skip silently
	}

	// Add the path itself
	if err := w.watcher.Add(path); err != nil {
		return err
	}

	// Walk the directory tree
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip files
		if !info.IsDir() {
			return nil
		}

		// Skip ignored directories
		if w.shouldIgnore(walkPath) {
			return filepath.SkipDir
		}

		// Add directory to watcher
		return w.watcher.Add(walkPath)
	})
}

// containsPath checks if a path contains a specific component
func containsPath(path, component string) bool {
	parts := filepath.SplitList(path)
	for _, part := range parts {
		if part == component {
			return true
		}
	}
	// Also check with filepath separator
	return filepath.Base(path) == component ||
		filepath.Dir(path) == component ||
		strings.Contains(path, string(filepath.Separator)+component+string(filepath.Separator))
}

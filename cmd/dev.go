package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/firecrown-media/stax/pkg/build"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	devTheme  string
	devWatch  bool
	devNoOpen bool
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "✓ Start development mode",
	Long: `Start development mode with file watching and auto-rebuild.

This runs 'npm start' in the specified theme directory, which typically
includes Hot Module Reloading (HMR) for faster development.

The process runs in the foreground. Press Ctrl+C to stop.`,
	Example: `  # Start dev mode for parent theme (default)
  stax dev

  # Start dev mode for specific theme
  stax dev --theme=firecrown-child

  # Start dev mode without file watching
  stax dev --no-watch`,
	RunE: runDev,
}

var devThemeCmd = &cobra.Command{
	Use:   "theme [theme-name]",
	Short: "Start dev mode for a specific theme",
	Long: `Start development mode for a specific theme.

Available themes:
  - firecrown-parent (default)
  - firecrown-child`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDevTheme,
}

var devStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop development mode",
	Long: `Stop any running development processes.

This looks for background npm start processes and stops them.`,
	RunE: runDevStop,
}

var devWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for file changes and rebuild",
	Long: `Watch source files for changes and automatically rebuild.

This watches:
  - Theme src directories
  - MU plugin src directory

When changes are detected, the appropriate build step is triggered.`,
	RunE: runDevWatch,
}

func init() {
	rootCmd.AddCommand(devCmd)

	// Subcommands
	devCmd.AddCommand(devThemeCmd)
	devCmd.AddCommand(devStopCmd)
	devCmd.AddCommand(devWatchCmd)

	// Flags
	devCmd.Flags().StringVar(&devTheme, "theme", "firecrown-parent", "theme to develop")
	devCmd.Flags().BoolVar(&devWatch, "watch", true, "enable file watching")
	devCmd.Flags().BoolVar(&devNoOpen, "no-open", false, "don't open browser")
}

func runDev(cmd *cobra.Command, args []string) error {
	ui.PrintHeader(fmt.Sprintf("Starting Development Mode: %s", devTheme))

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	themePath := filepath.Join(projectDir, "wp-content", "themes", devTheme)

	// Check if theme exists
	if _, err := os.Stat(filepath.Join(themePath, "package.json")); os.IsNotExist(err) {
		return fmt.Errorf("theme not found or has no package.json: %s", devTheme)
	}

	npm := build.NewNPM(themePath)

	// Check if npm start script exists
	scripts, err := npm.ListScripts()
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	if _, exists := scripts["start"]; !exists {
		return fmt.Errorf("no 'start' script found in package.json")
	}

	ui.Info("Running npm start in %s", devTheme)
	ui.Info("Press Ctrl+C to stop\n")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run npm start
	go func() {
		options := build.NPMOptions{
			WorkingDir: themePath,
			Verbose:    verbose,
		}
		if err := npm.Start(false, options); err != nil {
			ui.Error("npm start failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	ui.Info("\nStopping development mode...")

	return nil
}

func runDevTheme(cmd *cobra.Command, args []string) error {
	themeName := devTheme
	if len(args) > 0 {
		themeName = args[0]
	}

	// Set the theme and run main dev command
	devTheme = themeName
	return runDev(cmd, args)
}

func runDevStop(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Stopping Development Mode")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	themes := []string{"firecrown-parent", "firecrown-child"}
	stopped := 0

	for _, theme := range themes {
		themePath := filepath.Join(projectDir, "wp-content", "themes", theme)
		npm := build.NewNPM(themePath)

		pidFile := filepath.Join(themePath, ".npm-start.pid")
		if _, err := os.Stat(pidFile); err == nil {
			ui.Info("Stopping %s...", theme)
			if err := npm.StopBackground(); err != nil {
				ui.Warning("Failed to stop %s: %v", theme, err)
			} else {
				ui.Success("Stopped %s", theme)
				stopped++
			}
		}
	}

	if stopped == 0 {
		ui.Info("No development processes running")
	} else {
		ui.Success("Stopped %d process(es)", stopped)
	}

	return nil
}

func runDevWatch(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Watching for File Changes")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	mgr := build.NewManager(projectDir)
	mgr.SetVerbose(verbose)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create rebuild callback
	rebuildCallback := func() {
		ui.Info("Rebuilding...")
		if _, err := mgr.RunBuildScript(); err != nil {
			ui.Error("Build failed: %v", err)
		} else {
			ui.Success("Build complete!")
		}
	}

	// Start watching in goroutine
	go func() {
		if err := mgr.WatchForChanges(rebuildCallback); err != nil {
			ui.Error("Watch error: %v", err)
		}
	}()

	ui.Info("Watching for changes... (Press Ctrl+C to stop)")

	// Wait for interrupt
	<-sigChan
	ui.Info("\nStopping file watcher...")

	return nil
}

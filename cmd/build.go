package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/firecrown-media/stax/pkg/build"
	"github.com/firecrown-media/stax/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	buildVerbose bool
	buildForce   bool
	buildClean   bool
	buildTimeout int
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long: `Build the project by running all build scripts.

This command executes scripts/build.sh which orchestrates:
  - Composer install for MU plugins
  - NPM install and build for themes
  - Any custom build steps

The build process preserves Firecrown's existing workflow while
providing better feedback and error handling.`,
	Example: `  # Run full build
  stax build

  # Force rebuild even if not needed
  stax build --force

  # Clean and rebuild
  stax build --clean

  # Verbose output
  stax build --verbose`,
	RunE: runBuild,
}

var buildComposerCmd = &cobra.Command{
	Use:   "composer",
	Short: "Run composer install only",
	Long: `Run composer install for MU plugins and themes.

This installs PHP dependencies without running npm build steps.`,
	RunE: runBuildComposer,
}

var buildNpmCmd = &cobra.Command{
	Use:   "npm",
	Short: "Run npm install and build only",
	Long: `Run npm install and npm run build for themes.

This builds frontend assets without installing PHP dependencies.`,
	RunE: runBuildNPM,
}

var buildThemeCmd = &cobra.Command{
	Use:   "theme [theme-name]",
	Short: "Build a specific theme",
	Long: `Build a specific theme by name.

Available themes:
  - firecrown-parent
  - firecrown-child

This runs npm install && npm run build for the specified theme.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBuildTheme,
}

var buildCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean build artifacts",
	Long: `Remove all build artifacts including:
  - vendor directories (composer)
  - node_modules directories (npm)
  - build output directories

This is useful when you want to start fresh or troubleshoot build issues.`,
	RunE: runBuildClean,
}

var buildStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check build status",
	Long: `Check if a build is needed and display build status information.

This shows:
  - Whether dependencies are installed
  - Whether dependencies need updating
  - Last build time
  - Reasons why a build might be needed`,
	RunE: runBuildStatus,
}

var buildScriptsCmd = &cobra.Command{
	Use:   "scripts",
	Short: "List available build scripts",
	Long: `List all build scripts found in the project.

Shows scripts in scripts/build/ and their execution order.`,
	RunE: runBuildScripts,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Subcommands
	buildCmd.AddCommand(buildComposerCmd)
	buildCmd.AddCommand(buildNpmCmd)
	buildCmd.AddCommand(buildThemeCmd)
	buildCmd.AddCommand(buildCleanCmd)
	buildCmd.AddCommand(buildStatusCmd)
	buildCmd.AddCommand(buildScriptsCmd)

	// Flags
	buildCmd.PersistentFlags().BoolVarP(&buildVerbose, "verbose", "v", false, "verbose output")
	buildCmd.Flags().BoolVar(&buildForce, "force", false, "force rebuild even if not needed")
	buildCmd.Flags().BoolVar(&buildClean, "clean", false, "clean before building")
	buildCmd.Flags().IntVar(&buildTimeout, "timeout", 600, "build timeout in seconds")
}

func runBuild(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Building Project")

	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	mgr := build.NewManager(projectDir)
	mgr.SetVerbose(buildVerbose)

	// Check if build is needed (unless forced)
	if !buildForce {
		ui.Info("Checking build status...")
		status, err := mgr.GetBuildStatus()
		if err != nil {
			ui.Warning("Failed to check build status: %v", err)
		} else if !status.NeedsBuild {
			ui.Success("Build is up to date! (use --force to rebuild anyway)")
			return nil
		} else if len(status.Reasons) > 0 {
			ui.Info("Build needed:")
			for _, reason := range status.Reasons {
				ui.Info("  - %s", reason)
			}
		}
	}

	// Clean if requested
	if buildClean {
		spinner := ui.NewSpinner("Cleaning build artifacts...")
		spinner.Start()
		if err := mgr.Clean(); err != nil {
			spinner.Error("Failed to clean")
			return err
		}
		spinner.Success("Cleaned build artifacts")
	}

	// Run build
	spinner := ui.NewSpinner("Running build scripts...")
	spinner.Start()
	startTime := time.Now()

	result, err := mgr.RunBuildScript()
	duration := time.Since(startTime)

	if err != nil {
		spinner.Error("Build failed")
		ui.Error("Build failed after %s", duration.Round(time.Second))
		if buildVerbose && result != nil {
			ui.Print("\n%s\n", result.Output)
		}
		return err
	}

	spinner.Success(fmt.Sprintf("Build completed in %s", duration.Round(time.Second)))

	// Validate build output
	checker := build.NewStatusChecker(projectDir)
	if err := checker.ValidateBuild(); err != nil {
		ui.Warning("Build validation warnings: %v", err)
	}

	ui.Success("Project built successfully!")

	return nil
}

func runBuildComposer(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Running Composer Install")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	mgr := build.NewManager(projectDir)
	mgr.SetVerbose(buildVerbose)

	spinner := ui.NewSpinner("Installing composer dependencies...")
	spinner.Start()

	options := build.ComposerOptions{
		NoDev:              false,
		IgnorePlatformReqs: true,
		PreferDist:         true,
		Timeout:            buildTimeout,
		Verbose:            buildVerbose,
	}

	if err := mgr.BuildMUPlugins(options); err != nil {
		spinner.Error("Composer install failed")
		return err
	}

	spinner.Success("Composer dependencies installed")
	return nil
}

func runBuildNPM(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Running NPM Install & Build")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	mgr := build.NewManager(projectDir)
	mgr.SetVerbose(buildVerbose)

	spinner := ui.NewSpinner("Building themes...")
	spinner.Start()

	options := build.NPMOptions{
		LegacyPeerDeps: true,
		Timeout:        buildTimeout,
		Verbose:        buildVerbose,
	}

	if err := mgr.BuildThemes(options); err != nil {
		spinner.Error("NPM build failed")
		return err
	}

	spinner.Success("Themes built successfully")
	return nil
}

func runBuildTheme(cmd *cobra.Command, args []string) error {
	themeName := "firecrown-parent"
	if len(args) > 0 {
		themeName = args[0]
	}

	ui.PrintHeader(fmt.Sprintf("Building Theme: %s", themeName))

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	mgr := build.NewManager(projectDir)
	mgr.SetVerbose(buildVerbose)

	spinner := ui.NewSpinner(fmt.Sprintf("Building %s...", themeName))
	spinner.Start()

	options := build.NPMOptions{
		LegacyPeerDeps: true,
		Timeout:        buildTimeout,
		Verbose:        buildVerbose,
	}

	if err := mgr.BuildTheme(themeName, options); err != nil {
		spinner.Error(fmt.Sprintf("Failed to build %s", themeName))
		return err
	}

	spinner.Success(fmt.Sprintf("Theme %s built successfully", themeName))
	return nil
}

func runBuildClean(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Cleaning Build Artifacts")

	if !ui.Confirm("This will remove all build artifacts. Continue?") {
		ui.Info("Cancelled")
		return nil
	}

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	mgr := build.NewManager(projectDir)
	mgr.SetVerbose(buildVerbose)

	spinner := ui.NewSpinner("Removing build artifacts...")
	spinner.Start()

	if err := mgr.Clean(); err != nil {
		spinner.Error("Failed to clean")
		return err
	}

	spinner.Success("Build artifacts cleaned")
	return nil
}

func runBuildStatus(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Build Status")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	mgr := build.NewManager(projectDir)
	status, err := mgr.GetBuildStatus()
	if err != nil {
		return err
	}

	// Display status
	if status.BuildScriptExists {
		ui.Success("Build script found: scripts/build.sh")
	} else {
		ui.Warning("No build script found")
	}

	if len(status.CustomBuildScripts) > 0 {
		ui.Info("Custom build scripts (%d):", len(status.CustomBuildScripts))
		for _, script := range status.CustomBuildScripts {
			ui.Info("  - %s", script)
		}
	}

	// Composer status
	ui.Section("Composer Dependencies")
	if status.ComposerStatus.Installed {
		ui.Success("Installed")
		ui.Info("  Lock file: %s", status.ComposerStatus.LockModified.Format("2006-01-02 15:04:05"))
		ui.Info("  Vendor dir: %s", status.ComposerStatus.VendorModified.Format("2006-01-02 15:04:05"))
	} else {
		ui.Warning("Not installed")
	}

	if status.ComposerStatus.NeedsUpdate {
		ui.Warning("Needs update (composer.json modified)")
	}

	// NPM status
	ui.Section("NPM Dependencies")
	if status.NPMStatus.Installed {
		ui.Success("Installed")
		ui.Info("  Lock file: %s", status.NPMStatus.LockModified.Format("2006-01-02 15:04:05"))
		ui.Info("  node_modules: %s", status.NPMStatus.VendorModified.Format("2006-01-02 15:04:05"))
	} else {
		ui.Warning("Not installed")
	}

	if status.NPMStatus.NeedsUpdate {
		ui.Warning("Needs update (package.json modified)")
	}

	// Last build time
	if !status.LastBuildTime.IsZero() {
		ui.Section("Last Build")
		ui.Info("Time: %s", status.LastBuildTime.Format("2006-01-02 15:04:05"))
		ui.Info("Age: %s ago", time.Since(status.LastBuildTime).Round(time.Second))
	}

	// Build needed?
	ui.Section("Build Status")
	if status.NeedsBuild {
		ui.Warning("Build needed")
		if len(status.Reasons) > 0 {
			ui.Info("Reasons:")
			for _, reason := range status.Reasons {
				ui.Info("  - %s", reason)
			}
		}
	} else {
		ui.Success("Build is up to date")
	}

	return nil
}

func runBuildScripts(cmd *cobra.Command, args []string) error {
	ui.PrintHeader("Build Scripts")

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	mgr := build.NewManager(projectDir)
	scripts, err := mgr.DetectBuildScripts()
	if err != nil {
		return err
	}

	if len(scripts) == 0 {
		ui.Warning("No build scripts found")
		ui.Info("Run 'stax build:generate' to create default build scripts")
		return nil
	}

	ui.Info("Found %d build script(s):\n", len(scripts))

	for _, script := range scripts {
		ui.Info("  [%02d] %s", script.Order, script.Name)
		ui.Verbose("       Type: %s", script.Type)
		ui.Verbose("       Description: %s", script.Description)
		ui.Verbose("       Path: %s", script.Path)
		ui.Info("")
	}

	return nil
}

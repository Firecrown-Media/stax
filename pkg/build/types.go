package build

import "time"

// BuildOptions represents options for build operations
type BuildOptions struct {
	// ProjectPath is the path to the project
	ProjectPath string

	// Verbose enables verbose output
	Verbose bool

	// Force forces a rebuild even if not needed
	Force bool

	// Clean removes build artifacts before building
	Clean bool

	// SkipComposer skips composer install
	SkipComposer bool

	// SkipNPM skips npm install/build
	SkipNPM bool

	// Parallel enables parallel builds (where possible)
	Parallel bool

	// Timeout in seconds (0 = no timeout)
	Timeout int
}

// ComposerOptions represents options for composer operations
type ComposerOptions struct {
	// WorkingDir is the directory to run composer in
	WorkingDir string

	// NoDev skips installing dev dependencies
	NoDev bool

	// NoScripts skips running composer scripts
	NoScripts bool

	// IgnorePlatformReqs ignores platform requirements
	IgnorePlatformReqs bool

	// PreferDist prefers dist packages over source
	PreferDist bool

	// PreferSource prefers source packages over dist
	PreferSource bool

	// Optimize optimizes the autoloader
	Optimize bool

	// Verbose enables verbose output
	Verbose bool

	// Timeout in seconds
	Timeout int
}

// NPMOptions represents options for npm operations
type NPMOptions struct {
	// WorkingDir is the directory to run npm in
	WorkingDir string

	// Production only installs production dependencies
	Production bool

	// LegacyPeerDeps uses legacy peer deps
	LegacyPeerDeps bool

	// Clean removes node_modules before installing
	Clean bool

	// Verbose enables verbose output
	Verbose bool

	// Timeout in seconds
	Timeout int
}

// BuildStatus represents the current build status
type BuildStatus struct {
	// NeedsBuild indicates if a build is needed
	NeedsBuild bool

	// LastBuildTime is the timestamp of the last build
	LastBuildTime time.Time

	// Reasons lists why a build is needed
	Reasons []string

	// ComposerStatus is the status of composer dependencies
	ComposerStatus DependencyStatus

	// NPMStatus is the status of npm dependencies
	NPMStatus DependencyStatus

	// BuildScriptExists indicates if scripts/build.sh exists
	BuildScriptExists bool

	// CustomBuildScripts lists custom build scripts found
	CustomBuildScripts []string
}

// DependencyStatus represents the status of dependencies (composer/npm)
type DependencyStatus struct {
	// Installed indicates if dependencies are installed
	Installed bool

	// ConfigFile is the path to the config file (composer.json/package.json)
	ConfigFile string

	// ConfigExists indicates if the config file exists
	ConfigExists bool

	// ConfigModified is the last modified time of config file
	ConfigModified time.Time

	// LockFile is the path to the lock file (composer.lock/package-lock.json)
	LockFile string

	// LockExists indicates if the lock file exists
	LockExists bool

	// LockModified is the last modified time of lock file
	LockModified time.Time

	// VendorDir is the dependencies directory (vendor/node_modules)
	VendorDir string

	// VendorExists indicates if the dependencies directory exists
	VendorExists bool

	// VendorModified is the last modified time of vendor directory
	VendorModified time.Time

	// NeedsUpdate indicates if dependencies need updating
	NeedsUpdate bool
}

// PHPCSOptions represents options for PHPCS operations
type PHPCSOptions struct {
	// WorkingDir is the directory to run PHPCS in
	WorkingDir string

	// ConfigFile is the path to phpcs.xml or phpcs.xml.dist
	ConfigFile string

	// Standard is the coding standard to use (overrides config)
	Standard string

	// Extensions are the file extensions to check
	Extensions []string

	// Ignore is the pattern to ignore
	Ignore string

	// Files are specific files/directories to check
	Files []string

	// Report is the report format (full, summary, json, etc.)
	Report string

	// ShowSniffs shows sniff codes in report
	ShowSniffs bool

	// Severity is the minimum severity level
	Severity int

	// ErrorSeverity is the minimum error severity level
	ErrorSeverity int

	// WarningSeverity is the minimum warning severity level
	WarningSeverity int
}

// PHPCSResult represents the result of a PHPCS check
type PHPCSResult struct {
	// Success indicates if PHPCS passed
	Success bool

	// Errors is the number of errors found
	Errors int

	// Warnings is the number of warnings found
	Warnings int

	// Fixable is the number of fixable issues
	Fixable int

	// Files is the list of files with issues
	Files []PHPCSFileResult

	// Output is the raw PHPCS output
	Output string

	// ExitCode is the PHPCS exit code
	ExitCode int
}

// PHPCSFileResult represents PHPCS results for a single file
type PHPCSFileResult struct {
	// File is the path to the file
	File string

	// Errors is the number of errors in this file
	Errors int

	// Warnings is the number of warnings in this file
	Warnings int

	// Messages are the specific messages for this file
	Messages []PHPCSMessage
}

// PHPCSMessage represents a single PHPCS message
type PHPCSMessage struct {
	// Line is the line number
	Line int

	// Column is the column number
	Column int

	// Type is the message type (ERROR or WARNING)
	Type string

	// Message is the message text
	Message string

	// Source is the sniff source (e.g., PSR2.Classes.PropertyDeclaration)
	Source string

	// Severity is the severity level
	Severity int

	// Fixable indicates if this is auto-fixable
	Fixable bool
}

// ScriptInfo represents information about a build script
type ScriptInfo struct {
	// Name is the script name
	Name string

	// Path is the full path to the script
	Path string

	// Type is the script type (build, composer, npm, custom)
	Type string

	// Description is a description of what the script does
	Description string

	// Order is the execution order (for numbered scripts like 10-mu-plugins.sh)
	Order int
}

// WatchOptions represents options for file watching
type WatchOptions struct {
	// Paths are the paths to watch
	Paths []string

	// IgnorePatterns are patterns to ignore
	IgnorePatterns []string

	// Command is the command to run on changes
	Command string

	// Debounce is the debounce duration in milliseconds
	Debounce int

	// Recursive enables recursive watching
	Recursive bool

	// Verbose enables verbose output
	Verbose bool
}

// BuildResult represents the result of a build operation
type BuildResult struct {
	// Success indicates if the build succeeded
	Success bool

	// Duration is how long the build took
	Duration time.Duration

	// Steps are the individual build steps that were executed
	Steps []BuildStep

	// Output is the combined build output
	Output string

	// Error is any error that occurred
	Error error
}

// BuildStep represents a single step in the build process
type BuildStep struct {
	// Name is the step name
	Name string

	// Command is the command that was executed
	Command string

	// Success indicates if this step succeeded
	Success bool

	// Duration is how long this step took
	Duration time.Duration

	// Output is the output from this step
	Output string

	// Error is any error from this step
	Error error
}

// HuskyConfig represents Husky configuration
type HuskyConfig struct {
	// Enabled indicates if Husky is enabled
	Enabled bool

	// ConfigFile is the path to the husky config
	ConfigFile string

	// PreCommit is the pre-commit hook command
	PreCommit string

	// PrePush is the pre-push hook command
	PrePush string

	// CommitMsg is the commit-msg hook command
	CommitMsg string
}

// ComposerJSON represents a parsed composer.json file
type ComposerJSON struct {
	// Name is the package name
	Name string `json:"name"`

	// Type is the package type
	Type string `json:"type"`

	// Description is the package description
	Description string `json:"description"`

	// Require lists required packages
	Require map[string]string `json:"require"`

	// RequireDev lists dev packages
	RequireDev map[string]string `json:"require-dev"`

	// Scripts lists composer scripts
	Scripts map[string]interface{} `json:"scripts"`

	// Config contains composer configuration
	Config map[string]interface{} `json:"config"`

	// Autoload contains autoload configuration
	Autoload map[string]interface{} `json:"autoload"`
}

// PackageJSON represents a parsed package.json file
type PackageJSON struct {
	// Name is the package name
	Name string `json:"name"`

	// Version is the package version
	Version string `json:"version"`

	// Description is the package description
	Description string `json:"description"`

	// Scripts lists npm scripts
	Scripts map[string]string `json:"scripts"`

	// Dependencies lists runtime dependencies
	Dependencies map[string]string `json:"dependencies"`

	// DevDependencies lists dev dependencies
	DevDependencies map[string]string `json:"devDependencies"`

	// Engines specifies engine requirements
	Engines map[string]string `json:"engines"`
}

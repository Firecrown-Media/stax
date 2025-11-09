package system

// SystemInfo represents system information
type SystemInfo struct {
	OS           string
	Architecture string
	Hostname     string
	User         string
}

// Permission represents a file or directory permission
type Permission struct {
	Path  string
	Mode  string
	Owner string
	Group string
}

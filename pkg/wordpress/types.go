package wordpress

// Site represents a WordPress site
type Site struct {
	ID          int
	Name        string
	URL         string
	AdminEmail  string
	Description string
}

// Theme represents a WordPress theme
type Theme struct {
	Name      string
	Status    string
	Version   string
	Directory string
}

// Plugin represents a WordPress plugin
type Plugin struct {
	Name      string
	Status    string
	Version   string
	Directory string
}

// User represents a WordPress user
type User struct {
	ID    int
	Login string
	Email string
	Roles []string
}

// WPConfig represents WordPress configuration
type WPConfig struct {
	DBName        string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPrefix      string
	AuthKey       string
	SecureAuthKey string
	LoggedInKey   string
	NonceKey      string
}

package utils

// allowedSubcommands defines the allowed Tailscale subcommands list.
var allowedSubcommands = map[string]bool{
	"up":        true,
	"down":      true,
	"set":       true,
	"login":     true,
	"logout":    true,
	"switch":    true,
	"configure": true,
	"netcheck":  true,
	"ip":        true,
	"status":    true,
	"ping":      true,
	"nc":        true,
	"ssh":       true,
	"funnel":    true,
	"serve":     true,
	"version":   true,
	"web":       true,
	"file":      true,
	"bugreport": true,
	"cert":      true,
	"lock":      true,
	"licenses":  true,
	"exit-node": true,
	"update":    true,
}

const (
	// ESC represents the identifier for the escape key, used for UI control and shortcuts
	ESC = "ESC"

	// GET_KEY_URL defines an API endpoint for user login process
	// This URL points to an external service for authentication and retrieval of user credentials
	GET_KEY_URL = "https://sky-tailscale.sky1218.com/api/logIn"
)

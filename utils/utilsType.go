package utils

// AllowedSubcommands defines the allowed Tailscale subcommands list.
var AllowedSubcommands = map[string]bool{
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
	// KeyEsc represents the identifier for the escape key, used for UI control and shortcuts
	KeyEsc = "ESC"

	// LoginAPIEndpoint defines an API endpoint for user login process
	LoginAPIEndpoint = "https://sky-tailscale.sky1218.com/api/logIn"
)

// WaitAndExitConfig contains configuration options for waitAndExit function
type WaitAndExitConfig struct {
	ShouldExit bool // Controls whether to exit after waiting
	Countdown  int  // Timeout duration in seconds
}

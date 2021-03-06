package httpserver

const (
	// ModeDebug debug mode
	ModeDebug = "debug"
	// ModeTest test mode
	ModeTest = "test"
)

var modeRuning = ModeDebug

// SetMode set mode
func SetMode(mode string) {
	modeRuning = mode
}

package cotton

const (
	// ModeDebug debug mode
	ModeDebug = iota
	// ModeTest test mode
	ModeTest
	ModeProduct
)

var modeRuning = ModeDebug

// SetMode set mode
func SetMode(mode int) {
	modeRuning = mode
}

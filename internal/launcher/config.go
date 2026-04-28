package launcher

const (
	usageFile         = "/tmp/launcher_usage.json"
	envCloseWarpFloat = "LAUNCHER_CLOSE_WARP_FLOAT"

	// appsVisibleMax is how many app rows to show at once in the Apps section (Raycast-like).
	appsVisibleMax = 10

	// WarpWindowWidth and WarpWindowHeight are used by launch-warp-launcher.sh (see warp_window.env).
	// Keep in sync: numbers must match warp_window.env and internal/launcher/config.go.
	WarpWindowWidth  = 640
	WarpWindowHeight = 720
)

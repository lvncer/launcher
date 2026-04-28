package launcher

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type yabaiWindow struct {
	ID int `json:"id"`
}

func captureCurrentWindowIDIfRequested() int {
	if os.Getenv(envCloseWarpFloat) != "1" {
		return 0
	}

	out, err := exec.Command("yabai", "-m", "query", "--windows", "--window").Output()
	if err != nil {
		return 0
	}

	var w yabaiWindow
	if err := json.Unmarshal(out, &w); err != nil {
		return 0
	}
	return w.ID
}

func closeWarpFloatIfRequested(windowID int) {
	if os.Getenv(envCloseWarpFloat) != "1" {
		return
	}
	if windowID == 0 {
		return
	}

	windowIDStr := fmt.Sprintf("%d", windowID)

	if err := exec.Command("yabai", "-m", "window", windowIDStr, "--close").Run(); err == nil {
		return
	}

	_ = exec.Command("yabai", "-m", "window", windowIDStr, "--focus").Run()
	_ = exec.Command("osascript",
		"-e", `tell application "System Events" to keystroke "w" using command down`,
	).Run()
}

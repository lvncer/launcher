package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func loadItems() []item {
	items := loadApps()
	items = append(items, defaultCommands()...)
	return items
}

func loadApps() []item {
	var items []item
	files, _ := os.ReadDir("/Applications")

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".app") {
			name := strings.TrimSuffix(f.Name(), ".app")
			cmd := fmt.Sprintf("open -a '%s'", name)
			appPath := filepath.Join("/Applications", f.Name())
			items = append(items, item{
				title:    name,
				cmd:      cmd,
				typ:      appItem,
				iconPath: appIconPath(appPath),
			})
		}
	}
	return items
}

func defaultCommands() []item {
	return []item{
		{title: "Git Status", cmd: "git status", typ: commandItem},
		{title: "Docker PS", cmd: "docker ps", typ: commandItem},
		{title: "Open GitHub", cmd: "open https://github.com", typ: commandItem},
	}
}

func appIconPath(appPath string) string {
	infoPlist := filepath.Join(appPath, "Contents", "Info.plist")
	out, err := exec.Command("/usr/libexec/PlistBuddy", "-c", "Print :CFBundleIconFile", infoPlist).Output()
	if err != nil {
		return fallbackIconPath(appPath)
	}

	iconName := strings.TrimSpace(string(out))
	if iconName == "" {
		return fallbackIconPath(appPath)
	}
	if !strings.HasSuffix(strings.ToLower(iconName), ".icns") {
		iconName += ".icns"
	}

	iconPath := filepath.Join(appPath, "Contents", "Resources", iconName)
	if _, err := os.Stat(iconPath); err == nil {
		return iconPath
	}
	return fallbackIconPath(appPath)
}

func fallbackIconPath(appPath string) string {
	matches, err := filepath.Glob(filepath.Join(appPath, "Contents", "Resources", "*.icns"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	return matches[0]
}

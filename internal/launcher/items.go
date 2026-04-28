package launcher

import (
	"fmt"
	"os"
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
			items = append(items, item{name, cmd, appItem})
		}
	}
	return items
}

func defaultCommands() []item {
	return []item{
		{"Git Status", "git status", commandItem},
		{"Docker PS", "docker ps", commandItem},
		{"Open GitHub", "open https://github.com", commandItem},
	}
}

package launcher

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	left := m.input.View() + "\n\n"

	for i, item := range m.matches {
		cursor := "  "
		if i == m.index {
			cursor = "> "
		}
		left += fmt.Sprintf("%s%s\n", cursor, item.title)
	}

	right := "\n\n--- Preview ---\n"
	if len(m.matches) > 0 {
		item := m.matches[m.index]
		right += fmt.Sprintf("Type: %v\nCmd: %s\n", item.typ, item.cmd)
	}

	return splitView(left, right)
}

func splitView(left, right string) string {
	linesL := strings.Split(left, "\n")
	linesR := strings.Split(right, "\n")

	max := len(linesL)
	if len(linesR) > max {
		max = len(linesR)
	}

	var out string
	for i := 0; i < max; i++ {
		l, r := "", ""
		if i < len(linesL) {
			l = linesL[i]
		}
		if i < len(linesR) {
			r = linesR[i]
		}
		out += fmt.Sprintf("%-40s %s\n", l, r)
	}
	return out
}

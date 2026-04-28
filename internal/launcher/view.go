package launcher

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	left := m.buildLeftLines()
	right := m.buildRightLines(len(left))
	return splitViewLines(left, right)
}

func (m model) buildLeftLines() []string {
	var lines []string
	lines = append(lines, strings.Split(m.input.View(), "\n")...)
	lines = append(lines, "")

	recentN := recentSectionSize(m.matches, m.usage)
	if recentN > 0 {
		lines = append(lines, "Recent(5 items)")
		for i := 0; i < recentN; i++ {
			cur := "  "
			if m.index == i {
				cur = "> "
			}
			lines = append(lines, cur+m.matches[i].title)
		}
		lines = append(lines, "")
	}

	appTotal := len(m.matches) - recentN
	if appTotal == 0 {
		return lines
	}

	scrollable := appTotal > appsVisibleMax
	visibleEnd := min(m.appsScroll+appsVisibleMax, appTotal)
	if scrollable {
		lines = append(lines, fmt.Sprintf("Apps  %d-%d of %d  ·  max %d rows  ·  ↑↓", m.appsScroll+1, visibleEnd, appTotal, appsVisibleMax))
	} else {
		lines = append(lines, "Apps")
	}

	if scrollable && m.appsScroll > 0 {
		lines = append(lines, fmt.Sprintf("  ▲  %d more above", m.appsScroll))
	}

	for j := m.appsScroll; j < visibleEnd; j++ {
		gi := recentN + j
		if gi >= len(m.matches) {
			break
		}
		cur := "  "
		if m.index == gi {
			cur = "> "
		}
		line := cur + m.matches[gi].title
		if scrollable {
			line = fmt.Sprintf("%-36s %s", line, appsScrollBarCell(j-m.appsScroll, visibleEnd-m.appsScroll, m.appsScroll, appTotal))
		}
		lines = append(lines, line)
	}

	if scrollable {
		below := appTotal - visibleEnd
		if below > 0 {
			lines = append(lines, fmt.Sprintf("  ▼  %d more below  (use ↑↓)", below))
		} else {
			lines = append(lines, "  ▼  end of list")
		}
	}
	return lines
}

// appsScrollBarCell draws one cell of a vertical scrollbar for the visible Apps rows.
func appsScrollBarCell(row, visibleRows, fromTop, total int) string {
	if total <= visibleRows || visibleRows <= 0 {
		return ""
	}

	thumb := max(1, (visibleRows*visibleRows+total-1)/total)
	maxPos := visibleRows - thumb
	pos := 0
	if total > visibleRows {
		pos = fromTop * maxPos / (total - visibleRows)
		if pos > maxPos {
			pos = maxPos
		}
	}

	if row >= pos && row < pos+thumb {
		return "│█│"
	}
	return "│░│"
}

func (m model) buildRightLines(leftLineCount int) []string {
	right := make([]string, leftLineCount)
	if len(m.matches) == 0 {
		return right
	}
	item := m.matches[m.index]
	preview := []string{
		"--- Preview ---",
		fmt.Sprintf("Type: %v", item.typ),
		fmt.Sprintf("Cmd: %s", item.cmd),
		"",
		"Move: ↑ / ↓  (list scrolls with selection)",
	}
	for i, line := range preview {
		target := i + previewTopPadding
		if target < len(right) {
			right[target] = line
		}
	}
	return right
}

func recentSectionSize(items []item, usageByTitle map[string]usage) int {
	count := 0
	for _, item := range items {
		if count >= recentLimit || usageByTitle[item.title].Last == 0 {
			break
		}
		count++
	}
	return count
}

func splitViewLines(left, right []string) string {
	maxH := len(left)
	if len(right) > maxH {
		maxH = len(right)
	}
	for len(left) < maxH {
		left = append(left, "")
	}
	for len(right) < maxH {
		right = append(right, "")
	}
	var b strings.Builder
	const leftWidth = 48
	for i := 0; i < maxH; i++ {
		_, _ = b.WriteString(fmt.Sprintf("%-*s %s\n", leftWidth, left[i], right[i]))
	}
	return b.String()
}

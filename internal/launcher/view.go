package launcher

import (
	"fmt"
	"strings"
	"unicode/utf8"
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
	lines = append(lines, "Move: ↑ / ↓  (list scrolls with selection)")
	lines = append(lines, "")

	recentN := recentSectionSize(m.matches, m.usage)
	if recentN > 0 {
		lines = append(lines, "Recent(5 items)")
		for i := 0; i < recentN; i++ {
			cur := "  "
			if m.index == i {
				cur = "> "
			}
			lines = append(lines, itemDisplayLines(m.matches[i], cur)...)
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

	if scrollable {
		if m.appsScroll > 0 {
			lines = append(lines, fmt.Sprintf("  ▲  %d more above", m.appsScroll))
		} else {
			lines = append(lines, "  ▲  start of list")
		}
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
		itemLines := itemDisplayLines(m.matches[gi], cur)
		lines = append(lines, itemLines...)
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
		_, _ = b.WriteString(padANSI(left[i], leftWidth) + " " + right[i] + "\n")
	}
	return b.String()
}

func padANSI(s string, width int) string {
	visible := ansiVisibleWidth(s)
	if visible >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visible)
}

func ansiVisibleWidth(s string) int {
	width := 0
	for i := 0; i < len(s); {
		if s[i] == '\x1b' {
			i = skipANSISequence(s, i)
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 0 {
			break
		}
		width += runeDisplayWidth(r)
		i += size
	}
	return width
}

func skipANSISequence(s string, i int) int {
	i++
	if i < len(s) && s[i] == '[' {
		i++
		for i < len(s) {
			if s[i] >= '@' && s[i] <= '~' {
				return i + 1
			}
			i++
		}
	}
	return i
}

func runeDisplayWidth(r rune) int {
	if r == '\t' {
		return 4
	}
	if r < 32 {
		return 0
	}
	return 1
}

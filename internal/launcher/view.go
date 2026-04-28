package launcher

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	left := m.buildLeftLines()
	right := m.buildRightLines(left)
	return splitViewLines(left, right)
}

func (m model) buildLeftLines() []string {
	var lines []string
	lines = append(lines, strings.Split(m.input.View(), "\n")...)
	lines = append(lines, "")

	recentN := recentSectionSize(m.matches, m.usage)
	if recentN > 0 {
		lines = append(lines, padVBar("Recent", ' '))
		for i := 0; i < recentN; i++ {
			cur := "  "
			if m.index == i {
				cur = "> "
			}
			lines = append(lines, padVBar(cur+m.matches[i].title, ' '))
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
		lines = append(lines, padVBar(
			fmt.Sprintf("Apps  %d–%d of %d  ·  max %d  ·  ↕", m.appsScroll+1, visibleEnd, appTotal, appsVisibleMax),
			' ',
		))
	} else {
		lines = append(lines, padVBar("Apps", ' '))
	}

	if scrollable && m.appsScroll > 0 {
		lines = append(lines, padVBar(fmt.Sprintf("  ▲  %d more above", m.appsScroll), ' '))
	}

	nVis := visibleEnd - m.appsScroll
	thumbH, thumbStart := vThumb(nVis, appTotal, m.appsScroll, scrollable)
	for j := m.appsScroll; j < visibleEnd; j++ {
		gi := recentN + j
		if gi >= len(m.matches) {
			break
		}
		cur := "  "
		if m.index == gi {
			cur = "> "
		}
		rk := j - m.appsScroll
		v := vBarRune(rk, thumbStart, thumbH)
		lines = append(lines, cur+m.matches[gi].title+formatVCol(v))
	}

	if scrollable {
		below := appTotal - visibleEnd
		if below > 0 {
			lines = append(lines, padVBar(fmt.Sprintf("  ▼  %d more below  (use ↑↓)", below), ' '))
		} else {
			lines = append(lines, padVBar("  ▼  end of list", ' '))
		}
	}
	return lines
}

func padVBar(s string, fill rune) string {
	return s + formatVCol(fill)
}

func formatVCol(c rune) string {
	// one visual column: space + glyph (align with █/░)
	if c == ' ' {
		return "   "
	}
	return fmt.Sprintf("  %c", c)
}

func vBarRune(rk, thumbStart, thumbH int) rune {
	if thumbH <= 0 {
		return '░'
	}
	if rk >= thumbStart && rk < thumbStart+thumbH {
		return '█'
	}
	return '░'
}

func vThumb(nVis, appTotal, appsScroll int, scrollable bool) (thumbH, thumbStart int) {
	if nVis <= 0 {
		return 1, 0
	}
	if !scrollable || appTotal <= nVis {
		return nVis, 0
	}
	thumbH = max(1, (nVis*nVis+appTotal-1)/appTotal)
	if thumbH > nVis {
		thumbH = nVis
	}
	maxScroll := appTotal - nVis
	if maxScroll <= 0 {
		return thumbH, 0
	}
	thumbStart = (appsScroll*(nVis-thumbH) + maxScroll - 1) / maxScroll
	if thumbStart < 0 {
		thumbStart = 0
	}
	if thumbStart > nVis-thumbH {
		thumbStart = nVis - thumbH
	}
	return thumbH, thumbStart
}

func (m model) buildRightLines(left []string) []string {
	right := make([]string, len(left))
	inputLines := strings.Split(m.input.View(), "\n")
	previewStart := len(inputLines) + 1
	if len(left) == 0 {
		return right
	}
	if previewStart > len(left)-1 {
		previewStart = 0
	}
	if len(m.matches) == 0 {
		return right
	}
	item := m.matches[m.index]
	preview := []string{
		"--- Preview ---",
		fmt.Sprintf("Type: %v", item.typ),
		fmt.Sprintf("Cmd: %s", item.cmd),
		"",
		"Move: ↑ / ↓",
	}
	for i, line := range preview {
		idx := previewStart + i
		if idx < len(right) {
			right[idx] = line
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
	const leftWidth = 50
	for i := 0; i < maxH; i++ {
		_, _ = b.WriteString(fmt.Sprintf("%-*s %s\n", leftWidth, left[i], right[i]))
	}
	return b.String()
}

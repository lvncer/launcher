package launcher

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type model struct {
	input         textinput.Model
	items         []item
	matches       []item
	index         int
	usage         map[string]usage
	mode          mode
	lastFilterKey string
	appsScroll    int
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()

	m := model{
		input: ti,
		items: loadItems(),
		usage: loadUsage(),
		mode:  appMode,
	}
	m.filter()
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape, tea.KeyCtrlC:
			return m, tea.Quit
		}

		switch msg.String() {
		case "enter":
			if len(m.matches) > 0 {
				item := m.matches[m.index]
				_ = exec.Command("bash", "-c", item.cmd).Start()
				m.updateUsage(item.title)
			}
			return m, tea.Quit

		case "up", "ctrl+p":
			if m.index > 0 {
				m.index--
				m.ensureAppsScrollVisible()
			}
			return m, nil

		case "down", "ctrl+n":
			if m.index < len(m.matches)-1 {
				m.index++
				m.ensureAppsScrollVisible()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	query := m.input.Value()
	m.detectMode(query)
	m.filter()

	return m, cmd
}

func (m *model) detectMode(q string) {
	if strings.HasPrefix(q, ">") {
		m.mode = cmdMode
	} else if strings.HasPrefix(q, "/") {
		m.mode = fileMode
	} else {
		m.mode = appMode
	}
}

func (m *model) filter() {
	query := m.input.Value()
	query = strings.TrimPrefix(query, ">")
	query = strings.TrimPrefix(query, "/")

	var filtered []item
	for _, item := range m.items {
		if !m.matchesMode(item) {
			continue
		}

		if fuzzy.Match(query, item.title) {
			filtered = append(filtered, item)
		}
	}

	m.matches = sortByRecentThenName(filtered, m.usage)
	m.updateIndex()
	m.clampAppsScroll()
	m.ensureAppsScrollVisible()
}

func (m model) matchesMode(item item) bool {
	if m.mode == appMode && item.typ != appItem {
		return false
	}
	if m.mode == cmdMode && item.typ != commandItem {
		return false
	}
	return true
}

func (m *model) updateIndex() {
	full := m.input.Value()
	if full != m.lastFilterKey {
		m.index = 0
		m.appsScroll = 0
		m.lastFilterKey = full
	}
	switch {
	case len(m.matches) == 0:
		m.index = 0
	case m.index >= len(m.matches):
		m.index = len(m.matches) - 1
	}
}

func (m *model) clampAppsScroll() {
	recent := recentSectionSize(m.matches, m.usage)
	appCount := len(m.matches) - recent
	if appCount <= appsVisibleMax {
		m.appsScroll = 0
		return
	}
	maxScroll := appCount - appsVisibleMax
	if m.appsScroll > maxScroll {
		m.appsScroll = maxScroll
	}
	if m.appsScroll < 0 {
		m.appsScroll = 0
	}
}

func (m *model) ensureAppsScrollVisible() {
	recent := recentSectionSize(m.matches, m.usage)
	n := len(m.matches)
	if m.index < recent {
		return
	}
	appCount := n - recent
	if appCount <= appsVisibleMax {
		m.appsScroll = 0
		return
	}
	rel := m.index - recent
	if rel < m.appsScroll {
		m.appsScroll = rel
	}
	if rel >= m.appsScroll+appsVisibleMax {
		m.appsScroll = rel - appsVisibleMax + 1
	}
	m.clampAppsScroll()
}

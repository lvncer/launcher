package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type ItemType int

const (
	App ItemType = iota
	Command
	File
)

type Item struct {
	Title string
	Cmd   string
	Type  ItemType
}

type Usage struct {
	Count int
	Last  int64
}

type Model struct {
	input   textinput.Model
	items   []Item
	matches []Item
	index   int
	usage   map[string]Usage
	mode    string
}

const usageFile = "/tmp/launcher_usage.json"

// -------------------- init --------------------

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()

	items := loadApps()
	items = append(items, defaultCommands()...)

	return Model{
		input:   ti,
		items:   items,
		matches: items,
		usage:   loadUsage(),
		mode:    "app",
	}
}

func loadApps() []Item {
	var items []Item
	files, _ := os.ReadDir("/Applications")

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".app") {
			name := strings.TrimSuffix(f.Name(), ".app")
			cmd := fmt.Sprintf("open -a '%s'", name)
			items = append(items, Item{name, cmd, App})
		}
	}
	return items
}

func defaultCommands() []Item {
	return []Item{
		{"Git Status", "git status", Command},
		{"Docker PS", "docker ps", Command},
		{"Open GitHub", "open https://github.com", Command},
	}
}

// -------------------- usage --------------------

func loadUsage() map[string]Usage {
	m := map[string]Usage{}
	data, err := os.ReadFile(usageFile)
	if err != nil {
		return m
	}
	json.Unmarshal(data, &m)
	return m
}

func saveUsage(m map[string]Usage) {
	data, _ := json.Marshal(m)
	os.WriteFile(usageFile, data, 0644)
}

func (m *Model) updateUsage(key string) {
	u := m.usage[key]
	u.Count++
	u.Last = time.Now().Unix()
	m.usage[key] = u
	saveUsage(m.usage)
}

// -------------------- update --------------------

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if len(m.matches) > 0 {
				item := m.matches[m.index]
				exec.Command("bash", "-c", item.Cmd).Start()
				m.updateUsage(item.Title)
			}
			return m, tea.Quit

		case "up":
			if m.index > 0 {
				m.index--
			}
		case "down":
			if m.index < len(m.matches)-1 {
				m.index++
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	query := m.input.Value()
	m.detectMode(query)
	m.filter()

	return m, cmd
}

// -------------------- logic --------------------

func (m *Model) detectMode(q string) {
	if strings.HasPrefix(q, ">") {
		m.mode = "cmd"
	} else if strings.HasPrefix(q, "/") {
		m.mode = "file"
	} else {
		m.mode = "app"
	}
}

func (m *Model) filter() {
	query := m.input.Value()
	query = strings.TrimPrefix(query, ">")
	query = strings.TrimPrefix(query, "/")

	var filtered []Item

	for _, item := range m.items {
		if m.mode == "app" && item.Type != App {
			continue
		}
		if m.mode == "cmd" && item.Type != Command {
			continue
		}

		if fuzzy.Match(query, item.Title) {
			filtered = append(filtered, item)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		ui := m.usage[filtered[i].Title]
		uj := m.usage[filtered[j].Title]

		if ui.Count != uj.Count {
			return ui.Count > uj.Count
		}
		return ui.Last > uj.Last
	})

	m.matches = filtered
	m.index = 0
}

// -------------------- view --------------------

func (m Model) View() string {
	left := m.input.View() + "\n\n"

	for i, item := range m.matches {
		cursor := "  "
		if i == m.index {
			cursor = "> "
		}
		left += fmt.Sprintf("%s%s\n", cursor, item.Title)
	}

	right := "\n\n--- Preview ---\n"
	if len(m.matches) > 0 {
		item := m.matches[m.index]
		right += fmt.Sprintf("Type: %v\nCmd: %s\n", item.Type, item.Cmd)
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

// -------------------- main --------------------

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
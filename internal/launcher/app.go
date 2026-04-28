package launcher

import tea "github.com/charmbracelet/bubbletea"

func Run() error {
	windowID := captureCurrentWindowIDIfRequested()
	p := tea.NewProgram(newModel())
	if err := p.Start(); err != nil {
		return err
	}
	closeWarpFloatIfRequested(windowID)
	return nil
}

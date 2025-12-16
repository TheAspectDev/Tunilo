package tui

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width - 2)

		tableHeight := msg.Height - 10
		if tableHeight < 5 {
			tableHeight = 5
		}
		m.table.SetHeight(tableHeight)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "backspace":
			selected := m.table.SelectedRow()

			if len(selected) == 0 {
				return m, nil
			}
			id := selected[0]

			m.srv.SessionsMu.RLock()
			session, exists := m.srv.Sessions[id]
			m.srv.SessionsMu.RUnlock()

			if exists {
				session.Close()
				m.srv.SessionsMu.Lock()
				delete(m.srv.Sessions, id)
				m.srv.SessionsMu.Unlock()
				m.updateClientTable()
			}
			return m, nil
		}

	case refreshMsg:
		m.updateClientTable()
		return m, refreshTick()
	}

	m.table, _ = m.table.Update(msg)

	var cmd tea.Cmd

	return m, cmd
}

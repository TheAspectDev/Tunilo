package tui

import (
	"sort"

	"github.com/TheAspectDev/tunio/internal/tui"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateClientTable() {
	m.srv.SessionsMu.RLock()
	keys := make([]string, 0, len(m.srv.Sessions))

	for k := range m.srv.Sessions {
		keys = append(keys, k)
	}
	m.srv.SessionsMu.RUnlock()

	sort.Strings(keys)
	rows := make([]table.Row, 0, len(keys))

	if len(keys) == 0 {
		m.table.Blur()
		m.table.SetStyles(unselectableStyle)
		rows = append(rows, table.Row{
			tui.Color("No clients connected", lipgloss.Color("#5e5e5eff")),
			"",
		})
		m.table.SetRows(rows)
		return
	}

	// default color-override doesn't work here as it plays with components width
	// so colors are going to be removed if the row is selected
	for _, k := range keys {
		if len(m.table.SelectedRow()) > 0 && m.table.SelectedRow()[0] == k {
			rows = append(rows, table.Row{
				k,
				"● Connected",
			})
		} else {
			rows = append(rows, table.Row{
				k,
				tui.Color("● Connected", lipgloss.Color("#00c8ffff")),
			})
		}
	}
	m.table.Focus()
	m.table.SetStyles(defaultStyle)
	cursor := m.table.Cursor()
	m.table.SetRows(rows)
	m.table.SetCursor(cursor)
}

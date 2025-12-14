package components

import (
	"fmt"
	"sort"

	"github.com/TheAspectDev/tunio/internal/server"
	"github.com/TheAspectDev/tunio/internal/tui"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defaultStyle = table.Styles{
	Header: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c1c1c1ff")).
		Padding(0, 1).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottomForeground(lipgloss.Color("#3c3c3cff")),
	Selected: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#00bddaff")),
	Cell: lipgloss.NewStyle().MarginLeft(1),
}

var unselectableStyle = table.Styles{
	Header: defaultStyle.Header,
	// Set Selected style to be entirely transparent/neutral
	Selected: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("")),
	Cell:     defaultStyle.Cell,
}

type model struct {
	spinner  spinner.Model
	quitting bool
	err      error
	srv      *server.Server
	table    table.Model
}

func newClientTable() table.Model {
	columns := []table.Column{
		{Title: "Hostname", Width: 45},
		{Title: "Status", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	t.SetStyles(defaultStyle)

	return t
}

func SpinnerModel(srv *server.Server) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00c8ffff"))
	t := newClientTable()
	return model{spinner: s, srv: srv, table: t}
}

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
		m.table.SetCursor(-1)
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
	// m.table.SetHeight((len(rows) + 2) * 2)
	cursor := m.table.Cursor()
	m.table.SetRows(rows)
	m.table.SetCursor(cursor)
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

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
			session := m.srv.Sessions[m.table.SelectedRow()[0]]
			session.Close()
			return m, nil
		}
	}

	m.table, _ = m.table.Update(msg)

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	m.updateClientTable()

	return m, cmd
}

func (m model) HelpNotes() string {
	var keyNote = lipgloss.NewStyle().Foreground(lipgloss.Color("#338f8cff"))
	var valNote = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000e3"))

	return (keyNote.Render("↑/k") + " 	" + valNote.Render("move up")) + "			" + (keyNote.Render("q/esc") + " " + valNote.Render("kill tunnel")) + "\n" +
		(keyNote.Render("↓/j") + "     " + valNote.Render("move down")) + "\n" +
		(keyNote.Render("delete") + "  " + valNote.Render("close connection")) + "\n"

}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var title = lipgloss.NewStyle().
		SetString("  Tunilo Tunnel").
		Width(1000).
		Height(1).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true).
		Background(lipgloss.Color("#008599ff"))

	tableBox := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("#3c3c3cff")).
		Border(lipgloss.RoundedBorder()).
		Render(m.table.View())

	return fmt.Sprintf(
		"%s"+"\n\n"+"%s"+"\n\n"+"%s",
		title.Render(),
		tableBox,
		m.HelpNotes(),
	)
}

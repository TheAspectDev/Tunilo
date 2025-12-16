package tui

import (
	"time"

	"github.com/TheAspectDev/tunio/internal/server"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type refreshMsg struct{}

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

func ServerModel(srv *server.Server) model {
	t := newClientTable()
	return model{srv: srv, table: t}
}

func refreshTick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg {
		return refreshMsg{}
	})
}

func (m model) Init() tea.Cmd { return refreshTick() }

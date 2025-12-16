package tui

import (
	"github.com/TheAspectDev/tunio/internal/client"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	quitting bool
	err      error
	session  *client.Session

	logs     []LogMsg
	viewport viewport.Model
}

func ClientModel(session *client.Session) model {
	vp := viewport.New(100, 8)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#555555"))

	return model{
		session:  session,
		logs:     []LogMsg{},
		viewport: vp,
	}
}

func (m model) Init() tea.Cmd {
	return waitForLog()
}

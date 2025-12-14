package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	spinner        spinner.Model
	publicAddress  string
	controlAddress string
	quitting       bool
	err            error
}

func SpinnerModel(publicAddress string, controlAddress string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00c8ffff"))
	return model{spinner: s, publicAddress: publicAddress, controlAddress: controlAddress}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
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

	var serverRunning = lipgloss.NewStyle().
		SetString("Public server running at " + m.publicAddress).
		Width(54).
		Height(1).
		Foreground(lipgloss.Color("#89efffff"))

	var controlRunning = lipgloss.NewStyle().
		SetString("Control server running at " + m.controlAddress).
		Width(54).
		Height(1).
		Foreground(lipgloss.Color("#268492ff"))

	var pressQtoQuit = lipgloss.NewStyle().
		SetString("   ( press q to quit )" + "\n").
		Width(54).
		Height(1).
		PaddingLeft(9).
		Foreground(lipgloss.Color("#3a3a3aff"))

	str := fmt.Sprintf("%s\n\n%s%s\n%s\n\n%s", title.Render(), m.spinner.View(), serverRunning.Render(), m.spinner.View()+controlRunning.Render(), pressQtoQuit.Render())
	if m.quitting {
		return str + "\n"
	}
	return str
}

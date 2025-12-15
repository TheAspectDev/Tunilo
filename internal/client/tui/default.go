package tui

import (
	"fmt"
	"strings"

	"github.com/TheAspectDev/tunio/internal/client"
	"github.com/TheAspectDev/tunio/internal/tui"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var LogChan = make(chan LogMsg, 256)

func Logf(format string, args ...any) {
	LogChan <- LogMsg{
		Text: fmt.Sprintf(format, args...),
		Type: LogInfo,
	}
}

func Errorf(err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if err != nil {
		msg += ": " + err.Error()
	}

	LogChan <- LogMsg{
		Text: msg,
		Err:  err,
		Type: LogError,
	}
}

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

func waitForLog() tea.Cmd {
	return func() tea.Msg {
		return <-LogChan
	}
}

func (m model) Init() tea.Cmd {
	return waitForLog()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case LogMsg:
		m.appendLog(msg)
		return m, waitForLog()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *model) appendLog(log LogMsg) {
	m.logs = append(m.logs, log)

	var content strings.Builder

	for _, l := range m.logs {
		switch l.Type {
		case LogInfo:
			content.WriteString("• " + l.Text + "\n")

		case LogError:
			content.WriteString(
				tui.Color("• ", lipgloss.Color("#ff0000")) +
					tui.Color(l.Text, lipgloss.Color("#5e3c3cff")) +
					"\n",
			)
		}
	}

	m.viewport.SetContent(content.String())
	m.viewport.GotoBottom()
}

func (m model) HelpNotes() string {
	var keyNote = lipgloss.NewStyle().Foreground(lipgloss.Color("#338f8cff"))
	var valNote = lipgloss.NewStyle().Foreground(lipgloss.Color("#5f5f5fff"))

	return (keyNote.Render("↑/k") + " 	" + valNote.Render("scroll up")) + "			" + (keyNote.Render("q/esc") + " " + valNote.Render("kill connection")) + "\n" +
		(keyNote.Render("↓/j") + "     " + valNote.Render("scroll down")) + "\n"
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	title := lipgloss.NewStyle().
		SetString("  Tunilo Tunnel - CLIENT").
		Width(m.viewport.Width).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true).
		Background(lipgloss.Color("#008599ff"))

	help := m.HelpNotes()

	logTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Render(" Logs ")

	return fmt.Sprintf(
		"%s\n\n%s\n%s\n%s",
		title.Render(),
		logTitle,
		m.viewport.View(),
		help,
	)
}

package tui

import (
	"fmt"
	"strings"

	"github.com/TheAspectDev/tunio/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogType int8

const (
	LogInfo LogType = iota
	LogError
)

type LogMsg struct {
	Text string
	Err  error
	Type LogType
}

type UILogger struct{}

func (UILogger) Logf(format string, args ...any) {
	Logf(format, args...)
}

func (UILogger) Errorf(err error, format string, args ...any) {
	Errorf(err, format, args...)
}

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

func waitForLog() tea.Cmd {
	return func() tea.Msg {
		return <-LogChan
	}
}

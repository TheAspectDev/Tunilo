package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) HelpNotes() string {
	keyNote := lipgloss.NewStyle().Foreground(lipgloss.Color("#338f8cff"))
	valNote := lipgloss.NewStyle().Foreground(lipgloss.Color("#5f5f5fff"))

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

	return fmt.Sprintf(
		"%s\n%s\n%s",
		title.Render(),
		m.viewport.View(),
		help,
	)
}

package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) HelpNotes() string {
	var keyNote = lipgloss.NewStyle().Foreground(lipgloss.Color("#338f8cff"))
	var valNote = lipgloss.NewStyle().Foreground(lipgloss.Color("#5f5f5fff"))

	return (keyNote.Render("↑/k") + " 	" + valNote.Render("move up")) + "			" + (keyNote.Render("q/esc") + " " + valNote.Render("kill tunnel")) + "\n" +
		(keyNote.Render("↓/j") + "     " + valNote.Render("move down")) + "\n" +
		(keyNote.Render("delete") + "  " + valNote.Render("close connection")) + "\n"

}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var title = lipgloss.NewStyle().
		SetString("  Tunilo Tunnel - SERVER").
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

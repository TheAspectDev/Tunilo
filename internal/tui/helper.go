package tui

import "github.com/charmbracelet/lipgloss"

func Color(text string, hex lipgloss.Color) string {
	return lipgloss.NewStyle().
		Foreground(hex).
		Render(text)
}

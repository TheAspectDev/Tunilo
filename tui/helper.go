package tui

import "github.com/charmbracelet/lipgloss"

// Convert text to colored text
func Color(text string, hex lipgloss.Color) string {
	return lipgloss.NewStyle().
		Foreground(hex).
		Render(text)
}

package components

import tea "github.com/charmbracelet/bubbletea"

type Component interface {
	Render(width int) string
	OnKeyPress(key string) (redirect int, cmd tea.Cmd)
	Focus() tea.Cmd
	Blur()
	Focused() bool
}

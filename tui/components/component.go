package components

import tea "github.com/charmbracelet/bubbletea"

const (
	RedirectNone   = 0
	RedirectSelect = 1
)

type Component interface {
	Render(width int) string
	OnKeyPress(key string) (redirect int, cmd tea.Cmd)
	Update(msg tea.Msg) tea.Cmd
	Focus() tea.Cmd
	Blur()
	Focused() bool
}

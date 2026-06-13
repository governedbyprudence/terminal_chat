package components

import (
	"terminal_chat/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type Footer struct {
	Text string
}

func NewFooter(text string) *Footer {
	return &Footer{Text: text}
}

func (f *Footer) Render(_ int) string {
	return styles.FooterText.Render(f.Text)
}

func (f *Footer) OnKeyPress(key string) (int, tea.Cmd) {
	return 0, nil
}

func (f *Footer) Focus() tea.Cmd { return nil }
func (f *Footer) Blur()          {}
func (f *Footer) Focused() bool  { return false }

package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var footerTextStyle = lipgloss.NewStyle().
	Foreground(DimText).
	Padding(0, 2)

type Footer struct {
	Text string
}

func NewFooter(text string) *Footer {
	return &Footer{Text: text}
}

func (f *Footer) Render(_ int) string {
	return footerTextStyle.Render(f.Text)
}

func (f *Footer) OnKeyPress(key string) (int, tea.Cmd) {
	return 0, nil
}

func (f *Footer) Update(msg tea.Msg) tea.Cmd { return nil }
func (f *Footer) Focus() tea.Cmd             { return nil }
func (f *Footer) Blur()                      {}
func (f *Footer) Focused() bool              { return false }

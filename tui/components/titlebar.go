package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var titleBarStyle = lipgloss.NewStyle().
	Background(Accent).
	Foreground(DeepBg).
	Bold(true).
	Padding(0, 3)

type TitleBar struct {
	Title string
}

func NewTitleBar(title string) *TitleBar {
	return &TitleBar{Title: title}
}

func (t *TitleBar) Render(width int) string {
	return titleBarStyle.Width(width).Render("  " + t.Title + "  ")
}

func (t *TitleBar) OnKeyPress(key string) (int, tea.Cmd) {
	return 0, nil
}

func (t *TitleBar) Update(msg tea.Msg) tea.Cmd { return nil }
func (t *TitleBar) Focus() tea.Cmd             { return nil }
func (t *TitleBar) Blur()                      {}
func (t *TitleBar) Focused() bool              { return false }

package components

import (
	"terminal_chat/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type TitleBar struct {
	Title string
}

func NewTitleBar(title string) *TitleBar {
	return &TitleBar{Title: title}
}

func (t *TitleBar) Render(_ int) string {
	return styles.TitleBar.Render("  " + t.Title + "  ")
}

func (t *TitleBar) OnKeyPress(key string) (int, tea.Cmd) {
	return 0, nil
}

func (t *TitleBar) Focus() tea.Cmd  { return nil }
func (t *TitleBar) Blur()           {}
func (t *TitleBar) Focused() bool   { return false }

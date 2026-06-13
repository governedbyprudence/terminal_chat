package components

import (
	"terminal_chat/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BorderFrame struct {
	Inner Component
}

func NewBorderFrame(inner Component) *BorderFrame {
	return &BorderFrame{Inner: inner}
}

func (f *BorderFrame) Render(width int) string {
	contentW := width - styles.OuterFrameHSize
	body := f.Inner.Render(contentW)
	wrapper := lipgloss.NewStyle().Width(contentW)
	return styles.OuterFrame.Render(wrapper.Render(body)) + "\n"
}

func (f *BorderFrame) OnKeyPress(key string) (int, tea.Cmd) {
	return f.Inner.OnKeyPress(key)
}

func (f *BorderFrame) Focus() tea.Cmd {
	return f.Inner.Focus()
}

func (f *BorderFrame) Blur() {
	f.Inner.Blur()
}

func (f *BorderFrame) Focused() bool {
	return f.Inner.Focused()
}

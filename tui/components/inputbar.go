package components

import (
	"terminal_chat/tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputBar struct {
	Model textinput.Model
}

func NewInputBar(placeholder string) *InputBar {
	ti := styles.NewTextInput(false)
	ti.Placeholder = placeholder
	return &InputBar{Model: ti}
}

func (b *InputBar) Render(width int) string {
	b.Model.Focus()
	b.Model.Width = width - 7

	return styles.InputWrapper.Render(
		styles.InputPrompt.Render("> ") + b.Model.View(),
	)
}

func (b *InputBar) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	b.Model, cmd = b.Model.Update(msg)
	return cmd
}

func (b *InputBar) Value() string {
	return b.Model.Value()
}

func (b *InputBar) SetValue(v string) {
	b.Model.SetValue(v)
}

func (b *InputBar) OnKeyPress(key string) (int, tea.Cmd) {
	return 0, nil
}

func (b *InputBar) Focus() tea.Cmd {
	b.Model.Focus()
	return textinput.Blink
}

func (b *InputBar) Blur() {
	b.Model.Blur()
}

func (b *InputBar) Focused() bool {
	return b.Model.Focused()
}

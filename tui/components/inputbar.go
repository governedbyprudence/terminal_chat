package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var textInputStyle = lipgloss.NewStyle().Foreground(Text)

var inputWrapper = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true, false, true, false).
	BorderForeground(Accent).
	Padding(0, 2)

var inputPrompt = lipgloss.NewStyle().
	Foreground(Accent).
	Bold(true)

func newTextInput() textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Focus()
	ti.TextStyle = textInputStyle
	return ti
}

// textInputWidthFix compensates for bubbles textinput Width+1 bug
// plus prompt (2 chars) + wrapper border/padding (4 chars) = 7
const textInputWidthFix = 7

type InputBar struct {
	Model textinput.Model
}

func NewInputBar(placeholder string) *InputBar {
	ti := newTextInput()
	ti.Placeholder = placeholder
	return &InputBar{Model: ti}
}

func (b *InputBar) Render(width int) string {
	b.Model.Focus()
	b.Model.Width = width - textInputWidthFix

	return inputWrapper.Render(
		inputPrompt.Render("> ") + b.Model.View(),
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

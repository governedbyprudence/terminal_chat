package components

import (
	"strings"

	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"
	"terminal_chat/tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MessageList struct {
	Messages []states.Message
}

func NewMessageList() *MessageList {
	return &MessageList{}
}

func (m *MessageList) AddMessage(msg states.Message) {
	m.Messages = append(m.Messages, msg)
}

func (m *MessageList) Render(width int) string {
	if len(m.Messages) == 0 {
		return styles.EmptyState.Render("  no messages yet")
	}

	rightPad := func(s string) string {
		w := lipgloss.Width(s)
		if w >= width {
			return s
		}
		return strings.Repeat(" ", width-w) + s
	}

	var b strings.Builder
	for _, msg := range m.Messages {
		shortID := msg.UserID
		if len(shortID) > 6 {
			shortID = shortID[:6]
		}

		if msg.UserID == utils.CurrentUserID() {
			meta := styles.MessageTimestamp.Render("you")
			line := styles.SentMessage.Render(msg.Text)
			b.WriteString(meta)
			b.WriteString("\n")
			b.WriteString(line)
			b.WriteString("\n")
		} else {
			line := rightPad(styles.ReceivedMessage.Render(msg.Text))
			sender := rightPad(styles.SenderLabel.Render(shortID))
			b.WriteString(line)
			b.WriteString("\n")
			b.WriteString(sender)
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m *MessageList) OnKeyPress(key string) (int, tea.Cmd) {
	return 0, nil
}

func (m *MessageList) Focus() tea.Cmd  { return nil }
func (m *MessageList) Blur()           {}
func (m *MessageList) Focused() bool   { return false }

package components

import (
	"strings"

	"terminal_chat/tui/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var emptyState = lipgloss.NewStyle().
	Foreground(DimText).
	Italic(true).
	Padding(1, 2)

var messageTimestamp = lipgloss.NewStyle().
	Foreground(DimText).
	Italic(true)

var sentMessage = lipgloss.NewStyle().
	Foreground(Text).
	Bold(true)

var receivedMessage = lipgloss.NewStyle().
	Foreground(Text)

var senderLabel = lipgloss.NewStyle().
	Foreground(DimText)

type MessageList struct {
	Messages []models.Message
}

func NewMessageList() *MessageList {
	return &MessageList{}
}

func (m *MessageList) AddMessage(msg models.Message) {
	m.Messages = append(m.Messages, msg)
}

func rightPad(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return strings.Repeat(" ", width-w) + s
}

func fitWidth(s string, width int) string {
	if lipgloss.Width(s) <= width {
		return s
	}
	var b strings.Builder
	var w int
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if w+rw > width {
			break
		}
		b.WriteRune(r)
		w += rw
	}
	return b.String()
}

func (m *MessageList) Render(width int) string {
	if len(m.Messages) == 0 {
		return emptyState.Render("  no messages yet")
	}

	var b strings.Builder
	for _, msg := range m.Messages {
		shortID := msg.UserID
		if len(shortID) > 6 {
			shortID = shortID[:6]
		}

		text := fitWidth(msg.Text, width)
		if msg.UserID == models.CurrentUserID() {
			meta := messageTimestamp.Render("you")
			line := sentMessage.Render(text)
			b.WriteString(meta)
			b.WriteString("\n")
			b.WriteString(line)
			b.WriteString("\n")
		} else {
			line := rightPad(receivedMessage.Render(text), width)
			sender := rightPad(senderLabel.Render(shortID), width)
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

func (m *MessageList) Update(msg tea.Msg) tea.Cmd { return nil }
func (m *MessageList) Focus() tea.Cmd             { return nil }
func (m *MessageList) Blur()                      {}
func (m *MessageList) Focused() bool              { return false }

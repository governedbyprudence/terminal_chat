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
	Messages     []models.Message
	scrollOffset int
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

func (m *MessageList) ScrollTrack(maxHeight int) []string {
	track := make([]string, maxHeight)
	for i := range track {
		track[i] = "░░"
	}

	linesPerMsg := 2
	maxMsgs := maxHeight / linesPerMsg
	if maxMsgs <= 0 {
		maxMsgs = 1
	}
	maxOffset := max(0, len(m.Messages)-maxMsgs)
	if maxOffset <= 0 {
		return track
	}

	thumbPos := maxHeight - 1 - int(float64(m.scrollOffset)/float64(maxOffset)*float64(maxHeight-1))
	if thumbPos < 0 {
		thumbPos = 0
	}
	if thumbPos >= maxHeight {
		thumbPos = maxHeight - 1
	}
	track[thumbPos] = "██"
	return track
}

func (m *MessageList) Render(width, maxHeight int) string {
	if len(m.Messages) == 0 {
		m.scrollOffset = 0
		empty := emptyState.Render("  no messages yet")
		emptyLines := strings.Count(empty, "\n") + 1
		padLines := maxHeight - emptyLines
		if padLines > 0 {
			return strings.Repeat("\n", padLines) + empty
		}
		return empty
	}

	linesPerMsg := 2
	maxMsgs := maxHeight / linesPerMsg
	if maxMsgs <= 0 {
		maxMsgs = 1
	}

	maxOffset := max(0, len(m.Messages)-maxMsgs)
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}

	start := maxOffset - m.scrollOffset
	if start < 0 {
		start = 0
	}
	end := start + maxMsgs
	if end > len(m.Messages) {
		end = len(m.Messages)
	}

	var msgBuf strings.Builder
	for _, msg := range m.Messages[start:end] {
		shortID := msg.UserID
		if len(shortID) > 6 {
			shortID = shortID[:6]
		}

		text := fitWidth(msg.Text, width)
		if msg.UserID == models.CurrentUserID() {
			meta := messageTimestamp.Render("you")
			line := sentMessage.Render(text)
			msgBuf.WriteString(meta)
			msgBuf.WriteString("\n")
			msgBuf.WriteString(line)
			msgBuf.WriteString("\n")
		} else {
			line := rightPad(receivedMessage.Render(text), width)
			sender := rightPad(senderLabel.Render(shortID), width)
			msgBuf.WriteString(line)
			msgBuf.WriteString("\n")
			msgBuf.WriteString(sender)
			msgBuf.WriteString("\n")
		}
	}

	rendered := msgBuf.String()
	renderedLines := strings.Count(rendered, "\n")
	padLines := maxHeight - renderedLines
	if padLines > 0 {
		rendered = strings.Repeat("\n", padLines) + rendered
	}
	return strings.TrimRight(rendered, "\n")
}

func (m *MessageList) ScrollUp(n int) {
	m.scrollOffset += n
}

func (m *MessageList) ScrollDown(n int) {
	m.scrollOffset -= n
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

func (m *MessageList) OnKeyPress(key string) (int, tea.Cmd) {
	return 0, nil
}

func (m *MessageList) Update(msg tea.Msg) tea.Cmd { return nil }
func (m *MessageList) Focus() tea.Cmd             { return nil }
func (m *MessageList) Blur()                      {}
func (m *MessageList) Focused() bool              { return false }

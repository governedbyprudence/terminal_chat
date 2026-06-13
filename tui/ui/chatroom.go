package ui

import (
	"encoding/json"
	"log"
	"strings"

	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"
	"terminal_chat/tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatRoomScreen struct {
	State *states.ChatState
}

func (s ChatRoomScreen) Render(width, height int) string {
	contentW := width - styles.OuterFrameHSize

	title := styles.TitleBar.Render("  " + s.State.RoomName + "  ")

	messages := ""
	for _, msg := range s.State.Messages {
		shortID := msg.UserID
		if len(shortID) > 6 {
			shortID = shortID[:6]
		}

		if msg.UserID == utils.CurrentUserID() {
			line := styles.SentMessage.Render(msg.Text)
			meta := styles.MessageTimestamp.Render("you")
			messages += meta + "\n" + line + "\n"
		} else {
			rightPad := func(s string) string {
				w := lipgloss.Width(s)
				if w >= contentW {
					return s
				}
				return strings.Repeat(" ", contentW-w) + s
			}
			messages += rightPad(styles.ReceivedMessage.Render(msg.Text)) + "\n"
			messages += rightPad(styles.SenderLabel.Render(shortID)) + "\n"
		}
	}

	if messages == "" {
		messages = styles.EmptyState.Render("  no messages yet")
	}

	s.State.TextInput.Focus()
	s.State.TextInput.Width = contentW - 7

	input := styles.InputWrapper.Render(
		styles.InputPrompt.Render("> ") + s.State.TextInput.View(),
	)

	footer := styles.FooterText.Render(
		"ENTER send  ·  ESC back to rooms",
	)

	body := title + "\n" + messages + "\n" + input + "\n" + footer

	return styles.OuterFrame.Render(body)
}

func (s ChatRoomScreen) OnPressKey(key string) (int, tea.Cmd) {
	switch key {
	case "esc":
		return 1, nil
	case "enter":
		if s.State.TextInput.Value() != "" {
			msg := states.Message{
				RoomName: s.State.RoomName,
				Text:     s.State.TextInput.Value(),
				UserID:   utils.CurrentUserID(),
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				return 0, nil
			}

			messageStore.DefaultBroker().Publish("chat", string(data))
			s.State.TextInput.SetValue("")
			return 0, nil
		}
	}
	return 0, nil
}

func (s *ChatRoomScreen) UpdateMessage(msg states.Message) {
	s.State.Messages = append(s.State.Messages, msg)
}

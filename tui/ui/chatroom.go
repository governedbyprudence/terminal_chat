package ui

import (
	"encoding/json"
	"log"

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

func (s ChatRoomScreen) Render() string {
	title := styles.TitleStyle.Render("Chat Room: " + s.State.RoomName)
	body := ""

	for _, msg := range s.State.Messages {
		if msg.UserID == utils.CurrentUserID() {
			m := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render("You: " + msg.Text)
			body += lipgloss.NewStyle().Align(lipgloss.Right).Render(m) + "\n"
		} else {
			m := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render(msg.UserID[:6] + ": " + msg.Text)
			body += lipgloss.NewStyle().Align(lipgloss.Left).Render(m) + "\n"
		}
	}

	body = lipgloss.JoinVertical(lipgloss.Top, body)

	s.State.TextInput.Focus()
	body += "\n> " + s.State.TextInput.View()

	footer := styles.FooterStyle.Render("Type message and press `enter` to send.\n\nPress `esc` to go back to room list.")
	return lipgloss.JoinVertical(lipgloss.Top, title, body, footer)
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

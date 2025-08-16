package ui

import (
	"context"
	"encoding/json"
	"log"
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

	for _, message := range s.State.Messages {
		// Check if the message is from the current user
		if message.UserID == s.CurrentUserID() {
			// Right-align and style user's own messages differently
			msg := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render("You: " + message.Message)
			body += lipgloss.NewStyle().Align(lipgloss.Right).Render(msg) + "\n"
		} else {
			// Left-align other users' messages
			msg := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render(message.UserID[:6] + ": " + message.Message)
			body += lipgloss.NewStyle().Align(lipgloss.Left).Render(msg) + "\n"
		}
	}

	body = lipgloss.JoinVertical(lipgloss.Top, body)

	// Text input field
	s.State.TextInput.Focus()
	body += "\n> " + s.State.TextInput.View()

	footer := styles.FooterStyle.Render("Type message and press `enter` to send.\n\nPress `esc` to go back to room list.")
	return lipgloss.JoinVertical(lipgloss.Top, title, body, footer)

}

// CurrentUserID returns the current user's ID
func (s ChatRoomScreen) CurrentUserID() string {
	return utils.CurrentUserID()
}

func (s ChatRoomScreen) OnPressKey(key string) (int, tea.Cmd) {
	switch key {
	case "esc":
		return 1, nil
	case "enter":
		if s.State.TextInput.Value() != "" {
			// Create a message with user ID
			message := struct {
				RoomName string `json:"room_name"`
				Message  string `json:"message"`
				UserID   string `json:"user_id"`
			}{
				RoomName: s.State.RoomName,
				Message:  s.State.TextInput.Value(),
				UserID:   s.CurrentUserID(),
			}

			// Send the message to the server
			jsonEncodedMessage, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return 0, nil
			}

			ctx := context.Background()
			s.State.Redis.Publish(ctx, "chat", jsonEncodedMessage)

			// Clear the text input
			s.State.TextInput.SetValue("")
			return 0, nil
		}
	}
	return 0, nil
}

func (s *ChatRoomScreen) UpdateMessage(msg string, userID string) {
	s.State.Messages = append(s.State.Messages, states.Messages{
		Message:  msg,
		RoomName: s.State.RoomName,
		UserID:   userID,
	})
}

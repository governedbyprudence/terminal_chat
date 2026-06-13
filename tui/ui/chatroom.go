package ui

import (
	"encoding/json"
	"log"

	"terminal_chat/tui/components"
	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"
	"terminal_chat/tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatRoomScreen struct {
	RoomName    string
	MessageList *components.MessageList
	InputBar    *components.InputBar
	TitleBar    *components.TitleBar
	Footer      *components.Footer
	backRef     states.AppScreen
}

func NewChatRoomScreen(backRef states.AppScreen, roomName string, messages ...states.Message) *ChatRoomScreen {
	s := &ChatRoomScreen{
		RoomName:    roomName,
		backRef:     backRef,
		MessageList: components.NewMessageList(),
		InputBar:    components.NewInputBar(""),
		TitleBar:    components.NewTitleBar(roomName),
		Footer:      components.NewFooter("ENTER send  ·  ESC back to rooms"),
	}
	for _, msg := range messages {
		s.MessageList.AddMessage(msg)
	}
	return s
}

func (s *ChatRoomScreen) Render(width, _ int) string {
	contentW := width - styles.OuterFrameHSize

	title := s.TitleBar.Render(contentW)
	messages := s.MessageList.Render(contentW)
	input := s.InputBar.Render(contentW)
	footer := s.Footer.Render(contentW)

	body := title + "\n" + messages + "\n" + input + "\n" + footer

	wrapper := lipgloss.NewStyle().Width(contentW)
	return styles.OuterFrame.Render(wrapper.Render(body)) + "\n"
}

func (s *ChatRoomScreen) Update(msg tea.Msg) (states.AppScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s.backRef, nil
		case "enter":
			if s.InputBar.Value() != "" {
				m := states.Message{
					RoomName: s.RoomName,
					Text:     s.InputBar.Value(),
					UserID:   utils.CurrentUserID(),
				}
				data, err := json.Marshal(m)
				if err != nil {
					log.Println(err)
					return s, nil
				}
				messageStore.DefaultBroker().Publish("chat", string(data))
				s.InputBar.SetValue("")
			}
			return s, nil
		default:
			return s, s.InputBar.Update(msg)
		}
	case states.Message:
		if msg.RoomName == s.RoomName {
			s.MessageList.AddMessage(msg)
		}
		return s, nil
	}
	return s, nil
}

package ui

import (
	"encoding/json"
	"log"
	"strings"

	"terminal_chat/tui/components"
	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatRoomScreen struct {
	RoomName    string
	MessageList *components.MessageList
	InputBar    *components.InputBar
	TitleBar    *components.TitleBar
	Footer      *components.Footer
	UserPanel   *components.UserPanel
	backRef     models.AppScreen
	broker      messageStore.Broker
}

func NewChatRoomScreen(broker messageStore.Broker, backRef models.AppScreen, roomName string, messages ...models.Message) *ChatRoomScreen {
	s := &ChatRoomScreen{
		RoomName:    roomName,
		backRef:     backRef,
		broker:      broker,
		UserPanel:   components.NewUserPanel(),
		MessageList: components.NewMessageList(),
		InputBar:    components.NewInputBar(""),
		TitleBar:    components.NewTitleBar(roomName),
		Footer:      components.NewFooter("ENTER send  ·  ESC back to rooms"),
	}
	s.UserPanel.UserList.AddUser(models.CurrentUserID())
	s.publishPresence("join")
	for _, msg := range messages {
		s.MessageList.AddMessage(msg)
	}
	return s
}

func (s *ChatRoomScreen) publishPresence(action string) {
	evt := models.PresenceEvent{
		UserID:   models.CurrentUserID(),
		RoomName: s.RoomName,
		Action:   action,
	}
	data, err := json.Marshal(evt)
	if err != nil {
		log.Printf("presence marshal error: %v", err)
		return
	}
	s.broker.Publish(messageStore.ChannelPresence, string(data))
}

func (s *ChatRoomScreen) Render(width, height int) string {
	contentW := width - outerFrameHSize()
	contentH := height - outerFrameVSize()
	leftW := contentW - components.RightPanelWidth - components.PanelGap

	title := s.TitleBar.Render(leftW)

	msgHeight := contentH - 5
	if msgHeight < 2 {
		msgHeight = 2
	}
	messages := s.MessageList.Render(leftW, msgHeight)

	input := s.InputBar.Render(leftW)
	footer := s.Footer.Render(leftW)

	leftBody := title + "\n" + messages + "\n" + input + "\n" + footer
	leftPadded := lipgloss.NewStyle().Width(leftW).Render(leftBody)

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPadded, s.UserPanel.RenderFull(contentH))

	bodyLines := strings.Count(body, "\n") + 1
	if bodyLines > contentH {
		body = strings.Join(strings.Split(body, "\n")[bodyLines-contentH:], "\n")
	}

	wrapper := lipgloss.NewStyle().Width(contentW)
	return outerFrame.Render(wrapper.Render(body))
}

func (s *ChatRoomScreen) Update(msg tea.Msg) (models.AppScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			s.publishPresence("leave")
			return s.backRef, nil
		case "enter":
			text := strings.TrimSpace(s.InputBar.Value())
			if text != "" {
				m := models.Message{
					RoomName: s.RoomName,
					Text:     text,
					UserID:   models.CurrentUserID(),
				}
				data, err := json.Marshal(m)
				if err != nil {
					log.Printf("message marshal error: %v", err)
					return s, nil
				}
				s.broker.Publish(messageStore.ChannelChat, string(data))
				s.InputBar.SetValue("")
			}
			return s, nil
		case "pgup":
			s.MessageList.ScrollUp(5)
			return s, nil
		case "pgdown":
			s.MessageList.ScrollDown(5)
			return s, nil
		default:
			return s, s.InputBar.Update(msg)
		}
	case models.Message:
		if msg.RoomName == s.RoomName {
			s.UserPanel.UserList.AddUser(msg.UserID)
			s.MessageList.AddMessage(msg)
		}
		return s, nil
	case models.PresenceEvent:
		if msg.RoomName == s.RoomName {
			switch msg.Action {
			case "join":
				s.UserPanel.UserList.AddUser(msg.UserID)
			case "leave":
				s.UserPanel.UserList.RemoveUser(msg.UserID)
			}
		}
		return s, nil
	}
	return s, nil
}

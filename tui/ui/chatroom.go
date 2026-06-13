package ui

import (
	"encoding/json"
	"log"
	"strings"

	"terminal_chat/tui/components"
	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"
	"terminal_chat/tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const rightPanelWidth = 22

type ChatRoomScreen struct {
	RoomName    string
	UserList    *components.UserList
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
		UserList:    components.NewUserList(),
		MessageList: components.NewMessageList(),
		InputBar:    components.NewInputBar(""),
		TitleBar:    components.NewTitleBar(roomName),
		Footer:      components.NewFooter("ENTER send  ·  ESC back to rooms"),
	}
	s.UserList.AddUser(utils.CurrentUserID())
	s.publishPresence("join")
	for _, msg := range messages {
		s.MessageList.AddMessage(msg)
	}
	return s
}

func (s *ChatRoomScreen) publishPresence(action string) {
	evt := states.PresenceEvent{
		UserID:   utils.CurrentUserID(),
		RoomName: s.RoomName,
		Action:   action,
	}
	data, err := json.Marshal(evt)
	if err != nil {
		log.Println(err)
		return
	}
	messageStore.DefaultBroker().Publish("presence", string(data))
}

func (s *ChatRoomScreen) Render(width, _ int) string {
	contentW := width - styles.OuterFrameHSize
	leftW := contentW - rightPanelWidth - 2

	title := s.TitleBar.Render(leftW)
	messages := s.MessageList.Render(leftW)
	input := s.InputBar.Render(leftW)
	footer := s.Footer.Render(leftW)

	leftBody := title + "\n" + messages + "\n" + input + footer
	leftPadded := lipgloss.NewStyle().Width(leftW).Render(leftBody)

	leftLines := strings.Count(leftPadded, "\n")
	if leftPadded != "" && !strings.HasSuffix(leftPadded, "\n") {
		leftLines++
	}

	var gapBuilder strings.Builder
	for i := 0; i < leftLines-1; i++ {
		gapBuilder.WriteString("  \n")
	}
	gapBuilder.WriteString("  ")
	gap := gapBuilder.String()

	var rightBuilder strings.Builder
	rightBuilder.WriteString(styles.UserPanelTitle.Render("ONLINE"))
	rightBuilder.WriteString("\n")
	if s.UserList.Count() == 0 {
		rightBuilder.WriteString("nobody here")
		rightBuilder.WriteString("\n")
	} else {
		rightBuilder.WriteString(s.UserList.Render(rightPanelWidth - 5))
	}

	rightText := strings.TrimSuffix(rightBuilder.String(), "\n")
	rightLines := strings.Count(rightText, "\n") + 1
	for i := rightLines; i < leftLines; i++ {
		rightText += "\n"
	}

	rightPadded := styles.UserPanel.Width(rightPanelWidth-3).Render(rightText)

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPadded, gap, rightPadded)

	wrapper := lipgloss.NewStyle().Width(contentW)
	return styles.OuterFrame.Render(wrapper.Render(body)) + "\n"
}

func (s *ChatRoomScreen) Update(msg tea.Msg) (states.AppScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			s.publishPresence("leave")
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
			s.UserList.AddUser(msg.UserID)
			s.MessageList.AddMessage(msg)
		}
		return s, nil
	case states.PresenceEvent:
		if msg.RoomName == s.RoomName {
			switch msg.Action {
			case "join":
				s.UserList.AddUser(msg.UserID)
			case "leave":
				s.UserList.RemoveUser(msg.UserID)
			}
		}
		return s, nil
	}
	return s, nil
}

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
	Rooms         []string
	ActiveRoomIdx int
	InputBar      *components.InputBar
	Footer        *components.Footer
	UserPanel     *components.UserPanel
	RoomTabs      *components.RoomTabs
	messageLists  map[string]*components.MessageList
	usersPerRoom  map[string]map[string]bool
	backRef       models.AppScreen
	broker        messageStore.Broker
}

func NewChatRoomScreen(broker messageStore.Broker, backRef models.AppScreen, rooms []string, activeIdx int, messages ...models.Message) *ChatRoomScreen {
	if activeIdx < 0 || activeIdx >= len(rooms) {
		activeIdx = 0
	}

	ml := components.NewMessageList()
	for _, msg := range messages {
		ml.AddMessage(msg)
	}

	activeRoom := rooms[activeIdx]
	messageLists := map[string]*components.MessageList{activeRoom: ml}
	usersPerRoom := map[string]map[string]bool{activeRoom: {models.CurrentUserID(): true}}

	s := &ChatRoomScreen{
		Rooms:         rooms,
		ActiveRoomIdx: activeIdx,
		backRef:       backRef,
		broker:        broker,
		UserPanel:     components.NewUserPanel(),
		InputBar:      components.NewInputBar(""),
		Footer:        components.NewFooter("ENTER send  ·  TAB/←→ switch room  ·  ESC back"),
		RoomTabs:      components.NewRoomTabs(rooms),
		messageLists:  messageLists,
		usersPerRoom:  usersPerRoom,
	}
	s.RoomTabs.ActiveIndex = activeIdx
	s.syncUserPanel()
	s.publishPresence("join")
	return s
}

func (s *ChatRoomScreen) publishPresence(action string) {
	evt := models.PresenceEvent{
		UserID:   models.CurrentUserID(),
		RoomName: s.activeRoom(),
		Action:   action,
	}
	data, err := json.Marshal(evt)
	if err != nil {
		log.Printf("presence marshal error: %v", err)
		return
	}
	s.broker.Publish(messageStore.ChannelPresence, string(data))
}

func (s *ChatRoomScreen) activeRoom() string {
	if s.ActiveRoomIdx < 0 || s.ActiveRoomIdx >= len(s.Rooms) {
		return ""
	}
	return s.Rooms[s.ActiveRoomIdx]
}

func (s *ChatRoomScreen) activeMessageList() *components.MessageList {
	room := s.activeRoom()
	if room == "" {
		return nil
	}
	ml, ok := s.messageLists[room]
	if !ok {
		ml = components.NewMessageList()
		s.messageLists[room] = ml
	}
	return ml
}

func (s *ChatRoomScreen) syncUserPanel() {
	room := s.activeRoom()
	s.UserPanel.UserList = components.NewUserList()
	users, ok := s.usersPerRoom[room]
	if ok {
		for uid := range users {
			s.UserPanel.UserList.AddUser(uid)
		}
	} else {
		s.usersPerRoom[room] = map[string]bool{models.CurrentUserID(): true}
		s.UserPanel.UserList.AddUser(models.CurrentUserID())
	}
}

func (s *ChatRoomScreen) switchRoom(idx int) {
	if idx < 0 || idx >= len(s.Rooms) || idx == s.ActiveRoomIdx {
		return
	}

	s.publishPresence("leave")

	s.ActiveRoomIdx = idx
	s.RoomTabs.ActiveIndex = idx
	s.syncUserPanel()

	s.publishPresence("join")
}

func (s *ChatRoomScreen) Render(width, height int) string {
	contentW := width - outerFrameHSize()
	contentH := height - outerFrameVSize()
	leftW := contentW - components.RightPanelWidth - components.PanelGap

	tabs := s.RoomTabs.Render(contentW)
	sep := sepLine(contentW)

	belowSepH := contentH - 2
	msgHeight := belowSepH - 4
	if msgHeight < 2 {
		msgHeight = 2
	}
	ml := s.activeMessageList()
	var messages string
	if ml != nil {
		messages = ml.Render(leftW, msgHeight)
	}

	input := s.InputBar.Render(leftW)
	footer := s.Footer.Render(leftW)

	leftBody := messages + "\n" + input + "\n" + footer
	leftPadded := lipgloss.NewStyle().Width(leftW).Render(leftBody)

	combined := lipgloss.JoinHorizontal(lipgloss.Top, leftPadded, s.UserPanel.RenderFull(belowSepH))

	combSplit := strings.Split(combined, "\n")
	if len(combSplit) > belowSepH {
		combSplit = combSplit[len(combSplit)-belowSepH:]
	}
	for i, line := range combSplit {
		if w := lipgloss.Width(line); w < contentW {
			combSplit[i] = line + strings.Repeat(" ", contentW-w)
		}
	}
	combined = strings.Join(combSplit, "\n")

	body := tabs + "\n" + sep + "\n" + combined

	bodyLines := strings.Count(body, "\n") + 1
	if bodyLines > contentH {
		body = strings.Join(strings.Split(body, "\n")[bodyLines-contentH:], "\n")
	}

	return outerFrame.Render(body)
}

func (s *ChatRoomScreen) Update(msg tea.Msg) (models.AppScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if s.backRef != nil {
				s.publishPresence("leave")
				return s.backRef, nil
			}
			return s, nil
		case "enter":
			text := strings.TrimSpace(s.InputBar.Value())
			if text != "" {
				m := models.Message{
					RoomName: s.activeRoom(),
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
			if ml := s.activeMessageList(); ml != nil {
				ml.ScrollUp(5)
			}
			return s, nil
		case "pgdown":
			if ml := s.activeMessageList(); ml != nil {
				ml.ScrollDown(5)
			}
			return s, nil
		case "left":
			s.switchRoom(s.ActiveRoomIdx - 1)
			return s, nil
		case "right":
			s.switchRoom(s.ActiveRoomIdx + 1)
			return s, nil
		case "tab":
			s.switchRoom((s.ActiveRoomIdx + 1) % len(s.Rooms))
			return s, nil
		case "shift+tab":
			idx := s.ActiveRoomIdx - 1
			if idx < 0 {
				idx = len(s.Rooms) - 1
			}
			s.switchRoom(idx)
			return s, nil
		default:
			return s, s.InputBar.Update(msg)
		}
	case models.Message:
		if ml, ok := s.messageLists[msg.RoomName]; ok {
			ml.AddMessage(msg)
		} else {
			ml := components.NewMessageList()
			ml.AddMessage(msg)
			s.messageLists[msg.RoomName] = ml
		}
		users, ok := s.usersPerRoom[msg.RoomName]
		if !ok {
			users = map[string]bool{}
			s.usersPerRoom[msg.RoomName] = users
		}
		users[msg.UserID] = true
		if msg.RoomName == s.activeRoom() {
			s.syncUserPanel()
		}
		return s, nil
	case models.PresenceEvent:
		users, ok := s.usersPerRoom[msg.RoomName]
		if !ok {
			users = map[string]bool{}
			s.usersPerRoom[msg.RoomName] = users
		}
		switch msg.Action {
		case "join":
			users[msg.UserID] = true
		case "leave":
			delete(users, msg.UserID)
		}
		if msg.RoomName == s.activeRoom() {
			s.syncUserPanel()
		}
		return s, nil
	}
	return s, nil
}

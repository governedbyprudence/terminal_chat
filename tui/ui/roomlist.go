package ui

import (
	"strings"

	"terminal_chat/tui/components"
	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type RoomListScreen struct {
	TitleBar *components.TitleBar
	RoomList *components.RoomList
	Footer   *components.Footer
	broker   messageStore.Broker
}

func NewRoomListScreen(broker messageStore.Broker, rooms []string) *RoomListScreen {
	return &RoomListScreen{
		broker:   broker,
		TitleBar: components.NewTitleBar("ROOMS"),
		RoomList: components.NewRoomList(rooms),
		Footer:   components.NewFooter("UP/DOWN navigate  ·  ENTER select  ·  q quit"),
	}
}

func (s *RoomListScreen) Render(width int) string {
	contentW := width - outerFrameHSize()

	title := s.TitleBar.Render(contentW)
	items := s.RoomList.Render(contentW)
	sep := strings.Repeat("─", contentW-4)
	footer := s.Footer.Render(contentW)

	body := title + "\n" + items + "\n" + sep + "\n" + footer

	wrapper := lipgloss.NewStyle().Width(contentW)
	return outerFrame.Render(wrapper.Render(body)) + "\n"
}

func (s *RoomListScreen) Update(msg tea.Msg) (models.AppScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		redirect, cmd := s.RoomList.OnKeyPress(msg.String())
		if cmd != nil {
			return s, cmd
		}
		if redirect == components.RedirectSelect {
			roomName := s.RoomList.Rooms[s.RoomList.SelectedRoom]
			chat := NewChatRoomScreen(s.broker, s, roomName)
			return chat, textinput.Blink
		}
	}
	return s, nil
}

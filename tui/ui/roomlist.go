package ui

import (
	"strings"

	"terminal_chat/tui/components"
	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type RoomListScreen struct {
	TitleBar *components.TitleBar
	RoomList *components.RoomList
	Footer   *components.Footer
}

func NewRoomListScreen(rooms []string) *RoomListScreen {
	return &RoomListScreen{
		TitleBar: components.NewTitleBar("ROOMS"),
		RoomList: components.NewRoomList(rooms),
		Footer:   components.NewFooter("UP/DOWN navigate  ·  ENTER select  ·  q quit"),
	}
}

func (s *RoomListScreen) Render(width, _ int) string {
	contentW := width - styles.OuterFrameHSize

	title := s.TitleBar.Render(contentW)
	items := s.RoomList.Render(contentW)
	sep := strings.Repeat("─", contentW-4)
	footer := s.Footer.Render(contentW)

	body := title + "\n" + items + "\n" + sep + "\n\n" + footer

	wrapper := lipgloss.NewStyle().Width(contentW)
	return styles.OuterFrame.Render(wrapper.Render(body))
}

func (s *RoomListScreen) Update(msg tea.Msg) (states.AppScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		redirect, cmd := s.RoomList.OnKeyPress(msg.String())
		if cmd != nil {
			return s, cmd
		}
		if redirect == 1 {
			roomName := s.RoomList.Rooms[s.RoomList.SelectedRoom]
			chat := NewChatRoomScreen(s, roomName)
			return chat, textinput.Blink
		}
	}
	return s, nil
}

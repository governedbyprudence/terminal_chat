package ui

import (
	"strings"

	"terminal_chat/tui/components"
	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
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
Footer: components.NewFooter("UP/DOWN navigate  ·  ENTER select  ·  qq quit"),
	}
}

func (s *RoomListScreen) Render(width, height int) string {
	contentW := width - outerFrameHSize()
	contentH := height - outerFrameVSize()

	title := s.TitleBar.Render(contentW)
	items := s.RoomList.Render(contentW)
	sep := sepLine(contentW)
	footer := s.Footer.Render(contentW)

	body := title + "\n" + items + "\n" + sep + "\n" + footer

	bodyLines := strings.Count(body, "\n") + 1
	if bodyLines < contentH {
		body += strings.Repeat("\n", contentH-bodyLines)
	} else if bodyLines > contentH {
		body = strings.Join(strings.Split(body, "\n")[bodyLines-contentH:], "\n")
	}

	return outerFrame.Render(body)
}

func (s *RoomListScreen) Update(msg tea.Msg) (models.AppScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		redirect, cmd := s.RoomList.OnKeyPress(msg.String())
		if cmd != nil {
			return s, cmd
		}
		if redirect == components.RedirectSelect {
			chat := NewChatRoomScreen(s.broker, s, s.RoomList.Rooms, s.RoomList.SelectedRoom)
			return chat, textinput.Blink
		}
	}
	return s, nil
}

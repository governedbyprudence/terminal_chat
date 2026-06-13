package ui

import (
	"strings"

	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type RoomListScreen struct {
	State *states.RoomListState
}

func (s RoomListScreen) Render(width, height int) string {
	contentW := width - styles.OuterFrameHSize

	frame := styles.OuterFrame.Width(contentW)

	title := styles.TitleBar.Render("  ROOMS  ")

	items := ""
	for i, room := range s.State.Rooms {
		if i == s.State.SelectedRoom {
			items += styles.SelectedMenuItem.Render("  > " + room) + "\n"
		} else {
			items += styles.MenuItem.Render("    " + room) + "\n"
		}
	}

	sep := strings.Repeat("─", contentW-4)

	footer := styles.FooterText.Render(
		"UP/DOWN navigate  ·  ENTER select  ·  q quit",
	)

	body := title + "\n" + items + "\n" + sep + "\n\n" + footer

	return frame.Render(body)
}

func (s RoomListScreen) OnPressKey(key string) (int, tea.Cmd) {
	switch key {
	case "down":
		if s.State.SelectedRoom < len(s.State.Rooms)-1 {
			s.State.SelectedRoom++
		}
	case "up":
		if s.State.SelectedRoom > 0 {
			s.State.SelectedRoom--
		}
	case "enter":
		return 1, nil
	}
	return 0, nil
}
